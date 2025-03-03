package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lazyjean/sla2/config"
	"github.com/lazyjean/sla2/internal/wire"
	"github.com/lazyjean/sla2/pkg/logger"
	"go.uber.org/zap"
)

// @title        生词本 API
// @version      1.0
// @description  生词本服务 API 文档 (gRPC-Gateway)

// @contact.name   LazyJean
// @contact.email  lazyjean@foxmail.com

// @host      localhost:8080
// @BasePath  /api/v1
// @schemes   http https

// @securityDefinitions.apikey  Bearer
// @in                         header
// @name                       Authorization
// @description               Bearer token for authentication

func main() {
	// 加载配置
	err := config.InitConfig()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}
	cfg := config.GetConfig()

	// 初始化日志
	err = logger.InitLogger(&cfg.Log)
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Log.Sync()

	// 初始化应用
	app, err := wire.InitializeApp(cfg)
	if err != nil {
		logger.Log.Fatal("Failed to initialize app: " + err.Error())
	}

	// 创建一个通道用于等待服务启动完成
	started := make(chan struct{})

	// 启动 gRPC 服务
	go func() {
		if err := app.StartGRPCServer(cfg.GRPC.Port); err != nil {
			logger.Log.Fatal("Failed to start gRPC server: " + err.Error())
		}
		started <- struct{}{}
	}()

	// 启动 gRPC Gateway
	go func() {
		if err := app.StartGRPCGateway(cfg.GRPC.GatewayPort, cfg.GRPC.Port); err != nil {
			logger.Log.Fatal("Failed to start gRPC gateway: " + err.Error())
		}
		started <- struct{}{}
	}()

	// 等待所有服务启动完成
	for i := 0; i < 2; i++ {
		<-started
	}

	logger.Log.Info("All servers started successfully")
	logger.Log.Info("API Documentation",
		zap.String("swagger_url", fmt.Sprintf("http://localhost:%d/swagger/", cfg.GRPC.GatewayPort)),
		zap.String("gateway_url", fmt.Sprintf("http://localhost:%d", cfg.GRPC.GatewayPort)),
		zap.String("grpc_port", fmt.Sprintf(":%d", cfg.GRPC.Port)),
		zap.String("swagger_username", "admin"),
		zap.String("swagger_password", "swagger"),
	)

	// 监听系统信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Log.Info("Shutting down servers...")

	// 创建一个5秒的上下文用于超时控制
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 关闭应用
	if err := app.Stop(ctx); err != nil {
		logger.Log.Fatal("Server forced to shutdown: " + err.Error())
	}

	logger.Log.Info("Servers exiting")
}
