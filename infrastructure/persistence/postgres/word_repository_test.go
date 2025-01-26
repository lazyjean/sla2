package postgres

import (
	"context"
	"io"
	"log"
	"testing"
	"time"

	"github.com/lazyjean/sla2/domain/entity"
	"github.com/lazyjean/sla2/domain/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupTestDB(t *testing.T) *gorm.DB {
	dsn := "postgres://sla:sla1234@localhost:5432/sla2_test?sslmode=disable"

	// 配置测试专用的日志设置
	logConfig := logger.Config{
		SlowThreshold:             time.Second,
		LogLevel:                  logger.Silent, // 测试时禁用日志
		IgnoreRecordNotFoundError: true,          // 忽略记录未找到的错误
		Colorful:                  false,
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.New(
			log.New(io.Discard, "", 0), // 丢弃所有日志输出
			logConfig,
		),
	})
	require.NoError(t, err)

	// 删除现有表
	db.Exec("DROP TABLE IF EXISTS word_tags CASCADE")
	db.Exec("DROP TABLE IF EXISTS tags CASCADE")
	db.Exec("DROP TABLE IF EXISTS examples CASCADE")
	db.Exec("DROP TABLE IF EXISTS words CASCADE")

	// 创建表
	err = db.AutoMigrate(
		&entity.Word{},
		&entity.Example{},
		&entity.Tag{},
	)
	require.NoError(t, err)

	return db
}

func TestWordRepository_Save(t *testing.T) {
	db := setupTestDB(t)
	repo := NewWordRepository(db)
	ctx := context.Background()

	word := &entity.Word{
		Text:        "test",
		Translation: "测试",
		Phonetic:    "test",
		Difficulty:  1,
		Examples: []entity.Example{
			{Text: "This is a test."},
		},
		Tags: []entity.Tag{
			{Name: "common"},
		},
	}

	err := repo.Save(ctx, word)
	require.NoError(t, err)
	require.NotEmpty(t, word.ID)

	// 验证保存的数据
	saved, err := repo.FindByID(ctx, word.ID)
	require.NoError(t, err)
	require.NotNil(t, saved)

	assert.Equal(t, word.Text, saved.Text)
	assert.Equal(t, word.Translation, saved.Translation)
	require.NotEmpty(t, saved.Examples)
	assert.Equal(t, word.Examples[0].Text, saved.Examples[0].Text)
	require.NotEmpty(t, saved.Tags)
	assert.Equal(t, word.Tags[0].Name, saved.Tags[0].Name)
}

func TestWordRepository_FindByText(t *testing.T) {
	db := setupTestDB(t)
	repo := NewWordRepository(db)
	ctx := context.Background()

	// 先保存一个单词
	word := &entity.Word{
		Text:        "unique_test",
		Translation: "唯一测试",
	}
	err := repo.Save(ctx, word)
	require.NoError(t, err)

	tests := []struct {
		name    string
		text    string
		wantErr error
	}{
		{
			name:    "existing word",
			text:    "unique_test",
			wantErr: nil,
		},
		{
			name:    "non-existing word",
			text:    "nonexistent",
			wantErr: errors.ErrWordNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			found, err := repo.FindByText(ctx, tt.text)
			assert.Equal(t, tt.wantErr, err)
			if err == nil {
				assert.Equal(t, tt.text, found.Text)
			}
		})
	}
}
