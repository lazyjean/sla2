package entity

import (
	"time"
)

// MemoryUnitType 记忆单元类型
type MemoryUnitType uint8

const (
	MemoryUnitTypeUnspecified MemoryUnitType = 0 // 未指定
	MemoryUnitTypeHanChar     MemoryUnitType = 1 // 汉字
	MemoryUnitTypeWord        MemoryUnitType = 2 // 单词
)

// MasteryLevel 掌握程度
type MasteryLevel uint8

const (
	MasteryLevelUnspecified MasteryLevel = 0 // 未指定
	MasteryLevelUnlearned   MasteryLevel = 1 // 未学习
	MasteryLevelBeginner    MasteryLevel = 2 // 初学
	MasteryLevelFamiliar    MasteryLevel = 3 // 熟悉
	MasteryLevelMastered    MasteryLevel = 4 // 掌握
	MasteryLevelExpert      MasteryLevel = 5 // 精通
)

// MemoryUnit 记忆单元
// 用于表示一个可记忆的学习内容，如汉字、单词等
// 表名：memory_units
// 注释：记忆单元表，存储各种类型的学习内容
type MemoryUnit struct {
	ID        uint32         `gorm:"primaryKey;comment:主键ID"`
	UserID    uint32         `gorm:"not null;index:idx_user_content,unique;comment:用户ID"`
	Type      MemoryUnitType `gorm:"not null;comment:记忆单元类型，0-未指定，1-汉字，2-单词"`
	ContentID uint32         `gorm:"not null;index:idx_user_content,unique;comment:内容ID，关联到具体的内容表（如汉字表、单词表等）"`
	CreatedAt time.Time      `gorm:"not null;comment:记录创建时间，由数据库自动维护"`
	UpdatedAt time.Time      `gorm:"not null;comment:记录更新时间，由数据库自动维护"`

	// 学习状态
	MasteryLevel       MasteryLevel `gorm:"not null;default:1;comment:掌握程度，0-未指定，1-未学习，2-初学，3-熟悉，4-掌握，5-精通"`
	ReviewCount        uint32       `gorm:"not null;default:0;comment:复习次数，累计复习的总次数"`
	NextReviewAt       time.Time    `gorm:"not null;comment:下次复习时间，根据记忆曲线计算得出"`
	LastReviewAt       time.Time    `gorm:"not null;comment:上次复习时间，记录最近一次复习的时间"`
	StudyDuration      uint32       `gorm:"not null;default:0;comment:学习时长（秒），累计学习该记忆单元的总时间"`
	RetentionRate      float32      `gorm:"not null;default:0;comment:记忆保持率（0-1），表示记忆的牢固程度"`
	ConsecutiveCorrect uint32       `gorm:"not null;default:0;comment:连续正确次数，用于评估记忆稳定性"`
	ConsecutiveWrong   uint32       `gorm:"not null;default:0;comment:连续错误次数，用于评估记忆难度"`
}

// NewMemoryUnit 创建新的记忆单元
func NewMemoryUnit(userID uint32, unitType MemoryUnitType, contentID uint32) *MemoryUnit {
	now := time.Now()
	return &MemoryUnit{
		UserID:        userID,
		Type:          unitType,
		ContentID:     contentID,
		CreatedAt:     now,
		UpdatedAt:     now,
		MasteryLevel:  MasteryLevelUnlearned,
		NextReviewAt:  now,
		LastReviewAt:  now,
		RetentionRate: 0,
	}
}

// Update 更新记忆单元
func (m *MemoryUnit) Update() {
	m.UpdatedAt = time.Now()
}

// IsDueForReview 检查是否需要复习
func (m *MemoryUnit) IsDueForReview() bool {
	return time.Now().After(m.NextReviewAt)
}

// UpdateReviewStats 更新复习统计
func (m *MemoryUnit) UpdateReviewStats(isCorrect bool, responseTime uint32) {
	m.ReviewCount++
	m.LastReviewAt = time.Now()
	m.StudyDuration += responseTime / 1000 // 转换为秒

	if isCorrect {
		m.ConsecutiveCorrect++
		m.ConsecutiveWrong = 0
	} else {
		m.ConsecutiveWrong++
		m.ConsecutiveCorrect = 0
	}

	// 更新记忆保持率
	if m.ReviewCount > 0 {
		m.RetentionRate = float32(m.ConsecutiveCorrect) / float32(m.ReviewCount)
	}

	// 更新掌握程度
	m.updateMasteryLevel()
}

// updateMasteryLevel 根据复习情况更新掌握程度
func (m *MemoryUnit) updateMasteryLevel() {
	switch {
	case m.ConsecutiveCorrect >= 10:
		m.MasteryLevel = MasteryLevelExpert
	case m.ConsecutiveCorrect >= 5:
		m.MasteryLevel = MasteryLevelMastered
	case m.ConsecutiveCorrect >= 3:
		m.MasteryLevel = MasteryLevelFamiliar
	case m.ConsecutiveCorrect >= 1:
		m.MasteryLevel = MasteryLevelBeginner
	case m.ConsecutiveWrong >= 3:
		m.MasteryLevel = MasteryLevelUnlearned
	}
}
