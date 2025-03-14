package security

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"sync"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/lazyjean/sla2/internal/domain/entity"
	rbacembed "github.com/lazyjean/sla2/internal/domain/security/embed"
	"github.com/lazyjean/sla2/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// CasbinPermissionManager 基于Casbin的权限管理器实现
type CasbinPermissionManager struct {
	enforcer  *casbin.SyncedEnforcer
	configDir string
	db        *gorm.DB
	mu        sync.RWMutex
}

// NewCasbinPermissionManager 创建新的Casbin权限管理器
func NewCasbinPermissionManager(db *gorm.DB, configDir string) (*CasbinPermissionManager, error) {
	// 创建数据库适配器
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		return nil, fmt.Errorf("failed to create GORM adapter: %w", err)
	}

	// 尝试从嵌入的文件中加载模型配置
	var m model.Model
	modelData, err := rbacembed.GetRBACModelBytes()
	if err == nil {
		// 从嵌入的文件中创建模型
		m, err = model.NewModelFromString(string(modelData))
		if err != nil {
			return nil, fmt.Errorf("failed to create model from embedded data: %w", err)
		}
		logger.Log.Info("Loaded RBAC model from embedded file")
	} else {
		// 如果无法从嵌入的文件中加载，尝试从文件系统加载
		modelPath := filepath.Join(configDir, "rbac", "model.conf")
		m, err = model.NewModelFromFile(modelPath)
		if err != nil {
			return nil, fmt.Errorf("failed to load model configuration from file: %w", err)
		}
		logger.Log.Info("Loaded RBAC model from file system", zap.String("path", modelPath))
	}

	// 创建enforcer
	enforcer, err := casbin.NewSyncedEnforcer(m, adapter)
	if err != nil {
		return nil, fmt.Errorf("failed to create enforcer: %w", err)
	}

	// 加载策略
	if err := enforcer.LoadPolicy(); err != nil {
		return nil, fmt.Errorf("failed to load policy: %w", err)
	}

	return &CasbinPermissionManager{
		enforcer:  enforcer,
		configDir: configDir,
		db:        db,
	}, nil
}

// userOrRoleKey 获取用户或角色的字符串键
// 用户使用 "u:用户ID"，角色使用 "r:角色名称"
func userOrRoleKey(id interface{}) string {
	switch v := id.(type) {
	case entity.UID:
		return fmt.Sprintf("u:%d", v)
	case uint:
		return fmt.Sprintf("u:%d", v)
	case int:
		return fmt.Sprintf("u:%d", v)
	case int64:
		return fmt.Sprintf("u:%d", v)
	case string:
		// 如果已经是 r: 或 u: 开头，直接返回
		if len(v) > 2 && (v[:2] == "r:" || v[:2] == "u:") {
			return v
		}
		// 默认当作角色名
		return fmt.Sprintf("r:%s", v)
	default:
		return fmt.Sprintf("%v", v)
	}
}

// CheckPermission 检查主体是否有权限执行操作
func (pm *CasbinPermissionManager) CheckPermission(ctx context.Context, sub string, obj string, act string) (bool, error) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	log := logger.GetLogger(ctx)
	log.Debug("Checking permission",
		zap.String("subject", sub),
		zap.String("object", obj),
		zap.String("action", act))

	// 先检查是否有管理员角色，管理员拥有所有权限
	if strings.HasPrefix(sub, "u:") {
		adminRoleKey := "r:admin"
		hasAdminRole, err := pm.enforcer.HasRoleForUser(sub, adminRoleKey)
		if err != nil {
			log.Error("Failed to check admin role",
				zap.String("subject", sub),
				zap.String("admin_role", adminRoleKey),
				zap.Error(err))
		} else if hasAdminRole {
			log.Debug("User has admin role, granting permission",
				zap.String("subject", sub),
				zap.String("object", obj),
				zap.String("action", act))
			return true, nil
		}

		// 记录角色
		roles, err := pm.enforcer.GetRolesForUser(sub)
		if err != nil {
			log.Warn("Failed to get roles for user",
				zap.String("subject", sub),
				zap.Error(err))
		} else {
			log.Debug("User roles",
				zap.String("subject", sub),
				zap.Strings("roles", roles))
		}
	}

	// 标准权限检查
	res, err := pm.enforcer.Enforce(sub, obj, act)
	if err != nil {
		logger.Log.Error("Failed to check permission",
			zap.String("subject", sub),
			zap.String("object", obj),
			zap.String("action", act),
			zap.Error(err))
		return false, err
	}

	if res {
		logger.Log.Debug("Permission granted",
			zap.String("subject", sub),
			zap.String("object", obj),
			zap.String("action", act))
	} else {
		logger.Log.Warn("Permission denied",
			zap.String("subject", sub),
			zap.String("object", obj),
			zap.String("action", act))
	}

	return res, nil
}

// AddPolicy 添加权限策略
func (pm *CasbinPermissionManager) AddPolicy(ctx context.Context, sub string, obj string, act string) (bool, error) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	added, err := pm.enforcer.AddPolicy(sub, obj, act)
	if err != nil {
		return false, err
	}

	if added {
		if err := pm.enforcer.SavePolicy(); err != nil {
			return true, err
		}
	}

	return added, nil
}

// RemovePolicy 移除权限策略
func (pm *CasbinPermissionManager) RemovePolicy(ctx context.Context, sub string, obj string, act string) (bool, error) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	removed, err := pm.enforcer.RemovePolicy(sub, obj, act)
	if err != nil {
		return false, err
	}

	if removed {
		if err := pm.enforcer.SavePolicy(); err != nil {
			return true, err
		}
	}

	return removed, nil
}

