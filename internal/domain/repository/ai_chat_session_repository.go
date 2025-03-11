// File: chat_history_repository.go
// This file is temporarily commented out due to missing pb.ChatHistory type
// TODO: Define ChatHistory type in proto files or implement it in the project

package repository

import (
	"context"

	"github.com/lazyjean/sla2/internal/domain/entity"
)

// AiChatSessionRepository 定义了聊天历史记录的仓储接口
type AiChatSessionRepository interface {
	// 会话管理相关方法
	// CreateSession 创建新的聊天会话
	CreateSession(ctx context.Context, userID entity.UID, title string) (entity.SessionID, error)

	// ListSessions 获取用户的会话列表
	ListSessions(ctx context.Context, userID entity.UID, page uint32, pageSize uint32) ([]entity.AiChatSession, uint32, error)

	// DeleteSession 删除会话及其关联的历史记录
	DeleteSession(ctx context.Context, sessionID entity.SessionID) error

	// GetSession 获取会话
	GetSession(ctx context.Context, sessionID entity.SessionID) (*entity.AiChatSession, error)

	// UpdateSession 更新会话
	UpdateSession(ctx context.Context, session *entity.AiChatSession) error
}
