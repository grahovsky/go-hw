package logger

import (
	"io"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	core   *zap.Logger
	level  string
	writer io.Writer
}

var myLog *Logger

const (
	ErrorLevel = "ERROR"
	WarnLevel  = "WARN"
	InfoLevel  = "INFO"
	DebugLevel = "DEBUG"
)

func init() {
	myLog = &Logger{core: zap.Must(zap.NewDevelopment()), level: InfoLevel, writer: os.Stdout}
}

func SetLogLevel(level string) {
	myLog.level = level
	initCore()
}

func SetWriter(writer io.Writer) {
	myLog.writer = writer
	initCore()
}

func initCore() {
	var zapLevel zapcore.Level

	switch myLog.level {
	case "WARN":
		zapLevel = zap.WarnLevel
	case "INFO":
		zapLevel = zap.InfoLevel
	case "DEBUG":
		zapLevel = zap.DebugLevel
	default:
		zapLevel = zap.ErrorLevel
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.AddSync(myLog.writer),
		zap.NewAtomicLevelAt(zapLevel))

	myLog.core = zap.New(core)
}

func GetLogger() *Logger {
	return myLog
}

func Error(msg string) {
	myLog.core.Error(msg)
}

func Warn(msg string) {
	myLog.core.Warn(msg)
}

func Info(msg string) {
	myLog.core.Info(msg)
}

func Debug(msg string) {
	myLog.core.Debug(msg)
}
