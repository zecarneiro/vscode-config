package libs

import (
	"fmt"
	"jnoronhautils"
	"main/entities"
	"strconv"
	"strings"
	"time"
)

/* --------------------------------- PRIVATE -------------------------------- */
const (
	MAX_INSTALL_EXTENSIONS = 5
)

var profiles []entities.Profile
var profileAlreadyCreated = []string{}

func isValidProfileName(name string) bool {
	isValid := true
	if strings.Contains(name, " ") {
		isValid = false
		jnoronhautils.ErrorLog("Invalid profile name: "+name, false)
		jnoronhautils.InfoLog("Name must not contains: space", false)
	}
	return isValid
}

func askTryInstallAllExtensions(extensions []string) bool {
	fmt.Println("\n####### NOT INSTALLED EXTENSIONS ID'S #######")
	for index, id := range extensions {
		fmt.Println(strconv.Itoa(index) + " - " + id)
	}
	if !jnoronhautils.Confirm("Continue", false) {
		return false
	}
	return true
}

func isExtensionInstalled(id string, extensionsInstalled []string) bool {
	return jnoronhautils.InArray[string](extensionsInstalled, strings.ToLower(id))
}

func getInstalledExtensions(profileName string) []string {
	listExtensions := []string{}
	command := codeCommand
	command.Args = append(command.Args, profileName, "--list-extensions")
	command.Verbose = false
	response := jnoronhautils.Exec(command)
	if !response.HasError() && response.HasData() {
		for _, extension := range strings.Split(response.Data, "\n") {
			listExtensions = append(listExtensions, strings.ToLower(strings.TrimSpace(extension)))
		}
	} else {
		if response.HasData() {
			jnoronhautils.ErrorLog(response.Data, false)
		}
		if response.HasError() {
			jnoronhautils.ErrorLog(response.Error.Error(), false)
		}
	}
	return listExtensions
}

func getExtensionsToInstall(profileName string, extensions []entities.ProfileData) []string {
	listExtensions := getInstalledExtensions(profileName)
	listToInstal := []string{}
	for _, data := range extensions {
		for _, id := range data.Ids {
			if !isExtensionInstalled(id, listExtensions) {
				listToInstal = append(listToInstal, id)
			}
		}
	}
	return listToInstal
}

func installExtensionNoRetries(profileName string, id string, cwd string) {
	command := codeCommand
	command.Args = append(command.Args, profileName, "--install-extension", id)
	command.Verbose = true
	command.Cwd = cwd
	jnoronhautils.ExecRealTime(command)
}

func installExtension(profileName string, id string, cwdVsix string) {
	if len(cwdVsix) > 0 {
		installExtensionNoRetries(profileName, id+".vsix", cwdVsix)
	} else {
		installExtensionNoRetries(profileName, id, "")
	}
}

func downloadAllExtensions(ids []string) {
	if len(ids) > 0 {
		extensionsToDownload := []string{}
		counter := 0
		for {
			for _, id := range ids {
				idFile := getExtensionVsixFile(id)
				if !jnoronhautils.FileExist(idFile) {
					if !download(id) || !jnoronhautils.FileExist(idFile) {
						extensionsToDownload = append(extensionsToDownload, id)
					}
				}
			}
			if len(extensionsToDownload) == 0 {
				break
			}
			if counter > MAX_INSTALL_EXTENSIONS && len(extensionsToDownload) > 0 {
				if askTryInstallAllExtensions(extensionsToDownload) {
					counter = 1
				}
			}
			counter++
			time.Sleep(3 * time.Second)
			concatenated := fmt.Sprintf("#%d Try to download all extensions again", counter)
			jnoronhautils.WarnLog(concatenated, false)
		}
	}
}

func processProfile(profile entities.Profile) {
	var profileName = profile.Name
	if profile.IsSettingName {
		profileName = configurations.SettingsName
	}
	if !jnoronhautils.InArray[string](profileAlreadyCreated, profileName) {
		if !profileExistOnVscode(profileName) {
			if profileName != configurations.SettingsName {
				copyFrom := profile.CopyFrom
				if len(copyFrom) == 0 {
					copyFrom = configurations.SettingsName
				}
				jnoronhautils.Separatorlog(45)
				jnoronhautils.LogLog("Create Prefile on VS Code: "+profileName, false)
				jnoronhautils.LogLog("1 - Run command 'Profile: Create Profile' - CTRL+SHIFT+P", false)
				jnoronhautils.LogLog("2 - Insert '"+profileName+"' on PROFILE_NAME_FIELD", false)
				jnoronhautils.LogLog("3 - Insert '"+copyFrom+"' on 'Copy From' field", false)
				jnoronhautils.LogLog("4 - Only keep 'Extensions' on 'Choose what to configure in your Profile'", false)
			}
			openVscodeWithNewProfile(configurations.SettingsName)
		}
	}
	counter := 1
	for {
		extensionsToInstall := getExtensionsToInstall(profileName, profile.Extensions)
		if len(extensionsToInstall) == 0 {
			break
		}
		downloadAllExtensions(extensionsToInstall)
		if counter > MAX_INSTALL_EXTENSIONS && len(extensionsToInstall) > 0 {
			if askTryInstallAllExtensions(extensionsToInstall) {
				counter = 1
			}
		}
		for _, data := range profile.Extensions {
			for _, id := range data.Ids {
				if jnoronhautils.InArray[string](extensionsToInstall, id) {
					fmt.Println("\n#####")
					fmt.Println("INSTALL FOR PROFILE: " + profileName)
					fmt.Println("DESCRIPTION: " + data.Descriptions)
					installExtension(profileName, id, getExtensionVsixPath())
					fmt.Println("#####")
				}
			}
		}
		counter++
	}
}

func processDependsProfiles(profileName string, dependsName []string) []entities.ProfileData {
	var listExtensions = []entities.ProfileData{}
	if len(dependsName) > 0 {
		for _, profile := range profiles {
			if profile.Name != profileName && isValidProfileName(profile.Name) && jnoronhautils.InArray[string](dependsName, profile.Name) {
				jnoronhautils.DebugLog("Append all extensions from Depends Profile: "+profile.Name, false)
				listExtensions = append(listExtensions, profile.Extensions...)
			}
		}
	}
	return listExtensions
}

/* --------------------------------- PUBLIC --------------------------------- */
func InitInstallProcessor(values []entities.Profile) {
	profiles = values
	for _, profile := range profiles {
		jnoronhautils.LogLog("", false)
		jnoronhautils.InfoLog("====== Process Profile: "+profile.Name+" ======", false)
		profile.Extensions = append(profile.Extensions, processDependsProfiles(profile.Name, profile.DependsProfile)...)
		if isValidProfileName(profile.Name) {
			processProfile(profile)
		}
	}
}
