package repository

import (
	"context"

	"github.com/lazyjean/sla2/internal/domain/entity"
)

// AdminRepository 管理员仓储接口
type AdminRepository interface {
	GenericRepository[*entity.Admin, entity.UID]

	// IsInitialized 系统状态
	IsInitialized(ctx context.Context) (bool, error)

	// FindByUsername 管理员特定操作
	FindByUsername(ctx context.Context, username string) (*entity.Admin, error)
}
