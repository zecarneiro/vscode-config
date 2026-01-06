package libs

import (
	"golangutils"
	"main/entities"
)

var LoggerUtils = golangutils.NewLoggerUtils()
var SystemUtils = golangutils.NewSystemUtils(LoggerUtils)
var ConsoleUtils = golangutils.NewConsoleUtils(LoggerUtils, SystemUtils)

var ConfigDir = golangutils.ResolvePath(SystemUtils.HomeDir() + "/.config/vscode-config")
var JsonFile = golangutils.ResolvePath(ConfigDir + "/config.json")

var JsonInfo *entities.JsonInfo
var Configurations *entities.Configurations
