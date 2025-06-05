package repository

import (
	"context"

	"github.com/lazyjean/sla2/internal/domain/entity"
)

// CourseRepository 课程仓库接口
type CourseRepository interface {
	GenericRepository[*entity.Course, entity.CourseID]

	// ListWithFilters 获取课程列表（带过滤条件）
	ListWithFilters(ctx context.Context, offset, limit int, filters map[string]interface{}) ([]*entity.Course, int64, error)

	// Search 搜索课程
	Search(ctx context.Context, keyword string, offset, limit int, filters map[string]interface{}) ([]*entity.Course, int64, error)
}
