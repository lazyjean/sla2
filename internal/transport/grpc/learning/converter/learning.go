package converter

import (
	pb "github.com/lazyjean/sla2/api/proto/v1"
	"github.com/lazyjean/sla2/internal/application/service"
	"github.com/lazyjean/sla2/internal/domain/entity"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type LearningConverter struct{}

func NewLearningConverter() *LearningConverter {
	return &LearningConverter{}
}

// ToPBLearningProgress 将领域实体转换为 PB 学习进度
func (c *LearningConverter) ToPBLearningProgress(progress float64, completedItems, totalItems int) *pb.LearningProgress {
	return &pb.LearningProgress{
		Progress:       float32(progress),
		CompletedItems: uint32(completedItems),
		TotalItems:     uint32(totalItems),
	}
}

// ToPBMemoryUnit 将领域实体转换为 PB 格式
func (c *LearningConverter) ToPBMemoryUnit(memoryUnit *entity.MemoryUnit) *pb.MemoryUnit {
	return &pb.MemoryUnit{
		Id:                 uint32(memoryUnit.ID),
		Type:               pb.MemoryUnitType(memoryUnit.Type),
		ContentId:          memoryUnit.ContentID,
		MasteryLevel:       pb.MasteryLevel(memoryUnit.MasteryLevel),
		ReviewCount:        memoryUnit.ReviewCount,
		StudyDuration:      memoryUnit.StudyDuration,
		ConsecutiveCorrect: memoryUnit.ConsecutiveCorrect,
		ConsecutiveWrong:   memoryUnit.ConsecutiveWrong,
		LastReviewAt:       timestamppb.New(memoryUnit.LastReviewAt),
		NextReviewAt:       timestamppb.New(memoryUnit.NextReviewAt),
		RetentionRate:      memoryUnit.RetentionRate,
	}
}

// ToPBReviewInterval 将领域实体转换为 PB 复习间隔
func (c *LearningConverter) ToPBReviewInterval(interval *service.ReviewInterval) *pb.ReviewInterval {
	return &pb.ReviewInterval{
		Days:    interval.Days,
		Hours:   interval.Hours,
		Minutes: interval.Minutes,
	}
}

// ToEntityMasteryLevel 将 PB 掌握程度转换为领域实体
func (c *LearningConverter) ToEntityMasteryLevel(level pb.MasteryLevel) entity.MasteryLevel {
	return entity.MasteryLevel(level)
}

// ToEntityMemoryUnitType 将 PB 记忆单元类型转换为领域实体
func (c *LearningConverter) ToEntityMemoryUnitType(unitType pb.MemoryUnitType) entity.MemoryUnitType {
	return entity.MemoryUnitType(unitType)
}

// ToEntityReviewResult 将 PB 复习结果转换为领域实体
func (c *LearningConverter) ToEntityReviewResult(result pb.ReviewResult) entity.ReviewResult {
	return entity.ReviewResult(result)
}

// ToPBMemoryUnits 将领域实体切片转换为 PB 格式切片
func (c *LearningConverter) ToPBMemoryUnits(memoryUnits []*entity.MemoryUnit) []*pb.MemoryUnit {
	if memoryUnits == nil {
		return nil
	}
	pbUnits := make([]*pb.MemoryUnit, len(memoryUnits))
	for i, unit := range memoryUnits {
		pbUnits[i] = c.ToPBMemoryUnit(unit)
	}
	return pbUnits
}
