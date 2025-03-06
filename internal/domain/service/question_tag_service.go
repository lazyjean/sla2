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
	ValidateQuestionTag(ctx context.Context, tag *entity.QuestionTag) error
}
