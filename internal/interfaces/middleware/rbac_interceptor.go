package middleware

import (
	"context"
	"fmt"

	"github.com/lazyjean/sla2/internal/domain/entity"
	"github.com/lazyjean/sla2/internal/domain/security"
	"github.com/lazyjean/sla2/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// RBACInterceptor gRPC权限拦截器
type RBACInterceptor struct {
	permissionHelper *security.PermissionHelper
	// 方法权限映射表，格式: {"service.method": {"resource": "xxx", "action": "xxx"}}
	methodPermissionMap map[string]map[string]string
	// 白名单方法，不需要权限验证
	whitelistMethods map[string]bool
}

// NewRBACInterceptor 创建新的gRPC权限拦截器
func NewRBACInterceptor(permissionHelper *security.PermissionHelper) *RBACInterceptor {
	return &RBACInterceptor{
		permissionHelper:    permissionHelper,
		methodPermissionMap: make(map[string]map[string]string),
		whitelistMethods:    make(map[string]bool),
	}
}

// RegisterMethodPermission 注册方法对应的权限
func (r *RBACInterceptor) RegisterMethodPermission(fullMethod, resource, action string) {
	r.methodPermissionMap[fullMethod] = map[string]string{
		"resource": resource,
		"action":   action,
	}
	logger.Log.Debug("Registered method permission",
		zap.String("method", fullMethod),
		zap.String("resource", resource),
		zap.String("action", action))
}

// AddToWhitelist 添加方法到白名单
func (r *RBACInterceptor) AddToWhitelist(fullMethod string) {
	r.whitelistMethods[fullMethod] = true
	logger.Log.Debug("Added method to whitelist", zap.String("method", fullMethod))
}

// UnaryServerInterceptor 一元RPC拦截器
func (r *RBACInterceptor) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		log := logger.GetLogger(ctx)

		// 检查方法是否在白名单中
		if r.whitelistMethods[info.FullMethod] {
			log.Debug("Method is in whitelist, skipping permission check", zap.String("method", info.FullMethod))
			return handler(ctx, req)
		}

		// 获取用户ID
		userID, err := extractUserID(ctx)
		if err != nil {
			log.Error("Failed to extract user ID",
				zap.String("method", info.FullMethod),
				zap.Error(err))
			return nil, status.Errorf(codes.Unauthenticated, "invalid authentication: %v", err)
		}

		// 获取方法对应的权限信息
		permInfo, exists := r.methodPermissionMap[info.FullMethod]
		if !exists {
			log.Warn("No permission defined for method", zap.String("method", info.FullMethod))
			return handler(ctx, req)
		}

		log.Debug("Checking permission",
			zap.String("method", info.FullMethod),
			zap.Uint("user_id", userID),
			zap.String("resource", permInfo["resource"]),
			zap.String("action", permInfo["action"]))

		// 检查权限
		hasPermission, err := r.permissionHelper.CheckUserPermission(ctx, entity.UID(userID), permInfo["resource"], permInfo["action"])
		if err != nil {
			log.Error("Error checking permission",
				zap.Error(err),
				zap.String("method", info.FullMethod),
				zap.Uint("user_id", userID))
			return nil, status.Errorf(codes.Internal, "error checking permission: %v", err)
		}

		if !hasPermission {
			log.Warn("Permission denied",
				zap.String("method", info.FullMethod),
				zap.Uint("user_id", userID),
				zap.String("resource", permInfo["resource"]),
				zap.String("action", permInfo["action"]))
			return nil, status.Error(codes.PermissionDenied, "permission denied")
		}

		log.Debug("Permission granted",
			zap.String("method", info.FullMethod),
			zap.Uint("user_id", userID),
			zap.String("resource", permInfo["resource"]),
			zap.String("action", permInfo["action"]))

		// 权限验证通过，继续处理请求
		return handler(ctx, req)
	}
}

