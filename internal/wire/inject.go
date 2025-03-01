package wire

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lazyjean/sla2/config"
	"github.com/lazyjean/sla2/internal/interfaces/grpc"
	"github.com/lazyjean/sla2/internal/interfaces/http/handler"
	"github.com/lazyjean/sla2/internal/interfaces/http/routes"
	"github.com/lazyjean/sla2/pkg/auth"
	"github.com/lazyjean/sla2/pkg/logger"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type App struct {
	router     *gin.Engine
	handlers   *handler.Handlers
	grpcServer *grpc.Server
	httpServer *http.Server
}

func NewApp(handlers *handler.Handlers, grpcServer *grpc.Server, cfg *config.Config, jwtService *auth.JWTService) *App {
	// 创建 gin router
	gin.DefaultWriter = logger.NewGinLogger()
	gin.DefaultErrorWriter = logger.NewGinLogger()
	gin.SetMode(cfg.Server.Mode)

	r := gin.New()
	r.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		Output: logger.NewGinLogger(),
	}))
	r.Use(gin.Recovery())

	// 注册 Swagger 路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 注册业务路由
	routes.SetupRoutes(r, handlers, jwtService)

	// 注册用户路由
	r.POST("/api/v1/register", handlers.UserHandler.Register)
	r.POST("/api/v1/login", handlers.UserHandler.Login)

	return &App{
		router:     r,
		handlers:   handlers,
		grpcServer: grpcServer,
	}
}

func (a *App) StartHTTPServer(port string) error {
	a.httpServer = &http.Server{
		Addr:    ":" + port,
		Handler: a.router,
	}

	logger.Log.Info("HTTP server starting on port " + port)
	return a.httpServer.ListenAndServe()
}

func (a *App) StartGRPCServer(port int) error {
	logger.Log.Info("gRPC server starting on port " + strconv.Itoa(port))
	return a.grpcServer.Start(port)
}

func (a *App) StartGRPCGateway(gatewayPort int, grpcPort int) error {
	logger.Log.Info("gRPC-Gateway server starting on port " + strconv.Itoa(gatewayPort))
	return a.grpcServer.StartGateway(gatewayPort, grpcPort)
}

func (a *App) Stop(ctx context.Context) error {
	// 关闭 HTTP 服务器
	if err := a.httpServer.Shutdown(ctx); err != nil {
		return err
	}

	// 关闭 gRPC 服务器
	a.grpcServer.Stop()
	return nil
}
