package logger

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
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
	return ctxzap.Extract(ctx)
}

// WithContext 将 logger 注入到上下文中
func WithContext(ctx context.Context, logger *zap.Logger) context.Context {
	return ctxzap.ToContext(ctx, logger)
}

// Log 全局logger实例
var Log *zap.Logger

// 添加一个全局变量来跟踪日志文件
var logFile *os.File

// 定义日志级别的颜色
const (
	colorRed     = "\033[31m"
	colorGreen   = "\033[32m"
	colorYellow  = "\033[33m"
	colorBlue    = "\033[34m"
	colorMagenta = "\033[35m"
	colorCyan    = "\033[36m"
	colorReset   = "\033[0m"
)

// 彩色日志级别编码器
func coloredLevelEncoder(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	var color string
	switch l {
	case zapcore.DebugLevel:
		color = colorBlue
	case zapcore.InfoLevel:
		color = colorGreen
	case zapcore.WarnLevel:
		color = colorYellow
	case zapcore.ErrorLevel:
		color = colorRed
	case zapcore.DPanicLevel, zapcore.PanicLevel, zapcore.FatalLevel:
		color = colorMagenta
	default:
		color = colorReset
	}
	enc.AppendString(color + l.String() + colorReset)
}

// 自定义时间编码器
func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("15:04:05.000"))
}

// 自定义调用者编码器
func customCallerEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	// 只显示文件名和行号
	enc.AppendString(fmt.Sprintf("%s:%d", filepath.Base(caller.File), caller.Line))
}

// InitBaseLogger 初始化基础日志配置
func InitBaseLogger() {
	var err error
	Log, err = zap.NewProduction()
	if err != nil {
		panic("failed to initialize base logger: " + err.Error())
	}
}

func init() {
	InitBaseLogger()
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

	// 根据配置选择控制台输出格式
	var consoleEncoder zapcore.Encoder
	if cfg.Format == "console" {
		consoleEncoder = zapcore.NewConsoleEncoder(encoderConfig)
	} else {
		consoleEncoder = zapcore.NewJSONEncoder(encoderConfig)
	}
	consoleOutput := zapcore.AddSync(os.Stdout)

	// 文件输出始终使用 JSON 格式
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

// 优化错误日志格式
func formatError(err error) string {
	if err == nil {
		return ""
	}
	return fmt.Sprintf("\n%s%s%s", colorRed, err.Error(), colorReset)
}

func NewConsoleLogger() *zap.Logger {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    coloredLevelEncoder,
		EncodeTime:     customTimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   customCallerEncoder,
	}

	// 使用控制台格式，更适合开发环境
	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
	consoleOutput := zapcore.AddSync(os.Stdout)

	// 创建 core
	core := zapcore.NewCore(consoleEncoder, consoleOutput, zapcore.DebugLevel)

	// 创建 logger
	logger := zap.New(core,
		zap.AddCaller(),                       // 添加调用者信息
		zap.AddStacktrace(zapcore.ErrorLevel), // Error 级别以上添加堆栈信息
		zap.AddCallerSkip(1),                  // 跳过一层调用栈
	)

	return logger
}

func NewAppFileLogger(cfg *config.LogConfig) *zap.Logger {
	if logFile != nil {
		if err := logFile.Close(); err != nil {
			log.Fatalf("关闭现有日志文件失败: %v", err)
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

	// 根据配置选择控制台输出格式
	var consoleEncoder zapcore.Encoder
	if cfg.Format == "console" {
		consoleEncoder = zapcore.NewConsoleEncoder(encoderConfig)
	} else {
		consoleEncoder = zapcore.NewJSONEncoder(encoderConfig)
	}
	consoleOutput := zapcore.AddSync(os.Stdout)

	// 文件输出始终使用 JSON 格式
	jsonEncoder := zapcore.NewJSONEncoder(encoderConfig)

	var cores []zapcore.Core
	cores = append(cores, zapcore.NewCore(consoleEncoder, consoleOutput, getLogLevel(cfg.Level)))

	// 只有在提供了有效的文件路径时才启用文件日志
	if cfg.FilePath != "" {
		// 确保日志目录存在
		logDir := filepath.Dir(cfg.FilePath)
		if err := os.MkdirAll(logDir, 0755); err != nil {
			log.Fatalf("创建日志目录失败: %v", err)
		}

		var err error
		logFile, err = os.OpenFile(cfg.FilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("打开日志文件失败: %v", err)
		}

		fileOutput := zapcore.AddSync(logFile)
		cores = append(cores, zapcore.NewCore(jsonEncoder, fileOutput, getLogLevel(cfg.Level)))
	}

	core := zapcore.NewTee(cores...)

	Log = zap.New(core,
		zap.AddCaller(),                       // 添加调用者信息
		zap.AddStacktrace(zapcore.ErrorLevel), // Error 级别以上添加堆栈信息
	)

	return Log
}

// v2
func NewAppLogger(cfg *config.LogConfig) *zap.Logger {
	// todo: 直接初始化, 不通过子函数调用
	// 开发环境使用控制台输出
	if cfg.Format == "console" {
		return NewConsoleLogger()
	}
	// 生产环境使用 JSON 格式输出到文件
	return NewAppFileLogger(cfg)
}
