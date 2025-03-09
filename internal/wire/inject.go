package wire

import (
	"context"

	"github.com/lazyjean/sla2/config"
	"github.com/lazyjean/sla2/internal/domain/security"
	"github.com/lazyjean/sla2/internal/infrastructure/oauth"
	"github.com/lazyjean/sla2/internal/interfaces/grpc"
	"github.com/lazyjean/sla2/pkg/logger"
)

type App struct {
	grpcServer *grpc.Server
}

func NewApp(
	grpcServer *grpc.Server,
	cfg *config.Config,
	tokenService security.TokenService,
	appleAuthService *oauth.AppleAuthService,
) *App {
	return &App{
		grpcServer: grpcServer,
	}
}

func (a *App) StartGRPCServer() error {
	logger.Log.Info("Starting gRPC server and HTTP gateway")
	return a.grpcServer.Start()
}

func (a *App) Stop(ctx context.Context) error {
	// 关闭 gRPC 服务器
	return a.grpcServer.Stop()
}
