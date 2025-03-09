package service

import (
	"context"
	"errors"
	"testing"
	"time"

	pb "github.com/lazyjean/sla2/api/proto/v1"
	"github.com/lazyjean/sla2/internal/infrastructure/ai"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// MockChatHistoryRepo 是 ChatHistoryRepository 的 Mock 实现
type MockChatHistoryRepo struct {
	mock.Mock
}

func (m *MockChatHistoryRepo) GetHistory(ctx context.Context, userID string, sessionID string) ([]*pb.ChatHistory, error) {
	args := m.Called(ctx, userID, sessionID)
	if v, ok := args.Get(0).([]*pb.ChatHistory); ok {
		return v, args.Error(1)
	}
	return nil, args.Error(1)
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
	if v, ok := args.Get(0).([]*pb.SessionResponse); ok {
		return v, args.Get(1).(uint32), args.Error(2)
	}
	return nil, 0, args.Error(2)
}

func (m *MockChatHistoryRepo) DeleteSession(ctx context.Context, userID string, sessionID string) error {
	args := m.Called(ctx, userID, sessionID)
	return args.Error(0)
}

func (m *MockChatHistoryRepo) CountSessionMessages(ctx context.Context, sessionID string) (uint64, error) {
	args := m.Called(ctx, sessionID)
	if v, ok := args.Get(0).(uint64); ok {
		return v, args.Error(1)
	}
	return 0, args.Error(1)
}

// MockDeepSeekService 是 DeepSeekService 的 Mock 实现
type MockDeepSeekService struct {
	mock.Mock
}

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
	mockRepo := &MockChatHistoryRepo{}
	mockDeepSeek := &MockDeepSeekService{}

	// 创建服务实例
	service := NewAIService(mockDeepSeek, mockRepo)

	// 测试数据
	userID := "123"
	sessionID := "test-session-123"
	title := "测试会话"
	description := "会话描述"

	// 测试 CreateSession
	t.Run("CreateSession_Success", func(t *testing.T) {
		expectedResponse := &pb.SessionResponse{
			SessionId:    sessionID,
			Title:        title,
			Description:  description,
			CreatedAt:    timestamppb.Now(),
			UpdatedAt:    timestamppb.Now(),
			MessageCount: 0,
		}
		mockRepo.On("CreateSession", mock.Anything, userID, title, description).Return(expectedResponse, nil).Once()

		response, err := service.CreateSession(context.Background(), userID, title, description)

		require.NoError(t, err)
		assert.Equal(t, expectedResponse, response)
		mockRepo.AssertExpectations(t)
	})

	t.Run("CreateSession_Error", func(t *testing.T) {
		mockRepo.On("CreateSession", mock.Anything, userID, title, description).Return(nil, errors.New("database error")).Once()

		response, err := service.CreateSession(context.Background(), userID, title, description)

		require.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), "database error")
		mockRepo.AssertExpectations(t)
	})

	// 测试 GetSession
	t.Run("GetSession_Success", func(t *testing.T) {
		expectedResponse := &pb.SessionResponse{
			SessionId:    sessionID,
			Title:        title,
			Description:  description,
			CreatedAt:    timestamppb.Now(),
			UpdatedAt:    timestamppb.Now(),
			MessageCount: 5,
		}
		mockRepo.On("GetSession", mock.Anything, userID, sessionID).Return(expectedResponse, nil).Once()

		response, err := service.GetSession(context.Background(), userID, sessionID)

		require.NoError(t, err)
		assert.Equal(t, expectedResponse, response)
		mockRepo.AssertExpectations(t)
	})

	t.Run("GetSession_NotFound", func(t *testing.T) {
		mockRepo.On("GetSession", mock.Anything, userID, "non-existent").Return(nil, errors.New("session not found")).Once()

		response, err := service.GetSession(context.Background(), userID, "non-existent")

		require.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), "session not found")
		mockRepo.AssertExpectations(t)
	})

	// 测试 ListSessions
	t.Run("ListSessions_Success", func(t *testing.T) {
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

		sessions, total, err := service.ListSessions(context.Background(), userID, 1, 10)

		require.NoError(t, err)
		assert.Equal(t, expectedSessions, sessions)
		assert.Equal(t, expectedTotal, total)
		mockRepo.AssertExpectations(t)
	})

	t.Run("ListSessions_Empty", func(t *testing.T) {
		mockRepo.On("ListSessions", mock.Anything, userID, uint32(1), uint32(10)).Return([]*pb.SessionResponse{}, uint32(0), nil).Once()

		sessions, total, err := service.ListSessions(context.Background(), userID, 1, 10)

		require.NoError(t, err)
		assert.Empty(t, sessions)
		assert.Equal(t, uint32(0), total)
		mockRepo.AssertExpectations(t)
	})

	// 测试 DeleteSession
	t.Run("DeleteSession_Success", func(t *testing.T) {
		mockRepo.On("DeleteSession", mock.Anything, userID, sessionID).Return(nil).Once()

		err := service.DeleteSession(context.Background(), userID, sessionID)

		require.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("DeleteSession_Error", func(t *testing.T) {
		mockRepo.On("DeleteSession", mock.Anything, userID, "non-existent").Return(errors.New("session not found")).Once()

		err := service.DeleteSession(context.Background(), userID, "non-existent")

		require.Error(t, err)
		assert.Contains(t, err.Error(), "session not found")
		mockRepo.AssertExpectations(t)
	})
}

func TestAIService_Chat(t *testing.T) {
	mockRepo := &MockChatHistoryRepo{}
	mockDeepSeek := &MockDeepSeekService{}
	service := NewAIService(mockDeepSeek, mockRepo)

	userID := "123"
	sessionID := "test-session-456"
	userMessage := "你好，AI助手"

	t.Run("Chat_WithoutSession", func(t *testing.T) {
		expectedResponse := &ai.ChatCompletionResponse{
			ID:      "resp1",
			Object:  "chat.completion",
			Created: time.Now().Unix(),
			Model:   "deepseek-chat",
			Choices: []struct {
				Message      ai.ChatMessage `json:"message"`
				FinishReason string         `json:"finish_reason"`
			}{
				{
					Message: ai.ChatMessage{
						Role:    "assistant",
						Content: "你好！我是AI助手",
					},
					FinishReason: "stop",
				},
			},
			Usage: struct {
				PromptTokens     int `json:"prompt_tokens"`
				CompletionTokens int `json:"completion_tokens"`
				TotalTokens      int `json:"total_tokens"`
			}{
				PromptTokens:     10,
				CompletionTokens: 10,
				TotalTokens:      20,
			},
		}

		mockDeepSeek.On("ChatCompletion", mock.Anything, mock.MatchedBy(func(messages []ai.ChatMessage) bool {
			return len(messages) == 1 && messages[0].Content == userMessage
		})).Return(expectedResponse, nil).Once()

		response, err := service.Chat(context.Background(), userID, userMessage, nil)

		require.NoError(t, err)
		assert.Equal(t, expectedResponse.Choices[0].Message.Content, response.Message)
		assert.Equal(t, expectedResponse.Created, response.CreatedAt)
		mockDeepSeek.AssertExpectations(t)
	})

	t.Run("Chat_WithSession", func(t *testing.T) {
		history := []*pb.ChatHistory{
			{
				SessionId: sessionID,
				Role:      "user",
				Content:   "历史消息1",
				Timestamp: 1000,
			},
			{
				SessionId: sessionID,
				Role:      "assistant",
				Content:   "AI回复1",
				Timestamp: 1001,
			},
		}

		mockRepo.On("GetHistory", mock.Anything, userID, sessionID).Return(history, nil).Once()
		mockRepo.On("SaveHistory", mock.Anything, userID, mock.MatchedBy(func(msg *pb.ChatHistory) bool {
			return msg.Role == "user" && msg.Content == userMessage
		})).Return(nil).Once()

		expectedResponse := &ai.ChatCompletionResponse{
			ID:      "resp2",
			Object:  "chat.completion",
			Created: time.Now().Unix(),
			Model:   "deepseek-chat",
			Choices: []struct {
				Message      ai.ChatMessage `json:"message"`
				FinishReason string         `json:"finish_reason"`
			}{
				{
					Message: ai.ChatMessage{
						Role:    "assistant",
						Content: "你好！我是AI助手",
					},
					FinishReason: "stop",
				},
			},
			Usage: struct {
				PromptTokens     int `json:"prompt_tokens"`
				CompletionTokens int `json:"completion_tokens"`
				TotalTokens      int `json:"total_tokens"`
			}{
				PromptTokens:     20,
				CompletionTokens: 10,
				TotalTokens:      30,
			},
		}

		mockDeepSeek.On("ChatCompletion", mock.Anything, mock.MatchedBy(func(messages []ai.ChatMessage) bool {
			return len(messages) == 3 // 历史消息 + 新消息
		})).Return(expectedResponse, nil).Once()

		mockRepo.On("SaveHistory", mock.Anything, userID, mock.MatchedBy(func(msg *pb.ChatHistory) bool {
			return msg.Role == "assistant" && msg.Content == expectedResponse.Choices[0].Message.Content
		})).Return(nil).Once()

		response, err := service.Chat(context.Background(), userID, userMessage, &ChatContext{SessionID: sessionID})

		require.NoError(t, err)
		assert.Equal(t, expectedResponse.Choices[0].Message.Content, response.Message)
		assert.Equal(t, expectedResponse.Created, response.CreatedAt)
		mockRepo.AssertExpectations(t)
		mockDeepSeek.AssertExpectations(t)
	})

	t.Run("Chat_GetHistoryError", func(t *testing.T) {
		mockRepo.On("GetHistory", mock.Anything, userID, sessionID).Return(nil, errors.New("database error")).Once()

		response, err := service.Chat(context.Background(), userID, userMessage, &ChatContext{SessionID: sessionID})

		require.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), "获取聊天历史失败")
		mockRepo.AssertExpectations(t)
	})

	t.Run("Chat_AIError", func(t *testing.T) {
		mockDeepSeek.On("ChatCompletion", mock.Anything, mock.Anything).Return(nil, errors.New("AI service error")).Once()

		response, err := service.Chat(context.Background(), userID, userMessage, nil)

		require.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), "AI聊天失败")
		mockDeepSeek.AssertExpectations(t)
	})

	t.Run("Chat_SaveHistoryError", func(t *testing.T) {
		mockRepo.On("GetHistory", mock.Anything, userID, sessionID).Return([]*pb.ChatHistory{}, nil).Once()
		mockRepo.On("SaveHistory", mock.Anything, userID, mock.Anything).Return(errors.New("database error")).Once()

		response, err := service.Chat(context.Background(), userID, userMessage, &ChatContext{SessionID: sessionID})

		require.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), "保存用户消息失败")
		mockRepo.AssertExpectations(t)
	})
}

