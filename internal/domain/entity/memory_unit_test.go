package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewMemoryUnit(t *testing.T) {
	userID := uint32(1)
	unitType := MemoryUnitTypeHanChar
	contentID := uint32(100)
	now := time.Now()

	unit := NewMemoryUnit(userID, unitType, contentID)

	assert.Equal(t, userID, unit.UserID)
	assert.Equal(t, unitType, unit.Type)
	assert.Equal(t, contentID, unit.ContentID)
	assert.Equal(t, MasteryLevelUnlearned, unit.MasteryLevel)
	assert.Equal(t, uint32(0), unit.ReviewCount)
	assert.Equal(t, uint32(0), unit.StudyDuration)
	assert.Equal(t, float32(0), unit.RetentionRate)
	assert.Equal(t, uint32(0), unit.ConsecutiveCorrect)
	assert.Equal(t, uint32(0), unit.ConsecutiveWrong)
	assert.True(t, unit.CreatedAt.After(now) || unit.CreatedAt.Equal(now))
	assert.True(t, unit.UpdatedAt.After(now) || unit.UpdatedAt.Equal(now))
	assert.True(t, unit.NextReviewAt.After(now) || unit.NextReviewAt.Equal(now))
	assert.True(t, unit.LastReviewAt.After(now) || unit.LastReviewAt.Equal(now))
}

func TestUpdate(t *testing.T) {
	unit := NewMemoryUnit(1, MemoryUnitTypeHanChar, 100)
	oldUpdatedAt := unit.UpdatedAt

	// 等待一小段时间以确保时间戳会变化
	time.Sleep(time.Millisecond)
	unit.Update()

	assert.True(t, unit.UpdatedAt.After(oldUpdatedAt))
}

func TestIsDueForReview(t *testing.T) {
	unit := NewMemoryUnit(1, MemoryUnitTypeHanChar, 100)

	// 设置下次复习时间为过去
	unit.NextReviewAt = time.Now().Add(-time.Hour)
	assert.True(t, unit.IsDueForReview())

	// 设置下次复习时间为未来
	unit.NextReviewAt = time.Now().Add(time.Hour)
	assert.False(t, unit.IsDueForReview())
}

func TestUpdateReviewStats(t *testing.T) {
	tests := []struct {
		name           string
		isCorrect      bool
		responseTime   uint32
		initialCorrect uint32
		initialWrong   uint32
		expectedLevel  MasteryLevel
	}{
		{
			name:           "第一次正确",
			isCorrect:      true,
			responseTime:   5000,
			initialCorrect: 0,
			initialWrong:   0,
			expectedLevel:  MasteryLevelBeginner,
		},
		{
			name:           "连续5次正确",
			isCorrect:      true,
			responseTime:   5000,
			initialCorrect: 4,
			initialWrong:   0,
			expectedLevel:  MasteryLevelMastered,
		},
		{
			name:           "连续10次正确",
			isCorrect:      true,
			responseTime:   5000,
			initialCorrect: 9,
			initialWrong:   0,
			expectedLevel:  MasteryLevelExpert,
		},
		{
			name:           "第一次错误",
			isCorrect:      false,
			responseTime:   5000,
			initialCorrect: 0,
			initialWrong:   0,
			expectedLevel:  MasteryLevelUnlearned,
		},
		{
			name:           "连续3次错误",
			isCorrect:      false,
			responseTime:   5000,
			initialCorrect: 0,
			initialWrong:   2,
			expectedLevel:  MasteryLevelUnlearned,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			unit := NewMemoryUnit(1, MemoryUnitTypeHanChar, 100)
			unit.ConsecutiveCorrect = tt.initialCorrect
			unit.ConsecutiveWrong = tt.initialWrong
			oldReviewCount := unit.ReviewCount
			oldStudyDuration := unit.StudyDuration

			unit.UpdateReviewStats(tt.isCorrect, tt.responseTime)

			// 检查基本统计
			assert.Equal(t, oldReviewCount+1, unit.ReviewCount)
			assert.Equal(t, oldStudyDuration+5, unit.StudyDuration) // 5000ms = 5s
			assert.True(t, unit.LastReviewAt.After(time.Now().Add(-time.Second)))

			// 检查连续正确/错误次数
			if tt.isCorrect {
				assert.Equal(t, tt.initialCorrect+1, unit.ConsecutiveCorrect)
				assert.Equal(t, uint32(0), unit.ConsecutiveWrong)
			} else {
				assert.Equal(t, tt.initialWrong+1, unit.ConsecutiveWrong)
				assert.Equal(t, uint32(0), unit.ConsecutiveCorrect)
			}

			// 检查掌握程度
			assert.Equal(t, tt.expectedLevel, unit.MasteryLevel)

			// 检查记忆保持率
			expectedRetentionRate := float32(unit.ConsecutiveCorrect) / float32(unit.ReviewCount)
			assert.Equal(t, expectedRetentionRate, unit.RetentionRate)
		})
	}
}

func TestUpdateMasteryLevel(t *testing.T) {
	tests := []struct {
		name               string
		consecutiveCorrect uint32
		consecutiveWrong   uint32
		expectedLevel      MasteryLevel
	}{
		{"未学习", 0, 0, MasteryLevelUnlearned},
		{"初学", 1, 0, MasteryLevelBeginner},
		{"熟悉", 3, 0, MasteryLevelFamiliar},
		{"掌握", 5, 0, MasteryLevelMastered},
		{"精通", 10, 0, MasteryLevelExpert},
		{"连续错误", 0, 3, MasteryLevelUnlearned},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			unit := NewMemoryUnit(1, MemoryUnitTypeHanChar, 100)
			unit.ConsecutiveCorrect = tt.consecutiveCorrect
			unit.ConsecutiveWrong = tt.consecutiveWrong

			unit.updateMasteryLevel()

			assert.Equal(t, tt.expectedLevel, unit.MasteryLevel)
		})
	}
}
