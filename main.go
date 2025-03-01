package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lazyjean/sla2/config"
	_ "github.com/lazyjean/sla2/docs" // 导入 swagger docs
	"github.com/lazyjean/sla2/internal/wire"
	"github.com/lazyjean/sla2/pkg/logger"
)

// @title        生词本 API
// @version      1.0
// @description  生词本服务 API 文档

// @contact.name   LazyJean
// @contact.email  lazyjean@foxmail.com

// @host      localhost:9000
// @BasePath  /api/v1
// @schemes   http

// @securityDefinitions.apikey  Bearer
// @in                         header
// @name                       Authorization
// @description               Bearer token for authentication

func main() {
	// 加载配置
	err := config.InitConfig()
	if err != nil {
		logger.Log.Fatal("Failed to load config: " + err.Error())
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

	// 优雅关闭
	go func() {
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
	}()

	// 启动 HTTP 和 gRPC 服务
	go func() {
		if err := app.StartHTTPServer(cfg.Server.Port); err != nil {
			logger.Log.Fatal("Failed to start HTTP server: " + err.Error())
		}
	}()

	go func() {
		if err := app.StartGRPCServer(cfg.GRPC.Port); err != nil {
			logger.Log.Fatal("Failed to start gRPC server: " + err.Error())
		}
	}()

	// 启动 gRPC Gateway
	go func() {
		if err := app.StartGRPCGateway(cfg.GRPC.GatewayPort, cfg.GRPC.Port); err != nil {
			logger.Log.Fatal("Failed to start gRPC gateway: " + err.Error())
		}
	}()

	// 等待所有服务启动完成
	select {}
}
