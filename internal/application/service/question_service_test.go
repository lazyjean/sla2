package service

import (
	"context"
	"errors"
	"testing"

	"github.com/lazyjean/sla2/internal/application/dto"
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

func (m *MockQuestionRepository) CreateTag(ctx context.Context, tag *entity.QuestionTag) (*entity.QuestionTag, error) {
	args := m.Called(ctx, tag)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.QuestionTag), args.Error(1)
}

func (m *MockQuestionRepository) GetTag(ctx context.Context, id string) (*entity.QuestionTag, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.QuestionTag), args.Error(1)
}

func (m *MockQuestionRepository) UpdateTag(ctx context.Context, tag *entity.QuestionTag) (*entity.QuestionTag, error) {
	args := m.Called(ctx, tag)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.QuestionTag), args.Error(1)
}

func (m *MockQuestionRepository) DeleteTag(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockQuestionRepository) FindAllTags(ctx context.Context, limit int) ([]*entity.QuestionTag, error) {
	args := m.Called(ctx, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.QuestionTag), args.Error(1)
}

// TestQuestionService_Get 测试获取问题详情
func TestQuestionService_Get(t *testing.T) {
	mockQuestionRepo := new(MockQuestionRepository)
	service := NewQuestionService(mockQuestionRepo)
	ctx := context.Background()

	t.Run("成功获取问题", func(t *testing.T) {
		expectedQuestion := &entity.Question{ID: entity.QuestionID(1), Title: "测试问题"}
		mockQuestionRepo.On("Get", ctx, "1").Return(expectedQuestion, nil).Once()

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
		mockQuestionRepo.On("Get", ctx, "999").Return(nil, errors.New("问题不存在")).Once()

		question, err := service.Get(ctx, "999")
		assert.Error(t, err)
		assert.Nil(t, question)
	})
}

// TestQuestionService_Create 测试创建新问题
func TestQuestionService_Create(t *testing.T) {
	mockQuestionRepo := new(MockQuestionRepository)
	service := NewQuestionService(mockQuestionRepo)
	ctx := context.Background()

	t.Run("成功创建问题", func(t *testing.T) {
		mockQuestionRepo.On("Create", ctx, mock.AnythingOfType("*entity.Question")).Return(nil).Once()

		createDTO := &dto.CreateQuestionDTO{
			Title:          "测试标题",
			Content:        []byte("测试内容"),
			SimpleQuestion: "测试内容",
			Type:           "single_choice",
			Difficulty:     "easy",
			Options:        []byte("[]"),
			OptionTuples:   []byte("[]"),
			Answers:        []string{"A"},
			Category:       "grammar",
			Labels:         []string{"标签1"},
			Explanation:    "解析",
			Attachments:    []string{},
			TimeLimit:      60,
		}

		question, err := service.Create(ctx, createDTO)
		assert.NoError(t, err)
		assert.NotNil(t, question)
		assert.Equal(t, "测试标题", question.Title)
		assert.Equal(t, []byte("测试内容"), question.Content)
	})

	t.Run("标题为空", func(t *testing.T) {
		createDTO := &dto.CreateQuestionDTO{
			Title:          "",
			Content:        []byte("测试内容"),
			SimpleQuestion: "测试内容",
			Type:           "single_choice",
			Difficulty:     "easy",
			Options:        []byte("[]"),
			OptionTuples:   []byte("[]"),
			Answers:        []string{"A"},
			Category:       "grammar",
			Labels:         []string{"标签1"},
			Explanation:    "解析",
			Attachments:    []string{},
			TimeLimit:      60,
		}

		question, err := service.Create(ctx, createDTO)
		assert.Error(t, err)
		assert.Nil(t, question)
		assert.Equal(t, "标题和内容不能为空", err.Error())
	})

	t.Run("内容为空", func(t *testing.T) {
		createDTO := &dto.CreateQuestionDTO{
			Title:          "测试标题",
			Content:        nil,
			SimpleQuestion: "测试内容",
			Type:           "single_choice",
			Difficulty:     "easy",
			Options:        []byte("[]"),
			OptionTuples:   []byte("[]"),
			Answers:        []string{"A"},
			Category:       "grammar",
			Labels:         []string{"标签1"},
			Explanation:    "解析",
			Attachments:    []string{},
			TimeLimit:      60,
		}

		question, err := service.Create(ctx, createDTO)
		assert.Error(t, err)
		assert.Nil(t, question)
		assert.Equal(t, "标题和内容不能为空", err.Error())
	})

	t.Run("创建失败", func(t *testing.T) {
		mockQuestionRepo.On("Create", ctx, mock.AnythingOfType("*entity.Question")).Return(errors.New("创建失败")).Once()

		createDTO := &dto.CreateQuestionDTO{
			Title:          "测试标题",
			Content:        []byte("测试内容"),
			SimpleQuestion: "测试内容",
			Type:           "single_choice",
			Difficulty:     "easy",
			Options:        []byte("[]"),
			OptionTuples:   []byte("[]"),
			Answers:        []string{"A"},
			Category:       "grammar",
			Labels:         []string{"标签1"},
			Explanation:    "解析",
			Attachments:    []string{},
			TimeLimit:      60,
		}

		question, err := service.Create(ctx, createDTO)
		assert.Error(t, err)
		assert.Nil(t, question)
	})
}

