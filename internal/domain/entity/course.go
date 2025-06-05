package entity

import (
	"time"

	"gorm.io/gorm"
)

type CourseID uint32

// CourseCategory 课程分类
type CourseCategory string

const (
	CourseCategoryUnspecified CourseCategory = "unspecified" // 未指定
	CourseCategoryEnglish     CourseCategory = "english"     // 英语类
	CourseCategoryChinese     CourseCategory = "chinese"     // 汉语类
	CourseCategoryOther       CourseCategory = "other"       // 其他类型
)

// Course 课程实体
type Course struct {
	ID             CourseID         `gorm:"primaryKey"`
	Title          string           `gorm:"type:varchar(255);not null"`
	Description    string           `gorm:"type:text"`
	CoverURL       string           `gorm:"type:varchar(255)"`
	Level          string           `gorm:"type:varchar(50);not null"`
	Category       CourseCategory   `gorm:"type:varchar(50);not null;default:'unspecified'"` // 课程分类
	Tags           []string         `gorm:"type:jsonb;serializer:json;not null;default:'[]'"`
	Status         string           `gorm:"type:varchar(50);not null;default:'draft'"`
	Prompt         string           `gorm:"type:text"`                                        // AI 提示词
	Resources      []string         `gorm:"type:jsonb;serializer:json;not null;default:'[]'"` // 推荐学习资源 URL 列表
	RecommendedAge string           `gorm:"type:varchar(50)"`                                 // 推荐年龄范围，例如："7-12岁"
	StudyPlan      string           `gorm:"type:text"`                                        // 建议学习计划，包括学习时长、频率等建议
	CreatedAt      time.Time        `gorm:"type:timestamptz;not null"`
	UpdatedAt      time.Time        `gorm:"type:timestamptz;not null"`
	Sections       []*CourseSection `gorm:"-" json:"sections,omitempty"` // 课程章节列表，不存储在数据库中
	DeletedAt      gorm.DeletedAt   `gorm:"index"`
}

// TableName 指定表名
func (Course) TableName() string {
	return "courses"
}

// GetID 获取ID
func (c *Course) GetID() CourseID {
	return c.ID
}

// SetID 设置ID
func (c *Course) SetID(id CourseID) {
	c.ID = id
}

// GetCreatedAt 获取创建时间
func (c *Course) GetCreatedAt() time.Time {
	return c.CreatedAt
}

// SetCreatedAt 设置创建时间
func (c *Course) SetCreatedAt(t time.Time) {
	c.CreatedAt = t
}

// GetUpdatedAt 获取更新时间
func (c *Course) GetUpdatedAt() time.Time {
	return c.UpdatedAt
}

// SetUpdatedAt 设置更新时间
func (c *Course) SetUpdatedAt(t time.Time) {
	c.UpdatedAt = t
}

// GetDeletedAt 获取删除时间
func (c *Course) GetDeletedAt() gorm.DeletedAt {
	return c.DeletedAt
}

// SetDeletedAt 设置删除时间
func (c *Course) SetDeletedAt(t gorm.DeletedAt) {
	c.DeletedAt = t
}

// BatchCreateCourseInput 批量创建课程的输入结构
type BatchCreateCourseInput struct {
	Title          string
	Description    string
	CoverURL       string
	Level          string
	Category       CourseCategory
	Tags           []string
	Prompt         string
	Resources      []string
	RecommendedAge string
	StudyPlan      string
	Sections       []BatchCreateSectionInput
}

// BatchCreateSectionInput 批量创建章节的输入结构
type BatchCreateSectionInput struct {
	Title      string
	Desc       string
	OrderIndex int32
	Units      []BatchCreateUnitInput
}

// BatchCreateUnitInput 批量创建单元的输入结构
type BatchCreateUnitInput struct {
	Title       string
	Desc        string
	QuestionIds []uint32
	OrderIndex  int32
	Tags        []string
	Prompt      string
}
