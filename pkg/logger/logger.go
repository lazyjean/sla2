package logger

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lazyjean/sla2/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ctxKey struct{}

var loggerCtxKey = ctxKey{}

// GetLogger 从上下文中获取 logger
func GetLogger(ctx context.Context) *zap.Logger {
	if ctx == nil {
		return Log
	}
	if logger, ok := ctx.Value(loggerCtxKey).(*zap.Logger); ok {
		return logger
	}
	return Log
}

// WithContext 将 logger 注入到上下文中
func WithContext(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, loggerCtxKey, logger)
}

// NewRequestLogger 创建带有请求 ID 的 logger
func NewRequestLogger() (*zap.Logger, string) {
	traceID := uuid.New().String()
	return Log.With(zap.String("trace_id", traceID)), traceID
}

// Log 全局logger实例
var Log *zap.Logger

// 添加一个全局变量来跟踪日志文件
var logFile *os.File

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
	// 如果已经有打开的日志文件，先关闭它
	if logFile != nil {
		if err := logFile.Close(); err != nil {
			return fmt.Errorf("关闭现有日志文件失败: %w", err)
		}
		logFile = nil
	}

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

	// 控制台输出
	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
	consoleOutput := zapcore.AddSync(os.Stdout)

	// 文件输出
	jsonEncoder := zapcore.NewJSONEncoder(encoderConfig)

	var cores []zapcore.Core
	cores = append(cores, zapcore.NewCore(consoleEncoder, consoleOutput, getLogLevel(cfg.Level)))

	// 只有在提供了有效的文件路径时才启用文件日志
	if cfg.FilePath != "" {
		// 确保日志目录存在
		logDir := filepath.Dir(cfg.FilePath)
		if err := os.MkdirAll(logDir, 0755); err != nil {
			return fmt.Errorf("创建日志目录失败: %w", err)
		}

		var err error
		logFile, err = os.OpenFile(cfg.FilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("打开日志文件失败: %w", err)
		}

		fileOutput := zapcore.AddSync(logFile)
		cores = append(cores, zapcore.NewCore(jsonEncoder, fileOutput, getLogLevel(cfg.Level)))
	}

	core := zapcore.NewTee(cores...)

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
