package repository

import (
	"context"

	"github.com/lazyjean/sla2/domain/entity"
)

// WordRepository 单词仓储接口
type WordRepository interface {
	Save(ctx context.Context, word *entity.Word) error
	FindByID(ctx context.Context, id uint) (*entity.Word, error)
	FindByText(ctx context.Context, text string) (*entity.Word, error)
	List(ctx context.Context, offset, limit int) ([]*entity.Word, int64, error)
	Update(ctx context.Context, word *entity.Word) error
	Delete(ctx context.Context, id uint) error
	Search(ctx context.Context, query *WordQuery) ([]*entity.Word, int64, error)
}

type WordQuery struct {
	Text          string
	Tags          []string
	MinDifficulty int
	MaxDifficulty int
	Offset        int
	Limit         int
}
