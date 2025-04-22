package middleware

import (
	"context"
	"fmt"
	"strings"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/lazyjean/sla2/internal/application/service"
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
	"/proto.v1.UserService/Login":                               true,
	"/proto.v1.UserService/Register":                            true,
	"/proto.v1.UserService/AppleLogin":                          true,
	"/proto.v1.UserService/ResetPassword":                       true,
	"/proto.v1.AdminService/CheckSystemStatus":                  true,
	"/proto.v1.AdminService/InitializeSystem":                   true,
	"/proto.v1.AdminService/AdminLogin":                         true,
	"/grpc.reflection.v1.ServerReflection/ServerReflectionInfo": true,
	"/grpc.health.v1.Health/Check":                              true,
	"/grpc.health.v1.Health/Watch":                              true,
}

// UnaryServerInterceptor 一元 RPC 认证中间件
func UnaryServerInterceptor(tokenService security.TokenService) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		log := logger.GetLogger(ctx)

		// 检查是否需要认证
		if noAuthMethods[info.FullMethod] || isPublicService(info.FullMethod) {
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
		newCtx := service.WithUserID(ctx, userID)
		newCtx = service.WithRoles(newCtx, roles)

		return handler(newCtx, req)
	}
}

// extractToken 从上下文中提取 token
func extractToken(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "未授权(从grpc metadata中获取token失败)")
	}

	// 从 Authorization 头中获取 token
	values := md.Get(MDHeaderAccessToken)
	if len(values) == 0 {
		return "", status.Error(codes.Unauthenticated, "未授权(从Authorization头中获取token失败)")
	}

	return values[0], nil
}

// StreamServerInterceptor 流式 RPC 认证中间件
func StreamServerInterceptor(tokenService security.TokenService) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		// 跳过认证的服务
		if isPublicService(info.FullMethod) {
			return handler(srv, ss)
		}

		// 从上下文中获取元数据
		md, ok := metadata.FromIncomingContext(ss.Context())
		if !ok {
			return status.Errorf(codes.Unauthenticated, "未授权")
		}

		fmt.Println("Metadata", md)

		// 从元数据中获取 token
		token := extractTokenFromMD(md)
		if token == "" {
			return status.Errorf(codes.Unauthenticated, "未授权")
		}

		// 验证 token
		userID, roles, err := tokenService.ValidateToken(token)
		if err != nil {
			logger.GetLogger(ss.Context()).Error("Failed to validate token", zap.Error(err))
			return status.Errorf(codes.Unauthenticated, "未授权")
		}

		// 将用户信息添加到上下文中
		newCtx := service.WithUserID(ss.Context(), userID)
		newCtx = service.WithRoles(newCtx, roles)
		wrapped := grpc_middleware.WrapServerStream(ss)
		wrapped.WrappedContext = newCtx

		return handler(srv, wrapped)
	}
}

// isPublicService 检查是否是公开的服务
func isPublicService(fullMethod string) bool {
	// 跳过反射服务
	if strings.HasPrefix(fullMethod, "/grpc.reflection.v1.ServerReflection/") {
		return true
	}

	// 跳过健康检查服务
	if strings.HasPrefix(fullMethod, "/grpc.health.v1.Health/") {
		return true
	}

	// 跳过初始化和系统状态检查
	if fullMethod == "/proto.v1.AdminService/CheckSystemStatus" ||
		fullMethod == "/proto.v1.AdminService/InitializeSystem" {
		return true
	}

	return false
}

// extractTokenFromMD 从元数据中提取 token
func extractTokenFromMD(md metadata.MD) string {
	// 从 Authorization header 获取 token
	if values := md.Get("authorization"); len(values) > 0 {
		auth := values[0]
		if strings.HasPrefix(auth, "Bearer ") {
			return auth[7:]
		}
		return auth
	}

	return ""
}
