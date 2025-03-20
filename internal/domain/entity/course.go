package entity

import (
	"time"
)

type CourseID uint32

// Course 课程实体
type Course struct {
	ID             CourseID         `gorm:"primaryKey"`
	Title          string           `gorm:"type:varchar(255);not null"`
	Description    string           `gorm:"type:text"`
	CoverURL       string           `gorm:"type:varchar(255)"`
	Level          string           `gorm:"type:varchar(50);not null"`
	Tags           []string         `gorm:"type:jsonb;serializer:json;not null;default:'[]'"`
	Status         string           `gorm:"type:varchar(50);not null;default:'draft'"`
	Prompt         string           `gorm:"type:text"`                                        // AI 提示词
	Resources      []string         `gorm:"type:jsonb;serializer:json;not null;default:'[]'"` // 推荐学习资源 URL 列表
	RecommendedAge string           `gorm:"type:varchar(50)"`                                 // 推荐年龄范围，例如："7-12岁"
	StudyPlan      string           `gorm:"type:text"`                                        // 建议学习计划，包括学习时长、频率等建议
	CreatedAt      time.Time        `gorm:"not null"`
	UpdatedAt      time.Time        `gorm:"not null"`
	Sections       []*CourseSection `gorm:"-" json:"sections,omitempty"` // 课程章节列表，不存储在数据库中
}

// TableName 指定表名
func (Course) TableName() string {
	return "courses"
}
