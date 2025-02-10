package entity

import "time"

// LearningProgress 学习进度实体
type LearningProgress struct {
	ID             uint      `gorm:"primaryKey"`
	UserID         uint      `gorm:"not null;index:idx_user_word"`
	WordID         uint      `gorm:"not null;index:idx_user_word"`
	Familiarity    int       `gorm:"not null;default:0"`
	NextReviewAt   time.Time `gorm:"not null"`
	LastReviewedAt time.Time `gorm:"not null"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// LearningStats 学习统计信息
type LearningStats struct {
	UserID          uint      // 用户ID
	TotalWords      int64     // 总单词数
	MasteredWords   int64     // 已掌握的单词数
	LearningWords   int64     // 正在学习的单词数
	ReviewDueCount  int64     // 待复习的单词数
	LastStudyTime   time.Time // 最后学习时间
	TodayStudyCount int64     // 今日学习数量
	ContinuousDays  int64     // 连续学习天数
}
