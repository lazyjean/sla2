package grpc

import (
	"fmt"
	"net"

	pb "github.com/lazyjean/sla2/api/proto/v1"
	"github.com/lazyjean/sla2/config"
	"github.com/lazyjean/sla2/internal/application/service"
	"github.com/lazyjean/sla2/internal/interfaces/grpc/course"
	"github.com/lazyjean/sla2/internal/interfaces/grpc/learning"
	"github.com/lazyjean/sla2/internal/interfaces/grpc/user"
	"github.com/lazyjean/sla2/internal/interfaces/grpc/word"
	"github.com/lazyjean/sla2/pkg/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Server gRPC 服务器
type Server struct {
	server      *grpc.Server
	userSvc     *user.UserService
	wordSvc     *word.WordService
	learningSvc *learning.LearningService
	courseSvc   *course.CourseService
}

// NewServer 创建 gRPC 服务器
func NewServer(userService *service.UserService, wordService *service.WordService, learningService *service.LearningService, courseService *service.CourseService, authSvc auth.JWTServicer, cfg *config.Config) *Server {
	server := grpc.NewServer()

	s := &Server{
		server:      server,
		userSvc:     user.NewUserService(userService, authSvc),
		wordSvc:     word.NewWordService(wordService),
		courseSvc:   course.NewCourseService(courseService),
		learningSvc: learning.NewLearningService(learningService),
	}

	// 注册服务
	pb.RegisterUserServiceServer(server, s.userSvc)
	pb.RegisterWordServiceServer(server, s.wordSvc)
	pb.RegisterCourseServiceServer(server, s.courseSvc)
	pb.RegisterLearningServiceServer(server, s.learningSvc)

	// todo: 开发环境注册反射服务
	if cfg.GRPC.Reflection {
		reflection.Register(server)
	}

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
