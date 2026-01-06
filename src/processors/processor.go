package processors

import (
	"golangutils"
	entityutils "golangutils/entity"
	"main/entities"
	"main/libs"
	"reflect"

	"github.com/spf13/cobra"
)

type Processor struct {
	rootCmd *cobra.Command

	installProcessor *InstallProcessor
	profileProfessor *ProfileProcessor
}

/* ------------------------------ PRIVATE AREA ------------------------------ */
func (p *Processor) loadData() {
	libs.JsonInfo = &entities.JsonInfo{}
	golangutils.CreateDirectory(libs.ConfigDir, true)
	libs.FillJsonFile(false)
	p.profileProfessor = newProfileProcessor()
	p.installProcessor = newInstallProcessor(p.profileProfessor)
}

func (p *Processor) installJsonFile(jsonFile string) {
	p.installProcessor.installJsonFile(jsonFile)
	libs.FillJsonFile(true)
	p.loadData()
}

func (p *Processor) extractSettings() {
	data, err := golangutils.ObjectToString(libs.Configurations.Settings)
	golangutils.ProcessError(err)
	libs.LoggerUtils.Log(data)
}

func (p *Processor) extractProfile(name string) {
	data, err := golangutils.ObjectToString(p.profileProfessor.getAllExtensionsFromProfile(name))
	golangutils.ProcessError(err)
	libs.LoggerUtils.Log(data)
}

func (p *Processor) listProfiles() {
	for _, profile := range p.profileProfessor.getAllProfile() {
		if profile.IsInstalled {
			libs.LoggerUtils.Log("- " + profile.Name)
		} else {
			libs.LoggerUtils.Log("- " + profile.Name + " (Not Installed)")
		}
	}
}

func (p *Processor) generateDevContainer(name string, postCreateCommand string) {
	devContainer := entities.DevContainer{}
	currentDir, err := golangutils.GetCurrentDir()
	golangutils.ProcessError(err)
	for _, profile := range libs.Configurations.Profiles {
		if profile.Name == name {
			devContainer = entities.DevContainer{
				Name:              profile.Name + " Dev Container",
				DockerFile:        "Dockerfile",
				Context:           "..",
				RemoteUser:        "root",
				WorkspaceFolder:   "/workspace",
				WorkspaceMount:    "source=${localWorkspaceFolder},target=/workspace,type=bind,consistency=cached",
				Settings:          libs.Configurations.Settings,
				Extensions:        p.profileProfessor.getAllExtensionsFromProfile(profile.Name),
				PostCreateCommand: postCreateCommand,
			}
			break
		}
	}
	if !reflect.DeepEqual(devContainer, entities.DevContainer{}) {
		libs.LoggerUtils.Info("Generate dev container JSON file for profile: " + name)
		err = golangutils.WriteJsonFile(golangutils.ResolvePath(currentDir+"/devcontainer.json"), devContainer, true)
		golangutils.ProcessError(err)
	}
}

func (p *Processor) resetVscode() {
	var pathsCmd []string
	if libs.SystemUtils.IsLinux() {
		pathsCmd = []string{
			"rm -rf ~/.config/Code",
			"rm -rf ~/.vscode",
			"rm -rf ~/.cache/code",
		}
	} else if libs.SystemUtils.IsWindows() {
		pathsCmd = []string{
			"Remove-Item -Recurse -Force $env:APPDATA\\Code",
			"Remove-Item -Recurse -Force $env:USERPROFILE\\.vscode",
			"Remove-Item -Recurse -Force $env:LOCALAPPDATA\\Code",
		}
	} else if libs.SystemUtils.IsDarwin() {
		pathsCmd = []string{
			"rm -rf ~/Library/Application\\ Support/Code",
			"rm -rf ~/.vscode",
			"rm -rf ~/Library/Caches/com.microsoft.VSCode",
			"rm -rf ~/Library/Preferences/com.microsoft.VSCode.plist",
		}
	} else {
		libs.LoggerUtils.Error(golangutils.GetInvalidPlatformMsg())
		pathsCmd = []string{}
	}
	for _, pathCmd := range pathsCmd {
		libs.ConsoleUtils.ExecRealTime(entityutils.Command{Cmd: pathCmd, UseShell: true, Verbose: true})
	}
}

func (p *Processor) parseArgs() {
	p.rootCmd = &cobra.Command{
		Short: "Process some VSCode configurations",
		Run: func(cmd *cobra.Command, args []string) {
			resetVscode, _ := cmd.Flags().GetBool("reset-vscode")
			profileName, _ := cmd.Flags().GetString("profile-exists")
			if resetVscode {
				p.resetVscode()
			} else {
				status := "No"
				if p.profileProfessor.profileExists(profileName) {
					status = "Yes"
				}
				libs.LoggerUtils.Log("VSCode profile " + profileName + " Exists: " + status)
			}
		},
	}
	p.rootCmd.Flags().BoolP("reset-vscode", "r", false, "Reset VSCode")
	p.rootCmd.Flags().StringP("profile-exists", "e", "", "Check if given profile exists")

	var installCmd = &cobra.Command{
		Use:   "install",
		Short: "Install JSON file and start install process(Optional)",
		Run: func(cmd *cobra.Command, args []string) {
			jsonFile, _ := cmd.Flags().GetString("json-file")
			isProcessInstall, _ := cmd.Flags().GetBool("process-install")
			p.installJsonFile(jsonFile)
			if isProcessInstall {
				p.installProcessor.processInstall()
			}
		},
	}
	installCmd.Flags().StringP("json-file", "j", "", "Copy JSON File with all configs and profiles")
	installCmd.Flags().BoolP("process-install", "p", false, "Process install flag for profiles, settings and extensions")

	var extractCmd = &cobra.Command{
		Use:   "extract",
		Short: "Extract Data from JSON File",
		Run: func(cmd *cobra.Command, args []string) {
			isExtractSettings, _ := cmd.Flags().GetBool("settings")
			isListProfiles, _ := cmd.Flags().GetBool("list-profiles")
			extractProfileName, _ := cmd.Flags().GetString("profile")
			switch {
			case isExtractSettings:
				p.extractSettings()
			case extractProfileName != "":
				p.extractProfile(extractProfileName)
			case isListProfiles:
				p.listProfiles()
			}
		},
	}
	extractCmd.Flags().BoolP("settings", "s", false, "Extract settings from installed JSON file")
	extractCmd.Flags().BoolP("list-profiles", "l", false, "List all profiles names from installed JSON file")
	extractCmd.Flags().StringP("profile", "n", "", "Extract all extensions by profile name from installed JSON file")
	extractCmd.MarkFlagsMutuallyExclusive("settings", "profile", "list-profiles")

	var devContainerCmd = &cobra.Command{
		Use:   "dev-container [name]",
		Short: "Generate dev container file by profile name",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			devContainerPostCreateCommand, _ := cmd.Flags().GetString("post-create-command")
			devContainerName := args[0]
			p.generateDevContainer(devContainerName, devContainerPostCreateCommand)
		},
	}
	devContainerCmd.Flags().StringP("post-create-command", "c", "", "Post create command for dev container")

	p.rootCmd.AddCommand(installCmd, extractCmd, devContainerCmd)
}

/* ------------------------------- PUBLIC AREA ------------------------------ */
func StartProcessor() {
	processor := &Processor{}
	processor.loadData()
	processor.parseArgs()
	if err := processor.rootCmd.Execute(); err != nil {
		golangutils.ProcessError(err)
	}
}
