package apiserver

// Server config
type Config struct {
	Port        string `toml:"port"`
	DatabaseURL string `toml:"database_url"`
}

// Init new config
func NewConfig() *Config {
	return &Config{
		Port: ":8080",
	}
}
