package types

type App struct {
	Name  string          `yaml:"name"`
	Parts map[string]Part `yaml:"parts"`
}

type Part struct {
	Src  string `yaml:"src"`
	Path string `yaml:"path,omitempty"`
	Dev  string `yaml:"dev"`
}

type ApiAppResponse struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	InstalledAt string `json:"installedAt"`
}
