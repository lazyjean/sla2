package repository

import (
	"context"
	"errors"

	"github.com/lazyjean/sla2/internal/domain/entity"
)

var (
	ErrNotFound = errors.New("not found")
)

type QuestionTagRepository interface {
	Create(ctx context.Context, tag *entity.QuestionTag) (*entity.QuestionTag, error)
	Get(ctx context.Context, id string) (*entity.QuestionTag, error)
	Update(ctx context.Context, tag *entity.QuestionTag) (*entity.QuestionTag, error)
	Delete(ctx context.Context, id string) error
}
