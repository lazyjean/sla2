package learning

import (
	"context"
	"time"

	"github.com/lazyjean/sla2/internal/transport/grpc/learning/converter"

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
	converter       *converter.LearningConverter
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

// ListMemoriesForReview 获取需要复习的记忆单元列表
func (s *LearningService) ListMemoriesForReview(ctx context.Context, req *pb.LearningServiceListMemoriesForReviewRequest) (*pb.LearningServiceListMemoriesForReviewResponse, error) {
	// 确保 page 和 page_size 合理
	page := req.Page
	if page < 1 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 10 // 默认值
	}

	// 调用 memory service 获取数据
	domainUnits, totalUnits, err := s.memoryService.ListMemoriesForReview(ctx, page, pageSize, ToEntityMemoryUnitTypes(req.Types), nil)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to list memories for review")
	}

	// 转换结果为 PB 格式
	total := uint32(totalUnits)

	return &pb.LearningServiceListMemoriesForReviewResponse{
		MemoryUnits: s.converter.ToPBMemoryUnits(domainUnits),
		Total:       total,
	}, nil
}

// ReviewMemoryUnits 复习记忆单元
func (s *LearningService) ReviewMemoryUnits(ctx context.Context, req *pb.LearningServiceReviewMemoryUnitsRequest) (*pb.LearningServiceReviewMemoryUnitsResponse, error) {
	// 从 context 获取 userID
	_, err := GetUserID(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid user context")
	}

	// 处理每个复习项
	updatedUnits := make([]*entity.MemoryUnit, len(req.Items))
	for i, item := range req.Items {
		// 获取记忆单元
		memoryUnit, err := s.memoryService.GetMemoryUnit(ctx, item.MemoryUnitId)
		if err != nil {
			return nil, status.Error(codes.NotFound, "memory unit not found")
		}

		// 更新记忆单元状态
		now := time.Now()
		memoryUnit.LastReviewAt = now
		memoryUnit.ReviewCount++
		if item.Result == pb.ReviewResult_REVIEW_RESULT_CORRECT {
			memoryUnit.ConsecutiveCorrect++
			memoryUnit.ConsecutiveWrong = 0
		} else {
			memoryUnit.ConsecutiveCorrect = 0
			memoryUnit.ConsecutiveWrong++
		}

		// 更新学习时长
		memoryUnit.StudyDuration += item.ReviewDuration

		// 计算下次复习时间
		interval := s.memoryService.CalculateNextReviewInterval(memoryUnit)
		nextReviewAt := now.Add(time.Duration(interval.Days)*24*time.Hour +
			time.Duration(interval.Hours)*time.Hour +
			time.Duration(interval.Minutes)*time.Minute)
		memoryUnit.NextReviewAt = nextReviewAt

		// 保存更新
		if err := s.memoryService.UpdateMemoryUnit(ctx, memoryUnit); err != nil {
			return nil, status.Error(codes.Internal, "failed to update memory unit")
		}

		updatedUnits[i] = memoryUnit
	}

	return &pb.LearningServiceReviewMemoryUnitsResponse{}, nil
}

// ToEntityMemoryUnitTypes 辅助函数: 将 PB 类型列表转换为领域实体类型列表
func ToEntityMemoryUnitTypes(pbTypes []pb.MemoryUnitType) []entity.MemoryUnitType {
	if len(pbTypes) == 0 {
		return nil
	}
	entityTypes := make([]entity.MemoryUnitType, len(pbTypes))
	for i, t := range pbTypes {
		entityTypes[i] = entity.MemoryUnitType(t)
	}
	return entityTypes
}

// GetMemoryStats 获取记忆单元学习统计
func (s *LearningService) GetMemoryStats(ctx context.Context, req *pb.LearningServiceGetMemoryStatsRequest) (*pb.LearningServiceGetMemoryStatsResponse, error) {
	log := logger.GetLogger(ctx)
	log.Info("GetMemoryStats called", zap.Any("request", req))

	// 处理可选的类型过滤
	var unitType *entity.MemoryUnitType
	if req.Type != nil {
		t := s.converter.ToEntityMemoryUnitType(req.GetType())
		unitType = &t
	}

	// 调用 memory service 获取统计信息
	stats, err := s.memoryService.GetMemoryStats(ctx, unitType)
	if err != nil {
		log.Error("Failed to get memory stats from service", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to get memory stats")
	}

	if stats == nil {
		log.Warn("Memory service returned nil stats")
		return &pb.LearningServiceGetMemoryStatsResponse{}, nil
	}

	// 将领域统计信息映射到 PB 响应
	pbResponse := &pb.LearningServiceGetMemoryStatsResponse{
		TotalLearned:    uint32(stats.TotalCount),
		MasteredCount:   uint32(stats.MasteredCount),
		NeedReviewCount: uint32(stats.LearningCount + stats.NewCount),
		TotalStudyTime:  0,
		LevelStats:      make(map[uint32]uint32),
		RetentionRates:  make(map[uint32]float32),
	}

	// 填充 level_stats
	pbResponse.LevelStats[uint32(entity.MasteryLevelMastered)] = uint32(stats.MasteredCount)
	pbResponse.LevelStats[uint32(entity.MasteryLevelBeginner)] = 0
	pbResponse.LevelStats[uint32(entity.MasteryLevelFamiliar)] = 0
	if stats.LearningCount > 0 {
		pbResponse.LevelStats[uint32(entity.MasteryLevelBeginner)] = uint32(stats.LearningCount / 2)
		pbResponse.LevelStats[uint32(entity.MasteryLevelFamiliar)] = uint32(stats.LearningCount) - pbResponse.LevelStats[uint32(entity.MasteryLevelBeginner)]
	}
	pbResponse.LevelStats[uint32(entity.MasteryLevelUnlearned)] = uint32(stats.NewCount)

	// 填充 retention_rates
	if stats.TotalCount > 0 && stats.RetentionRate > 0 {
		if unitType != nil {
			pbResponse.RetentionRates[uint32(*unitType)] = float32(stats.RetentionRate)
		} else {
			pbResponse.RetentionRates[0] = float32(stats.RetentionRate)
		}
	}

	log.Info("GetMemoryStats successful", zap.Any("response", pbResponse))
	return pbResponse, nil
}

// InitializeMemoryUnits 初始化记忆单元
func (s *LearningService) InitializeMemoryUnits(ctx context.Context, req *pb.LearningServiceInitializeMemoryUnitRequest) (*pb.LearningServiceInitializeMemoryUnitResponse, error) {
	// 从 context 获取 userID
	userID, err := GetUserID(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid user context")
	}

	// 创建记忆单元列表
	memoryUnits := make([]*entity.MemoryUnit, len(req.Items))
	for i, item := range req.Items {
		memoryUnit := entity.NewMemoryUnit(entity.UID(userID), entity.MemoryUnitType(item.Type), item.ContentId)
		memoryUnit.MasteryLevel = entity.MasteryLevel(item.MasteryLevel)
		memoryUnits[i] = memoryUnit
	}

	// 批量保存记忆单元
	memoryUnitIDs, err := s.memoryService.CreateMemoryUnits(ctx, memoryUnits)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create memory units")
	}

	return &pb.LearningServiceInitializeMemoryUnitResponse{
		MemoryUnitIds: memoryUnitIDs,
	}, nil
}

// GetUserID 从 context 中获取用户ID
func GetUserID(ctx context.Context) (uint, error) {
	// TODO: 从 context 中获取用户ID
	// 这里需要根据实际的认证机制来实现
	return 0, nil
}
