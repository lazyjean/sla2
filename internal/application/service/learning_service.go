package service

import (
	"context"

	"github.com/lazyjean/sla2/internal/domain/entity"
	"github.com/lazyjean/sla2/internal/domain/repository"
	"github.com/lazyjean/sla2/internal/domain/valueobject"
)

type LearningService struct {
	learningRepo      repository.LearningRepository
	memoryService     MemoryService
	vocabularyService VocabularyService
}

func NewLearningService(learningRepo repository.LearningRepository, memoryService MemoryService, vocabularyService VocabularyService) *LearningService {
	return &LearningService{
		learningRepo:      learningRepo,
		memoryService:     memoryService,
		vocabularyService: vocabularyService,
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

// GetHanCharTest 获取汉字测试题目
func (s *LearningService) GetHanCharTest(ctx context.Context, level valueobject.WordDifficultyLevel) ([]*entity.HanChar, error) {
	// 1. 从 context 获取 userID
	_, err := GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	// 2. 获取用户已学习的汉字ID
	memoryUnits, _, err := s.memoryService.ListMemoriesForReview(ctx, 1, 100, []entity.MemoryUnitType{entity.MemoryUnitTypeHanChar}, nil)
	if err != nil {
		return nil, err
	}

	// 3. 收集已学习的汉字ID
	excludeIDs := make([]uint, 0, len(memoryUnits))
	for _, unit := range memoryUnits {
		excludeIDs = append(excludeIDs, uint(unit.ContentID))
	}

	// 4. 从 vocabulary service 获取汉字列表,排除已学习的汉字
	hanChars, _, err := s.vocabularyService.ListHanChars(ctx, 1, 50, level, nil, nil, excludeIDs)
	if err != nil {
		return nil, err
	}

	return hanChars, nil
}

// SubmitHanCharTest 提交汉字测试答案
func (s *LearningService) SubmitHanCharTest(ctx context.Context, hanCharID uint32, isCorrect bool) error {
	// 1. 从 context 获取 userID
	userID, err := GetUserID(ctx)
	if err != nil {
		return err
	}

	// 2. 更新记忆单元状态
	masteryLevel := entity.MasteryLevelUnspecified
	if isCorrect {
		masteryLevel = entity.MasteryLevelFamiliar
	} else {
		masteryLevel = entity.MasteryLevelBeginner
	}

	// 3. 获取或创建记忆单元
	memoryUnit, err := s.memoryService.GetMemoryUnit(ctx, hanCharID)
	if err != nil {
		// 如果记忆单元不存在,创建新的
		memoryUnit = entity.NewMemoryUnit(entity.UID(userID), entity.MemoryUnitTypeHanChar, hanCharID)
		if err := s.memoryService.UpdateMemoryStatus(ctx, uint32(memoryUnit.ID), masteryLevel, 0); err != nil {
			return err
		}
	} else {
		// 更新现有记忆单元
		if err := s.memoryService.UpdateMemoryStatus(ctx, uint32(memoryUnit.ID), masteryLevel, 0); err != nil {
			return err
		}
	}

	return nil
}
