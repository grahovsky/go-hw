package config

import (
	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Sender struct {
	Log struct {
		Level string `mapstructure:"level" env:"LOG_LEVEL"`
	} `mapstructure:"log"`
	Rmq RMQ `mapstructure:"rmq"`
}

var SenderSettings *Sender

func InitSenderSettings() {
	defaultSenderSettings := defaultSenderSettings()
	SenderSettings = &defaultSenderSettings

	pflag.String("loglevel", "INFO", "log level app")
	pflag.String("config", "./config/sender.yaml", "Path to Senderuration file")
	pflag.String("rmq_host", "0.0.0.0", "rmq hostname")
	pflag.String("rmq_port", "5672", "rmq port")

	pflag.Parse()

	viper.BindPFlags(pflag.CommandLine)

	viper.SetConfigFile(viper.Get("Sender").(string))
	viper.AutomaticEnv()
	viper.SetEnvPrefix("APP")

	if err := viper.ReadInConfig(); err != nil {
		logger.Error(err.Error())
	}

	if err := viper.Unmarshal(&SenderSettings); err != nil {
		logger.Error(err.Error())
	}

	envLogLevel := viper.Get("LOG_LEVEL")
	if envLogLevel != nil {
		SenderSettings.Log.Level = envLogLevel.(string)
	}
}

func defaultSenderSettings() Sender {
	return Sender{
		Log: struct {
			Level string "mapstructure:\"level\" env:\"LOG_LEVEL\""
		}{Level: "DEBUG"},
	}
}
