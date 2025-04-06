package service

import (
	"context"
	"time"

	"github.com/lazyjean/sla2/internal/domain/entity"
	"github.com/lazyjean/sla2/internal/domain/repository"
)

// MemoryService 记忆服务接口
type MemoryService interface {
	// ReviewWord 复习单词
	ReviewWord(ctx context.Context, wordID entity.WordID, result bool, responseTime uint32) error
	// GetNextReviewWords 获取下一批需要复习的单词
	GetNextReviewWords(ctx context.Context, limit int) ([]*entity.Word, error)
	// GetWordStats 获取单词的学习统计信息
	GetWordStats(ctx context.Context, wordID entity.WordID) (*WordStats, error)
	// GetLearningProgress 获取学习进度
	GetLearningProgress(ctx context.Context) (*LearningProgress, error)
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

// memoryService 记忆服务实现
type memoryService struct {
	wordRepo   repository.WordRepository
	memoryRepo repository.MemoryUnitRepository
}

// NewMemoryService 创建记忆服务实例
func NewMemoryService(wordRepo repository.WordRepository, memoryRepo repository.MemoryUnitRepository) MemoryService {
	return &memoryService{
		wordRepo:   wordRepo,
		memoryRepo: memoryRepo,
	}
}

// ReviewWord 复习单词
func (s *memoryService) ReviewWord(ctx context.Context, wordID entity.WordID, result bool, responseTime uint32) error {
	// 获取或创建记忆单元
	memoryUnit, err := s.memoryRepo.GetByTypeAndContentID(ctx, entity.MemoryUnitTypeWord, uint32(wordID))
	if err != nil {
		return err
	}

	if memoryUnit == nil {
		memoryUnit = entity.NewMemoryUnit(0, entity.MemoryUnitTypeWord, uint32(wordID))
		if err := s.memoryRepo.Create(ctx, memoryUnit); err != nil {
			return err
		}
	}

	// 更新记忆统计
	memoryUnit.UpdateReviewStats(result, responseTime)

	// 保存更新
	return s.memoryRepo.Update(ctx, memoryUnit)
}

// GetNextReviewWords 获取下一批需要复习的单词
func (s *memoryService) GetNextReviewWords(ctx context.Context, limit int) ([]*entity.Word, error) {
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
func (s *memoryService) GetWordStats(ctx context.Context, wordID entity.WordID) (*WordStats, error) {
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
func (s *memoryService) GetLearningProgress(ctx context.Context) (*LearningProgress, error) {
	// TODO: 实现学习进度统计
	return &LearningProgress{}, nil
}
