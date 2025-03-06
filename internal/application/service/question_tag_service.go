package service

import (
	"context"

	"github.com/google/wire"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/lazyjean/sla2/internal/domain/entity"
	"github.com/lazyjean/sla2/internal/domain/repository"
	"github.com/lazyjean/sla2/internal/domain/service"
)

var QuestionTagServiceSet = wire.NewSet(
	wire.Struct(new(QuestionTagService), "*"),
)

type QuestionTagService struct {
	logger *zap.Logger
	repo   repository.QuestionTagRepository
	domain service.QuestionTagService
}

func NewQuestionTagService(
	logger *zap.Logger,
	repo repository.QuestionTagRepository,
	domain service.QuestionTagService,
) *QuestionTagService {
	return &QuestionTagService{
		logger: logger,
		repo:   repo,
		domain: domain,
	}
}

func (s *QuestionTagService) CreateQuestionTag(ctx context.Context, req *entity.QuestionTag) (*entity.QuestionTag, error) {
	if err := s.domain.ValidateQuestionTag(ctx, req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	tag, err := s.repo.Create(ctx, req)
	if err != nil {
		s.logger.Error("failed to create question tag", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to create question tag")
	}

	return tag, nil
}

func (s *QuestionTagService) GetQuestionTag(ctx context.Context, id string) (*entity.QuestionTag, error) {
	tag, err := s.repo.Get(ctx, id)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, status.Error(codes.NotFound, "question tag not found")
		}
		s.logger.Error("failed to get question tag", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to get question tag")
	}

	return tag, nil
}

func (s *QuestionTagService) UpdateQuestionTag(ctx context.Context, req *entity.QuestionTag) (*entity.QuestionTag, error) {
	if err := s.domain.ValidateQuestionTag(ctx, req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	tag, err := s.repo.Update(ctx, req)
	if err != nil {
		s.logger.Error("failed to update question tag", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to update question tag")
	}

	return tag, nil
}

func (s *QuestionTagService) DeleteQuestionTag(ctx context.Context, id string) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Error("failed to delete question tag", zap.Error(err))
		return status.Error(codes.Internal, "failed to delete question tag")
	}

	return nil
}
