package grpc

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	pb "github.com/lazyjean/sla2/api/proto/v1"
	"github.com/lazyjean/sla2/config"
	"github.com/lazyjean/sla2/internal/application/service"
	"github.com/lazyjean/sla2/internal/domain/security"
	"github.com/lazyjean/sla2/internal/interfaces/grpc/admin"
	"github.com/lazyjean/sla2/internal/interfaces/grpc/course"
	"github.com/lazyjean/sla2/internal/interfaces/grpc/learning"
	"github.com/lazyjean/sla2/internal/interfaces/grpc/middleware"
	"github.com/lazyjean/sla2/internal/interfaces/grpc/question"
	"github.com/lazyjean/sla2/internal/interfaces/grpc/user"
	"github.com/lazyjean/sla2/internal/interfaces/grpc/word"
	"github.com/lazyjean/sla2/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

// Server gRPC 服务器
type Server struct {
	grpcServer     *grpc.Server
	httpServer     *http.Server
	adminSvc       *admin.AdminService
	userSvc        *user.UserService
	wordSvc        *word.WordService
	learningSvc    *learning.LearningService
	courseSvc      *course.CourseService
	questionSvc    *question.QuestionService
	questionTagSvc *question.QuestionTagServiceGRPC
	tokenService   security.TokenService
	config         *config.Config
	stopCh         chan struct{} // 用于通知所有 goroutine 停止
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
	tokenService security.TokenService,
	cfg *config.Config,
) *Server {
	s := &Server{
		tokenService: tokenService,
		config:       cfg,
		stopCh:       make(chan struct{}),
	}

	// 创建带有拦截器的 gRPC 服务器
	s.grpcServer = grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			middleware.ValidatorInterceptor(),
			s.loggingInterceptor(),
			s.authInterceptor(),
		),
		grpc.ChainStreamInterceptor(
			middleware.StreamValidatorInterceptor(),
		),
	)

	s.adminSvc = admin.NewAdminService(adminService)
	s.userSvc = user.NewUserService(userService)
	s.wordSvc = word.NewWordService(wordService)
	s.courseSvc = course.NewCourseService(courseService)
	s.learningSvc = learning.NewLearningService(learningService)
	s.questionSvc = question.NewQuestionService(questionService)
	s.questionTagSvc = question.NewQuestionTagServiceGRPC(questionTagService)

	// 注册服务
	pb.RegisterAdminServiceServer(s.grpcServer, s.adminSvc)
	pb.RegisterUserServiceServer(s.grpcServer, s.userSvc)
	pb.RegisterWordServiceServer(s.grpcServer, s.wordSvc)
	pb.RegisterCourseServiceServer(s.grpcServer, s.courseSvc)
	pb.RegisterLearningServiceServer(s.grpcServer, s.learningSvc)
	pb.RegisterQuestionServiceServer(s.grpcServer, s.questionSvc)
	pb.RegisterQuestionTagServiceServer(s.grpcServer, s.questionTagSvc)

	// 开发环境注册反射服务
	if cfg.GRPC.Reflection {
		reflection.Register(s.grpcServer)
	}

	return s
}

// 不需要鉴权的方法列表
var noAuthMethods = map[string]bool{
	"/proto.v1.AdminService/CheckSystemStatus": true,
	"/proto.v1.AdminService/InitializeSystem":  true,
	"/proto.v1.AdminService/AdminLogin":        true,
	"/proto.v1.AdminService/RefreshToken":      true,
	"/proto.v1.UserService/Login":              true,
	"/proto.v1.UserService/UserRegister":       true,
	"/proto.v1.UserService/RefreshToken":       true,
	"/proto.v1.UserService/LoginWithApple":     true,
}

// authInterceptor 创建鉴权拦截器
func (s *Server) authInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// 检查是否需要鉴权
		if noAuthMethods[info.FullMethod] {
			return handler(ctx, req)
		}

		// 从元数据中获取令牌
		token, err := s.extractToken(ctx)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "missing auth token")
		}

		// Admin 服务统一使用 admin 的验证方法
		if strings.HasPrefix(info.FullMethod, "/proto.v1.AdminService/") {
			userID, roles, err := s.tokenService.ValidateToken(token)
			if err != nil {
				return nil, status.Error(codes.Unauthenticated, "invalid admin token")
			}
			// 检查是否有管理员角色
			hasAdminRole := false
			for _, role := range roles {
				if role == "admin" {
					hasAdminRole = true
					break
				}
			}
			if !hasAdminRole {
				return nil, status.Error(codes.PermissionDenied, "requires admin role")
			}
			newCtx := context.WithValue(ctx, service.AdminIDKey, userID)
			newCtx = context.WithValue(newCtx, service.UserRolesKey, roles)
			return handler(newCtx, req)
		}

		// 其他服务使用 TokenService 验证
		userID, roles, err := s.tokenService.ValidateToken(token)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "invalid token")
		}

		// 将用户ID和角色信息存入上下文
		newCtx := context.WithValue(ctx, service.UserIDKey, userID)
		newCtx = context.WithValue(newCtx, service.UserRolesKey, roles)
		return handler(newCtx, req)
	}
}

