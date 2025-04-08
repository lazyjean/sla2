package service

import (
	"context"
	"errors"
	"time"

	"github.com/lazyjean/sla2/internal/domain/entity"
	"github.com/lazyjean/sla2/internal/domain/repository"
	"github.com/lazyjean/sla2/pkg/logger"
	"go.uber.org/zap"
)

// MemoryService 记忆服务接口
type MemoryService interface {
	// ReviewWord 复习单词
	ReviewWord(ctx context.Context, wordID entity.WordID, result bool, responseTime uint32) error
	// ReviewHanChar 复习汉字
	ReviewHanChar(ctx context.Context, hanCharID uint32, result bool, responseTime uint32) error
	// GetNextReviewWords 获取下一批需要复习的单词
	GetNextReviewWords(ctx context.Context, limit int) ([]*entity.Word, error)
	// GetWordStats 获取单词的学习统计信息
	GetWordStats(ctx context.Context, wordID entity.WordID) (*WordStats, error)
	// GetLearningProgress 获取学习进度
	GetLearningProgress(ctx context.Context) (*LearningProgress, error)
	// UpdateMemoryStatus 更新记忆单元状态
	UpdateMemoryStatus(ctx context.Context, memoryUnitID uint32, masteryLevel entity.MasteryLevel, studyDuration uint32) error
	// GetMemoryUnit 获取记忆单元
	GetMemoryUnit(ctx context.Context, memoryUnitID uint32) (*entity.MemoryUnit, error)
	// ListMemoriesForReview 获取需要复习的记忆单元列表
	ListMemoriesForReview(ctx context.Context, page, pageSize uint32, types []entity.MemoryUnitType) ([]*entity.MemoryUnit, int, error)
}

// WordStats 单词学习统计信息
type WordStats struct {
	WordID           entity.WordID       `json:"word_id"`
	Text             string              `json:"text"`
	MasteryLevel     entity.MasteryLevel `json:"mastery_level"`
	ReviewCount      uint32              `json:"review_count"`
	LastReviewAt     time.Time           `json:"last_review_at"`
	NextReviewAt     time.Time           `json:"next_review_at"`
	RetentionRate    float32             `json:"retention_rate"`
	ConsecutiveRight uint32              `json:"consecutive_right"`
	ConsecutiveWrong uint32              `json:"consecutive_wrong"`
}

// LearningProgress 学习进度
type LearningProgress struct {
	TotalWords      int     `json:"total_words"`
	MasteredWords   int     `json:"mastered_words"`
	LearningWords   int     `json:"learning_words"`
	NewWords        int     `json:"new_words"`
	MasteryRate     float64 `json:"mastery_rate"`
	RetentionRate   float64 `json:"retention_rate"`
	DailyReviewGoal int     `json:"daily_review_goal"`
	DailyNewGoal    int     `json:"daily_new_goal"`
}

// ReviewInterval 复习间隔
type ReviewInterval struct {
	Days    uint32
	Hours   uint32
	Minutes uint32
}

// MemoryServiceImpl 记忆服务实现
type MemoryServiceImpl struct {
	wordRepo   repository.WordRepository
	memoryRepo repository.MemoryUnitRepository
}

// NewMemoryService 创建记忆服务实例
func NewMemoryService(wordRepo repository.WordRepository, memoryRepo repository.MemoryUnitRepository) MemoryService {
	return &MemoryServiceImpl{
		wordRepo:   wordRepo,
		memoryRepo: memoryRepo,
	}
}

