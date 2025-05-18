package postgres

import (
	"context"
	"errors"
	"fmt"
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

// CreateBatch 批量创建记忆单元
func (r *memoryUnitRepository) CreateBatch(ctx context.Context, units []*entity.MemoryUnit) error {
	return r.db.WithContext(ctx).Create(units).Error
}

// Update 更新记忆单元
func (r *memoryUnitRepository) Update(ctx context.Context, unit *entity.MemoryUnit) error {
	return r.db.WithContext(ctx).Save(unit).Error
}

// GetByID 根据ID获取记忆单元
func (r *memoryUnitRepository) GetByID(ctx context.Context, userID entity.UID, id uint32) (*entity.MemoryUnit, error) {
	var unit entity.MemoryUnit
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND id = ?", userID, id).
		First(&unit).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &unit, nil
}

// GetByTypeAndContentID 根据类型和内容ID获取记忆单元
func (r *memoryUnitRepository) GetByTypeAndContentID(ctx context.Context, userID entity.UID, unitType entity.MemoryUnitType, contentID uint32) (*entity.MemoryUnit, error) {
	var unit entity.MemoryUnit
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND type = ? AND content_id = ?", userID, unitType, contentID).
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
func (r *memoryUnitRepository) ListNeedReview(ctx context.Context, userID entity.UID, unitType entity.MemoryUnitType, before time.Time, limit int) ([]*entity.MemoryUnit, error) {
	// 构建基础查询
	query := r.db.WithContext(ctx).Model(&entity.MemoryUnit{}).
		Where("user_id = ? AND type = ? AND next_review_at <= ?::timestamptz", userID, unitType, before) // 使用 timestamptz 类型

	// 执行查询
	var units []*entity.MemoryUnit
	err := query.
		Order("next_review_at ASC"). // 按复习时间升序排序，优先复习时间早的
		Limit(limit).
		Find(&units).Error
	if err != nil {
		return nil, err
	}
	return units, nil
}

// ListByUserID 获取用户的所有记忆单元
func (r *memoryUnitRepository) ListByUserID(ctx context.Context, userID entity.UID) ([]*entity.MemoryUnit, error) {
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
func (r *memoryUnitRepository) ListByUserIDAndType(ctx context.Context, userID entity.UID, unitType entity.MemoryUnitType) ([]*entity.MemoryUnit, error) {
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
func (r *memoryUnitRepository) ListNeedReviewByTypes(ctx context.Context, userID entity.UID, types []entity.MemoryUnitType, before time.Time, offset uint32, limit int, masteryLevels []entity.MasteryLevel) ([]*entity.MemoryUnit, error) {
	// 构建基础查询
	query := r.db.WithContext(ctx).Model(&entity.MemoryUnit{}).
		Where("user_id = ? AND type IN ? AND next_review_at <= ?::timestamptz", userID, types, before)

	// 添加掌握程度条件
	if len(masteryLevels) > 0 {
		query = query.Where("mastery_level IN ?", masteryLevels)
	}

	// 执行查询
	var units []*entity.MemoryUnit
	err := query.
		Order("next_review_at ASC"). // 按复习时间升序排序，优先复习时间早的
		Offset(int(offset)).
		Limit(limit).
		Find(&units).Error
	if err != nil {
		return nil, fmt.Errorf("failed to list memory units: %w", err)
	}

	// 检查是否找到记录
	if len(units) == 0 {
		return nil, nil
	}

	return units, nil
}

// CountNeedReviewByTypes 根据类型列表计算需要复习的记忆单元总数
func (r *memoryUnitRepository) CountNeedReviewByTypes(ctx context.Context, userID entity.UID, types []entity.MemoryUnitType, before time.Time, masteryLevels []entity.MasteryLevel) (int64, error) {
	// 构建基础查询
	query := r.db.WithContext(ctx).Model(&entity.MemoryUnit{}).
		Where("user_id = ? AND type IN ? AND next_review_at <= ?::timestamptz", userID, types, before)

	if len(masteryLevels) > 0 {
		query = query.Where("mastery_level IN ?", masteryLevels)
	}

	var count int64
	err := query.Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// GetStats 获取统计信息
func (r *memoryUnitRepository) GetStats(ctx context.Context, userID entity.UID, unitType entity.MemoryUnitType) (*repository.MemoryUnitStats, error) {
	var stats repository.MemoryUnitStats
	var totalLearned, masteredCount, learningCount, newCount int64

	// 获取已学习总数（按 content_id 去重）
	if err := r.db.WithContext(ctx).
		Model(&entity.MemoryUnit{}).
		Select("COUNT(DISTINCT content_id)").
		Where("user_id = ? AND type = ?", userID, unitType).
		Count(&totalLearned).Error; err != nil {
		return nil, err
	}
	stats.TotalCount = int(totalLearned)

	// 获取已掌握数量（掌握和精通，按 content_id 去重）
	if err := r.db.WithContext(ctx).
		Model(&entity.MemoryUnit{}).
		Select("COUNT(DISTINCT content_id)").
		Where("user_id = ? AND type = ? AND mastery_level IN ?", userID, unitType, []entity.MasteryLevel{entity.MasteryLevelMastered, entity.MasteryLevelExpert}).
		Count(&masteredCount).Error; err != nil {
		return nil, err
	}
	stats.MasteredCount = int(masteredCount)

	// 获取学习中数量（初学和熟悉，按 content_id 去重）
	if err := r.db.WithContext(ctx).
		Model(&entity.MemoryUnit{}).
		Select("COUNT(DISTINCT content_id)").
		Where("user_id = ? AND type = ? AND mastery_level IN ?", userID, unitType, []entity.MasteryLevel{entity.MasteryLevelBeginner, entity.MasteryLevelFamiliar}).
		Count(&learningCount).Error; err != nil {
		return nil, err
	}
	stats.LearningCount = int(learningCount)

	// 获取新内容数量（未学习，按 content_id 去重）
	if err := r.db.WithContext(ctx).
		Model(&entity.MemoryUnit{}).
		Select("COUNT(DISTINCT content_id)").
		Where("user_id = ? AND type = ? AND mastery_level = ?", userID, unitType, entity.MasteryLevelUnlearned).
		Count(&newCount).Error; err != nil {
		return nil, err
	}
	stats.NewCount = int(newCount)

	// 计算掌握率
	if stats.TotalCount > 0 {
		stats.MasteryRate = float64(stats.MasteredCount) / float64(stats.TotalCount)
	}

	// 计算记忆保持率
	// TODO: 实现记忆保持率计算
	stats.RetentionRate = 1.0

	// 设置每日目标
	// TODO: 从配置或用户设置中获取
	stats.DailyReviewGoal = 20
	stats.DailyNewGoal = 5

	return &stats, nil
}
