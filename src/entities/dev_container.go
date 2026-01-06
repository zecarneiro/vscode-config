package entities

type DevContainer struct {
	Name              string                 `json:"name"`
	DockerFile        string                 `json:"dockerFile"`
	Context           string                 `json:"context"`
	RemoteUser        string                 `json:"remoteUser"`
	WorkspaceFolder   string                 `json:"workspaceFolder"`
	WorkspaceMount    string                 `json:"workspaceMount"`
	Settings          map[string]interface{} `json:"settings"`
	Extensions        []string               `json:"extensions"`
	PostCreateCommand string                 `json:"postCreateCommand"`
}