// StreamServerInterceptor 流式RPC拦截器
func (r *RBACInterceptor) StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := ss.Context()
		log := logger.GetLogger(ctx)

		// 检查方法是否在白名单中
		if r.whitelistMethods[info.FullMethod] {
			log.Debug("Stream method is in whitelist, skipping permission check", zap.String("method", info.FullMethod))
			return handler(srv, ss)
		}

		// 获取用户ID
		userID, err := extractUserID(ctx)
		if err != nil {
			log.Error("Failed to extract user ID for stream",
				zap.String("method", info.FullMethod),
				zap.Error(err))
			return status.Errorf(codes.Unauthenticated, "invalid authentication: %v", err)
		}

		// 获取方法对应的权限信息
		permInfo, exists := r.methodPermissionMap[info.FullMethod]
		if !exists {
			log.Warn("No permission defined for stream method", zap.String("method", info.FullMethod))
			return handler(srv, ss)
		}

		log.Debug("Checking permission for stream",
			zap.String("method", info.FullMethod),
			zap.Uint("user_id", userID),
			zap.String("resource", permInfo["resource"]),
			zap.String("action", permInfo["action"]))

		// 检查权限
		hasPermission, err := r.permissionHelper.CheckUserPermission(ctx, entity.UID(userID), permInfo["resource"], permInfo["action"])
		if err != nil {
			log.Error("Error checking permission for stream",
				zap.Error(err),
				zap.String("method", info.FullMethod),
				zap.Uint("user_id", userID))
			return status.Errorf(codes.Internal, "error checking permission: %v", err)
		}

		if !hasPermission {
			log.Warn("Permission denied for stream",
				zap.String("method", info.FullMethod),
				zap.Uint("user_id", userID),
				zap.String("resource", permInfo["resource"]),
				zap.String("action", permInfo["action"]))
			return status.Error(codes.PermissionDenied, "permission denied")
		}

		log.Debug("Permission granted for stream",
			zap.String("method", info.FullMethod),
			zap.Uint("user_id", userID),
			zap.String("resource", permInfo["resource"]),
			zap.String("action", permInfo["action"]))

		// 权限验证通过，继续处理请求
		return handler(srv, ss)
	}
}

// extractUserID 从上下文中提取用户ID
func extractUserID(ctx context.Context) (uint, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		logger.Log.Debug("No metadata found in context")
		return 0, fmt.Errorf("no metadata in context")
	}

	// 尝试从不同的元数据字段获取用户ID
	// 1. 首先尝试从 user_id 字段获取
	userIDValues := md.Get("user_id")
	if len(userIDValues) > 0 {
		return parseUserID(userIDValues[0])
	}

	// 2. 尝试从 authorization 头中获取
	authValues := md.Get("authorization")
	if len(authValues) > 0 {
		// 这里假设 JWT 中包含了用户ID信息
		// 实际中应该使用 tokenService 进行验证并提取用户ID
		logger.Log.Debug("Found authorization header, but no direct user_id")
		return 0, fmt.Errorf("authorization header found, but user_id extraction not implemented")
	}

	// 3. 尝试从 x-user-id 头中获取
	xUserIDValues := md.Get("x-user-id")
	if len(xUserIDValues) > 0 {
		return parseUserID(xUserIDValues[0])
	}

	logger.Log.Debug("User ID not found in any metadata field")
	return 0, fmt.Errorf("user_id not found in metadata")
}

// parseUserID 解析用户ID字符串
func parseUserID(userIDStr string) (uint, error) {
	var userID uint
	_, err := fmt.Sscanf(userIDStr, "%d", &userID)
	if err != nil {
		logger.Log.Debug("Invalid user ID format", zap.String("user_id", userIDStr))
		return 0, fmt.Errorf("invalid user_id format: %s", userIDStr)
	}
	return userID, nil
}

// 包装ServerStream以传递修改后的上下文
type wrappedServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (w *wrappedServerStream) Context() context.Context {
	return w.ctx
}
