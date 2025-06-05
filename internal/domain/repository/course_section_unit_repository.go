package repository

import (
	"context"

	"github.com/lazyjean/sla2/internal/domain/entity"
)

// CourseSectionUnitRepository 课程章节单元仓库接口
type CourseSectionUnitRepository interface {
	GenericRepository[*entity.CourseSectionUnit, entity.CourseSectionUnitID]

	// ListBySectionID 获取章节的所有单元
	ListBySectionID(ctx context.Context, sectionID entity.CourseSectionID) ([]*entity.CourseSectionUnit, error)
}
