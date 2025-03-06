package repository

import (
	"context"

	"github.com/lazyjean/sla2/internal/domain/entity"
)

// CourseSectionRepository 课程章节仓库接口
type CourseSectionRepository interface {
	// Create 创建章节
	Create(ctx context.Context, section *entity.CourseSection) error

	// Update 更新章节
	Update(ctx context.Context, section *entity.CourseSection) error

	// Delete 删除章节
	Delete(ctx context.Context, id entity.CourseSectionID) error

	// GetByID 根据ID获取章节
	GetByID(ctx context.Context, id entity.CourseSectionID) (*entity.CourseSection, error)

	// ListByCourseID 获取课程的所有章节
	ListByCourseID(ctx context.Context, courseID entity.CourseID) ([]*entity.CourseSection, error)

	// CreateUnit 创建单元
	CreateUnit(ctx context.Context, unit *entity.CourseSectionUnit) error

	// UpdateUnit 更新单元
	UpdateUnit(ctx context.Context, unit *entity.CourseSectionUnit) error

	// DeleteUnit 删除单元
	DeleteUnit(ctx context.Context, id entity.CourseSectionUnitID) error

	// GetUnitByID 根据ID获取单元
	GetUnitByID(ctx context.Context, id entity.CourseSectionUnitID) (*entity.CourseSectionUnit, error)

	// ListUnitsBySectionID 获取章节的所有单元
	ListUnitsBySectionID(ctx context.Context, sectionID entity.CourseSectionID) ([]*entity.CourseSectionUnit, error)
}
