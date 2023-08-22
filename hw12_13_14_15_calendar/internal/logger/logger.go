package logger

import (
	"go.uber.org/zap"
)

type Logger struct {
	core  *zap.Logger
	LEVEL string
}

var DefaultLog *Logger

func init() {
	DefaultLog = &Logger{core: zap.Must(zap.NewDevelopment()), LEVEL: "INFO"}
	DefaultLog.Info("inited logger")
}

func New(level string) *Logger {
	DefaultLog = &Logger{core: zap.Must(zap.NewDevelopment()), LEVEL: level}
	return DefaultLog
}

func (l Logger) Info(msg string) {
	l.core.Info(msg)
}

func (l Logger) Debug(msg string) {
	if l.LEVEL == "DEBUG" {
		l.core.Debug(msg)
	}
}

func (l Logger) Error(msg string) {
	l.core.Error(msg)
}

func (l Logger) Fatal(msg string) {
	l.core.Fatal(msg)
}

// TODO
