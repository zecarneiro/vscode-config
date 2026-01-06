package processors

import (
	"fmt"
	"golangutils/pkg/console"
	"golangutils/pkg/exe"
	"golangutils/pkg/file"
	"golangutils/pkg/logger"
	"golangutils/pkg/logic"
	"golangutils/pkg/obj"
	"golangutils/pkg/platform"
	"golangutils/pkg/system"
	"main/entities"
	"main/libs"
	"slices"
	"strconv"
	"strings"
)

type InstallProcessor struct {
	profileProcessor      *ProfileProcessor
	downloadProcessor     *DownloadProcessor
	profileAlreadyCreated []string
}

func newInstallProcessor(profileProcessor *ProfileProcessor) *InstallProcessor {
	return &InstallProcessor{
		profileProcessor:      profileProcessor,
		downloadProcessor:     &DownloadProcessor{},
		profileAlreadyCreated: []string{},
	}
}

/* --------------------------------- PRIVATE -------------------------------- */
const (
	MAX_INSTALL_EXTENSIONS = 5
)

func (ip *InstallProcessor) installJsonFile(jsonFile string) {
	jsonFile = file.ResolvePath(jsonFile)
	if file.IsFile(jsonFile) {
		logger.Info("Copy given file to: " + libs.JsonFile)
		logic.ProcessError(file.CopyFile(jsonFile, libs.JsonFile))
	}
}

func (ip *InstallProcessor) setSettingConfigurations() {
	var settingsDir string
	logger.Separator()
	logger.Header("Set settings")
	const keyAllSettings = "workbench.settings.applyToAllProfiles"
	arrayAllSettings := []string{}
	settings := ip.profileProcessor.getAllInstallSettings()

	if platform.IsWindows() {
		settingsDir = file.ResolvePath(system.HomeDir(), "AppData\\Roaming\\Code\\User")
	} else if platform.IsLinux() {
		settingsDir = file.ResolvePath(system.HomeDir(), ".config/Code/User")
	}
	for key := range settings {
		arrayAllSettings = append(arrayAllSettings, key)
	}
	// settings[keyAllSettings] = arrayAllSettings # Problems with duplicated settings with local settings and dev containers settings
	if len(settingsDir) > 0 {
		file.CreateDirectory(settingsDir, true)
		logic.ProcessError(file.WriteJsonFile(file.ResolvePath(settingsDir, "settings.json"), settings, true))
	} else {
		logger.Debug("\n\nGo to setting and open json settings")
		logger.Debug("Append this setting bellow on json settings")
		settingsStr, err := obj.ObjectToString(settings)
		logic.ProcessError(err)
		fmt.Println(settingsStr)
		libs.OpenVscodeWithNewProfile(libs.Configurations.SettingsName)
	}
}

func (ip *InstallProcessor) processInstall() {
	skipProfiles := []string{}
	console.WaitForAnyKeyPressed("Please, close all instance of VSCode and PRESS ANY KEY TO CONTINUE...")
	for _, profile := range ip.profileProcessor.profiles {
		if !profile.CanInstall {
			skipProfiles = append(skipProfiles, profile.Name)
			continue
		}
		logger.Log("")
		logger.Title("Process Profile: " + profile.Name)
		profile.Extensions = ip.profileProcessor.getAllProfileData(profile.Name)
		if ip.profileProcessor.isValidProfileName(profile.Name) {
			ip.processByProfile(profile)
		}
		if profile.CanInstall {
			logger.Header("Processing Profile: " + profile.Name + ", Done.")
		}
	}
	ip.setSettingConfigurations()
	if len(skipProfiles) > 0 {
		logger.Header("SKIPED PROFILES TO INSTALL")
		for _, profile := range skipProfiles {
			logger.Log("- " + profile)
		}
	}
}

func (ip *InstallProcessor) askTryInstallAllExtensions(extensions []string) bool {
	fmt.Println("\n####### NOT INSTALLED EXTENSIONS ID'S #######")
	for index, id := range extensions {
		fmt.Println(strconv.Itoa(index) + " - " + id)
	}
	return console.Confirm("Continue", false)
}

func (ip *InstallProcessor) isExtensionInstalled(id string, extensionsInstalled []string) bool {
	return slices.Contains(extensionsInstalled, strings.ToLower(id))
}

func (ip *InstallProcessor) getInstalledExtensions(profileName string) []string {
	listExtensions := []string{}
	command := libs.GetVscodeListExtCommand(profileName)
	output, err := exe.Exec(command)
	if err != nil {
		logger.Error(err)
		return listExtensions
	}
	if len(output) > 0 {
		for _, extension := range strings.Split(output, "\n") {
			listExtensions = append(listExtensions, strings.ToLower(strings.TrimSpace(extension)))
		}
	}
	return listExtensions
}

func (ip *InstallProcessor) getExtensionsToInstall(profileName string, extensions []entities.ProfileData) []string {
	listExtensions := ip.getInstalledExtensions(profileName)
	listToInstal := []string{}
	for _, data := range extensions {
		for _, id := range data.Ids {
			if !ip.isExtensionInstalled(id, listExtensions) {
				listToInstal = append(listToInstal, id)
			}
		}
	}
	return listToInstal
}

func (ip *InstallProcessor) installExtension(profileName string, id string) {
	command := libs.GetCodeCommand()
	fileId := ip.downloadProcessor.getExtensionVsixFile(id)
	if !file.IsFile(fileId) {
		ip.downloadProcessor.download(id)
		if file.IsFile(fileId) {
			id = fileId
		}
	} else {
		id = fileId
	}
	command.Args = append(command.Args, profileName, "--install-extension", id)
	command.Verbose = true
	resp, err := exe.Exec(command)
	logger.Error(err)
	if !strings.Contains(resp, "was successfully installed") {
		logger.ErrorStr(resp)
	}
}

func (ip *InstallProcessor) processByProfile(profile entities.Profile) {
	var profileName = profile.Name
	ip.downloadProcessor.force = false
	if profile.IsSettingName {
		profileName = libs.Configurations.SettingsName
	}
	ip.profileProcessor.createProfile(profileName)
	ip.profileAlreadyCreated = append(ip.profileAlreadyCreated, profileName)
	counter := 1
	for {
		extensionsToInstall := ip.getExtensionsToInstall(profileName, profile.Extensions)
		if len(extensionsToInstall) == 0 {
			break
		} else {
			logger.Header("Download Extensions")
			ip.downloadProcessor.downloadList(extensionsToInstall)
		}
		if counter > MAX_INSTALL_EXTENSIONS && len(extensionsToInstall) > 0 {
			if ip.askTryInstallAllExtensions(extensionsToInstall) {
				ip.downloadProcessor.force = true
				if !profile.IsSettingName {
					ip.profileProcessor.createProfile(profileName)
				}
				counter = 1
			} else {
				break
			}
		}
		logger.Header("Install Extensions")
		for _, data := range profile.Extensions {
			for _, id := range data.Ids {
				if slices.Contains(extensionsToInstall, id) {
					ip.installExtension(profileName, id)
				}
			}
		}
		counter++
	}
}
