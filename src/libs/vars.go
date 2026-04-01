package libs

import (
	"golangutils/pkg/file"
	"golangutils/pkg/system"
	"vscodeconfig/core/entities"
)

var ConfigDir = file.JoinPath(system.HomeDir(), ".config/vscode-config")
var JsonFile = file.JoinPath(ConfigDir, "config.json")

var JsonInfo *entities.JsonInfo
var Configurations *entities.Configurations
