package grpc

import (
	"context"
	"fmt"
	"net"
	"time"

	pb "github.com/lazyjean/sla2/api/proto/v1"
	"github.com/lazyjean/sla2/config"
	"github.com/lazyjean/sla2/internal/application/service"
	"github.com/lazyjean/sla2/internal/interfaces/grpc/course"
	"github.com/lazyjean/sla2/internal/interfaces/grpc/learning"
	"github.com/lazyjean/sla2/internal/interfaces/grpc/user"
	"github.com/lazyjean/sla2/internal/interfaces/grpc/word"
	"github.com/lazyjean/sla2/pkg/auth"
	"github.com/lazyjean/sla2/pkg/logger"
	"go.uber.org/zap"
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
	// 创建带有拦截器的 gRPC 服务器
	server := grpc.NewServer(
		grpc.UnaryInterceptor(loggingInterceptor()),
	)

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

// loggingInterceptor 创建日志拦截器
func loggingInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		log := logger.GetLogger(ctx)
		// 记录请求开始时间
		start := time.Now()

		// 记录请求信息
		log.Info("Request", zap.String("method", info.FullMethod), zap.Any("payload", req))

		// 处理请求
		resp, err := handler(ctx, req)

		// 计算处理时间
		duration := time.Since(start)

		// 记录响应信息
		if err != nil {
			log.Info("Response",
				zap.String("method", info.FullMethod),
				zap.Error(err),
				zap.Duration("duration", duration),
			)
		} else {
			log.Info("Response",
				zap.String("method", info.FullMethod),
				zap.Any("payload", resp),
				zap.Duration("duration", duration),
			)
		}

		return resp, err
	}
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
