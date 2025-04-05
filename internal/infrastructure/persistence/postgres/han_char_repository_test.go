package postgres

import (
	"context"
	"testing"

	"github.com/lazyjean/sla2/internal/domain/entity"
	"github.com/lazyjean/sla2/internal/domain/valueobject"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHanCharRepository(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewHanCharRepository(db)

	t.Run("Create and Get", func(t *testing.T) {
		// 创建测试数据
		hanChar := entity.NewHanChar("测试", "cè shì", valueobject.WORD_DIFFICULTY_LEVEL_A1)
		hanChar.Tags = []string{"test", "example"}
		hanChar.Categories = []string{"test", "example"}
		hanChar.Examples = []string{"这是一个测试", "这是一个例子"}

		// 创建汉字
		id, err := repo.Create(context.Background(), hanChar)
		require.NoError(t, err)
		assert.NotZero(t, id)

		// 获取汉字
		got, err := repo.GetByID(context.Background(), id)
		require.NoError(t, err)
		assert.Equal(t, hanChar.Character, got.Character)
		assert.Equal(t, hanChar.Pinyin, got.Pinyin)
		assert.Equal(t, hanChar.Level, got.Level)
		assert.Equal(t, hanChar.Tags, got.Tags)
		assert.Equal(t, hanChar.Categories, got.Categories)
		assert.Equal(t, hanChar.Examples, got.Examples)
	})

	t.Run("GetByCharacter", func(t *testing.T) {
		// 创建测试数据
		hanChar := entity.NewHanChar("测试2", "cè shì 2", valueobject.WORD_DIFFICULTY_LEVEL_A1)
		hanChar.Tags = []string{"test", "example"}
		hanChar.Categories = []string{"test", "example"}
		hanChar.Examples = []string{"这是一个测试", "这是一个例子"}

		// 创建汉字
		_, err := repo.Create(context.Background(), hanChar)
		require.NoError(t, err)

		// 获取汉字
		got, err := repo.GetByCharacter(context.Background(), hanChar.Character)
		require.NoError(t, err)
		assert.Equal(t, hanChar.Character, got.Character)
		assert.Equal(t, hanChar.Pinyin, got.Pinyin)
		assert.Equal(t, hanChar.Level, got.Level)
		assert.Equal(t, hanChar.Tags, got.Tags)
		assert.Equal(t, hanChar.Categories, got.Categories)
		assert.Equal(t, hanChar.Examples, got.Examples)
	})

	t.Run("List", func(t *testing.T) {
		// 创建测试数据
		hanChar := entity.NewHanChar("测试3", "cè shì 3", valueobject.WORD_DIFFICULTY_LEVEL_A1)
		hanChar.Tags = []string{"test", "example"}
		hanChar.Categories = []string{"test", "example"}
		hanChar.Examples = []string{"这是一个测试", "这是一个例子"}

		// 创建汉字
		_, err := repo.Create(context.Background(), hanChar)
		require.NoError(t, err)

		// 获取列表
		hanChars, total, err := repo.List(context.Background(), 0, 10, nil)
		require.NoError(t, err)
		assert.NotZero(t, total)
		assert.NotEmpty(t, hanChars)
	})

	t.Run("Search", func(t *testing.T) {
		// 创建测试数据
		hanChar := entity.NewHanChar("测试4", "cè shì 4", valueobject.WORD_DIFFICULTY_LEVEL_A1)
		hanChar.Tags = []string{"test", "example"}
		hanChar.Categories = []string{"test", "example"}
		hanChar.Examples = []string{"这是一个测试", "这是一个例子"}

		// 创建汉字
		_, err := repo.Create(context.Background(), hanChar)
		require.NoError(t, err)

		// 搜索汉字
		hanChars, total, err := repo.Search(context.Background(), "测试", 0, 10, nil)
		require.NoError(t, err)
		assert.NotZero(t, total)
		assert.NotEmpty(t, hanChars)
	})

	t.Run("Update", func(t *testing.T) {
		// 创建测试数据
		hanChar := entity.NewHanChar("测试5", "cè shì 5", valueobject.WORD_DIFFICULTY_LEVEL_A1)
		hanChar.Tags = []string{"test", "example"}
		hanChar.Categories = []string{"test", "example"}
		hanChar.Examples = []string{"这是一个测试", "这是一个例子"}

		// 创建汉字
		id, err := repo.Create(context.Background(), hanChar)
		require.NoError(t, err)

		// 更新汉字
		hanChar.Pinyin = "cè shì 6"
		hanChar.Tags = []string{"test", "example", "updated"}
		hanChar.Categories = []string{"test", "example", "updated"}
		hanChar.Examples = []string{"这是一个测试", "这是一个例子", "这是一个更新"}
		err = repo.Update(context.Background(), hanChar)
		require.NoError(t, err)

		// 获取更新后的汉字
		got, err := repo.GetByID(context.Background(), id)
		require.NoError(t, err)
		assert.Equal(t, hanChar.Pinyin, got.Pinyin)
		assert.Equal(t, hanChar.Tags, got.Tags)
		assert.Equal(t, hanChar.Categories, got.Categories)
		assert.Equal(t, hanChar.Examples, got.Examples)
	})

	t.Run("Delete", func(t *testing.T) {
		// 创建测试数据
		hanChar := entity.NewHanChar("测试7", "cè shì 7", valueobject.WORD_DIFFICULTY_LEVEL_A1)
		hanChar.Tags = []string{"test", "example"}
		hanChar.Categories = []string{"test", "example"}
		hanChar.Examples = []string{"这是一个测试", "这是一个例子"}

		// 创建汉字
		id, err := repo.Create(context.Background(), hanChar)
		require.NoError(t, err)

		// 删除汉字
		err = repo.Delete(context.Background(), id)
		require.NoError(t, err)

		// 验证汉字已被删除
		got, err := repo.GetByID(context.Background(), id)
		require.Error(t, err)
		assert.Nil(t, got)
	})
}
