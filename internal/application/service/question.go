package service

import (
	"context"
	"errors"

	"github.com/lazyjean/sla2/internal/application/dto"
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
func (s *QuestionService) Get(ctx context.Context, id string) (*entity.Question, error) {
	if id == "" {
		return nil, errors.New("问题ID不能为空")
	}
	return s.questionRepo.Get(ctx, id)
}

// Create 创建新问题
func (s *QuestionService) Create(ctx context.Context, createDTO *dto.CreateQuestionDTO) (*entity.Question, error) {
	if createDTO.Title == "" || createDTO.Content == nil {
		return nil, errors.New("标题和内容不能为空")
	}

	question := entity.NewQuestion(
		createDTO.Title,
		createDTO.Content,
		createDTO.SimpleQuestion,
		createDTO.Type,
		createDTO.Difficulty,
		createDTO.Options,
		createDTO.OptionTuples,
		createDTO.Answers,
		createDTO.Category,
		createDTO.Labels,
		createDTO.Explanation,
		createDTO.Attachments,
		createDTO.TimeLimit,
	)

	if err := s.questionRepo.Create(ctx, question); err != nil {
		return nil, err
	}
	return question, nil
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
func (s *QuestionService) Update(ctx context.Context, updateDTO *dto.UpdateQuestionDTO) (*entity.Question, error) {
	if updateDTO.ID == "" {
		return nil, errors.New("问题ID不能为空")
	}
	if updateDTO.Title == "" || updateDTO.Content == nil {
		return nil, errors.New("标题和内容不能为空")
	}

	question, err := s.questionRepo.Get(ctx, updateDTO.ID)
	if err != nil {
		return nil, err
	}

	question.Update(
		updateDTO.Title,
		updateDTO.Content,
		updateDTO.SimpleQuestion,
		updateDTO.Type,
		updateDTO.Difficulty,
		updateDTO.Options,
		updateDTO.OptionTuples,
		updateDTO.Answers,
		updateDTO.Category,
		updateDTO.Labels,
		updateDTO.Explanation,
		updateDTO.Attachments,
		updateDTO.TimeLimit,
	)
	if err := s.questionRepo.Update(ctx, question); err != nil {
		return nil, err
	}
	return question, nil
}

// Delete 删除问题
func (s *QuestionService) Delete(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("问题ID不能为空")
	}

	question, err := s.questionRepo.Get(ctx, id)
	if err != nil {
		return err
	}

	question.Delete()
	return s.questionRepo.Update(ctx, question)
}

// Publish 发布问题
func (s *QuestionService) Publish(ctx context.Context, id string) (*entity.Question, error) {
	if id == "" {
		return nil, errors.New("问题ID不能为空")
	}

	question, err := s.questionRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	question.Publish()
	if err := s.questionRepo.Update(ctx, question); err != nil {
		return nil, err
	}
	return question, nil
}
