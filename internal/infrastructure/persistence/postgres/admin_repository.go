package postgres

import (
	"context"
	"errors"

	"github.com/lazyjean/sla2/internal/domain/entity"
	"github.com/lazyjean/sla2/internal/domain/repository"
	"github.com/lazyjean/sla2/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// adminRepository 管理员仓储实现
type adminRepository struct {
	db *gorm.DB
}

// NewAdminRepository 创建管理员仓储
func NewAdminRepository(db *gorm.DB) repository.AdminRepository {
	return &adminRepository{
		db: db,
	}
}

// IsInitialized 检查系统是否已初始化
func (r *adminRepository) IsInitialized(ctx context.Context) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&entity.Admin{}).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// Create 创建管理后续需要扩展
func (r *adminRepository) Create(ctx context.Context, admin *entity.Admin) error {
	log := logger.GetLogger(ctx)
	if admin == nil {
		return errors.New("admin is nil")
	}

	// 确保 Roles 字段是有效的 JSON 数组
	if admin.Roles == nil {
		admin.Roles = []string{}
	}

	// 调试：打印实体字段值
	log.Info("Creating admin with email_verified", zap.Any("admin", admin))

	return r.db.WithContext(ctx).Select("*").Create(admin).Error
}

// FindByID 根据ID查找管理员
func (r *adminRepository) FindByID(ctx context.Context, adminID entity.UID) (*entity.Admin, error) {
	var admin entity.Admin
	if err := r.db.WithContext(ctx).Where("id = ?", adminID).First(&admin).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("admin not found")
		}
		return nil, err
	}
	return &admin, nil
}

// FindByUsername 根据用户名查找管理员
func (r *adminRepository) FindByUsername(ctx context.Context, username string) (*entity.Admin, error) {
	var admin entity.Admin
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&admin).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("admin not found")
		}
		return nil, err
	}
	return &admin, nil
}

// Delete 删除管理员
func (r *adminRepository) Delete(ctx context.Context, adminID entity.UID) error {
	result := r.db.WithContext(ctx).Delete(&entity.Admin{}, adminID)
	if result.Error != nil {
		return result.Error
	}

	// 如果没有找到记录，返回错误
	if result.RowsAffected == 0 {
		return errors.New("admin not found")
	}

	return nil
}
