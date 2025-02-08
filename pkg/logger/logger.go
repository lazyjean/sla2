package logger

import (
	"os"
	"strings"
	"time"

	"github.com/lazyjean/sla2/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Log 全局logger实例
var Log *zap.Logger

// InitBaseLogger 初始化基础日志配置
func InitBaseLogger() {
	var err error
	Log, err = zap.NewProduction()
	if err != nil {
		panic("failed to initialize base logger: " + err.Error())
	}
}

// InitLogger 初始化完整的日志配置
func InitLogger(cfg *config.LogConfig) error {
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

	return nil
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

// GinLogger 实现 io.Writer 接口，用于 gin 框架的日志集成
type GinLogger struct {
	logger *zap.Logger
}

// Write 实现 io.Writer 接口
func (g *GinLogger) Write(p []byte) (n int, err error) {
	// 去除末尾的换行符
	msg := string(p)
	if len(msg) > 0 && msg[len(msg)-1] == '\n' {
		msg = msg[:len(msg)-1]
	}

	// 根据前缀选择日志级别
	switch {
	case strings.HasPrefix(msg, "[GIN-debug]"):
		msg = strings.TrimPrefix(msg, "[GIN-debug] ")
		g.logger.Debug(msg)
	case strings.HasPrefix(msg, "[GIN-info]"):
		msg = strings.TrimPrefix(msg, "[GIN-info] ")
		g.logger.Info(msg)
	case strings.HasPrefix(msg, "[GIN-warning]"):
		msg = strings.TrimPrefix(msg, "[GIN-warning] ")
		g.logger.Warn(msg)
	case strings.HasPrefix(msg, "[GIN-error]"):
		msg = strings.TrimPrefix(msg, "[GIN-error] ")
		g.logger.Error(msg)
	default:
		// 默认使用 Info 级别
		g.logger.Info(msg)
	}

	return len(p), nil
}

// NewGinLogger 创建一个新的 gin logger
func NewGinLogger() *GinLogger {
	return &GinLogger{
		logger: Log,
	}
}
