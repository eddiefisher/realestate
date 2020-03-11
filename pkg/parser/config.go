package parser

import (
	"github.com/BurntSushi/toml"
	"github.com/sirupsen/logrus"
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
		logrus.Fatal(err)
	}
	return conf
}
