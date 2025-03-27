package wire

import (
	"context"

	"github.com/lazyjean/sla2/config"
	grpcserver "github.com/lazyjean/sla2/internal/interfaces/grpc"
)

// Application 应用程序结构体
type Application struct {
	config     *config.Config
	grpcServer *grpcserver.GRPCServer
}

// NewApplication 创建新的应用程序
func NewApplication(
	config *config.Config,
	grpcServer *grpcserver.GRPCServer,
) *Application {
	return &Application{
		config:     config,
		grpcServer: grpcServer,
	}
}

// Start 启动应用程序
func (a *Application) Start(ctx context.Context) error {
	return a.grpcServer.Start()
}

// Stop 停止应用程序
func (a *Application) Stop(ctx context.Context) error {
	return a.grpcServer.Stop()
}
