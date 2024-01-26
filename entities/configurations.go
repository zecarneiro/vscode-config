package entities

type Configurations struct {
	Settings     map[string]any    `json:"settings"`
	SettingsName string `json:"settingsProfileName"`
	Profiles     []Profile
}
