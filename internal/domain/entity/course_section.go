package entity

import (
	"time"
)

// CourseSectionID 章节ID类型
type CourseSectionID uint32

// CourseSection 课程章节实体
type CourseSection struct {
	ID         CourseSectionID `gorm:"primaryKey;autoIncrement"`
	CourseID   CourseID        `gorm:"not null;index"`                              // 所属课程ID
	Title      string          `gorm:"type:varchar(100);not null"`                  // 章节标题
	Desc       string          `gorm:"type:text"`                                   // 章节描述
	OrderIndex int32           `gorm:"not null;default:0"`                          // 显示顺序
	Status     string          `gorm:"type:varchar(20);not null;default:'enabled'"` // 状态：enabled-启用，disabled-禁用
	CreatedAt  time.Time       `gorm:"not null"`                                    // 创建时间
	UpdatedAt  time.Time       `gorm:"not null"`                                    // 更新时间
}

// TableName 指定表名
func (CourseSection) TableName() string {
	return "course_sections"
}
