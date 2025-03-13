package entity

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// ChatMessage represents a single message in a chat conversation
type ChatMessage struct {
	Role    string `json:"role"`    // "user" or "assistant"
	Content string `json:"content"` // The actual message content
}

// ChatHistory is a slice of ChatMessage objects
type ChatHistory []ChatMessage

// Value implements the driver.Valuer interface for JSONB storage
func (ch ChatHistory) Value() (driver.Value, error) {
	if len(ch) == 0 {
		return json.Marshal([]ChatMessage{})
	}
	return json.Marshal(ch)
}

// Scan implements the sql.Scanner interface for JSONB retrieval
func (ch *ChatHistory) Scan(value interface{}) error {
	if value == nil {
		*ch = ChatHistory{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, ch)
}

// AiChatSession represents a conversation session with an AI
type AiChatSession struct {
	ID        SessionID   `gorm:"primaryKey;autoIncrement"`
	UserID    UID         `gorm:"not null;index"`             // Reference to the user who owns this session
	Title     string      `gorm:"type:varchar(255);not null"` // Title of the conversation
	History   ChatHistory `gorm:"type:jsonb;not null"`        // Chat history stored as JSONB
	CreatedAt time.Time   `gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time   `gorm:"not null;default:CURRENT_TIMESTAMP"`
}

// AddMessage adds a new message to the chat history
func (s *AiChatSession) AddMessage(role, content string) {
	s.History = append(s.History, ChatMessage{
		Role:    role,
		Content: content,
	})
	s.UpdatedAt = time.Now()
}

// GenerateTitle creates a title for the session based on the first user message
// This can be called when a title hasn't been set yet
func (s *AiChatSession) GenerateTitle() {
	if s.Title != "" || len(s.History) == 0 {
		return
	}

	// Find the first user message
	for _, msg := range s.History {
		if msg.Role == "user" && len(msg.Content) > 0 {
			// Use the first 30 characters of the first user message as the title
			if len(msg.Content) > 30 {
				s.Title = msg.Content[:30] + "..."
			} else {
				s.Title = msg.Content
			}
			return
		}
	}
}

// TableName 指定表名
func (AiChatSession) TableName() string {
	return "ai_chat_sessions"
}
