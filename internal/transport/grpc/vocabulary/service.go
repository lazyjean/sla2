package vocabulary

import (
	"context"

	"github.com/lazyjean/sla2/internal/domain/entity"

	pb "github.com/lazyjean/sla2/api/proto/v1"
	"github.com/lazyjean/sla2/internal/application/dto"
	"github.com/lazyjean/sla2/internal/application/service"
	"github.com/lazyjean/sla2/internal/transport/grpc/vocabulary/converter"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Service 词汇服务实现
type Service struct {
	pb.UnimplementedVocabularyServiceServer
	service   service.VocabularyService
	converter *converter.VocabularyConverter
}

// NewVocabularyService 创建词汇服务实例
func NewVocabularyService(service service.VocabularyService) *Service {
	return &Service{
		service:   service,
		converter: converter.NewVocabularyConverter(),
	}
}

// Get 获取单词详情
func (s *Service) Get(ctx context.Context, req *pb.VocabularyServiceGetRequest) (*pb.VocabularyServiceGetResponse, error) {
	word, err := s.service.GetWord(ctx, uint(req.Id))
	if err != nil {
		return nil, err
	}

	return &pb.VocabularyServiceGetResponse{
		Word: s.converter.ToProtoWord(word),
	}, nil
}

// List 获取单词列表
func (s *Service) List(ctx context.Context, req *pb.VocabularyServiceListRequest) (*pb.VocabularyServiceListResponse, error) {
	words, total, err := s.service.ListWords(
		ctx,
		int(req.Page),
		int(req.PageSize),
		s.converter.ConvertLevelToValueObject(req.Level),
		req.Tags,
		req.Categories,
	)
	if err != nil {
		return nil, err
	}

	var pbWords []*pb.Word
	for _, word := range words {
		pbWords = append(pbWords, s.converter.ToProtoWord(word))
	}

	return &pb.VocabularyServiceListResponse{
		Words: pbWords,
		Total: uint32(total),
	}, nil
}

// GetAllMetadata 获取所有标签和分类信息
func (s *Service) GetAllMetadata(ctx context.Context, req *pb.VocabularyServiceGetAllMetadataRequest) (*pb.VocabularyServiceGetAllMetadataResponse, error) {
	tags, categories, err := s.service.GetAllMetadata(ctx)
	if err != nil {
		return nil, err
	}

	return &pb.VocabularyServiceGetAllMetadataResponse{
		Tags:       tags,
		Categories: categories,
	}, nil
}

// ListHanChar 获取汉字列表
func (s *Service) ListHanChar(ctx context.Context, req *pb.VocabularyServiceListHanCharRequest) (*pb.VocabularyServiceListHanCharResponse, error) {
	// 将 protobuf 枚举值转换为内部枚举值
	level := s.converter.ConvertLevelToValueObject(req.Level)

	request := dto.ListHanCharsDTO{
		Page:       int(req.Page),
		PageSize:   int(req.PageSize),
		Level:      level,
		Tags:       req.Tags,
		Categories: req.Categories,
	}

	hanChars, total, err := s.service.ListHanChars(
		ctx,
		request.Page,
		request.PageSize,
		request.Level,
		request.Tags,
		request.Categories,
		nil, // excludeIDs is optional, passing nil for now
	)
	if err != nil {
		return nil, err
	}

	var pbHanChars []*pb.HanChar
	for _, hanChar := range hanChars {
		pbHanChars = append(pbHanChars, s.converter.ToProtoHanChar(hanChar))
	}

	return &pb.VocabularyServiceListHanCharResponse{
		HanChars: pbHanChars,
		Total:    uint32(total),
	}, nil
}

// BatchCreate 批量创建英文单词
func (s *Service) BatchCreate(ctx context.Context, req *pb.VocabularyServiceBatchCreateRequest) (*pb.VocabularyServiceBatchCreateResponse, error) {
	words := s.converter.PbWordsToEntities(req.Words)
	err := s.service.BatchCreateWords(ctx, words)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.VocabularyServiceBatchCreateResponse{}, nil
}

// BatchCreateHanChar 批量创建汉字
func (s *Service) BatchCreateHanChar(ctx context.Context, req *pb.VocabularyServiceBatchCreateHanCharRequest) (*pb.VocabularyServiceBatchCreateHanCharResponse, error) {
	var hanChars []*entity.HanChar
	for _, hanCharPb := range req.HanChars {
		hanChars = append(hanChars, s.converter.PbToEntityHanChar(hanCharPb))
	}
	ids, err := s.service.BatchCreateHanChars(ctx, hanChars)
	if err != nil {
		return nil, err
	}

	// 转换 []uint 为 []uint32
	uint32Ids := make([]uint32, len(ids))
	for i, id := range ids {
		uint32Ids[i] = uint32(id)
	}

	return &pb.VocabularyServiceBatchCreateHanCharResponse{
		Ids: uint32Ids,
	}, nil
}
