package service

import (
	"context"

	"github.com/lazyjean/sla2/internal/application/dto"
	"github.com/lazyjean/sla2/internal/domain/entity"
	"github.com/lazyjean/sla2/internal/domain/repository"
)

type WordService struct {
	wordRepo repository.WordRepository
}

func NewWordService(wordRepo repository.WordRepository) *WordService {
	return &WordService{
		wordRepo: wordRepo,
	}
}

// CreateWord 创建生词
func (s *WordService) CreateWord(ctx context.Context, createDTO *dto.WordCreateDTO, userID uint) (*entity.Word, error) {
	word, err := createDTO.ToEntity(userID)
	if err != nil {
		return nil, err
	}

	if err := s.wordRepo.Save(ctx, word); err != nil {
		return nil, err
	}

	return word, nil
}

// GetWord 获取生词
func (s *WordService) GetWord(ctx context.Context, id uint) (*entity.Word, error) {
	return s.wordRepo.FindByID(ctx, id)
}

// ListWords 获取生词列表
func (s *WordService) ListWords(ctx context.Context, userID uint, offset, limit int) ([]*entity.Word, int64, error) {
	return s.wordRepo.List(ctx, userID, offset, limit)
}

// DeleteWord 删除生词
func (s *WordService) DeleteWord(ctx context.Context, id uint) error {
	return s.wordRepo.Delete(ctx, id)
}

func (s *WordService) SearchWords(ctx context.Context, query *repository.WordQuery) ([]*dto.WordResponseDTO, int64, error) {
	words, total, err := s.wordRepo.Search(ctx, query)
	if err != nil {
		return nil, 0, err
	}

	// 转换为 DTO
	dtos := make([]*dto.WordResponseDTO, len(words))
	for i, word := range words {
		dtos[i] = dto.WordResponseDTOFromEntity(word)
	}

	return dtos, total, nil
}
