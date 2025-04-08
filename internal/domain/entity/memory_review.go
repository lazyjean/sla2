package entity

import (
	"time"
)

// ReviewResult 复习结果
type ReviewResult uint32

const (
	ReviewResultUnspecified ReviewResult = 0 // 未指定
	ReviewResultCorrect     ReviewResult = 1 // 正确
	ReviewResultWrong       ReviewResult = 2 // 错误
	ReviewResultSkip        ReviewResult = 3 // 跳过
)

// MemoryReview 记忆复习记录, todo: 这个记录保留一年即可
// 用于记录用户对记忆单元的复习情况，包括复习结果、响应时间等信息
// 表名：memory_reviews
// 注释：记忆复习记录表，记录用户对记忆单元的复习情况
type MemoryReview struct {
	ID           uint32       `gorm:"primaryKey;comment:主键ID"`
	MemoryUnitID uint32       `gorm:"not null;index;comment:记忆单元ID，关联到记忆单元表"`
	UserID       uint32       `gorm:"not null;index;comment:用户ID，关联到用户表"`
	Result       ReviewResult `gorm:"not null;comment:复习结果，0-未指定，1-正确，2-错误，3-跳过"`
	ResponseTime uint32       `gorm:"not null;comment:响应时间，单位毫秒，表示用户从看到题目到做出回答的时间"`
	ReviewTime   time.Time    `gorm:"not null;comment:实际的复习时间，表示用户进行复习的具体时间点"`
	CreatedAt    time.Time    `gorm:"not null;comment:记录创建时间，由数据库自动维护"`
}

// NewMemoryReview 创建新的复习记录
// memoryUnitID: 记忆单元ID
// userID: 用户ID
// result: 复习结果
// responseTime: 响应时间（毫秒）
func NewMemoryReview(
	memoryUnitID uint32,
	userID uint32,
	result ReviewResult,
	responseTime uint32,
) *MemoryReview {
	now := time.Now()
	return &MemoryReview{
		MemoryUnitID: memoryUnitID,
		UserID:       userID,
		Result:       result,
		ResponseTime: responseTime,
		ReviewTime:   now,
		CreatedAt:    now,
	}
}

// IsCorrect 判断是否正确
// 返回：如果复习结果为正确，返回true；否则返回false
func (m *MemoryReview) IsCorrect() bool {
	return m.Result == ReviewResultCorrect
}

// IsWrong 判断是否错误
// 返回：如果复习结果为错误，返回true；否则返回false
func (m *MemoryReview) IsWrong() bool {
	return m.Result == ReviewResultWrong
}

// IsSkip 判断是否跳过
// 返回：如果复习结果为跳过，返回true；否则返回false
func (m *MemoryReview) IsSkip() bool {
	return m.Result == ReviewResultSkip
}
