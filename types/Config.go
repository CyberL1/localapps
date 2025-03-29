package types

type Config struct {
	Domain string `default:"\"apps.localhost:8080\""`
	ApiKey string `default:"\"\""`
}
