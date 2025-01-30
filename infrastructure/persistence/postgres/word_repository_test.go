package postgres

import (
	"context"
	"testing"

	"github.com/lazyjean/sla2/domain/entity"
	"github.com/lazyjean/sla2/domain/errors"
	"github.com/lazyjean/sla2/domain/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWordRepository_Save(t *testing.T) {
	db := setupTestDB(t)
	repo := NewWordRepository(db)
	ctx := context.WithValue(context.Background(), repository.UserIDKey, 1)

	word := &entity.Word{
		Text:        "test",
		Translation: "测试",
		Phonetic:    "test",
		Difficulty:  1,
		Examples:    []string{"This is a test."},
		Tags:        []string{"common"},
		UserID:      1,
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
	assert.Equal(t, word.Examples[0], saved.Examples[0])
	require.NotEmpty(t, saved.Tags)
	assert.Equal(t, word.Tags[0], saved.Tags[0])
}

func TestWordRepository_FindByText(t *testing.T) {
	db := setupTestDB(t)
	repo := NewWordRepository(db)
	ctx := context.Background()

	// 先保存一个单词
	word := &entity.Word{
		Text:        "unique_test",
		Translation: "唯一测试",
		UserID:      1,
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