// ReviewWord 复习单词
func (s *MemoryServiceImpl) ReviewWord(ctx context.Context, wordID entity.WordID, result bool, responseTime uint32) error {
	log := logger.GetLogger(ctx)

	// 获取或创建记忆单元
	memoryUnit, err := s.memoryRepo.GetByTypeAndContentID(ctx, entity.MemoryUnitTypeWord, uint32(wordID))
	if err != nil {
		log.Error("Failed to get memory unit by type and content ID", zap.Error(err), zap.Uint32("wordID", uint32(wordID)))
		return err
	}

	if memoryUnit == nil {
		// TODO: 需要确定 UserID 的来源，暂时使用 0
		userID := uint32(0)
		memoryUnit = entity.NewMemoryUnit(userID, entity.MemoryUnitTypeWord, uint32(wordID))
		if err := s.memoryRepo.Create(ctx, memoryUnit); err != nil {
			log.Error("Failed to create new memory unit for word", zap.Error(err), zap.Uint32("wordID", uint32(wordID)), zap.Uint32("userID", userID))
			return err
		}
		log.Info("Created new memory unit for word", zap.Uint32("unitID", memoryUnit.ID), zap.Uint32("wordID", uint32(wordID)))
	}

	// 更新记忆统计
	memoryUnit.UpdateReviewStats(result, responseTime)

	// 计算下次复习时间
	interval := s.calculateNextReviewInterval(memoryUnit)
	now := time.Now()
	nextReviewAt := now.Add(time.Duration(interval.Days)*24*time.Hour +
		time.Duration(interval.Hours)*time.Hour +
		time.Duration(interval.Minutes)*time.Minute)
	memoryUnit.NextReviewAt = nextReviewAt

	// 添加日志
	log.Info("[Service ReviewWord] Calculated review interval",
		zap.Uint32("UnitID", memoryUnit.ID),
		zap.Any("Interval", interval),
		zap.Time("LastReviewAt", memoryUnit.LastReviewAt),
		zap.Time("CalculatedNextReviewAt", nextReviewAt),
	)

	// 保存更新
	if err := s.memoryRepo.Update(ctx, memoryUnit); err != nil {
		log.Error("Failed to update memory unit after word review", zap.Error(err), zap.Uint32("unitID", memoryUnit.ID))
		return err
	}
	return nil
}

// ReviewHanChar 复习汉字
func (s *MemoryServiceImpl) ReviewHanChar(ctx context.Context, hanCharID uint32, result bool, responseTime uint32) error {
	log := logger.GetLogger(ctx)

	// 获取或创建记忆单元
	memoryUnit, err := s.memoryRepo.GetByTypeAndContentID(ctx, entity.MemoryUnitTypeHanChar, hanCharID)
	if err != nil {
		log.Error("Failed to get memory unit by type and content ID", zap.Error(err), zap.Uint32("hanCharID", hanCharID))
		return err
	}

	if memoryUnit == nil {
		// TODO: 需要确定 UserID 的来源，暂时使用 0
		userID := uint32(0)
		memoryUnit = entity.NewMemoryUnit(userID, entity.MemoryUnitTypeHanChar, hanCharID)
		if err := s.memoryRepo.Create(ctx, memoryUnit); err != nil {
			log.Error("Failed to create new memory unit for han char", zap.Error(err), zap.Uint32("hanCharID", hanCharID), zap.Uint32("userID", userID))
			return err
		}
		log.Info("Created new memory unit for han char", zap.Uint32("unitID", memoryUnit.ID), zap.Uint32("hanCharID", hanCharID))
	}

	// 更新记忆统计
	memoryUnit.UpdateReviewStats(result, responseTime)

	// 计算下次复习时间
	interval := s.calculateNextReviewInterval(memoryUnit)
	now := time.Now()
	nextReviewAt := now.Add(time.Duration(interval.Days)*24*time.Hour +
		time.Duration(interval.Hours)*time.Hour +
		time.Duration(interval.Minutes)*time.Minute)
	memoryUnit.NextReviewAt = nextReviewAt

	// 添加日志
	log.Info("[Service ReviewHanChar] Calculated review interval",
		zap.Uint32("UnitID", memoryUnit.ID),
		zap.Any("Interval", interval),
		zap.Time("LastReviewAt", memoryUnit.LastReviewAt),
		zap.Time("CalculatedNextReviewAt", nextReviewAt),
	)

	// 保存更新
	if err := s.memoryRepo.Update(ctx, memoryUnit); err != nil {
		log.Error("Failed to update memory unit after han char review", zap.Error(err), zap.Uint32("unitID", memoryUnit.ID))
		return err
	}
	return nil
}

// GetNextReviewWords 获取下一批需要复习的单词
func (s *MemoryServiceImpl) GetNextReviewWords(ctx context.Context, limit int) ([]*entity.Word, error) {
	// 获取需要复习的记忆单元
	units, err := s.memoryRepo.ListNeedReview(ctx, entity.MemoryUnitTypeWord, time.Now(), limit)
	if err != nil {
		return nil, err
	}

	// 获取对应的单词
	var words []*entity.Word
	for _, unit := range units {
		word, err := s.wordRepo.GetByID(ctx, entity.WordID(unit.ContentID))
		if err != nil {
			continue
		}
		words = append(words, word)
	}

	return words, nil
}

