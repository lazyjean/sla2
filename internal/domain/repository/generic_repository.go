package repository

import (
	"context"
	"time"

	"gorm.io/gorm"
)

// Entity 实体接口
type Entity[ID comparable] interface {
	GetID() ID
	SetID(id ID)
	GetCreatedAt() time.Time
	SetCreatedAt(t time.Time)
	GetUpdatedAt() time.Time
	SetUpdatedAt(t time.Time)
	GetDeletedAt() gorm.DeletedAt
	SetDeletedAt(t gorm.DeletedAt)
}

// ValidatableEntity 可验证实体接口
type ValidatableEntity interface {
	Validate() error
}

// GenericRepository 通用仓储接口
type GenericRepository[T Entity[ID], ID comparable] interface {
	// Create 创建实体
	Create(ctx context.Context, entity T) error
	// CreateBatch 批量创建实体
	CreateBatch(ctx context.Context, entities []T) error
	// Update 更新实体
	Update(ctx context.Context, entity T) error
	// UpdateBatch 批量更新实体
	UpdateBatch(ctx context.Context, entities []T) error
	// Delete 删除实体
	Delete(ctx context.Context, id ID) error
	// GetByID 根据ID获取实体
	GetByID(ctx context.Context, id ID) (T, error)
	// List 获取实体列表
	List(ctx context.Context, offset, limit int) ([]T, error)
	// Count 获取实体总数
	Count(ctx context.Context) (int64, error)
}

// GenericRepositoryImpl 通用仓储实现
type GenericRepositoryImpl[T Entity[ID], ID comparable] struct {
	DB *gorm.DB
}

// NewGenericRepository 创建通用仓储
func NewGenericRepository[T Entity[ID], ID comparable](db *gorm.DB) *GenericRepositoryImpl[T, ID] {
	return &GenericRepositoryImpl[T, ID]{
		DB: db,
	}
}

// Create 创建实体
func (r *GenericRepositoryImpl[T, ID]) Create(ctx context.Context, entity T) error {
	// 如果实体实现了 ValidatableEntity 接口，进行验证
	if validatable, ok := any(entity).(ValidatableEntity); ok {
		if err := validatable.Validate(); err != nil {
			return err
		}
	}

	now := time.Now()
	entity.SetCreatedAt(now)
	entity.SetUpdatedAt(now)
	return r.DB.WithContext(ctx).Create(&entity).Error
}

// CreateBatch 批量创建实体
func (r *GenericRepositoryImpl[T, ID]) CreateBatch(ctx context.Context, entities []T) error {
	now := time.Now()
	for i := range entities {
		// 如果实体实现了 ValidatableEntity 接口，进行验证
		if validatable, ok := any(entities[i]).(ValidatableEntity); ok {
			if err := validatable.Validate(); err != nil {
				return err
			}
		}
		entities[i].SetCreatedAt(now)
		entities[i].SetUpdatedAt(now)
	}
	return r.DB.WithContext(ctx).Create(&entities).Error
}

// Update 更新实体
func (r *GenericRepositoryImpl[T, ID]) Update(ctx context.Context, entity T) error {
	// 如果实体实现了 ValidatableEntity 接口，进行验证
	if validatable, ok := any(entity).(ValidatableEntity); ok {
		if err := validatable.Validate(); err != nil {
			return err
		}
	}

	entity.SetUpdatedAt(time.Now())
	return r.DB.WithContext(ctx).Save(&entity).Error
}

// UpdateBatch 批量更新实体
func (r *GenericRepositoryImpl[T, ID]) UpdateBatch(ctx context.Context, entities []T) error {
	now := time.Now()
	for i := range entities {
		// 如果实体实现了 ValidatableEntity 接口，进行验证
		if validatable, ok := any(entities[i]).(ValidatableEntity); ok {
			if err := validatable.Validate(); err != nil {
				return err
			}
		}
		entities[i].SetUpdatedAt(now)
	}
	return r.DB.WithContext(ctx).Save(&entities).Error
}

// Delete 删除实体
func (r *GenericRepositoryImpl[T, ID]) Delete(ctx context.Context, id ID) error {
	return r.DB.WithContext(ctx).Delete(new(T), "id = ?", id).Error
}

// GetByID 根据ID获取实体
func (r *GenericRepositoryImpl[T, ID]) GetByID(ctx context.Context, id ID) (T, error) {
	var entity T
	err := r.DB.WithContext(ctx).First(&entity, "id = ?", id).Error
	return entity, err
}

// List 获取实体列表
func (r *GenericRepositoryImpl[T, ID]) List(ctx context.Context, offset, limit int) ([]T, error) {
	var entities []T
	err := r.DB.WithContext(ctx).Offset(offset).Limit(limit).Find(&entities).Error
	return entities, err
}

// Count 获取实体总数
func (r *GenericRepositoryImpl[T, ID]) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.DB.WithContext(ctx).Model(new(T)).Count(&count).Error
	return count, err
}
