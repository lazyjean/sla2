package entity

import (
	"time"
)

// MemoryStats 记忆统计
// 用于记录和展示用户的学习统计信息
// 表名：memory_stats
// 注释：记忆统计表，存储用户的学习进度和效果数据
type MemoryStats struct {
	TotalLearned    uint32             `json:"total_learned" gorm:"comment:已学习总数，包括所有已开始学习的记忆单元"`
	MasteredCount   uint32             `json:"mastered_count" gorm:"comment:已掌握数量，达到最高掌握程度的记忆单元数量"`
	NeedReviewCount uint32             `json:"need_review_count" gorm:"comment:需要复习数量，当前需要复习的记忆单元数量"`
	TotalStudyTime  uint32             `json:"total_study_time" gorm:"comment:总学习时长（秒），累计学习时间"`
	LevelStats      map[uint32]uint32  `json:"level_stats" gorm:"type:jsonb;comment:各掌握程度统计，key为掌握程度，value为对应数量"`
	RetentionRates  map[uint32]float32 `json:"retention_rates" gorm:"type:jsonb;comment:各类型记忆保持率，key为记忆单元类型，value为保持率"`
	DailyStats      []DailyStat        `json:"daily_stats" gorm:"type:jsonb;comment:每日学习统计，记录每天的学习情况"`
}

// DailyStat 每日学习统计
// 用于记录用户每天的学习情况
type DailyStat struct {
	Date           time.Time `json:"date" gorm:"type:timestamptz;comment:统计日期"`
	NewLearned     uint32    `json:"new_learned" gorm:"comment:新学习数量，当天新学习的记忆单元数量"`
	ReviewCount    uint32    `json:"review_count" gorm:"comment:复习数量，当天复习的记忆单元数量"`
	CorrectCount   uint32    `json:"correct_count" gorm:"comment:正确数量，当天复习正确的记忆单元数量"`
	StudyTime      uint32    `json:"study_time" gorm:"comment:学习时长（秒），当天的学习时间"`
	MasteredCount  uint32    `json:"mastered_count" gorm:"comment:掌握数量，当天达到掌握程度的记忆单元数量"`
	RetentionRate  float32   `json:"retention_rate" gorm:"comment:保持率，当天复习的正确率"`
	ReviewInterval uint32    `json:"review_interval" gorm:"comment:平均复习间隔（天），当天复习的记忆单元的平均间隔"`
}

// NewMemoryStats 创建新的记忆统计
func NewMemoryStats() *MemoryStats {
	return &MemoryStats{
		LevelStats:     make(map[uint32]uint32),
		RetentionRates: make(map[uint32]float32),
		DailyStats:     make([]DailyStat, 0),
	}
}

// AddDailyStat 添加每日统计
func (m *MemoryStats) AddDailyStat(stat DailyStat) {
	m.DailyStats = append(m.DailyStats, stat)
}

// UpdateLevelStats 更新掌握程度统计
func (m *MemoryStats) UpdateLevelStats(level uint32, count uint32) {
	m.LevelStats[level] = count
}

// UpdateRetentionRate 更新记忆保持率
func (m *MemoryStats) UpdateRetentionRate(unitType uint32, rate float32) {
	m.RetentionRates[unitType] = rate
}

// CalculateMasteryRate 计算掌握率
func (m *MemoryStats) CalculateMasteryRate() float32 {
	if m.TotalLearned == 0 {
		return 0
	}
	return float32(m.MasteredCount) / float32(m.TotalLearned)
}

// CalculateAverageStudyTime 计算平均学习时长
func (m *MemoryStats) CalculateAverageStudyTime() float32 {
	if m.TotalLearned == 0 {
		return 0
	}
	return float32(m.TotalStudyTime) / float32(m.TotalLearned)
}

// CalculateAverageRetentionRate 计算平均记忆保持率
func (m *MemoryStats) CalculateAverageRetentionRate() float32 {
	if len(m.RetentionRates) == 0 {
		return 0
	}
	var total float32
	for _, rate := range m.RetentionRates {
		total += rate
	}
	return total / float32(len(m.RetentionRates))
}
