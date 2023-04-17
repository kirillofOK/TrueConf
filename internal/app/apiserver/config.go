package apiserver

// Config
type Config struct {
	BindAddr string `toml:"bind_addr"`
	LogLevel string `toml:"log_level"`
	StoreURL string `toml:"store_url"`
}

// Default config
func NewConfig() *Config {
	return &Config{
		BindAddr: "8080",
		LogLevel: "debug",
		StoreURL: "users.json",
	}
}
