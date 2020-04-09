package apiserver

// Server config
type Config struct {
	Port           string `toml:"port"`
	DatabaseURL    string `toml:"database_url"`
	TelegramToken  string `toml:"telegram_token"`
	TelegramUserId int    `toml:"telegram_user_id"`
}

// Init new config
func NewConfig() *Config {
	return &Config{
		Port: ":8080",
	}
}
