package logger

import (
	"go.uber.org/zap"
)

type Logger struct {
	core  *zap.Logger
	level string
}

var myLog *Logger

func init() {
	myLog = &Logger{core: zap.Must(zap.NewDevelopment()), level: "INFO"}
}

func SetLogLevel(level string) {
	myLog.level = level
}

func GetLogger() *Logger {
	return myLog
}

func Error(msg string) {
	myLog.core.Error(msg)
}

func Warn(msg string) {
	if myLog.level == "WARN" || myLog.level == "INFO" || myLog.level == "DEBUG" {
		myLog.core.Warn(msg)
	}
}

func Info(msg string) {
	if myLog.level == "INFO" || myLog.level == "DEBUG" {
		myLog.core.Info(msg)
	}
}

func Debug(msg string) {
	if myLog.level == "DEBUG" {
		myLog.core.Debug(msg)
	}
}
