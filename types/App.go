package types

type App struct {
	Name  string          `yaml:"name"`
	Parts map[string]Part `yaml:"parts"`
}

type Part struct {
	Src  string `yaml:"src"`
	Path string `yaml:"path,omitempty"`
	Run  string `yaml:"run"`
	Dev  string `yaml:"dev"`
}

type AppNameWithSubdomain struct {
	Name      string `json:"name"`
	Subdomain string `json:"subdomain"`
}
