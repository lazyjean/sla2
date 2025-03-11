package grpc

import (
	"context"
	"crypto/subtle"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	pb "github.com/lazyjean/sla2/api/proto/v1"
	"github.com/lazyjean/sla2/config"
	_ "github.com/lazyjean/sla2/docs" // 导入 Swagger 文档
	"github.com/lazyjean/sla2/internal/application/service"
	"github.com/lazyjean/sla2/internal/domain/security"
	"github.com/lazyjean/sla2/internal/interfaces/grpc/admin"
	"github.com/lazyjean/sla2/internal/interfaces/grpc/ai"
	"github.com/lazyjean/sla2/internal/interfaces/grpc/course"
	"github.com/lazyjean/sla2/internal/interfaces/grpc/learning"
	"github.com/lazyjean/sla2/internal/interfaces/grpc/middleware"
	"github.com/lazyjean/sla2/internal/interfaces/grpc/question"
	"github.com/lazyjean/sla2/internal/interfaces/grpc/user"
	"github.com/lazyjean/sla2/internal/interfaces/grpc/word"
	"github.com/lazyjean/sla2/internal/interfaces/websocket"
	"github.com/lazyjean/sla2/pkg/logger"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
)

// Server gRPC 服务器
type Server struct {
	grpcServer        *grpc.Server
	httpServer        *http.Server
	adminSvc          *admin.AdminService
	userSvc           *user.UserService
	wordSvc           *word.WordService
	learningSvc       *learning.LearningService
	courseSvc         *course.CourseService
	questionSvc       *question.QuestionService
	questionTagSvc    *question.QuestionTagServiceGRPC
	aiSvc             *ai.AIService
	tokenService      security.TokenService
	config            *config.Config
	stopCh            chan struct{}                 // 用于通知所有 goroutine 停止
	chatHandler       *websocket.UnifiedChatHandler // 统一的 WebSocket 聊天处理器
	unaryInterceptor  grpc.UnaryServerInterceptor
	streamInterceptor grpc.StreamServerInterceptor
}

// NewServer 创建 gRPC 服务器
func NewServer(
	adminService *service.AdminService,
	userService *service.UserService,
	wordService *service.WordService,
	learningService *service.LearningService,
	courseService *service.CourseService,
	questionService *service.QuestionService,
	questionTagService *service.QuestionTagService,
	aiService *service.AIService,
	tokenService security.TokenService,
	config *config.Config,
) *Server {
	// 创建拦截器
	unaryInterceptor := middleware.UnaryServerInterceptor(tokenService)

	// 创建 gRPC 服务器
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(unaryInterceptor),
	)

	// 创建统一的 WebSocket 聊天处理器
	chatHandler := websocket.NewUnifiedChatHandler(aiService, tokenService)

	// 创建服务器实例
	server := &Server{
		grpcServer:       grpcServer,
		adminSvc:         admin.NewAdminService(adminService),
		userSvc:          user.NewUserService(userService),
		wordSvc:          word.NewWordService(wordService),
		learningSvc:      learning.NewLearningService(learningService),
		courseSvc:        course.NewCourseService(courseService),
		questionSvc:      question.NewQuestionService(questionService),
		questionTagSvc:   question.NewQuestionTagServiceGRPC(questionTagService),
		aiSvc:            ai.NewAIService(aiService, tokenService),
		tokenService:     tokenService,
		config:           config,
		stopCh:           make(chan struct{}),
		chatHandler:      chatHandler,
		unaryInterceptor: unaryInterceptor,
	}

	// 注册服务
	server.registerServices()

	// 开发环境注册反射服务
	if config.GRPC.Reflection {
		reflection.Register(grpcServer)
	}

	return server
}

// registerServices 注册所有服务
func (s *Server) registerServices() {
	// 注册 gRPC 服务
	pb.RegisterAdminServiceServer(s.grpcServer, s.adminSvc)
	pb.RegisterUserServiceServer(s.grpcServer, s.userSvc)
	pb.RegisterWordServiceServer(s.grpcServer, s.wordSvc)
	pb.RegisterLearningServiceServer(s.grpcServer, s.learningSvc)
	pb.RegisterCourseServiceServer(s.grpcServer, s.courseSvc)
	pb.RegisterQuestionServiceServer(s.grpcServer, s.questionSvc)
	pb.RegisterQuestionTagServiceServer(s.grpcServer, s.questionTagSvc)
	pb.RegisterAIChatServiceServer(s.grpcServer, s.aiSvc)
}

