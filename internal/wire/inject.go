package wire

import (
	"context"
	"strconv"

	"github.com/lazyjean/sla2/config"
	"github.com/lazyjean/sla2/internal/domain/security"
	"github.com/lazyjean/sla2/internal/interfaces/grpc"
	"github.com/lazyjean/sla2/internal/infrastructure/oauth"
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

func (a *App) StartGRPCServer(port int) error {
	logger.Log.Info("gRPC server starting on port " + strconv.Itoa(port))
	return a.grpcServer.Start(port)
}

func (a *App) StartGRPCGateway(gatewayPort int, grpcPort int) error {
	logger.Log.Info("gRPC-Gateway server starting on port " + strconv.Itoa(gatewayPort))
	return a.grpcServer.StartGateway(gatewayPort, grpcPort)
}

func (a *App) Stop(ctx context.Context) error {
	// 关闭 gRPC 服务器
	a.grpcServer.Stop()
	return nil
}
