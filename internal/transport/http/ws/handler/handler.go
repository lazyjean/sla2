package handler

import (
	"net/http"

	"github.com/gorilla/websocket"
)

// WebSocketHandler WebSocket 处理器
type WebSocketHandler struct {
	upgrader websocket.Upgrader
}

// NewWebSocketHandler 创建 WebSocket 处理器
func NewWebSocketHandler() *WebSocketHandler {
	return &WebSocketHandler{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // 开发环境允许所有来源
			},
		},
	}
}