// basicAuth 中间件
func (s *Server) basicAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok || subtle.ConstantTimeCompare([]byte(username), []byte(s.config.Swagger.Username)) != 1 || subtle.ConstantTimeCompare([]byte(password), []byte(s.config.Swagger.Password)) != 1 {
			w.Header().Set("WWW-Authenticate", `Basic realm="Swagger UI"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Start 启动服务器
func (s *Server) Start() error {
	// 启动 gRPC 服务器
	go func() {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.config.GRPC.Port))
		if err != nil {
			logger.GetLogger(context.Background()).Fatal("failed to listen", zap.Error(err))
		}

		logger.GetLogger(context.Background()).Info("gRPC server is running", zap.Int("port", s.config.GRPC.Port))
		if err := s.grpcServer.Serve(lis); err != nil {
			logger.GetLogger(context.Background()).Fatal("failed to serve", zap.Error(err))
		}
	}()

	// 启动 HTTP 服务器
	go func() {
		ctx := context.Background()
		mux := runtime.NewServeMux(
			runtime.WithMetadata(func(ctx context.Context, req *http.Request) metadata.MD {
				md := make(map[string]string)
				if auth := req.Header.Get("Authorization"); auth != "" {
					md["authorization"] = auth
				}
				return metadata.New(md)
			}),
		)

		opts := []grpc.DialOption{grpc.WithInsecure()}
		endpoint := fmt.Sprintf("localhost:%d", s.config.GRPC.Port)

		// 注册 HTTP 处理器
		if err := pb.RegisterAdminServiceHandlerFromEndpoint(ctx, mux, endpoint, opts); err != nil {
			logger.GetLogger(ctx).Fatal("failed to register admin service handler", zap.Error(err))
		}
		if err := pb.RegisterUserServiceHandlerFromEndpoint(ctx, mux, endpoint, opts); err != nil {
			logger.GetLogger(ctx).Fatal("failed to register user service handler", zap.Error(err))
		}
		if err := pb.RegisterWordServiceHandlerFromEndpoint(ctx, mux, endpoint, opts); err != nil {
			logger.GetLogger(ctx).Fatal("failed to register word service handler", zap.Error(err))
		}
		if err := pb.RegisterLearningServiceHandlerFromEndpoint(ctx, mux, endpoint, opts); err != nil {
			logger.GetLogger(ctx).Fatal("failed to register learning service handler", zap.Error(err))
		}
		if err := pb.RegisterCourseServiceHandlerFromEndpoint(ctx, mux, endpoint, opts); err != nil {
			logger.GetLogger(ctx).Fatal("failed to register course service handler", zap.Error(err))
		}
		if err := pb.RegisterQuestionServiceHandlerFromEndpoint(ctx, mux, endpoint, opts); err != nil {
			logger.GetLogger(ctx).Fatal("failed to register question service handler", zap.Error(err))
		}
		if err := pb.RegisterQuestionTagServiceHandlerFromEndpoint(ctx, mux, endpoint, opts); err != nil {
			logger.GetLogger(ctx).Fatal("failed to register question tag service handler", zap.Error(err))
		}
		if err := pb.RegisterAIChatServiceHandlerFromEndpoint(ctx, mux, endpoint, opts); err != nil {
			logger.GetLogger(ctx).Fatal("failed to register ai chat service handler", zap.Error(err))
		}

		// 创建 HTTP 路由
		router := http.NewServeMux()

		// 注册 WebSocket 处理器
		router.HandleFunc("/ws/chat", s.chatHandler.HandleWebSocket)

		// 注册 Swagger UI 处理器
		swaggerHandler := httpSwagger.Handler(
			httpSwagger.URL(fmt.Sprintf("http://localhost:%d/swagger/doc.json", s.config.GRPC.GatewayPort)),
			httpSwagger.DeepLinking(true),
			httpSwagger.DocExpansion("none"),
			httpSwagger.DomID("swagger-ui"),
		)
		router.Handle("/swagger/", s.basicAuth(swaggerHandler))

		// 注册 swagger.json 文件处理
		router.HandleFunc("/swagger/doc.json", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			http.ServeFile(w, r, "api/swagger/swagger.json")
		})

		// 注册 gRPC-Gateway 处理器
		router.Handle("/", mux)

		// 创建 HTTP 服务器
		s.httpServer = &http.Server{
			Addr:    fmt.Sprintf(":%d", s.config.GRPC.GatewayPort),
			Handler: router,
		}

		logger.GetLogger(ctx).Info("HTTP server is running", zap.Int("port", s.config.GRPC.GatewayPort))
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.GetLogger(ctx).Fatal("failed to serve HTTP", zap.Error(err))
		}
	}()

	return nil
}

// Stop 停止服务器
func (s *Server) Stop() error {
	close(s.stopCh)

	// 优雅关闭 HTTP 服务器
	if s.httpServer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := s.httpServer.Shutdown(ctx); err != nil {
			logger.GetLogger(context.Background()).Error("failed to shutdown HTTP server", zap.Error(err))
		}
	}

	// 优雅关闭 gRPC 服务器
	s.grpcServer.GracefulStop()

	return nil
}
