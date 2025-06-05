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
	*repository.GenericRepositoryImpl[*entity.MemoryUnit, entity.MemoryUnitID]
	db *gorm.DB
}

// NewMemoryUnitRepository 创建记忆单元仓储实例
func NewMemoryUnitRepository(db *gorm.DB) repository.MemoryUnitRepository {
	return &memoryUnitRepository{
		GenericRepositoryImpl: repository.NewGenericRepository[*entity.MemoryUnit, entity.MemoryUnitID](db),
		db:                    db,
	}
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

// ListNeedReviewByTypes 根据类型列表获取需要复习的记忆单元列表（分页）
func (r *memoryUnitRepository) ListNeedReviewByTypes(ctx context.Context, userID entity.UID, types []entity.MemoryUnitType, before time.Time, offset, limit int, masteryLevels []entity.MasteryLevel) ([]*entity.MemoryUnit, error) {
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
		Offset(offset).
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
func (r *memoryUnitRepository) GetStats(ctx context.Context, userID entity.UID, unitType entity.MemoryUnitType) (*entity.MemoryUnitStats, error) {
	var stats entity.MemoryUnitStats
	var totalLearned, masteredCount, learningCount, newCount int64

	// 获取已学习总数（按 content_id 去重）
	if err := r.db.WithContext(ctx).
		Model(&entity.MemoryUnit{}).
		Select("COUNT(DISTINCT content_id)").
		Where("user_id = ? AND type = ?", userID, unitType).
		Count(&totalLearned).Error; err != nil {
		return nil, err
	}
	stats.TotalCount = totalLearned

	// 获取已掌握数量（掌握和精通，按 content_id 去重）
	if err := r.db.WithContext(ctx).
		Model(&entity.MemoryUnit{}).
		Select("COUNT(DISTINCT content_id)").
		Where("user_id = ? AND type = ? AND mastery_level IN ?", userID, unitType, []entity.MasteryLevel{entity.MasteryLevelMastered, entity.MasteryLevelExpert}).
		Count(&masteredCount).Error; err != nil {
		return nil, err
	}
	stats.MasteredCount = masteredCount

	// 获取学习中数量（初学和熟悉，按 content_id 去重）
	if err := r.db.WithContext(ctx).
		Model(&entity.MemoryUnit{}).
		Select("COUNT(DISTINCT content_id)").
		Where("user_id = ? AND type = ? AND mastery_level IN ?", userID, unitType, []entity.MasteryLevel{entity.MasteryLevelBeginner, entity.MasteryLevelFamiliar}).
		Count(&learningCount).Error; err != nil {
		return nil, err
	}
	stats.LearningCount = learningCount

	// 获取新内容数量（未学习，按 content_id 去重）
	if err := r.db.WithContext(ctx).
		Model(&entity.MemoryUnit{}).
		Select("COUNT(DISTINCT content_id)").
		Where("user_id = ? AND type = ? AND mastery_level = ?", userID, unitType, entity.MasteryLevelUnlearned).
		Count(&newCount).Error; err != nil {
		return nil, err
	}
	stats.NewCount = newCount

	// 计算记忆保持率
	if stats.TotalCount > 0 {
		stats.RetentionRate = float64(stats.MasteredCount) / float64(stats.TotalCount)
	}

	return &stats, nil
}

// GetOrCreate 获取或创建记忆单元
func (r *memoryUnitRepository) GetOrCreate(ctx context.Context, userID entity.UID, unitType entity.MemoryUnitType, contentID uint32) (*entity.MemoryUnit, error) {
	var unit entity.MemoryUnit
	result := r.db.WithContext(ctx).Where("user_id = ? AND type = ? AND content_id = ?", userID, unitType, contentID).
		FirstOrCreate(&unit, entity.MemoryUnit{
			UserID:    userID,
			Type:      unitType,
			ContentID: contentID,
		})
	if result.Error != nil {
		return nil, result.Error
	}
	return &unit, nil
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

func (r *memoryUnitRepository) CountNeedReviewByType(ctx context.Context, userID entity.UID, unitType entity.MemoryUnitType, before time.Time) (int64, error) {
	//TODO implement me
	panic("implement me")
}

var _ repository.MemoryUnitRepository = (*memoryUnitRepository)(nil)