// extractToken 从元数据中提取令牌
func (s *Server) extractToken(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "missing metadata")
	}

	values := md.Get("authorization")
	if len(values) == 0 {
		return "", status.Error(codes.Unauthenticated, "missing authorization header")
	}

	auth := values[0]
	if !strings.HasPrefix(auth, "Bearer ") {
		return "", status.Error(codes.Unauthenticated, "invalid authorization format")
	}

	return strings.TrimPrefix(auth, "Bearer "), nil
}

// loggingInterceptor 创建日志拦截器
func (s *Server) loggingInterceptor() grpc.UnaryServerInterceptor {
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

// Start 启动 gRPC 服务
func (s *Server) Start(port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("failed to listen gRPC: %v", err)
	}

	// 启动 gRPC 服务
	go func() {
		if err := s.grpcServer.Serve(lis); err != nil {
			logger.GetLogger(context.Background()).Error("gRPC server failed", zap.Error(err))
		}
	}()

	// 等待停止信号
	<-s.stopCh

	return nil
}

// StartGateway 启动 gRPC HTTP 网关服务
func (s *Server) StartGateway(gatewayPort int, grpcPort int) error {
	ctx := context.Background()
	mux := runtime.NewServeMux(
		runtime.WithIncomingHeaderMatcher(customHeaderMatcher),
		runtime.WithOutgoingHeaderMatcher(customHeaderMatcher),
	)
	opts := []grpc.DialOption{grpc.WithInsecure()}

	// 注册 HTTP 网关
	if err := pb.RegisterAdminServiceHandlerFromEndpoint(ctx, mux, fmt.Sprintf("localhost:%d", grpcPort), opts); err != nil {
		return fmt.Errorf("failed to register admin service gateway: %v", err)
	}
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
	if err := pb.RegisterQuestionServiceHandlerFromEndpoint(ctx, mux, fmt.Sprintf("localhost:%d", grpcPort), opts); err != nil {
		return fmt.Errorf("failed to register question service gateway: %v", err)
	}
	if err := pb.RegisterQuestionTagServiceHandlerFromEndpoint(ctx, mux, fmt.Sprintf("localhost:%d", grpcPort), opts); err != nil {
		return fmt.Errorf("failed to register question tag service gateway: %v", err)
	}

	// 创建一个新的 http.ServeMux 来处理所有请求
	rootMux := http.NewServeMux()

	// 注册 API 文档，添加基本认证
	rootMux.HandleFunc("/swagger/", func(w http.ResponseWriter, r *http.Request) {
		// todo: disable basic auth tempory
		// username, password, ok := r.BasicAuth()
		// if !ok || !s.validateSwaggerAuth(username, password) {
		// 	w.Header().Set("WWW-Authenticate", `Basic realm="Swagger API Documentation"`)
		// 	http.Error(w, "Unauthorized", http.StatusUnauthorized)
		// 	return
		// }
		http.StripPrefix("/swagger/", http.FileServer(http.Dir("./api/swagger"))).ServeHTTP(w, r)
	})

	// 注册 gRPC-Gateway 处理程序
	rootMux.Handle("/", mux)

	s.httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", gatewayPort),
		Handler: rootMux,
	}

	// 启动 HTTP 网关
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.GetLogger(context.Background()).Error("HTTP gateway failed", zap.Error(err))
		}
	}()

	// 等待停止信号
	<-s.stopCh

	return nil
}

// validateSwaggerAuth 验证 swagger 文档的访问凭证
func (s *Server) validateSwaggerAuth(username, password string) bool {
	if s.config.Swagger.Username == "" || s.config.Swagger.Password == "" {
		// 如果没有配置凭证，使用默认值
		return username == "admin" && password == "swagger"
	}
	return username == s.config.Swagger.Username && password == s.config.Swagger.Password
}

// customHeaderMatcher 自定义请求头匹配器
func customHeaderMatcher(key string) (string, bool) {
	switch strings.ToLower(key) {
	case "authorization", "set-cookie":
		return key, true
	default:
		return runtime.DefaultHeaderMatcher(key)
	}
}

// Stop 停止 gRPC 和 HTTP 网关服务
func (s *Server) Stop() {
	// 发送停止信号
	close(s.stopCh)

	// 创建一个带超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 先停止 HTTP 网关
	if s.httpServer != nil {
		if err := s.httpServer.Shutdown(ctx); err != nil {
			logger.GetLogger(ctx).Error("HTTP gateway shutdown failed", zap.Error(err))
		}
	}

	// 停止 gRPC 服务
	if s.grpcServer != nil {
		s.grpcServer.GracefulStop()
	}
}
