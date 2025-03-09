//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/google/wire"
	"github.com/lazyjean/sla2/config"
	"github.com/lazyjean/sla2/internal/application/service"
	"github.com/lazyjean/sla2/internal/infrastructure/ai"
	"github.com/lazyjean/sla2/internal/infrastructure/cache/redis"
	"github.com/lazyjean/sla2/internal/infrastructure/oauth"
	"github.com/lazyjean/sla2/internal/infrastructure/persistence/postgres"
	"github.com/lazyjean/sla2/internal/infrastructure/security"
	"github.com/lazyjean/sla2/internal/interfaces/grpc"
	"github.com/sirupsen/logrus"
)

// 提供 Logger 实例
func ProvideLogger() *logrus.Logger {
	return logrus.New()
}

// 提供 DeepSeekConfig 配置
func ProvideDeepSeekConfig(cfg *config.Config) *ai.DeepSeekConfig {
	return &ai.DeepSeekConfig{
		APIKey:      cfg.DeepSeek.APIKey,
		BaseURL:     cfg.DeepSeek.BaseURL,
		Timeout:     cfg.DeepSeek.Timeout,
		MaxRetries:  cfg.DeepSeek.MaxRetries,
		Temperature: cfg.DeepSeek.Temperature,
		MaxTokens:   cfg.DeepSeek.MaxTokens,
	}
}

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
	postgres.NewCourseRepository,
	postgres.NewCourseSectionRepository,
	postgres.NewAdminRepository,
	postgres.NewQuestionTagRepository,
	postgres.NewQuestionRepository,
	postgres.NewChatHistoryRepository,
)

// AI 服务集
var aiSet = wire.NewSet(
	service.NewAIService,
	ai.NewDeepSeekService,
	ProvideLogger,
	ProvideDeepSeekConfig,
)

// 服务集
var serviceSet = wire.NewSet(
	service.NewWordService,
	service.NewLearningService,
	service.NewUserService,
	service.NewCourseService,
	service.NewAdminService,
	service.NewQuestionService,
	service.NewQuestionTagService,
)

// 认证集
var authSet = wire.NewSet(
	oauth.NewAppleConfig,
	oauth.NewAppleAuthService,
)

// 安全服务集
var securitySet = wire.NewSet(
	// 密码服务
	security.NewBCryptPasswordService,

	// 令牌服务
	security.NewJWTTokenService,
)

func InitializeApp(cfg *config.Config) (*App, error) {
	wire.Build(
		NewApp,
		configSet,
		dbSet,
		repositorySet,
		serviceSet,
		authSet,
		cacheSet,
		securitySet,
		aiSet,
		grpc.NewServer,
	)
	return nil, nil
}
