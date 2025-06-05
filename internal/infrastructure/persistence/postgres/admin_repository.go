package postgres

import (
	"context"

	"github.com/lazyjean/sla2/internal/domain/entity"
	"github.com/lazyjean/sla2/internal/domain/repository"
	"gorm.io/gorm"
)

// adminRepository 管理员仓储实现
type adminRepository struct {
	*repository.GenericRepositoryImpl[*entity.Admin, entity.UID]
}

// NewAdminRepository 创建管理员仓储
func NewAdminRepository(db *gorm.DB) repository.AdminRepository {
	return &adminRepository{
		GenericRepositoryImpl: repository.NewGenericRepository[*entity.Admin, entity.UID](db),
	}
}

// IsInitialized 检查系统是否已初始化
func (r *adminRepository) IsInitialized(ctx context.Context) (bool, error) {
	cnt, err := r.Count(ctx)
	return cnt > 0, err
}

// FindByUsername 根据用户名查找管理员
func (r *adminRepository) FindByUsername(ctx context.Context, username string) (*entity.Admin, error) {
	var admin entity.Admin
	err := r.DB.WithContext(ctx).Where("username = ?", username).First(&admin).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

// Delete 删除管理员
func (r *adminRepository) Delete(ctx context.Context, id entity.UID) error {
	return r.DB.WithContext(ctx).Delete(&entity.Admin{}, "id = ?", id).Error
}
