package libs

import (
	"golangutils"
	entityutils "golangutils/entity"
	"main/entities"
)

func GetCodeCommand() entityutils.Command {
	command := entityutils.Command{
		Cmd:      "code",
		Args:     []string{"--profile"},
		Verbose:  true,
		UseShell: false,
		IsThrow:  false,
	}
	return command
}

func OpenVscodeWithNewProfile(name string) {
	command := GetCodeCommand()
	command.Args = append(command.Args, name, "--wait")
	command.Verbose = true
	ConsoleUtils.ExecRealTime(command)
}

func GetVscodeListExtCommand(profileName string) entityutils.Command {
	command := GetCodeCommand()
	command.Args = append(command.Args, profileName, "--list-extensions")
	command.Verbose = false
	return command
}

func IsJsonFileInstalled() bool {
	return golangutils.IsFile(JsonFile)
}

func FillJsonFile(verbose bool) {
	if IsJsonFileInstalled() {
		if verbose {
			LoggerUtils.Info("Reading file: " + JsonFile)
		}
		data, err := golangutils.ReadFile(JsonFile)
		golangutils.ProcessError(err)
		JsonInfo = &entities.JsonInfo{File: JsonFile, Data: data}
		object, err := golangutils.StringToObject[entities.Configurations](JsonInfo.Data)
		golangutils.ProcessError(err)
		Configurations = &object
	} else {
		LoggerUtils.Warn("JSON File is not installed!!")
	}
}
