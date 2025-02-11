//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/google/wire"
	"github.com/lazyjean/sla2/config"
	"github.com/lazyjean/sla2/internal/application/service"
	"github.com/lazyjean/sla2/internal/infrastructure/cache/redis"
	"github.com/lazyjean/sla2/internal/infrastructure/persistence/postgres"
	"github.com/lazyjean/sla2/internal/interfaces/grpc"
	"github.com/lazyjean/sla2/internal/interfaces/http/handler"
	"github.com/lazyjean/sla2/pkg/auth"
)

// 配置集
var configSet = wire.NewSet(
	wire.FieldsOf(new(*config.Config), "Database", "Redis", "JWT", "Apple"),
)

// 数据库集
var dbSet = wire.NewSet(
	postgres.NewDB,
)

// redis仓库集
var cacheSet = wire.NewSet(
	redis.NewRedisCache,
)

// 仓储集
var repositorySet = wire.NewSet(
	postgres.NewWordRepository,
	postgres.NewCachedWordRepository,
	postgres.NewLearningRepository,
	postgres.NewUserRepository,
)

// 服务集
var serviceSet = wire.NewSet(
	service.NewWordService,
	service.NewLearningService,
	service.NewUserService,
)

// 处理器集
var handlerSet = wire.NewSet(
	handler.NewWordHandler,
	handler.NewUserHandler,
	handler.NewLearningHandler,
	handler.NewHealthHandler,
	handler.NewHandlers,
)

// 认证集
var authSet = wire.NewSet(
	auth.NewJWTConfig,
	auth.NewJWTService,
	wire.Bind(new(auth.JWTServicer), new(*auth.JWTService)),
)

func InitializeApp(cfg *config.Config) (*App, error) {
	wire.Build(
		NewApp,
		configSet,
		dbSet,
		repositorySet,
		serviceSet,
		handlerSet,
		authSet,
		cacheSet,
		grpc.NewServer,
	)
	return nil, nil
}
