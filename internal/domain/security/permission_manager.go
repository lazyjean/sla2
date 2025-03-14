package security

import (
	"context"

	"github.com/lazyjean/sla2/internal/domain/entity"
)

// PermissionManager 权限管理器接口
type PermissionManager interface {
	// CheckPermission 检查主体是否有权限执行操作
	// sub: 主体（用户ID、角色名等）
	// obj: 资源对象
	// act: 操作动作
	CheckPermission(ctx context.Context, sub string, obj string, act string) (bool, error)

	// AddPolicy 添加权限策略
	// sub: 主体（用户ID、角色名等）
	// obj: 资源对象
	// act: 操作动作
	AddPolicy(ctx context.Context, sub string, obj string, act string) (bool, error)

	// RemovePolicy 移除权限策略
	RemovePolicy(ctx context.Context, sub string, obj string, act string) (bool, error)

	// AddRoleForUser 为用户添加角色
	AddRoleForUser(ctx context.Context, user entity.UID, role string) (bool, error)

	// DeleteRoleForUser 删除用户的角色
	DeleteRoleForUser(ctx context.Context, user entity.UID, role string) (bool, error)

	// GetRolesForUser 获取用户的所有角色
	GetRolesForUser(ctx context.Context, user entity.UID) ([]string, error)

	// GetUsersForRole 获取拥有指定角色的所有用户
	GetUsersForRole(ctx context.Context, role string) ([]string, error)

	// HasRoleForUser 检查用户是否拥有指定角色
	HasRoleForUser(ctx context.Context, user entity.UID, role string) (bool, error)

	// GetAllRoles 获取所有角色
	GetAllRoles(ctx context.Context) ([]string, error)

	// GetPermissionsForUser 获取用户的所有权限
	GetPermissionsForUser(ctx context.Context, user entity.UID) ([][]string, error)

	// GetPermissionsForRole 获取角色的所有权限
	GetPermissionsForRole(ctx context.Context, role string) ([][]string, error)

	// AddUserPermission 直接添加用户权限
	AddUserPermission(ctx context.Context, userId entity.UID, obj string, act string) (bool, error)

	// RemoveUserPermission 移除用户权限
	RemoveUserPermission(ctx context.Context, userId entity.UID, obj string, act string) (bool, error)

	// GetAllPermissions 获取所有权限策略
	GetAllPermissions(ctx context.Context) ([][]string, error)

	// GetAllRolesToPermissions 获取所有角色及其权限的映射
	GetAllRolesToPermissions(ctx context.Context) (map[string][][]string, error)

	// LoadPolicy 从存储中加载策略
	LoadPolicy(ctx context.Context) error

	// SavePolicy 保存策略到存储
	SavePolicy(ctx context.Context) error
}
