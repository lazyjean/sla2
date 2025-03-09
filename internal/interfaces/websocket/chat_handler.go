package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/lazyjean/sla2/internal/application/service"
	"github.com/lazyjean/sla2/internal/domain/security"
	"github.com/lazyjean/sla2/pkg/logger"
	"go.uber.org/zap"
)

// ChatHandler WebSocket 聊天处理器
type ChatHandler struct {
	aiService    *service.AIService
	tokenService security.TokenService
	upgrader     websocket.Upgrader
	connections  sync.Map // 存储所有活跃连接
}

// NewChatHandler 创建新的聊天处理器
func NewChatHandler(aiService *service.AIService, tokenService security.TokenService) *ChatHandler {
	return &ChatHandler{
		aiService:    aiService,
		tokenService: tokenService,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // 允许所有来源，生产环境应该配置具体的域名
			},
		},
	}
}

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

// HandleWebSocket 处理 WebSocket 连接
func (h *ChatHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger(r.Context())

	// 验证用户身份
	userID, _, err := h.tokenService.ValidateTokenFromContext(r.Context())
	if err != nil {
		log.Error("无效的令牌", zap.Error(err))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// 升级 HTTP 连接为 WebSocket
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error("升级到 WebSocket 失败", zap.Error(err))
		return
	}
	defer conn.Close()

	// 保存连接
	connID := fmt.Sprintf("%d-%d", userID, time.Now().UnixNano())
	h.connections.Store(connID, conn)
	defer h.connections.Delete(connID)

	// 处理连接
	h.handleConnection(r.Context(), conn, fmt.Sprintf("%d", userID))
}

// handleConnection 处理单个 WebSocket 连接
func (h *ChatHandler) handleConnection(ctx context.Context, conn *websocket.Conn, userID string) {
	log := logger.GetLogger(ctx)

	for {
		// 读取消息
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Error("读取 WebSocket 消息失败", zap.Error(err))
			}
			break
		}

		// 解析消息
		var chatMsg ChatMessage
		if err := json.Unmarshal(message, &chatMsg); err != nil {
			log.Error("解析消息失败", zap.Error(err))
			h.sendErrorMessage(conn, "消息格式错误")
			continue
		}

		// 处理消息
		go h.handleChatMessage(ctx, conn, userID, &chatMsg)
	}
}

// handleChatMessage 处理聊天消息
func (h *ChatHandler) handleChatMessage(ctx context.Context, conn *websocket.Conn, userID string, msg *ChatMessage) {
	log := logger.GetLogger(ctx)

	// 发送开始消息
	startMsg := ChatMessage{
		Type:      "stream",
		Action:    "start",
		SessionID: msg.SessionID,
		StreamID:  msg.StreamID,
		Timestamp: time.Now().Unix(),
	}
	if err := conn.WriteJSON(startMsg); err != nil {
		log.Error("发送开始消息失败", zap.Error(err))
		return
	}

	// 调用 AI 服务进行聊天
	responseChan, err := h.aiService.StreamChat(ctx, userID, msg.Message, msg.Context)
	if err != nil {
		log.Error("AI 聊天失败", zap.Error(err))
		h.sendErrorMessage(conn, "AI 服务暂时不可用")
		return
	}

	// 发送流式响应
	for response := range responseChan {
		chatResp := ChatMessage{
			Type:      "stream",
			Action:    "message",
			Message:   response.Message,
			SessionID: msg.SessionID,
			StreamID:  msg.StreamID,
			Role:      "assistant",
			Timestamp: time.Now().Unix(),
			IsFinal:   false,
		}

		if err := conn.WriteJSON(chatResp); err != nil {
			log.Error("发送消息失败", zap.Error(err))
			return
		}
	}

	// 发送结束消息
	endMsg := ChatMessage{
		Type:      "stream",
		Action:    "end",
		SessionID: msg.SessionID,
		StreamID:  msg.StreamID,
		Timestamp: time.Now().Unix(),
		IsFinal:   true,
	}
	if err := conn.WriteJSON(endMsg); err != nil {
		log.Error("发送结束消息失败", zap.Error(err))
	}
}

// sendErrorMessage 发送错误消息
func (h *ChatHandler) sendErrorMessage(conn *websocket.Conn, errMsg string) {
	errorResp := ChatMessage{
		Type:      "stream",
		Action:    "error",
		Error:     errMsg,
		Timestamp: time.Now().Unix(),
		IsFinal:   true,
	}
	conn.WriteJSON(errorResp)
}

// Broadcast 向所有连接广播系统消息
func (h *ChatHandler) Broadcast(message string) {
	sysMsg := ChatMessage{
		Type:      "system",
		Action:    "message",
		Message:   message,
		Role:      "system",
		Timestamp: time.Now().Unix(),
	}

	h.connections.Range(func(key, value interface{}) bool {
		if conn, ok := value.(*websocket.Conn); ok {
			conn.WriteJSON(sysMsg)
		}
		return true
	})
}
