package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

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
	// 初始化配置
	if err := config.InitConfig(); err != nil {
		fmt.Printf("Failed to initialize config: %v\n", err)
		os.Exit(1)
	}

	// 初始化日志
	cfg := config.GetConfig()
	if err := logger.InitLogger(&cfg.Log); err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}

	// 创建上下文
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 初始化应用程序
	app, err := wire.InitializeApp()
	if err != nil {
		logger.Log.Fatal("Failed to initialize application", zap.Error(err))
	}

	// 启动应用程序
	if err := app.Start(ctx); err != nil {
		logger.Log.Fatal("Failed to start application", zap.Error(err))
	}

	// 等待中断信号
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	// 优雅关闭
	if err := app.Stop(ctx); err != nil {
		logger.Log.Error("Failed to stop application", zap.Error(err))
	}
}
