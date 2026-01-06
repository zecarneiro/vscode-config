package processors

import (
	"fmt"
	"golangutils"
	"golangutils/entity"
	"main/entities"
	"main/libs"
	"strings"
)

type ProfileProcessor struct {
	profiles []entities.Profile
}

func newProfileProcessor() *ProfileProcessor {
	processor := &ProfileProcessor{}
	processor.loadData()
	return processor
}

func (pp *ProfileProcessor) loadData() {
	if libs.IsJsonFileInstalled() {
		pp.profiles = libs.Configurations.Profiles
	}
}

func (pp *ProfileProcessor) isValidProfileName(name string) bool {
	isValid := true
	if strings.Contains(name, " ") {
		isValid = false
		libs.LoggerUtils.Error("Invalid profile name: " + name)
		libs.LoggerUtils.Info("Name must not contains: space")
	}
	return isValid
}

func (pp *ProfileProcessor) profileExists(name string) bool {
	command := libs.GetVscodeListExtCommand(name)
	command.Verbose = false
	result, _ := libs.ConsoleUtils.Exec(command)
	if strings.Contains(result, "Profile '"+name+"' not found.") {
		return false
	}
	return true
}

func (pp *ProfileProcessor) createProfile(name string) {
	if !pp.profileExists(name) {
		var commandStr string
		if libs.SystemUtils.IsLinux() {
			commandStr = "(code --profile \"%s\" &) && sleep 2 && killall -9 code && sleep 2"
		} else if libs.SystemUtils.IsWindows() {
			commandStr = "Start-Process code -ArgumentList '--profile \"%s\"'; Start-Sleep -Seconds 2; Get-Process code -ErrorAction SilentlyContinue | Stop-Process -Force; Start-Sleep -Seconds 2"
		} else {
			commandStr = ""
			libs.LoggerUtils.Error("Can not create this profile, because this SO is not supported")
			libs.LoggerUtils.Info(fmt.Sprintf("Please open VSCode and create this profile: %s", name))
			golangutils.Pause("")
		}
		if len(commandStr) > 0 {
			libs.LoggerUtils.Header("Creating Profile")
			command := entity.Command{Cmd: fmt.Sprintf(commandStr, name), UseShell: true, Verbose: false}
			libs.ConsoleUtils.ExecRealTime(command)
		}
	}
}

func (pp *ProfileProcessor) processDependsProfiles(profileName string, dependsName []string, verbose bool) []entities.ProfileData {
	var listExtensions = []entities.ProfileData{}
	if len(dependsName) > 0 {
		for _, profile := range pp.profiles {
			if profile.Name != profileName && pp.isValidProfileName(profile.Name) && golangutils.InArray(dependsName, profile.Name) {
				if verbose {
					libs.LoggerUtils.Header("Append extensions from Profile: " + profile.Name)
				}
				listExtensions = append(listExtensions, profile.Extensions...)
			}
		}
	}
	return listExtensions
}

func (pp *ProfileProcessor) getAllProfileData(name string) []entities.ProfileData {
	profileData := []entities.ProfileData{}
	for _, profile := range libs.Configurations.Profiles {
		if profile.Name == name {
			profileData = append(profile.Extensions, pp.processDependsProfiles(profile.Name, profile.DependsProfile, false)...)
			break
		}
	}
	return profileData
}

func (pp *ProfileProcessor) getAllExtensionsFromProfile(name string) []string {
	extensions := []string{}
	for _, data := range pp.getAllProfileData(name) {
		extensions = append(extensions, data.Ids...)
	}
	return extensions
}

func (pp *ProfileProcessor) getAllProfile() []entities.ProfileStatus {
	status := []entities.ProfileStatus{}
	for _, profile := range libs.Configurations.Profiles {
		name := profile.Name
		if profile.IsSettingName {
			name = libs.Configurations.SettingsName
		}
		status = append(status, entities.ProfileStatus{Name: name, IsInstalled: pp.profileExists(name)})
	}
	return status
}
