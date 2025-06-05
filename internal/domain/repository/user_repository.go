package repository

import (
	"context"

	"github.com/lazyjean/sla2/internal/domain/entity"
)

// UserRepository 用户仓储接口
type UserRepository interface {
	// GenericRepository 继承通用仓储接口
	GenericRepository[*entity.User, entity.UID]
	// FindByUsername 根据用户名查找用户
	FindByUsername(ctx context.Context, username string) (*entity.User, error)
	// FindByEmail 根据邮箱查找用户
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	// FindByPhone 根据手机号查找用户
	FindByPhone(ctx context.Context, phone string) (*entity.User, error)
	// FindByAppleID 根据苹果用户ID查找用户
	FindByAppleID(ctx context.Context, appleID string) (*entity.User, error)
	// FindByAccount 根据账号（用户名/邮箱/手机号）查找用户
	FindByAccount(ctx context.Context, account string) (*entity.User, error)
}
