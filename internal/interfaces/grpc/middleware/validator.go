package middleware

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Validator 定义了验证接口
type Validator interface {
	Validate() error
}

// ValidatorInterceptor 创建一个验证拦截器
func ValidatorInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// 检查请求是否实现了Validator接口
		if v, ok := req.(Validator); ok {
			// 调用验证方法
			if err := v.Validate(); err != nil {
				return nil, status.Error(codes.InvalidArgument, err.Error())
			}
		}

		// 如果验证通过或者没有实现Validator接口，继续处理请求
		return handler(ctx, req)
	}
}

// StreamValidatorInterceptor 创建一个流验证拦截器
func StreamValidatorInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		// 创建一个包装的ServerStream，它会验证所有接收到的消息
		wrapper := &recvWrapper{
			ServerStream: ss,
		}
		return handler(srv, wrapper)
	}
}

// recvWrapper 包装ServerStream以验证接收到的消息
type recvWrapper struct {
	grpc.ServerStream
}

// RecvMsg 重写RecvMsg方法以添加验证
func (s *recvWrapper) RecvMsg(m interface{}) error {
	if err := s.ServerStream.RecvMsg(m); err != nil {
		return err
	}

	// 检查消息是否实现了Validator接口
	if v, ok := m.(Validator); ok {
		if err := v.Validate(); err != nil {
			return status.Error(codes.InvalidArgument, err.Error())
		}
	}

	return nil
}
