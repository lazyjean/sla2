package postgres

import (
	"context"
	"time"

	"github.com/lazyjean/sla2/internal/domain/entity"
	"github.com/lazyjean/sla2/internal/domain/repository"
	"gorm.io/gorm"
)

type wordLearningRepository struct {
	db *gorm.DB
}

// NewWordLearningRepository 创建单词学习进度仓储
func NewWordLearningRepository(db *gorm.DB) repository.WordLearningRepository {
	return &wordLearningRepository{db: db}
}

// SaveProgress 保存学习进度
func (r *wordLearningRepository) SaveProgress(ctx context.Context, progress *entity.WordLearningProgress) error {
	return r.db.WithContext(ctx).Create(progress).Error
}

// GetProgress 获取学习进度
func (r *wordLearningRepository) GetProgress(ctx context.Context, userID, wordID uint) (*entity.WordLearningProgress, error) {
	var progress entity.WordLearningProgress
	err := r.db.WithContext(ctx).Where("user_id = ? AND word_id = ?", userID, wordID).First(&progress).Error
	if err != nil {
		return nil, err
	}
	return &progress, nil
}

// UpdateProgress 更新学习进度
func (r *wordLearningRepository) UpdateProgress(ctx context.Context, userID, wordID uint, familiarity int, nextReviewAt time.Time) (*entity.WordLearningProgress, error) {
	var progress entity.WordLearningProgress
	err := r.db.WithContext(ctx).Where("user_id = ? AND word_id = ?", userID, wordID).First(&progress).Error
	if err == gorm.ErrRecordNotFound {
		// 如果记录不存在，创建新记录
		progress = *entity.NewWordLearningProgress(userID, wordID, familiarity, nextReviewAt)
		if err := r.SaveProgress(ctx, &progress); err != nil {
			return nil, err
		}
		return &progress, nil
	}
	if err != nil {
		return nil, err
	}

	// 更新现有记录
	progress.UpdateProgress(familiarity, nextReviewAt)
	if err := r.db.WithContext(ctx).Save(&progress).Error; err != nil {
		return nil, err
	}
	return &progress, nil
}

// ListProgress 获取学习进度列表
func (r *wordLearningRepository) ListProgress(ctx context.Context, userID uint, offset, limit int) ([]*entity.WordLearningProgress, int64, error) {
	var total int64
	if err := r.db.WithContext(ctx).Model(&entity.WordLearningProgress{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var progresses []*entity.WordLearningProgress
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Offset(offset).Limit(limit).Find(&progresses).Error
	if err != nil {
		return nil, 0, err
	}

	return progresses, total, nil
}

// ListNeedReviewWords 获取需要复习的单词列表
func (r *wordLearningRepository) ListNeedReviewWords(ctx context.Context, userID uint, before time.Time, offset, limit int) ([]*entity.Word, int64, error) {
	var total int64
	subQuery := r.db.Model(&entity.WordLearningProgress{}).
		Select("word_id").
		Where("user_id = ? AND next_review_at <= ?", userID, before)

	if err := r.db.Model(&entity.Word{}).Where("id IN (?)", subQuery).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var words []*entity.Word
	err := r.db.Where("id IN (?)", subQuery).Offset(offset).Limit(limit).Find(&words).Error
	if err != nil {
		return nil, 0, err
	}

	return words, total, nil
}
