// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package wire

import (
	"github.com/google/wire"
	"github.com/lazyjean/sla2/config"
	"github.com/lazyjean/sla2/internal/application/service"
	security2 "github.com/lazyjean/sla2/internal/domain/security"
	"github.com/lazyjean/sla2/internal/infrastructure/cache/redis"
	"github.com/lazyjean/sla2/internal/infrastructure/oauth"
	"github.com/lazyjean/sla2/internal/infrastructure/persistence/postgres"
	"github.com/lazyjean/sla2/internal/infrastructure/security"
	"github.com/lazyjean/sla2/internal/interfaces/grpc"
	"github.com/lazyjean/sla2/internal/interfaces/middleware"
	"github.com/lazyjean/sla2/pkg/logger"
	"go.uber.org/zap"
)

// Injectors from wire.go:

func InitializeApp(cfg *config.Config) (*App, error) {
	databaseConfig := &cfg.Database
	db, err := postgres.NewDB(databaseConfig)
	if err != nil {
		return nil, err
	}
	adminRepository := postgres.NewAdminRepository(db)
	passwordService := security.NewBCryptPasswordService()
	tokenService := security.NewJWTTokenService(cfg)
	adminService := service.NewAdminService(adminRepository, passwordService, tokenService)
	userRepository := postgres.NewUserRepository(db)
	appleConfig := oauth.NewAppleConfig(cfg)
	appleAuthService := oauth.NewAppleAuthService(appleConfig)
	userService := service.NewUserService(userRepository, tokenService, passwordService, appleAuthService)
	wordRepository := postgres.NewWordRepository(db)
	redisConfig := &cfg.Redis
	cache, err := redis.NewRedisCache(redisConfig)
	if err != nil {
		return nil, err
	}
	cachedWordRepository := postgres.NewCachedWordRepository(wordRepository, cache)
	wordService := service.NewWordService(cachedWordRepository)
	learningRepository := postgres.NewLearningRepository(db)
	learningService := service.NewLearningService(learningRepository)
	courseRepository := postgres.NewCourseRepository(db)
	courseSectionRepository := postgres.NewCourseSectionRepository(db)
	courseService := service.NewCourseService(courseRepository, courseSectionRepository)
	questionRepository := postgres.NewQuestionRepository(db)
	questionService := service.NewQuestionService(questionRepository)
	questionTagRepository := postgres.NewQuestionTagRepository(db)
	questionTagService := service.NewQuestionTagService(questionTagRepository)
	rbacConfig := &cfg.RBAC
	rbacProvider, err := security2.NewRBACProvider(db, rbacConfig)
	if err != nil {
		return nil, err
	}
	permissionHelper := providePermissionHelper(rbacProvider)
	rbacInterceptor := middleware.NewRBACInterceptor(permissionHelper)
	server := grpc.NewServer(adminService, userService, wordService, learningService, courseService, questionService, questionTagService, tokenService, rbacInterceptor, cfg)
	app := NewApp(server, cfg, tokenService, appleAuthService)
	return app, nil
}

// wire.go:

// 提供 Logger 实例
func ProvideLogger() *zap.Logger {

	return logger.Log
}

// 配置集
var configSet = wire.NewSet(wire.FieldsOf(new(*config.Config), "Database", "Redis", "JWT", "Apple", "RBAC"))

// 数据库集
var dbSet = wire.NewSet(postgres.NewDB)

// redis仓库集
var cacheSet = wire.NewSet(redis.NewRedisCache)

// 仓储集
var repositorySet = wire.NewSet(postgres.NewWordRepository, postgres.NewCachedWordRepository, postgres.NewLearningRepository, postgres.NewUserRepository, postgres.NewCourseRepository, postgres.NewCourseSectionRepository, postgres.NewAdminRepository, postgres.NewQuestionTagRepository, postgres.NewQuestionRepository)

// 服务集
var serviceSet = wire.NewSet(service.NewWordService, service.NewLearningService, service.NewUserService, service.NewCourseService, service.NewAdminService, service.NewQuestionService, service.NewQuestionTagService)

// 认证集
var authSet = wire.NewSet(oauth.NewAppleConfig, oauth.NewAppleAuthService)

// 安全服务集
var securitySet = wire.NewSet(security.NewBCryptPasswordService, security.NewJWTTokenService)
