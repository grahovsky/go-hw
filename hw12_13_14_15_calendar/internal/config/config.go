package config

import (
	"os"

	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/logger"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Log struct {
		Level string `yaml:"level"`
	} `yaml:"log"`
	DebugMessage string `yaml:"debugMessage"`
	Some         string `yaml:"some"`
}

var Settings *Config

func Read(configPath string) {
	yamlFile, err := os.ReadFile(configPath)
	if err != nil {
		logger.Fatal(err.Error())
	}
	err = yaml.Unmarshal(yamlFile, &Settings)
	if err != nil {
		logger.Fatal(err.Error())
	}
}
