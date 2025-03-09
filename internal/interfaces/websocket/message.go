package websocket

import "github.com/lazyjean/sla2/internal/application/service"

// ChatMessage WebSocket 消息结构
type ChatMessage struct {
	Type      string               `json:"type"`      // 消息类型：chat, error, system, stream
	Action    string               `json:"action"`    // 动作类型：start, message, end, error
	Message   string               `json:"message"`   // 消息内容
	SessionID string               `json:"sessionId"` // 会话ID
	StreamID  string               `json:"streamId"`  // 流ID
	IsFinal   bool                 `json:"isFinal"`   // 是否是最后一条消息
	Error     string               `json:"error"`     // 错误信息
	Context   *service.ChatContext `json:"context"`   // 聊天上下文
	Role      string               `json:"role"`      // 消息角色：user, assistant, system
	Timestamp int64                `json:"timestamp"` // 消息时间戳
}
