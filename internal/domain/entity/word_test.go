package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewWord(t *testing.T) {
	tests := []struct {
		name        string
		userID      uint
		text        string
		definitions []Definition
		phonetic    string
		examples    []string
		tags        []string
		wantErr     error
	}{
		{
			name:   "valid word with common English word",
			userID: 1,
			text:   "appreciate",
			definitions: []Definition{
				{
					PartOfSpeech: "verb",
					Meaning:      "欣赏；感激；领会",
					Example:      "I really appreciate your help.",
					Synonyms:     []string{"value", "cherish"},
				},
			},
			phonetic: "əˈpriːʃieɪt",
			examples: []string{"I really appreciate your help.", "We should appreciate the beauty of nature."},
			tags:     []string{"verb", "advanced"},
			wantErr:  nil,
		},
		{
			name:   "valid word with idiom",
			userID: 1,
			text:   "break the ice",
			definitions: []Definition{
				{
					PartOfSpeech: "idiom",
					Meaning:      "打破僵局；消除陌生感",
					Example:      "The party games helped to break the ice between the guests.",
					Synonyms:     []string{"start a conversation", "make friends"},
				},
			},
			phonetic: "", // 习语通常不需要音标
			examples: []string{"The party games helped to break the ice between the guests."},
			tags:     []string{"idiom", "social"},
			wantErr:  nil,
		},
		{
			name:   "valid word with multiple meanings",
			userID: 1,
			text:   "light",
			definitions: []Definition{
				{
					PartOfSpeech: "noun",
					Meaning:      "光",
					Example:      "Please turn on the light.",
					Synonyms:     []string{"illumination", "brightness"},
				},
				{
					PartOfSpeech: "adjective",
					Meaning:      "轻的；浅色的",
					Example:      "This bag is very light.",
					Synonyms:     []string{"weightless", "pale"},
				},
				{
					PartOfSpeech: "verb",
					Meaning:      "点燃",
					Example:      "He lit a cigarette.",
					Synonyms:     []string{"ignite", "kindle"},
				},
			},
			phonetic: "laɪt",
			examples: []string{
				"Please turn on the light.",
				"This bag is very light.",
				"He lit a cigarette.",
			},
			tags:    []string{"noun", "adjective", "verb"},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			word, err := NewWord(UID(tt.userID), tt.text, tt.phonetic, tt.definitions, tt.examples, tt.tags)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, UID(tt.userID), word.UserID)
			assert.Equal(t, tt.text, word.Text)
			assert.Equal(t, tt.definitions, word.Definitions)
			assert.Equal(t, tt.phonetic, word.Phonetic)
			assert.Equal(t, tt.examples, word.Examples)
			assert.Equal(t, tt.tags, word.Tags)
		})
	}
}

func TestWord_AddExample(t *testing.T) {
	word := &Word{
		UserID: 1,
		Text:   "test",
		Definitions: []Definition{
			{
				PartOfSpeech: "noun",
				Meaning:      "测试",
				Example:      "This is a test.",
			},
		},
		Examples: []string{"This is a test."},
	}

	// 添加新的例句
	word.AddExample("This is another test.")
	assert.Equal(t, []string{"This is a test.", "This is another test."}, word.Examples)

	// 添加重复的例句
	word.AddExample("This is a test.")
	assert.Equal(t, []string{"This is a test.", "This is another test."}, word.Examples)
}

func TestWord_AddTag(t *testing.T) {
	word := &Word{
		UserID: 1,
		Text:   "test",
		Definitions: []Definition{
			{
				PartOfSpeech: "noun",
				Meaning:      "测试",
				Example:      "This is a test.",
			},
		},
		Tags: []string{"test"},
	}

	// 添加新的标签
	word.AddTag("example")
	assert.Equal(t, []string{"test", "example"}, word.Tags)

	// 添加重复的标签
	word.AddTag("test")
	assert.Equal(t, []string{"test", "example"}, word.Tags)
}
