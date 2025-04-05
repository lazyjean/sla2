package service

import (
	"context"

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
func (s *VocabularyService) ListWords(ctx context.Context, page, pageSize int, level string, tags, categories []string) ([]*entity.Word, int64, error) {
	offset := (page - 1) * pageSize

	filters := make(map[string]interface{})
	if level != "" {
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
func (s *VocabularyService) CreateWord(ctx context.Context, userID entity.UID, text, phonetic string, definitions []entity.Definition, examples, tags []string) (*entity.Word, error) {
	// 检查单词是否已存在
	existing, err := s.wordRepository.GetByWord(ctx, text)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.ErrWordAlreadyExists
	}

	// 创建新的单词实体
	word, err := entity.NewWord(userID, text, phonetic, definitions, examples, tags)
	if err != nil {
		return nil, err
	}

	// 保存到数据库
	if err := s.wordRepository.Create(ctx, word); err != nil {
		return nil, err
	}

	return word, nil
}

// BatchCreateWords 批量创建单词
func (s *VocabularyService) BatchCreateWords(ctx context.Context, userID entity.UID, words []dto.BatchCreateWordRequest) ([]uint, error) {
	var ids []uint
	for _, word := range words {
		// 转换释义
		var definitions []entity.Definition
		for _, def := range word.Definitions {
			definitions = append(definitions, entity.Definition{
				PartOfSpeech: def.PartOfSpeech,
				Meaning:      def.Meaning,
				Example:      def.Example,
				Synonyms:     def.Synonyms,
				Antonyms:     def.Antonyms,
			})
		}

		newWord, err := s.CreateWord(ctx, userID, word.Word, "", definitions, word.Examples, word.Tags)
		if err != nil {
			return nil, err
		}
		ids = append(ids, uint(newWord.ID))
	}
	return ids, nil
}

// BatchCreateHanChars 批量创建汉字
func (s *VocabularyService) BatchCreateHanChars(ctx context.Context, hanChars []struct {
	Character  string
	Pinyin     string
	Level      string
	Tags       []string
	Categories []string
	Examples   []string
}) ([]uint, error) {
	var ids []uint
	for _, hanChar := range hanChars {
		level, err := valueobject.ParseWordDifficultyLevel(hanChar.Level)
		if err != nil {
			return nil, err
		}

		// 创建新的汉字实体
		newHanChar := entity.NewHanChar(hanChar.Character, hanChar.Pinyin, level)
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
