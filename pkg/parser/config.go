package parser

import (
	"log"

	"github.com/BurntSushi/toml"
)

// Config ...
type Config struct {
	Vladis        string `toml:"vladis"`
	Domofond      string `toml:"domofond"`
	Cian          string `toml:"cian"`
	Avito         string `toml:"avito"`
	Ioffepartners string `toml:"ioffepartners"`
}

// NewConfig ...
func NewConfig(path string) Config {
	conf := Config{}
	_, err := toml.DecodeFile(path, &conf)
	if err != nil {
		log.Fatal(err)
	}
	return conf
}
