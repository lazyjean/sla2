package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
)

type DeepSeekService struct {
	client    *http.Client
	rateLimit *rate.Limiter
	config    *DeepSeekConfig
	logger    *logrus.Logger
}

type DeepSeekConfig struct {
	APIKey      string
	BaseURL     string
	Timeout     int
	MaxRetries  int
	Temperature float64
	MaxTokens   int
}

func NewDeepSeekService(config *DeepSeekConfig, logger *logrus.Logger) *DeepSeekService {
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

func (s *DeepSeekService) ChatCompletion(ctx context.Context, messages []ChatMessage) (*ChatCompletionResponse, error) {
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

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatCompletionRequest struct {
	Model       string        `json:"model"`
	Messages    []ChatMessage `json:"messages"`
	Temperature float64       `json:"temperature,omitempty"`
	MaxTokens   int           `json:"max_tokens,omitempty"`
}

type ChatCompletionResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Message      ChatMessage `json:"message"`
		FinishReason string      `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}
