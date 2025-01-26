package entity

import (
	"time"

	"github.com/lazyjean/sla2/domain/errors"
	"gorm.io/gorm"
)

// Word 单词实体
type Word struct {
	// ID 单词唯一标识
	ID uint `gorm:"primaryKey;autoIncrement"`
	// Text 单词文本
	Text string `gorm:"type:varchar(100);not null;index"`
	// Phonetic 音标
	Phonetic string `gorm:"type:varchar(100)"`
	// Translation 翻译
	Translation string `gorm:"type:text;not null"`
	// Examples 例句列表
	Examples []Example `gorm:"foreignKey:WordID;constraint:OnDelete:CASCADE"`
	// Tags 标签列表
	Tags []Tag `gorm:"many2many:word_tags;constraint:OnDelete:CASCADE"`
	// Difficulty 难度等级（1-5）
	Difficulty int `gorm:"type:int;default:1"`
	// CreatedAt 创建时间
	CreatedAt time.Time
	// UpdatedAt 更新时间
	UpdatedAt time.Time
	// DeletedAt 删除时间（软删除）
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// Example 例句实体
type Example struct {
	// ID 例句唯一标识
	ID uint `gorm:"primaryKey"`
	// WordID 关联的单词ID
	WordID uint `gorm:"not null"`
	// Text 例句文本
	Text string `gorm:"type:text;not null"`
	// CreatedAt 创建时间
	CreatedAt time.Time
}

// Tag 标签实体
type Tag struct {
	// ID 标签唯一标识
	ID uint `gorm:"primaryKey"`
	// Name 标签名称
	Name string `gorm:"type:varchar(50);uniqueIndex"`
	// CreatedAt 创建时间
	CreatedAt time.Time
}

type WordTag struct {
	// WordID 关联的单词ID
	WordID uint `gorm:"not null"`
	// TagID 关联的标签ID
	TagID uint `gorm:"not null"`
}

// NewWord 创建新单词
func NewWord(text string, translation string) (*Word, error) {
	if text == "" {
		return nil, errors.ErrEmptyWordText
	}
	if translation == "" {
		return nil, errors.ErrEmptyTranslation
	}

	return &Word{
		Text:        text,
		Translation: translation,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil
}

// AddExample 添加例句
func (w *Word) AddExample(example string) error {
	if example == "" {
		return errors.ErrEmptyExample
	}
	w.Examples = append(w.Examples, Example{Text: example, WordID: w.ID})
	w.UpdatedAt = time.Now()
	return nil
}

// AddTag 添加标签
func (w *Word) AddTag(tag string) error {
	if tag == "" {
		return errors.ErrEmptyTag
	}
	w.Tags = append(w.Tags, Tag{Name: tag})
	w.UpdatedAt = time.Now()
	return nil
}
