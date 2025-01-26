package entity

import (
	"testing"
	"time"

	"github.com/lazyjean/sla2/domain/errors"
	"github.com/stretchr/testify/assert"
)

func TestNewWord(t *testing.T) {
	tests := []struct {
		name        string
		text        string
		translation string
		wantErr     error
	}{
		{
			name:        "valid word",
			text:        "hello",
			translation: "你好",
			wantErr:     nil,
		},
		{
			name:        "empty text",
			text:        "",
			translation: "测试",
			wantErr:     errors.ErrEmptyWordText,
		},
		{
			name:        "empty translation",
			text:        "test",
			translation: "",
			wantErr:     errors.ErrEmptyTranslation,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			word, err := NewWord(tt.text, tt.translation)
			assert.Equal(t, tt.wantErr, err)
			if err == nil {
				assert.NotEmpty(t, word.CreatedAt)
				assert.NotEmpty(t, word.UpdatedAt)
				assert.Equal(t, tt.text, word.Text)
				assert.Equal(t, tt.translation, word.Translation)
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
				assert.Equal(t, tt.example, word.Examples[0].Text)
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
				assert.Equal(t, tt.tag, word.Tags[0].Name)
				assert.True(t, word.UpdatedAt.After(oldUpdatedAt))
			}
		})
	}
}
