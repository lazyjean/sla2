package service

import (
	"context"
	"time"

	"github.com/lazyjean/sla2/internal/domain/entity"
	"github.com/lazyjean/sla2/internal/domain/repository"
)

// MemoryService 定义记忆单元相关的业务逻辑
type MemoryService interface {
	// CalculateNextReviewInterval 计算下次复习间隔
	CalculateNextReviewInterval(unit *entity.MemoryUnit) time.Duration
	// CalculateRetentionRate 计算记忆保持率
	CalculateRetentionRate(ctx context.Context, userID entity.UID, unitType entity.MemoryUnitType) (float32, error)
}

// MemoryServiceImpl 记忆服务实现
type MemoryServiceImpl struct {
	memoryUnitRepo repository.MemoryUnitRepository
}

// NewMemoryService 创建记忆服务实例
func NewMemoryService(
	memoryUnitRepo repository.MemoryUnitRepository,
) MemoryService {
	return &MemoryServiceImpl{
		memoryUnitRepo: memoryUnitRepo,
	}
}

// CalculateNextReviewInterval 计算下次复习间隔
func (s *MemoryServiceImpl) CalculateNextReviewInterval(unit *entity.MemoryUnit) time.Duration {
	// 基础间隔（小时）
	var baseInterval float64
	switch unit.MasteryLevel {
	case entity.MasteryLevelUnlearned:
		baseInterval = 1 // 1小时
	case entity.MasteryLevelBeginner:
		baseInterval = 4 // 4小时
	case entity.MasteryLevelFamiliar:
		baseInterval = 24 // 1天
	case entity.MasteryLevelMastered:
		baseInterval = 72 // 3天
	case entity.MasteryLevelExpert:
		baseInterval = 168 // 7天
	default:
		baseInterval = 1
	}

	// 根据连续正确次数调整间隔
	if unit.ConsecutiveCorrect > 0 {
		baseInterval *= 1.2 * float64(unit.ConsecutiveCorrect)
	}

	// 根据连续错误次数减少间隔
	if unit.ConsecutiveWrong > 0 {
		baseInterval /= 2.0 * float64(unit.ConsecutiveWrong)
	}

	// 转换为时间间隔
	return time.Duration(baseInterval * float64(time.Hour))
}

// CalculateRetentionRate 计算记忆保持率
func (s *MemoryServiceImpl) CalculateRetentionRate(ctx context.Context, userID entity.UID, unitType entity.MemoryUnitType) (float32, error) {
	// 1. 获取指定类型的所有记忆单元
	units, err := s.memoryUnitRepo.ListByUserIDAndType(ctx, userID, unitType)
	if err != nil {
		return 0, err
	}

	// 2. 计算平均记忆保持率
	var totalRetentionRate float32
	var count uint32
	for _, unit := range units {
		if unit.RetentionRate > 0 {
			totalRetentionRate += unit.RetentionRate
			count++
		}
	}

	if count == 0 {
		return 0, nil
	}

	return totalRetentionRate / float32(count), nil
}

var _ MemoryService = (*MemoryServiceImpl)(nil)
