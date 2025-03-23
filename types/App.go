package types

type App struct {
	Id    string          `yaml:"id"`
	Name  string          `yaml:"name"`
	Parts map[string]Part `yaml:"parts"`
}

type Part struct {
	Src  string `yaml:"src"`
	Path string `yaml:"path,omitempty"`
	Dev  string `yaml:"dev"`
}

type ApiAppListResponse struct {
	Id          string            `json:"id"`
	Name        string            `json:"name"`
	InstalledAt string            `json:"installedAt"`
	Parts       map[string]string `json:"parts"`
}

type ApiAppInstallRequestBody struct {
	File string `json:"file"`
}

type ApiAppUninstallRequestBody struct {
	Id string `json:"id"`
}
