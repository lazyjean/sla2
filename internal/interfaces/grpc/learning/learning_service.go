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

// ReviewWord 复习单词
func (s *LearningService) ReviewWord(ctx context.Context, req *pb.ReviewWordRequest) (*pb.ReviewWordResponse, error) {
	// 将 ReviewResult 转换为 bool
	result := req.Result == pb.ReviewResult_REVIEW_RESULT_CORRECT

	// 调用 memory service 的 ReviewWord
	err := s.memoryService.ReviewWord(ctx, entity.WordID(req.WordId), result, req.ResponseTime)
	if err != nil {
		return nil, err
	}

	return &pb.ReviewWordResponse{}, nil
}

// ReviewHanChar 复习汉字
func (s *LearningService) ReviewHanChar(ctx context.Context, req *pb.ReviewHanCharRequest) (*pb.ReviewHanCharResponse, error) {
	// 将 ReviewResult 转换为 bool
	result := req.Result == pb.ReviewResult_REVIEW_RESULT_CORRECT

	// 调用 memory service 的 ReviewHanChar
	err := s.memoryService.ReviewHanChar(ctx, req.HanCharId, result, req.ResponseTime)
	if err != nil {
		return nil, err
	}

	return &pb.ReviewHanCharResponse{}, nil
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

// ListMemoriesForReview 获取需要复习的记忆单元列表
func (s *LearningService) ListMemoriesForReview(ctx context.Context, req *pb.ListMemoriesForReviewRequest) (*pb.ListMemoriesForReviewResponse, error) {
	// 确保 page 和 page_size 合理，例如 page >= 1, page_size > 0
	page := req.Page
	if page < 1 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 10 // 默认值
	}

	// 调用 memory service 获取数据
	// 注意：这里的 service 方法需要返回总数
	domainUnits, totalUnits, err := s.memoryService.ListMemoriesForReview(ctx, page, pageSize, ToEntityMemoryUnitTypes(req.Types))
	if err != nil {
		// TODO: 更细致的错误处理，例如区分 not found 和其他错误
		return nil, err // 返回通用错误，或者转换为 gRPC status error
	}

	// 转换结果为 PB 格式
	statuses := ToPBMemoryUnits(domainUnits)
	total := uint32(totalUnits)

	return &pb.ListMemoriesForReviewResponse{
		Statuses: statuses,
		Total:    total,
	}, nil
}

// ToEntityMemoryUnitTypes 辅助函数: 将 PB 类型列表转换为领域实体类型列表
func ToEntityMemoryUnitTypes(pbTypes []pb.MemoryUnitType) []entity.MemoryUnitType {
	// 如果输入为空，可能表示不过滤，返回 nil 或者根据业务逻辑处理
	if len(pbTypes) == 0 {
		return nil
	}
	entityTypes := make([]entity.MemoryUnitType, len(pbTypes))
	for i, t := range pbTypes {
		entityTypes[i] = entity.MemoryUnitType(t)
	}
	return entityTypes
}

// ToPBMemoryUnits 辅助函数: 将领域实体列表转换为 PB 列表
func ToPBMemoryUnits(units []*entity.MemoryUnit) []*pb.MemoryUnit {
	pbUnits := make([]*pb.MemoryUnit, len(units))
	for i, u := range units {
		pbUnits[i] = ToPBMemoryUnit(u) // 假设已存在 ToPBMemoryUnit
	}
	return pbUnits
}
