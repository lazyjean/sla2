package postgres

import (
	"context"
	"time"

	"github.com/lazyjean/sla2/internal/domain/entity"
	"github.com/lazyjean/sla2/internal/domain/repository"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, username, password, email, nickname string) (*entity.User, error) {
	user := &entity.User{
		Username:  username,
		Password:  password,
		Email:     email,
		Nickname:  nickname,
		Status:    entity.UserStatusActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := r.db.WithContext(ctx).Create(user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepository) Update(ctx context.Context, user *entity.User) error {
	user.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *userRepository) FindByID(ctx context.Context, id uint) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByUsername(ctx context.Context, username string) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByPhone(ctx context.Context, phone string) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).Where("phone = ?", phone).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByAccount(ctx context.Context, account string) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).Where("username = ? OR email = ? OR phone = ?",
		account, account, account).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entity.User{}).Where("username = ?", username).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *userRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entity.User{}).Where("email = ?", email).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *userRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&entity.User{}, id).Error
}

func (r *userRepository) CreateWithAppleID(ctx context.Context, user *entity.User) error {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) FindByAppleID(ctx context.Context, appleID string) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).Where("apple_id = ?", appleID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