// GetWordStats 获取单词的学习统计信息
func (s *MemoryServiceImpl) GetWordStats(ctx context.Context, wordID entity.WordID) (*WordStats, error) {
	// 获取单词
	word, err := s.wordRepo.GetByID(ctx, wordID)
	if err != nil {
		return nil, err
	}

	// 获取记忆单元
	unit, err := s.memoryRepo.GetByTypeAndContentID(ctx, entity.MemoryUnitTypeWord, uint32(wordID))
	if err != nil {
		return nil, err
	}
	if unit == nil {
		unit = entity.NewMemoryUnit(0, entity.MemoryUnitTypeWord, uint32(wordID))
	}

	return &WordStats{
		WordID:           wordID,
		Text:             word.Text,
		MasteryLevel:     unit.MasteryLevel,
		ReviewCount:      unit.ReviewCount,
		LastReviewAt:     unit.LastReviewAt,
		NextReviewAt:     unit.NextReviewAt,
		RetentionRate:    unit.RetentionRate,
		ConsecutiveRight: unit.ConsecutiveCorrect,
		ConsecutiveWrong: unit.ConsecutiveWrong,
	}, nil
}

// GetLearningProgress 获取学习进度
func (s *MemoryServiceImpl) GetLearningProgress(ctx context.Context) (*LearningProgress, error) {
	// TODO: 实现学习进度统计
	return &LearningProgress{}, nil
}

// UpdateMemoryStatus 更新记忆单元状态
func (s *MemoryServiceImpl) UpdateMemoryStatus(ctx context.Context, memoryUnitID uint32, masteryLevel entity.MasteryLevel, studyDuration uint32) error {
	// 获取记忆单元
	memoryUnit, err := s.memoryRepo.GetByID(ctx, memoryUnitID)
	if err != nil {
		return err
	}
	if memoryUnit == nil {
		return errors.New("memory unit not found")
	}

	// 更新状态
	memoryUnit.MasteryLevel = masteryLevel
	memoryUnit.StudyDuration += studyDuration
	memoryUnit.Update()

	// 保存更新
	return s.memoryRepo.Update(ctx, memoryUnit)
}

// GetMemoryUnit 获取记忆单元
func (s *MemoryServiceImpl) GetMemoryUnit(ctx context.Context, memoryUnitID uint32) (*entity.MemoryUnit, error) {
	return s.memoryRepo.GetByID(ctx, memoryUnitID)
}

// ListMemoriesForReview 获取需要复习的记忆单元列表
func (s *MemoryServiceImpl) ListMemoriesForReview(ctx context.Context, page, pageSize uint32, types []entity.MemoryUnitType) ([]*entity.MemoryUnit, int, error) {
	// 计算 offset
	offset := (page - 1) * pageSize
	limit := int(pageSize)
	now := time.Now()

	// 调用 repository 获取数据
	units, err := s.memoryRepo.ListNeedReviewByTypes(ctx, types, now, offset, limit)
	if err != nil {
		// 考虑记录日志
		// log.Printf("Error listing memories for review: %v", err)
		return nil, 0, err // 直接返回错误
	}

	// 调用 repository 获取总数
	total, err := s.memoryRepo.CountNeedReviewByTypes(ctx, types, now)
	if err != nil {
		// 考虑记录日志
		// log.Printf("Error counting memories for review: %v", err)
		return nil, 0, err // 直接返回错误
	}

	return units, int(total), nil
}

// calculateNextReviewInterval 计算下次复习间隔
func (s *MemoryServiceImpl) calculateNextReviewInterval(unit *entity.MemoryUnit) *ReviewInterval {
	// 基础间隔（小时）
	var baseInterval float64
	switch unit.MasteryLevel {
	case entity.MasteryLevelUnlearned:
		baseInterval = 1 // 1小时
	case entity.MasteryLevelBeginner:
		baseInterval = 4 // 4小时
	case entity.MasteryLevelFamiliar:
		baseInterval = 24 // 1天
	case entity.MasteryLevelMastered:
		baseInterval = 72 // 3天
	case entity.MasteryLevelExpert:
		baseInterval = 168 // 7天
	default:
		baseInterval = 1
	}

	// 根据连续正确次数调整间隔
	if unit.ConsecutiveCorrect > 0 {
		baseInterval *= 1.2 * float64(unit.ConsecutiveCorrect)
	}

	// 根据连续错误次数减少间隔
	if unit.ConsecutiveWrong > 0 {
		baseInterval /= 2.0 * float64(unit.ConsecutiveWrong)
	}

	// 转换为天、小时、分钟
	totalMinutes := int(baseInterval * 60)
	days := uint32(totalMinutes / (24 * 60))
	totalMinutes = totalMinutes % (24 * 60)
	hours := uint32(totalMinutes / 60)
	minutes := uint32(totalMinutes % 60)

	return &ReviewInterval{
		Days:    days,
		Hours:   hours,
		Minutes: minutes,
	}
}

var _ MemoryService = (*MemoryServiceImpl)(nil)
