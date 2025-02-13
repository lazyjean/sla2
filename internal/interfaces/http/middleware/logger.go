package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/lazyjean/sla2/pkg/logger"
)

// LoggerMiddleware 日志中间件
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// start := time.Now()
		// path := c.Request.URL.Path
		// raw := c.Request.URL.RawQuery

		// 为每个请求创建带有 trace ID 的 logger
		reqLogger, traceID := logger.NewRequestLogger()

		// 将 logger 注入到上下文中
		ctx := logger.WithContext(c.Request.Context(), reqLogger)
		c.Request = c.Request.WithContext(ctx)

		// 设置 trace ID 到响应头
		c.Header("X-Trace-ID", traceID)

		// 处理请求
		c.Next()

		// 记录请求日志
		// if raw != "" {
		// 	path = path + "?" + raw
		// }

		// reqLogger.Info("http request",
		// 	zap.String("method", c.Request.Method),
		// 	zap.String("path", path),
		// 	zap.Int("status", c.Writer.Status()),
		// 	zap.Duration("latency", time.Since(start)),
		// 	zap.String("client_ip", c.ClientIP()),
		// 	zap.String("user_agent", c.Request.UserAgent()),
		// )
	}
}
