package entities

type Profile struct {
	Name           string `json:"name"`
	Extensions     []ProfileData
	DependsProfile []string `json:"dependsProfile,omitempty"`
	IsSettingName  bool     `json:"settingName"`
	CopyFrom       string   `json:"copyFrom"`
	CanInstall     bool     `json:"canInstall"`
}
