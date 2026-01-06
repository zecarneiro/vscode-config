package libs

import (
	"golangutils/pkg/exe"
	"golangutils/pkg/file"
	"golangutils/pkg/logger"
	"golangutils/pkg/logic"
	"golangutils/pkg/models"
	"golangutils/pkg/obj"

	"main/entities"
)

func GetCodeCommand() models.Command {
	command := models.Command{
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
	exe.ExecRealTime(command)
}

func GetVscodeListExtCommand(profileName string) models.Command {
	command := GetCodeCommand()
	command.Args = append(command.Args, profileName, "--list-extensions")
	command.Verbose = false
	return command
}

func IsJsonFileInstalled() bool {
	return file.IsFile(JsonFile)
}

func FillJsonFile(verbose bool) {
	if IsJsonFileInstalled() {
		if verbose {
			logger.Info("Reading file: " + JsonFile)
		}
		data, err := file.ReadFile(JsonFile)
		logic.ProcessError(err)
		JsonInfo = &entities.JsonInfo{File: JsonFile, Data: data}
		object, err := obj.StringToObject[entities.Configurations](JsonInfo.Data)
		logic.ProcessError(err)
		Configurations = &object
	} else {
		logger.Warn("JSON File is not installed!!")
	}
}
