package dto

import (
	"time"
)

// CourseProgressDTO 课程进度DTO
type CourseProgressDTO struct {
	ID          uint       `json:"id"`
	CourseID    uint       `json:"course_id"`
	Status      string     `json:"status"`
	Score       int        `json:"score"`
	StartedAt   time.Time  `json:"started_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// SectionProgressDTO 章节进度DTO
type SectionProgressDTO struct {
	ID          uint       `json:"id"`
	CourseID    uint       `json:"course_id"`
	SectionID   uint       `json:"section_id"`
	Status      string     `json:"status"`
	Progress    float64    `json:"progress"`
	StartedAt   time.Time  `json:"started_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// UnitProgressDTO 单元进度DTO
type UnitProgressDTO struct {
	ID          uint       `json:"id"`
	UnitID      uint       `json:"unit_id"`
	Status      string     `json:"status"`
	StartedAt   time.Time  `json:"started_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}
