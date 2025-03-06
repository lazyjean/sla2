package service

import (
	"context"
	"errors"
	"testing"

	"github.com/lazyjean/sla2/internal/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockQuestionRepository 是 QuestionRepository 的模拟实现
type MockQuestionRepository struct {
	mock.Mock
}

func (m *MockQuestionRepository) Get(ctx context.Context, id string) (*entity.Question, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Question), args.Error(1)
}

func (m *MockQuestionRepository) Create(ctx context.Context, question *entity.Question) error {
	args := m.Called(ctx, question)
	return args.Error(0)
}

func (m *MockQuestionRepository) Update(ctx context.Context, question *entity.Question) error {
	args := m.Called(ctx, question)
	return args.Error(0)
}

func (m *MockQuestionRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockQuestionRepository) Search(ctx context.Context, keyword string, tags []string, page, pageSize int) ([]*entity.Question, int64, error) {
	args := m.Called(ctx, keyword, tags, page, pageSize)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*entity.Question), args.Get(1).(int64), args.Error(2)
}

// TestQuestionService_Get 测试获取问题详情
func TestQuestionService_Get(t *testing.T) {
	mockRepo := new(MockQuestionRepository)
	service := NewQuestionService(mockRepo)
	ctx := context.Background()

	t.Run("成功获取问题", func(t *testing.T) {
		expectedQuestion := &entity.Question{ID: entity.QuestionID(1), Title: "测试问题"}
		mockRepo.On("Get", ctx, "1").Return(expectedQuestion, nil).Once()

		question, err := service.Get(ctx, "1")
		assert.NoError(t, err)
		assert.Equal(t, expectedQuestion, question)
	})

	t.Run("问题ID为空", func(t *testing.T) {
		question, err := service.Get(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, question)
		assert.Equal(t, "问题ID不能为空", err.Error())
	})

	t.Run("问题不存在", func(t *testing.T) {
		mockRepo.On("Get", ctx, "999").Return(nil, errors.New("问题不存在")).Once()

		question, err := service.Get(ctx, "999")
		assert.Error(t, err)
		assert.Nil(t, question)
	})
}

// TestQuestionService_Create 测试创建新问题
func TestQuestionService_Create(t *testing.T) {
	mockRepo := new(MockQuestionRepository)
	service := NewQuestionService(mockRepo)
	ctx := context.Background()

	t.Run("成功创建问题", func(t *testing.T) {
		mockRepo.On("Create", ctx, mock.AnythingOfType("*entity.Question")).Return(nil).Once()

		question, err := service.Create(ctx, "测试标题", "测试内容", []string{"标签1"}, "user1")
		assert.NoError(t, err)
		assert.NotNil(t, question)
		assert.Equal(t, "测试标题", question.Title)
		assert.Equal(t, "测试内容", question.Content)
	})

	t.Run("标题为空", func(t *testing.T) {
		question, err := service.Create(ctx, "", "测试内容", []string{"标签1"}, "user1")
		assert.Error(t, err)
		assert.Nil(t, question)
		assert.Equal(t, "标题和内容不能为空", err.Error())
	})

	t.Run("内容为空", func(t *testing.T) {
		question, err := service.Create(ctx, "测试标题", "", []string{"标签1"}, "user1")
		assert.Error(t, err)
		assert.Nil(t, question)
		assert.Equal(t, "标题和内容不能为空", err.Error())
	})

	t.Run("创建失败", func(t *testing.T) {
		mockRepo.On("Create", ctx, mock.AnythingOfType("*entity.Question")).Return(errors.New("创建失败")).Once()

		question, err := service.Create(ctx, "测试标题", "测试内容", []string{"标签1"}, "user1")
		assert.Error(t, err)
		assert.Nil(t, question)
	})
}

// TestQuestionService_Search 测试搜索问题
func TestQuestionService_Search(t *testing.T) {
	mockRepo := new(MockQuestionRepository)
	service := NewQuestionService(mockRepo)
	ctx := context.Background()

	t.Run("成功搜索问题", func(t *testing.T) {
		expectedQuestions := []*entity.Question{{ID: entity.QuestionID(1), Title: "测试问题1"}, {ID: entity.QuestionID(2), Title: "测试问题2"}}
		mockRepo.On("Search", ctx, "测试", []string{"标签1"}, 1, 10).Return(expectedQuestions, int64(2), nil).Once()

		questions, total, err := service.Search(ctx, "测试", []string{"标签1"}, 1, 10)
		assert.NoError(t, err)
		assert.Equal(t, expectedQuestions, questions)
		assert.Equal(t, int64(2), total)
	})

	t.Run("页码小于1时自动修正", func(t *testing.T) {
		expectedQuestions := []*entity.Question{}
		mockRepo.On("Search", ctx, "测试", []string{"标签1"}, 1, 10).Return(expectedQuestions, int64(0), nil).Once()

		questions, total, err := service.Search(ctx, "测试", []string{"标签1"}, 0, 10)
		assert.NoError(t, err)
		assert.Empty(t, questions)
		assert.Equal(t, int64(0), total)
	})

	t.Run("每页数量超出限制时自动修正", func(t *testing.T) {
		expectedQuestions := []*entity.Question{}
		mockRepo.On("Search", ctx, "测试", []string{"标签1"}, 1, 10).Return(expectedQuestions, int64(0), nil).Once()

		questions, total, err := service.Search(ctx, "测试", []string{"标签1"}, 1, 200)
		assert.NoError(t, err)
		assert.Empty(t, questions)
		assert.Equal(t, int64(0), total)
	})
}

