package entity

import "time"

// CourseLearningProgress 课程学习进度
type CourseLearningProgress struct {
	ID        uint      `gorm:"primaryKey" comment:"主键ID"`
	UserID    UID       `gorm:"not null;index" comment:"用户ID"`
	CourseID  uint      `gorm:"not null" comment:"课程ID"`
	Status    string    `gorm:"type:varchar(20);not null;default:'not_started'" comment:"学习状态：未开始、进行中、已完成"`
	Score     int       `gorm:"not null;default:0" comment:"习题得分"`
	Progress  float64   `gorm:"not null;default:0" comment:"进度百分比"`
	CreatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" comment:"创建时间"`
	UpdatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" comment:"更新时间"`
}

// Complete 完成课程学习
func (p *CourseLearningProgress) Complete(score int) {
	p.Status = "completed"
	p.Score = score
	p.Progress = 100
}

// SectionProgress 章节学习进度
type CourseSectionProgress struct {
	ID        uint      `gorm:"primaryKey" comment:"主键ID"`
	UserID    uint      `gorm:"not null;index" comment:"用户ID"`
	CourseID  uint      `gorm:"not null" comment:"课程ID"`
	SectionID uint      `gorm:"not null" comment:"章节ID"`
	Status    string    `gorm:"type:varchar(20);not null;default:'not_started'" comment:"学习状态：未开始、进行中、已完成"`
	Progress  float64   `gorm:"not null;default:0" comment:"进度百分比"`
	CreatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" comment:"创建时间"`
	UpdatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" comment:"更新时间"`
}

// UnitProgress 单元学习进度
type CourseSectionUnitProgress struct {
	ID            uint      `gorm:"primaryKey" comment:"主键ID"`
	UserID        uint      `gorm:"not null;index:idx_user_section_unit,unique" comment:"用户ID"`
	SectionID     uint      `gorm:"not null;index:idx_user_section_unit,unique" comment:"章节ID"`
	UnitID        uint      `gorm:"not null;index:idx_user_section_unit,unique" comment:"单元ID"`
	Status        string    `gorm:"type:varchar(20);not null;default:'not_started'" comment:"学习状态：未开始、进行中、已完成"`
	CompleteCount uint      `gorm:"not null;default:0" comment:"完成次数"`
	CreatedAt     time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" comment:"创建时间"`
	UpdatedAt     time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" comment:"更新时间"`
}