// AddRoleForUser 为用户添加角色
func (pm *CasbinPermissionManager) AddRoleForUser(ctx context.Context, user entity.UID, role string) (bool, error) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	userKey := userOrRoleKey(user)
	roleKey := userOrRoleKey(role)

	added, err := pm.enforcer.AddGroupingPolicy(userKey, roleKey)
	if err != nil {
		return false, err
	}

	if added {
		if err := pm.enforcer.SavePolicy(); err != nil {
			return true, err
		}
	}

	return added, nil
}

// DeleteRoleForUser 删除用户的角色
func (pm *CasbinPermissionManager) DeleteRoleForUser(ctx context.Context, user entity.UID, role string) (bool, error) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	userKey := userOrRoleKey(user)
	roleKey := userOrRoleKey(role)

	removed, err := pm.enforcer.RemoveGroupingPolicy(userKey, roleKey)
	if err != nil {
		return false, err
	}

	if removed {
		if err := pm.enforcer.SavePolicy(); err != nil {
			return true, err
		}
	}

	return removed, nil
}

// GetRolesForUser 获取用户的所有角色
func (pm *CasbinPermissionManager) GetRolesForUser(ctx context.Context, user entity.UID) ([]string, error) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	userKey := userOrRoleKey(user)
	roles, err := pm.enforcer.GetRolesForUser(userKey)
	if err != nil {
		return nil, err
	}

	// 过滤掉角色前缀
	result := make([]string, 0, len(roles))
	for _, role := range roles {
		if len(role) > 2 && role[:2] == "r:" {
			result = append(result, role[2:])
		} else {
			result = append(result, role)
		}
	}

	return result, nil
}

// GetUsersForRole 获取拥有指定角色的所有用户
func (pm *CasbinPermissionManager) GetUsersForRole(ctx context.Context, role string) ([]string, error) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	roleKey := userOrRoleKey(role)
	users, err := pm.enforcer.GetUsersForRole(roleKey)
	if err != nil {
		return nil, err
	}

	// 过滤掉用户前缀，并转换为用户ID
	result := make([]string, 0, len(users))
	for _, user := range users {
		if len(user) > 2 && user[:2] == "u:" {
			result = append(result, user[2:])
		} else {
			result = append(result, user)
		}
	}

	return result, nil
}

// HasRoleForUser 检查用户是否拥有指定角色
func (pm *CasbinPermissionManager) HasRoleForUser(ctx context.Context, user entity.UID, role string) (bool, error) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	userKey := userOrRoleKey(user)
	roleKey := userOrRoleKey(role)

	return pm.enforcer.HasRoleForUser(userKey, roleKey)
}

// GetAllRoles 获取所有角色
func (pm *CasbinPermissionManager) GetAllRoles(ctx context.Context) ([]string, error) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	return pm.enforcer.GetAllRoles()
}

// GetPermissionsForUser 获取用户的所有权限
func (pm *CasbinPermissionManager) GetPermissionsForUser(ctx context.Context, user entity.UID) ([][]string, error) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	userKey := userOrRoleKey(user)
	return pm.enforcer.GetPermissionsForUser(userKey)
}

// GetPermissionsForRole 获取角色的所有权限
func (pm *CasbinPermissionManager) GetPermissionsForRole(ctx context.Context, role string) ([][]string, error) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	roleKey := userOrRoleKey(role)
	return pm.enforcer.GetPermissionsForUser(roleKey)
}

// AddUserPermission 直接添加用户权限
func (pm *CasbinPermissionManager) AddUserPermission(ctx context.Context, userId entity.UID, obj string, act string) (bool, error) {
	userKey := userOrRoleKey(userId)
	return pm.AddPolicy(ctx, userKey, obj, act)
}

// RemoveUserPermission 移除用户权限
func (pm *CasbinPermissionManager) RemoveUserPermission(ctx context.Context, userId entity.UID, obj string, act string) (bool, error) {
	userKey := userOrRoleKey(userId)
	return pm.RemovePolicy(ctx, userKey, obj, act)
}

// GetAllPermissions 获取所有权限策略
func (pm *CasbinPermissionManager) GetAllPermissions(ctx context.Context) ([][]string, error) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	return pm.enforcer.GetPolicy()
}

// GetAllRolesToPermissions 获取所有角色及其权限的映射
func (pm *CasbinPermissionManager) GetAllRolesToPermissions(ctx context.Context) (map[string][][]string, error) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	// 获取所有角色
	roles, err := pm.GetAllRoles(ctx)
	if err != nil {
		return nil, err
	}

	// 构建映射
	result := make(map[string][][]string)
	for _, role := range roles {
		// 过滤掉角色前缀
		roleName := role
		if len(role) > 2 && role[:2] == "r:" {
			roleName = role[2:]
		}

		// 获取角色权限
		permissions, err := pm.GetPermissionsForRole(ctx, role)
		if err != nil {
			return nil, err
		}

		result[roleName] = permissions
	}

	return result, nil
}

// LoadPolicy 从存储中加载策略
func (pm *CasbinPermissionManager) LoadPolicy(ctx context.Context) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	return pm.enforcer.LoadPolicy()
}

// SavePolicy 保存策略到存储
func (pm *CasbinPermissionManager) SavePolicy(ctx context.Context) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	return pm.enforcer.SavePolicy()
}
