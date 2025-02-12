package repository

import (
	"context"

	"github.com/lazyjean/sla2/internal/domain/entity"
)

// LearningRepository 定义学习进度相关的仓储接口
type LearningRepository interface {
	// 课程进度
	SaveCourseProgress(ctx context.Context, progress *entity.CourseLearningProgress) error
	GetCourseProgress(ctx context.Context, userID, courseID uint) (*entity.CourseLearningProgress, error)
	ListCourseProgress(ctx context.Context, userID uint, offset, limit int) ([]*entity.CourseLearningProgress, int64, error)

	// 章节进度
	SaveSectionProgress(ctx context.Context, progress *entity.CourseSectionProgress) error
	GetSectionProgress(ctx context.Context, userID, sectionID uint) (*entity.CourseSectionProgress, error)
	ListSectionProgress(ctx context.Context, userID, courseID uint) ([]*entity.CourseSectionProgress, error)

	// 单元进度
	SaveUnitProgress(ctx context.Context, progress *entity.CourseSectionUnitProgress) error
	GetUnitProgress(ctx context.Context, userID, unitID uint) (*entity.CourseSectionUnitProgress, error)
	ListUnitProgress(ctx context.Context, userID, sectionID uint) ([]*entity.CourseSectionUnitProgress, error)
}
