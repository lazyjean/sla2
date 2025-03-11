package middleware

import (
	"context"
	"time"

	"github.com/lazyjean/sla2/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// LoggingUnaryServerInterceptor 创建一个用于记录请求日志的拦截器
func LoggingUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		log := logger.GetLogger(ctx)
		startTime := time.Now()

		// 记录请求开始
		log.Info("gRPC request started",
			zap.String("method", info.FullMethod),
			zap.Any("request", req),
		)

		// 调用处理器
		resp, err := handler(ctx, req)

		// 计算处理时间
		duration := time.Since(startTime)

		// 获取状态码
		statusCode := codes.OK
		if err != nil {
			if st, ok := status.FromError(err); ok {
				statusCode = st.Code()
			} else {
				statusCode = codes.Internal
			}
		}

		// 记录请求结束
		log.Info("gRPC request completed",
			zap.String("method", info.FullMethod),
			zap.String("status", statusCode.String()),
			zap.Duration("duration", duration),
			zap.Error(err),
		)

		return resp, err
	}
}