// TestQuestionService_Search 测试搜索问题
func TestQuestionService_Search(t *testing.T) {
	mockQuestionRepo := new(MockQuestionRepository)
	service := NewQuestionService(mockQuestionRepo)
	ctx := context.Background()

	t.Run("成功搜索问题", func(t *testing.T) {
		expectedQuestions := []*entity.Question{{ID: entity.QuestionID(1), Title: "测试问题1"}, {ID: entity.QuestionID(2), Title: "测试问题2"}}
		mockQuestionRepo.On("Search", ctx, "测试", []string{"标签1"}, 1, 10).Return(expectedQuestions, int64(2), nil).Once()

		questions, total, err := service.Search(ctx, "测试", []string{"标签1"}, 1, 10)
		assert.NoError(t, err)
		assert.Equal(t, expectedQuestions, questions)
		assert.Equal(t, int64(2), total)
	})

	t.Run("页码小于1时自动修正", func(t *testing.T) {
		expectedQuestions := []*entity.Question{}
		mockQuestionRepo.On("Search", ctx, "测试", []string{"标签1"}, 1, 10).Return(expectedQuestions, int64(0), nil).Once()

		questions, total, err := service.Search(ctx, "测试", []string{"标签1"}, 0, 10)
		assert.NoError(t, err)
		assert.Empty(t, questions)
		assert.Equal(t, int64(0), total)
	})

	t.Run("每页数量超出限制时自动修正", func(t *testing.T) {
		expectedQuestions := []*entity.Question{}
		mockQuestionRepo.On("Search", ctx, "测试", []string{"标签1"}, 1, 10).Return(expectedQuestions, int64(0), nil).Once()

		questions, total, err := service.Search(ctx, "测试", []string{"标签1"}, 1, 200)
		assert.NoError(t, err)
		assert.Empty(t, questions)
		assert.Equal(t, int64(0), total)
	})
}

