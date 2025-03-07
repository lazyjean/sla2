package ai

import (
	"context"
	"fmt"
	"time"

	pb "github.com/lazyjean/sla2/api/proto/v1"
	"github.com/lazyjean/sla2/internal/application/service"
	"github.com/lazyjean/sla2/internal/domain/security"
	"github.com/lazyjean/sla2/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AIService struct {
	pb.UnimplementedAIChatServiceServer
	aiService    service.AIService
	tokenService security.TokenService
}

func NewAIService(aiService service.AIService, tokenService security.TokenService) *AIService {
	return &AIService{
		aiService:    aiService,
		tokenService: tokenService,
	}
}

func (s *AIService) Chat(ctx context.Context, req *pb.ChatRequest) (*pb.ChatResponse, error) {
	log := logger.GetLogger(ctx)

	// Validate user
	userID, _, err := s.tokenService.ValidateTokenFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	// Call AI service
	response, err := s.aiService.Chat(ctx, fmt.Sprintf("%d", userID), req.Message, toDomainChatContext(req.Context))
	if err != nil {
		log.Error("AI chat failed", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to process chat request")
	}

	return &pb.ChatResponse{
		Message:   response.Message,
		CreatedAt: timestamppb.New(time.Unix(response.CreatedAt, 0)),
	}, nil
}

func (s *AIService) StreamChat(req *pb.StreamChatRequest, stream pb.AIChatService_StreamChatServer) error {
	ctx := stream.Context()
	log := logger.GetLogger(ctx)

	// Validate user
	userID, _, err := s.tokenService.ValidateTokenFromContext(ctx)
	if err != nil {
		return status.Error(codes.Unauthenticated, "invalid token")
	}

	// Call AI service
	responseStream, err := s.aiService.StreamChat(ctx, fmt.Sprintf("%d", userID), req.Message, toDomainChatContext(req.Context))
	if err != nil {
		log.Error("AI stream chat failed", zap.Error(err))
		return status.Error(codes.Internal, "failed to process stream chat request")
	}

	for response := range responseStream {
		if err := stream.Send(&pb.ChatResponse{
			Message:   response.Message,
			CreatedAt: timestamppb.New(time.Unix(response.CreatedAt, 0)),
		}); err != nil {
			return err
		}
	}

	return nil
}

func toDomainChatContext(ctx *pb.ChatContext) *service.ChatContext {
	if ctx == nil {
		return nil
	}
	return &service.ChatContext{
		History: ctx.History,
	}
}
