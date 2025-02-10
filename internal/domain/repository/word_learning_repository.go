package repository

import (
	"context"
	"time"

	"github.com/lazyjean/sla2/internal/domain/entity"
)

// WordLearningRepository 单词学习进度仓储接口
type WordLearningRepository interface {
	// SaveProgress 保存学习进度
	SaveProgress(ctx context.Context, progress *entity.WordLearningProgress) error
	// GetProgress 获取学习进度
	GetProgress(ctx context.Context, userID, wordID uint) (*entity.WordLearningProgress, error)
	// UpdateProgress 更新学习进度
	UpdateProgress(ctx context.Context, userID, wordID uint, familiarity int, nextReviewAt time.Time) (*entity.WordLearningProgress, error)
	// ListProgress 获取学习进度列表
	ListProgress(ctx context.Context, userID uint, offset, limit int) ([]*entity.WordLearningProgress, int64, error)
	// ListNeedReviewWords 获取需要复习的单词列表
	ListNeedReviewWords(ctx context.Context, userID uint, before time.Time, offset, limit int) ([]*entity.Word, int64, error)
}
