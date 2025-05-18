package middleware

import (
	"context"

	"github.com/lazyjean/sla2/internal/application/service"
	"github.com/lazyjean/sla2/internal/domain/security"
	"github.com/lazyjean/sla2/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
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

		// 获取方法对应的权限信息
		permInfo, exists := r.methodPermissionMap[info.FullMethod]
		if !exists {
			log.Warn("No permission defined for method", zap.String("method", info.FullMethod))
			return handler(ctx, req)
		}

		// 获取用户ID
		userID, err := service.GetUserID(ctx)
		if err != nil {
			log.Error("Failed to extract user ID",
				zap.String("method", info.FullMethod),
				zap.Error(err))
			return nil, status.Errorf(codes.Unauthenticated, "invalid authentication: %v", err)
		}

		log.Info("Checking permission with detailed info",
			zap.String("method", info.FullMethod),
			zap.Uint64("user_id", uint64(userID)),
			zap.String("resource", permInfo["resource"]),
			zap.String("action", permInfo["action"]))

		// 添加额外的角色检查日志
		roles, roleErr := r.permissionHelper.GetUserRoles(ctx, userID)
		if roleErr != nil {
			log.Warn("Failed to get user roles",
				zap.Uint64("user_id", uint64(userID)),
				zap.Error(roleErr))
		} else {
			log.Info("User roles",
				zap.Uint64("user_id", uint64(userID)),
				zap.Strings("roles", roles))
		}

		// 检查用户是否有admin角色
		hasAdminRole, adminErr := r.permissionHelper.HasUserRole(ctx, userID, security.RoleAdmin)
		if adminErr != nil {
			log.Warn("Failed to check admin role",
				zap.Uint64("user_id", uint64(userID)),
				zap.Error(adminErr))
		} else {
			log.Info("Admin role check result",
				zap.Uint64("user_id", uint64(userID)),
				zap.Bool("has_admin_role", hasAdminRole))
		}

		// 检查权限
		hasPermission, err := r.permissionHelper.CheckUserPermission(ctx, userID, permInfo["resource"], permInfo["action"])
		if err != nil {
			log.Error("Error checking permission",
				zap.Error(err),
				zap.String("method", info.FullMethod),
				zap.Uint64("user_id", uint64(userID)))
			return nil, status.Errorf(codes.Internal, "error checking permission: %v", err)
		}

		if !hasPermission {
			log.Warn("Permission denied",
				zap.String("method", info.FullMethod),
				zap.Uint64("user_id", uint64(userID)),
				zap.String("resource", permInfo["resource"]),
				zap.String("action", permInfo["action"]))
			return nil, status.Error(codes.PermissionDenied, "permission denied")
		}

		log.Debug("Permission granted",
			zap.String("method", info.FullMethod),
			zap.Uint64("user_id", uint64(userID)),
			zap.String("resource", permInfo["resource"]),
			zap.String("action", permInfo["action"]))

		// 权限验证通过，继续处理请求
		return handler(ctx, req)
	}
}
