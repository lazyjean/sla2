package service

import (
	"context"
	"fmt"
	"time"

	"github.com/lazyjean/sla2/internal/application/dto"
	"github.com/lazyjean/sla2/internal/domain/entity"
	"github.com/lazyjean/sla2/internal/domain/repository"
	"github.com/lazyjean/sla2/internal/domain/service"
	"github.com/lazyjean/sla2/pkg/logger"
	"go.uber.org/zap"
)

type LearningService struct {
	learningRepo      repository.LearningRepository
	memoryRepo        repository.MemoryUnitRepository
	wordRepo          repository.WordRepository
	hanCharRepo       repository.HanCharRepository
	memoryService     service.MemoryService
	vocabularyService VocabularyService
}

func NewLearningService(
	learningRepo repository.LearningRepository,
	memoryRepo repository.MemoryUnitRepository,
	wordRepo repository.WordRepository,
	hanCharRepo repository.HanCharRepository,
	memoryService service.MemoryService,
	vocabularyService VocabularyService,
) *LearningService {
	return &LearningService{
		learningRepo:      learningRepo,
		memoryRepo:        memoryRepo,
		wordRepo:          wordRepo,
		hanCharRepo:       hanCharRepo,
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

// GetSectionProgressWithCompletedUnits 获取章节进度及已完成的单元ID
func (s *LearningService) GetSectionProgressWithCompletedUnits(ctx context.Context, sectionID uint) (float64, int, int, []uint32, error) {
	// 获取章节进度及统计信息
	progress, completedUnits, totalUnits, err := s.GetSectionProgressWithStats(ctx, sectionID)
	if err != nil {
		return 0, 0, 0, nil, err
	}

	// 获取已完成的单元ID列表
	unitProgresses, err := s.ListUnitProgress(ctx, sectionID)
	if err != nil {
		return 0, 0, 0, nil, err
	}

	// 收集已完成的单元ID
	completedUnitIDs := make([]uint32, 0)
	for _, unit := range unitProgresses {
		if unit.Status == "completed" {
			completedUnitIDs = append(completedUnitIDs, uint32(unit.UnitID))
		}
	}

	return progress, completedUnits, totalUnits, completedUnitIDs, nil
}

// GetReviewContent 获取需要复习的内容（汉字和单词）
func (s *LearningService) GetReviewContent(ctx context.Context, count int) ([]*entity.HanChar, []*entity.Word, bool, error) {
	units, _, err := s.ListMemoriesForReview(ctx, entity.MemoryUnitTypeUnspecified, count)
	if err != nil {
		return nil, nil, false, err
	}

	var hanChars []*entity.HanChar
	var words []*entity.Word
	for _, unit := range units {
		switch unit.Type {
		case entity.MemoryUnitTypeHanChar:
			hanChar, err := s.vocabularyService.GetHanChar(ctx, entity.HanCharID(unit.ContentID))
			if err != nil {
				continue
			}
			hanChars = append(hanChars, hanChar)
		case entity.MemoryUnitTypeWord:
			word, err := s.vocabularyService.GetWord(ctx, uint(unit.ContentID))
			if err != nil {
				continue
			}
			words = append(words, word)
		}
	}

	hasMore := len(units) == count
	return hanChars, words, hasMore, nil
}

// ListMemoriesForReview 获取需要复习的记忆单元列表
func (s *LearningService) ListMemoriesForReview(ctx context.Context, unitType entity.MemoryUnitType, limit int) ([]*entity.MemoryUnit, int, error) {
	// 从 context 获取 userID
	userID, err := GetUserID(ctx)
	if err != nil {
		return nil, 0, err
	}

	now := time.Now()

	// 调用 repository 获取数据
	units, err := s.memoryRepo.ListNeedReview(ctx, userID, unitType, now, limit)
	if err != nil {
		return nil, 0, err
	}

	// 调用 repository 获取总数
	total, err := s.memoryRepo.CountNeedReviewByType(ctx, userID, unitType, time.Now())
	if err != nil {
		return nil, 0, err
	}

	return units, int(total), nil
}

// GetMemoryStats 获取记忆统计信息
func (s *LearningService) GetMemoryStats(ctx context.Context, unitType *entity.MemoryUnitType) (*entity.MemoryUnitStats, error) {
	log := logger.GetLogger(ctx)

	// 从 context 获取 userID
	userID, err := GetUserID(ctx)
	if err != nil {
		log.Error("获取用户ID失败", zap.Error(err))
		return nil, err
	}

	// 如果未指定类型，使用默认类型
	if unitType == nil {
		defaultType := entity.MemoryUnitTypeHanChar
		unitType = &defaultType
	}

	// 调用 repository 获取统计信息
	stats, err := s.memoryRepo.GetStats(ctx, entity.UID(userID), *unitType)
	if err != nil {
		log.Error("获取统计信息失败", zap.Error(err))
		return nil, err
	}

	return stats, nil
}

// ReviewMemoryUnits 复习记忆单元
func (s *LearningService) ReviewMemoryUnits(ctx context.Context, items []*dto.ReviewItem) ([]*entity.MemoryUnit, error) {
	log := logger.GetLogger(ctx)

	// 从 context 获取 userID
	userID, err := GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	// 收集所有需要更新的记忆单元
	updatedUnits := make([]*entity.MemoryUnit, 0, len(items))
	now := time.Now()
	var updateErrors []error

	for i, item := range items {
		// 获取记忆单元
		memoryUnit, err := s.memoryRepo.GetByID(ctx, entity.MemoryUnitID(item.MemoryUnitID))
		if err != nil {
			log.Error("获取记忆单元失败",
				zap.Error(err),
				zap.Uint32("memory_unit_id", item.MemoryUnitID),
				zap.Int("item_index", i))
			updateErrors = append(updateErrors, fmt.Errorf("failed to get memory unit at index %d: %w", i, err))
			continue
		}
		if memoryUnit == nil {
			log.Error("记忆单元不存在",
				zap.Uint32("memory_unit_id", item.MemoryUnitID),
				zap.Int("item_index", i))
			updateErrors = append(updateErrors, fmt.Errorf("memory unit not found at index %d", i))
			continue
		}
		if memoryUnit.UserID != userID {
			updateErrors = append(updateErrors, fmt.Errorf("invalid memory unit %d not beyond to user %d", memoryUnit.ID, memoryUnit.UserID))
			continue
		}

		// 更新记忆单元状态
		memoryUnit.LastReviewAt = now
		memoryUnit.ReviewCount++
		if item.Result {
			memoryUnit.ConsecutiveCorrect++
			memoryUnit.ConsecutiveWrong = 0
		} else {
			memoryUnit.ConsecutiveCorrect = 0
			memoryUnit.ConsecutiveWrong++
		}

		// 更新学习时长
		memoryUnit.StudyDuration += uint32(item.ReviewDuration.Seconds())

		// 计算下次复习时间
		memoryUnit.NextReviewAt = now.Add(s.memoryService.CalculateNextReviewInterval(memoryUnit))

		updatedUnits = append(updatedUnits, memoryUnit)
	}

	// 如果有需要更新的单元，执行批量更新
	if len(updatedUnits) > 0 {
		if err := s.memoryRepo.UpdateBatch(ctx, updatedUnits); err != nil {
			log.Error("批量更新记忆单元失败",
				zap.Error(err),
				zap.Int("total_units", len(updatedUnits)))
			updateErrors = append(updateErrors, fmt.Errorf("failed to update memory units: %w", err))
		}
	}

	// 如果有错误，返回部分成功的结果和错误信息
	if len(updateErrors) > 0 {
		log.Warn("部分记忆单元更新失败",
			zap.Int("success_count", len(updatedUnits)),
			zap.Int("error_count", len(updateErrors)))
		return updatedUnits, fmt.Errorf("部分更新失败: %v", updateErrors)
	}
	return updatedUnits, nil
}

// InitializeMemoryUnits 初始化记忆单元
func (s *LearningService) InitializeMemoryUnits(ctx context.Context, items []*entity.MemoryUnitInitItem) ([]uint32, error) {
	log := logger.GetLogger(ctx)

	// 从 context 获取 userID
	userID, err := GetUserID(ctx)
	if err != nil {
		log.Error("获取用户ID失败", zap.Error(err))
		return nil, err
	}

	// 创建记忆单元列表
	memoryUnits := make([]*entity.MemoryUnit, len(items))
	for i, item := range items {
		memoryUnit := entity.NewMemoryUnit(userID, item.Type, item.ContentID)
		memoryUnit.MasteryLevel = item.MasteryLevel
		memoryUnits[i] = memoryUnit
	}

	// 批量保存记忆单元
	if err := s.memoryRepo.CreateBatch(ctx, memoryUnits); err != nil {
		log.Error("创建记忆单元失败", zap.Error(err))
		return nil, err
	}

	// 收集创建的单元ID
	ids := make([]uint32, len(memoryUnits))
	for i, unit := range memoryUnits {
		ids[i] = uint32(unit.ID)
	}

	return ids, nil
}
