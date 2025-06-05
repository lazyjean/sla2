package entity

import (
	"time"

	"gorm.io/gorm"
)

// CourseSectionID 章节ID类型
type CourseSectionID uint32

// CourseSectionUnitID 单元ID类型
type CourseSectionUnitID uint32

// CourseSection 课程章节实体
type CourseSection struct {
	ID         CourseSectionID      `gorm:"primaryKey;autoIncrement"`
	CourseID   CourseID             `gorm:"not null;index"`                              // 所属课程ID
	Title      string               `gorm:"type:varchar(100);not null"`                  // 章节标题
	Desc       string               `gorm:"type:text"`                                   // 章节描述
	OrderIndex int32                `gorm:"not null;default:0"`                          // 显示顺序
	Status     string               `gorm:"type:varchar(20);not null;default:'enabled'"` // 状态：enabled-启用，disabled-禁用
	Units      []*CourseSectionUnit `gorm:"-"`                                           // 章节单元列表，不存储在数据库中
	CreatedAt  time.Time            `gorm:"type:timestamptz;not null"`                   // 创建时间
	UpdatedAt  time.Time            `gorm:"type:timestamptz;not null"`                   // 更新时间
	DeletedAt  gorm.DeletedAt       `gorm:"index"`
}

// CourseSectionUnit 课程章节单元实体
type CourseSectionUnit struct {
	ID          CourseSectionUnitID `gorm:"primaryKey;autoIncrement"`
	SectionID   CourseSectionID     `gorm:"not null;index"`             // 所属章节ID
	Title       string              `gorm:"type:varchar(100);not null"` // 单元标题
	Desc        string              `gorm:"type:text"`                  // 单元内容
	QuestionIds string              `gorm:"type:text"`                  // 问题ID列表，存储关联题目的ID数组
	OrderIndex  int32               `gorm:"not null;default:0"`         // 显示顺序
	Status      int32               `gorm:"not null;default:1"`         // 状态：0-禁用，1-启用
	Tags        string              `gorm:"type:text"`                  // 标签，多个标签用逗号分隔
	Prompt      string              `gorm:"type:text"`                  // AI 提示词
	CreatedAt   time.Time           `gorm:"type:timestamptz;not null"`  // 创建时间
	UpdatedAt   time.Time           `gorm:"type:timestamptz;not null"`  // 更新时间
}

// TableName 指定表名
func (CourseSection) TableName() string {
	return "course_sections"
}

// TableName 指定表名
func (CourseSectionUnit) TableName() string {
	return "course_section_units"
}

// GetID 获取ID
func (s *CourseSection) GetID() CourseSectionID {
	return s.ID
}

// SetID 设置ID
func (s *CourseSection) SetID(id CourseSectionID) {
	s.ID = id
}

// GetCreatedAt 获取创建时间
func (s *CourseSection) GetCreatedAt() time.Time {
	return s.CreatedAt
}

// SetCreatedAt 设置创建时间
func (s *CourseSection) SetCreatedAt(t time.Time) {
	s.CreatedAt = t
}

// GetUpdatedAt 获取更新时间
func (s *CourseSection) GetUpdatedAt() time.Time {
	return s.UpdatedAt
}

// SetUpdatedAt 设置更新时间
func (s *CourseSection) SetUpdatedAt(t time.Time) {
	s.UpdatedAt = t
}

// GetDeletedAt 获取删除时间
func (s *CourseSection) GetDeletedAt() gorm.DeletedAt {
	return s.DeletedAt
}

// SetDeletedAt 设置删除时间
func (s *CourseSection) SetDeletedAt(t gorm.DeletedAt) {
	s.DeletedAt = t
}

// GetID 获取ID
func (u *CourseSectionUnit) GetID() CourseSectionUnitID {
	return u.ID
}

// SetID 设置ID
func (u *CourseSectionUnit) SetID(id CourseSectionUnitID) {
	u.ID = id
}

// GetCreatedAt 获取创建时间
func (u *CourseSectionUnit) GetCreatedAt() time.Time {
	return u.CreatedAt
}

// SetCreatedAt 设置创建时间
func (u *CourseSectionUnit) SetCreatedAt(t time.Time) {
	u.CreatedAt = t
}

// GetUpdatedAt 获取更新时间
func (u *CourseSectionUnit) GetUpdatedAt() time.Time {
	return u.UpdatedAt
}

// SetUpdatedAt 设置更新时间
func (u *CourseSectionUnit) SetUpdatedAt(t time.Time) {
	u.UpdatedAt = t
}

// GetDeletedAt 获取删除时间
func (u *CourseSectionUnit) GetDeletedAt() gorm.DeletedAt {
	return gorm.DeletedAt{}
}

// SetDeletedAt 设置删除时间
func (u *CourseSectionUnit) SetDeletedAt(t gorm.DeletedAt) {
	// CourseSectionUnit 不支持软删除
}
