package postgres

import (
	"context"
	"testing"

	"github.com/lazyjean/sla2/internal/domain/entity"
	"github.com/lazyjean/sla2/internal/domain/errors"
	"github.com/lazyjean/sla2/internal/domain/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWordRepository_Save(t *testing.T) {
	db := setupTestDB(t)
	repo := NewWordRepository(db)
	ctx := context.WithValue(context.Background(), repository.UserIDKey, 1)

	word := &entity.Word{
		Text:        "resilient",
		Translation: "有弹性的；能快速恢复的；适应力强的",
		Phonetic:    "rɪˈzɪliənt",
		Examples: []string{
			"Children are generally more resilient than adults.",
			"The company proved resilient during the economic crisis.",
		},
		Tags:   []string{"adjective", "personality", "advanced"},
		UserID: 1,
	}

	// 第一次保存
	err := repo.Save(ctx, word)
	require.NoError(t, err)
	require.NotEmpty(t, word.ID)

	// 验证保存的数据
	saved, err := repo.FindByID(ctx, word.ID)
	require.NoError(t, err)
	require.NotNil(t, saved)

	assert.Equal(t, word.Text, saved.Text)
	assert.Equal(t, word.Translation, saved.Translation)
	assert.Equal(t, word.Examples, saved.Examples)
	assert.Equal(t, word.Tags, saved.Tags)

	// 尝试再次保存相同的单词，但有不同的翻译和例句
	duplicateWord := &entity.Word{
		Text:        "resilient",
		Translation: "坚韧的；有复原力的", // 不同的翻译
		Phonetic:    "rɪˈzɪliənt",
		Examples: []string{
			"She is remarkably resilient in the face of adversity.",
			"A resilient economy can withstand external shocks.",
		},
		Tags:   []string{"adjective", "character", "psychology"},
		UserID: 1,
	}

	// 保存重复单词应该成功，但不会创建新记录
	err = repo.Save(ctx, duplicateWord)
	require.NoError(t, err)

	// 验证数据库中仍然只有一条记录
	var count int64
	db.Model(&entity.Word{}).Count(&count)
	assert.Equal(t, int64(1), count)

	// 验证原始数据保持不变
	saved, err = repo.FindByID(ctx, word.ID)
	require.NoError(t, err)
	assert.Equal(t, "有弹性的；能快速恢复的；适应力强的", saved.Translation)
	assert.Equal(t, []string{"adjective", "personality", "advanced"}, saved.Tags)
	assert.Equal(t, []string{
		"Children are generally more resilient than adults.",
		"The company proved resilient during the economic crisis.",
	}, saved.Examples)
}

func TestWordRepository_FindByText(t *testing.T) {
	db := setupTestDB(t)
	repo := NewWordRepository(db)
	ctx := context.Background()

	// 先保存一个单词
	word := &entity.Word{
		Text:        "resilient",
		Translation: "有弹性的；能快速恢复的；适应力强的",
		Phonetic:    "rɪˈzɪliənt",
		Difficulty:  4,
		Examples: []string{
			"Children are generally more resilient than adults.",
			"The company proved resilient during the economic crisis.",
		},
		Tags:   []string{"adjective", "personality"},
		UserID: 1,
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
			text:    "resilient",
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
