package types

type App struct {
	Name  string `yaml:"name"`
	Parts map[string]Part `yaml:"parts"`
}

type Part struct {
	Src     string `yaml:"src"`
	Path    string `yaml:"path,omitempty"`
	Primary bool   `yaml:"primary,omitempty"`
	Run     string `yaml:"run"`
}
