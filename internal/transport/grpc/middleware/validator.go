package middleware

import (
	"context"

	"github.com/lazyjean/sla2/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// Validator 定义了验证接口
type Validator interface {
	Validate() error
}

// ValidatorInterceptor 创建一个验证拦截器
func ValidatorInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Preserve incoming metadata
		md, _ := metadata.FromIncomingContext(ctx)
		if md != nil {
			ctx = metadata.NewIncomingContext(ctx, md)
		}
		log := logger.GetLogger(ctx)
		// 检查请求是否实现了Validator接口
		if v, ok := req.(Validator); ok {
			// 调用验证方法
			if err := v.Validate(); err != nil {
				// 序列化请求体为JSON
				if msg, ok := req.(proto.Message); ok {
					jsonData, marshalErr := protojson.Marshal(msg)
					if marshalErr != nil {
						log.Error("marshal request failed", zap.Error(marshalErr))
					} else {
						// 截断过长的JSON数据
						jsonStr := string(jsonData)
						if len(jsonStr) > 1000 {
							jsonStr = jsonStr[:1000] + "...(truncated)"
						}
						log.Error("validate failed",
							zap.String("error", err.Error()),
							zap.String("request_body", jsonStr))
					}
				} else {
					log.Error("validate failed (non-proto request)",
						zap.String("error", err.Error()),
						zap.Any("request", req))
				}
				return nil, status.Error(codes.InvalidArgument, err.Error())
			}
		}

		// 如果验证通过或者没有实现Validator接口，继续处理请求
		return handler(ctx, req)
	}
}
