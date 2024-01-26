package libs

import (
	"fmt"
	utils "jnoronha_golangutils"
	utilsEntities "jnoronha_golangutils/entities"
	"log"
	"main/entities"
	"strings"
)

/* ------------------------------ PRIVATE AREA ------------------------------ */

const ()

var jsonInfo entities.JsonInfo
var codeCommand utilsEntities.CommandInfo
var configurations entities.Configurations

func readJsonFile(jsonFile string) {
	utils.InfoLog("Reading file: "+jsonFile, false)
	data, err := utils.ReadFile(jsonFile)
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
	response := utils.Exec(command)
	if response.HasData() && strings.Contains(response.Data, dataError) {
		return false
	}
	return true
}

func openVscodeWithNewProfile(name string) {
	utils.InfoLog("Please verify if new profile are created and then close the VsCode", false)
	command := codeCommand
	command.Args = append(command.Args, name, "--wait")
	command.Verbose = true
	utils.ExecRealTime(command)
}

func setSettingConfigurations() {
	const keyAllSettings = "workbench.settings.applyToAllProfiles"
	var arrayAllSettings = []string{}
	for key := range configurations.Settings {
		arrayAllSettings = append(arrayAllSettings, key)
	}
	configurations.Settings[keyAllSettings] = arrayAllSettings
	utils.DebugLog("\n\nGo to setting and open json settings", false)
	utils.DebugLog("Append this setting bellow on json settings", false)
	fmt.Println(utils.ObjectToString(configurations.Settings))
	openVscodeWithNewProfile(configurations.SettingsName)
}

/* ------------------------------- PUBLIC AREA ------------------------------ */
func Start(jsonFile string) {
	codeCommand = utilsEntities.CommandInfo{
		Cmd:     "code",
		Args:    []string{"--profile"},
		Verbose: true,
		IsThrow: false,
	}
	jsonInfo = entities.JsonInfo{File: jsonFile}
	readJsonFile(jsonFile)
	object, err := utils.StringToObject[entities.Configurations](jsonInfo.Data)
	utils.ProcessError(err)
	configurations = object
	InitInstallProcessor(configurations.Profiles)
	setSettingConfigurations()
}
