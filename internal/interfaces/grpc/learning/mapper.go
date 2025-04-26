package learning

import (
	pb "github.com/lazyjean/sla2/api/proto/v1"
	"github.com/lazyjean/sla2/internal/application/service"
	"github.com/lazyjean/sla2/internal/domain/entity"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ToPBLearningProgress 将领域实体转换为 PB 学习进度
func ToPBLearningProgress(progress float64, completedItems, totalItems int) *pb.LearningProgress {
	return &pb.LearningProgress{
		Progress:       float32(progress),
		CompletedItems: uint32(completedItems),
		TotalItems:     uint32(totalItems),
	}
}

// ToPBMemoryUnit 将领域实体转换为 PB 记忆单元
func ToPBMemoryUnit(memoryUnit *entity.MemoryUnit) *pb.MemoryUnit {
	return &pb.MemoryUnit{
		Id:                 uint32(memoryUnit.ID),
		UserId:             uint32(memoryUnit.UserID),
		Type:               pb.MemoryUnitType(memoryUnit.Type),
		ContentId:          memoryUnit.ContentID,
		CreatedAt:          timestamppb.New(memoryUnit.CreatedAt),
		UpdatedAt:          timestamppb.New(memoryUnit.UpdatedAt),
		MasteryLevel:       pb.MasteryLevel(memoryUnit.MasteryLevel),
		ReviewCount:        memoryUnit.ReviewCount,
		NextReviewAt:       timestamppb.New(memoryUnit.NextReviewAt),
		LastReviewAt:       timestamppb.New(memoryUnit.LastReviewAt),
		StudyDuration:      memoryUnit.StudyDuration,
		RetentionRate:      memoryUnit.RetentionRate,
		ConsecutiveCorrect: memoryUnit.ConsecutiveCorrect,
		ConsecutiveWrong:   memoryUnit.ConsecutiveWrong,
	}
}

// ToPBReviewInterval 将领域实体转换为 PB 复习间隔
func ToPBReviewInterval(interval *service.ReviewInterval) *pb.ReviewInterval {
	return &pb.ReviewInterval{
		Days:    interval.Days,
		Hours:   interval.Hours,
		Minutes: interval.Minutes,
	}
}

// ToEntityMasteryLevel 将 PB 掌握程度转换为领域实体
func ToEntityMasteryLevel(level pb.MasteryLevel) entity.MasteryLevel {
	return entity.MasteryLevel(level)
}

// ToEntityMemoryUnitType 将 PB 记忆单元类型转换为领域实体
func ToEntityMemoryUnitType(unitType pb.MemoryUnitType) entity.MemoryUnitType {
	return entity.MemoryUnitType(unitType)
}

// ToEntityReviewResult 将 PB 复习结果转换为领域实体
func ToEntityReviewResult(result pb.ReviewResult) entity.ReviewResult {
	return entity.ReviewResult(result)
}
