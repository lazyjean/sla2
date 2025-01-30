package repository

import (
	"context"

	"github.com/lazyjean/sla2/domain/entity"
)

// LearningRepository 学习进度仓储接口
type LearningRepository interface {
	// 课程进度
	SaveCourseProgress(ctx context.Context, progress *entity.CourseLearningProgress) error
	GetCourseProgress(ctx context.Context, userID, courseID uint) (*entity.CourseLearningProgress, error)
	ListCourseProgress(ctx context.Context, userID uint, offset, limit int) ([]*entity.CourseLearningProgress, int64, error)

	// 章节进度
	SaveSectionProgress(ctx context.Context, progress *entity.SectionProgress) error
	GetSectionProgress(ctx context.Context, userID, sectionID uint) (*entity.SectionProgress, error)
	ListSectionProgress(ctx context.Context, userID, courseID uint) ([]*entity.SectionProgress, error)

	// 单元进度
	SaveUnitProgress(ctx context.Context, progress *entity.UnitProgress) error
	GetUnitProgress(ctx context.Context, userID, unitID uint) (*entity.UnitProgress, error)
	ListUnitProgress(ctx context.Context, userID, sectionID uint) ([]*entity.UnitProgress, error)
}
