package repository

import (
	"context"
	"time"

	"github.com/lazyjean/sla2/internal/domain/entity"
)

// MemoryUnitRepository 记忆单元仓储接口
type MemoryUnitRepository interface {
	// Create 创建记忆单元
	Create(ctx context.Context, unit *entity.MemoryUnit) error
	// CreateBatch 批量创建记忆单元
	CreateBatch(ctx context.Context, units []*entity.MemoryUnit) error
	// Update 更新记忆单元
	Update(ctx context.Context, unit *entity.MemoryUnit) error
	// GetByID 通过ID获取记忆单元
	GetByID(ctx context.Context, userID entity.UID, id uint32) (*entity.MemoryUnit, error)
	// GetByTypeAndContentID 通过类型和内容ID获取记忆单元
	GetByTypeAndContentID(ctx context.Context, userID entity.UID, unitType entity.MemoryUnitType, contentID uint32) (*entity.MemoryUnit, error)
	// ListNeedReview 获取需要复习的记忆单元列表 (DEPRECATED? Consider removing if ListNeedReviewByTypes covers all cases)
	ListNeedReview(ctx context.Context, userID entity.UID, unitType entity.MemoryUnitType, before time.Time, limit int) ([]*entity.MemoryUnit, error)
	// ListNeedReviewByTypes 根据类型列表获取需要复习的记忆单元列表（分页）
	ListNeedReviewByTypes(ctx context.Context, userID entity.UID, types []entity.MemoryUnitType, before time.Time, offset uint32, limit int, masteryLevels []entity.MasteryLevel) ([]*entity.MemoryUnit, error)
	// CountNeedReviewByTypes 根据类型列表计算需要复习的记忆单元总数
	CountNeedReviewByTypes(ctx context.Context, userID entity.UID, types []entity.MemoryUnitType, before time.Time, masteryLevels []entity.MasteryLevel) (int64, error)
	// ListByUserID 获取用户的所有记忆单元
	ListByUserID(ctx context.Context, userID entity.UID) ([]*entity.MemoryUnit, error)
	// ListByUserIDAndType 获取用户指定类型的记忆单元
	ListByUserIDAndType(ctx context.Context, userID entity.UID, unitType entity.MemoryUnitType) ([]*entity.MemoryUnit, error)
	// GetStats 获取指定用户的统计信息
	GetStats(ctx context.Context, userID entity.UID, unitType entity.MemoryUnitType) (*MemoryUnitStats, error)
}

// MemoryUnitStats 记忆单元统计信息
type MemoryUnitStats struct {
	TotalCount      int     `json:"total_count"`       // 总数
	MasteredCount   int     `json:"mastered_count"`    // 已掌握数量
	LearningCount   int     `json:"learning_count"`    // 学习中数量
	NewCount        int     `json:"new_count"`         // 新内容数量
	MasteryRate     float64 `json:"mastery_rate"`      // 掌握率
	RetentionRate   float64 `json:"retention_rate"`    // 记忆保持率
	DailyReviewGoal int     `json:"daily_review_goal"` // 每日复习目标
	DailyNewGoal    int     `json:"daily_new_goal"`    // 每日新内容目标
}
