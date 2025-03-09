package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	pb "github.com/lazyjean/sla2/api/proto/v1"
	"github.com/lazyjean/sla2/internal/application/service"
	"github.com/lazyjean/sla2/internal/domain/entity"
	"github.com/lazyjean/sla2/internal/domain/security"
	"github.com/lazyjean/sla2/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// StopStreamRequest 定义停止流请求结构体
type StopStreamRequest struct {
	StreamId string `json:"stream_id"` // Stream ID to be stopped
}

// AIServiceInterface 定义 AIService 接口，方便测试
type AIServiceInterface interface {
	Chat(ctx context.Context, userID string, message string, chatContext *service.ChatContext) (*service.ChatResponse, error)
	StreamChat(ctx context.Context, userID string, message string, chatContext *service.ChatContext) (<-chan *service.ChatResponse, error)
}

// TokenServiceInterface 定义 TokenService 接口，方便测试
type TokenServiceInterface interface {
	ValidateTokenFromRequest(r *http.Request) (entity.UID, []string, error)
}

// UnifiedChatHandler 统一的 WebSocket 聊天处理器
type UnifiedChatHandler struct {
	aiService     AIServiceInterface
	tokenService  TokenServiceInterface
	upgrader      websocket.Upgrader
	connections   sync.Map // 存储所有活跃连接
	activeStreams sync.Map // 存储活跃的流ID和对应的取消函数
}

// NewUnifiedChatHandler 创建新的统一聊天处理器
func NewUnifiedChatHandler(aiService *service.AIService, tokenService security.TokenService) *UnifiedChatHandler {
	return &UnifiedChatHandler{
		aiService:    aiService,
		tokenService: tokenService,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // 允许所有来源，生产环境应该配置具体的域名
			},
		},
	}
}

// HandleWebSocket 处理 WebSocket 连接
func (h *UnifiedChatHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger(r.Context())

	// 验证用户身份
	userID, _, err := h.tokenService.ValidateTokenFromRequest(r)
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
func (h *UnifiedChatHandler) handleConnection(ctx context.Context, conn *websocket.Conn, userID string) {
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

		// 尝试解析为 protobuf 消息
		var protoReq pb.StreamChatRequest
		if err := json.Unmarshal(message, &protoReq); err == nil && protoReq.StreamId != "" {
			// 处理 protobuf 格式的消息
			go h.handleProtobufMessage(ctx, conn, userID, &protoReq)
			continue
		}

		// 尝试解析为普通聊天消息
		var chatMsg ChatMessage
		if err := json.Unmarshal(message, &chatMsg); err != nil {
			log.Error("解析消息失败", zap.Error(err))
			h.sendErrorMessage(conn, "消息格式错误")
			continue
		}

		// 检查是否是流终止消息
		if chatMsg.Type == "stream" && chatMsg.Action == "stop" && chatMsg.StreamID != "" {
			if h.cancelStream(chatMsg.StreamID) {
				log.Info("流已终止", zap.String("streamId", chatMsg.StreamID))
				// 发送终止确认消息
				endMsg := ChatMessage{
					Type:      "stream",
					Action:    "end",
					SessionID: chatMsg.SessionID,
					StreamID:  chatMsg.StreamID,
					Timestamp: time.Now().Unix(),
					IsFinal:   true,
					Message:   "Stream terminated by client request",
				}
				if err := conn.WriteJSON(endMsg); err != nil {
					log.Error("发送终止确认消息失败", zap.Error(err))
				}
			} else {
				log.Warn("尝试终止不存在的流", zap.String("streamId", chatMsg.StreamID))
				h.sendErrorMessage(conn, "找不到指定的流ID")
			}
			continue
		}

		// 检查是否是protobuf格式的终止消息
		var stopReq StopStreamRequest
		if err := json.Unmarshal(message, &stopReq); err == nil && stopReq.StreamId != "" {
			if h.cancelStream(stopReq.StreamId) {
				log.Info("流已终止(protobuf)", zap.String("streamId", stopReq.StreamId))
				// 发送终止确认消息
				conn.WriteJSON(&pb.ChatResponse{
					StreamId:  stopReq.StreamId,
					IsFinal:   true,
					Message:   "Stream terminated by client request",
					Code:      pb.StatusCode_STATUS_OK,
					CreatedAt: timestamppb.New(time.Now()),
				})
			} else {
				log.Warn("尝试终止不存在的流(protobuf)", zap.String("streamId", stopReq.StreamId))
				conn.WriteJSON(&pb.ChatResponse{
					StreamId:  stopReq.StreamId,
					IsFinal:   true,
					Message:   "Stream not found",
					Code:      pb.StatusCode_STREAM_ERROR,
					CreatedAt: timestamppb.New(time.Now()),
				})
			}
			continue
		}

		// 处理普通聊天消息
		go h.handleChatMessage(ctx, conn, userID, &chatMsg)
	}
}

// handleProtobufMessage 处理 protobuf 格式的消息
func (h *UnifiedChatHandler) handleProtobufMessage(ctx context.Context, conn *websocket.Conn, userID string, req *pb.StreamChatRequest) {
	log := logger.GetLogger(ctx)

	// 调用流式处理
	responseChan, err := h.aiService.StreamChat(ctx, userID, req.Message, &service.ChatContext{
		SessionID: req.Context.GetSessionId(),
		History:   req.Context.GetHistory(),
	})

	if err != nil {
		log.Error("启动流式聊天失败", zap.Error(err))
		conn.WriteJSON(&pb.ChatResponse{
			StreamId:  req.GetStreamId(),
			IsFinal:   true,
			Code:      pb.StatusCode_INTERNAL_ERROR,
			ErrorMsg:  err.Error(),
			CreatedAt: timestamppb.New(time.Now()),
		})
		return
	}

	// 处理流式响应
	for response := range responseChan {
		resp := &pb.ChatResponse{
			Message:   response.Message,
			StreamId:  req.GetStreamId(),
			IsFinal:   false,
			Code:      pb.StatusCode_STATUS_OK,
			CreatedAt: timestamppb.New(time.Unix(response.CreatedAt, 0)),
		}

		if err := conn.WriteJSON(resp); err != nil {
			log.Error("发送消息失败", zap.Error(err))
			return
		}
	}

	// 发送结束消息
	conn.WriteJSON(&pb.ChatResponse{
		StreamId:  req.GetStreamId(),
		IsFinal:   true,
		Code:      pb.StatusCode_STATUS_OK,
		CreatedAt: timestamppb.New(time.Now()),
	})
}

// handleChatMessage 处理普通聊天消息
func (h *UnifiedChatHandler) handleChatMessage(ctx context.Context, conn *websocket.Conn, userID string, msg *ChatMessage) {
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
func (h *UnifiedChatHandler) sendErrorMessage(conn *websocket.Conn, errMsg string) {
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
func (h *UnifiedChatHandler) Broadcast(message string) {
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

// registerStream 注册一个新的流
func (h *UnifiedChatHandler) registerStream(streamID string, cancelFunc context.CancelFunc) {
	h.activeStreams.Store(streamID, cancelFunc)
}

// unregisterStream 注销一个流
func (h *UnifiedChatHandler) unregisterStream(streamID string) {
	h.activeStreams.Delete(streamID)
}

// cancelStream 取消一个流
func (h *UnifiedChatHandler) cancelStream(streamID string) bool {
	if value, ok := h.activeStreams.Load(streamID); ok {
		if cancelFunc, ok := value.(context.CancelFunc); ok {
			cancelFunc()
			h.unregisterStream(streamID)
			return true
		}
	}
	return false
}
