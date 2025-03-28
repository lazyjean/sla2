package grpc

import (
	"context"
	"crypto/subtle"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	pb "github.com/lazyjean/sla2/api/proto/v1"
	"github.com/lazyjean/sla2/config"
	_ "github.com/lazyjean/sla2/docs" // 导入 Swagger 文档
	"github.com/lazyjean/sla2/internal/application/service"
	"github.com/lazyjean/sla2/internal/domain/security"
	"github.com/lazyjean/sla2/internal/interfaces/grpc/admin"
	"github.com/lazyjean/sla2/internal/interfaces/grpc/course"
	"github.com/lazyjean/sla2/internal/interfaces/grpc/learning"
	grpcmiddleware "github.com/lazyjean/sla2/internal/interfaces/grpc/middleware"
	"github.com/lazyjean/sla2/internal/interfaces/grpc/question"
	"github.com/lazyjean/sla2/internal/interfaces/grpc/user"
	"github.com/lazyjean/sla2/internal/interfaces/grpc/vocabulary"
	"github.com/lazyjean/sla2/internal/interfaces/http/ws/handler"
	"github.com/lazyjean/sla2/pkg/logger"
	"github.com/soheilhy/cmux"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/proto"
)

// GRPCServer gRPC 服务器
type GRPCServer struct {
	config            *config.Config
	grpcServer        *grpc.Server
	httpServer        *http.Server
	userService       *service.UserService
	questionService   *service.QuestionService
	vocabularyService *service.VocabularyService
	courseService     *service.CourseService
	learningService   *service.LearningService
	adminService      *service.AdminService
	wsHandler         *handler.WebSocketHandler
	unaryInterceptor  grpc.UnaryServerInterceptor
	upgrader          websocket.Upgrader
	mux               *runtime.ServeMux
	stopCh            chan struct{}
	tokenService      security.TokenService
}

// healthServer 健康检查服务实现
type healthServer struct {
	grpc_health_v1.UnimplementedHealthServer
}

func (s *healthServer) Check(ctx context.Context, req *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	return &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	}, nil
}

func (s *healthServer) Watch(req *grpc_health_v1.HealthCheckRequest, stream grpc_health_v1.Health_WatchServer) error {
	// 创建一个定时器，每5秒发送一次健康状态
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	// 定期发送健康状态
	for {
		select {
		case <-stream.Context().Done():
			return stream.Context().Err()
		case <-ticker.C:
			if err := stream.Send(&grpc_health_v1.HealthCheckResponse{
				Status: grpc_health_v1.HealthCheckResponse_SERVING,
			}); err != nil {
				return err
			}
		}
	}
}

