//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/google/wire"
	"github.com/lazyjean/sla2/config"
	"github.com/lazyjean/sla2/internal/application/service"
	"github.com/lazyjean/sla2/internal/domain/repository"
	domainsecurity "github.com/lazyjean/sla2/internal/domain/security"
	"github.com/lazyjean/sla2/internal/infrastructure/cache/redis"
	"github.com/lazyjean/sla2/internal/infrastructure/oauth"
	"github.com/lazyjean/sla2/internal/infrastructure/persistence/postgres"
	"github.com/lazyjean/sla2/internal/infrastructure/security"
	"github.com/lazyjean/sla2/internal/interfaces/grpc"
	"github.com/lazyjean/sla2/pkg/logger"
	"go.uber.org/zap"
)

// 提供 Logger 实例
func ProvideLogger() *zap.Logger {
	// 使用项目中的标准 logger 实例
	return logger.Log
}

// 配置集
var configSet = wire.NewSet(
	wire.FieldsOf(new(*config.Config), "Database", "Redis", "JWT", "Apple", "RBAC"),
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
)

// 服务集
var serviceSet = wire.NewSet(
	service.NewWordService,
	service.NewLearningService,
	service.NewUserService,
	service.NewCourseService,
	provideAdminService,
	service.NewQuestionService,
	service.NewQuestionTagService,
)

// provideAdminService 提供管理员服务
func provideAdminService(
	adminRepo repository.AdminRepository,
	passwordService domainsecurity.PasswordService,
	tokenService domainsecurity.TokenService,
	permissionHelper *domainsecurity.PermissionHelper,
) *service.AdminService {
	return service.NewAdminService(
		adminRepo,
		passwordService,
		tokenService,
		permissionHelper,
	)
}

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
		rbacSet,
		grpc.NewServer,
	)
	return nil, nil
}
