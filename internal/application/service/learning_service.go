package service

import (
	"context"

	"github.com/lazyjean/sla2/internal/domain/entity"
	"github.com/lazyjean/sla2/internal/domain/repository"
)

type LearningService struct {
	learningRepo  repository.LearningRepository
	memoryService MemoryService
}

func NewLearningService(learningRepo repository.LearningRepository, memoryService MemoryService) *LearningService {
	return &LearningService{
		learningRepo:  learningRepo,
		memoryService: memoryService,
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
func (s *LearningService) GetCourseProgress(ctx context.Context, courseID uint) (*entity.CourseLearningProgress, error) {
	userID, err := GetUserID(ctx)
	if err != nil {
		return nil, err
	}
	return s.learningRepo.GetCourseProgress(ctx, uint(userID), courseID)
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
func (s *LearningService) GetSectionProgress(ctx context.Context, sectionID uint) (*entity.CourseSectionProgress, error) {
	userID, err := GetUserID(ctx)
	if err != nil {
		return nil, err
	}
	return s.learningRepo.GetSectionProgress(ctx, uint(userID), sectionID)
}

// ListSectionProgress 获取课程的章节学习进度列表
func (s *LearningService) ListSectionProgress(ctx context.Context, courseID uint) ([]*entity.CourseSectionProgress, error) {
	userID, err := GetUserID(ctx)
	if err != nil {
		return nil, err
	}
	return s.learningRepo.ListSectionProgress(ctx, uint(userID), courseID)
}

// SaveUnitProgress 保存单元学习进度
func (s *LearningService) SaveUnitProgress(ctx context.Context, sectionID, unitID uint, status string, progress float64) (*entity.CourseSectionUnitProgress, error) {
	userID, err := GetUserID(ctx)
	if err != nil {
		return nil, err
	}
	unitProgress := &entity.CourseSectionUnitProgress{
		UserID:    uint(userID),
		SectionID: sectionID,
		UnitID:    unitID,
		Status:    status,
	}

	if err := s.learningRepo.UpsertUnitProgress(ctx, unitProgress); err != nil {
		return nil, err
	}

	return unitProgress, nil
}

// ListUnitProgress 获取章节的单元学习进度列表
func (s *LearningService) ListUnitProgress(ctx context.Context, sectionID uint) ([]*entity.CourseSectionUnitProgress, error) {
	userID, err := GetUserID(ctx)
	if err != nil {
		return nil, err
	}
	return s.learningRepo.ListUnitProgress(ctx, uint(userID), sectionID)
}

// UpdateUnitProgress 更新单元学习进度
func (s *LearningService) UpdateUnitProgress(ctx context.Context, unitID uint, sectionID uint, completed bool) error {
	userID, err := GetUserID(ctx)
	if err != nil {
		return err
	}

	// 设置状态
	status := "in_progress"
	if completed {
		status = "completed"
	}

	// 保存更新后的进度
	unitProgress := &entity.CourseSectionUnitProgress{
		UserID:        uint(userID),
		SectionID:     sectionID,
		UnitID:        unitID,
		Status:        status,
		CompleteCount: 0, // 初始完成次数为0，数据库会自动增加
	}

	if err := s.learningRepo.UpsertUnitProgress(ctx, unitProgress); err != nil {
		return err
	}

	return nil
}

// RecordLearningResult 记录学习结果
func (s *LearningService) RecordLearningResult(ctx context.Context, memoryUnitID uint32, result bool, responseTime uint32, userNotes []string) error {
	return s.memoryService.RecordLearningResult(ctx, memoryUnitID, result, responseTime, userNotes)
}

// GetCourseProgressWithStats 获取课程学习进度及统计信息
func (s *LearningService) GetCourseProgressWithStats(ctx context.Context, courseID uint) (float64, int, int, error) {
	// 获取课程进度
	progress, err := s.GetCourseProgress(ctx, courseID)
	if err != nil {
		return 0, 0, 0, err
	}

	// 获取章节总数
	sections, err := s.ListSectionProgress(ctx, courseID)
	if err != nil {
		return 0, 0, 0, err
	}

	// 计算完成章节数
	completedSections := 0
	for _, section := range sections {
		if section.Status == "completed" {
			completedSections++
		}
	}

	return progress.Progress, completedSections, len(sections), nil
}

// GetSectionProgressWithStats 获取章节学习进度及统计信息
func (s *LearningService) GetSectionProgressWithStats(ctx context.Context, sectionID uint) (float64, int, int, error) {
	// 获取章节进度
	progress, err := s.GetSectionProgress(ctx, sectionID)
	if err != nil {
		return 0, 0, 0, err
	}

	// 获取单元总数
	units, err := s.ListUnitProgress(ctx, sectionID)
	if err != nil {
		return 0, 0, 0, err
	}

	// 计算完成单元数
	completedUnits := 0
	for _, unit := range units {
		if unit.Status == "completed" {
			completedUnits++
		}
	}

	return progress.Progress, completedUnits, len(units), nil
}

// UpdateMemoryStatus 更新记忆单元状态
func (s *LearningService) UpdateMemoryStatus(ctx context.Context, memoryUnitID uint32, masteryLevel entity.MasteryLevel, studyDuration uint32) error {
	return s.memoryService.UpdateMemoryStatus(ctx, memoryUnitID, masteryLevel, studyDuration)
}
