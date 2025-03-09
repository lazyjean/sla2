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

// @title        SLA2 API
// @version      1.0
// @description  SLA2 API Documentation

// @contact.name   LazyJean
// @contact.email  lazyjean@foxmail.com

// @host      localhost:9102
// @BasePath  /v1
// @schemes   http https

// @securityDefinitions.apikey  Bearer
// @in                         header
// @name                       Authorization
// @description               Bearer token for authentication

// @securityDefinitions.basic  BasicAuth
// @description               Basic authentication for Swagger UI

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

	// 启动服务器
	if err := app.StartGRPCServer(); err != nil {
		logger.Log.Fatal("Failed to start servers: " + err.Error())
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
