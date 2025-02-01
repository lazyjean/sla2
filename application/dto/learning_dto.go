package dto

import (
	"time"

	"github.com/lazyjean/sla2/domain/entity"
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
	SectionID   uint       `json:"section_id"`
	UnitID      uint       `json:"unit_id"`
	Status      string     `json:"status"`
	Progress    float64    `json:"progress"`
	StartedAt   time.Time  `json:"started_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	LastWordID  *uint      `json:"last_word_id,omitempty"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// 转换函数
func CourseProgressToDTO(progress *entity.CourseLearningProgress) *CourseProgressDTO {
	return &CourseProgressDTO{
		ID:          progress.ID,
		CourseID:    progress.CourseID,
		Status:      progress.Status,
		Score:       progress.Score,
		StartedAt:   progress.StartedAt,
		CompletedAt: progress.CompletedAt,
		UpdatedAt:   progress.UpdatedAt,
	}
}

func SectionProgressToDTO(progress *entity.SectionProgress) *SectionProgressDTO {
	return &SectionProgressDTO{
		ID:          progress.ID,
		CourseID:    progress.CourseID,
		SectionID:   progress.SectionID,
		Status:      progress.Status,
		Progress:    progress.Progress,
		StartedAt:   progress.StartedAt,
		CompletedAt: progress.CompletedAt,
		UpdatedAt:   progress.UpdatedAt,
	}
}

func UnitProgressToDTO(progress *entity.UnitProgress) *UnitProgressDTO {
	return &UnitProgressDTO{
		ID:          progress.ID,
		SectionID:   progress.SectionID,
		UnitID:      progress.UnitID,
		Status:      progress.Status,
		Progress:    progress.Progress,
		StartedAt:   progress.StartedAt,
		CompletedAt: progress.CompletedAt,
		LastWordID:  progress.LastWordID,
		UpdatedAt:   progress.UpdatedAt,
	}
}
