// File: chat_history_repository.go
// This file is temporarily commented out due to missing pb.ChatHistory type
// TODO: Define ChatHistory type in proto files or implement it in the project

package repository

import (
	"context"

	pb "github.com/lazyjean/sla2/api/proto/v1"
)

// ChatHistoryRepository 定义了聊天历史记录的仓储接口
type ChatHistoryRepository interface {
	// GetHistory 根据用户ID和会话ID获取聊天历史
	GetHistory(ctx context.Context, userID string, sessionID string) ([]*pb.ChatHistory, error)

	// SaveHistory 保存聊天历史记录
	SaveHistory(ctx context.Context, userID string, record *pb.ChatHistory) error

	// 会话管理相关方法
	// CreateSession 创建新的聊天会话
	CreateSession(ctx context.Context, userID string, title string, description string) (*pb.SessionResponse, error)

	// GetSession 获取会话详情
	GetSession(ctx context.Context, userID string, sessionID string) (*pb.SessionResponse, error)

	// ListSessions 获取用户的会话列表
	ListSessions(ctx context.Context, userID string, page uint32, pageSize uint32) ([]*pb.SessionResponse, uint32, error)

	// DeleteSession 删除会话及其关联的历史记录
	DeleteSession(ctx context.Context, userID string, sessionID string) error

	// CountSessionMessages 统计会话中的消息数量
	CountSessionMessages(ctx context.Context, sessionID string) (uint64, error)
}
