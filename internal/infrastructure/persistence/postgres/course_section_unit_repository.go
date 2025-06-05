package postgres

import (
	"context"

	"github.com/lazyjean/sla2/internal/domain/entity"
	"github.com/lazyjean/sla2/internal/domain/repository"
	"gorm.io/gorm"
)

// courseSectionUnitRepository PostgreSQL 课程章节单元仓库实现
type courseSectionUnitRepository struct {
	*repository.GenericRepositoryImpl[*entity.CourseSectionUnit, entity.CourseSectionUnitID]
}

// NewCourseSectionUnitRepository 创建课程章节单元仓库实例
func NewCourseSectionUnitRepository(db *gorm.DB) repository.CourseSectionUnitRepository {
	return &courseSectionUnitRepository{
		GenericRepositoryImpl: repository.NewGenericRepository[*entity.CourseSectionUnit, entity.CourseSectionUnitID](db),
	}
}

// ListBySectionID 获取章节的所有单元
func (r *courseSectionUnitRepository) ListBySectionID(ctx context.Context, sectionID entity.CourseSectionID) ([]*entity.CourseSectionUnit, error) {
	var units []*entity.CourseSectionUnit
	err := r.DB.WithContext(ctx).Where("section_id = ?", sectionID).Order("order_index asc").Find(&units).Error
	if err != nil {
		return nil, err
	}
	return units, nil
}
