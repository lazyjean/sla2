package wire

import (
	"errors"
	"os"
	"os/signal"
	"syscall"

	"github.com/lazyjean/sla2/config"
	"github.com/lazyjean/sla2/internal/infrastructure/listen"
	grpcserver "github.com/lazyjean/sla2/internal/transport/grpc"
	"github.com/lazyjean/sla2/pkg/banner"
	"go.uber.org/zap"
)

// Application 应用程序结构体
type Application struct {
	config     *config.Config
	grpcServer *grpcserver.Server
	initLogger *zap.Logger
}

// NewApplication 创建新的应用程序
func NewApplication(
	config *config.Config,
	grpcServer *grpcserver.Server,
	logger *zap.Logger,
) *Application {
	return &Application{
		config:     config,
		grpcServer: grpcServer,
		initLogger: logger,
	}
}

// Start 启动应用程序
func (a *Application) Start() error {
	return a.grpcServer.Start()
}

// Stop 停止应用程序
func (a *Application) Stop() error {
	return a.grpcServer.Stop()
}

// GetGRPCListener 获取 gRPC 监听器
func (a *Application) GetGRPCListener() listen.Listener {
	return a.grpcServer.GetListener()
}

// GetInitLogger GetLogger 获取初始化 logger
func (a *Application) GetInitLogger() *zap.Logger {
	return a.initLogger
}

func (a *Application) Run() error {
	// 显示启动 banner
	banner.PrintBanner("1.0.0")

	// 初始化配置
	if err := config.InitConfig(); err != nil {
		return errors.Join(err, errors.New("failed to initialize config"))
	}

	// 启动应用程序
	if err := a.Start(); err != nil {
		return errors.Join(err, errors.New("failed to start application"))
	}

	// 等待中断信号
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	// 优雅关闭
	if err := a.Stop(); err != nil {
		return errors.Join(err, errors.New("failed to stop application"))
	}
	return nil
}

func Run() error {
	application, err := InitializeApp()
	if err != nil {
		return errors.Join(err, errors.New("failed to initialize application"))
	}
	if err := application.Run(); err != nil {
		return errors.Join(err, errors.New("failed to run application"))
	}
	return nil
}
