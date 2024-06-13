package libs

import (
	"errors"
	"fmt"
	"jnoronhautils"
	jnoronhautilsEntities "jnoronhautils/entities"
	"log"
	"main/entities"
	"strings"
)

/* ------------------------------ PRIVATE AREA ------------------------------ */

const ()

var jsonInfo entities.JsonInfo
var codeCommand jnoronhautilsEntities.CommandInfo
var configurations entities.Configurations

func readJsonFile(jsonFile string) {
	jnoronhautils.InfoLog("Reading file: "+jsonFile, false)
	data, err := jnoronhautils.ReadFile(jsonFile)
	if err != nil {
		log.Fatalf(err.Error())
	}
	jsonInfo.Data = data
}

func profileExistOnVscode(name string) bool {
	dataError := "Profile '" + name + "' not found."
	command := codeCommand
	command.Args = append(command.Args, name, "--list-extensions")
	command.Verbose = false
	response := jnoronhautils.Exec(command)
	if response.HasData() && strings.Contains(response.Data, dataError) {
		return false
	}
	return true
}

func openVscodeWithNewProfile(name string) {
	jnoronhautils.InfoLog("Please verify if new profile are created and then close the VsCode", false)
	command := codeCommand
	command.Args = append(command.Args, name, "--wait")
	command.Verbose = true
	jnoronhautils.ExecRealTime(command)
}

func setSettingConfigurations() {
	const keyAllSettings = "workbench.settings.applyToAllProfiles"
	var settingsDir string
	var arrayAllSettings = []string{}
	if jnoronhautils.IsWindows() {
		settingsDir = jnoronhautils.ResolvePath(jnoronhautils.SystemInfo().HomeDir + "\\AppData\\Roaming\\Code\\User")
	} else if jnoronhautils.IsLinux() {
		settingsDir = jnoronhautils.ResolvePath(jnoronhautils.SystemInfo().HomeDir + "/.config/Code/User")
	}
	for key := range configurations.Settings {
		arrayAllSettings = append(arrayAllSettings, key)
	}
	configurations.Settings[keyAllSettings] = arrayAllSettings
	if len(settingsDir) > 0 {
		jnoronhautils.CreateDirectory(settingsDir, true)
		jnoronhautils.WriteJsonFile(jnoronhautils.ResolvePath(settingsDir+"/settings.json"), configurations.Settings)
	} else {
		jnoronhautils.DebugLog("\n\nGo to setting and open json settings", false)
		jnoronhautils.DebugLog("Append this setting bellow on json settings", false)
		fmt.Println(jnoronhautils.ObjectToString(configurations.Settings))
		openVscodeWithNewProfile(configurations.SettingsName)
	}
}

/* ------------------------------- PUBLIC AREA ------------------------------ */
func Start(jsonFile string) {
	codeCommand = jnoronhautilsEntities.CommandInfo{
		Cmd:     "code",
		Args:    []string{"--profile"},
		Verbose: true,
		IsThrow: false,
	}
	jsonInfo = entities.JsonInfo{File: jsonFile}
	readJsonFile(jsonFile)
	object, err := jnoronhautils.StringToObject[entities.Configurations](jsonInfo.Data)
	jnoronhautils.ProcessError(err)
	jnoronhautils.CreateDirectory(getExtensionVsixPath(), true)
	if !jnoronhautils.FileExist(getExtensionVsixPath()) {
		jnoronhautils.ProcessError(errors.New(getExtensionVsixPath() + " does not exist. Please, please create before running."))
	}
	configurations = object
	InitInstallProcessor(configurations.Profiles)
	setSettingConfigurations()
}
