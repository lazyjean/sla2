package middleware

import (
	"net/http"
	"time"

	"github.com/lazyjean/sla2/pkg/logger"
	"go.uber.org/zap"
)

// HTTPLoggingMiddleware 创建一个用于记录 HTTP 请求日志的中间件
func HTTPLoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		log := logger.GetLogger(r.Context())

		// 记录请求开始
		log.Info("HTTP request started",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.String("remote_addr", r.RemoteAddr),
			zap.String("user_agent", r.UserAgent()),
		)

		// 创建一个响应记录器来捕获状态码
		rw := newResponseWriter(w)

		// 调用下一个处理器
		next.ServeHTTP(rw, r)

		// 计算处理时间
		duration := time.Since(startTime)

		// 记录请求结束
		log.Info("HTTP request completed",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.Int("status", rw.status),
			zap.Duration("duration", duration),
		)
	})
}

// responseWriter 是一个自定义的 ResponseWriter，用于捕获状态码
type responseWriter struct {
	http.ResponseWriter
	status int
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w, status: http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}