// NewGRPCServer 创建新的 gRPC 服务器
func NewGRPCServer(
	userService *service.UserService,
	questionService *service.QuestionService,
	vocabularyService *service.VocabularyService,
	courseService *service.CourseService,
	learningService *service.LearningService,
	adminService *service.AdminService,
	wsHandler *handler.WebSocketHandler,
	tokenService security.TokenService,
) *GRPCServer {
	ctx := context.Background()
	mux := runtime.NewServeMux(
		runtime.WithMetadata(func(ctx context.Context, req *http.Request) metadata.MD {
			log := logger.GetLogger(ctx)
			log.Info("gRPC-Gateway request",
				zap.String("method", req.Method),
				zap.String("path", req.URL.Path),
				zap.String("query", req.URL.RawQuery),
				zap.String("remote_addr", req.RemoteAddr),
				zap.String("user_agent", req.UserAgent()),
			)

			md := metadata.MD{}

			// 从 Authorization header 获取 token
			if auth := req.Header.Get("Authorization"); auth != "" {
				md = metadata.Join(md, metadata.Pairs("authorization", auth))
			}

			// 从 cookie 获取 token
			if cookie, err := req.Cookie(grpcmiddleware.TokenCookieName); err == nil {
				md = metadata.Join(md, metadata.Pairs("authorization", "Bearer "+cookie.Value))
			}

			return md
		}),
		runtime.WithForwardResponseOption(func(ctx context.Context, w http.ResponseWriter, resp proto.Message) error {
			log := logger.GetLogger(ctx)
			md, ok := runtime.ServerMetadataFromContext(ctx)
			if !ok {
				log.Warn("No server metadata in context")
				return nil
			}

			// 处理 cookie 设置
			if vals := md.HeaderMD.Get("set-cookie-token"); len(vals) > 0 {
				grpcmiddleware.SetTokenCookie(ctx, w, vals[0])
			}

			return nil
		}),
		runtime.WithErrorHandler(func(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {
			log := logger.GetLogger(ctx)
			log.Error("gRPC-Gateway error",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.Error(err),
			)
			runtime.DefaultHTTPErrorHandler(ctx, mux, marshaler, w, r, err)
		}),
		runtime.WithOutgoingHeaderMatcher(func(key string) (string, bool) {
			// 不需要转换任何 metadata 为 header，因为我们使用 SetCookie 直接设置
			return "", false
		}),
	)

	// todo: 整理这里的拦截器
	// 创建拦截器链
	unaryInterceptors := []grpc.UnaryServerInterceptor{
		grpc_recovery.UnaryServerInterceptor(),
		grpc_ctxtags.UnaryServerInterceptor(),
		grpc_prometheus.UnaryServerInterceptor,
		grpc_zap.UnaryServerInterceptor(logger.GetLogger(context.Background())),
		grpcmiddleware.MetadataLoggerUnaryInterceptor(),
		grpcmiddleware.UnaryServerInterceptor(tokenService),
		grpc_validator.UnaryServerInterceptor(),
	}

	streamInterceptors := []grpc.StreamServerInterceptor{
		grpc_recovery.StreamServerInterceptor(),
		grpc_ctxtags.StreamServerInterceptor(),
		grpc_prometheus.StreamServerInterceptor,
		grpc_zap.StreamServerInterceptor(logger.GetLogger(context.Background())),
		grpcmiddleware.MetadataLoggerStreamInterceptor(),
		grpcmiddleware.StreamServerInterceptor(tokenService),
		grpc_validator.StreamServerInterceptor(),
	}

	// 创建 gRPC 服务器
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(unaryInterceptors...),
		grpc.ChainStreamInterceptor(streamInterceptors...),
	)

	// 注册 gRPC-Gateway 路由
	opts := []grpc.DialOption{grpc.WithInsecure()}
	endpoint := fmt.Sprintf("localhost:%d", config.GetConfig().GRPC.Port)

	// 注册 HTTP 处理器
	if err := pb.RegisterUserServiceHandlerFromEndpoint(ctx, mux, endpoint, opts); err != nil {
		logger.GetLogger(ctx).Fatal("failed to register user service handler", zap.Error(err))
	}
	if err := pb.RegisterQuestionServiceHandlerFromEndpoint(ctx, mux, endpoint, opts); err != nil {
		logger.GetLogger(ctx).Fatal("failed to register question service handler", zap.Error(err))
	}
	if err := pb.RegisterVocabularyServiceHandlerFromEndpoint(ctx, mux, endpoint, opts); err != nil {
		logger.GetLogger(ctx).Fatal("failed to register vocabulary service handler", zap.Error(err))
	}
	if err := pb.RegisterCourseServiceHandlerFromEndpoint(ctx, mux, endpoint, opts); err != nil {
		logger.GetLogger(ctx).Fatal("failed to register course service handler", zap.Error(err))
	}
	if err := pb.RegisterLearningServiceHandlerFromEndpoint(ctx, mux, endpoint, opts); err != nil {
		logger.GetLogger(ctx).Fatal("failed to register learning service handler", zap.Error(err))
	}
	if err := pb.RegisterAdminServiceHandlerFromEndpoint(ctx, mux, endpoint, opts); err != nil {
		logger.GetLogger(ctx).Fatal("failed to register admin service handler", zap.Error(err))
	}

	// 创建 HTTP 路由
	router := http.NewServeMux()

	// 注册 Swagger UI
	swaggerHandler := httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("/swagger/doc.json")),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	)

	// 创建服务器实例
	srv := &GRPCServer{
		config:            config.GetConfig(),
		grpcServer:        grpcServer,
		userService:       userService,
		questionService:   questionService,
		vocabularyService: vocabularyService,
		courseService:     courseService,
		learningService:   learningService,
		adminService:      adminService,
		wsHandler:         wsHandler,
		mux:               mux,
		stopCh:            make(chan struct{}),
		tokenService:      tokenService,
	}

	// 需要 Basic Auth 的路由
	router.Handle("/swagger/", srv.basicAuth(swaggerHandler))

	// 直接提供swagger.json文件
	router.HandleFunc("/swagger/doc.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// 读取swagger.json文件
		swaggerFile := "./api/swagger/swagger.json"
		data, err := os.ReadFile(swaggerFile)
		if err != nil {
			logger.GetLogger(context.Background()).Error("failed to read swagger.json", zap.Error(err))
			http.Error(w, "Swagger file not found", http.StatusNotFound)
			return
		}

		w.Write(data)
	})

	// API 路由
	router.Handle("/api/", srv.mux)

	// 健康检查
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	// 注册健康检查服务
	grpc_health_v1.RegisterHealthServer(grpcServer, &healthServer{})

	// 注册反射服务
	reflection.Register(grpcServer)

	// 注册服务
	srv.registerServices()

	return srv
}

