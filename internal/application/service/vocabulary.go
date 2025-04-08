package service

import (
	"context"
	"fmt"

	"github.com/lazyjean/sla2/internal/application/dto"
	"github.com/lazyjean/sla2/internal/domain/entity"
	"github.com/lazyjean/sla2/internal/domain/errors"
	"github.com/lazyjean/sla2/internal/domain/repository"
	"github.com/lazyjean/sla2/internal/domain/valueobject"
)

// VocabularyService 词汇服务
type VocabularyService struct {
	hanCharRepository repository.HanCharRepository
	wordRepository    repository.WordRepository
}

// NewVocabularyService 创建词汇服务实例
func NewVocabularyService(hanCharRepository repository.HanCharRepository, wordRepository repository.WordRepository) *VocabularyService {
	return &VocabularyService{
		hanCharRepository: hanCharRepository,
		wordRepository:    wordRepository,
	}
}

// CreateHanChar 创建汉字
func (s *VocabularyService) CreateHanChar(ctx context.Context, character, pinyin string, level valueobject.WordDifficultyLevel, tags, categories, examples []string) (*entity.HanChar, error) {
	// 检查汉字是否已存在
	existing, err := s.hanCharRepository.GetByCharacter(ctx, character)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.ErrHanCharAlreadyExists
	}

	// 创建新的汉字实体
	hanChar := entity.NewHanChar(character, pinyin, level)
	hanChar.Tags = tags
	hanChar.Categories = categories
	hanChar.Examples = examples

	// 保存到数据库
	id, err := s.hanCharRepository.Create(ctx, hanChar)
	if err != nil {
		return nil, err
	}

	hanChar.ID = id
	return hanChar, nil
}

// UpdateHanChar 更新汉字
func (s *VocabularyService) UpdateHanChar(ctx context.Context, id entity.HanCharID, character, pinyin string, level valueobject.WordDifficultyLevel, tags, categories, examples []string) (*entity.HanChar, error) {
	// 获取现有汉字
	hanChar, err := s.hanCharRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 更新汉字信息
	hanChar.Update(character, pinyin, level)
	hanChar.Tags = tags
	hanChar.Categories = categories
	hanChar.Examples = examples

	// 保存更新
	if err := s.hanCharRepository.Update(ctx, hanChar); err != nil {
		return nil, err
	}

	return hanChar, nil
}

// DeleteHanChar 删除汉字
func (s *VocabularyService) DeleteHanChar(ctx context.Context, id entity.HanCharID) error {
	return s.hanCharRepository.Delete(ctx, id)
}

// GetHanChar 获取汉字详情
func (s *VocabularyService) GetHanChar(ctx context.Context, id entity.HanCharID) (*entity.HanChar, error) {
	return s.hanCharRepository.GetByID(ctx, id)
}

// GetHanCharByCharacter 根据字符获取汉字
func (s *VocabularyService) GetHanCharByCharacter(ctx context.Context, character string) (*entity.HanChar, error) {
	return s.hanCharRepository.GetByCharacter(ctx, character)
}

// ListHanChars 获取汉字列表
func (s *VocabularyService) ListHanChars(ctx context.Context, page, pageSize int, level valueobject.WordDifficultyLevel, tags, categories []string) ([]*entity.HanChar, int64, error) {
	offset := (page - 1) * pageSize

	filters := make(map[string]interface{})
	if level != valueobject.WORD_DIFFICULTY_LEVEL_UNSPECIFIED {
		filters["level"] = level
	}
	if len(tags) > 0 {
		filters["tags"] = tags
	}
	if len(categories) > 0 {
		filters["categories"] = categories
	}

	return s.hanCharRepository.List(ctx, offset, pageSize, filters)
}

// SearchHanChars 搜索汉字
func (s *VocabularyService) SearchHanChars(ctx context.Context, keyword string, page, pageSize int, level valueobject.WordDifficultyLevel, tags, categories []string) ([]*entity.HanChar, int64, error) {
	offset := (page - 1) * pageSize

	filters := make(map[string]interface{})
	if level != valueobject.WORD_DIFFICULTY_LEVEL_UNSPECIFIED {
		filters["level"] = level
	}
	if len(tags) > 0 {
		filters["tags"] = tags
	}
	if len(categories) > 0 {
		filters["categories"] = categories
	}

	return s.hanCharRepository.Search(ctx, keyword, offset, pageSize, filters)
}

