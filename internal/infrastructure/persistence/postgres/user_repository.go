package postgres

import (
	"context"
	"github.com/lazyjean/sla2/internal/domain/entity"
	"github.com/lazyjean/sla2/internal/domain/errors"
	"github.com/lazyjean/sla2/internal/domain/repository"
	"gorm.io/gorm"
)

type userRepository struct {
	*repository.GenericRepositoryImpl[*entity.User, entity.UID]
}

func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &userRepository{
		GenericRepositoryImpl: repository.NewGenericRepository[*entity.User, entity.UID](db),
	}
}

// FindByUsername 根据用户名查找用户
func (r *userRepository) FindByUsername(ctx context.Context, username string) (*entity.User, error) {
	var user entity.User
	if err := r.DB.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NewError(errors.CodeUserNotFound, "用户不存在")
		}
		return nil, errors.NewError(errors.CodeInternalError, "查询用户失败")
	}
	return &user, nil
}

// FindByEmail 根据邮箱查找用户
func (r *userRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	if err := r.DB.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NewError(errors.CodeUserNotFound, "用户不存在")
		}
		return nil, errors.NewError(errors.CodeInternalError, "查询用户失败")
	}
	return &user, nil
}

// FindByPhone 根据手机号查找用户
func (r *userRepository) FindByPhone(ctx context.Context, phone string) (*entity.User, error) {
	var user entity.User
	if err := r.DB.WithContext(ctx).Where("phone = ?", phone).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NewError(errors.CodeUserNotFound, "用户不存在")
		}
		return nil, errors.NewError(errors.CodeInternalError, "查询用户失败")
	}
	return &user, nil
}

// FindByAccount 根据账号（用户名/邮箱/手机号）查找用户
func (r *userRepository) FindByAccount(ctx context.Context, account string) (*entity.User, error) {
	var user entity.User
	err := r.DB.WithContext(ctx).Where("username = ? OR email = ? OR phone = ?",
		account, account, account).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NewError(errors.CodeUserNotFound, "用户不存在")
		}
		return nil, errors.NewError(errors.CodeInternalError, "查询用户失败")
	}
	return &user, nil
}

// FindByAppleID 根据苹果用户ID查找用户
func (r *userRepository) FindByAppleID(ctx context.Context, appleID string) (*entity.User, error) {
	var user entity.User
	if err := r.DB.WithContext(ctx).Where("apple_id = ?", appleID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NewError(errors.CodeUserNotFound, "用户不存在")
		}
		return nil, errors.NewError(errors.CodeInternalError, "查询用户失败")
	}
	return &user, nil
}

var _ repository.UserRepository = (*userRepository)(nil)
