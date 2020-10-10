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

func GetLogLevel() (level log.Level) {
	var err error
	level, err = log.ParseLevel(Conf.Loglevel)
	if err != nil {
		level = log.ErrorLevel
	}
	return
}

type Config struct {
	Loglevel string `yaml:"loglevel"`
	Steam SteamConfig `yaml:"steam"`
	Hooks HookSecrets `yaml:"hooks"`
}

type SteamConfig struct {
	ApiKey string `yaml:"apikey"`
}

type HookSecrets struct {
	ExampleSenderSimple string `yaml:"exampleSenderSimple"`
	ExampleSenderEmbeds string `yaml:"exampleSenderEmbeds"`
	RustNewItems string `yaml:"rustNewItems"`
	GeekHackNews string `yaml:"geekHackNews"`
}
