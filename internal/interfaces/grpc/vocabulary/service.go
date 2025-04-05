package vocabulary

import (
	"context"

	pb "github.com/lazyjean/sla2/api/proto/v1"
	"github.com/lazyjean/sla2/internal/application/dto"
	"github.com/lazyjean/sla2/internal/application/service"
	"github.com/lazyjean/sla2/internal/domain/entity"
	"github.com/lazyjean/sla2/internal/domain/valueobject"
)

// VocabularyService 词汇服务实现
type VocabularyService struct {
	pb.UnimplementedVocabularyServiceServer
	service *service.VocabularyService
}

// NewVocabularyService 创建词汇服务实例
func NewVocabularyService(service *service.VocabularyService) *VocabularyService {
	return &VocabularyService{
		service: service,
	}
}

// Get 获取单词详情
func (s *VocabularyService) Get(ctx context.Context, req *pb.VocabularyServiceGetRequest) (*pb.VocabularyServiceGetResponse, error) {
	word, err := s.service.GetWord(ctx, uint(req.Id))
	if err != nil {
		return nil, err
	}

	return &pb.VocabularyServiceGetResponse{
		Word: ToProtoWord(word),
	}, nil
}

// List 获取单词列表
func (s *VocabularyService) List(ctx context.Context, req *pb.VocabularyServiceListRequest) (*pb.VocabularyServiceListResponse, error) {
	words, total, err := s.service.ListWords(
		ctx,
		int(req.Page),
		int(req.PageSize),
		ConvertLevelToString(req.Difficulty),
		req.Tags,
		req.Categories,
	)
	if err != nil {
		return nil, err
	}

	var pbWords []*pb.Word
	for _, word := range words {
		pbWords = append(pbWords, ToProtoWord(word))
	}

	return &pb.VocabularyServiceListResponse{
		Words: pbWords,
		Total: uint32(total),
	}, nil
}

// GetAllMetadata 获取所有标签和分类信息
func (s *VocabularyService) GetAllMetadata(ctx context.Context, req *pb.VocabularyServiceGetAllMetadataRequest) (*pb.VocabularyServiceGetAllMetadataResponse, error) {
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
func (s *VocabularyService) ListHanChar(ctx context.Context, req *pb.VocabularyServiceListHanCharRequest) (*pb.VocabularyServiceListHanCharResponse, error) {
	request := dto.ListHanCharsDTO{
		Page:       int(req.Page),
		PageSize:   int(req.PageSize),
		Level:      valueobject.WordDifficultyLevel(req.Level),
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
	)
	if err != nil {
		return nil, err
	}

	var pbHanChars []*pb.HanChar
	for _, hanChar := range hanChars {
		pbHanChars = append(pbHanChars, ToProtoHanChar(hanChar))
	}

	return &pb.VocabularyServiceListHanCharResponse{
		HanChars: pbHanChars,
		Total:    uint32(total),
	}, nil
}

// BatchCreate 批量创建英文单词
func (s *VocabularyService) BatchCreate(ctx context.Context, req *pb.VocabularyServiceBatchCreateRequest) (*pb.VocabularyServiceBatchCreateResponse, error) {
	var words []dto.BatchCreateWordRequest
	for _, wordPb := range req.Words {
		word := dto.BatchCreateWordRequest{
			Word:     wordPb.Word,
			Level:    ConvertLevelToString(wordPb.Difficulty),
			Tags:     wordPb.Tags,
			Examples: wordPb.Examples,
		}

		for _, def := range wordPb.Definitions {
			word.Definitions = append(word.Definitions, struct {
				PartOfSpeech string
				Meaning      string
				Example      string
				Synonyms     []string
				Antonyms     []string
			}{
				PartOfSpeech: def.PartOfSpeech.String(),
				Meaning:      def.Meaning,
				Example:      def.Example,
				Synonyms:     def.Synonyms,
				Antonyms:     def.Antonyms,
			})
		}

		words = append(words, word)
	}

	ids, err := s.service.BatchCreateWords(ctx, entity.UID(0), words)
	if err != nil {
		return nil, err
	}

	// 转换 []uint 为 []uint32
	uint32Ids := make([]uint32, len(ids))
	for i, id := range ids {
		uint32Ids[i] = uint32(id)
	}

	return &pb.VocabularyServiceBatchCreateResponse{
		Ids: uint32Ids,
	}, nil
}

// BatchCreateHanChar 批量创建汉字
func (s *VocabularyService) BatchCreateHanChar(ctx context.Context, req *pb.VocabularyServiceBatchCreateHanCharRequest) (*pb.VocabularyServiceBatchCreateHanCharResponse, error) {
	var hanChars []struct {
		Character  string
		Pinyin     string
		Level      string
		Tags       []string
		Categories []string
		Examples   []string
	}

	for _, hanCharPb := range req.HanChars {
		hanChars = append(hanChars, struct {
			Character  string
			Pinyin     string
			Level      string
			Tags       []string
			Categories []string
			Examples   []string
		}{
			Character:  hanCharPb.Character,
			Pinyin:     hanCharPb.Pinyin,
			Level:      ConvertLevelToString(hanCharPb.Level),
			Tags:       hanCharPb.Tags,
			Categories: hanCharPb.Categories,
			Examples:   hanCharPb.Examples,
		})
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
