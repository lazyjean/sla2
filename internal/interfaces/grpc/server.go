package grpc

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
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
	grpcServer  *grpc.Server
	httpServer  *http.Server
	userSvc     *user.UserService
	wordSvc     *word.WordService
	learningSvc *learning.LearningService
	courseSvc   *course.CourseService
}

// NewServer 创建 gRPC 服务器
func NewServer(userService *service.UserService, wordService *service.WordService,
	learningService *service.LearningService, courseService *service.CourseService,
	authSvc auth.JWTServicer, cfg *config.Config) *Server {

	// 创建带有拦截器的 gRPC 服务器
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(loggingInterceptor()),
	)

	s := &Server{
		grpcServer:  grpcServer,
		userSvc:     user.NewUserService(userService, authSvc),
		wordSvc:     word.NewWordService(wordService),
		courseSvc:   course.NewCourseService(courseService),
		learningSvc: learning.NewLearningService(learningService),
	}

	// 注册服务
	pb.RegisterUserServiceServer(grpcServer, s.userSvc)
	pb.RegisterWordServiceServer(grpcServer, s.wordSvc)
	pb.RegisterCourseServiceServer(grpcServer, s.courseSvc)
	pb.RegisterLearningServiceServer(grpcServer, s.learningSvc)

	// 开发环境注册反射服务
	if cfg.GRPC.Reflection {
		reflection.Register(grpcServer)
	}

	return s
}

// Start 启动 gRPC 服务
func (s *Server) Start(port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("failed to listen gRPC: %v", err)
	}

	go func() {
		if err := s.grpcServer.Serve(lis); err != nil {
			logger.GetLogger(context.Background()).Error("gRPC server failed", zap.Error(err))
		}
	}()

	return nil
}

// StartGateway 启动 gRPC HTTP 网关服务
func (s *Server) StartGateway(gatewayPort int, grpcPort int) error {
	ctx := context.Background()
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	// 注册 HTTP 网关
	if err := pb.RegisterUserServiceHandlerFromEndpoint(ctx, mux, fmt.Sprintf("localhost:%d", grpcPort), opts); err != nil {
		return fmt.Errorf("failed to register user service gateway: %v", err)
	}
	if err := pb.RegisterWordServiceHandlerFromEndpoint(ctx, mux, fmt.Sprintf("localhost:%d", grpcPort), opts); err != nil {
		return fmt.Errorf("failed to register word service gateway: %v", err)
	}
	if err := pb.RegisterCourseServiceHandlerFromEndpoint(ctx, mux, fmt.Sprintf("localhost:%d", grpcPort), opts); err != nil {
		return fmt.Errorf("failed to register course service gateway: %v", err)
	}
	if err := pb.RegisterLearningServiceHandlerFromEndpoint(ctx, mux, fmt.Sprintf("localhost:%d", grpcPort), opts); err != nil {
		return fmt.Errorf("failed to register learning service gateway: %v", err)
	}

	s.httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", gatewayPort),
		Handler: mux,
	}

	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.GetLogger(context.Background()).Error("HTTP gateway failed", zap.Error(err))
		}
	}()

	return nil
}

// Stop 停止 gRPC 和 HTTP 网关服务
func (s *Server) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 先停止 HTTP 网关
	if err := s.httpServer.Shutdown(ctx); err != nil {
		logger.GetLogger(ctx).Error("HTTP gateway shutdown failed", zap.Error(err))
	}

	// 停止 gRPC 服务
	s.grpcServer.GracefulStop()
}

// loggingInterceptor 创建日志拦截器
func loggingInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		log := logger.GetLogger(ctx)
		start := time.Now()

		log.Info("Request", zap.String("method", info.FullMethod), zap.Any("payload", req))

		resp, err := handler(ctx, req)

		duration := time.Since(start)

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
