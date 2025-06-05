package middleware

import (
	"context"
	"fmt"
	"strings"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/lazyjean/sla2/internal/application/service"
	"github.com/lazyjean/sla2/internal/domain/security"
	"github.com/lazyjean/sla2/pkg/logger"
)

// 不需要认证的方法列表
var noAuthMethods = map[string]bool{
	// 用户服务
	"/proto.v1.UserService/Login":         true,
	"/proto.v1.UserService/Register":      true,
	"/proto.v1.UserService/AppleLogin":    true,
	"/proto.v1.UserService/ResetPassword": true,

	// 管理员服务
	"/proto.v1.AdminService/IsSystemInitialized": true,
	"/proto.v1.AdminService/InitializeSystem":    true,
	"/proto.v1.AdminService/AdminLogin":          true,

	// gRPC 反射服务
	"/grpc.reflection.v1.ServerReflection/ServerReflectionInfo": true,

	// 健康检查服务
	"/grpc.health.v1.Health/Check": true,
	"/grpc.health.v1.Health/Watch": true,
}

func AuthFunc(tokenService security.TokenService) auth.AuthFunc {
	return func(ctx context.Context) (context.Context, error) {
		return authenticate(ctx, tokenService)
	}
}

// authenticate 通用的认证逻辑
func authenticate(ctx context.Context, tokenService security.TokenService) (context.Context, error) {
	// 从上下文中获取 token
	token, err := extractToken(ctx)
	if err != nil {
		log := logger.GetLogger(ctx)
		log.Error("Failed to extract token", zap.Error(err))
		return nil, status.Error(codes.Unauthenticated, "未授权")
	}

	// 验证 token
	userID, roles, err := tokenService.ValidateToken(token)
	if err != nil {
		log := logger.GetLogger(ctx)
		log.Error("Failed to validate token", zap.Error(err))
		return nil, status.Error(codes.Unauthenticated, "未授权")
	}

	// 将用户信息添加到上下文
	newCtx := service.WithUserID(ctx, userID)
	newCtx = service.WithRoles(newCtx, roles)

	return newCtx, nil
}

func RequireAuth(ctx context.Context, callMeta interceptors.CallMeta) bool {
	if noAuthMethods[callMeta.FullMethod()] {
		return false
	}
	return true
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
		// 尝试从其他可能的头部获取 token
		headers := []string{
			"authorization",  // Bearer token
			"jwt",            // JWT token
			"access-token",   // Access token
			"x-access-token", // X-Access-Token
		}

		for _, header := range headers {
			values = md.Get(header)
			if len(values) > 0 {
				// 如果是 authorization 头，需要提取 Bearer token
				if header == "authorization" {
					authValue := values[0]
					if strings.HasPrefix(authValue, "Bearer ") {
						values[0] = strings.TrimPrefix(authValue, "Bearer ")
					}
				}
				break
			}
		}
	}

	if len(values) == 0 {
		log := logger.GetLogger(ctx)
		log.Error("Failed to extract token", zap.String("metadata", fmt.Sprintf("%v", md)))
		return "", status.Error(codes.Unauthenticated, "未授权")
	}

	return values[0], nil
}
