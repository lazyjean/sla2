package service

import (
	"context"

	"github.com/lazyjean/sla2/internal/domain/entity"
	"github.com/lazyjean/sla2/internal/domain/repository"
	"github.com/lazyjean/sla2/pkg/logger"
	"go.uber.org/zap"
)

// QuestionTagService 提供问题标签相关的应用层服务
type QuestionTagService struct {
	tagRepo repository.QuestionTagRepository
}

// NewQuestionTagService 创建问题标签服务
func NewQuestionTagService(tagRepo repository.QuestionTagRepository) *QuestionTagService {
	return &QuestionTagService{
		tagRepo: tagRepo,
	}
}

// FindAll 查询所有标签，限制返回数量
func (s *QuestionTagService) FindAll(ctx context.Context, limit int) ([]*entity.QuestionTag, error) {
	log := logger.GetLogger(ctx)
	log.Debug("查询所有标签", zap.Int("limit", limit))
	return s.tagRepo.FindAll(ctx, limit)
}

// Create 创建标签
func (s *QuestionTagService) Create(ctx context.Context, name string) (*entity.QuestionTag, error) {
	log := logger.GetLogger(ctx)
	log.Debug("创建标签", zap.String("name", name))

	tag := &entity.QuestionTag{
		Name: name,
	}

	return s.tagRepo.Create(ctx, tag)
}

// Update 更新标签
func (s *QuestionTagService) Update(ctx context.Context, id string, name string) (*entity.QuestionTag, error) {
	log := logger.GetLogger(ctx)
	log.Debug("更新标签", zap.String("id", id), zap.String("name", name))

	tag := &entity.QuestionTag{
		ID:   id,
		Name: name,
	}

	return s.tagRepo.Update(ctx, tag)
}

// Delete 删除标签
func (s *QuestionTagService) Delete(ctx context.Context, id string) error {
	log := logger.GetLogger(ctx)
	log.Debug("删除标签", zap.String("id", id))

	return s.tagRepo.Delete(ctx, id)
}
