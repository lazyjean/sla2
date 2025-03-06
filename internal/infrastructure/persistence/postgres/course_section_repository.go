package postgres

import (
	"context"

	"github.com/lazyjean/sla2/internal/domain/entity"
	"github.com/lazyjean/sla2/internal/domain/repository"
	"gorm.io/gorm"
)

// CourseSectionRepository PostgreSQL 课程章节仓库实现
type courseSectionRepository struct {
	db *gorm.DB
}

// NewCourseSectionRepository 创建课程章节仓库实例
func NewCourseSectionRepository(db *gorm.DB) repository.CourseSectionRepository {
	return &courseSectionRepository{
		db: db,
	}
}

// Create 创建章节
func (r *courseSectionRepository) Create(ctx context.Context, section *entity.CourseSection) error {
	return r.db.WithContext(ctx).Create(section).Error
}

// Update 更新章节
func (r *courseSectionRepository) Update(ctx context.Context, section *entity.CourseSection) error {
	return r.db.WithContext(ctx).Save(section).Error
}

// Delete 删除章节
func (r *courseSectionRepository) Delete(ctx context.Context, id entity.CourseSectionID) error {
	return r.db.WithContext(ctx).Delete(&entity.CourseSection{}, id).Error
}

// GetByID 根据ID获取章节
func (r *courseSectionRepository) GetByID(ctx context.Context, id entity.CourseSectionID) (*entity.CourseSection, error) {
	var section entity.CourseSection
	err := r.db.WithContext(ctx).First(&section, id).Error
	if err != nil {
		return nil, err
	}
	return &section, nil
}

// ListByCourseID 获取课程的所有章节
func (r *courseSectionRepository) ListByCourseID(ctx context.Context, courseID entity.CourseID) ([]*entity.CourseSection, error) {
	var sections []*entity.CourseSection
	err := r.db.WithContext(ctx).Where("course_id = ?", courseID).Order("order_index asc").Find(&sections).Error
	if err != nil {
		return nil, err
	}
	return sections, nil
}
