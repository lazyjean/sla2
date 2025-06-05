package converter

import (
	pb "github.com/lazyjean/sla2/api/proto/v1"
	"github.com/lazyjean/sla2/internal/domain/entity"
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

// ToEntityMasteryLevel 将 PB 掌握程度转换为领域实体
func (c *LearningConverter) ToEntityMasteryLevel(level pb.MasteryLevel) entity.MasteryLevel {
	return entity.MasteryLevel(level)
}

// ToEntityMemoryUnitType 将 PB 记忆单元类型转换为领域实体
func (c *LearningConverter) ToEntityMemoryUnitType(reqType *pb.MemoryUnitType) *entity.MemoryUnitType {
	if reqType == nil {
		return nil
	}
	t := entity.MemoryUnitType(*reqType)
	return &t
}
