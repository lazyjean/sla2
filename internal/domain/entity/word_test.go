package entity

import (
	"testing"
	"time"

	"github.com/lazyjean/sla2/internal/domain/errors"
	"github.com/stretchr/testify/assert"
)

func TestNewWord(t *testing.T) {
	tests := []struct {
		name        string
		userID      uint
		text        string
		translation string
		phonetic    string
		examples    []string
		tags        []string
		wantErr     error
	}{
		{
			name:        "valid word with common English word",
			userID:      1,
			text:        "appreciate",
			translation: "欣赏；感激；领会",
			phonetic:    "əˈpriːʃieɪt",
			examples:    []string{"I really appreciate your help.", "We should appreciate the beauty of nature."},
			tags:        []string{"verb", "advanced"},
			wantErr:     nil,
		},
		{
			name:        "valid word with idiom",
			userID:      1,
			text:        "break the ice",
			translation: "打破僵局；消除陌生感",
			phonetic:    "", // 习语通常不需要音标
			examples:    []string{"The party games helped to break the ice between the guests."},
			tags:        []string{"idiom", "social"},
			wantErr:     nil,
		},
		{
			name:        "valid word with multiple meanings",
			userID:      1,
			text:        "light",
			translation: "光；轻的；浅色的；点燃",
			phonetic:    "laɪt",
			examples: []string{
				"Please turn on the light.",
				"This bag is very light.",
				"She prefers light colors for summer clothing.",
			},
			tags:    []string{"basic", "multiple-meanings", "noun", "adjective"},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			word, err := NewWord(UID(tt.userID), tt.text, tt.phonetic, tt.translation, tt.examples, tt.tags)
			assert.Equal(t, tt.wantErr, err)
			if err == nil {
				assert.NotEmpty(t, word.CreatedAt)
				assert.NotEmpty(t, word.UpdatedAt)
				assert.Equal(t, UID(tt.userID), word.UserID)
				assert.Equal(t, tt.text, word.Text)
				assert.Equal(t, tt.translation, word.Translation)
				assert.Equal(t, tt.examples, word.Examples)
				assert.Equal(t, tt.tags, word.Tags)
			}
		})
	}
}

func TestWord_AddExample(t *testing.T) {
	word := &Word{
		Text:        "test",
		Translation: "测试",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	tests := []struct {
		name    string
		example string
		wantErr error
	}{
		{
			name:    "valid example",
			example: "This is a test.",
			wantErr: nil,
		},
		{
			name:    "empty example",
			example: "",
			wantErr: errors.ErrEmptyExample,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldUpdatedAt := word.UpdatedAt
			err := word.AddExample(tt.example)
			assert.Equal(t, tt.wantErr, err)
			if err == nil {
				assert.Len(t, word.Examples, 1)
				assert.Equal(t, tt.example, word.Examples[0])
				assert.True(t, word.UpdatedAt.After(oldUpdatedAt))
			}
		})
	}
}

func TestWord_AddTag(t *testing.T) {
	word := &Word{
		Text:        "test",
		Translation: "测试",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	tests := []struct {
		name    string
		tag     string
		wantErr error
	}{
		{
			name:    "valid tag",
			tag:     "common",
			wantErr: nil,
		},
		{
			name:    "empty tag",
			tag:     "",
			wantErr: errors.ErrEmptyTag,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldUpdatedAt := word.UpdatedAt
			err := word.AddTag(tt.tag)
			assert.Equal(t, tt.wantErr, err)
			if err == nil {
				assert.Len(t, word.Tags, 1)
				assert.Equal(t, tt.tag, word.Tags[0])
				assert.True(t, word.UpdatedAt.After(oldUpdatedAt))
			}
		})
	}
}