func TestAIService_StreamChat(t *testing.T) {
	mockRepo := &MockChatHistoryRepo{}
	mockDeepSeek := &MockDeepSeekService{}
	service := NewAIService(mockDeepSeek, mockRepo)

	userID := "123"
	sessionID := "test-session-789"
	userMessage := "你好，AI助手"

	t.Run("StreamChat_Success", func(t *testing.T) {
		// 创建响应通道
		chunkChan := make(chan *ai.ChatCompletionChunk)
		mockRepo.On("GetHistory", mock.Anything, userID, sessionID).Return([]*pb.ChatHistory{}, nil).Once()
		mockRepo.On("SaveHistory", mock.Anything, userID, mock.MatchedBy(func(msg *pb.ChatHistory) bool {
			return msg.Role == "user" && msg.Content == userMessage
		})).Return(nil).Once()
		mockDeepSeek.On("StreamChatCompletion", mock.Anything, mock.MatchedBy(func(messages []ai.ChatMessage) bool {
			return len(messages) == 1 && messages[0].Content == userMessage
		})).Return((<-chan *ai.ChatCompletionChunk)(chunkChan), nil).Once()
		mockRepo.On("SaveHistory", mock.Anything, userID, mock.MatchedBy(func(msg *pb.ChatHistory) bool {
			return msg.Role == "assistant" && msg.Content == "你好"
		})).Return(nil).Once()

		// 启动流式聊天
		responseChan, err := service.StreamChat(context.Background(), userID, userMessage, &ChatContext{SessionID: sessionID})
		require.NoError(t, err)

		// 模拟发送流式响应
		go func() {
			chunkChan <- &ai.ChatCompletionChunk{
				ID:      "chunk1",
				Object:  "chat.completion.chunk",
				Created: time.Now().Unix(),
				Model:   "deepseek-chat",
				Choices: []struct {
					Delta struct {
						Content string `json:"content"`
					} `json:"delta"`
					FinishReason string `json:"finish_reason"`
				}{
					{
						Delta: struct {
							Content string `json:"content"`
						}{Content: "你"},
						FinishReason: "",
					},
				},
			}
			chunkChan <- &ai.ChatCompletionChunk{
				ID:      "chunk2",
				Object:  "chat.completion.chunk",
				Created: time.Now().Unix(),
				Model:   "deepseek-chat",
				Choices: []struct {
					Delta struct {
						Content string `json:"content"`
					} `json:"delta"`
					FinishReason string `json:"finish_reason"`
				}{
					{
						Delta: struct {
							Content string `json:"content"`
						}{Content: "好"},
						FinishReason: "stop",
					},
				},
			}
			close(chunkChan)
		}()

		// 接收并验证响应
		var responses []*ChatResponse
		for response := range responseChan {
			responses = append(responses, response)
		}

		assert.Len(t, responses, 2)
		assert.Equal(t, "你", responses[0].Message)
		assert.Equal(t, "你好", responses[1].Message)
		mockRepo.AssertExpectations(t)
		mockDeepSeek.AssertExpectations(t)
	})

	t.Run("StreamChat_GetHistoryError", func(t *testing.T) {
		mockRepo.On("GetHistory", mock.Anything, userID, sessionID).Return(nil, errors.New("database error")).Once()

		responseChan, err := service.StreamChat(context.Background(), userID, userMessage, &ChatContext{SessionID: sessionID})

		require.Error(t, err)
		assert.Nil(t, responseChan)
		assert.Contains(t, err.Error(), "获取聊天历史失败")
		mockRepo.AssertExpectations(t)
	})

	t.Run("StreamChat_SaveHistoryError", func(t *testing.T) {
		mockRepo.On("GetHistory", mock.Anything, userID, sessionID).Return([]*pb.ChatHistory{}, nil).Once()
		mockRepo.On("SaveHistory", mock.Anything, userID, mock.Anything).Return(errors.New("database error")).Once()

		responseChan, err := service.StreamChat(context.Background(), userID, userMessage, &ChatContext{SessionID: sessionID})

		require.Error(t, err)
		assert.Nil(t, responseChan)
		assert.Contains(t, err.Error(), "保存用户消息失败")
		mockRepo.AssertExpectations(t)
	})

	t.Run("StreamChat_AIError", func(t *testing.T) {
		mockRepo.On("GetHistory", mock.Anything, userID, sessionID).Return([]*pb.ChatHistory{}, nil).Once()
		mockRepo.On("SaveHistory", mock.Anything, userID, mock.MatchedBy(func(msg *pb.ChatHistory) bool {
			return msg.Role == "user" && msg.Content == userMessage
		})).Return(nil).Once()
		mockDeepSeek.On("StreamChatCompletion", mock.Anything, mock.MatchedBy(func(messages []ai.ChatMessage) bool {
			return len(messages) == 1 && messages[0].Content == userMessage
		})).Return(nil, errors.New("AI service error")).Once()

		responseChan, err := service.StreamChat(context.Background(), userID, userMessage, &ChatContext{SessionID: sessionID})

		require.Error(t, err)
		assert.Nil(t, responseChan)
		assert.Contains(t, err.Error(), "流式聊天失败")
		mockRepo.AssertExpectations(t)
		mockDeepSeek.AssertExpectations(t)
	})
}
