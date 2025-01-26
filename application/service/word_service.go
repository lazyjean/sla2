package service

import (
	"context"

	"github.com/lazyjean/sla2/application/dto"
	"github.com/lazyjean/sla2/domain/repository"
)

type WordService struct {
	wordRepo repository.WordRepository
}

func NewWordService(wordRepo repository.WordRepository) *WordService {
	return &WordService{wordRepo: wordRepo}
}

func (s *WordService) CreateWord(ctx context.Context, createDTO *dto.WordCreateDTO) (*dto.WordResponseDTO, error) {
	word, err := createDTO.ToEntity()
	if err != nil {
		return nil, err
	}

	if err := s.wordRepo.Save(ctx, word); err != nil {
		return nil, err
	}

	return dto.WordResponseDTOFromEntity(word), nil
}

func (s *WordService) GetWord(ctx context.Context, id uint) (*dto.WordResponseDTO, error) {
	word, err := s.wordRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return dto.WordResponseDTOFromEntity(word), nil
}

func (s *WordService) UpdateWord(ctx context.Context, id uint, updateDTO *dto.WordCreateDTO) (*dto.WordResponseDTO, error) {
	word, err := updateDTO.ToEntity()
	if err != nil {
		return nil, err
	}

	word.ID = id
	if err := s.wordRepo.Update(ctx, word); err != nil {
		return nil, err
	}

	return dto.WordResponseDTOFromEntity(word), nil
}

func (s *WordService) DeleteWord(ctx context.Context, id uint) error {
	return s.wordRepo.Delete(ctx, id)
}

func (s *WordService) ListWords(ctx context.Context, page, pageSize int) ([]*dto.WordResponseDTO, int64, error) {
	offset := (page - 1) * pageSize
	words, total, err := s.wordRepo.List(ctx, offset, pageSize)
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
