package apiserver

// Config ...
type Config struct {
	Port        string `toml:"port"`
	DatabaseURL string `toml:"database_url"`
}

// NewConfig ...
func NewConfig() *Config {
	return &Config{
		Port: ":8080",
	}
}
