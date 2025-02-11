// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

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

// Injectors from wire.go:

func InitializeApp(cfg *config.Config) (*App, error) {
	databaseConfig := &cfg.Database
	db, err := postgres.NewDB(databaseConfig)
	if err != nil {
		return nil, err
	}
	wordRepository := postgres.NewWordRepository(db)
	redisConfig := &cfg.Redis
	cache, err := redis.NewRedisCache(redisConfig)
	if err != nil {
		return nil, err
	}
	cachedWordRepository := postgres.NewCachedWordRepository(wordRepository, cache)
	wordService := service.NewWordService(cachedWordRepository)
	wordHandler := handler.NewWordHandler(wordService)
	userRepository := postgres.NewUserRepository(db)
	jwtConfig := auth.NewJWTConfig(cfg)
	jwtService := auth.NewJWTService(jwtConfig)
	userService := service.NewUserService(userRepository, jwtService)
	userHandler := handler.NewUserHandler(userService)
	learningRepository := postgres.NewLearningRepository(db)
	learningService := service.NewLearningService(learningRepository)
	learningHandler := handler.NewLearningHandler(learningService)
	healthHandler := handler.NewHealthHandler()
	handlers := handler.NewHandlers(wordHandler, userHandler, learningHandler, healthHandler)
	server := grpc.NewServer(userService, wordService, learningService, jwtService)
	app := NewApp(handlers, server, cfg)
	return app, nil
}

// wire.go:

// 配置集
var configSet = wire.NewSet(wire.FieldsOf(new(*config.Config), "Database", "Redis", "JWT", "Apple"))

// 数据库集
var dbSet = wire.NewSet(postgres.NewDB)

// redis仓库集
var cacheSet = wire.NewSet(redis.NewRedisCache)

// 仓储集
var repositorySet = wire.NewSet(postgres.NewWordRepository, postgres.NewCachedWordRepository, postgres.NewLearningRepository, postgres.NewUserRepository)

// 服务集
var serviceSet = wire.NewSet(service.NewWordService, service.NewLearningService, service.NewUserService)

// 处理器集
var handlerSet = wire.NewSet(handler.NewWordHandler, handler.NewUserHandler, handler.NewLearningHandler, handler.NewHealthHandler, handler.NewHandlers)

// 认证集
var authSet = wire.NewSet(auth.NewJWTConfig, auth.NewJWTService, wire.Bind(new(auth.JWTServicer), new(*auth.JWTService)))
