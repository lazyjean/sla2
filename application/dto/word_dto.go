package dto

import (
	"github.com/lazyjean/sla2/domain/entity"
)

// WordCreateDTO 创建单词的数据传输对象
type WordCreateDTO struct {
	Text        string   `json:"text" binding:"required" example:"hello"`
	Translation string   `json:"translation" binding:"required" example:"你好"`
	Phonetic    string   `json:"phonetic" example:"həˈləʊ"`
	Examples    []string `json:"examples" example:"Hello, world!"`
	Tags        []string `json:"tags" example:"common,greeting"`
}

// WordResponseDTO 单词响应的数据传输对象
type WordResponseDTO struct {
	ID          uint     `json:"id" example:"1"`
	Text        string   `json:"text" example:"hello"`
	Translation string   `json:"translation" example:"你好"`
	Phonetic    string   `json:"phonetic" example:"həˈləʊ"`
	Examples    []string `json:"examples" example:"Hello, world!"`
	Tags        []string `json:"tags" example:"common,greeting"`
	CreatedAt   string   `json:"created_at" example:"2025-01-26 18:00:00"`
	UpdatedAt   string   `json:"updated_at" example:"2025-01-26 18:00:00"`
}

// ToEntity 将DTO转换为领域实体
func (dto *WordCreateDTO) ToEntity() (*entity.Word, error) {
	word, err := entity.NewWord(dto.Text, dto.Translation)
	if err != nil {
		return nil, err
	}

	word.Phonetic = dto.Phonetic

	// 添加例句
	for _, example := range dto.Examples {
		if err := word.AddExample(example); err != nil {
			return nil, err
		}
	}

	// 添加标签
	for _, tag := range dto.Tags {
		if err := word.AddTag(tag); err != nil {
			return nil, err
		}
	}

	return word, nil
}

// FromEntity 从领域实体转换为DTO
func WordResponseDTOFromEntity(word *entity.Word) *WordResponseDTO {
	examples := make([]string, len(word.Examples))
	for i, example := range word.Examples {
		examples[i] = example.Text
	}

	tags := make([]string, len(word.Tags))
	for i, tag := range word.Tags {
		tags[i] = tag.Name
	}

	return &WordResponseDTO{
		ID:          word.ID,
		Text:        word.Text,
		Translation: word.Translation,
		Phonetic:    word.Phonetic,
		Examples:    examples,
		Tags:        tags,
		CreatedAt:   word.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   word.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
