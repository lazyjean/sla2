package repository

import (
	"context"
	"time"

	"github.com/lazyjean/sla2/internal/domain/entity"
)

// MemoryUnitRepository 记忆单元仓储接口
type MemoryUnitRepository interface {
	// GenericRepository 继承通用仓储接口的基础功能
	GenericRepository[*entity.MemoryUnit, entity.MemoryUnitID]

	// GetOrCreate 获取或创建记忆单元
	GetOrCreate(ctx context.Context, userID entity.UID, unitType entity.MemoryUnitType, contentID uint32) (*entity.MemoryUnit, error)

	// ListNeedReview 获取需要复习的记忆单元列表
	ListNeedReview(ctx context.Context, userID entity.UID, unitType entity.MemoryUnitType, before time.Time, limit int) ([]*entity.MemoryUnit, error)

	// ListByUserIDAndType 获取用户指定类型的所有记忆单元
	ListByUserIDAndType(ctx context.Context, userID entity.UID, unitType entity.MemoryUnitType) ([]*entity.MemoryUnit, error)

	// CountNeedReviewByType 根据类型列表计算需要复习的记忆单元总数
	CountNeedReviewByType(ctx context.Context, userID entity.UID, unitType entity.MemoryUnitType, before time.Time) (int64, error)

	// GetStats 获取统计信息
	GetStats(ctx context.Context, userID entity.UID, unitType entity.MemoryUnitType) (*entity.MemoryUnitStats, error)
}
