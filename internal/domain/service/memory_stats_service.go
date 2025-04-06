package service

import (
	"context"
	"sort"
	"time"

	"github.com/lazyjean/sla2/internal/domain/entity"
	"github.com/lazyjean/sla2/internal/domain/repository"
)

// MemoryStatsService 记忆统计服务
type MemoryStatsService struct {
	memoryUnitRepo   repository.MemoryUnitRepository
	memoryReviewRepo repository.MemoryReviewRepository
}

// NewMemoryStatsService 创建记忆统计服务
func NewMemoryStatsService(
	memoryUnitRepo repository.MemoryUnitRepository,
	memoryReviewRepo repository.MemoryReviewRepository,
) *MemoryStatsService {
	return &MemoryStatsService{
		memoryUnitRepo:   memoryUnitRepo,
		memoryReviewRepo: memoryReviewRepo,
	}
}

// CalculateUserStats 计算用户记忆统计
func (s *MemoryStatsService) CalculateUserStats(ctx context.Context, userID uint32) (*entity.MemoryStats, error) {
	stats := entity.NewMemoryStats()

	// 1. 获取用户的所有记忆单元
	units, err := s.memoryUnitRepo.ListByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 2. 计算基础统计数据
	for _, unit := range units {
		stats.TotalLearned++
		if unit.MasteryLevel >= entity.MasteryLevelMastered {
			stats.MasteredCount++
		}
		if time.Now().After(unit.NextReviewAt) {
			stats.NeedReviewCount++
		}
		stats.TotalStudyTime += unit.StudyDuration

		// 更新掌握程度统计
		stats.UpdateLevelStats(uint32(unit.MasteryLevel), stats.LevelStats[uint32(unit.MasteryLevel)]+1)
	}

	// 3. 获取最近30天的复习记录
	startTime := time.Now().AddDate(0, 0, -30)
	reviews, err := s.memoryReviewRepo.ListByUserIDAndTimeRange(ctx, userID, startTime, time.Now())
	if err != nil {
		return nil, err
	}

	// 4. 计算每日统计
	dailyMap := make(map[string]*entity.DailyStat)
	for _, review := range reviews {
		date := review.ReviewTime.Format("2006-01-02")
		if _, exists := dailyMap[date]; !exists {
			dailyMap[date] = &entity.DailyStat{
				Date: review.ReviewTime.Truncate(24 * time.Hour),
			}
		}
		dailyStat := dailyMap[date]

		dailyStat.ReviewCount++
		if review.IsCorrect() {
			dailyStat.CorrectCount++
			dailyStat.RetentionRate = float32(dailyStat.CorrectCount) / float32(dailyStat.ReviewCount)
		}
	}

	// 5. 计算记忆保持率
	for _, unit := range units {
		stats.UpdateRetentionRate(uint32(unit.ID), unit.RetentionRate)
	}

	// 6. 添加每日统计
	for _, dailyStat := range dailyMap {
		stats.AddDailyStat(*dailyStat)
	}

	return stats, nil
}

// CalculateDailyStats 计算每日统计
func (s *MemoryStatsService) CalculateDailyStats(ctx context.Context, userID uint32, startDate, endDate time.Time) ([]*entity.DailyStat, error) {
	var stats []*entity.DailyStat

	// 获取日期范围内的所有复习记录
	reviews, err := s.memoryReviewRepo.ListByUserIDAndTimeRange(ctx, userID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	// 按日期分组统计
	dailyMap := make(map[time.Time]*entity.DailyStat)
	for _, review := range reviews {
		date := time.Date(review.ReviewTime.Year(), review.ReviewTime.Month(), review.ReviewTime.Day(), 0, 0, 0, 0, time.Local)
		if _, ok := dailyMap[date]; !ok {
			dailyMap[date] = &entity.DailyStat{
				Date:           date,
				NewLearned:     0,
				ReviewCount:    0,
				CorrectCount:   0,
				StudyTime:      0,
				MasteredCount:  0,
				RetentionRate:  0,
				ReviewInterval: 0,
			}
		}

		stat := dailyMap[date]
		stat.ReviewCount++
		if review.IsCorrect() {
			stat.CorrectCount++
		}
		stat.StudyTime += review.ResponseTime / 1000 // 转换为秒

		// 更新保持率
		if stat.ReviewCount > 0 {
			stat.RetentionRate = float32(stat.CorrectCount) / float32(stat.ReviewCount)
		}
	}

	// 转换为数组并排序
	for _, stat := range dailyMap {
		stats = append(stats, stat)
	}

	// 按日期排序
	sort.Slice(stats, func(i, j int) bool {
		return stats[i].Date.Before(stats[j].Date)
	})

	return stats, nil
}

// CalculateRetentionRate 计算记忆保持率
func (s *MemoryStatsService) CalculateRetentionRate(ctx context.Context, userID uint32, unitType entity.MemoryUnitType) (float32, error) {
	// 1. 获取指定类型的所有记忆单元
	units, err := s.memoryUnitRepo.ListByUserIDAndType(ctx, userID, unitType)
	if err != nil {
		return 0, err
	}

	// 2. 计算平均记忆保持率
	var totalRetentionRate float32
	var count uint32
	for _, unit := range units {
		if unit.RetentionRate > 0 {
			totalRetentionRate += unit.RetentionRate
			count++
		}
	}

	if count == 0 {
		return 0, nil
	}

	return totalRetentionRate / float32(count), nil
}
