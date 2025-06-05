package service

import (
	"context"
	"github.com/lazyjean/sla2/internal/domain/entity"
	"github.com/lazyjean/sla2/internal/domain/repository"
)

// QuestionService 问题服务
type QuestionService struct {
	questionRepo repository.QuestionRepository
}

// NewQuestionService 创建问题服务实例
func NewQuestionService(questionRepo repository.QuestionRepository) *QuestionService {
	return &QuestionService{
		questionRepo: questionRepo,
	}
}

// Get 获取问题详情
func (s *QuestionService) Get(ctx context.Context, id entity.QuestionID) (*entity.Question, error) {
	return s.questionRepo.GetByID(ctx, id)
}

// Create 创建新问题
func (s *QuestionService) Create(ctx context.Context, question *entity.Question) (entity.QuestionID, error) {
	if err := s.questionRepo.Create(ctx, question); err != nil {
		return entity.NullQuestionID, err
	}
	return question.ID, nil
}

// Search 搜索问题
func (s *QuestionService) Search(ctx context.Context, keyword string, labels []string, page, pageSize int) ([]*entity.Question, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	return s.questionRepo.Search(ctx, keyword, labels, page, pageSize)
}

// Update 更新问题
func (s *QuestionService) Update(ctx context.Context, question *entity.Question) error {
	return s.questionRepo.Update(ctx, question)
}

// Delete 删除问题
func (s *QuestionService) Delete(ctx context.Context, id entity.QuestionID) error {
	return s.questionRepo.Delete(ctx, id)
}

// Publish 发布问题
func (s *QuestionService) Publish(ctx context.Context, id entity.QuestionID) error {
	question, err := s.questionRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	question.Published()
	if err := s.questionRepo.Update(ctx, question); err != nil {
		return err
	}
	return nil
}
