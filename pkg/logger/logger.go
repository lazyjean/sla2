package logger

import (
	"context"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/lazyjean/sla2/config"
)

type loggerKey struct{}

var zapLogger *zap.Logger

// GetLogger 从 Context 获取 Logger
func GetLogger(ctx context.Context) *zap.Logger {
	if logger, ok := ctx.Value(loggerKey{}).(*zap.Logger); ok {
		return logger
	}
	return zapLogger
}

// WithContext 将 logger 注入到上下文中
func WithContext(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, loggerKey{}, logger)
}

func NewAppLogger(cfg *config.LogConfig) *zap.Logger {
	var logger *zap.Logger
	if cfg.Production {
		logger, _ = zap.NewProduction()
	} else {
		config := zap.NewDevelopmentConfig()

		// 自定义编码器
		config.EncoderConfig.EncodeLevel = func(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
			var color string
			switch level {
			case zapcore.DebugLevel:
				color = "\x1b[36m" // 青色
			case zapcore.InfoLevel:
				color = "\x1b[32m" // 绿色
			case zapcore.WarnLevel:
				color = "\x1b[33m" // 黄色
			case zapcore.ErrorLevel, zapcore.DPanicLevel, zapcore.PanicLevel, zapcore.FatalLevel:
				color = "\x1b[31m" // 红色
			default:
				color = "\x1b[0m" // 默认
			}
			enc.AppendString(color + level.CapitalString() + "\x1b[0m")
		}

		// 自定义时间格式
		config.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
		}

		// 创建自定义编码器
		encoder := zapcore.NewConsoleEncoder(config.EncoderConfig)
		core := zapcore.NewCore(
			encoder,
			zapcore.AddSync(os.Stdout),
			zapcore.DebugLevel,
		)

		logger = zap.New(core)
	}
	zapLogger = logger
	return logger
}