// TestQuestionService_Update 测试更新问题
func TestQuestionService_Update(t *testing.T) {
	mockQuestionRepo := new(MockQuestionRepository)
	service := NewQuestionService(mockQuestionRepo)
	ctx := context.Background()

	t.Run("成功更新问题", func(t *testing.T) {
		originalQuestion := &entity.Question{ID: entity.QuestionID(1), Title: "原标题", Content: []byte("原内容")}
		mockQuestionRepo.On("Get", ctx, "1").Return(originalQuestion, nil).Once()
		mockQuestionRepo.On("Update", ctx, mock.AnythingOfType("*entity.Question")).Return(nil).Once()

		updateDTO := &dto.UpdateQuestionDTO{
			ID:             "1",
			Title:          "新标题",
			Content:        []byte("新内容"),
			SimpleQuestion: "新内容",
			Type:           "single_choice",
			Difficulty:     "easy",
			Options:        []byte("[]"),
			OptionTuples:   []byte("[]"),
			Answers:        []string{"A"},
			Category:       "grammar",
			Labels:         []string{"新标签"},
			Explanation:    "解析",
			Attachments:    []string{},
			TimeLimit:      60,
		}

		question, err := service.Update(ctx, updateDTO)
		assert.NoError(t, err)
		assert.Equal(t, "新标题", question.Title)
		assert.Equal(t, []byte("新内容"), question.Content)
	})

	t.Run("问题ID为空", func(t *testing.T) {
		updateDTO := &dto.UpdateQuestionDTO{
			ID:             "",
			Title:          "新标题",
			Content:        []byte("新内容"),
			SimpleQuestion: "新内容",
			Type:           "single_choice",
			Difficulty:     "easy",
			Options:        []byte("[]"),
			OptionTuples:   []byte("[]"),
			Answers:        []string{"A"},
			Category:       "grammar",
			Labels:         []string{"新标签"},
			Explanation:    "解析",
			Attachments:    []string{},
			TimeLimit:      60,
		}

		question, err := service.Update(ctx, updateDTO)
		assert.Error(t, err)
		assert.Nil(t, question)
		assert.Equal(t, "问题ID不能为空", err.Error())
	})

	t.Run("标题为空", func(t *testing.T) {
		updateDTO := &dto.UpdateQuestionDTO{
			ID:             "1",
			Title:          "",
			Content:        []byte("新内容"),
			SimpleQuestion: "新内容",
			Type:           "single_choice",
			Difficulty:     "easy",
			Options:        []byte("[]"),
			OptionTuples:   []byte("[]"),
			Answers:        []string{"A"},
			Category:       "grammar",
			Labels:         []string{"新标签"},
			Explanation:    "解析",
			Attachments:    []string{},
			TimeLimit:      60,
		}

		question, err := service.Update(ctx, updateDTO)
		assert.Error(t, err)
		assert.Nil(t, question)
		assert.Equal(t, "标题和内容不能为空", err.Error())
	})

	t.Run("问题不存在", func(t *testing.T) {
		mockQuestionRepo.On("Get", ctx, "999").Return(nil, errors.New("问题不存在")).Once()

		updateDTO := &dto.UpdateQuestionDTO{
			ID:             "999",
			Title:          "新标题",
			Content:        []byte("新内容"),
			SimpleQuestion: "新内容",
			Type:           "single_choice",
			Difficulty:     "easy",
			Options:        []byte("[]"),
			OptionTuples:   []byte("[]"),
			Answers:        []string{"A"},
			Category:       "grammar",
			Labels:         []string{"新标签"},
			Explanation:    "解析",
			Attachments:    []string{},
			TimeLimit:      60,
		}

		question, err := service.Update(ctx, updateDTO)
		assert.Error(t, err)
		assert.Nil(t, question)
	})
}

// TestQuestionService_Delete 测试删除问题
func TestQuestionService_Delete(t *testing.T) {
	mockQuestionRepo := new(MockQuestionRepository)
	service := NewQuestionService(mockQuestionRepo)
	ctx := context.Background()

	t.Run("成功删除问题", func(t *testing.T) {
		question := &entity.Question{ID: entity.QuestionID(1), Title: "测试问题"}
		mockQuestionRepo.On("Get", ctx, "1").Return(question, nil).Once()
		mockQuestionRepo.On("Update", ctx, mock.AnythingOfType("*entity.Question")).Return(nil).Once()

		err := service.Delete(ctx, "1")
		assert.NoError(t, err)
	})

	t.Run("问题ID为空", func(t *testing.T) {
		err := service.Delete(ctx, "")
		assert.Error(t, err)
		assert.Equal(t, "问题ID不能为空", err.Error())
	})

	t.Run("问题不存在", func(t *testing.T) {
		mockQuestionRepo.On("Get", ctx, "999").Return(nil, errors.New("问题不存在")).Once()

		err := service.Delete(ctx, "999")
		assert.Error(t, err)
	})
}

// TestQuestionService_Publish 测试发布问题
func TestQuestionService_Publish(t *testing.T) {
	mockQuestionRepo := new(MockQuestionRepository)
	service := NewQuestionService(mockQuestionRepo)
	ctx := context.Background()

	t.Run("成功发布问题", func(t *testing.T) {
		question := &entity.Question{ID: entity.QuestionID(1), Title: "测试问题", Status: "draft"}
		mockQuestionRepo.On("Get", ctx, "1").Return(question, nil).Once()
		mockQuestionRepo.On("Update", ctx, mock.AnythingOfType("*entity.Question")).Return(nil).Once()

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
		mockQuestionRepo.On("Get", ctx, "999").Return(nil, errors.New("问题不存在")).Once()

		question, err := service.Publish(ctx, "999")
		assert.Error(t, err)
		assert.Nil(t, question)
	})

	t.Run("更新失败", func(t *testing.T) {
		question := &entity.Question{ID: entity.QuestionID(1), Title: "测试问题", Status: "draft"}
		mockQuestionRepo.On("Get", ctx, "1").Return(question, nil).Once()
		mockQuestionRepo.On("Update", ctx, mock.AnythingOfType("*entity.Question")).Return(errors.New("更新失败")).Once()

		updatedQuestion, err := service.Publish(ctx, "1")
		assert.Error(t, err)
		assert.Nil(t, updatedQuestion)
	})
}

