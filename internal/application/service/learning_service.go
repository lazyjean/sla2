package service

import (
	"context"

	pb "github.com/lazyjean/sla2/api/proto/v1"
	"github.com/lazyjean/sla2/internal/domain/entity"
	"github.com/lazyjean/sla2/internal/domain/repository"
)

type LearningService struct {
	pb.UnimplementedLearningServiceServer
	learningRepo repository.LearningRepository
}

func NewLearningService(learningRepo repository.LearningRepository) *LearningService {
	return &LearningService{
		learningRepo: learningRepo,
	}
}

// SaveCourseProgress 保存课程学习进度
func (s *LearningService) SaveCourseProgress(ctx context.Context, userID entity.UID, courseID uint, status string, score int) (*entity.CourseLearningProgress, error) {
	progress := &entity.CourseLearningProgress{
		UserID:   userID,
		CourseID: courseID,
		Status:   status,
		Score:    score,
	}

	if status == "completed" {
		progress.Score = score
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
func (s *LearningService) SaveSectionProgress(ctx context.Context, userID, courseID, sectionID uint, status string, progress float64) (*entity.CourseSectionProgress, error) {
	sectionProgress := &entity.CourseSectionProgress{
		UserID:    userID,
		CourseID:  courseID,
		SectionID: sectionID,
		Status:    status,
		Progress:  progress,
	}

	if status == "completed" {
		sectionProgress.Progress = 100
	}

	if err := s.learningRepo.SaveSectionProgress(ctx, sectionProgress); err != nil {
		return nil, err
	}

	return sectionProgress, nil
}

// GetSectionProgress 获取章节学习进度
func (s *LearningService) GetSectionProgress(ctx context.Context, userID, sectionID uint) (*entity.CourseSectionProgress, error) {
	return s.learningRepo.GetSectionProgress(ctx, userID, sectionID)
}

// ListSectionProgress 获取课程的章节学习进度列表
func (s *LearningService) ListSectionProgress(ctx context.Context, userID, courseID uint) ([]*entity.CourseSectionProgress, error) {
	return s.learningRepo.ListSectionProgress(ctx, userID, courseID)
}

// SaveUnitProgress 保存单元学习进度
func (s *LearningService) SaveUnitProgress(ctx context.Context, userID, sectionID, unitID uint, status string, progress float64, lastWordID *uint) (*entity.CourseSectionUnitProgress, error) {
	unitProgress := &entity.CourseSectionUnitProgress{
		UserID:    userID,
		SectionID: sectionID,
		UnitID:    unitID,
		Status:    status,
		Progress:  progress,
	}

	if status == "completed" {
		unitProgress.Progress = 100
	}

	if err := s.learningRepo.SaveUnitProgress(ctx, unitProgress); err != nil {
		return nil, err
	}

	return unitProgress, nil
}

// GetUnitProgress 获取单元学习进度
func (s *LearningService) GetUnitProgress(ctx context.Context, userID, unitID uint) (*entity.CourseSectionUnitProgress, error) {
	return s.learningRepo.GetUnitProgress(ctx, userID, unitID)
}

// ListUnitProgress 获取章节的单元学习进度列表
func (s *LearningService) ListUnitProgress(ctx context.Context, userID, sectionID uint) ([]*entity.CourseSectionUnitProgress, error) {
	return s.learningRepo.ListUnitProgress(ctx, userID, sectionID)
}