// TestQuestionService_Update 测试更新问题
func TestQuestionService_Update(t *testing.T) {
	mockRepo := new(MockQuestionRepository)
	service := NewQuestionService(mockRepo)
	ctx := context.Background()

	t.Run("成功更新问题", func(t *testing.T) {
		originalQuestion := &entity.Question{ID: entity.QuestionID(1), Title: "原标题", Content: "原内容"}
		mockRepo.On("Get", ctx, "1").Return(originalQuestion, nil).Once()
		mockRepo.On("Update", ctx, mock.AnythingOfType("*entity.Question")).Return(nil).Once()

		question, err := service.Update(ctx, "1", "新标题", "新内容", []string{"新标签"})
		assert.NoError(t, err)
		assert.Equal(t, "新标题", question.Title)
		assert.Equal(t, "新内容", question.Content)
	})

	t.Run("问题ID为空", func(t *testing.T) {
		question, err := service.Update(ctx, "", "新标题", "新内容", []string{"新标签"})
		assert.Error(t, err)
		assert.Nil(t, question)
		assert.Equal(t, "问题ID不能为空", err.Error())
	})

	t.Run("标题为空", func(t *testing.T) {
		question, err := service.Update(ctx, "1", "", "新内容", []string{"新标签"})
		assert.Error(t, err)
		assert.Nil(t, question)
		assert.Equal(t, "标题和内容不能为空", err.Error())
	})

	t.Run("问题不存在", func(t *testing.T) {
		mockRepo.On("Get", ctx, "999").Return(nil, errors.New("问题不存在")).Once()

		question, err := service.Update(ctx, "999", "新标题", "新内容", []string{"新标签"})
		assert.Error(t, err)
		assert.Nil(t, question)
	})
}

// TestQuestionService_Delete 测试删除问题
func TestQuestionService_Delete(t *testing.T) {
	mockRepo := new(MockQuestionRepository)
	service := NewQuestionService(mockRepo)
	ctx := context.Background()

	t.Run("成功删除问题", func(t *testing.T) {
		question := &entity.Question{ID: entity.QuestionID(1), Title: "测试问题"}
		mockRepo.On("Get", ctx, "1").Return(question, nil).Once()
		mockRepo.On("Update", ctx, mock.AnythingOfType("*entity.Question")).Return(nil).Once()

		err := service.Delete(ctx, "1")
		assert.NoError(t, err)
	})

	t.Run("问题ID为空", func(t *testing.T) {
		err := service.Delete(ctx, "")
		assert.Error(t, err)
		assert.Equal(t, "问题ID不能为空", err.Error())
	})

	t.Run("问题不存在", func(t *testing.T) {
		mockRepo.On("Get", ctx, "999").Return(nil, errors.New("问题不存在")).Once()

		err := service.Delete(ctx, "999")
		assert.Error(t, err)
	})
}

// TestQuestionService_Publish 测试发布问题
func TestQuestionService_Publish(t *testing.T) {
	mockRepo := new(MockQuestionRepository)
	service := NewQuestionService(mockRepo)
	ctx := context.Background()

	t.Run("成功发布问题", func(t *testing.T) {
		question := &entity.Question{ID: entity.QuestionID(1), Title: "测试问题", Status: "draft"}
		mockRepo.On("Get", ctx, "1").Return(question, nil).Once()
		mockRepo.On("Update", ctx, mock.AnythingOfType("*entity.Question")).Return(nil).Once()

		updatedQuestion, err := service.Publish(ctx, "1")
		assert.NoError(t, err)
		assert.Equal(t, "published", updatedQuestion.Status)
	})

	t.Run("问题ID为空", func(t *testing.T) {
		question, err := service.Publish(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, question)
		assert.Equal(t, "问题ID不能为空", err.Error())
	})

	t.Run("问题不存在", func(t *testing.T) {
		mockRepo.On("Get", ctx, "999").Return(nil, errors.New("问题不存在")).Once()

		question, err := service.Publish(ctx, "999")
		assert.Error(t, err)
		assert.Nil(t, question)
	})

	t.Run("更新失败", func(t *testing.T) {
		question := &entity.Question{ID: entity.QuestionID(1), Title: "测试问题", Status: "draft"}
		mockRepo.On("Get", ctx, "1").Return(question, nil).Once()
		mockRepo.On("Update", ctx, mock.AnythingOfType("*entity.Question")).Return(errors.New("更新失败")).Once()

		updatedQuestion, err := service.Publish(ctx, "1")
		assert.Error(t, err)
		assert.Nil(t, updatedQuestion)
	})
}
