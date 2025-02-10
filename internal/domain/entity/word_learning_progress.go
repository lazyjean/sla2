package entity

import "time"

// WordLearningProgress 单词学习进度
type WordLearningProgress struct {
	ID             uint      `gorm:"primaryKey"`
	UserID         uint      `gorm:"not null;index;uniqueIndex:idx_user_word,priority:1"`
	WordID         uint      `gorm:"not null;uniqueIndex:idx_user_word,priority:2"`
	Familiarity    int       `gorm:"not null;default:0"` // 熟悉度 0-5
	NextReviewAt   time.Time `gorm:"not null"`           // 下次复习时间
	LastReviewedAt time.Time `gorm:"not null"`           // 上次复习时间
	ReviewCount    int       `gorm:"not null;default:0"` // 复习次数
	CreatedAt      time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt      time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
}

// NewWordLearningProgress 创建新的单词学习进度
func NewWordLearningProgress(userID, wordID uint, familiarity int, nextReviewAt time.Time) *WordLearningProgress {
	now := time.Now()
	return &WordLearningProgress{
		UserID:         userID,
		WordID:         wordID,
		Familiarity:    familiarity,
		NextReviewAt:   nextReviewAt,
		LastReviewedAt: now,
		ReviewCount:    1,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

// UpdateProgress 更新学习进度
func (p *WordLearningProgress) UpdateProgress(familiarity int, nextReviewAt time.Time) {
	now := time.Now()
	p.Familiarity = familiarity
	p.NextReviewAt = nextReviewAt
	p.LastReviewedAt = now
	p.ReviewCount++
	p.UpdatedAt = now
}
