package entity

import (
	"time"

	"github.com/lazyjean/sla2/internal/domain/valueobject"
	"gorm.io/gorm"
)

// HanCharID 汉字ID类型
type HanCharID uint32

// HanChar 汉字实体
type HanChar struct {
	// ID 唯一标识符
	ID HanCharID `gorm:"primaryKey;autoIncrement;comment:唯一标识符"`
	// Character 汉字字符
	Character string `gorm:"type:varchar(10);not null;uniqueIndex;comment:汉字字符"`
	// Pinyin 拼音
	Pinyin string `gorm:"type:varchar(50);not null;comment:拼音"`
	// Tags 标签列表
	Tags []string `gorm:"type:jsonb;not null;comment:标签列表"`
	// Categories 分类列表
	Categories []string `gorm:"type:jsonb;not null;comment:分类列表"`
	// Examples 例句列表
	Examples []string `gorm:"type:jsonb;not null;comment:例句列表"`
	// Level 难度等级
	Level valueobject.WordDifficultyLevel `gorm:"type:integer;not null;comment:难度等级"`
	// CreatedAt 创建时间
	CreatedAt time.Time `gorm:"type:timestamp with time zone;not null;comment:创建时间"`
	// UpdatedAt 更新时间
	UpdatedAt time.Time `gorm:"type:timestamp with time zone;not null;comment:更新时间"`
	// DeletedAt 删除时间
	DeletedAt gorm.DeletedAt `gorm:"index;comment:删除时间"`
}

// NewHanChar 创建新的汉字实体
func NewHanChar(character, pinyin string, level valueobject.WordDifficultyLevel) *HanChar {
	now := time.Now()
	return &HanChar{
		Character: character,
		Pinyin:    pinyin,
		Level:     level,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// AddTag 添加标签
func (h *HanChar) AddTag(tag string) {
	for _, t := range h.Tags {
		if t == tag {
			return
		}
	}
	h.Tags = append(h.Tags, tag)
}

// RemoveTag 移除标签
func (h *HanChar) RemoveTag(tag string) {
	for i, t := range h.Tags {
		if t == tag {
			h.Tags = append(h.Tags[:i], h.Tags[i+1:]...)
			return
		}
	}
}

// AddCategory 添加分类
func (h *HanChar) AddCategory(category string) {
	for _, c := range h.Categories {
		if c == category {
			return
		}
	}
	h.Categories = append(h.Categories, category)
}

// RemoveCategory 移除分类
func (h *HanChar) RemoveCategory(category string) {
	for i, c := range h.Categories {
		if c == category {
			h.Categories = append(h.Categories[:i], h.Categories[i+1:]...)
			return
		}
	}
}

// AddExample 添加例句
func (h *HanChar) AddExample(example string) {
	for _, e := range h.Examples {
		if e == example {
			return
		}
	}
	h.Examples = append(h.Examples, example)
}

// RemoveExample 移除例句
func (h *HanChar) RemoveExample(example string) {
	for i, e := range h.Examples {
		if e == example {
			h.Examples = append(h.Examples[:i], h.Examples[i+1:]...)
			return
		}
	}
}

// Update 更新汉字信息
func (h *HanChar) Update(character, pinyin string, level valueobject.WordDifficultyLevel) {
	h.Character = character
	h.Pinyin = pinyin
	h.Level = level
	h.UpdatedAt = time.Now()
}
