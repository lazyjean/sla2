package postgres

import (
	"context"
	"testing"

	"github.com/lazyjean/sla2/internal/domain/entity"
	"github.com/lazyjean/sla2/internal/domain/errors"
	"github.com/lazyjean/sla2/internal/domain/valueobject"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWordRepository_Create(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	repo := NewVocabularyRepository(db)
	ctx := context.Background()

	word := &entity.Word{
		Text: "resilient",
		Definitions: []entity.Definition{
			{
				PartOfSpeech: "adj.",
				Meaning:      "有弹性的；能复原的；适应力强的",
				Example:      "Children are generally more resilient than adults.",
				Synonyms:     []string{"flexible", "adaptable"},
				Antonyms:     []string{"fragile", "brittle"},
			},
		},
		Phonetic: "rɪˈzɪliənt",
		Examples: []string{
			"Children are generally more resilient than adults.",
			"The company proved resilient during the economic crisis.",
		},
		Tags:  []string{"adjective", "personality", "advanced"},
		Level: valueobject.WORD_DIFFICULTY_LEVEL_B2,
	}

	// 第一次创建
	err := repo.Create(ctx, word)
	require.NoError(t, err)
	require.NotEmpty(t, word.ID)

	// 验证创建的数据
	saved, err := repo.GetByID(ctx, word.ID)
	require.NoError(t, err)
	require.NotNil(t, saved)

	assert.Equal(t, word.Text, saved.Text)
	assert.Equal(t, word.Definitions, saved.Definitions)
	assert.Equal(t, word.Phonetic, saved.Phonetic)
	assert.Equal(t, word.Examples, saved.Examples)
	assert.Equal(t, word.Tags, saved.Tags)
	assert.Equal(t, word.Level, saved.Level)

	// 尝试创建重复的单词
	duplicate := &entity.Word{
		Text: "resilient",
		Definitions: []entity.Definition{
			{
				PartOfSpeech: "adjective",
				Meaning:      "有弹性的；能快速恢复的；适应力强的",
				Example:      "Children are generally more resilient than adults.",
				Synonyms:     []string{"flexible", "adaptable"},
			},
		},
		Phonetic: "rɪˈzɪliənt",
		Examples: []string{
			"Children are generally more resilient than adults.",
			"The company proved resilient during the economic crisis.",
		},
		Tags:  []string{"adjective", "personality", "advanced"},
		Level: valueobject.WORD_DIFFICULTY_LEVEL_B2,
	}
	err = repo.Create(ctx, duplicate)
	assert.Error(t, err)
	assert.Equal(t, errors.ErrWordAlreadyExists, err)
}

func TestWordRepository_GetByWord(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	repo := NewVocabularyRepository(db)
	ctx := context.Background()

	word := &entity.Word{
		Text: "resilient",
		Definitions: []entity.Definition{
			{
				PartOfSpeech: "adj.",
				Meaning:      "有弹性的；能复原的；适应力强的",
				Example:      "Children are generally more resilient than adults.",
				Synonyms:     []string{"flexible", "adaptable"},
				Antonyms:     []string{"fragile", "brittle"},
			},
		},
		Phonetic: "rɪˈzɪliənt",
		Examples: []string{
			"Children are generally more resilient than adults.",
			"The company proved resilient during the economic crisis.",
		},
		Tags:  []string{"adjective", "personality", "advanced"},
		Level: valueobject.WORD_DIFFICULTY_LEVEL_B2,
	}

	// 创建单词
	err := repo.Create(ctx, word)
	require.NoError(t, err)

	// 通过单词文本查找
	found, err := repo.GetByWord(ctx, word.Text)
	require.NoError(t, err)
	require.NotNil(t, found)

	assert.Equal(t, word.Text, found.Text)
	assert.Equal(t, word.Definitions, found.Definitions)
	assert.Equal(t, word.Phonetic, found.Phonetic)
	assert.Equal(t, word.Examples, found.Examples)
	assert.Equal(t, word.Tags, found.Tags)
	assert.Equal(t, word.Level, found.Level)
}

func TestWordRepository_GetByID(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	repo := NewVocabularyRepository(db)
	ctx := context.Background()

	word := &entity.Word{
		Text: "resilient",
		Definitions: []entity.Definition{
			{
				PartOfSpeech: "adj.",
				Meaning:      "有弹性的；能复原的；适应力强的",
				Example:      "Children are generally more resilient than adults.",
				Synonyms:     []string{"flexible", "adaptable"},
				Antonyms:     []string{"fragile", "brittle"},
			},
		},
		Phonetic: "rɪˈzɪliənt",
		Examples: []string{
			"Children are generally more resilient than adults.",
			"The company proved resilient during the economic crisis.",
		},
		Tags:  []string{"adjective", "personality", "advanced"},
		Level: valueobject.WORD_DIFFICULTY_LEVEL_B2,
	}

	// 创建单词
	err := repo.Create(ctx, word)
	require.NoError(t, err)

	// 通过ID查找
	found, err := repo.GetByID(ctx, word.ID)
	require.NoError(t, err)
	require.NotNil(t, found)

	assert.Equal(t, word.Text, found.Text)
	assert.Equal(t, word.Definitions, found.Definitions)
	assert.Equal(t, word.Phonetic, found.Phonetic)
	assert.Equal(t, word.Examples, found.Examples)
	assert.Equal(t, word.Tags, found.Tags)
	assert.Equal(t, word.Level, found.Level)
}
