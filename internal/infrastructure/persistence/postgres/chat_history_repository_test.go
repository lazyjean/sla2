package postgres

import (
	"context"
	"testing"
	"time"

	pb "github.com/lazyjean/sla2/api/proto/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChatHistoryRepository(t *testing.T) {
	// 设置测试数据库
	db, cleanup := setupTestDB(t)
	defer cleanup()

	// 创建仓储实例
	repo := NewChatHistoryRepository(db)

	// 设置测试用户ID和会话标题
	userID := "123"
	title := "测试会话"
	description := "这是一个测试会话"

	// 测试创建会话
	t.Run("CreateSession", func(t *testing.T) {
		// 创建会话
		session, err := repo.CreateSession(context.Background(), userID, title, description)
		require.NoError(t, err)
		require.NotNil(t, session)

		// 验证会话数据
		assert.NotEmpty(t, session.SessionId)
		assert.Equal(t, title, session.Title)
		assert.Equal(t, description, session.Description)
		assert.NotNil(t, session.CreatedAt)
		assert.NotNil(t, session.UpdatedAt)
		assert.Equal(t, uint64(0), session.MessageCount)

		// 保存会话ID用于后续测试
		sessionID := session.SessionId

		// 测试保存聊天历史
		t.Run("SaveHistory", func(t *testing.T) {
			// 创建用户消息
			userMessage := &pb.ChatHistory{
				SessionId: sessionID,
				Role:      "user",
				Content:   "你好，AI助手",
				Timestamp: time.Now().Unix(),
			}

			// 保存用户消息
			err := repo.SaveHistory(context.Background(), userID, userMessage)
			require.NoError(t, err)

			// 创建AI回复
			aiMessage := &pb.ChatHistory{
				SessionId: sessionID,
				Role:      "assistant",
				Content:   "你好！有什么我可以帮助你的吗？",
				Timestamp: time.Now().Unix(),
			}

			// 保存AI回复
			err = repo.SaveHistory(context.Background(), userID, aiMessage)
			require.NoError(t, err)

			// 测试获取聊天历史
			t.Run("GetHistory", func(t *testing.T) {
				// 获取历史记录
				history, err := repo.GetHistory(context.Background(), userID, sessionID)
				require.NoError(t, err)
				require.Len(t, history, 2)

				// 验证第一条消息（用户消息）
				assert.Equal(t, "user", history[0].Role)
				assert.Equal(t, "你好，AI助手", history[0].Content)
				assert.Equal(t, sessionID, history[0].SessionId)

				// 验证第二条消息（AI回复）
				assert.Equal(t, "assistant", history[1].Role)
				assert.Equal(t, "你好！有什么我可以帮助你的吗？", history[1].Content)
				assert.Equal(t, sessionID, history[1].SessionId)
			})

			// 测试消息数量统计
			t.Run("CountSessionMessages", func(t *testing.T) {
				count, err := repo.CountSessionMessages(context.Background(), sessionID)
				require.NoError(t, err)
				assert.Equal(t, uint64(2), count)
			})

			// 测试获取会话
			t.Run("GetSession", func(t *testing.T) {
				session, err := repo.GetSession(context.Background(), userID, sessionID)
				require.NoError(t, err)
				require.NotNil(t, session)

				assert.Equal(t, sessionID, session.SessionId)
				assert.Equal(t, title, session.Title)
				assert.Equal(t, description, session.Description)
				assert.Equal(t, uint64(2), session.MessageCount)
			})
		})

		// 测试列出会话
		t.Run("ListSessions", func(t *testing.T) {
			// 再创建两个会话，确保有多个会话可以列出
			_, err := repo.CreateSession(context.Background(), userID, "第二个测试会话", "描述2")
			require.NoError(t, err)

			_, err = repo.CreateSession(context.Background(), userID, "第三个测试会话", "描述3")
			require.NoError(t, err)

			// 列出会话
			sessions, total, err := repo.ListSessions(context.Background(), userID, 1, 10)
			require.NoError(t, err)

			// 验证会话列表
			assert.Equal(t, uint32(3), total) // 总共3个会话
			assert.Len(t, sessions, 3)        // 返回3个会话

			// 测试分页
			sessions, total, err = repo.ListSessions(context.Background(), userID, 1, 2)
			require.NoError(t, err)
			assert.Equal(t, uint32(3), total) // 总共3个会话
			assert.Len(t, sessions, 2)        // 第一页返回2个会话

			sessions, total, err = repo.ListSessions(context.Background(), userID, 2, 2)
			require.NoError(t, err)
			assert.Equal(t, uint32(3), total) // 总共3个会话
			assert.Len(t, sessions, 1)        // 第二页返回1个会话
		})

		// 测试删除会话
		t.Run("DeleteSession", func(t *testing.T) {
			// 先获取会话确认存在
			_, err := repo.GetSession(context.Background(), userID, sessionID)
			require.NoError(t, err)

			// 删除会话
			err = repo.DeleteSession(context.Background(), userID, sessionID)
			require.NoError(t, err)

			// 尝试再次获取会话，应该报错
			_, err = repo.GetSession(context.Background(), userID, sessionID)
			require.Error(t, err)

			// 尝试获取聊天历史，应该为空
			history, err := repo.GetHistory(context.Background(), userID, sessionID)
			require.NoError(t, err)
			assert.Len(t, history, 0)
		})
	})
}

// 测试自动创建会话的功能
func TestChatHistoryRepository_AutoCreateSession(t *testing.T) {
	// 设置测试数据库
	db, cleanup := setupTestDB(t)
	defer cleanup()

	// 创建仓储实例
	repo := NewChatHistoryRepository(db)

	// 设置测试数据
	userID := "456"
	sessionID := "auto-session-123"

	// 创建消息（不先创建会话）
	message := &pb.ChatHistory{
		SessionId: sessionID,
		Role:      "user",
		Content:   "直接发送消息，应该自动创建会话",
		Timestamp: time.Now().Unix(),
	}

	// 保存消息
	err := repo.SaveHistory(context.Background(), userID, message)
	require.NoError(t, err)

	// 验证会话已自动创建
	session, err := repo.GetSession(context.Background(), userID, sessionID)
	require.NoError(t, err)
	require.NotNil(t, session)

	// 验证会话数据
	assert.Equal(t, sessionID, session.SessionId)
	assert.Contains(t, session.Title, "Chat") // 标题应该包含"Chat"
	assert.Equal(t, "自动创建的会话", session.Description)
	assert.Equal(t, uint64(1), session.MessageCount)
}
