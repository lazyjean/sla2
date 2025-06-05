package repository

import (
	"context"

	"github.com/lazyjean/sla2/internal/domain/entity"
)

// CourseSectionRepository 课程章节仓库接口
type CourseSectionRepository interface {
	GenericRepository[*entity.CourseSection, entity.CourseSectionID]

	// ListByCourseID 获取课程的所有章节
	ListByCourseID(ctx context.Context, courseID entity.CourseID) ([]*entity.CourseSection, error)
}
