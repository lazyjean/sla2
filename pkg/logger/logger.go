package logger

import (
	"os"
	"time"

	"github.com/lazyjean/sla2/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger

func InitLogger(cfg *config.LogConfig) {
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
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	encoder := zapcore.NewJSONEncoder(encoderConfig)
	writeSyncer := zapcore.AddSync(os.Stdout)
	logLevel := getLogLevel(cfg.Level)

	core := zapcore.NewCore(encoder, writeSyncer, logLevel)

	Log = zap.New(core,
		zap.AddCaller(),                       // 添加调用者信息
		zap.AddStacktrace(zapcore.ErrorLevel), // Error 级别以上添加堆栈信息
	)
}

func getLogLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

func LogRequest(method, path string, latency time.Duration, status int) {
	Log.Info("HTTP Request",
		zap.String("method", method),
		zap.String("path", path),
		zap.Duration("latency", latency),
		zap.Int("status", status),
	)
}

func LogError(err error, msg string, fields ...zap.Field) {
	Log.Error(msg,
		append(fields, zap.Error(err))...,
	)
}
