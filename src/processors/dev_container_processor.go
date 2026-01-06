package processors

import (
	"fmt"
	"golangutils/pkg/common"
	"golangutils/pkg/file"
	"golangutils/pkg/logger"
	"golangutils/pkg/logic"
	"golangutils/pkg/models"
	"reflect"

	"main/entities"
	"main/libs"
)

type DevContainerProcessor struct {
	profileProcessor *ProfileProcessor
	devContainerDir  string
	jsonFile         string
	dockerFile       string
	validTypeChoices map[string]bool
}

func newDevContainerProcessor(profileProcessor *ProfileProcessor) *DevContainerProcessor {
	return &DevContainerProcessor{
		profileProcessor: profileProcessor,
	}
}

func (d *DevContainerProcessor) loadData() {
	currentDir, err := file.GetCurrentDir()
	logic.ProcessError(err)
	d.devContainerDir = file.ResolvePath(currentDir, ".devcontainer-generated")
	d.jsonFile = file.ResolvePath(d.devContainerDir, "devcontainer.json")
	d.dockerFile = file.ResolvePath(d.devContainerDir, "Dockerfile")
	d.validTypeChoices = map[string]bool{"go": true}
}

func (d *DevContainerProcessor) validateType(containerType string) {
	if !d.validTypeChoices[containerType] {
		logic.ProcessError(fmt.Errorf("invalid type: %s. Valid choices are: go", containerType))
	}
}

func (d *DevContainerProcessor) createDir() {
	if file.IsDir(d.devContainerDir) {
		logic.ProcessError(file.DeleteDirectory(d.devContainerDir))
	}
	file.CreateDirectory(d.devContainerDir, true)
	if !file.IsDir(d.devContainerDir) {
		logic.ProcessError(fmt.Errorf("Creating directory: %s", d.devContainerDir))
	}
}

func (d *DevContainerProcessor) generateDockerfile(containerType string) {
	var template string
	switch containerType {
	case "go":
		template = libs.DockerfileGoTemplate
	default:
		template = ""
	}
	if len(template) > 0 {
		logger.Info("Generate dev container Dockerfile for profile: " + containerType)
		template = fmt.Sprintf("%s %s%s", template, common.Eol(), libs.WorkspaceTemplate)
		fileConfig := models.FileWriterConfig{
			File:        d.dockerFile,
			Data:        template,
			IsAppend:    false,
			IsCreateDir: true,
			WithUtf8BOM: false,
		}
		logic.ProcessError(file.WriteFile(fileConfig))
	}
}

func (d *DevContainerProcessor) generateJsonFile(name string, postCreateCommand string) {
	devContainer := entities.DevContainer{}
	for _, profile := range libs.Configurations.Profiles {
		if profile.Name == name {
			devContainer = entities.DevContainer{
				Name:              profile.Name + " Dev Container",
				DockerFile:        "Dockerfile",
				Context:           "..",
				RemoteUser:        "ubuntu",
				WorkspaceFolder:   "/workspace",
				WorkspaceMount:    "source=${localWorkspaceFolder},target=/workspace,type=bind,consistency=cached",
				Settings:          d.profileProcessor.getSettings(profile.Name),
				Extensions:        d.profileProcessor.getAllExtensionsFromProfile(profile.Name),
				PostCreateCommand: postCreateCommand,
			}
			break
		}
	}
	if !reflect.DeepEqual(devContainer, entities.DevContainer{}) {
		logger.Info("Generate dev container JSON file for profile: " + name)
		logic.ProcessError(file.WriteJsonFile(d.jsonFile, devContainer, true))
	}
}

func (d *DevContainerProcessor) generate(name string, postCreateCommand string, containerType string) {
	d.loadData()
	d.validateType(containerType)
	d.createDir()
	d.generateDockerfile(containerType)
	d.generateJsonFile(name, postCreateCommand)
}
