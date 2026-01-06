package entities

type Configurations struct {
	SettingsName string `json:"settingsProfileName"`
	Profiles     []Profile
}
