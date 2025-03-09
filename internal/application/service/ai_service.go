package service

import (
	"context"
	"fmt"
	"time"

	pb "github.com/lazyjean/sla2/api/proto/v1"
	"github.com/lazyjean/sla2/internal/domain/repository"
	"github.com/lazyjean/sla2/internal/infrastructure/ai"
)

// ChatContext 聊天上下文
type ChatContext struct {
	History   []string
	SessionID string
}

// ChatResponse 聊天响应
type ChatResponse struct {
	Message   string
	CreatedAt int64
}

// AIService AI助手服务
type AIService struct {
	deepseekService *ai.DeepSeekService
	chatHistoryRepo repository.ChatHistoryRepository
}

// NewAIService 创建AI助手服务
func NewAIService(deepseekService *ai.DeepSeekService, chatHistoryRepo repository.ChatHistoryRepository) *AIService {
	return &AIService{
		deepseekService: deepseekService,
		chatHistoryRepo: chatHistoryRepo,
	}
}

// Chat 单次聊天
func (s *AIService) Chat(ctx context.Context, userID string, message string, chatContext *ChatContext) (*ChatResponse, error) {
	// 构建聊天消息
	messages := []ai.ChatMessage{
		{
			Role:    "user",
			Content: message,
		},
	}

	sessionID := ""
	if chatContext != nil {
		sessionID = chatContext.SessionID
	}

	// 如果提供了会话ID，则尝试从存储中获取历史记录
	if sessionID != "" {
		history, err := s.chatHistoryRepo.GetHistory(ctx, userID, sessionID)
		if err != nil {
			return nil, fmt.Errorf("获取聊天历史失败: %w", err)
		}

		// 将历史记录添加到消息中
		for _, msg := range history {
			messages = append(messages, ai.ChatMessage{
				Role:    msg.Role,
				Content: msg.Content,
			})
		}
	} else if chatContext != nil && len(chatContext.History) > 0 {
		// 兼容旧的方式：如果没有会话ID但有历史记录，添加到消息中
		for i, historyMsg := range chatContext.History {
			role := "user"
			if i%2 == 1 {
				role = "assistant"
			}
			messages = append(messages, ai.ChatMessage{
				Role:    role,
				Content: historyMsg,
			})
		}
	}

	// 调用底层服务
	response, err := s.deepseekService.ChatCompletion(ctx, messages)
	if err != nil {
		return nil, fmt.Errorf("AI聊天失败: %w", err)
	}

	// 处理响应
	if len(response.Choices) == 0 {
		return nil, fmt.Errorf("没有收到AI响应")
	}

	// 保存用户消息到历史记录
	if sessionID != "" {
		userMsg := &pb.ChatHistory{
			SessionId: sessionID,
			Role:      "user",
			Content:   message,
			Timestamp: time.Now().Unix(),
		}
		if err := s.chatHistoryRepo.SaveHistory(ctx, userID, userMsg); err != nil {
			return nil, fmt.Errorf("保存用户消息失败: %w", err)
		}

		// 保存AI回复到历史记录
		aiMsg := &pb.ChatHistory{
			SessionId: sessionID,
			Role:      "assistant",
			Content:   response.Choices[0].Message.Content,
			Timestamp: time.Now().Unix(),
		}
		if err := s.chatHistoryRepo.SaveHistory(ctx, userID, aiMsg); err != nil {
			return nil, fmt.Errorf("保存AI回复失败: %w", err)
		}
	}

	return &ChatResponse{
		Message:   response.Choices[0].Message.Content,
		CreatedAt: time.Now().Unix(),
	}, nil
}

// StreamChat 流式聊天
func (s *AIService) StreamChat(ctx context.Context, userID string, message string, chatContext *ChatContext) (<-chan *ChatResponse, error) {
	// 构建聊天消息
	messages := []ai.ChatMessage{
		{
			Role:    "user",
			Content: message,
		},
	}

	sessionID := ""
	if chatContext != nil {
		sessionID = chatContext.SessionID
	}

	// 如果提供了会话ID，则尝试从存储中获取历史记录
	if sessionID != "" {
		history, err := s.chatHistoryRepo.GetHistory(ctx, userID, sessionID)
		if err != nil {
			return nil, fmt.Errorf("获取聊天历史失败: %w", err)
		}

		// 将历史记录添加到消息中
		for _, msg := range history {
			messages = append(messages, ai.ChatMessage{
				Role:    msg.Role,
				Content: msg.Content,
			})
		}

		// 保存用户消息到历史记录
		userMsg := &pb.ChatHistory{
			SessionId: sessionID,
			Role:      "user",
			Content:   message,
			Timestamp: time.Now().Unix(),
		}
		if err := s.chatHistoryRepo.SaveHistory(ctx, userID, userMsg); err != nil {
			return nil, fmt.Errorf("保存用户消息失败: %w", err)
		}
	} else if chatContext != nil && len(chatContext.History) > 0 {
		// 兼容旧的方式：如果没有会话ID但有历史记录，添加到消息中
		for i, historyMsg := range chatContext.History {
			role := "user"
			if i%2 == 1 {
				role = "assistant"
			}
			messages = append(messages, ai.ChatMessage{
				Role:    role,
				Content: historyMsg,
			})
		}
	}

	// 调用流式 API
	chunkChan, err := s.deepseekService.StreamChatCompletion(ctx, messages)
	if err != nil {
		return nil, fmt.Errorf("流式聊天失败: %w", err)
	}

	// 创建一个响应通道
	responseChan := make(chan *ChatResponse)

	// 启动 goroutine 处理流式响应
	go func() {
		defer close(responseChan)

		// 累积的响应内容
		var accumulatedContent string

		for chunk := range chunkChan {
			if len(chunk.Choices) == 0 {
				continue
			}

			// 获取增量内容
			deltaContent := chunk.Choices[0].Delta.Content
			if deltaContent == "" {
				continue
			}

			// 累积内容
			accumulatedContent += deltaContent

			// 发送当前累积的内容
			responseChan <- &ChatResponse{
				Message:   accumulatedContent,
				CreatedAt: time.Now().Unix(),
			}
		}

		// 在流式回复完成后，保存完整的AI回复到历史记录
		if sessionID != "" {
			aiMsg := &pb.ChatHistory{
				SessionId: sessionID,
				Role:      "assistant",
				Content:   accumulatedContent,
				Timestamp: time.Now().Unix(),
			}
			_ = s.chatHistoryRepo.SaveHistory(ctx, userID, aiMsg)
		}
	}()

	return responseChan, nil
}

// CreateSession 创建新的聊天会话
func (s *AIService) CreateSession(ctx context.Context, userID string, title string, description string) (*pb.SessionResponse, error) {
	return s.chatHistoryRepo.CreateSession(ctx, userID, title, description)
}

// GetSession 获取会话详情
func (s *AIService) GetSession(ctx context.Context, userID string, sessionID string) (*pb.SessionResponse, error) {
	return s.chatHistoryRepo.GetSession(ctx, userID, sessionID)
}

// ListSessions 获取用户的会话列表
func (s *AIService) ListSessions(ctx context.Context, userID string, page uint32, pageSize uint32) ([]*pb.SessionResponse, uint32, error) {
	return s.chatHistoryRepo.ListSessions(ctx, userID, page, pageSize)
}

// DeleteSession 删除会话
func (s *AIService) DeleteSession(ctx context.Context, userID string, sessionID string) error {
	return s.chatHistoryRepo.DeleteSession(ctx, userID, sessionID)
}
