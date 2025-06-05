package security

import (
	"context"
	"fmt"
	"strconv"

	"github.com/lazyjean/sla2/internal/domain/entity"
	"github.com/lazyjean/sla2/pkg/logger"
	"go.uber.org/zap"
)

// SubjectFromString 从字符串中解析主体(用户ID或角色名)
func SubjectFromString(subject string) string {
	// 如果是数字，假设是用户ID
	if _, err := strconv.ParseInt(subject, 10, 64); err == nil {
		return fmt.Sprintf("u:%s", subject)
	}
	// 否则假设是角色名
	return fmt.Sprintf("r:%s", subject)
}

// SubjectFromUserID 从用户ID创建主体标识符
func SubjectFromUserID(userID entity.UID) string {
	return fmt.Sprintf("u:%d", userID)
}

// SubjectFromRoleName 从角色名创建主体标识符
func SubjectFromRoleName(roleName string) string {
	return fmt.Sprintf("r:%s", roleName)
}

// PermissionHelper 权限验证辅助工具
type PermissionHelper struct {
	permissionManager PermissionManager
}

// NewPermissionHelper 创建新的权限验证辅助工具
func NewPermissionHelper(permissionManager PermissionManager) *PermissionHelper {
	return &PermissionHelper{
		permissionManager: permissionManager,
	}
}

// CheckUserPermission 检查用户对特定资源和操作的权限
func (ph *PermissionHelper) CheckUserPermission(ctx context.Context, userID entity.UID, resource, action string) (bool, error) {
	if userID == 0 {
		logger.GetLogger(ctx).Warn("Checking permission for invalid user ID",
			zap.String("resource", resource),
			zap.String("action", action))
		return false, nil
	}

	sub := SubjectFromUserID(userID)
	return ph.permissionManager.CheckPermission(ctx, sub, resource, action)
}

// RequireUserPermission 要求用户拥有特定权限，如果没有则返回错误
func (ph *PermissionHelper) RequireUserPermission(ctx context.Context, userID entity.UID, resource, action string) error {
	if userID == 0 {
		return fmt.Errorf("invalid user ID")
	}

	hasPermission, err := ph.CheckUserPermission(ctx, userID, resource, action)
	if err != nil {
		return fmt.Errorf("failed to check permission: %w", err)
	}

	if !hasPermission {
		return fmt.Errorf("user %d does not have permission to %s on %s", userID, action, resource)
	}

	return nil
}

// CheckRolePermission 检查角色对特定资源和操作的权限
func (ph *PermissionHelper) CheckRolePermission(ctx context.Context, roleName, resource, action string) (bool, error) {
	if roleName == "" {
		logger.GetLogger(ctx).Warn("Checking permission for empty role name",
			zap.String("resource", resource),
			zap.String("action", action))
		return false, nil
	}

	sub := SubjectFromRoleName(roleName)
	return ph.permissionManager.CheckPermission(ctx, sub, resource, action)
}

// HasUserRole 检查用户是否拥有特定角色
func (ph *PermissionHelper) HasUserRole(ctx context.Context, userID entity.UID, roleName string) (bool, error) {
	if userID == 0 || roleName == "" {
		return false, nil
	}

	role := SubjectFromRoleName(roleName)
	return ph.permissionManager.HasRoleForUser(ctx, userID, role)
}

// GetUserRoles 获取用户的所有角色
func (ph *PermissionHelper) GetUserRoles(ctx context.Context, userID entity.UID) ([]string, error) {
	if userID == 0 {
		return []string{}, nil
	}

	return ph.permissionManager.GetRolesForUser(ctx, userID)
}

// AssignRoleToUser 为用户分配角色
func (ph *PermissionHelper) AssignRoleToUser(ctx context.Context, userID entity.UID, roleName string) (bool, error) {
	if userID == 0 || roleName == "" {
		return false, fmt.Errorf("invalid user ID or role name")
	}

	role := SubjectFromRoleName(roleName)
	return ph.permissionManager.AddRoleForUser(ctx, userID, role)
}

// RemoveRoleFromUser 从用户移除角色
func (ph *PermissionHelper) RemoveRoleFromUser(ctx context.Context, userID entity.UID, roleName string) (bool, error) {
	if userID == 0 || roleName == "" {
		return false, fmt.Errorf("invalid user ID or role name")
	}

	role := SubjectFromRoleName(roleName)
	return ph.permissionManager.DeleteRoleForUser(ctx, userID, role)
}
