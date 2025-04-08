package learning

import (
	"context"

	pb "github.com/lazyjean/sla2/api/proto/v1"
	"github.com/lazyjean/sla2/internal/application/service"
	"github.com/lazyjean/sla2/internal/domain/entity"
)

type LearningService struct {
	pb.UnimplementedLearningServiceServer
	learningService *service.LearningService
	memoryService   service.MemoryService
}

func NewLearningService(learningService *service.LearningService, memoryService service.MemoryService) *LearningService {
	return &LearningService{
		learningService: learningService,
		memoryService:   memoryService,
	}
}

// GetCourseProgress 获取课程学习进度
func (s *LearningService) GetCourseProgress(ctx context.Context, req *pb.LearningServiceGetCourseProgressRequest) (*pb.LearningServiceGetCourseProgressResponse, error) {
	courseID := req.CourseId

	// 获取课程进度及统计信息
	progress, completedSections, totalSections, err := s.learningService.GetCourseProgressWithStats(ctx, uint(courseID))
	if err != nil {
		return nil, err
	}

	return &pb.LearningServiceGetCourseProgressResponse{
		Progress: &pb.LearningProgress{
			Progress:       float32(progress),
			CompletedItems: uint32(completedSections),
			TotalItems:     uint32(totalSections),
		},
	}, nil
}

// GetSectionProgress 获取章节学习进度
func (s *LearningService) GetSectionProgress(ctx context.Context, req *pb.LearningServiceGetSectionProgressRequest) (*pb.LearningServiceGetSectionProgressResponse, error) {
	sectionID := req.SectionId

	// 获取章节进度及统计信息
	progress, completedUnits, totalUnits, err := s.learningService.GetSectionProgressWithStats(ctx, uint(sectionID))
	if err != nil {
		return nil, err
	}

	// 获取已完成的单元ID列表
	unitProgresses, err := s.learningService.ListUnitProgress(ctx, uint(sectionID))
	if err != nil {
		return nil, err
	}

	// 收集已完成的单元ID
	completedUnitIDs := make([]uint32, 0)
	for _, unit := range unitProgresses {
		if unit.Status == "completed" {
			completedUnitIDs = append(completedUnitIDs, uint32(unit.UnitID))
		}
	}

	return &pb.LearningServiceGetSectionProgressResponse{
		Progress: &pb.LearningProgress{
			Progress:       float32(progress),
			CompletedItems: uint32(completedUnits),
			TotalItems:     uint32(totalUnits),
		},
		CompletedUnitIds: completedUnitIDs,
	}, nil
}

// UpdateUnitProgress 更新单元学习进度
func (s *LearningService) UpdateUnitProgress(ctx context.Context, req *pb.LearningServiceUpdateUnitProgressRequest) (*pb.LearningServiceUpdateUnitProgressResponse, error) {
	err := s.learningService.UpdateUnitProgress(ctx, uint(req.UnitId), uint(req.SectionId), req.Completed)
	if err != nil {
		return nil, err
	}
	return &pb.LearningServiceUpdateUnitProgressResponse{}, nil
}

// UpdateMemoryStatus 更新记忆单元状态
func (s *LearningService) UpdateMemoryStatus(ctx context.Context, req *pb.UpdateMemoryStatusRequest) (*pb.UpdateMemoryStatusResponse, error) {
	err := s.learningService.UpdateMemoryStatus(ctx, req.MemoryUnitId, entity.MasteryLevel(req.MasteryLevel), req.StudyDuration)
	if err != nil {
		return nil, err
	}

	return &pb.UpdateMemoryStatusResponse{}, nil
}

// RecordLearningResult 记录学习结果
func (s *LearningService) RecordLearningResult(ctx context.Context, req *pb.RecordLearningResultRequest) (*pb.RecordLearningResultResponse, error) {
	// 将 ReviewResult 转换为 bool
	result := req.Result == pb.ReviewResult_REVIEW_RESULT_CORRECT

	// 记录学习结果
	err := s.learningService.RecordLearningResult(ctx, req.MemoryUnitId, result, req.ResponseTime, req.UserNotes)
	if err != nil {
		return nil, err
	}

	return &pb.RecordLearningResultResponse{}, nil
}

// GetMemoryStatus 获取记忆单元状态
func (s *LearningService) GetMemoryStatus(ctx context.Context, req *pb.GetMemoryStatusRequest) (*pb.GetMemoryStatusResponse, error) {
	// 获取记忆单元
	memoryUnit, err := s.memoryService.GetMemoryUnit(ctx, req.MemoryUnitId)
	if err != nil {
		return nil, err
	}

	return &pb.GetMemoryStatusResponse{
		Status: ToPBMemoryUnit(memoryUnit),
	}, nil
}
