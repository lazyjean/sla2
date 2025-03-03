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
	// Translation 翻译
	Translation string `gorm:"type:text;not null"`
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

// NewWord 创建新生词
func NewWord(userID UID, text string, phonetic string, translation string, examples []string, tags []string) (*Word, error) {
	if userID == 0 {
		return nil, errors.ErrInvalidUserID
	}
	if text == "" {
		return nil, errors.ErrEmptyWordText
	}
	if translation == "" {
		return nil, errors.ErrEmptyTranslation
	}

	// 过滤空字符串
	filteredExamples := make([]string, 0, len(examples))
	for _, ex := range examples {
		if ex != "" {
			filteredExamples = append(filteredExamples, ex)
		}
	}

	filteredTags := make([]string, 0, len(tags))
	for _, t := range tags {
		if t != "" {
			filteredTags = append(filteredTags, t)
		}
	}

	return &Word{
		UserID:       userID,
		Text:         text,
		Phonetic:     phonetic,
		Translation:  translation,
		Difficulty:   1,
		MasteryLevel: 0,
		ReviewCount:  0,
		Examples:     filteredExamples,
		Tags:         filteredTags,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}, nil
}

// AddExample 添加例句
func (w *Word) AddExample(example string) error {
	if example == "" {
		return errors.ErrEmptyExample
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
		return errors.ErrEmptyWordText
	}
	if w.Translation == "" {
		return errors.ErrEmptyTranslation
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
