package dto

import "github.com/lazyjean/sla2/internal/domain/valueobject"

// ListHanCharsDTO 获取汉字列表请求
type ListHanCharsDTO struct {
	Page       int
	PageSize   int
	Level      valueobject.WordDifficultyLevel
	Tags       []string
	Categories []string
}
