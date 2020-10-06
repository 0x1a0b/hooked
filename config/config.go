package config

import (
	"github.com/sherifabdlnaby/configuro"
	log "github.com/sirupsen/logrus"
)

func Read() error {

	c := Config{}

	config, err := configuro.NewConfig()

	if err != nil {
		log.Fatalf("problem while creating config loader: %v", err)
	}

	err = config.Load(&c)

	if err != nil {
		log.Fatalf("problem while loading config: %v", err)
	}

	log.WithFields(log.Fields{
		"config": c,
	}).Debug("Initialized lykill with config")

	Conf = &c

	return err
}

var (
	Conf *Config
)

func GetConf() *Config {
	return Conf
}

type Config struct {
	Discord DiscordConfig `yaml:"discord"`
	Steam SteamConfig `yaml:"steam"`
}

type DiscordConfig struct {
	Url string `yaml:"url"`
}

type SteamConfig struct {
	ApiKey string `yaml:"apikey"`
}