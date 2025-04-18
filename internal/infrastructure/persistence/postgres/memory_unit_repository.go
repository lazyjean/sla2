package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/lazyjean/sla2/internal/domain/entity"
	"github.com/lazyjean/sla2/internal/domain/repository"
	"gorm.io/gorm"
)

// memoryUnitRepository 记忆单元仓储实现
type memoryUnitRepository struct {
	db *gorm.DB
}

// NewMemoryUnitRepository 创建记忆单元仓储实例
func NewMemoryUnitRepository(db *gorm.DB) repository.MemoryUnitRepository {
	return &memoryUnitRepository{
		db: db,
	}
}

// Create 创建记忆单元
func (r *memoryUnitRepository) Create(ctx context.Context, unit *entity.MemoryUnit) error {
	return r.db.WithContext(ctx).Create(unit).Error
}

// Update 更新记忆单元
func (r *memoryUnitRepository) Update(ctx context.Context, unit *entity.MemoryUnit) error {
	return r.db.WithContext(ctx).Save(unit).Error
}

// GetByID 根据ID获取记忆单元
func (r *memoryUnitRepository) GetByID(ctx context.Context, id uint32) (*entity.MemoryUnit, error) {
	var unit entity.MemoryUnit
	err := r.db.WithContext(ctx).First(&unit, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &unit, nil
}

// GetByTypeAndContentID 根据类型和内容ID获取记忆单元
func (r *memoryUnitRepository) GetByTypeAndContentID(ctx context.Context, unitType entity.MemoryUnitType, contentID uint32) (*entity.MemoryUnit, error) {
	var unit entity.MemoryUnit
	err := r.db.WithContext(ctx).
		Where("type = ? AND content_id = ?", unitType, contentID).
		First(&unit).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &unit, nil
}

// ListNeedReview 获取需要复习的记忆单元列表
func (r *memoryUnitRepository) ListNeedReview(ctx context.Context, unitType entity.MemoryUnitType, now time.Time, limit int) ([]*entity.MemoryUnit, error) {
	var units []*entity.MemoryUnit
	err := r.db.WithContext(ctx).
		Where("type = ? AND next_review_at <= ?", unitType, now).
		Order("next_review_at ASC").
		Limit(limit).
		Find(&units).Error
	if err != nil {
		return nil, err
	}
	return units, nil
}

// ListByUserID 获取用户的所有记忆单元
func (r *memoryUnitRepository) ListByUserID(ctx context.Context, userID uint32) ([]*entity.MemoryUnit, error) {
	var units []*entity.MemoryUnit
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Find(&units).Error
	if err != nil {
		return nil, err
	}
	return units, nil
}

// ListByUserIDAndType 获取用户指定类型的所有记忆单元
func (r *memoryUnitRepository) ListByUserIDAndType(ctx context.Context, userID uint32, unitType entity.MemoryUnitType) ([]*entity.MemoryUnit, error) {
	var units []*entity.MemoryUnit
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND type = ?", userID, unitType).
		Find(&units).Error
	if err != nil {
		return nil, err
	}
	return units, nil
}

// ListNeedReviewByTypes 根据类型列表获取需要复习的记忆单元列表（分页）
func (r *memoryUnitRepository) ListNeedReviewByTypes(ctx context.Context, types []entity.MemoryUnitType, before time.Time, offset uint32, limit int) ([]*entity.MemoryUnit, error) {
	var units []*entity.MemoryUnit
	query := r.db.WithContext(ctx).
		Where("next_review_at <= ?", before)

	if len(types) > 0 {
		query = query.Where("type IN ?", types)
	}

	err := query.Order("next_review_at ASC").
		Offset(int(offset)).
		Limit(limit).
		Find(&units).Error

	if err != nil {
		return nil, err
	}
	return units, nil
}

// CountNeedReviewByTypes 根据类型列表计算需要复习的记忆单元总数
func (r *memoryUnitRepository) CountNeedReviewByTypes(ctx context.Context, types []entity.MemoryUnitType, before time.Time) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).
		Model(&entity.MemoryUnit{}).
		Where("next_review_at <= ?", before)

	if len(types) > 0 {
		query = query.Where("type IN ?", types)
	}

	err := query.Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// GetStats 获取统计信息
func (r *memoryUnitRepository) GetStats(ctx context.Context, userID entity.UID, unitType entity.MemoryUnitType) (*repository.MemoryUnitStats, error) {
	var stats repository.MemoryUnitStats
	var totalCount, masteredCount, learningCount, newCount int64

	// 获取总数
	if err := r.db.WithContext(ctx).
		Model(&entity.MemoryUnit{}).
		Where("user_id = ? AND type = ?", userID, unitType).
		Count(&totalCount).Error; err != nil {
		return nil, err
	}
	stats.TotalCount = int(totalCount)

	// 获取已掌握数量（掌握和精通）
	if err := r.db.WithContext(ctx).
		Model(&entity.MemoryUnit{}).
		Where("user_id = ? AND type = ? AND mastery_level IN ?", userID, unitType, []entity.MasteryLevel{entity.MasteryLevelMastered, entity.MasteryLevelExpert}).
		Count(&masteredCount).Error; err != nil {
		return nil, err
	}
	stats.MasteredCount = int(masteredCount)

	// 获取学习中数量（初学和熟悉）
	if err := r.db.WithContext(ctx).
		Model(&entity.MemoryUnit{}).
		Where("user_id = ? AND type = ? AND mastery_level IN ?", userID, unitType, []entity.MasteryLevel{entity.MasteryLevelBeginner, entity.MasteryLevelFamiliar}).
		Count(&learningCount).Error; err != nil {
		return nil, err
	}
	stats.LearningCount = int(learningCount)

	// 获取新内容数量（未学习）
	if err := r.db.WithContext(ctx).
		Model(&entity.MemoryUnit{}).
		Where("user_id = ? AND type = ? AND mastery_level = ?", userID, unitType, entity.MasteryLevelUnlearned).
		Count(&newCount).Error; err != nil {
		return nil, err
	}
	stats.NewCount = int(newCount)

	// 计算掌握率
	if stats.TotalCount > 0 {
		stats.MasteryRate = float64(stats.MasteredCount) / float64(stats.TotalCount)
	}

	// 计算记忆保持率（这里使用一个简单的公式：掌握率 + 学习中的比例）
	if stats.TotalCount > 0 {
		stats.RetentionRate = float64(stats.MasteredCount+stats.LearningCount) / float64(stats.TotalCount)
	}

	// 设置每日目标（这里使用固定值，实际应用中可能需要根据用户学习情况动态调整）
	stats.DailyReviewGoal = 20
	stats.DailyNewGoal = 10

	return &stats, nil
}
