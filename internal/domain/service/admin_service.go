package service

import (
	"context"

	"github.com/lazyjean/sla2/internal/domain/entity"
)

// AdminService 管理员领域服务接口
type AdminService interface {
	// IsSystemInitialized 检查系统是否已初始化
	IsSystemInitialized(ctx context.Context) (bool, error)

	// InitializeSystem 初始化系统并创建初始管理员
	InitializeSystem(ctx context.Context, admin *entity.Admin) error

	// GetAdminByID 根据ID获取管理员信息
	GetAdminByID(ctx context.Context, adminID string) (*entity.Admin, error)

	// GetAdminByUsername 根据用户名获取管理员信息
	GetAdminByUsername(ctx context.Context, username string) (*entity.Admin, error)
} 