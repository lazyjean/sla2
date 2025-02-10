package repository

import (
	"context"
	"time"

	"github.com/lazyjean/sla2/internal/domain/entity"
)

// LearningRepository 定义学习进度相关的仓储接口
type LearningRepository interface {
	// UpdateProgress 更新学习进度
	UpdateProgress(ctx context.Context, userID, wordID uint, familiarity int, nextReviewAt time.Time) (*entity.LearningProgress, error)
	// ListByUserID 获取用户的学习进度列表
	ListByUserID(ctx context.Context, userID uint, page, pageSize int) ([]*entity.LearningProgress, int, error)
	// GetUserStats 获取用户的学习统计信息
	GetUserStats(ctx context.Context, userID uint) (*entity.LearningStats, error)
	// ListReviewWords 获取用户待复习的单词列表
	ListReviewWords(ctx context.Context, userID uint, page, pageSize int) ([]*entity.Word, []*entity.LearningProgress, int, error)

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
