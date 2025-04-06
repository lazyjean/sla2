package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewWord(t *testing.T) {
	text := "hello"
	phonetic := "həˈləʊ"
	definitions := []Definition{
		{
			PartOfSpeech: "int.",
			Meaning:      "你好",
			Example:      "Hello, how are you?",
		},
	}
	examples := []string{"Hello, how are you?"}
	tags := []string{"greeting"}

	word, err := NewWord(text, phonetic, definitions, examples, tags)
	assert.NoError(t, err)
	assert.NotNil(t, word)
	assert.Equal(t, text, word.Text)
	assert.Equal(t, phonetic, word.Phonetic)
	assert.Equal(t, definitions, word.Definitions)
	assert.Equal(t, examples, word.Examples)
	assert.Equal(t, tags, word.Tags)
	assert.Equal(t, 0, word.Difficulty)
	assert.True(t, word.CreatedAt.After(time.Now().Add(-time.Second)))
	assert.True(t, word.UpdatedAt.After(time.Now().Add(-time.Second)))
}

func TestWordValidation(t *testing.T) {
	tests := []struct {
		name    string
		word    *Word
		wantErr bool
	}{
		{
			name: "valid word",
			word: &Word{
				Text: "hello",
				Definitions: []Definition{
					{
						PartOfSpeech: "int.",
						Meaning:      "你好",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "empty text",
			word: &Word{
				Text: "",
				Definitions: []Definition{
					{
						PartOfSpeech: "int.",
						Meaning:      "你好",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "empty definitions",
			word: &Word{
				Text:        "hello",
				Definitions: []Definition{},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.word.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAddExample(t *testing.T) {
	word := &Word{
		Text: "hello",
		Definitions: []Definition{
			{
				PartOfSpeech: "int.",
				Meaning:      "你好",
			},
		},
	}

	// 添加新例句
	err := word.AddExample("Hello, how are you?")
	assert.NoError(t, err)
	assert.Len(t, word.Examples, 1)
	assert.Equal(t, "Hello, how are you?", word.Examples[0])

	// 添加重复例句
	err = word.AddExample("Hello, how are you?")
	assert.NoError(t, err)
	assert.Len(t, word.Examples, 1)

	// 添加空例句
	err = word.AddExample("")
	assert.Error(t, err)
}

func TestAddTag(t *testing.T) {
	word := &Word{
		Text: "hello",
		Definitions: []Definition{
			{
				PartOfSpeech: "int.",
				Meaning:      "你好",
			},
		},
	}

	// 添加新标签
	err := word.AddTag("greeting")
	assert.NoError(t, err)
	assert.Len(t, word.Tags, 1)
	assert.Equal(t, "greeting", word.Tags[0])

	// 添加重复标签
	err = word.AddTag("greeting")
	assert.NoError(t, err)
	assert.Len(t, word.Tags, 1)

	// 添加空标签
	err = word.AddTag("")
	assert.Error(t, err)
}

func TestRemoveTag(t *testing.T) {
	word := &Word{
		Text: "hello",
		Definitions: []Definition{
			{
				PartOfSpeech: "int.",
				Meaning:      "你好",
			},
		},
		Tags: []string{"greeting", "basic"},
	}

	// 移除存在的标签
	word.RemoveTag("greeting")
	assert.Len(t, word.Tags, 1)
	assert.Equal(t, "basic", word.Tags[0])

	// 移除不存在的标签
	word.RemoveTag("nonexistent")
	assert.Len(t, word.Tags, 1)
}

func TestHasTag(t *testing.T) {
	word := &Word{
		Text: "hello",
		Definitions: []Definition{
			{
				PartOfSpeech: "int.",
				Meaning:      "你好",
			},
		},
		Tags: []string{"greeting", "basic"},
	}

	assert.True(t, word.HasTag("greeting"))
	assert.True(t, word.HasTag("basic"))
	assert.False(t, word.HasTag("nonexistent"))
}
