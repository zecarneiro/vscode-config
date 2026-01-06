package processors

import (
	"fmt"
	"golangutils"
	"main/entities"
	"main/libs"
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
	jsonFile = golangutils.ResolvePath(jsonFile)
	if golangutils.IsFile(jsonFile) {
		libs.LoggerUtils.Info("Copy given file to: " + libs.JsonFile)
		err := golangutils.CopyFile(jsonFile, libs.JsonFile)
		golangutils.ProcessError(err)
	}
}

func (ip *InstallProcessor) setSettingConfigurations() {
	libs.LoggerUtils.Separator()
	libs.LoggerUtils.Header("Set settings")
	const keyAllSettings = "workbench.settings.applyToAllProfiles"
	var settingsDir string
	var arrayAllSettings = []string{}
	if libs.SystemUtils.IsWindows() {
		settingsDir = golangutils.ResolvePath(libs.SystemUtils.HomeDir() + "\\AppData\\Roaming\\Code\\User")
	} else if libs.SystemUtils.IsLinux() {
		settingsDir = golangutils.ResolvePath(libs.SystemUtils.HomeDir() + "/.config/Code/User")
	}
	for key := range libs.Configurations.Settings {
		arrayAllSettings = append(arrayAllSettings, key)
	}
	libs.Configurations.Settings[keyAllSettings] = arrayAllSettings
	if len(settingsDir) > 0 {
		golangutils.CreateDirectory(settingsDir, true)
		golangutils.WriteJsonFile(golangutils.ResolvePath(settingsDir+"/settings.json"), libs.Configurations.Settings, true)
	} else {
		libs.LoggerUtils.Debug("\n\nGo to setting and open json settings")
		libs.LoggerUtils.Debug("Append this setting bellow on json settings")
		fmt.Println(golangutils.ObjectToString(libs.Configurations.Settings))
		libs.OpenVscodeWithNewProfile(libs.Configurations.SettingsName)
	}
}

func (ip *InstallProcessor) processInstall() {
	skipProfiles := []string{}
	libs.ConsoleUtils.WaitForAnyKeyPressed("Please, close all instance of VSCode and PRESS ANY KEY TO CONTINUE...")
	for _, profile := range ip.profileProcessor.profiles {
		if !profile.CanInstall {
			skipProfiles = append(skipProfiles, profile.Name)
			continue
		}
		libs.LoggerUtils.Log("")
		libs.LoggerUtils.Title("Process Profile: " + profile.Name)
		profile.Extensions = ip.profileProcessor.getAllProfileData(profile.Name)
		if ip.profileProcessor.isValidProfileName(profile.Name) {
			ip.processByProfile(profile)
		}
		if profile.CanInstall {
			libs.LoggerUtils.Header("Processing Profile: " + profile.Name + ", Done.")
		}
	}
	ip.setSettingConfigurations()
	if len(skipProfiles) > 0 {
		libs.LoggerUtils.Header("SKIPED PROFILES TO INSTALL")
		for _, profile := range skipProfiles {
			libs.LoggerUtils.Log("- " + profile)
		}
	}
}

func (ip *InstallProcessor) askTryInstallAllExtensions(extensions []string) bool {
	fmt.Println("\n####### NOT INSTALLED EXTENSIONS ID'S #######")
	for index, id := range extensions {
		fmt.Println(strconv.Itoa(index) + " - " + id)
	}
	return golangutils.Confirm("Continue", false)
}

func (ip *InstallProcessor) isExtensionInstalled(id string, extensionsInstalled []string) bool {
	return golangutils.InArray(extensionsInstalled, strings.ToLower(id))
}

func (ip *InstallProcessor) getInstalledExtensions(profileName string) []string {
	listExtensions := []string{}
	command := libs.GetVscodeListExtCommand(profileName)
	output, err := libs.ConsoleUtils.Exec(command)
	if err != nil {
		libs.LoggerUtils.Error(err.Error())
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
	file := ip.downloadProcessor.getExtensionVsixFile(id)
	if !golangutils.IsFile(file) {
		ip.downloadProcessor.download(id)
		if golangutils.IsFile(file) {
			id = file
		}
	} else {
		id = file
	}
	command.Args = append(command.Args, profileName, "--install-extension", id)
	command.Verbose = true
	resp, err := libs.ConsoleUtils.Exec(command)
	if err != nil {
		libs.LoggerUtils.Error(err.Error())
	}
	if !strings.Contains(resp, "was successfully installed") {
		libs.LoggerUtils.Error(resp)
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
			libs.LoggerUtils.Header("Download Extensions")
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
		libs.LoggerUtils.Header("Install Extensions")
		for _, data := range profile.Extensions {
			for _, id := range data.Ids {
				if golangutils.InArray(extensionsToInstall, id) {
					ip.installExtension(profileName, id)
				}
			}
		}
		counter++
	}
}
