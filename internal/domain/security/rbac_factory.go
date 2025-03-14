package security

import (
	"context"
	"fmt"

	"github.com/lazyjean/sla2/config"
	"github.com/lazyjean/sla2/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// RBACProvider RBAC权限系统供应商
type RBACProvider struct {
	permissionManager PermissionManager
	permissionHelper  *PermissionHelper
	initializer       *PermissionInitializer
}

// NewRBACProvider 创建新的RBAC权限系统供应商
func NewRBACProvider(db *gorm.DB, cfg *config.RBACConfig) (*RBACProvider, error) {
	// 创建Casbin权限管理器
	permManager, err := NewCasbinPermissionManager(db, cfg.ConfigDir)
	if err != nil {
		return nil, fmt.Errorf("failed to create permission manager: %w", err)
	}

	// 创建权限辅助工具
	permHelper := NewPermissionHelper(permManager)

	// 创建权限初始化器
	initializer := NewPermissionInitializer(permManager)

	// 如果配置了自动初始化，则初始化权限
	logger.Log.Info("Auto-initializing RBAC permissions...")
	if err := initializer.Initialize(context.Background()); err != nil {
		logger.Log.Error("Failed to initialize permissions", zap.Error(err))
		return nil, fmt.Errorf("failed to initialize permissions: %w", err)
	}

	logger.Log.Info("RBAC permissions initialized successfully")

	return &RBACProvider{
		permissionManager: permManager,
		permissionHelper:  permHelper,
		initializer:       initializer,
	}, nil
}

// GetPermissionManager 获取权限管理器
func (p *RBACProvider) GetPermissionManager() PermissionManager {
	return p.permissionManager
}

// GetPermissionHelper 获取权限辅助工具
func (p *RBACProvider) GetPermissionHelper() *PermissionHelper {
	return p.permissionHelper
}

// GetInitializer 获取权限初始化器
func (p *RBACProvider) GetInitializer() *PermissionInitializer {
	return p.initializer
}

// InitializePermissions 初始化权限
func (p *RBACProvider) InitializePermissions(ctx context.Context) error {
	return p.initializer.Initialize(ctx)
}
