package learning

import (
	"context"

	pb "github.com/lazyjean/sla2/api/proto/v1"
	"github.com/lazyjean/sla2/internal/application/service"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type LearningService struct {
	pb.UnimplementedLearningServiceServer
	learningService *service.LearningService
}

func NewLearningService(learningService *service.LearningService) *LearningService {
	return &LearningService{
		learningService: learningService,
	}
}

// GetCourseProgress 获取课程学习进度
func (s *LearningService) GetCourseProgress(ctx context.Context, req *pb.GetCourseProgressRequest) (*pb.GetCourseProgressResponse, error) {
	courseID := req.CourseId

	// TODO: 从上下文中获取用户ID
	userID := uint(1)

	// 获取章节总数
	sections, err := s.learningService.ListSectionProgress(ctx, userID, uint(courseID))
	if err != nil {
		return nil, err
	}

	completedSections := 0
	for _, section := range sections {
		if section.Status == "completed" {
			completedSections++
		}
	}

	return &pb.GetCourseProgressResponse{
		Progress:         100.0, // 临时固定值，需要根据实际进度计算
		CompletedSection: uint32(completedSections),
		TotalSection:     uint32(len(sections)),
	}, nil
}

// GetSectionProgress 获取章节学习进度
func (s *LearningService) GetSectionProgress(ctx context.Context, req *pb.GetSectionProgressRequest) (*pb.GetSectionProgressResponse, error) {
	sectionID := req.SectionId

	// TODO: 从上下文中获取用户ID
	userID := uint(1)

	// 获取单元总数
	units, err := s.learningService.ListUnitProgress(ctx, userID, uint(sectionID))
	if err != nil {
		return nil, err
	}

	completedUnits := 0
	for _, unit := range units {
		if unit.Status == "completed" {
			completedUnits++
		}
	}

	return &pb.GetSectionProgressResponse{
		Progress:       float32(completedUnits) / float32(len(units)) * 100,
		CompletedUnits: uint32(completedUnits),
		TotalUnits:     uint32(len(units)),
	}, nil
}

// GetUnitProgress 获取单元学习进度
func (s *LearningService) GetUnitProgress(ctx context.Context, req *pb.GetUnitProgressRequest) (*pb.GetUnitProgressResponse, error) {
	unitID := req.UnitId

	// TODO: 从上下文中获取用户ID
	userID := uint(1)

	progress, err := s.learningService.GetUnitProgress(ctx, userID, uint(unitID))
	if err != nil {
		return nil, err
	}

	return &pb.GetUnitProgressResponse{
		Completed:      progress.Status == "completed",
		LastAccessTime: timestamppb.New(progress.UpdatedAt),
		StudyDuration:  0, // 临时固定值，需要根据实际学习时长计算
	}, nil
}

// UpdateUnitProgress 更新单元学习进度
func (s *LearningService) UpdateUnitProgress(ctx context.Context, req *pb.UpdateUnitProgressRequest) (*pb.UpdateUnitProgressResponse, error) {
	unitID := req.UnitId

	// TODO: 从上下文中获取用户ID
	userID := uint(1)

	status := "in_progress"
	if req.Completed {
		status = "completed"
	}

	// 获取当前单元所属的章节ID
	currentProgress, err := s.learningService.GetUnitProgress(ctx, userID, uint(unitID))
	if err != nil {
		return nil, err
	}

	_, err = s.learningService.SaveUnitProgress(ctx, userID, currentProgress.SectionID, uint(unitID), status, float64(req.StudyDuration), nil)
	if err != nil {
		return nil, err
	}

	return &pb.UpdateUnitProgressResponse{}, nil
}