// TestQuestionService_FindAllTags 测试查询所有标签
func TestQuestionService_FindAllTags(t *testing.T) {
	mockQuestionRepo := new(MockQuestionRepository)
	service := NewQuestionService(mockQuestionRepo)
	ctx := context.Background()

	t.Run("成功获取标签列表", func(t *testing.T) {
		expectedTags := []*entity.QuestionTag{
			{ID: "1", Name: "标签1"},
			{ID: "2", Name: "标签2"},
		}
		mockQuestionRepo.On("FindAllTags", ctx, 10).Return(expectedTags, nil).Once()

		tags, err := service.FindAllTags(ctx, 10)
		assert.NoError(t, err)
		assert.Equal(t, expectedTags, tags)
	})

	t.Run("获取标签失败", func(t *testing.T) {
		mockQuestionRepo.On("FindAllTags", ctx, 10).Return(nil, errors.New("获取标签失败")).Once()

		tags, err := service.FindAllTags(ctx, 10)
		assert.Error(t, err)
		assert.Nil(t, tags)
	})
}

// TestQuestionService_CreateTag 测试创建标签
func TestQuestionService_CreateTag(t *testing.T) {
	mockQuestionRepo := new(MockQuestionRepository)
	service := NewQuestionService(mockQuestionRepo)
	ctx := context.Background()

	t.Run("成功创建标签", func(t *testing.T) {
		expectedTag := &entity.QuestionTag{ID: "1", Name: "新标签"}
		mockQuestionRepo.On("CreateTag", ctx, mock.AnythingOfType("*entity.QuestionTag")).Return(expectedTag, nil).Once()

		tag, err := service.CreateTag(ctx, "新标签")
		assert.NoError(t, err)
		assert.Equal(t, expectedTag, tag)
	})

	t.Run("创建标签失败", func(t *testing.T) {
		mockQuestionRepo.On("CreateTag", ctx, mock.AnythingOfType("*entity.QuestionTag")).Return(nil, errors.New("创建标签失败")).Once()

		tag, err := service.CreateTag(ctx, "新标签")
		assert.Error(t, err)
		assert.Nil(t, tag)
	})
}

// TestQuestionService_UpdateTag 测试更新标签
func TestQuestionService_UpdateTag(t *testing.T) {
	mockQuestionRepo := new(MockQuestionRepository)
	service := NewQuestionService(mockQuestionRepo)
	ctx := context.Background()

	t.Run("成功更新标签", func(t *testing.T) {
		expectedTag := &entity.QuestionTag{ID: "1", Name: "更新后的标签"}
		mockQuestionRepo.On("UpdateTag", ctx, mock.AnythingOfType("*entity.QuestionTag")).Return(expectedTag, nil).Once()

		tag, err := service.UpdateTag(ctx, "1", "更新后的标签")
		assert.NoError(t, err)
		assert.Equal(t, expectedTag, tag)
	})

	t.Run("更新标签失败", func(t *testing.T) {
		mockQuestionRepo.On("UpdateTag", ctx, mock.AnythingOfType("*entity.QuestionTag")).Return(nil, errors.New("更新标签失败")).Once()

		tag, err := service.UpdateTag(ctx, "1", "更新后的标签")
		assert.Error(t, err)
		assert.Nil(t, tag)
	})
}

// TestQuestionService_DeleteTag 测试删除标签
func TestQuestionService_DeleteTag(t *testing.T) {
	mockQuestionRepo := new(MockQuestionRepository)
	service := NewQuestionService(mockQuestionRepo)
	ctx := context.Background()

	t.Run("成功删除标签", func(t *testing.T) {
		mockQuestionRepo.On("DeleteTag", ctx, "1").Return(nil).Once()

		err := service.DeleteTag(ctx, "1")
		assert.NoError(t, err)
	})

	t.Run("删除标签失败", func(t *testing.T) {
		mockQuestionRepo.On("DeleteTag", ctx, "1").Return(errors.New("删除标签失败")).Once()

		err := service.DeleteTag(ctx, "1")
		assert.Error(t, err)
	})
}
