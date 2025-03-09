package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	pb "github.com/lazyjean/sla2/api/proto/v1"
	"github.com/lazyjean/sla2/internal/domain/repository"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

// ChatHistoryRecordPG 是聊天历史记录的数据库模型
type ChatHistoryRecordPG struct {
	gorm.Model
	UserID    string `gorm:"index:idx_user_session"`
	SessionID string `gorm:"index:idx_user_session"`
	Role      string
	Content   string
}

// ChatSessionPG 是聊天会话的数据库模型
type ChatSessionPG struct {
	gorm.Model
	SessionID   string `gorm:"index;unique"`
	UserID      string `gorm:"index"`
	Title       string
	Description string
}

// TableName 设置表名
func (ChatHistoryRecordPG) TableName() string {
	return "chat_histories"
}

// TableName 设置表名
func (ChatSessionPG) TableName() string {
	return "chat_sessions"
}

// pgChatHistoryRepo 是ChatHistoryRepository的PostgreSQL实现
type pgChatHistoryRepo struct {
	db *gorm.DB
}

// NewChatHistoryRepository 创建一个新的ChatHistoryRepository实例
func NewChatHistoryRepository(db *gorm.DB) repository.ChatHistoryRepository {
	// 自动迁移表结构
	db.AutoMigrate(&ChatHistoryRecordPG{})
	db.AutoMigrate(&ChatSessionPG{})

	return &pgChatHistoryRepo{db: db}
}

// GetHistory 根据用户ID和会话ID获取聊天历史
func (r *pgChatHistoryRepo) GetHistory(ctx context.Context, userID, sessionID string) ([]*pb.ChatHistory, error) {
	var records []ChatHistoryRecordPG
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND session_id = ?", userID, sessionID).
		Order("created_at asc").
		Find(&records).Error

	history := make([]*pb.ChatHistory, len(records))
	for i, r := range records {
		history[i] = &pb.ChatHistory{
			SessionId: r.SessionID,
			Role:      r.Role,
			Content:   r.Content,
			Timestamp: r.CreatedAt.Unix(),
		}
	}
	return history, err
}

// SaveHistory 保存聊天历史记录
func (r *pgChatHistoryRepo) SaveHistory(ctx context.Context, userID string, record *pb.ChatHistory) error {
	// 确保会话存在
	var count int64
	r.db.WithContext(ctx).Model(&ChatSessionPG{}).Where("session_id = ? AND user_id = ?", record.SessionId, userID).Count(&count)
	if count == 0 {
		// 自动创建会话
		session := &ChatSessionPG{
			SessionID:   record.SessionId,
			UserID:      userID,
			Title:       "Chat " + time.Now().Format("2006-01-02 15:04:05"), // 自动生成标题
			Description: "自动创建的会话",
		}
		if err := r.db.WithContext(ctx).Create(session).Error; err != nil {
			return err
		}
	}

	return r.db.WithContext(ctx).Create(&ChatHistoryRecordPG{
		UserID:    userID,
		SessionID: record.SessionId,
		Role:      record.Role,
		Content:   record.Content,
	}).Error
}

// CreateSession 创建新的聊天会话
func (r *pgChatHistoryRepo) CreateSession(ctx context.Context, userID string, title string, description string) (*pb.SessionResponse, error) {
	sessionID := uuid.New().String()
	session := &ChatSessionPG{
		SessionID:   sessionID,
		UserID:      userID,
		Title:       title,
		Description: description,
	}

	if err := r.db.WithContext(ctx).Create(session).Error; err != nil {
		return nil, err
	}

	return &pb.SessionResponse{
		SessionId:    sessionID,
		Title:        title,
		Description:  description,
		CreatedAt:    timestamppb.New(session.CreatedAt),
		UpdatedAt:    timestamppb.New(session.UpdatedAt),
		MessageCount: 0,
	}, nil
}

// GetSession 获取会话详情
func (r *pgChatHistoryRepo) GetSession(ctx context.Context, userID string, sessionID string) (*pb.SessionResponse, error) {
	var session ChatSessionPG
	if err := r.db.WithContext(ctx).Where("session_id = ? AND user_id = ?", sessionID, userID).First(&session).Error; err != nil {
		return nil, err
	}

	// 获取消息数量
	messageCount, err := r.CountSessionMessages(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	return &pb.SessionResponse{
		SessionId:    session.SessionID,
		Title:        session.Title,
		Description:  session.Description,
		CreatedAt:    timestamppb.New(session.CreatedAt),
		UpdatedAt:    timestamppb.New(session.UpdatedAt),
		MessageCount: messageCount,
	}, nil
}

// ListSessions 获取用户的会话列表
func (r *pgChatHistoryRepo) ListSessions(ctx context.Context, userID string, page uint32, pageSize uint32) ([]*pb.SessionResponse, uint32, error) {
	var sessions []ChatSessionPG
	var total int64

	// 计算总数
	if err := r.db.WithContext(ctx).Model(&ChatSessionPG{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).
		Order("updated_at DESC").
		Offset(int(offset)).
		Limit(int(pageSize)).
		Find(&sessions).Error; err != nil {
		return nil, 0, err
	}

	// 转换为响应格式
	responses := make([]*pb.SessionResponse, len(sessions))
	for i, session := range sessions {
		// 获取消息数量
		messageCount, err := r.CountSessionMessages(ctx, session.SessionID)
		if err != nil {
			return nil, 0, err
		}

		responses[i] = &pb.SessionResponse{
			SessionId:    session.SessionID,
			Title:        session.Title,
			Description:  session.Description,
			CreatedAt:    timestamppb.New(session.CreatedAt),
			UpdatedAt:    timestamppb.New(session.UpdatedAt),
			MessageCount: messageCount,
		}
	}

	return responses, uint32(total), nil
}

// DeleteSession 删除会话及其关联的历史记录
func (r *pgChatHistoryRepo) DeleteSession(ctx context.Context, userID string, sessionID string) error {
	tx := r.db.WithContext(ctx).Begin()

	// 删除历史记录
	if err := tx.Where("user_id = ? AND session_id = ?", userID, sessionID).Delete(&ChatHistoryRecordPG{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 删除会话
	if err := tx.Where("user_id = ? AND session_id = ?", userID, sessionID).Delete(&ChatSessionPG{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// CountSessionMessages 统计会话中的消息数量
func (r *pgChatHistoryRepo) CountSessionMessages(ctx context.Context, sessionID string) (uint64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&ChatHistoryRecordPG{}).Where("session_id = ?", sessionID).Count(&count).Error
	return uint64(count), err
}
