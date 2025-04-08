package dto

import (
	"github.com/lazyjean/sla2/internal/domain/entity"
	"github.com/lazyjean/sla2/internal/domain/valueobject"
)

// CreateWordRequest 创建单词的请求数据
type CreateWordRequest struct {
	Text        string
	Phonetic    string
	Definitions []entity.Definition
	Examples    []string
	Tags        []string
	Level       valueobject.WordDifficultyLevel
}

// WordCreateDTO 创建单词的请求数据
type WordCreateDTO struct {
	Text        string
	Definitions []entity.Definition
	Phonetic    string
	Examples    []string
	Tags        []string
	Level       valueobject.WordDifficultyLevel
}

// WordResponseDTO 单词响应的数据传输对象
type WordResponseDTO struct {
	ID          uint32
	Text        string
	Definitions []entity.Definition
	Phonetic    string
	Examples    []string
	Tags        []string
	CreatedAt   string
	UpdatedAt   string
}

// BatchCreateWordRequest 批量创建单词请求
type BatchCreateWordRequest struct {
	Word        string
	Definitions []struct {
		PartOfSpeech string
		Meaning      string
		Example      string
		Synonyms     []string
		Antonyms     []string
	}
	Level    valueobject.WordDifficultyLevel
	Tags     []string
	Examples []string
}

// ToEntity 将DTO转换为领域实体
func (dto *WordCreateDTO) ToEntity() *entity.Word {
	// 转换释义
	var definitions []entity.Definition
	for _, def := range dto.Definitions {
		definitions = append(definitions, entity.Definition{
			PartOfSpeech: def.PartOfSpeech,
			Meaning:      def.Meaning,
			Example:      def.Example,
			Synonyms:     def.Synonyms,
			Antonyms:     def.Antonyms,
		})
	}

	return entity.NewWord(
		dto.Text,
		dto.Phonetic,
		definitions,
		dto.Examples,
		dto.Tags,
		dto.Level,
	)
}

// FromEntity 从领域实体转换为DTO
func WordResponseDTOFromEntity(word *entity.Word) *WordResponseDTO {
	return &WordResponseDTO{
		ID:          uint32(word.ID),
		Text:        word.Text,
		Definitions: word.Definitions,
		Phonetic:    word.Phonetic,
		Examples:    word.Examples,
		Tags:        word.Tags,
		CreatedAt:   word.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   word.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
