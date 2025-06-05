package postgres

import (
	"context"

	"github.com/lazyjean/sla2/internal/domain/entity"
	"github.com/lazyjean/sla2/internal/domain/repository"
	"gorm.io/gorm"
)

// courseRepository PostgreSQL 课程仓库实现
type courseRepository struct {
	*repository.GenericRepositoryImpl[*entity.Course, entity.CourseID]
}

// NewCourseRepository 创建课程仓库实例
func NewCourseRepository(db *gorm.DB) repository.CourseRepository {
	return &courseRepository{
		GenericRepositoryImpl: repository.NewGenericRepository[*entity.Course, entity.CourseID](db),
	}
}

// ListWithFilters 获取课程列表（带过滤条件）
func (r *courseRepository) ListWithFilters(ctx context.Context, offset, limit int, filters map[string]interface{}) ([]*entity.Course, int64, error) {
	var courses []*entity.Course
	var total int64

	query := r.DB.WithContext(ctx).Model(&entity.Course{})

	// 应用过滤条件
	for key, value := range filters {
		switch key {
		case "level":
			if v, ok := value.(string); ok && v != "" {
				query = query.Where("level = ?", v)
			}
		case "category":
			if v, ok := value.(string); ok && v != "" {
				query = query.Where("category = ?", v)
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

	query := r.DB.WithContext(ctx).
		Where("title ILIKE ? OR description ILIKE ?", "%"+keyword+"%", "%"+keyword+"%")

	// 应用过滤条件
	for key, value := range filters {
		switch key {
		case "level", "status", "category":
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
