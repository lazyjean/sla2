package ai

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/lazyjean/sla2/internal/domain/entity"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

type DeepSeekService struct {
	client    *http.Client
	rateLimit *rate.Limiter
	config    *DeepSeekConfig
	logger    *zap.Logger
}

type DeepSeekConfig struct {
	APIKey      string
	BaseURL     string
	Timeout     int
	MaxRetries  int
	Temperature float64
	MaxTokens   int
}

func NewDeepSeekService(config *DeepSeekConfig, logger *zap.Logger) *DeepSeekService {
	client := &http.Client{
		Timeout: time.Duration(config.Timeout) * time.Second,
	}

	// 限流：每秒5个请求
	limiter := rate.NewLimiter(rate.Limit(5), 1)

	return &DeepSeekService{
		client:    client,
		rateLimit: limiter,
		config:    config,
		logger:    logger,
	}
}

func (s *DeepSeekService) ChatCompletion(ctx context.Context, messages []entity.ChatMessage) (*ChatCompletionResponse, error) {
	if err := s.rateLimit.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit exceeded: %w", err)
	}

	reqBody := ChatCompletionRequest{
		Model:       "deepseek-chat",
		Messages:    messages,
		Temperature: s.config.Temperature,
		MaxTokens:   s.config.MaxTokens,
	}

	reqBodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", s.config.BaseURL+"/v1/chat/completions", bytes.NewBuffer(reqBodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.config.APIKey)

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	var response ChatCompletionResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}

// StreamChatCompletion 流式聊天完成
func (s *DeepSeekService) StreamChatCompletion(ctx context.Context, messages []entity.ChatMessage) (<-chan *ChatCompletionChunk, error) {
	if err := s.rateLimit.Wait(ctx); err != nil {
		return nil, fmt.Errorf("超出速率限制: %w", err)
	}

	reqBody := ChatCompletionRequest{
		Model:       "deepseek-chat",
		Messages:    messages,
		Temperature: s.config.Temperature,
		MaxTokens:   s.config.MaxTokens,
		Stream:      true, // 启用流式响应
	}

	reqBodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("请求序列化失败: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", s.config.BaseURL+"/v1/chat/completions", bytes.NewBuffer(reqBodyBytes))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.config.APIKey)

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("API 返回非预期状态码: %d, body: %s", resp.StatusCode, string(body))
	}

	// 创建结果通道
	chunkChan := make(chan *ChatCompletionChunk)

	// 启动 goroutine 处理流式响应
	go func() {
		defer resp.Body.Close()
		defer close(chunkChan)

		reader := bufio.NewReader(resp.Body)
		for {
			// 读取每行数据
			line, err := reader.ReadString('\n')
			if err != nil {
				if err != io.EOF {
					s.logger.Error("读取流式响应失败", zap.Error(err))
				}
				break
			}

			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}

			// 处理 SSE 格式
			if !strings.HasPrefix(line, "data: ") {
				continue
			}

			data := strings.TrimPrefix(line, "data: ")
			if data == "[DONE]" {
				break
			}

			var chunk ChatCompletionChunk
			if err := json.Unmarshal([]byte(data), &chunk); err != nil {
				s.logger.Error("解析流式响应失败", zap.Error(err))
				continue
			}

			select {
			case <-ctx.Done():
				return
			case chunkChan <- &chunk:
				// 成功发送
			}
		}
	}()

	return chunkChan, nil
}

type ChatCompletionRequest struct {
	Model       string               `json:"model"`
	Messages    []entity.ChatMessage `json:"messages"`
	Temperature float64              `json:"temperature,omitempty"`
	MaxTokens   int                  `json:"max_tokens,omitempty"`
	Stream      bool                 `json:"stream,omitempty"`
}

type ChatCompletionResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Message      entity.ChatMessage `json:"message"`
		FinishReason string             `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

// ChatCompletionChunk 流式聊天响应片段
type ChatCompletionChunk struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Delta struct {
			Content string `json:"content"`
		} `json:"delta"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
}
