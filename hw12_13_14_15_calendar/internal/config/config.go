package config

import (
	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Config struct {
	Log struct {
		Level string `mapstructure:"level" env:"LOG_LEVEL"`
	} `mapstructure:"log"`
	DebugMessage string `mapstructure:"debugMessage"`
	Some         string `mapstructure:"some"`
	PrintVersion bool
}

var Settings *Config

func init() {
	defaultSettings := defaultSettings()
	Settings = &defaultSettings

	versionFlag := pflag.Bool("version", false, "version app")
	pflag.String("loglevel", "INFO", "log level app")
	pflag.String("config", "configs/config.yaml", "Path to configuration file")
	pflag.Parse()

	if *versionFlag {
		defaultSettings.PrintVersion = *versionFlag
		return
	}

	viper.BindPFlags(pflag.CommandLine)

	viper.SetConfigFile(viper.Get("config").(string))
	viper.AutomaticEnv()
	viper.SetEnvPrefix("APP")

	if err := viper.ReadInConfig(); err != nil {
		logger.Fatal(err.Error())
	}

	if err := viper.Unmarshal(&Settings); err != nil {
		logger.Fatal(err.Error())
	}

	envLogLevel := viper.Get("LOG_LEVEL")
	if envLogLevel != nil {
		Settings.Log.Level = envLogLevel.(string)
	}
}

func defaultSettings() Config {
	return Config{
		Log: struct {
			Level string "mapstructure:\"level\" env:\"LOG_LEVEL\""
		}{Level: "DEBUG"},
		PrintVersion: false,
	}
}