// GetWord 获取单词详情
func (s *VocabularyService) GetWord(ctx context.Context, id uint) (*entity.Word, error) {
	return s.wordRepository.GetByID(ctx, entity.WordID(id))
}

// ListWords 获取单词列表
func (s *VocabularyService) ListWords(ctx context.Context, page, pageSize int, level valueobject.WordDifficultyLevel, tags, categories []string) ([]*entity.Word, int64, error) {
	offset := (page - 1) * pageSize

	filters := make(map[string]interface{})
	if level != valueobject.WORD_DIFFICULTY_LEVEL_UNSPECIFIED {
		filters["level"] = level
	}
	if len(tags) > 0 {
		filters["tags"] = tags
	}
	if len(categories) > 0 {
		filters["categories"] = categories
	}

	return s.wordRepository.List(ctx, offset, pageSize, filters)
}

// GetAllMetadata 获取所有标签和分类信息
func (s *VocabularyService) GetAllMetadata(ctx context.Context) ([]string, []string, error) {
	tags, err := s.wordRepository.GetAllTags(ctx)
	if err != nil {
		return nil, nil, err
	}

	categories, err := s.wordRepository.GetAllCategories(ctx)
	if err != nil {
		return nil, nil, err
	}

	return tags, categories, nil
}

// CreateWord 创建单词
func (s *VocabularyService) CreateWord(ctx context.Context, req dto.CreateWordRequest) (*entity.Word, error) {
	// 检查单词是否已存在
	existingWord, err := s.wordRepository.GetByWord(ctx, req.Text)
	if err != nil && !errors.Is(err, errors.ErrWordNotFound) {
		return nil, fmt.Errorf("failed to check word existence: %w", err)
	}
	if existingWord != nil {
		return nil, errors.ErrWordAlreadyExists
	}

	// 创建新单词
	word := entity.NewWord(
		req.Text,
		req.Phonetic,
		req.Definitions,
		req.Examples,
		req.Tags,
		req.Level,
	)

	// 保存单词
	if err := s.wordRepository.Create(ctx, word); err != nil {
		return nil, fmt.Errorf("failed to create word: %w", err)
	}

	return word, nil
}

// BatchCreateWords 批量创建单词
func (s *VocabularyService) BatchCreateWords(ctx context.Context, words []dto.BatchCreateWordRequest) error {
	for _, word := range words {
		// 解析难度等级
		level := word.Level

		// 创建单词请求
		req := dto.CreateWordRequest{
			Text:        word.Word,
			Definitions: make([]entity.Definition, len(word.Definitions)),
			Examples:    word.Examples,
			Tags:        word.Tags,
			Level:       level,
		}

		// 转换释义
		for i, def := range word.Definitions {
			req.Definitions[i] = entity.Definition{
				PartOfSpeech: def.PartOfSpeech,
				Meaning:      def.Meaning,
				Example:      def.Example,
				Synonyms:     def.Synonyms,
				Antonyms:     def.Antonyms,
			}
		}

		// 创建单词
		if _, err := s.CreateWord(ctx, req); err != nil {
			return fmt.Errorf("failed to create word %s: %w", word.Word, err)
		}
	}

	return nil
}

// BatchCreateHanChars 批量创建汉字
func (s *VocabularyService) BatchCreateHanChars(ctx context.Context, hanChars []struct {
	Character  string
	Pinyin     string
	Level      valueobject.WordDifficultyLevel
	Tags       []string
	Categories []string
	Examples   []string
}) ([]uint, error) {
	var ids []uint
	for _, hanChar := range hanChars {
		// 创建新的汉字实体
		newHanChar := entity.NewHanChar(hanChar.Character, hanChar.Pinyin, hanChar.Level)
		newHanChar.Tags = hanChar.Tags
		newHanChar.Categories = hanChar.Categories
		newHanChar.Examples = hanChar.Examples

		// 保存到数据库
		id, err := s.hanCharRepository.Create(ctx, newHanChar)
		if err != nil {
			return nil, err
		}

		ids = append(ids, uint(id))
	}
	return ids, nil
}
