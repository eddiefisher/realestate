package config

// Config ...
type Config struct {
	Vladis        string `toml:"vladis"`
	Domofond      string `toml:"domofond"`
	Cian          string `toml:"cian"`
	Avito         string `toml:"avito"`
	Ioffepartners string `toml:"ioffepartners"`
}

// NewConfig ...
func NewConfig() *Config {
	return &Config{}
}
