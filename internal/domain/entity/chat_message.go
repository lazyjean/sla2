package entity

// ChatMessage 表示聊天消息
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
