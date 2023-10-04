package config

import (
	"time"

	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Scheduler struct {
	Log struct {
		Level string `mapstructure:"level" env:"LOG_LEVEL"`
	} `mapstructure:"log"`
	Storage  Storage       `mapstructure:"storage"`
	Rmq      RMQ           `mapstructure:"rmq"`
	Schedule time.Duration `mapstructure:"schedule" env:"SHEDULE"`
}

var SchedulerSettings *Scheduler

func InitSchedulerSettings() {
	defaultSchedulerSettings := defaultSchedulerSettings()
	SchedulerSettings = &defaultSchedulerSettings

	pflag.String("loglevel", "INFO", "log level app")
	pflag.String("config", "./configs/scheduler.yaml", "path to sheduler config file")
	pflag.String("rmq_host", "0.0.0.0", "rmq hostname")
	pflag.String("rmq_port", "5672", "server port")

	pflag.Parse()

	viper.BindPFlags(pflag.CommandLine)

	viper.SetConfigFile(viper.Get("config").(string))
	viper.AutomaticEnv()
	viper.SetEnvPrefix("APP")

	if err := viper.ReadInConfig(); err != nil {
		logger.Error(err.Error())
	}

	if err := viper.Unmarshal(&SchedulerSettings); err != nil {
		logger.Error(err.Error())
	}

	envLogLevel := viper.Get("LOG_LEVEL")
	if envLogLevel != nil {
		SchedulerSettings.Log.Level = envLogLevel.(string)
	}
}

func defaultSchedulerSettings() Scheduler {
	return Scheduler{
		Log: struct {
			Level string "mapstructure:\"level\" env:\"LOG_LEVEL\""
		}{Level: "DEBUG"},
	}
}
