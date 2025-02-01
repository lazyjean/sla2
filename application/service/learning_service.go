package service

import (
	"context"
	"time"

	"github.com/lazyjean/sla2/domain/entity"
	"github.com/lazyjean/sla2/domain/repository"
)

type LearningService struct {
	learningRepo repository.LearningRepository
}

func NewLearningService(learningRepo repository.LearningRepository) *LearningService {
	return &LearningService{
		learningRepo: learningRepo,
	}
}

// SaveCourseProgress 保存课程学习进度
func (s *LearningService) SaveCourseProgress(ctx context.Context, userID, courseID uint, status string, score int) (*entity.CourseLearningProgress, error) {
	progress := &entity.CourseLearningProgress{
		UserID:    userID,
		CourseID:  courseID,
		Status:    status,
		Score:     score,
		StartedAt: time.Now(),
	}

	if status == "completed" {
		now := time.Now()
		progress.CompletedAt = &now
	}

	if err := s.learningRepo.SaveCourseProgress(ctx, progress); err != nil {
		return nil, err
	}

	return progress, nil
}

// GetCourseProgress 获取课程学习进度
func (s *LearningService) GetCourseProgress(ctx context.Context, userID, courseID uint) (*entity.CourseLearningProgress, error) {
	return s.learningRepo.GetCourseProgress(ctx, userID, courseID)
}

// ListCourseProgress 获取用户的课程学习进度列表
func (s *LearningService) ListCourseProgress(ctx context.Context, userID uint, page, pageSize int) ([]*entity.CourseLearningProgress, int64, error) {
	offset := (page - 1) * pageSize
	return s.learningRepo.ListCourseProgress(ctx, userID, offset, pageSize)
}

// SaveSectionProgress 保存章节学习进度
func (s *LearningService) SaveSectionProgress(ctx context.Context, userID, courseID, sectionID uint, status string, progress float64) (*entity.SectionProgress, error) {
	sectionProgress := &entity.SectionProgress{
		UserID:    userID,
		CourseID:  courseID,
		SectionID: sectionID,
		Status:    status,
		Progress:  progress,
		StartedAt: time.Now(),
	}

	if status == "completed" {
		now := time.Now()
		sectionProgress.CompletedAt = &now
	}

	if err := s.learningRepo.SaveSectionProgress(ctx, sectionProgress); err != nil {
		return nil, err
	}

	return sectionProgress, nil
}

// GetSectionProgress 获取章节学习进度
func (s *LearningService) GetSectionProgress(ctx context.Context, userID, sectionID uint) (*entity.SectionProgress, error) {
	return s.learningRepo.GetSectionProgress(ctx, userID, sectionID)
}

// ListSectionProgress 获取课程的章节学习进度列表
func (s *LearningService) ListSectionProgress(ctx context.Context, userID, courseID uint) ([]*entity.SectionProgress, error) {
	return s.learningRepo.ListSectionProgress(ctx, userID, courseID)
}

// SaveUnitProgress 保存单元学习进度
func (s *LearningService) SaveUnitProgress(ctx context.Context, userID, sectionID, unitID uint, status string, progress float64, lastWordID *uint) (*entity.UnitProgress, error) {
	unitProgress := &entity.UnitProgress{
		UserID:     userID,
		SectionID:  sectionID,
		UnitID:     unitID,
		Status:     status,
		Progress:   progress,
		StartedAt:  time.Now(),
		LastWordID: lastWordID,
	}

	if status == "completed" {
		now := time.Now()
		unitProgress.CompletedAt = &now
	}

	if err := s.learningRepo.SaveUnitProgress(ctx, unitProgress); err != nil {
		return nil, err
	}

	return unitProgress, nil
}

// GetUnitProgress 获取单元学习进度
func (s *LearningService) GetUnitProgress(ctx context.Context, userID, unitID uint) (*entity.UnitProgress, error) {
	return s.learningRepo.GetUnitProgress(ctx, userID, unitID)
}

// ListUnitProgress 获取章节的单元学习进度列表
func (s *LearningService) ListUnitProgress(ctx context.Context, userID, sectionID uint) ([]*entity.UnitProgress, error) {
	return s.learningRepo.ListUnitProgress(ctx, userID, sectionID)
}
