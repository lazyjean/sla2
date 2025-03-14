package repository

import (
	"context"

	"github.com/lazyjean/sla2/internal/domain/entity"
)

// AdminRepository 管理员仓储接口
type AdminRepository interface {
	// 系统状态
	IsInitialized(ctx context.Context) (bool, error)

	// 管理员操作
	Create(ctx context.Context, admin *entity.Admin) error
	FindByID(ctx context.Context, adminID entity.UID) (*entity.Admin, error)
	FindByUsername(ctx context.Context, username string) (*entity.Admin, error)
	Delete(ctx context.Context, adminID entity.UID) error
}
