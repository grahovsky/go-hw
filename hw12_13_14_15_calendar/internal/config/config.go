package config

import (
	"os"

	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/logger"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Log  LoggerConf `yaml: "log"`
	Some string     `yaml: "some"`
}

type LoggerConf struct {
	Level string `yaml: "level"`
}

var Settings *Config

func New(configPath string) Config {
	yamlFile, err := os.ReadFile(configPath)
	if err != nil {
		logger.DefaultLog.Fatal(err.Error())
	}
	err = yaml.Unmarshal(yamlFile, &Settings)
	if err != nil {
		logger.DefaultLog.Fatal(err.Error())
	}

	return *Settings
}
