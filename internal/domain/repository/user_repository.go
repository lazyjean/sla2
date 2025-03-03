package repository

import (
	"context"

	"github.com/lazyjean/sla2/internal/domain/entity"
)

// UserRepository 用户仓储接口
type UserRepository interface {
	// Create 创建用户
	Create(ctx context.Context, username, email, password, nickname, appleID string, emailVerified bool) (*entity.User, error)
	// Update 更新用户信息
	Update(ctx context.Context, user *entity.User) error
	// FindByID 根据ID查找用户
	FindByID(ctx context.Context, id entity.UID) (*entity.User, error)
	// FindByUsername 根据用户名查找用户
	FindByUsername(ctx context.Context, username string) (*entity.User, error)
	// FindByEmail 根据邮箱查找用户
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	// FindByPhone 根据手机号查找用户
	FindByPhone(ctx context.Context, phone string) (*entity.User, error)
	// FindByAppleID 根据苹果用户ID查找用户
	FindByAppleID(ctx context.Context, appleID string) (*entity.User, error)
	// ExistsByUsername 检查用户名是否存在
	ExistsByUsername(ctx context.Context, username string) (bool, error)
	// ExistsByEmail 检查邮箱是否存在
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	// ExistsByPhone 检查手机号是否存在
	ExistsByPhone(ctx context.Context, phone string) (bool, error)
	// ExistsByAppleID 检查苹果用户ID是否存在
	ExistsByAppleID(ctx context.Context, appleID string) (bool, error)
	// Delete 删除用户
	Delete(ctx context.Context, id entity.UID) error
}
