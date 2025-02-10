package postgres

import (
	"context"
	"time"

	"github.com/lazyjean/sla2/internal/domain/entity"
	domainErrors "github.com/lazyjean/sla2/internal/domain/errors"
	"gorm.io/gorm"
)

type learningRepository struct {
	db *gorm.DB
}

func NewLearningRepository(db *gorm.DB) *learningRepository {
	return &learningRepository{
		db: db,
	}
}

// UpdateProgress 更新学习进度
func (r *learningRepository) UpdateProgress(ctx context.Context, userID, wordID uint, familiarity int, nextReviewAt time.Time) (*entity.LearningProgress, error) {
	now := time.Now()
	progress := &entity.LearningProgress{
		UserID:         userID,
		WordID:         wordID,
		Familiarity:    familiarity,
		NextReviewAt:   nextReviewAt,
		LastReviewedAt: now,
	}

	err := r.db.WithContext(ctx).
		Where("user_id = ? AND word_id = ?", userID, wordID).
		Assign(progress).
		FirstOrCreate(progress).Error

	if err != nil {
		return nil, err
	}

	return progress, nil
}

// ListByUserID 获取用户的学习进度列表
func (r *learningRepository) ListByUserID(ctx context.Context, userID uint, page, pageSize int) ([]*entity.LearningProgress, int, error) {
	var progresses []*entity.LearningProgress
	var total int64

	offset := (page - 1) * pageSize
	err := r.db.WithContext(ctx).Model(&entity.LearningProgress{}).
		Where("user_id = ?", userID).
		Count(&total).
		Offset(offset).
		Limit(pageSize).
		Find(&progresses).Error

	if err != nil {
		return nil, 0, err
	}

	return progresses, int(total), nil
}

// GetUserStats 获取用户的学习统计信息
func (r *learningRepository) GetUserStats(ctx context.Context, userID uint) (*entity.LearningStats, error) {
	var stats entity.LearningStats
	stats.UserID = userID

	// 获取总单词数
	err := r.db.WithContext(ctx).Model(&entity.LearningProgress{}).
		Where("user_id = ?", userID).
		Count(&stats.TotalWords).Error
	if err != nil {
		return nil, err
	}

	// 获取已掌握的单词数（熟悉度 >= 80）
	err = r.db.WithContext(ctx).Model(&entity.LearningProgress{}).
		Where("user_id = ? AND familiarity >= ?", userID, 80).
		Count(&stats.MasteredWords).Error
	if err != nil {
		return nil, err
	}

	// 获取正在学习的单词数（0 < 熟悉度 < 80）
	err = r.db.WithContext(ctx).Model(&entity.LearningProgress{}).
		Where("user_id = ? AND familiarity > 0 AND familiarity < ?", userID, 80).
		Count(&stats.LearningWords).Error
	if err != nil {
		return nil, err
	}

	// 获取待复习的单词数
	now := time.Now()
	err = r.db.WithContext(ctx).Model(&entity.LearningProgress{}).
		Where("user_id = ? AND next_review_at <= ?", userID, now).
		Count(&stats.ReviewDueCount).Error
	if err != nil {
		return nil, err
	}

	// 获取最后学习时间
	var lastProgress entity.LearningProgress
	err = r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("last_reviewed_at DESC").
		First(&lastProgress).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if err == nil {
		stats.LastStudyTime = lastProgress.LastReviewedAt
	}

	// 获取今日学习数量
	today := time.Now().Truncate(24 * time.Hour)
	err = r.db.WithContext(ctx).Model(&entity.LearningProgress{}).
		Where("user_id = ? AND last_reviewed_at >= ?", userID, today).
		Count(&stats.TodayStudyCount).Error
	if err != nil {
		return nil, err
	}

	// 计算连续学习天数
	stats.ContinuousDays = r.calculateContinuousDays(ctx, userID)

	return &stats, nil
}

// ListReviewWords 获取用户待复习的单词列表
func (r *learningRepository) ListReviewWords(ctx context.Context, userID uint, page, pageSize int) ([]*entity.Word, []*entity.LearningProgress, int, error) {
	var progresses []*entity.LearningProgress
	var total int64
	now := time.Now()

	offset := (page - 1) * pageSize
	err := r.db.WithContext(ctx).Model(&entity.LearningProgress{}).
		Where("user_id = ? AND next_review_at <= ?", userID, now).
		Count(&total).
		Offset(offset).
		Limit(pageSize).
		Find(&progresses).Error

	if err != nil {
		return nil, nil, 0, err
	}

	if len(progresses) == 0 {
		return nil, nil, 0, nil
	}

	// 获取单词ID列表
	var wordIDs []uint
	for _, p := range progresses {
		wordIDs = append(wordIDs, p.WordID)
	}

	// 获取单词详情
	var words []*entity.Word
	err = r.db.WithContext(ctx).
		Where("id IN ?", wordIDs).
		Find(&words).Error
	if err != nil {
		return nil, nil, 0, err
	}

	return words, progresses, int(total), nil
}

