package postgres

import (
	"context"

	"github.com/lazyjean/sla2/internal/domain/entity"
	"github.com/lazyjean/sla2/internal/domain/repository"
	"gorm.io/gorm"
)

// CourseRepository PostgreSQL 课程仓库实现
type courseRepository struct {
	db *gorm.DB
}

// NewCourseRepository 创建课程仓库实例
func NewCourseRepository(db *gorm.DB) repository.CourseRepository {
	return &courseRepository{
		db: db,
	}
}

// Create 创建课程
func (r *courseRepository) Create(ctx context.Context, course *entity.Course) error {
	return r.db.WithContext(ctx).Create(course).Error
}

// Update 更新课程
func (r *courseRepository) Update(ctx context.Context, course *entity.Course) error {
	return r.db.WithContext(ctx).Save(course).Error
}

// Delete 删除课程
func (r *courseRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&entity.Course{}, id).Error
}

// GetByID 根据ID获取课程
func (r *courseRepository) GetByID(ctx context.Context, id uint) (*entity.Course, error) {
	var course entity.Course
	err := r.db.WithContext(ctx).First(&course, id).Error
	if err != nil {
		return nil, err
	}
	return &course, nil
}

// List 获取课程列表
func (r *courseRepository) List(ctx context.Context, offset, limit int, filters map[string]interface{}) ([]*entity.Course, int64, error) {
	var courses []*entity.Course
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.Course{})

	// 应用过滤条件
	for key, value := range filters {
		switch key {
		case "level":
			if v, ok := value.(string); ok && v != "" {
				query = query.Where("level = ?", v)
			}
		case "status":
			if v, ok := value.(string); ok && v != "" {
				query = query.Where("status = ?", v)
			}
		case "tags":
			if v, ok := value.(string); ok && v != "" {
				query = query.Where("tags @> ARRAY[?]::varchar(50)[]", v)
			}
		}
	}

	// 获取总数
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err = query.Session(&gorm.Session{PrepareStmt: true}).Offset(offset).Limit(limit).Find(&courses).Error
	if err != nil {
		return nil, 0, err
	}

	return courses, total, nil
}

// Search 搜索课程
func (r *courseRepository) Search(ctx context.Context, keyword string, offset, limit int, filters map[string]interface{}) ([]*entity.Course, int64, error) {
	var courses []*entity.Course
	var total int64

	query := r.db.WithContext(ctx).
		Where("title ILIKE ? OR description ILIKE ?", "%"+keyword+"%", "%"+keyword+"%")

	// 应用过滤条件
	for key, value := range filters {
		switch key {
		case "level", "status":
			if v, ok := value.(string); ok && v != "" {
				query = query.Where(key+" = ?", v)
			}
		case "tag":
			if v, ok := value.(string); ok && v != "" {
				query = query.Where("tags @> ARRAY[?]::varchar(50)[]", v)
			}
		}
	}

	if err := query.Model(&entity.Course{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Session(&gorm.Session{PrepareStmt: true}).Offset(offset).Limit(limit).Find(&courses).Error
	return courses, total, err
}

var _ repository.CourseRepository = (*courseRepository)(nil)
