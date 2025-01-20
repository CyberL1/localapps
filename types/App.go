package types

type App struct {
	Name  string `yaml:"name"`
	Parts []Part `yaml:"parts"`
}

type Part struct {
	// Name    string `yaml:"name"`
	Src     string `yaml:"src"`
	Path    string `yaml:"path,omitempty"`
	Primary bool   `yaml:"primary,omitempty"`
	Run     string `yaml:"run"`
}