func (s *GRPCServer) registerServices() {
	// 注册用户服务
	pb.RegisterUserServiceServer(s.grpcServer, user.NewUserService(s.userService))

	// 注册问题服务
	pb.RegisterQuestionServiceServer(s.grpcServer, question.NewQuestionService(s.questionService))

	// 注册词汇服务
	pb.RegisterVocabularyServiceServer(s.grpcServer, vocabulary.NewVocabularyService(s.vocabularyService))

	// 注册课程服务
	pb.RegisterCourseServiceServer(s.grpcServer, course.NewCourseService(s.courseService))

	// 注册学习服务
	pb.RegisterLearningServiceServer(s.grpcServer, learning.NewLearningService(s.learningService))

	// 注册管理员服务
	pb.RegisterAdminServiceServer(s.grpcServer, admin.NewAdminService(s.adminService))

	// 注册健康检查和反射服务已在 NewGRPCServer 中完成
}

// basicAuth 中间件
func (s *GRPCServer) basicAuth(next http.Handler) http.Handler {
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
func (s *GRPCServer) Start() error {
	// 创建 TCP 监听器
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.config.GRPC.Port))
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	// 创建 cmux
	m := cmux.New(lis)

	// 匹配 gRPC 流量 (HTTP/2)
	grpcL := m.Match(cmux.HTTP2())
	// 匹配 HTTP/1.1 流量
	httpL := m.Match(cmux.HTTP1Fast())

	// 创建 HTTP 路由
	router := http.NewServeMux()

	// 注册 Swagger UI
	swaggerHandler := httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("/swagger/doc.json")),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	)

	// 需要 Basic Auth 的路由
	router.Handle("/swagger/", s.basicAuth(swaggerHandler))

	// 直接提供swagger.json文件
	router.HandleFunc("/swagger/doc.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// 读取swagger.json文件
		swaggerFile := "./api/swagger/swagger.json"
		data, err := os.ReadFile(swaggerFile)
		if err != nil {
			logger.GetLogger(context.Background()).Error("failed to read swagger.json", zap.Error(err))
			http.Error(w, "Swagger file not found", http.StatusNotFound)
			return
		}

		w.Write(data)
	})

	// API 路由
	router.Handle("/api/", s.mux)

	// 健康检查
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	// 创建 HTTP 服务器
	s.httpServer = &http.Server{
		Addr: fmt.Sprintf(":%d", config.GetConfig().GRPC.Port),
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
				s.grpcServer.ServeHTTP(w, r)
			} else {
				router.ServeHTTP(w, r)
			}
		}),
	}

	// 启动 gRPC 服务器
	go func() {
		logger.GetLogger(context.Background()).Info("gRPC server is running", zap.Int("port", s.config.GRPC.Port))
		if err := s.grpcServer.Serve(grpcL); err != nil {
			logger.GetLogger(context.Background()).Fatal("failed to serve gRPC", zap.Error(err))
		}
	}()

	// 启动 HTTP 服务器
	go func() {
		logger.GetLogger(context.Background()).Info("HTTP server is running",
			zap.Int("port", s.config.GRPC.Port),
			zap.String("swagger_url", fmt.Sprintf("http://localhost:%d/swagger/", s.config.GRPC.Port)),
			zap.String("swagger_username", s.config.Swagger.Username),
			zap.String("swagger_password", s.config.Swagger.Password),
		)
		if err := s.httpServer.Serve(httpL); err != nil && err != http.ErrServerClosed {
			logger.GetLogger(context.Background()).Fatal("failed to serve HTTP", zap.Error(err))
		}
	}()

	// 启动 cmux
	logger.GetLogger(context.Background()).Info("Server is running", zap.Int("port", s.config.GRPC.Port))
	if err := m.Serve(); err != nil {
		return fmt.Errorf("failed to serve cmux: %w", err)
	}

	return nil
}

// Stop 停止服务器
func (s *GRPCServer) Stop() error {
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
