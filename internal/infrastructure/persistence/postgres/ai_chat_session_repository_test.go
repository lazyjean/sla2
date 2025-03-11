package postgres

import (
	"context"
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lazyjean/sla2/internal/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, *sql.DB) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 db,
		PreferSimpleProtocol: true,
	})

	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	require.NoError(t, err)

	return gormDB, mock, db
}

func TestCreateSession(t *testing.T) {
	// 设置模拟数据库
	gormDB, mock, sqlDB := setupMockDB(t)
	defer sqlDB.Close()

	repo := NewAiChatSessionRepository(gormDB)
	ctx := context.Background()

	// 测试数据
	userID := entity.UID(1)
	title := "测试会话"
	expectedID := entity.SessionID(123)

	// 期望SQL操作
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "ai_chat_sessions"`)).
		WithArgs(sqlmock.AnyArg(), userID, title, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(expectedID))
	mock.ExpectCommit()

	// 执行测试
	resultID, err := repo.CreateSession(ctx, userID, title)

	// 验证结果
	assert.NoError(t, err)
	assert.Equal(t, expectedID, resultID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestListSessions(t *testing.T) {
	// 设置模拟数据库
	gormDB, mock, sqlDB := setupMockDB(t)
	defer sqlDB.Close()

	repo := NewAiChatSessionRepository(gormDB)
	ctx := context.Background()

	// 测试数据
	userID := entity.UID(1)
	page := uint32(1)
	pageSize := uint32(10)

	mockTime := time.Now()
	expectedSessions := []entity.AiChatSession{
		{
			ID:        entity.SessionID(1),
			UserID:    userID,
			Title:     "会话1",
			CreatedAt: mockTime,
			UpdatedAt: mockTime,
		},
		{
			ID:        entity.SessionID(2),
			UserID:    userID,
			Title:     "会话2",
			CreatedAt: mockTime,
			UpdatedAt: mockTime,
		},
	}
	totalCount := 2

	// 期望SQL操作 - 计数查询
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "ai_chat_sessions"`)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(totalCount))

	// 期望SQL操作 - 获取记录
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "ai_chat_sessions"`)).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "title", "history", "created_at", "updated_at"}).
			AddRow(expectedSessions[0].ID, expectedSessions[0].UserID, expectedSessions[0].Title, "[]", mockTime, mockTime).
			AddRow(expectedSessions[1].ID, expectedSessions[1].UserID, expectedSessions[1].Title, "[]", mockTime, mockTime))

	// 执行测试
	sessions, total, err := repo.ListSessions(ctx, userID, page, pageSize)

	// 验证结果
	assert.NoError(t, err)
	assert.Equal(t, len(expectedSessions), len(sessions))
	assert.Equal(t, uint32(totalCount), total)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteSession(t *testing.T) {
	// 设置模拟数据库
	gormDB, mock, sqlDB := setupMockDB(t)
	defer sqlDB.Close()

	repo := NewAiChatSessionRepository(gormDB)
	ctx := context.Background()

	// 测试数据
	sessionID := entity.SessionID(1)

	// 期望SQL操作
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "ai_chat_sessions"`)).
		WithArgs(sessionID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	// 执行测试
	err := repo.DeleteSession(ctx, sessionID)

	// 验证结果
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteSession_Error(t *testing.T) {
	// 设置模拟数据库
	gormDB, mock, sqlDB := setupMockDB(t)
	defer sqlDB.Close()

	repo := NewAiChatSessionRepository(gormDB)
	ctx := context.Background()

	// 测试数据
	sessionID := entity.SessionID(1)

	// 期望SQL操作 - 返回错误
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "ai_chat_sessions"`)).
		WithArgs(sessionID).
		WillReturnError(sql.ErrNoRows)
	mock.ExpectRollback()

	// 执行测试
	err := repo.DeleteSession(ctx, sessionID)

	// 验证结果
	assert.Error(t, err)
	assert.Equal(t, sql.ErrNoRows, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
