package dto

import (
	"github.com/lazyjean/sla2/internal/domain/entity"
	"github.com/lazyjean/sla2/internal/domain/valueobject"
)

// CreateWordRequest 创建单词的请求数据
type CreateWordRequest struct {
	Text        string
	Phonetic    string
	Definitions []entity.Definition
	Examples    []string
	Tags        []string
	Level       valueobject.WordDifficultyLevel
}
