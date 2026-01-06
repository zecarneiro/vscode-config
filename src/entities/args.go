package entities

type Args struct {
	JsonFile       string
	ProcessInstall bool

	GenerateDevContainer          string
	DevContainerPostCreateCommand string
	ResetVSCode                   bool

	ExtractProfile  string
	ExtractSettings bool
	ListProfiles    bool
}
