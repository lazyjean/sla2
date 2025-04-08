package learning

import (
	"context"

	pb "github.com/lazyjean/sla2/api/proto/v1"
	"github.com/lazyjean/sla2/internal/application/service"
	"github.com/lazyjean/sla2/internal/domain/entity"
	"github.com/lazyjean/sla2/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

// GetMemoryStats 获取记忆单元学习统计
func (s *LearningService) GetMemoryStats(ctx context.Context, req *pb.GetMemoryStatsRequest) (*pb.GetMemoryStatsResponse, error) {
	log := logger.GetLogger(ctx)
	log.Info("GetMemoryStats called", zap.Any("request", req))

	// 处理可选的类型过滤
	var unitType *entity.MemoryUnitType
	if req.Type != nil {
		t := ToEntityMemoryUnitType(req.GetType())
		unitType = &t
	}

	// TODO: 处理 tag 和 category 过滤 (需要修改 service 和 repository)

	// 调用 memory service 获取统计信息
	stats, err := s.memoryService.GetMemoryStats(ctx, unitType)
	if err != nil {
		log.Error("Failed to get memory stats from service", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to get memory stats")
	}

	if stats == nil { // Handle case where service returns nil stats (e.g., no data)
		log.Warn("Memory service returned nil stats")
		return &pb.GetMemoryStatsResponse{}, nil // Return empty response
	}

	// 将领域统计信息映射到 PB 响应
	pbResponse := &pb.GetMemoryStatsResponse{
		TotalLearned:    uint32(stats.TotalCount),                     // 映射 TotalCount
		MasteredCount:   uint32(stats.MasteredCount),                  // 映射 MasteredCount
		NeedReviewCount: uint32(stats.LearningCount + stats.NewCount), // 估算: 学习中+新
		TotalStudyTime:  0,                                            // repository.MemoryUnitStats 没有这个字段
		LevelStats:      make(map[uint32]uint32),
		RetentionRates:  make(map[uint32]float32),
	}

	// 填充 level_stats
	pbResponse.LevelStats[uint32(entity.MasteryLevelMastered)] = uint32(stats.MasteredCount)
	pbResponse.LevelStats[uint32(entity.MasteryLevelBeginner)] = 0 // Placeholder
	pbResponse.LevelStats[uint32(entity.MasteryLevelFamiliar)] = 0 // Placeholder
	if stats.LearningCount > 0 {                                   // Approximate learning count split if needed, better handled in repo/service
		// This is just a guess, real breakdown needs more data
		pbResponse.LevelStats[uint32(entity.MasteryLevelBeginner)] = uint32(stats.LearningCount / 2)
		pbResponse.LevelStats[uint32(entity.MasteryLevelFamiliar)] = uint32(stats.LearningCount) - pbResponse.LevelStats[uint32(entity.MasteryLevelBeginner)]
	}
	pbResponse.LevelStats[uint32(entity.MasteryLevelUnlearned)] = uint32(stats.NewCount)
	// Add Expert level if tracked
	// pbResponse.LevelStats[uint32(entity.MasteryLevelExpert)] = ...

	// 填充 retention_rates
	if stats.TotalCount > 0 && stats.RetentionRate > 0 {
		if unitType != nil {
			pbResponse.RetentionRates[uint32(*unitType)] = float32(stats.RetentionRate)
		} else {
			// Maybe add an 'overall' key or calculate per type if repo provided richer data
			pbResponse.RetentionRates[0] = float32(stats.RetentionRate) // Use 0 for overall/unspecified?
		}
	}

	log.Info("GetMemoryStats successful", zap.Any("response", pbResponse))
	return pbResponse, nil
}
