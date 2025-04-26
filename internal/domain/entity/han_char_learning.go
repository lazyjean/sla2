package entity

import (
	"time"
)

// HanCharLearningID 汉字学习记录ID类型
type HanCharLearningID uint32

// HanCharLearning 汉字学习进度
type HanCharLearning struct {
	// ID 唯一标识符
	ID HanCharLearningID `gorm:"primaryKey;autoIncrement;comment:唯一标识符"`
	// UserID 用户ID
	UserID UID `gorm:"index;comment:用户ID"`
	// HanCharID 汉字ID
	HanCharID HanCharID `gorm:"index;comment:汉字ID"`
	// MemoryUnitID 记忆单元ID
	MemoryUnitID uint32 `gorm:"index;comment:记忆单元ID"`
	// FirstTryCorrect 第一次是否做对
	FirstTryCorrect bool `gorm:"comment:第一次是否做对"`
	// SecondTryCorrect 第二次是否做对
	SecondTryCorrect bool `gorm:"comment:第二次是否做对"`
	// ThirdTryCorrect 第三次是否做对
	ThirdTryCorrect bool `gorm:"comment:第三次是否做对"`
	// Mastered 是否掌握
	Mastered bool `gorm:"comment:是否掌握"`
	// ErrorCount 错误次数
	ErrorCount uint32 `gorm:"comment:错误次数"`
	// CorrectCount 正确次数
	CorrectCount uint32 `gorm:"comment:正确次数"`
	// StudyDuration 学习时长（秒）
	StudyDuration uint32 `gorm:"comment:学习时长（秒）"`
	// Notes 学习笔记
	Notes []string `gorm:"type:jsonb;serializer:json;comment:学习笔记"`
	// CreatedAt 创建时间
	CreatedAt time.Time `gorm:"type:timestamp with time zone;not null;comment:创建时间"`
	// UpdatedAt 更新时间
	UpdatedAt time.Time `gorm:"type:timestamp with time zone;not null;comment:更新时间"`
}

// NewHanCharLearning 创建新的汉字学习记录
func NewHanCharLearning(
	userID UID,
	hanCharID HanCharID,
	memoryUnitID uint32,
	firstTryCorrect bool,
	secondTryCorrect bool,
	thirdTryCorrect bool,
	mastered bool,
	errorCount uint32,
	correctCount uint32,
	studyDuration uint32,
	notes []string,
) *HanCharLearning {
	now := time.Now()
	return &HanCharLearning{
		UserID:           userID,
		HanCharID:        hanCharID,
		MemoryUnitID:     memoryUnitID,
		FirstTryCorrect:  firstTryCorrect,
		SecondTryCorrect: secondTryCorrect,
		ThirdTryCorrect:  thirdTryCorrect,
		Mastered:         mastered,
		ErrorCount:       errorCount,
		CorrectCount:     correctCount,
		StudyDuration:    studyDuration,
		Notes:            notes,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
}

// UpdateLearningResult 更新学习结果
func (h *HanCharLearning) UpdateLearningResult(
	firstTryCorrect bool,
	secondTryCorrect bool,
	thirdTryCorrect bool,
	mastered bool,
	errorCount uint32,
	correctCount uint32,
	studyDuration uint32,
	notes []string,
) {
	h.FirstTryCorrect = firstTryCorrect
	h.SecondTryCorrect = secondTryCorrect
	h.ThirdTryCorrect = thirdTryCorrect
	h.Mastered = mastered
	h.ErrorCount = errorCount
	h.CorrectCount = correctCount
	h.StudyDuration = studyDuration
	h.Notes = notes
	h.UpdatedAt = time.Now()
}

// IsFullyMastered 判断是否完全掌握
func (h *HanCharLearning) IsFullyMastered() bool {
	return h.FirstTryCorrect && h.SecondTryCorrect && h.ThirdTryCorrect && h.Mastered
}

// GetLearningEfficiency 获取学习效率
func (h *HanCharLearning) GetLearningEfficiency() float64 {
	if h.CorrectCount+h.ErrorCount == 0 {
		return 0
	}
	return float64(h.CorrectCount) / float64(h.CorrectCount+h.ErrorCount)
}

// GetLearningTimePerChar 获取每个汉字的平均学习时间
func (h *HanCharLearning) GetLearningTimePerChar() float64 {
	if h.CorrectCount+h.ErrorCount == 0 {
		return 0
	}
	return float64(h.StudyDuration) / float64(h.CorrectCount+h.ErrorCount)
}