// calculateContinuousDays 计算连续学习天数
func (r *learningRepository) calculateContinuousDays(ctx context.Context, userID uint) int64 {
	var dates []time.Time
	err := r.db.WithContext(ctx).Model(&entity.LearningProgress{}).
		Select("DATE(last_reviewed_at) as date").
		Where("user_id = ?", userID).
		Group("DATE(last_reviewed_at)").
		Order("date DESC").
		Pluck("date", &dates).Error

	if err != nil || len(dates) == 0 {
		return 0
	}

	continuousDays := int64(1)
	today := time.Now().Truncate(24 * time.Hour)

	// 如果最后一次学习不是今天，从昨天开始计算
	if dates[0].Truncate(24 * time.Hour).Before(today) {
		today = today.AddDate(0, 0, -1)
	}

	for i := 0; i < len(dates)-1; i++ {
		expectedDate := today.AddDate(0, 0, -i)
		if dates[i].Truncate(24 * time.Hour).Equal(expectedDate) {
			continuousDays++
		} else {
			break
		}
	}

	return continuousDays
}

// SaveCourseProgress 保存课程进度
func (r *learningRepository) SaveCourseProgress(ctx context.Context, progress *entity.CourseLearningProgress) error {
	if err := r.db.WithContext(ctx).Save(progress).Error; err != nil {
		return domainErrors.ErrFailedToSave
	}
	return nil
}

// GetCourseProgress 获取课程进度
func (r *learningRepository) GetCourseProgress(ctx context.Context, userID, courseID uint) (*entity.CourseLearningProgress, error) {
	var progress entity.CourseLearningProgress
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND course_id = ?", userID, courseID).
		First(&progress).Error

	if err == gorm.ErrRecordNotFound {
		return nil, domainErrors.ErrProgressNotFound
	}
	if err != nil {
		return nil, domainErrors.ErrFailedToQuery
	}
	return &progress, nil
}

// ListCourseProgress 列出用户的课程进度
func (r *learningRepository) ListCourseProgress(ctx context.Context, userID uint, offset, limit int) ([]*entity.CourseLearningProgress, int64, error) {
	var progresses []*entity.CourseLearningProgress
	var total int64

	err := r.db.WithContext(ctx).Model(&entity.CourseLearningProgress{}).
		Where("user_id = ?", userID).
		Count(&total).Error
	if err != nil {
		return nil, 0, domainErrors.ErrFailedToQuery
	}

	err = r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("updated_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&progresses).Error
	if err != nil {
		return nil, 0, domainErrors.ErrFailedToQuery
	}

	return progresses, total, nil
}

// SaveSectionProgress 保存章节进度
func (r *learningRepository) SaveSectionProgress(ctx context.Context, progress *entity.SectionProgress) error {
	if err := r.db.WithContext(ctx).Save(progress).Error; err != nil {
		return domainErrors.ErrFailedToSave
	}
	return nil
}

// GetSectionProgress 获取章节进度
func (r *learningRepository) GetSectionProgress(ctx context.Context, userID, sectionID uint) (*entity.SectionProgress, error) {
	var progress entity.SectionProgress
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND section_id = ?", userID, sectionID).
		First(&progress).Error

	if err == gorm.ErrRecordNotFound {
		return nil, domainErrors.ErrProgressNotFound
	}
	if err != nil {
		return nil, domainErrors.ErrFailedToQuery
	}
	return &progress, nil
}

// ListSectionProgress 列出章节进度
func (r *learningRepository) ListSectionProgress(ctx context.Context, userID, courseID uint) ([]*entity.SectionProgress, error) {
	var progresses []*entity.SectionProgress
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND course_id = ?", userID, courseID).
		Order("section_id ASC").
		Find(&progresses).Error
	if err != nil {
		return nil, domainErrors.ErrFailedToQuery
	}
	return progresses, nil
}

// SaveUnitProgress 保存单元进度
func (r *learningRepository) SaveUnitProgress(ctx context.Context, progress *entity.UnitProgress) error {
	if err := r.db.WithContext(ctx).Save(progress).Error; err != nil {
		return domainErrors.ErrFailedToSave
	}
	return nil
}

// GetUnitProgress 获取单元进度
func (r *learningRepository) GetUnitProgress(ctx context.Context, userID, unitID uint) (*entity.UnitProgress, error) {
	var progress entity.UnitProgress
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND unit_id = ?", userID, unitID).
		First(&progress).Error

	if err == gorm.ErrRecordNotFound {
		return nil, domainErrors.ErrProgressNotFound
	}
	if err != nil {
		return nil, domainErrors.ErrFailedToQuery
	}
	return &progress, nil
}

// ListUnitProgress 列出单元进度
func (r *learningRepository) ListUnitProgress(ctx context.Context, userID, sectionID uint) ([]*entity.UnitProgress, error) {
	var progresses []*entity.UnitProgress
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND section_id = ?", userID, sectionID).
		Order("unit_id ASC").
		Find(&progresses).Error
	if err != nil {
		return nil, domainErrors.ErrFailedToQuery
	}
	return progresses, nil
}
