package entity

import "time"

// CourseLearningProgress 课程学习进度
type CourseLearningProgress struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    UID       `gorm:"not null;index"`
	CourseID  uint      `gorm:"not null"`
	Status    string    `gorm:"type:varchar(20);not null;default:'not_started'"` // not_started, in_progress, completed
	Score     int       `gorm:"not null;default:0"`                              // 习题得分
	Progress  float64   `gorm:"not null;default:0"`                              // 进度百分比
	CreatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
}

// Complete 完成课程学习
func (p *CourseLearningProgress) Complete(score int) {
	p.Status = "completed"
	p.Score = score
	p.Progress = 100
}

// SectionProgress 章节学习进度
type CourseSectionProgress struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"not null;index"`
	CourseID  uint      `gorm:"not null"`
	SectionID uint      `gorm:"not null"`
	Status    string    `gorm:"type:varchar(20);not null;default:'not_started'"` // not_started, in_progress, completed
	Progress  float64   `gorm:"not null;default:0"`                              // 进度百分比
	CreatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
}

// UnitProgress 单元学习进度
type CourseSectionUnitProgress struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"not null;index"`
	SectionID uint      `gorm:"not null"`
	UnitID    uint      `gorm:"not null"`
	Status    string    `gorm:"type:varchar(20);not null;default:'not_started'"` // not_started, in_progress, completed
	Progress  float64   `gorm:"not null;default:0"`                              // 完成百分比 0-100
	CreatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
}
