package libs

import (
	"golangutils/pkg/file"
	"golangutils/pkg/system"
	"main/entities"
)

var ConfigDir = file.ResolvePath(system.HomeDir(), ".config/vscode-config")
var JsonFile = file.ResolvePath(ConfigDir, "config.json")

var JsonInfo *entities.JsonInfo
var Configurations *entities.Configurations
