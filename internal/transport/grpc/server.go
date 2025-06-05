package grpc

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/apolloconfig/agollo/v4/component/log"
	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/selector"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/validator"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/prometheus/client_golang/prometheus"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	pb "github.com/lazyjean/sla2/api/proto/v1"
	"github.com/lazyjean/sla2/config"
	_ "github.com/lazyjean/sla2/docs" // 导入 Swagger 文档
	"github.com/lazyjean/sla2/internal/domain/security"
	"github.com/lazyjean/sla2/internal/infrastructure/listen"
	"github.com/lazyjean/sla2/internal/transport/grpc/admin"
	"github.com/lazyjean/sla2/internal/transport/grpc/course"
	"github.com/lazyjean/sla2/internal/transport/grpc/learning"
	"github.com/lazyjean/sla2/internal/transport/grpc/middleware"
	"github.com/lazyjean/sla2/internal/transport/grpc/question"
	"github.com/lazyjean/sla2/internal/transport/grpc/user"
	"github.com/lazyjean/sla2/internal/transport/grpc/vocabulary"
	"github.com/lazyjean/sla2/internal/transport/http/ws/handler"
)

// Server gRPC 服务器
type Server struct {
	config       *config.Config
	listener     listen.Listener
	grpcServer   *grpc.Server
	httpServer   *http.Server
	wsHandler    *handler.WebSocketHandler
	mux          *runtime.ServeMux
	stopCh       chan struct{}
	tokenService security.TokenService
	log          *zap.Logger
}

// NewGRPCServer 创建新的 gRPC 服务器
func NewGRPCServer(
	wsHandler *handler.WebSocketHandler,
	tokenService security.TokenService,
	listener listen.Listener,
	appLogger *zap.Logger,
	grpcServer *grpc.Server,
	userService *user.Service,
	learningService *learning.Service,
	courseService *course.Service,
	adminService *admin.Service,
	questionService *question.Service,
	vocabularyService *vocabulary.Service,
) *Server {
	mux := runtime.NewServeMux(
		runtime.WithMetadata(WithMetadata),
		runtime.WithForwardResponseOption(WithForwardResponseOption),
		runtime.WithErrorHandler(WithErrorHandler),
		runtime.WithOutgoingHeaderMatcher(WithOutgoingHeaderMatcher),
	)

	ctx := context.Background()
	// 注册 gRPC-Gateway 路由
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	endpoint := fmt.Sprintf("localhost:%d", config.GetConfig().GRPC.Port)

	// 注册 HTTP 处理器
	if err := pb.RegisterUserServiceHandlerFromEndpoint(ctx, mux, endpoint, opts); err != nil {
		appLogger.Error("failed to register user service handler", zap.Error(err))
		os.Exit(1)
	}
	if err := pb.RegisterQuestionServiceHandlerFromEndpoint(ctx, mux, endpoint, opts); err != nil {
		appLogger.Error("failed to register question service handler", zap.Error(err))
		os.Exit(1)
	}
	if err := pb.RegisterVocabularyServiceHandlerFromEndpoint(ctx, mux, endpoint, opts); err != nil {
		appLogger.Error("failed to register vocabulary service handler", zap.Error(err))
		os.Exit(1)
	}
	if err := pb.RegisterCourseServiceHandlerFromEndpoint(ctx, mux, endpoint, opts); err != nil {
		appLogger.Error("failed to register course service handler", zap.Error(err))
		os.Exit(1)
	}
	if err := pb.RegisterLearningServiceHandlerFromEndpoint(ctx, mux, endpoint, opts); err != nil {
		appLogger.Error("failed to register learning service handler", zap.Error(err))
		os.Exit(1)
	}
	if err := pb.RegisterAdminServiceHandlerFromEndpoint(ctx, mux, endpoint, opts); err != nil {
		appLogger.Error("failed to register admin service handler", zap.Error(err))
		os.Exit(1)
	}

	// 创建 HTTP 路由
	router := http.NewServeMux()

	// 注册 Swagger UI
	swaggerHandler := httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	)

	// 创建服务器实例
	srv := &Server{
		config:       config.GetConfig(),
		grpcServer:   grpcServer,
		wsHandler:    wsHandler,
		mux:          mux,
		stopCh:       make(chan struct{}),
		tokenService: tokenService,
		listener:     listener,
		log:          appLogger,
	}

	// 需要 Basic Auth 的路由
	router.Handle("/swagger/", srv.basicAuth(swaggerHandler))

	// 直接提供swagger.json文件
	router.HandleFunc("/swagger/doc.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// 读取swagger.json文件
		swaggerFile := "./api/swagger/swagger.json"
		http.ServeFile(w, r, swaggerFile)
	})

	// API 路由
	router.Handle("/api/", srv.mux)

	// 健康检查 todo: 抽象为统一的健康检查接口
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if _, er := w.Write([]byte(`{"status":"ok"}`)); er != nil {
			log.Error("write health response", zap.Error(er))
		}
	})

	// 注册健康检查服务
	grpc_health_v1.RegisterHealthServer(grpcServer, &healthServer{})

	// 注册反射服务（只使用 v1 版本）
	reflection.RegisterV1(grpcServer)

	// 注册服务
	pb.RegisterUserServiceServer(grpcServer, userService)
	pb.RegisterQuestionServiceServer(grpcServer, questionService)
	pb.RegisterVocabularyServiceServer(grpcServer, vocabularyService)
	pb.RegisterCourseServiceServer(grpcServer, courseService)
	pb.RegisterLearningServiceServer(grpcServer, learningService)
	pb.RegisterAdminServiceServer(grpcServer, adminService)

	return srv
}

