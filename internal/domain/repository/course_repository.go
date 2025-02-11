package repository

import (
	"context"

	"github.com/lazyjean/sla2/internal/domain/entity"
)

// CourseRepository 课程仓库接口
type CourseRepository interface {
	// Create 创建课程
	Create(ctx context.Context, course *entity.Course) error

	// Update 更新课程
	Update(ctx context.Context, course *entity.Course) error

	// Delete 删除课程
	Delete(ctx context.Context, id uint) error

	// GetByID 根据ID获取课程
	GetByID(ctx context.Context, id uint) (*entity.Course, error)

	// List 获取课程列表
	List(ctx context.Context, offset, limit int, filters map[string]interface{}) ([]*entity.Course, int64, error)

	// Search 搜索课程
	Search(ctx context.Context, keyword string, offset, limit int, filters map[string]interface{}) ([]*entity.Course, int64, error)
}
