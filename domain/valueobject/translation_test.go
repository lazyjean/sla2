package valueobject

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTranslation(t *testing.T) {
	tests := []struct {
		name      string
		primary   string
		secondary []string
		wantErr   bool
	}{
		{
			name:      "valid translation",
			primary:   "你好",
			secondary: []string{"哈喽"},
			wantErr:   false,
		},
		{
			name:      "empty primary",
			primary:   "",
			secondary: []string{"test"},
			wantErr:   true,
		},
		{
			name:      "nil secondary",
			primary:   "test",
			secondary: nil,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			trans, err := NewTranslation(tt.primary, tt.secondary)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.primary, trans.Primary)
				if tt.secondary == nil {
					assert.Empty(t, trans.Secondary)
				} else {
					assert.Equal(t, tt.secondary, trans.Secondary)
				}
			}
		})
	}
}

func TestTranslation_Equals(t *testing.T) {
	tests := []struct {
		name     string
		trans1   Translation
		trans2   Translation
		expected bool
	}{
		{
			name: "equal translations",
			trans1: Translation{
				Primary:   "hello",
				Secondary: []string{"hi", "hey"},
			},
			trans2: Translation{
				Primary:   "hello",
				Secondary: []string{"hi", "hey"},
			},
			expected: true,
		},
		{
			name: "different primary",
			trans1: Translation{
				Primary:   "hello",
				Secondary: []string{"hi"},
			},
			trans2: Translation{
				Primary:   "goodbye",
				Secondary: []string{"hi"},
			},
			expected: false,
		},
		{
			name: "different secondary length",
			trans1: Translation{
				Primary:   "hello",
				Secondary: []string{"hi"},
			},
			trans2: Translation{
				Primary:   "hello",
				Secondary: []string{"hi", "hey"},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.trans1.Equals(tt.trans2)
			assert.Equal(t, tt.expected, result)
		})
	}
}
