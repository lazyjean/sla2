package service

import (
	"context"
	"fmt"
	"time"

	"github.com/lazyjean/sla2/internal/domain/entity"
	"github.com/lazyjean/sla2/internal/domain/repository"
	"github.com/lazyjean/sla2/internal/infrastructure/ai"
	"github.com/lazyjean/sla2/pkg/logger"
	"go.uber.org/zap"
)

// DeepSeekService 定义了 AI 服务接口
type DeepSeekService interface {
	ChatCompletion(ctx context.Context, messages []entity.ChatMessage) (*ai.ChatCompletionResponse, error)
	StreamChatCompletion(ctx context.Context, messages []entity.ChatMessage) (<-chan *ai.ChatCompletionChunk, error)
}

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
	deepseekService DeepSeekService
	chatSessionRepo repository.AiChatSessionRepository
}

// NewAIService 创建AI助手服务
func NewAIService(deepseekService DeepSeekService, chatSessionRepo repository.AiChatSessionRepository) *AIService {
	return &AIService{
		deepseekService: deepseekService,
		chatSessionRepo: chatSessionRepo,
	}
}

// Chat 单次聊天
func (s *AIService) Chat(ctx context.Context, sessionID entity.SessionID, message string) (*ChatResponse, error) {
	if sessionID == 0 {
		return nil, fmt.Errorf("会话ID不能为0")
	}

	// 如果提供了会话ID，则尝试从存储中获取历史记录
	session, err := s.chatSessionRepo.GetSession(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("获取聊天历史失败: %w", err)
	}

	// 保存用户消息到历史记录
	session.AddMessage("user", message)

	// 调用 AI 服务
	response, err := s.deepseekService.ChatCompletion(ctx, session.History)
	if err != nil {
		return nil, fmt.Errorf("AI聊天失败: %w", err)
	}

	session.History = append(session.History, response.Choices[0].Message)

	return &ChatResponse{
		Message:   response.Choices[0].Message.Content,
		CreatedAt: response.Created,
	}, nil
}

// StreamChat 流式聊天
func (s *AIService) StreamChat(ctx context.Context, sessionID entity.SessionID, message string) (<-chan *ChatResponse, error) {
	log := logger.GetLogger(ctx)
	// 如果提供了会话ID，则尝试从存储中获取历史记录
	session, err := s.chatSessionRepo.GetSession(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("获取聊天历史失败: %w", err)
	}

	// 保存用户消息到历史记录
	session.AddMessage("user", message)

	if err := s.chatSessionRepo.UpdateSession(ctx, session); err != nil {
		return nil, fmt.Errorf("保存用户消息失败: %w", err)
	}

	// 调用流式 API
	chunkChan, err := s.deepseekService.StreamChatCompletion(ctx, session.History)
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
		session.AddMessage("assistant", accumulatedContent)
		if err := s.chatSessionRepo.UpdateSession(ctx, session); err != nil {
			// 在goroutine中记录错误而不是返回
			log.Error("保存AI回复失败", zap.Error(err))
		}
	}()

	return responseChan, nil
}

// CreateSession 创建新的聊天会话
func (s *AIService) CreateSession(ctx context.Context, userID entity.UID, title string) (entity.SessionID, error) {
	return s.chatSessionRepo.CreateSession(ctx, userID, title)
}

// GetSession 获取会话详情
func (s *AIService) GetSession(ctx context.Context, sessionID entity.SessionID) (*entity.AiChatSession, error) {
	return s.chatSessionRepo.GetSession(ctx, sessionID)
}

// ListSessions 获取用户的会话列表
func (s *AIService) ListSessions(ctx context.Context, userID entity.UID, page uint32, pageSize uint32) ([]entity.AiChatSession, uint32, error) {
	return s.chatSessionRepo.ListSessions(ctx, userID, page, pageSize)
}

// DeleteSession 删除会话
func (s *AIService) DeleteSession(ctx context.Context, sessionID entity.SessionID) error {
	return s.chatSessionRepo.DeleteSession(ctx, sessionID)
}
