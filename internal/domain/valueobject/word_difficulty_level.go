package valueobject

import (
	"github.com/lazyjean/sla2/internal/domain/errors"
)

// WordDifficultyLevel 定义单词难度等级
type WordDifficultyLevel int32

const (
	// WORD_DIFFICULTY_LEVEL_UNSPECIFIED 未指定难度
	WORD_DIFFICULTY_LEVEL_UNSPECIFIED WordDifficultyLevel = 0
	// WORD_DIFFICULTY_LEVEL_A1 CEFR A1级别
	WORD_DIFFICULTY_LEVEL_A1 WordDifficultyLevel = 1
	// WORD_DIFFICULTY_LEVEL_A2 CEFR A2级别
	WORD_DIFFICULTY_LEVEL_A2 WordDifficultyLevel = 2
	// WORD_DIFFICULTY_LEVEL_B1 CEFR B1级别
	WORD_DIFFICULTY_LEVEL_B1 WordDifficultyLevel = 3
	// WORD_DIFFICULTY_LEVEL_B2 CEFR B2级别
	WORD_DIFFICULTY_LEVEL_B2 WordDifficultyLevel = 4
	// WORD_DIFFICULTY_LEVEL_C1 CEFR C1级别
	WORD_DIFFICULTY_LEVEL_C1 WordDifficultyLevel = 5
	// WORD_DIFFICULTY_LEVEL_C2 CEFR C2级别
	WORD_DIFFICULTY_LEVEL_C2 WordDifficultyLevel = 6
	// WORD_DIFFICULTY_LEVEL_HSK1 HSK 1级
	WORD_DIFFICULTY_LEVEL_HSK1 WordDifficultyLevel = 7
	// WORD_DIFFICULTY_LEVEL_HSK2 HSK 2级
	WORD_DIFFICULTY_LEVEL_HSK2 WordDifficultyLevel = 8
	// WORD_DIFFICULTY_LEVEL_HSK3 HSK 3级
	WORD_DIFFICULTY_LEVEL_HSK3 WordDifficultyLevel = 9
	// WORD_DIFFICULTY_LEVEL_HSK4 HSK 4级
	WORD_DIFFICULTY_LEVEL_HSK4 WordDifficultyLevel = 10
	// WORD_DIFFICULTY_LEVEL_HSK5 HSK 5级
	WORD_DIFFICULTY_LEVEL_HSK5 WordDifficultyLevel = 11
	// WORD_DIFFICULTY_LEVEL_HSK6 HSK 6级
	WORD_DIFFICULTY_LEVEL_HSK6 WordDifficultyLevel = 12
)

// String 返回难度等级的字符串表示
func (l WordDifficultyLevel) String() string {
	switch l {
	case WORD_DIFFICULTY_LEVEL_UNSPECIFIED:
		return "未指定"
	case WORD_DIFFICULTY_LEVEL_A1:
		return "A1"
	case WORD_DIFFICULTY_LEVEL_A2:
		return "A2"
	case WORD_DIFFICULTY_LEVEL_B1:
		return "B1"
	case WORD_DIFFICULTY_LEVEL_B2:
		return "B2"
	case WORD_DIFFICULTY_LEVEL_C1:
		return "C1"
	case WORD_DIFFICULTY_LEVEL_C2:
		return "C2"
	case WORD_DIFFICULTY_LEVEL_HSK1:
		return "HSK1"
	case WORD_DIFFICULTY_LEVEL_HSK2:
		return "HSK2"
	case WORD_DIFFICULTY_LEVEL_HSK3:
		return "HSK3"
	case WORD_DIFFICULTY_LEVEL_HSK4:
		return "HSK4"
	case WORD_DIFFICULTY_LEVEL_HSK5:
		return "HSK5"
	case WORD_DIFFICULTY_LEVEL_HSK6:
		return "HSK6"
	default:
		return "未知"
	}
}

// IsValid 检查难度等级是否有效
func (l WordDifficultyLevel) IsValid() bool {
	return l >= WORD_DIFFICULTY_LEVEL_UNSPECIFIED && l <= WORD_DIFFICULTY_LEVEL_HSK6
}

// IsCEFR 检查是否是CEFR标准
func (l WordDifficultyLevel) IsCEFR() bool {
	return l >= WORD_DIFFICULTY_LEVEL_A1 && l <= WORD_DIFFICULTY_LEVEL_C2
}

// IsHSK 检查是否是HSK标准
func (l WordDifficultyLevel) IsHSK() bool {
	return l >= WORD_DIFFICULTY_LEVEL_HSK1 && l <= WORD_DIFFICULTY_LEVEL_HSK6
}

// ParseWordDifficultyLevel 解析难度等级字符串
func ParseWordDifficultyLevel(level string) (WordDifficultyLevel, error) {
	switch level {
	case WORD_DIFFICULTY_LEVEL_A1.String():
		return WORD_DIFFICULTY_LEVEL_A1, nil
	case WORD_DIFFICULTY_LEVEL_A2.String():
		return WORD_DIFFICULTY_LEVEL_A2, nil
	case WORD_DIFFICULTY_LEVEL_B1.String():
		return WORD_DIFFICULTY_LEVEL_B1, nil
	case WORD_DIFFICULTY_LEVEL_B2.String():
		return WORD_DIFFICULTY_LEVEL_B2, nil
	case WORD_DIFFICULTY_LEVEL_C1.String():
		return WORD_DIFFICULTY_LEVEL_C1, nil
	case WORD_DIFFICULTY_LEVEL_C2.String():
		return WORD_DIFFICULTY_LEVEL_C2, nil
	case WORD_DIFFICULTY_LEVEL_HSK1.String():
		return WORD_DIFFICULTY_LEVEL_HSK1, nil
	case WORD_DIFFICULTY_LEVEL_HSK2.String():
		return WORD_DIFFICULTY_LEVEL_HSK2, nil
	case WORD_DIFFICULTY_LEVEL_HSK3.String():
		return WORD_DIFFICULTY_LEVEL_HSK3, nil
	case WORD_DIFFICULTY_LEVEL_HSK4.String():
		return WORD_DIFFICULTY_LEVEL_HSK4, nil
	case WORD_DIFFICULTY_LEVEL_HSK5.String():
		return WORD_DIFFICULTY_LEVEL_HSK5, nil
	case WORD_DIFFICULTY_LEVEL_HSK6.String():
		return WORD_DIFFICULTY_LEVEL_HSK6, nil
	default:
		return WORD_DIFFICULTY_LEVEL_UNSPECIFIED, errors.ErrInvalidDifficultyLevel
	}
}
