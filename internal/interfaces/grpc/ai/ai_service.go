package ai

import (
	"context"

	pb "github.com/lazyjean/sla2/api/proto/v1"
	"github.com/lazyjean/sla2/internal/application/service"
	"github.com/lazyjean/sla2/internal/domain/entity"
	"github.com/lazyjean/sla2/internal/domain/security"
	"github.com/lazyjean/sla2/internal/pkg/auth"
	"github.com/lazyjean/sla2/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// AIService AI 聊天服务
type AIService struct {
	pb.UnimplementedAIChatServiceServer
	aiService    *service.AIService
	tokenService security.TokenService
}

// NewAIService 创建 AI 聊天服务
func NewAIService(aiService *service.AIService, tokenService security.TokenService) *AIService {
	return &AIService{
		aiService:    aiService,
		tokenService: tokenService,
	}
}

// CreateSession 创建新的聊天会话
func (s *AIService) CreateSession(ctx context.Context, req *pb.AIChatServiceCreateSessionRequest) (*pb.AIChatServiceCreateSessionResponse, error) {
	log := logger.GetLogger(ctx)

	// Validate user
	userID, _, err := s.tokenService.ValidateTokenFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	// 校验请求参数
	if req.Title == "" {
		return nil, status.Error(codes.InvalidArgument, "title is required")
	}

	// 调用服务创建会话
	sessionID, err := s.aiService.CreateSession(ctx, userID, req.Title)
	if err != nil {
		log.Error("Failed to create chat session", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to create chat session")
	}

	return &pb.AIChatServiceCreateSessionResponse{
		SessionId: uint64(sessionID),
	}, nil
}

// ListSessions 获取用户的会话列表
func (s *AIService) ListSessions(ctx context.Context, req *pb.AIChatServiceListSessionsRequest) (*pb.AIChatServiceListSessionsResponse, error) {
	log := logger.GetLogger(ctx)
	userID, err := auth.GetUserID(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "get user id from context failed")
	}
	// 校验分页参数
	page := req.Page
	if page == 0 {
		page = 1 // 默认第一页
	}
	pageSize := req.PageSize
	if pageSize == 0 {
		pageSize = 10 // 默认每页10条
	}
	if pageSize > 100 {
		pageSize = 100 // 最大每页100条
	}

	// 调用服务获取会话列表
	sessions, total, err := s.aiService.ListSessions(ctx, userID, page, pageSize)
	if err != nil {
		log.Error("Failed to list chat sessions", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to list chat sessions")
	}
	return &pb.AIChatServiceListSessionsResponse{
		Sessions:   toSessionsPB(sessions),
		TotalCount: total,
		Page:       page,
		PageSize:   pageSize,
	}, nil
}

// GetSession 获取会话详情
func (s *AIService) GetSession(ctx context.Context, req *pb.AIChatServiceGetSessionRequest) (*pb.AIChatServiceGetSessionResponse, error) {
	log := logger.GetLogger(ctx)
	// 校验请求参数
	if req.SessionId == 0 {
		return nil, status.Error(codes.InvalidArgument, "session_id is required")
	}

	// 调用服务获取会话详情
	session, err := s.aiService.GetSession(ctx, entity.SessionID(req.SessionId))
	if err != nil {
		log.Error("Failed to get chat session", zap.Error(err), zap.Uint64("session_id", req.SessionId))
		return nil, status.Error(codes.Internal, "failed to get chat session")
	}

	// 转换为 protobuf 消息
	return &pb.AIChatServiceGetSessionResponse{
		Session: &pb.AIChatServiceSessionDetail{
			Session: &pb.AIChatServiceSession{
				Id:        uint64(session.ID),
				Title:     session.Title,
				CreatedAt: timestamppb.New(session.CreatedAt),
				UpdatedAt: timestamppb.New(session.UpdatedAt),
			},
			Messages: toMessagesPB(session.History),
		},
	}, nil
}

// DeleteSession 删除会话
func (s *AIService) DeleteSession(ctx context.Context, req *pb.AIChatServiceDeleteSessionRequest) (*pb.AIChatServiceDeleteSessionResponse, error) {
	log := logger.GetLogger(ctx)

	// 校验请求参数
	if req.SessionId == 0 {
		return nil, status.Error(codes.InvalidArgument, "session_id is required")
	}

	// 调用服务删除会话
	err := s.aiService.DeleteSession(ctx, entity.SessionID(req.SessionId))
	if err != nil {
		log.Error("Failed to delete chat session", zap.Error(err), zap.Uint64("session_id", req.SessionId))
		return nil, status.Error(codes.Internal, "failed to delete chat session")
	}

	return &pb.AIChatServiceDeleteSessionResponse{}, nil
}

// toSessionsPB 将实体会话列表转换为protobuf消息
func toSessionsPB(sessions []entity.AiChatSession) []*pb.AIChatServiceSession {
	result := make([]*pb.AIChatServiceSession, len(sessions))
	for i, session := range sessions {
		result[i] = &pb.AIChatServiceSession{
			Id:        uint64(session.ID),
			Title:     session.Title,
			CreatedAt: timestamppb.New(session.CreatedAt),
			UpdatedAt: timestamppb.New(session.UpdatedAt),
		}
	}
	return result
}

// toMessagesPB 将实体消息列表转换为 protobuf 消息
func toMessagesPB(messages []entity.ChatMessage) []*pb.AIChatServiceMessage {
	result := make([]*pb.AIChatServiceMessage, len(messages))
	for i, msg := range messages {
		result[i] = &pb.AIChatServiceMessage{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}
	return result
}
