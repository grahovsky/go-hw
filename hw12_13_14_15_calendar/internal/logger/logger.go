package logger

import (
	"go.uber.org/zap"
)

type Logger struct {
	core  *zap.Logger
	level string
}

var myLog *Logger

const (
	ErrorLevel = "ERROR"
	WarnLevel  = "WARN"
	InfoLevel  = "INFO"
	DebugLevel = "DEBUG"
)

func init() {
	myLog = &Logger{core: zap.Must(zap.NewDevelopment()), level: InfoLevel}
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
	if myLog.level == WarnLevel || myLog.level == InfoLevel || myLog.level == DebugLevel {
		myLog.core.Warn(msg)
	}
}

func Info(msg string) {
	if myLog.level == InfoLevel || myLog.level == DebugLevel {
		myLog.core.Info(msg)
	}
}

func Debug(msg string) {
	if myLog.level == DebugLevel {
		myLog.core.Debug(msg)
	}
}
