package websocket

import (
	"context"
	"time"

	"github.com/gorilla/websocket"
	pb "github.com/lazyjean/sla2/api/proto/v1"
	"github.com/lazyjean/sla2/internal/application/service"
	"github.com/lazyjean/sla2/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type StreamHandler struct {
	aiService *service.AIService
}

func NewStreamHandler(aiService *service.AIService) *StreamHandler {
	return &StreamHandler{aiService: aiService}
}

func (h *StreamHandler) HandleStream(conn *websocket.Conn, req *pb.StreamChatRequest) {
	log := logger.GetLogger(context.Background())

	// 调用流式处理
	responseChan, err := h.aiService.StreamChat(context.Background(), "0", req.Message, &service.ChatContext{
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