func NewServer(metrics *grpcprom.ServerMetrics, appLogger *zap.Logger, tokenService security.TokenService) *grpc.Server {
	exemplarFromContext := func(ctx context.Context) prometheus.Labels {
		if span := trace.SpanContextFromContext(ctx); span.IsSampled() {
			return prometheus.Labels{"traceID": span.TraceID().String()}
		}
		return nil
	}

	// 创建拦截器链
	unaryInterceptors := []grpc.UnaryServerInterceptor{
		recovery.UnaryServerInterceptor(),
		metrics.UnaryServerInterceptor(grpcprom.WithExemplarFromContext(exemplarFromContext)), // todo: prometheus 是否需要带上 traceID?
		logging.UnaryServerInterceptor(interceptorLogger(appLogger), loggingOpts()...),
		selector.UnaryServerInterceptor(auth.UnaryServerInterceptor(middleware.AuthFunc(tokenService)), selector.MatchFunc(middleware.RequireAuth)),
		validator.UnaryServerInterceptor(),
	}

	streamInterceptors := []grpc.StreamServerInterceptor{
		recovery.StreamServerInterceptor(),
		metrics.StreamServerInterceptor(grpcprom.WithExemplarFromContext(exemplarFromContext)),
		logging.StreamServerInterceptor(interceptorLogger(appLogger), loggingOpts()...),
		auth.StreamServerInterceptor(middleware.AuthFunc(tokenService)),
		validator.StreamServerInterceptor(),
	}

	// 创建 gRPC 服务器
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(unaryInterceptors...),
		grpc.ChainStreamInterceptor(streamInterceptors...),
	)
	return grpcServer
}

// Start 启动服务器
func (s *Server) Start() error {
	// 创建 gRPC 监听器
	if err := s.listener.ListenGRPC(s.config.GRPC.Port); err != nil {
		return fmt.Errorf("failed to listen gRPC: %w", err)
	}

	// 创建 HTTP 监听器
	if err := s.listener.ListenHTTP(s.config.GRPC.GatewayPort); err != nil {
		return fmt.Errorf("failed to listen HTTP: %w", err)
	}

	// 创建 HTTP 路由
	router := http.NewServeMux()

	// 注册 Swagger UI
	swaggerHandler := httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
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
			s.log.Error("failed to read swagger.json", zap.Error(err))
			http.Error(w, "Swagger file not found", http.StatusNotFound)
			return
		}

		if _, er := w.Write(data); er != nil {
			s.log.Error("failed to write swagger.json", zap.Error(er))
		}
	})

	// API 路由
	router.Handle("/api/", s.mux)

	// 健康检查
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if _, er := w.Write([]byte(`{"status":"ok"}`)); er != nil {
			s.log.Error("failed to write health response", zap.Error(er))
		}
	})

	// 创建 HTTP 服务器
	s.httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.config.GRPC.GatewayPort),
		Handler: router,
	}

	// 启动 gRPC 服务器
	go func() {
		defer func() {
			if r := recover(); r != nil {
				s.log.Error("gRPC server panic recovered",
					zap.Any("panic", r),
				)
			}
		}()
		s.log.Info("Starting gRPC server",
			zap.String("address", fmt.Sprintf(":%d", s.config.GRPC.Port)),
		)
		if err := s.grpcServer.Serve(s.listener.GetGRPCListener()); err != nil {
			log.Error("failed to serve gRPC", zap.Error(err))
			os.Exit(1)
		}
	}()

	// 启动 HTTP 服务器
	go func() {
		defer func() {
			if r := recover(); r != nil {
				s.log.Error("HTTP server panic recovered",
					zap.Any("panic", r),
				)
			}
		}()
		s.log.Info("Starting HTTP Gateway server",
			zap.String("address", fmt.Sprintf(":%d", s.config.GRPC.GatewayPort)),
			zap.String("swagger_url", fmt.Sprintf("http://localhost:%d/swagger/", s.config.GRPC.GatewayPort)),
		)
		if err := s.httpServer.Serve(s.listener.GetHTTPListener()); err != nil && err != http.ErrServerClosed {
			s.log.Error("failed to serve HTTP", zap.Error(err))
			os.Exit(1)
		}
	}()

	// 等待停止信号
	<-s.stopCh
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
			s.log.Error("failed to shutdown HTTP server", zap.Error(err))
		}
	}

	// 优雅关闭 gRPC 服务器
	s.grpcServer.GracefulStop()
	return nil
}

// GetListener 获取监听器
func (s *Server) GetListener() listen.Listener {
	return s.listener
}
