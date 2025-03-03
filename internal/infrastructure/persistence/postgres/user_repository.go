package postgres

import (
	"context"
	"strings"
	"time"

	"github.com/lazyjean/sla2/internal/domain/entity"
	"github.com/lazyjean/sla2/internal/domain/errors"
	"github.com/lazyjean/sla2/internal/domain/repository"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, username, email, password, nickname, appleID string, emailVerified bool) (*entity.User, error) {
	user := &entity.User{
		Username:      username,
		Password:      password,
		Email:         email,
		Nickname:      nickname,
		Status:        entity.UserStatusActive,
		EmailVerified: emailVerified,
		AppleID:       appleID,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return nil, errors.NewError(errors.CodeUserAlreadyExists, "用户已存在")
		}
		return nil, errors.NewError(errors.CodeInternalError, "创建用户失败")
	}

	return user, nil
}

func (r *UserRepository) Update(ctx context.Context, user *entity.User) error {
	user.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *UserRepository) FindByID(ctx context.Context, id entity.UID) (*entity.User, error) {
	var user entity.User
	if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NewError(errors.CodeUserNotFound, "用户不存在")
		}
		return nil, errors.NewError(errors.CodeInternalError, "查询用户失败")
	}
	return &user, nil
}

func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*entity.User, error) {
	var user entity.User
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NewError(errors.CodeUserNotFound, "用户不存在")
		}
		return nil, errors.NewError(errors.CodeInternalError, "查询用户失败")
	}
	return &user, nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NewError(errors.CodeUserNotFound, "用户不存在")
		}
		return nil, errors.NewError(errors.CodeInternalError, "查询用户失败")
	}
	return &user, nil
}

func (r *UserRepository) FindByPhone(ctx context.Context, phone string) (*entity.User, error) {
	var user entity.User
	if err := r.db.WithContext(ctx).Where("phone = ?", phone).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NewError(errors.CodeUserNotFound, "用户不存在")
		}
		return nil, errors.NewError(errors.CodeInternalError, "查询用户失败")
	}
	return &user, nil
}

func (r *UserRepository) FindByAccount(ctx context.Context, account string) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).Where("username = ? OR email = ? OR phone = ?",
		account, account, account).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&entity.User{}).Where("username = ?", username).Count(&count).Error; err != nil {
		return false, errors.NewError(errors.CodeInternalError, "查询用户失败")
	}
	return count > 0, nil
}

func (r *UserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&entity.User{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return false, errors.NewError(errors.CodeInternalError, "查询用户失败")
	}
	return count > 0, nil
}

func (r *UserRepository) ExistsByPhone(ctx context.Context, phone string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&entity.User{}).Where("phone = ?", phone).Count(&count).Error; err != nil {
		return false, errors.NewError(errors.CodeInternalError, "查询用户失败")
	}
	return count > 0, nil
}

func (r *UserRepository) Delete(ctx context.Context, id entity.UID) error {
	if err := r.db.WithContext(ctx).Delete(&entity.User{}, id).Error; err != nil {
		return errors.NewError(errors.CodeInternalError, "删除用户失败")
	}
	return nil
}

func (r *UserRepository) FindByAppleID(ctx context.Context, appleID string) (*entity.User, error) {
	var user entity.User
	if err := r.db.WithContext(ctx).Where("apple_id = ?", appleID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NewError(errors.CodeUserNotFound, "用户不存在")
		}
		return nil, errors.NewError(errors.CodeInternalError, "查询用户失败")
	}
	return &user, nil
}

func (r *UserRepository) ExistsByAppleID(ctx context.Context, appleID string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&entity.User{}).Where("apple_id = ?", appleID).Count(&count).Error; err != nil {
		return false, errors.NewError(errors.CodeInternalError, "查询用户失败")
	}
	return count > 0, nil
}
