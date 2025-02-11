package grpc

import (
	"fmt"
	"net"

	pb "github.com/lazyjean/sla2/api/proto/v1"
	"github.com/lazyjean/sla2/internal/application/service"
	"github.com/lazyjean/sla2/internal/interfaces/grpc/learning"
	"github.com/lazyjean/sla2/internal/interfaces/grpc/user"
	"github.com/lazyjean/sla2/internal/interfaces/grpc/word"
	"github.com/lazyjean/sla2/pkg/auth"
	"google.golang.org/grpc"
)

// Server gRPC 服务器
type Server struct {
	server      *grpc.Server
	userSvc     *user.UserService
	wordSvc     *word.WordService
	learningSvc *learning.LearningService
}

// NewServer 创建 gRPC 服务器
func NewServer(userService *service.UserService, wordService *service.WordService, learningService *service.LearningService, authSvc auth.JWTServicer) *Server {
	server := grpc.NewServer()

	s := &Server{
		server:      server,
		userSvc:     user.NewUserService(userService, authSvc),
		wordSvc:     word.NewWordService(wordService),
		learningSvc: learning.NewLearningService(learningService),
	}

	// 注册服务
	pb.RegisterUserServiceServer(server, s.userSvc)
	pb.RegisterWordServiceServer(server, s.wordSvc)
	pb.RegisterLearningServiceServer(server, s.learningSvc)
	return s
}

// Start 启动 gRPC 服务器
func (s *Server) Start(port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	if err := s.server.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}

	return nil
}

// Stop 停止 gRPC 服务器
func (s *Server) Stop() {
	s.server.GracefulStop()
}
