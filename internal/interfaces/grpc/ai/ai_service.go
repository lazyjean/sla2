package ai

import (
	"context"
	"fmt"

	pb "github.com/lazyjean/sla2/api/proto/v1"
	"github.com/lazyjean/sla2/internal/application/service"
	"github.com/lazyjean/sla2/internal/domain/security"
	"github.com/lazyjean/sla2/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// AIService AI 聊天服务
// @title AI Chat Service
// @description AI 聊天服务，提供会话管理功能
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
// @Summary 创建聊天会话
// @Description 创建一个新的聊天会话
// @Tags AI Chat
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param body body pb.CreateSessionRequest true "会话信息"
// @Success 200 {object} pb.SessionResponse
// @Failure 400 {object} status.Status "Invalid request"
// @Failure 401 {object} status.Status "Unauthorized"
// @Failure 500 {object} status.Status "Internal server error"
// @Router /v1/chat/sessions [post]
func (s *AIService) CreateSession(ctx context.Context, req *pb.CreateSessionRequest) (*pb.SessionResponse, error) {
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
	response, err := s.aiService.CreateSession(ctx, fmt.Sprintf("%d", userID), req.Title, req.Description)
	if err != nil {
		log.Error("Failed to create chat session", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to create chat session")
	}

	return response, nil
}

// ListSessions 获取用户的会话列表
// @Summary 获取会话列表
// @Description 获取用户的所有聊天会话列表
// @Tags AI Chat
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param page query int false "页码，从1开始" default(1)
// @Param page_size query int false "每页大小" default(10)
// @Success 200 {object} pb.ListSessionsResponse
// @Failure 401 {object} status.Status "Unauthorized"
// @Failure 500 {object} status.Status "Internal server error"
// @Router /v1/chat/sessions [get]
func (s *AIService) ListSessions(ctx context.Context, req *pb.ListSessionsRequest) (*pb.ListSessionsResponse, error) {
	log := logger.GetLogger(ctx)

	// Validate user
	userID, _, err := s.tokenService.ValidateTokenFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid token")
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
	sessions, total, err := s.aiService.ListSessions(ctx, fmt.Sprintf("%d", userID), page, pageSize)
	if err != nil {
		log.Error("Failed to list chat sessions", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to list chat sessions")
	}

	return &pb.ListSessionsResponse{
		Sessions:   sessions,
		TotalCount: total,
		Page:       page,
		PageSize:   pageSize,
	}, nil
}

// GetSession 获取会话详情
// @Summary 获取会话详情
// @Description 获取指定会话的详细信息
// @Tags AI Chat
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param session_id path string true "会话ID"
// @Success 200 {object} pb.SessionResponse
// @Failure 400 {object} status.Status "Invalid request"
// @Failure 401 {object} status.Status "Unauthorized"
// @Failure 404 {object} status.Status "Session not found"
// @Failure 500 {object} status.Status "Internal server error"
// @Router /v1/chat/sessions/{session_id} [get]
func (s *AIService) GetSession(ctx context.Context, req *pb.GetSessionRequest) (*pb.SessionResponse, error) {
	log := logger.GetLogger(ctx)

	// Validate user
	userID, _, err := s.tokenService.ValidateTokenFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	// 校验请求参数
	if req.SessionId == "" {
		return nil, status.Error(codes.InvalidArgument, "session_id is required")
	}

	// 调用服务获取会话详情
	session, err := s.aiService.GetSession(ctx, fmt.Sprintf("%d", userID), req.SessionId)
	if err != nil {
		log.Error("Failed to get chat session", zap.Error(err), zap.String("session_id", req.SessionId))
		return nil, status.Error(codes.Internal, "failed to get chat session")
	}

	return session, nil
}

// DeleteSession 删除会话
// @Summary 删除会话
// @Description 删除指定的聊天会话
// @Tags AI Chat
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param session_id path string true "会话ID"
// @Success 200 {object} emptypb.Empty
// @Failure 400 {object} status.Status "Invalid request"
// @Failure 401 {object} status.Status "Unauthorized"
// @Failure 404 {object} status.Status "Session not found"
// @Failure 500 {object} status.Status "Internal server error"
// @Router /v1/chat/sessions/{session_id} [delete]
func (s *AIService) DeleteSession(ctx context.Context, req *pb.DeleteSessionRequest) (*emptypb.Empty, error) {
	log := logger.GetLogger(ctx)

	// Validate user
	userID, _, err := s.tokenService.ValidateTokenFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	// 校验请求参数
	if req.SessionId == "" {
		return nil, status.Error(codes.InvalidArgument, "session_id is required")
	}

	// 调用服务删除会话
	err = s.aiService.DeleteSession(ctx, fmt.Sprintf("%d", userID), req.SessionId)
	if err != nil {
		log.Error("Failed to delete chat session", zap.Error(err), zap.String("session_id", req.SessionId))
		return nil, status.Error(codes.Internal, "failed to delete chat session")
	}

	return &emptypb.Empty{}, nil
}

func toDomainChatContext(ctx *pb.ChatContext) *service.ChatContext {
	if ctx == nil {
		return nil
	}
	return &service.ChatContext{
		History:   ctx.History,
		SessionID: ctx.SessionId,
	}
}
