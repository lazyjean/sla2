package websocket

// ChatResponse 是聊天响应结构
type ChatResponse struct {
	Message   string `json:"message"`   // 消息内容
	CreatedAt int64  `json:"createdAt"` // 创建时间戳
	Error     string `json:"error"`     // 错误信息
}
