package postgres

import (
	"context"

	"github.com/lazyjean/sla2/internal/domain/entity"
	"github.com/lazyjean/sla2/internal/domain/repository"
	"gorm.io/gorm"
)

// courseSectionRepository PostgreSQL 课程章节仓库实现
type courseSectionRepository struct {
	*repository.GenericRepositoryImpl[*entity.CourseSection, entity.CourseSectionID]
}

// NewCourseSectionRepository 创建课程章节仓库实例
func NewCourseSectionRepository(db *gorm.DB) repository.CourseSectionRepository {
	return &courseSectionRepository{
		GenericRepositoryImpl: repository.NewGenericRepository[*entity.CourseSection, entity.CourseSectionID](db),
	}
}

// ListByCourseID 获取课程的所有章节
func (r *courseSectionRepository) ListByCourseID(ctx context.Context, courseID entity.CourseID) ([]*entity.CourseSection, error) {
	var sections []*entity.CourseSection
	err := r.DB.WithContext(ctx).Where("course_id = ?", courseID).Order("order_index asc").Find(&sections).Error
	if err != nil {
		return nil, err
	}
	return sections, nil
}
