package service

import (
	"context"
	"errors"

	"github.com/lazyjean/sla2/internal/domain/entity"
)

var (
	ErrNotFound = errors.New("not found")
)

type QuestionTagService interface {
	Create(ctx context.Context, tag *entity.QuestionTag) error
	Get(ctx context.Context, id int64) (*entity.QuestionTag, error)
	Update(ctx context.Context, tag *entity.QuestionTag) error
	Delete(ctx context.Context, id int64) error
	ValidateQuestionTag(ctx context.Context, tag *entity.QuestionTag) error
}
