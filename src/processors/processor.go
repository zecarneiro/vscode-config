package processors

import (
	"fmt"
	"golangutils/pkg/file"
	"golangutils/pkg/logger"
	"golangutils/pkg/logic"
	"golangutils/pkg/obj"
	"golangutils/pkg/platform"
	"golangutils/pkg/system"
	"os"
	"vscodeconfig/core/entities"
	"vscodeconfig/core/libs"

	"github.com/spf13/cobra"
)

type Processor struct {
	rootCmd *cobra.Command

	installProcessor      *InstallProcessor
	profileProfessor      *ProfileProcessor
	devContainerProcessor *DevContainerProcessor
}

/* ------------------------------ PRIVATE AREA ------------------------------ */
func (p *Processor) loadData() {
	libs.JsonInfo = &entities.JsonInfo{}
	logic.ProcessError(file.CreateDirectory(libs.ConfigDir, true))
	libs.FillJsonFile(false)
	p.profileProfessor = newProfileProcessor()
	p.installProcessor = newInstallProcessor(p.profileProfessor)
	p.devContainerProcessor = newDevContainerProcessor(p.profileProfessor)
}

func (p *Processor) installJsonFile(jsonFile string) {
	p.installProcessor.installJsonFile(jsonFile)
	libs.FillJsonFile(true)
	p.loadData()
}

func (p *Processor) extractProfile(name string) {
	data, err := obj.ObjectToString(p.profileProfessor.getAllExtensionsFromProfile(name))
	logic.ProcessError(err)
	logger.Log(data)
}

func (p *Processor) listProfiles() {
	for _, profile := range p.profileProfessor.getAllProfile() {
		if profile.IsInstalled {
			logger.Log("- " + profile.Name)
		} else {
			logger.Log("- " + profile.Name + " (Not Installed)")
		}
	}
}

func (p *Processor) resetVscode() {
	var pathsCmd []string
	if platform.IsLinux() {
		pathsCmd = []string{
			file.JoinPath(system.HomeUserConfigDir(), "Code"),
			file.JoinPath(system.HomeDir(), ".vscode"),
			file.JoinPath(system.HomeDir(), ".cache/code"),
		}
	} else if platform.IsWindows() {
		pathsCmd = []string{
			file.JoinPath(os.Getenv("APPDATA"), "Code"),
			file.JoinPath(system.HomeDir(), ".vscode"),
			file.JoinPath(os.Getenv("LOCALAPPDATA"), "Code"),
		}
	} else if platform.IsDarwin() {
		pathsCmd = []string{
			file.JoinPath(system.HomeDir(), "Library/Application Support/Code"),
			file.JoinPath(system.HomeDir(), ".vscode"),
			file.JoinPath(system.HomeDir(), "Library/Caches/com.microsoft.VSCode"),
			file.JoinPath(system.HomeDir(), "Library/Preferences/com.microsoft.VSCode.plist"),
		}
	} else {
		logger.ErrorStr(platform.InvalidMSG)
		pathsCmd = []string{}
	}
	for _, pathCmd := range pathsCmd {
		if file.IsFile(pathCmd) {
			logger.Info(fmt.Sprintf("Deleting file: %s", pathCmd))
			logger.Error(file.DeleteFile(pathCmd))
		}
		if file.IsDir(pathCmd) {
			logger.Info(fmt.Sprintf("Deleting directory: %s", pathCmd))
			logger.Error(file.DeleteDirectory(pathCmd))
		}
	}
}

func (p *Processor) parseArgs() {
	var extractCmd *cobra.Command
	var installCmd *cobra.Command
	var devContainerCmd *cobra.Command

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
				logger.Log("VSCode profile " + profileName + " Exists: " + status)
			}
		},
	}
	p.rootCmd.Flags().BoolP("reset-vscode", "r", false, "Reset VSCode")
	p.rootCmd.Flags().StringP("profile-exists", "e", "", "Check if given profile exists")

	installCmd = &cobra.Command{
		Use:   "install",
		Short: "Install JSON file and start install process(Optional)",
		Run: func(cmd *cobra.Command, args []string) {
			jsonFile, _ := cmd.Flags().GetString("json-file")
			isProcessInstall, _ := cmd.Flags().GetBool("process-install")
			isSetSettings, _ := cmd.Flags().GetBool("settings")
			p.installJsonFile(jsonFile)
			if isProcessInstall {
				p.installProcessor.processInstall()
			}
			if isSetSettings {
				p.installProcessor.setSettingConfigurations()
			}
		},
	}
	installCmd.Flags().StringP("json-file", "j", "", "Copy JSON File with all configs and profiles")
	installCmd.Flags().BoolP("settings", "s", false, "Set settings from installed JSON file")
	installCmd.Flags().BoolP("process-install", "p", false, "Process install flag for profiles, settings and extensions")

	extractCmd = &cobra.Command{
		Use:   "extract",
		Short: "Extract Data from JSON File",
		Run: func(cmd *cobra.Command, args []string) {
			isListProfiles, _ := cmd.Flags().GetBool("list-profiles")
			extractProfileName, _ := cmd.Flags().GetString("profile")
			isSettings, _ := cmd.Flags().GetBool("settings")
			if isSettings {
				data, err := obj.ObjectToString(p.profileProfessor.getAllInstallSettings(extractProfileName, false))
				logic.ProcessError(err)
				logger.Log(data)
			} else if extractProfileName != "" {
				p.extractProfile(extractProfileName)
			} else if isListProfiles {
				p.listProfiles()
			} else {
				logger.Error(extractCmd.Help())
			}
		},
	}
	extractCmd.Flags().BoolP("settings", "s", false, "Extract settings from installed JSON file")
	extractCmd.Flags().BoolP("list-profiles", "l", false, "List all profiles names from installed JSON file")
	extractCmd.Flags().StringP("profile", "n", "", "Extract all extensions by profile name from installed JSON file")

	devContainerCmd = &cobra.Command{
		Use:   "dev-container [name]",
		Short: "Generate dev container file by profile name",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			devContainerPostCreateCommand, _ := cmd.Flags().GetString("post-create-command")
			containerType, _ := cmd.Flags().GetString("type")
			devContainerName := args[0]
			p.devContainerProcessor.generate(devContainerName, devContainerPostCreateCommand, containerType)
		},
	}
	devContainerCmd.Flags().StringP("type", "t", "go", "Type of dev container (go)")

	p.rootCmd.AddCommand(installCmd, extractCmd, devContainerCmd)
}

/* ------------------------------- PUBLIC AREA ------------------------------ */
func StartProcessor() {
	processor := &Processor{}
	processor.loadData()
	processor.parseArgs()
	if err := processor.rootCmd.Execute(); err != nil {
		logic.ProcessError(err)
	}
}
