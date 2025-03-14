package security

import (
	"context"
	"fmt"
	"strconv"

	"github.com/lazyjean/sla2/internal/domain/entity"
)

// StringToUID 将字符串转换为 entity.UID
func StringToUID(s string) (entity.UID, error) {
	// 处理 "r:rolename" 格式的字符串 - 角色不能直接转为 UID
	if len(s) > 2 && s[:2] == "r:" {
		return 0, fmt.Errorf("cannot convert role to UID: %s", s)
	}

	// 处理 "u:userid" 格式的字符串
	if len(s) > 2 && s[:2] == "u:" {
		s = s[2:] // 移除 "u:" 前缀
	}

	// 转换为 uint64
	id, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse %s as UID: %w", s, err)
	}

	return entity.UID(id), nil
}

// PermissionManagerAdapter 权限管理器适配器
// 用于桥接 string 类型与 entity.UID 类型
type PermissionManagerAdapter struct {
	manager PermissionManager
}

// NewPermissionManagerAdapter 创建新的权限管理器适配器
func NewPermissionManagerAdapter(manager PermissionManager) *PermissionManagerAdapter {
	return &PermissionManagerAdapter{
		manager: manager,
	}
}

// GetPermissionsForUser 获取用户的所有权限 (字符串版本)
func (a *PermissionManagerAdapter) GetPermissionsForUser(ctx context.Context, sub string) ([][]string, error) {
	// 如果是角色
	if len(sub) > 2 && sub[:2] == "r:" {
		return a.manager.GetPermissionsForRole(ctx, sub[2:])
	}

	// 普通用户
	uid, err := StringToUID(sub)
	if err != nil {
		return nil, err
	}

	return a.manager.GetPermissionsForUser(ctx, uid)
}

// AddPolicy 添加权限策略
func (a *PermissionManagerAdapter) AddPolicy(ctx context.Context, sub string, obj string, act string) (bool, error) {
	return a.manager.AddPolicy(ctx, sub, obj, act)
}
