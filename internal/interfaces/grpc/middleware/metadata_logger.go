package middleware

import (
	"context"

	"github.com/lazyjean/sla2/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// MetadataLoggerUnaryInterceptor 创建打印metadata的一元拦截器
func MetadataLoggerUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		log := logger.GetLogger(ctx)

		// 从上下文中获取metadata
		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			log.Info("请求metadata",
				zap.String("method", info.FullMethod),
				zap.Any("metadata", md),
			)
		} else {
			log.Info("请求无metadata", zap.String("method", info.FullMethod))
		}

		// 继续处理请求
		return handler(ctx, req)
	}
}

// MetadataLoggerStreamInterceptor 创建打印metadata的流式拦截器
func MetadataLoggerStreamInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		log := logger.GetLogger(ss.Context())

		// 从上下文中获取metadata
		md, ok := metadata.FromIncomingContext(ss.Context())
		if ok {
			log.Info("流式请求metadata",
				zap.String("method", info.FullMethod),
				zap.Any("metadata", md),
			)
		} else {
			log.Info("流式请求无metadata", zap.String("method", info.FullMethod))
		}

		// 继续处理请求
		return handler(srv, ss)
	}
}
