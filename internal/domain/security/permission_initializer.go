package security

import (
	"context"
	"fmt"

	"github.com/lazyjean/sla2/pkg/logger"
	"go.uber.org/zap"
)

// PermissionInitializer 权限初始化器
type PermissionInitializer struct {
	permissionManager PermissionManager
	adapter           *PermissionManagerAdapter
}

// NewPermissionInitializer 创建新的权限初始化器
func NewPermissionInitializer(permissionManager PermissionManager) *PermissionInitializer {
	return &PermissionInitializer{
		permissionManager: permissionManager,
		adapter:           NewPermissionManagerAdapter(permissionManager),
	}
}

// Initialize 初始化权限数据
func (pi *PermissionInitializer) Initialize(ctx context.Context) error {
	logger.Log.Info("Initializing RBAC permissions...")

	// 初始化角色
	if err := pi.initializeRoles(ctx); err != nil {
		return fmt.Errorf("failed to initialize roles: %w", err)
	}

	// 初始化权限策略
	if err := pi.initializePolicies(ctx); err != nil {
		return fmt.Errorf("failed to initialize policies: %w", err)
	}

	logger.Log.Info("RBAC permissions initialized successfully")
	return nil
}

// initializeRoles 初始化角色层级关系
func (pi *PermissionInitializer) initializeRoles(ctx context.Context) error {
	// 不再需要复杂的角色继承关系，因为只有有限的两种角色
	// 普通用户不再通过RBAC验证权限
	return nil
}

// initializePolicies 初始化权限策略
func (pi *PermissionInitializer) initializePolicies(ctx context.Context) error {
	// 定义角色权限策略
	policies := [][]string{
		// 管理员角色拥有所有权限
		{fmt.Sprintf("r:%s", RoleAdmin), ResourceAny, ActionAny},

		// 课程管理员角色权限
		{fmt.Sprintf("r:%s", RoleContentManager), ResourceCourse, ActionAny},
		{fmt.Sprintf("r:%s", RoleContentManager), ResourceQuestion, ActionAny},
	}

	// 添加权限策略
	for _, policy := range policies {
		if len(policy) != 3 {
			logger.Log.Warn("Invalid policy format", zap.Strings("policy", policy))
			continue
		}

		sub := policy[0]
		obj := policy[1]
		act := policy[2]

		// 检查策略是否已存在
		permissions, err := pi.adapter.GetPermissionsForUser(ctx, sub)
		if err != nil {
			return err
		}

		exists := false
		for _, p := range permissions {
			if len(p) == 3 && p[0] == sub && p[1] == obj && p[2] == act {
				exists = true
				break
			}
		}

		if !exists {
			_, err := pi.adapter.AddPolicy(ctx, sub, obj, act)
			if err != nil {
				return err
			}
			logger.Log.Info("Added permission policy",
				zap.String("subject", sub),
				zap.String("object", obj),
				zap.String("action", act))
		}
	}

	return nil
}
