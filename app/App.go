package app

import (
	"context"
	"fmt"
	"github.com/lazyjean/sla2/config"
	"github.com/lazyjean/sla2/internal/wire"
	"github.com/lazyjean/sla2/pkg/banner"
	"github.com/lazyjean/sla2/pkg/logger"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
)

func Run(app *wire.Application, err error) {
	// 显示启动 banner
	banner.PrintBanner("1.0.0")

	// 初始化配置
	if err := config.InitConfig(); err != nil {
		fmt.Printf("Failed to initialize config: %v\n", err)
		os.Exit(1)
	}

	// 创建上下文
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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
