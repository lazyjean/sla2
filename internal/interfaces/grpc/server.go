package grpc

import (
	"fmt"
	"net"

	pb "github.com/lazyjean/sla2/api/proto/v1"
	"github.com/lazyjean/sla2/internal/domain/repository"
	"github.com/lazyjean/sla2/internal/interfaces/grpc/learning"
	"github.com/lazyjean/sla2/internal/interfaces/grpc/word"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Server gRPC 服务器
type Server struct {
	server       *grpc.Server
	wordRepo     repository.WordRepository
	learningRepo repository.LearningRepository
}

// NewServer 创建 gRPC 服务器
func NewServer(wordRepo repository.WordRepository, learningRepo repository.LearningRepository) *Server {
	return &Server{
		server:       grpc.NewServer(),
		wordRepo:     wordRepo,
		learningRepo: learningRepo,
	}
}

// Start 启动 gRPC 服务器
func (s *Server) Start(port string) error {
	// 注册服务
	wordService := word.NewWordService(s.wordRepo)
	pb.RegisterWordServiceServer(s.server, wordService)

	learningService := learning.NewLearningService(s.learningRepo)
	pb.RegisterLearningServiceServer(s.server, learningService)

	// 启用反射服务
	reflection.Register(s.server)

	// 监听端口
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	// 启动服务器
	if err := s.server.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}

	return nil
}

// Stop 停止 gRPC 服务器
func (s *Server) Stop() {
	s.server.GracefulStop()
}
