package service

import (
	"context"
	"errors"
	"testing"

	pb "github.com/lazyjean/sla2/api/proto/v1"
	"github.com/lazyjean/sla2/internal/infrastructure/ai"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// 创建 ChatHistoryRepository 的 Mock
type MockChatHistoryRepo struct {
	mock.Mock
}

func (m *MockChatHistoryRepo) GetHistory(ctx context.Context, userID string, sessionID string) ([]*pb.ChatHistory, error) {
	args := m.Called(ctx, userID, sessionID)
	return args.Get(0).([]*pb.ChatHistory), args.Error(1)
}

func (m *MockChatHistoryRepo) SaveHistory(ctx context.Context, userID string, record *pb.ChatHistory) error {
	args := m.Called(ctx, userID, record)
	return args.Error(0)
}

func (m *MockChatHistoryRepo) CreateSession(ctx context.Context, userID string, title string, description string) (*pb.SessionResponse, error) {
	args := m.Called(ctx, userID, title, description)
	if v, ok := args.Get(0).(*pb.SessionResponse); ok {
		return v, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockChatHistoryRepo) GetSession(ctx context.Context, userID string, sessionID string) (*pb.SessionResponse, error) {
	args := m.Called(ctx, userID, sessionID)
	if v, ok := args.Get(0).(*pb.SessionResponse); ok {
		return v, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockChatHistoryRepo) ListSessions(ctx context.Context, userID string, page uint32, pageSize uint32) ([]*pb.SessionResponse, uint32, error) {
	args := m.Called(ctx, userID, page, pageSize)
	return args.Get(0).([]*pb.SessionResponse), args.Get(1).(uint32), args.Error(2)
}

func (m *MockChatHistoryRepo) DeleteSession(ctx context.Context, userID string, sessionID string) error {
	args := m.Called(ctx, userID, sessionID)
	return args.Error(0)
}

func (m *MockChatHistoryRepo) CountSessionMessages(ctx context.Context, sessionID string) (uint64, error) {
	args := m.Called(ctx, sessionID)
	return args.Get(0).(uint64), args.Error(1)
}

// 创建 DeepSeekService 的 Mock
type MockDeepSeekService struct {
	mock.Mock
}

// 实现 DeepSeekService 接口
func (m *MockDeepSeekService) ChatCompletion(ctx context.Context, messages []ai.ChatMessage) (*ai.ChatCompletionResponse, error) {
	args := m.Called(ctx, messages)
	if v, ok := args.Get(0).(*ai.ChatCompletionResponse); ok {
		return v, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockDeepSeekService) StreamChatCompletion(ctx context.Context, messages []ai.ChatMessage) (<-chan *ai.ChatCompletionChunk, error) {
	args := m.Called(ctx, messages)
	if v, ok := args.Get(0).(<-chan *ai.ChatCompletionChunk); ok {
		return v, args.Error(1)
	}
	return nil, args.Error(1)
}

func TestAIService_SessionManagement(t *testing.T) {
	// 创建模拟对象
	mockRepo := new(MockChatHistoryRepo)
	mockDeepSeek := new(MockDeepSeekService)

	// 创建服务实例
	service := NewAIService(mockDeepSeek, mockRepo)

	// 测试数据
	userID := "123"
	sessionID := "test-session-123"
	title := "测试会话"
	description := "会话描述"

	// 测试 CreateSession
	t.Run("CreateSession", func(t *testing.T) {
		// 设置模拟行为
		expectedResponse := &pb.SessionResponse{
			SessionId:    sessionID,
			Title:        title,
			Description:  description,
			CreatedAt:    timestamppb.Now(),
			UpdatedAt:    timestamppb.Now(),
			MessageCount: 0,
		}
		mockRepo.On("CreateSession", mock.Anything, userID, title, description).Return(expectedResponse, nil).Once()

		// 调用服务
		response, err := service.CreateSession(context.Background(), userID, title, description)

		// 断言结果
		require.NoError(t, err)
		assert.Equal(t, expectedResponse, response)
		mockRepo.AssertExpectations(t)
	})

	// 测试 GetSession
	t.Run("GetSession", func(t *testing.T) {
		// 设置模拟行为
		expectedResponse := &pb.SessionResponse{
			SessionId:    sessionID,
			Title:        title,
			Description:  description,
			CreatedAt:    timestamppb.Now(),
			UpdatedAt:    timestamppb.Now(),
			MessageCount: 5,
		}
		mockRepo.On("GetSession", mock.Anything, userID, sessionID).Return(expectedResponse, nil).Once()

		// 调用服务
		response, err := service.GetSession(context.Background(), userID, sessionID)

		// 断言结果
		require.NoError(t, err)
		assert.Equal(t, expectedResponse, response)
		mockRepo.AssertExpectations(t)
	})

	// 测试 ListSessions
	t.Run("ListSessions", func(t *testing.T) {
		// 设置模拟行为
		expectedSessions := []*pb.SessionResponse{
			{
				SessionId:    sessionID,
				Title:        title,
				Description:  description,
				CreatedAt:    timestamppb.Now(),
				UpdatedAt:    timestamppb.Now(),
				MessageCount: 5,
			},
		}
		expectedTotal := uint32(1)
		mockRepo.On("ListSessions", mock.Anything, userID, uint32(1), uint32(10)).Return(expectedSessions, expectedTotal, nil).Once()

		// 调用服务
		sessions, total, err := service.ListSessions(context.Background(), userID, 1, 10)

		// 断言结果
		require.NoError(t, err)
		assert.Equal(t, expectedSessions, sessions)
		assert.Equal(t, expectedTotal, total)
		mockRepo.AssertExpectations(t)
	})

	// 测试 DeleteSession
	t.Run("DeleteSession", func(t *testing.T) {
		// 设置模拟行为
		mockRepo.On("DeleteSession", mock.Anything, userID, sessionID).Return(nil).Once()

		// 调用服务
		err := service.DeleteSession(context.Background(), userID, sessionID)

		// 断言结果
		require.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	// 测试错误情况
	t.Run("ErrorHandling", func(t *testing.T) {
		// 设置模拟行为 - 返回错误
		expectedError := errors.New("数据库错误")
		mockRepo.On("GetSession", mock.Anything, userID, "non-existent").Return(nil, expectedError).Once()

		// 调用服务
		_, err := service.GetSession(context.Background(), userID, "non-existent")

		// 断言结果
		require.Error(t, err)
		assert.Equal(t, expectedError, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestAIService_ChatWithSession(t *testing.T) {
	// 创建模拟对象
	mockRepo := new(MockChatHistoryRepo)
	mockDeepSeek := new(MockDeepSeekService)

	// 创建服务实例
	service := NewAIService(mockDeepSeek, mockRepo)

	// 测试数据
	userID := "123"
	sessionID := "test-session-456"
	userMessage := "你好，AI助手"

	// 模拟现有的聊天历史
	existingHistory := []*pb.ChatHistory{
		{
			SessionId: sessionID,
			Role:      "user",
			Content:   "前一条用户消息",
			Timestamp: 1000,
		},
		{
			SessionId: sessionID,
			Role:      "assistant",
			Content:   "前一条AI回复",
			Timestamp: 1001,
		},
	}

	// 测试带会话ID的Chat方法
	t.Run("ChatWithSessionID", func(t *testing.T) {
		// 设置模拟行为
		mockRepo.On("GetHistory", mock.Anything, userID, sessionID).Return(existingHistory, nil).Once()

		// 模拟AI响应
		aiResponse := &ai.ChatCompletionResponse{
			ID:      "resp-123",
			Object:  "chat.completion",
			Created: 1000,
			Model:   "deepseek-chat",
			Choices: []struct {
				Message      ai.ChatMessage `json:"message"`
				FinishReason string         `json:"finish_reason"`
			}{
				{
					Message: ai.ChatMessage{
						Role:    "assistant",
						Content: "你好！我是AI助手，有什么可以帮助你的？",
					},
					FinishReason: "stop",
				},
			},
		}

		// 验证传递给DeepSeek的消息包含正确的历史记录
		mockDeepSeek.On("ChatCompletion", mock.Anything, mock.MatchedBy(func(messages []ai.ChatMessage) bool {
			// 验证消息数量是否正确（1个新消息 + 2个历史消息）
			if len(messages) != 3 {
				return false
			}
			// 验证新消息内容
			if messages[0].Role != "user" || messages[0].Content != userMessage {
				return false
			}
			// 验证历史消息
			if messages[1].Role != "user" || messages[1].Content != "前一条用户消息" {
				return false
			}
			if messages[2].Role != "assistant" || messages[2].Content != "前一条AI回复" {
				return false
			}
			return true
		})).Return(aiResponse, nil).Once()

		// 模拟保存用户消息和AI回复
		mockRepo.On("SaveHistory", mock.Anything, userID, mock.MatchedBy(func(msg *pb.ChatHistory) bool {
			return msg.Role == "user" && msg.Content == userMessage && msg.SessionId == sessionID
		})).Return(nil).Once()

		mockRepo.On("SaveHistory", mock.Anything, userID, mock.MatchedBy(func(msg *pb.ChatHistory) bool {
			return msg.Role == "assistant" && msg.Content == "你好！我是AI助手，有什么可以帮助你的？" && msg.SessionId == sessionID
		})).Return(nil).Once()

		// 调用服务
		ctx := &ChatContext{
			SessionID: sessionID,
		}
		response, err := service.Chat(context.Background(), userID, userMessage, ctx)

		// 断言结果
		require.NoError(t, err)
		assert.Equal(t, "你好！我是AI助手，有什么可以帮助你的？", response.Message)
		mockRepo.AssertExpectations(t)
		mockDeepSeek.AssertExpectations(t)
	})
}
