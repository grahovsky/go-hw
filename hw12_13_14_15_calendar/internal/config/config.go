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
	Storage struct {
		Type string `mapstructure:"type" env:"STORAGE_TYPE"`
	} `mapstructure:"storage"`
	DB struct {
		Name     string `mapstructure:"name" env:"DB_NAME"`
		Host     string `mapstructure:"host" env:"DB_HOST"`
		User     string `mapstructure:"user" env:"DB_USER"`
		Password string `mapstructure:"password" env:"DB_PASSWORD"`
	} `mapstructure:"db"`
	Server struct {
		Host string `mapstructure:"host" env:"SRV_HOST"`
		Port string `mapstructure:"port" env:"SRV_PORT"`
	} `mapstructure:"server"`
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
	pflag.String("config", "./configs/config.yaml", "Path to configuration file")
	pflag.String("server_host", "127.0.0.1", "server hostname")
	pflag.String("server_port", "8080", "server port")

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
		logger.Error(err.Error())
	}

	if err := viper.Unmarshal(&Settings); err != nil {
		logger.Error(err.Error())
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
