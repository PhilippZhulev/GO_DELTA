package confiiguration

// Config ...
// Протокол конфигураций
type Config struct {
	BindAddr    string `toml:"bind_addr"`
	LogLevel    string `toml:"log_level"`
	DatabaseURL string `toml:"database_url"`
	SessionKey  string `toml:"session_key"`
	Salt        string `toml:"salt"`
}

// NewConfig ...
// Задать конфигурации
func NewConfig() *Config {
	return &Config{
		BindAddr: ":4444",
		LogLevel: "debug",
	}
}
