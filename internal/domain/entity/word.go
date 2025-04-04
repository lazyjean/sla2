package entity

import (
	"time"

	"github.com/lazyjean/sla2/internal/domain/errors"
)

type WordID uint32

// Word 生词实体
type Word struct {
	// ID 单词唯一标识
	ID WordID `gorm:"primaryKey;autoIncrement"`
	// UserID 用户ID
	UserID UID `gorm:"not null;index;uniqueIndex:idx_user_text,priority:1"`
	// Text 单词文本
	Text string `gorm:"type:varchar(100);not null;index;uniqueIndex:idx_user_text,priority:2"`
	// Phonetic 音标
	Phonetic string `gorm:"type:varchar(100)"`
	// Definitions 释义列表
	Definitions []Definition `gorm:"type:jsonb;serializer:json;not null;default:'[]'"`
	// Examples 例句列表
	Examples []string `gorm:"type:jsonb;serializer:json;not null;default:'[]'"`
	// Tags 标签列表
	Tags []string `gorm:"type:jsonb;serializer:json;not null;default:'[]'"`
	// Difficulty 难度等级（1-5）
	Difficulty int `gorm:"type:int;not null;default:0"`
	// MasteryLevel 掌握程度（0-5）
	MasteryLevel int `gorm:"default:0"`
	// NextReviewAt 下次复习时间
	NextReviewAt *time.Time
	// ReviewCount 复习次数
	ReviewCount int `gorm:"default:0"`
	// CreatedAt 创建时间
	CreatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
	// UpdatedAt 更新时间
	UpdatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
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

// NewWord 创建新生词
func NewWord(userID UID, text, phonetic string, definitions []Definition, examples, tags []string) (*Word, error) {
	if userID == 0 {
		return nil, errors.ErrInvalidUserID
	}
	if text == "" {
		return nil, errors.ErrInvalidWord
	}

	return &Word{
		UserID:      userID,
		Text:        text,
		Phonetic:    phonetic,
		Definitions: definitions,
		Examples:    examples,
		Tags:        tags,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil
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

// UpdateMastery 更新掌握度
func (w *Word) UpdateMastery(level int) {
	w.MasteryLevel = level
	w.ReviewCount++

	// 根据掌握度计算下次复习时间
	var nextReview time.Time
	switch level {
	case 1:
		nextReview = time.Now().Add(24 * time.Hour)
	case 2:
		nextReview = time.Now().Add(3 * 24 * time.Hour)
	case 3:
		nextReview = time.Now().Add(7 * 24 * time.Hour)
	case 4:
		nextReview = time.Now().Add(14 * 24 * time.Hour)
	case 5:
		nextReview = time.Now().Add(30 * 24 * time.Hour)
	default:
		nextReview = time.Now().Add(12 * time.Hour)
	}
	w.NextReviewAt = &nextReview
	w.UpdatedAt = time.Now()
}

// Validate 验证单词数据
func (w *Word) Validate() error {
	if w.UserID == 0 {
		return errors.ErrInvalidUserID
	}
	if w.Text == "" {
		return errors.ErrInvalidWord
	}
	if len(w.Definitions) == 0 {
		return errors.ErrEmptyDefinition
	}
	return nil
}

// IsNeedReview 检查是否需要复习
func (w *Word) IsNeedReview() bool {
	if w.NextReviewAt == nil {
		return false
	}
	return w.NextReviewAt.Before(time.Now())
}
