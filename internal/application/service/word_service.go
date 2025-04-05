package service

import (
	"context"

	pb "github.com/lazyjean/sla2/api/proto/v1"
	"github.com/lazyjean/sla2/internal/application/dto"
	"github.com/lazyjean/sla2/internal/domain/entity"
	"github.com/lazyjean/sla2/internal/domain/repository"
)

type WordService struct {
	pb.UnimplementedWordServiceServer
	wordRepo repository.WordRepository
}

func NewWordService(wordRepo repository.WordRepository) *WordService {
	return &WordService{
		wordRepo: wordRepo,
	}
}

// CreateWord 创建生词
func (s *WordService) CreateWord(ctx context.Context, createDTO *dto.WordCreateDTO, userID entity.UID) (*entity.Word, error) {
	word, err := createDTO.ToEntity(userID)
	if err != nil {
		return nil, err
	}

	if err := s.wordRepo.Create(ctx, word); err != nil {
		return nil, err
	}

	return word, nil
}

// GetWord 获取生词
func (s *WordService) GetWord(ctx context.Context, id uint) (*entity.Word, error) {
	return s.wordRepo.GetByID(ctx, entity.WordID(id))
}

// ListWords 获取生词列表
func (s *WordService) ListWords(ctx context.Context, userID uint, offset, limit int) ([]*entity.Word, uint32, error) {
	filters := make(map[string]interface{})
	filters["user_id"] = userID

	words, total, err := s.wordRepo.List(ctx, offset, limit, filters)
	if err != nil {
		return nil, 0, err
	}

	return words, uint32(total), nil
}

// DeleteWord 删除生词
func (s *WordService) DeleteWord(ctx context.Context, id uint) error {
	return s.wordRepo.Delete(ctx, entity.WordID(id))
}

// SearchWords 搜索生词
func (s *WordService) SearchWords(ctx context.Context, keyword string, userID uint, offset, limit int, filters map[string]interface{}) ([]*dto.WordResponseDTO, uint32, error) {
	if filters == nil {
		filters = make(map[string]interface{})
	}
	filters["user_id"] = userID

	words, total, err := s.wordRepo.Search(ctx, keyword, offset, limit, filters)
	if err != nil {
		return nil, 0, err
	}

	// 转换为 DTO
	dtos := make([]*dto.WordResponseDTO, len(words))
	for i, word := range words {
		dtos[i] = dto.WordResponseDTOFromEntity(word)
	}

	return dtos, uint32(total), nil
}
