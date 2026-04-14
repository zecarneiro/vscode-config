package processors

import (
	"fmt"
	"golangutils/pkg/console"
	"golangutils/pkg/exe"
	"golangutils/pkg/logger"
	"golangutils/pkg/logic"
	"golangutils/pkg/models"
	"golangutils/pkg/platform"
	"golangutils/pkg/slice"
	"golangutils/pkg/str"
	"slices"
	"strings"
	"vscodeconfig/core/entities"
	"vscodeconfig/core/libs"
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
		logger.ErrorStr("Invalid profile name: " + name)
		logger.Info("Name must not contains: space")
	}
	return isValid
}

func (pp *ProfileProcessor) profileExists(name string) bool {
	command := libs.GetVscodeListExtCommand(name)
	command.Verbose = false
	result, _ := exe.Exec(command)
	if strings.Contains(result, "Profile '"+name+"' not found.") {
		return false
	}
	return true
}

func (pp *ProfileProcessor) createProfile(name string) {
	if !pp.profileExists(name) {
		var commandStr string
		if platform.IsLinux() {
			commandStr = "(code --profile \"%s\" &) && sleep 2 && killall -9 code && sleep 2"
		} else if platform.IsWindows() {
			commandStr = "Start-Process code -ArgumentList '--profile \"%s\"'; Start-Sleep -Seconds 2; Get-Process code -ErrorAction SilentlyContinue | Stop-Process -Force; Start-Sleep -Seconds 2"
		} else {
			commandStr = ""
			logger.ErrorStr("Can not create this profile, because this SO is not supported")
			logger.Info(fmt.Sprintf("Please open VSCode and create this profile: %s", name))
			console.Pause()
		}
		if len(commandStr) > 0 {
			logger.Header("Creating Profile")
			command := models.Command{Cmd: fmt.Sprintf(commandStr, name), UseShell: true, Verbose: false}
			logic.ProcessError(exe.ExecRealTime(command))
		}
	}
}

func (pp *ProfileProcessor) processDependsProfiles(profileName string, dependsName []string, verbose bool) []entities.ProfileData {
	listExtensions := []entities.ProfileData{}
	if len(dependsName) > 0 {
		for _, profile := range pp.profiles {
			if profile.Name != profileName && pp.isValidProfileName(profile.Name) && slices.Contains(dependsName, profile.Name) {
				if verbose {
					logger.Header("Append extensions from Profile: " + profile.Name)
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

func (pp *ProfileProcessor) isSettingsName(Name string) bool {
	for _, profile := range libs.Configurations.Profiles {
		if profile.Name == Name && profile.IsSettingName {
			return true
		}
	}
	return false
}

func (pp *ProfileProcessor) processDependsProfilesSettings(profileName string, dependsName []string, verbose bool) map[string]interface{} {
	settings := make(map[string]interface{})
	if len(dependsName) > 0 {
		for _, profile := range pp.profiles {
			if profile.Name != profileName && pp.isValidProfileName(profile.Name) && slices.Contains(dependsName, profile.Name) {
				if verbose {
					logger.Header("Append settings from Profile: " + profile.Name)
				}
				settings = slice.ConcatMap(settings, profile.Settings)
			}
		}
	}
	return settings
}

func (pp *ProfileProcessor) getSettings(name string) map[string]interface{} {
	settings := make(map[string]interface{})
	isSettingsName := pp.isSettingsName(name)
	for _, profile := range libs.Configurations.Profiles {
		if profile.Name == name {
			settings = slice.ConcatMap(settings, profile.Settings)
			if !isSettingsName {
				settings = slice.ConcatMap(settings, pp.processDependsProfilesSettings(profile.Name, profile.DependsProfile, true))
			}
			break
		}
	}
	return settings
}

func (pp *ProfileProcessor) getAllInstallSettings(profileName string, verbose bool) map[string]interface{} {
	settings := make(map[string]interface{})
	profileSelectedList := slice.FilterArray(libs.Configurations.Profiles, func(profile entities.Profile) bool {
		if str.IsEmpty(profileName) {
			return logic.Ternary(profile.IsSettingName, true, false)
		}
		return logic.Ternary(profile.Name == profileName, true, false)
	})
	profileSelected := logic.Ternary(len(profileSelectedList) > 0, profileSelectedList[0], entities.Profile{Name: "", DependsProfile: []string{}})
	for _, profile := range libs.Configurations.Profiles {
		if slices.Contains(profileSelected.DependsProfile, profile.Name) {
			if verbose {
				logger.Debug(fmt.Sprintf("Import settings from profile: %s", profile.Name))
			}
			settings = slice.ConcatMap(settings, profile.Settings)
		} else if profile.Name == profileSelected.Name {
			if verbose {
				logger.Debug(fmt.Sprintf("Import self settings: %s", profile.Name))
			}
			settings = slice.ConcatMap(settings, profile.Settings)
		}
	}
	return settings
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
