package entity

import (
	"time"

	pg "github.com/lib/pq"
)

type CourseID uint32

// Course 课程实体
type Course struct {
	ID          CourseID       `gorm:"primaryKey"`
	Title       string         `gorm:"type:varchar(255);not null"`
	Description string         `gorm:"type:text"`
	CoverURL    string         `gorm:"type:varchar(255)"`
	Level       string         `gorm:"type:varchar(50);not null"`
	Tags        pg.StringArray `gorm:"type:varchar(50)[]"`
	Status      string         `gorm:"type:varchar(50);not null;default:'draft'"`
	CreatedAt   time.Time      `gorm:"not null"`
	UpdatedAt   time.Time      `gorm:"not null"`
}

// TableName 指定表名
func (Course) TableName() string {
	return "courses"
}
