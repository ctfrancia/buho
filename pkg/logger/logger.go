package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func New(env string) *zap.Logger {
	var core zapcore.Core
	err := os.MkdirAll("./internal/logs", os.ModePerm)
	if err != nil {
		panic(err)
	}
	isDev := env == "dev" || env == "test"

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	logLevel := zap.InfoLevel
	if isDev {
		// In development: Use a console encoder and write to stderr
		logLevel = zap.DebugLevel
		consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
		core = zapcore.NewCore(
			consoleEncoder,
			zapcore.AddSync(os.Stderr),
			logLevel,
		)
	} else {
		// In production/other environments: Use JSON encoder and log rotation
		logRotator := &lumberjack.Logger{
			Filename:   "./internal/logs/buho.log",
			MaxSize:    10, // megabytes
			MaxBackups: 3,
			MaxAge:     28, // days
			Compress:   true,
		}

		core = zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.AddSync(logRotator),
			logLevel,
		)
	}

	zapLogger := zap.New(core, zap.AddCaller())
	return zapLogger
}
