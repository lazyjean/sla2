package entity

import (
	"time"

	"github.com/lazyjean/sla2/internal/domain/errors"
	"github.com/lazyjean/sla2/internal/domain/valueobject"
	"gorm.io/gorm"
)

type WordID uint32

// Word 单词实体
type Word struct {
	// ID 单词唯一标识
	ID WordID `gorm:"primaryKey;autoIncrement"`
	// Text 单词文本
	Text string `gorm:"type:varchar(100);not null;uniqueIndex"`
	// Phonetic 音标
	Phonetic string `gorm:"type:varchar(100)"`
	// Definitions 释义列表
	Definitions []Definition `gorm:"type:jsonb;serializer:json;not null;default:'[]'"`
	// Examples 例句列表
	Examples []string `gorm:"type:jsonb;serializer:json;not null;default:'[]'"`
	// Tags 标签列表
	Tags []string `gorm:"type:jsonb;serializer:json;not null;default:'[]'"`
	// Level 难度等级
	Level valueobject.WordDifficultyLevel `gorm:"type:integer;not null;comment:难度等级"`
	// CreatedAt 创建时间
	CreatedAt time.Time `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP"`
	// UpdatedAt 更新时间
	UpdatedAt time.Time `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP"`
}

// Definition 单词释义
type Definition struct {
	// PartOfSpeech 词性
	PartOfSpeech string `json:"part_of_speech"`
	// Meaning 含义
	Meaning string `json:"meaning"`
	// Example 例句
	Example string `json:"example"`
	// Synonyms 同义词
	Synonyms []string `json:"synonyms"`
	// Antonyms 反义词
	Antonyms []string `json:"antonyms"`
}

// NewWord 创建新的单词
func NewWord(
	text string,
	phonetic string,
	definitions []Definition,
	examples []string,
	tags []string,
	level valueobject.WordDifficultyLevel,
) *Word {
	return &Word{
		Text:        text,
		Phonetic:    phonetic,
		Definitions: definitions,
		Examples:    examples,
		Tags:        tags,
		Level:       level,
	}
}

// AddExample 添加例句
func (w *Word) AddExample(example string) error {
	if example == "" {
		return errors.ErrEmptyExample
	}
	// 检查例句是否已存在
	for _, e := range w.Examples {
		if e == example {
			return nil // 例句已存在，直接返回
		}
	}
	w.Examples = append(w.Examples, example)
	w.UpdatedAt = time.Now()
	return nil
}

// AddTag 添加标签
func (w *Word) AddTag(tagName string) error {
	if tagName == "" {
		return errors.ErrEmptyTag
	}
	// 检查标签是否已存在
	for _, tag := range w.Tags {
		if tag == tagName {
			return nil // 标签已存在，直接返回
		}
	}
	w.Tags = append(w.Tags, tagName)
	w.UpdatedAt = time.Now()
	return nil
}

// RemoveTag 移除标签
func (w *Word) RemoveTag(tagName string) {
	for i, tag := range w.Tags {
		if tag == tagName {
			w.Tags = append(w.Tags[:i], w.Tags[i+1:]...)
			w.UpdatedAt = time.Now()
			return
		}
	}
}

// HasTag 检查是否包含指定标签
func (w *Word) HasTag(tagName string) bool {
	for _, tag := range w.Tags {
		if tag == tagName {
			return true
		}
	}
	return false
}

// Validate 验证单词数据
func (w *Word) Validate() error {
	if w.Text == "" {
		return errors.ErrInvalidWord
	}
	if len(w.Definitions) == 0 {
		return errors.ErrEmptyDefinition
	}
	if !w.Level.IsValid() || !w.Level.IsCEFR() {
		return errors.ErrInvalidWord
	}
	return nil
}

// IsNeedReview 检查是否需要复习
func (w *Word) IsNeedReview() bool {
	return time.Now().After(w.UpdatedAt)
}

// Update 更新单词信息
func (w *Word) Update(
	text string,
	phonetic string,
	definitions []Definition,
	examples []string,
	tags []string,
	level valueobject.WordDifficultyLevel,
) {
	w.Text = text
	w.Phonetic = phonetic
	w.Definitions = definitions
	w.Examples = examples
	w.Tags = tags
	w.Level = level
	w.UpdatedAt = time.Now()
}

// GetID 获取ID
func (w *Word) GetID() WordID {
	return w.ID
}

// SetID 设置ID
func (w *Word) SetID(id WordID) {
	w.ID = id
}

// GetCreatedAt 获取创建时间
func (w *Word) GetCreatedAt() time.Time {
	return w.CreatedAt
}

// SetCreatedAt 设置创建时间
func (w *Word) SetCreatedAt(t time.Time) {
	w.CreatedAt = t
}

// GetUpdatedAt 获取更新时间
func (w *Word) GetUpdatedAt() time.Time {
	return w.UpdatedAt
}

// SetUpdatedAt 设置更新时间
func (w *Word) SetUpdatedAt(t time.Time) {
	w.UpdatedAt = t
}

// GetDeletedAt 获取删除时间
func (w *Word) GetDeletedAt() gorm.DeletedAt {
	return gorm.DeletedAt{}
}

// SetDeletedAt 设置删除时间
func (w *Word) SetDeletedAt(t gorm.DeletedAt) {
	// Word 实体不支持软删除，所以这个方法不做任何事
}
