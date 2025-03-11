package postgres

import (
	"context"
	"time"

	"github.com/lazyjean/sla2/internal/domain/entity"
	"github.com/lazyjean/sla2/internal/domain/repository"
	"gorm.io/gorm"
)

// pgChatHistoryRepo 是ChatHistoryRepository的PostgreSQL实现
type aiChatSessionRepository struct {
	db *gorm.DB
}

// NewAiChatSessionRepository 创建一个新的AiChatSessionRepository实例
func NewAiChatSessionRepository(db *gorm.DB) repository.AiChatSessionRepository {
	// 自动迁移表结构
	db.AutoMigrate(&entity.AiChatSession{})

	return &aiChatSessionRepository{
		db: db,
	}
}

// CreateSession 创建新的聊天会话
func (r *aiChatSessionRepository) CreateSession(ctx context.Context, userID entity.UID, title string) (entity.SessionID, error) {
	session := &entity.AiChatSession{
		UserID:    userID,
		Title:     title,
		History:   entity.ChatHistory{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := r.db.WithContext(ctx).Create(session).Error
	if err != nil {
		return 0, err
	}

	return session.ID, nil
}

// ListSessions 获取用户的会话列表
func (r *aiChatSessionRepository) ListSessions(ctx context.Context, userID entity.UID, page uint32, pageSize uint32) ([]entity.AiChatSession, uint32, error) {
	var sessions []entity.AiChatSession
	var total int64

	// 计算总数
	if err := r.db.WithContext(ctx).Model(&entity.AiChatSession{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset(int(offset)).
		Limit(int(pageSize)).
		Find(&sessions).Error; err != nil {
		return nil, 0, err
	}

	return sessions, uint32(total), nil
}

// DeleteSession 删除会话及其关联的历史记录
func (r *aiChatSessionRepository) DeleteSession(ctx context.Context, sessionID entity.SessionID) error {
	tx := r.db.WithContext(ctx).Begin()

	// 删除历史记录
	if err := tx.Where("session_id = ?", sessionID).Delete(&entity.AiChatSession{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// GetSession 通过ID获取会话详情
func (r *aiChatSessionRepository) GetSession(ctx context.Context, sessionID entity.SessionID) (*entity.AiChatSession, error) {
	var session entity.AiChatSession

	if err := r.db.WithContext(ctx).Where("id = ?", sessionID).First(&session).Error; err != nil {
		return nil, err
	}

	return &session, nil
}

// UpdateSession 更新会话
func (r *aiChatSessionRepository) UpdateSession(ctx context.Context, session *entity.AiChatSession) error {
	return r.db.WithContext(ctx).Save(session).Error
}

var _ repository.AiChatSessionRepository = (*aiChatSessionRepository)(nil)
