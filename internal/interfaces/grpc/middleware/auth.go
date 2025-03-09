package middleware

import (
	"context"
	"strings"

	"github.com/lazyjean/sla2/internal/domain/security"
	"github.com/lazyjean/sla2/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// 不需要认证的方法列表
var noAuthMethods = map[string]bool{
	"/proto.v1.AdminService/CheckSystemStatus": true,
	"/proto.v1.AdminService/InitializeSystem":  true,
	"/proto.v1.AdminService/AdminLogin":        true,
	"/proto.v1.AdminService/RefreshToken":      true,
	"/proto.v1.UserService/Register":           true,
	"/proto.v1.UserService/Login":              true,
	"/proto.v1.UserService/RefreshToken":       true,
	"/proto.v1.UserService/AppleLogin":         true,
	"/proto.v1.UserService/ResetPassword":      true,
}

// UnaryServerInterceptor 一元 RPC 认证中间件
func UnaryServerInterceptor(tokenService security.TokenService) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		log := logger.GetLogger(ctx)

		// 检查是否需要认证
		if noAuthMethods[info.FullMethod] {
			return handler(ctx, req)
		}

		// 从上下文中获取 token
		token, err := extractToken(ctx)
		if err != nil {
			log.Error("Failed to extract token", zap.Error(err))
			return nil, status.Error(codes.Unauthenticated, "未授权")
		}

		// 验证 token
		userID, roles, err := tokenService.ValidateToken(token)
		if err != nil {
			log.Error("Failed to validate token", zap.Error(err))
			return nil, status.Error(codes.Unauthenticated, "未授权")
		}

		// 将用户信息添加到上下文
		newCtx := context.WithValue(ctx, "user_id", userID)
		newCtx = context.WithValue(newCtx, "roles", roles)

		return handler(newCtx, req)
	}
}

// StreamServerInterceptor 流式 RPC 认证中间件
func StreamServerInterceptor(tokenService security.TokenService) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		log := logger.GetLogger(ss.Context())

		// 检查是否需要认证
		if noAuthMethods[info.FullMethod] {
			return handler(srv, ss)
		}

		// 从上下文中获取 token
		token, err := extractToken(ss.Context())
		if err != nil {
			log.Error("Failed to extract token", zap.Error(err))
			return status.Error(codes.Unauthenticated, "未授权")
		}

		// 验证 token
		userID, roles, err := tokenService.ValidateToken(token)
		if err != nil {
			log.Error("Failed to validate token", zap.Error(err))
			return status.Error(codes.Unauthenticated, "未授权")
		}

		// 将用户信息添加到上下文
		newCtx := context.WithValue(ss.Context(), "user_id", userID)
		newCtx = context.WithValue(newCtx, "roles", roles)

		// 包装 ServerStream 以使用新的上下文
		wrappedStream := &wrappedServerStream{
			ServerStream: ss,
			ctx:          newCtx,
		}

		return handler(srv, wrappedStream)
	}
}

// extractToken 从上下文中提取 token
func extractToken(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "未授权")
	}

	// 从 Authorization 头中获取 token
	values := md.Get("authorization")
	if len(values) == 0 {
		return "", status.Error(codes.Unauthenticated, "未授权")
	}

	token := values[0]
	if !strings.HasPrefix(token, "Bearer ") {
		return "", status.Error(codes.Unauthenticated, "无效的认证方式")
	}

	return strings.TrimPrefix(token, "Bearer "), nil
}

// wrappedServerStream 包装 ServerStream 以支持修改上下文
type wrappedServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

// Context 返回包装的上下文
func (w *wrappedServerStream) Context() context.Context {
	return w.ctx
}
