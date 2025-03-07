package service

import (
	"context"
)

type AIService interface {
	Chat(ctx context.Context, userID string, message string, context *ChatContext) (*ChatResponse, error)
	StreamChat(ctx context.Context, userID string, message string, context *ChatContext) (<-chan *ChatResponse, error)
}

type ChatContext struct {
	History []string
}

type ChatResponse struct {
	Message   string
	CreatedAt int64
}
