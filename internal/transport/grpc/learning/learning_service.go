package learning

import (
	"context"
	"time"

	"github.com/lazyjean/sla2/internal/transport/grpc/learning/converter"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/lazyjean/sla2/api/proto/v1"
	"github.com/lazyjean/sla2/internal/application/dto"
	"github.com/lazyjean/sla2/internal/application/service"
	"github.com/lazyjean/sla2/internal/domain/entity"
	vocabconverter "github.com/lazyjean/sla2/internal/transport/grpc/vocabulary/converter"
)

type Service struct {
	pb.UnimplementedLearningServiceServer
	learningService     *service.LearningService
	converter           *converter.LearningConverter
	vocabularyService   service.VocabularyService
	vocabularyConverter *vocabconverter.VocabularyConverter
}

func NewLearningService(learningService *service.LearningService, vocabularyService service.VocabularyService) *Service {
	return &Service{
		learningService:     learningService,
		converter:           converter.NewLearningConverter(),
		vocabularyService:   vocabularyService,
		vocabularyConverter: vocabconverter.NewVocabularyConverter(),
	}
}

// GetCourseProgress 获取课程学习进度
func (s *Service) GetCourseProgress(ctx context.Context, req *pb.LearningServiceGetCourseProgressRequest) (*pb.LearningServiceGetCourseProgressResponse, error) {
	courseID := req.CourseId

	// 获取课程进度及统计信息
	progress, completedSections, totalSections, err := s.learningService.GetCourseProgressWithStats(ctx, uint(courseID))
	if err != nil {
		return nil, err
	}

	return &pb.LearningServiceGetCourseProgressResponse{
		Progress: s.converter.ToPBLearningProgress(progress, completedSections, totalSections),
	}, nil
}

// GetSectionProgress 获取章节学习进度
func (s *Service) GetSectionProgress(ctx context.Context, req *pb.LearningServiceGetSectionProgressRequest) (*pb.LearningServiceGetSectionProgressResponse, error) {
	sectionID := req.SectionId

	// 获取章节进度及已完成的单元ID
	progress, completedUnits, totalUnits, completedUnitIDs, err := s.learningService.GetSectionProgressWithCompletedUnits(ctx, uint(sectionID))
	if err != nil {
		return nil, err
	}

	return &pb.LearningServiceGetSectionProgressResponse{
		Progress:         s.converter.ToPBLearningProgress(progress, completedUnits, totalUnits),
		CompletedUnitIds: completedUnitIDs,
	}, nil
}

// UpdateUnitProgress 更新单元学习进度
func (s *Service) UpdateUnitProgress(ctx context.Context, req *pb.LearningServiceUpdateUnitProgressRequest) (*pb.LearningServiceUpdateUnitProgressResponse, error) {
	err := s.learningService.UpdateUnitProgress(ctx, uint(req.UnitId), uint(req.SectionId), req.Completed)
	if err != nil {
		return nil, err
	}
	return &pb.LearningServiceUpdateUnitProgressResponse{}, nil
}

func (s *Service) InitializeMemoryUnits(ctx context.Context, req *pb.LearningServiceInitializeMemoryUnitRequest) (*pb.LearningServiceInitializeMemoryUnitResponse, error) {
	// 转换请求参数
	items := make([]*entity.MemoryUnitInitItem, 0, len(req.Items))
	for _, item := range req.Items {
		items = append(items, &entity.MemoryUnitInitItem{
			Type:         entity.MemoryUnitType(item.Type),
			ContentID:    item.ContentId,
			MasteryLevel: entity.MasteryLevel(item.MasteryLevel),
		})
	}

	// 调用应用服务层
	memoryUnitIDs, err := s.learningService.InitializeMemoryUnits(ctx, items)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.LearningServiceInitializeMemoryUnitResponse{
		MemoryUnitIds: memoryUnitIDs,
	}, nil
}

func (s *Service) ReviewMemoryUnits(ctx context.Context, req *pb.LearningServiceReviewMemoryUnitsRequest) (*pb.LearningServiceReviewMemoryUnitsResponse, error) {
	// 转换请求参数
	items := make([]*dto.ReviewItem, 0, len(req.Items))
	for _, item := range req.Items {
		items = append(items, &dto.ReviewItem{
			MemoryUnitID:   item.MemoryUnitId,
			Result:         item.Result == pb.ReviewResult_REVIEW_RESULT_CORRECT,
			ReviewDuration: time.Duration(item.ReviewDuration) * time.Second,
		})
	}

	// 调用应用服务层
	_, err := s.learningService.ReviewMemoryUnits(ctx, items)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.LearningServiceReviewMemoryUnitsResponse{}, nil
}

func (s *Service) GetMemoryStats(ctx context.Context, req *pb.LearningServiceGetMemoryStatsRequest) (*pb.LearningServiceGetMemoryStatsResponse, error) {
	// 转换请求参数
	unitType := s.converter.ToEntityMemoryUnitType(req.Type)

	// 调用应用服务层
	stats, err := s.learningService.GetMemoryStats(ctx, unitType)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.LearningServiceGetMemoryStatsResponse{
		TotalLearned:    uint32(stats.TotalCount),
		MasteredCount:   uint32(stats.MasteredCount),
		NeedReviewCount: uint32(stats.LearningCount),
		TotalStudyTime:  0,
		LevelStats:      map[uint32]uint32{},
		RetentionRates:  map[uint32]float32{},
	}, nil
}

func (s *Service) GetReviewContent(ctx context.Context, req *pb.GetReviewContentRequest) (*pb.GetReviewContentResponse, error) {
	hanChars, words, hasMore, err := s.learningService.GetReviewContent(ctx, int(req.Count))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var pbHanChars []*pb.HanChar
	for _, hanChar := range hanChars {
		pbHanChars = append(pbHanChars, s.vocabularyConverter.ToProtoHanChar(hanChar))
	}

	var pbWords []*pb.Word
	for _, word := range words {
		pbWords = append(pbWords, s.vocabularyConverter.ToProtoWord(word))
	}

	return &pb.GetReviewContentResponse{
		HanChars: pbHanChars,
		Words:    pbWords,
		HasMore:  hasMore,
	}, nil
}
