//go:build wireinject
// +build wireinject

package wire

import (
	"testing"

	"github.com/google/wire"
	"go.uber.org/zap"

	"github.com/lazyjean/sla2/config"
	"github.com/lazyjean/sla2/pkg/logger"

	"github.com/lazyjean/sla2/internal/application/service"
	"github.com/lazyjean/sla2/internal/domain/repository"
	domainsecurity "github.com/lazyjean/sla2/internal/domain/security"

	"github.com/lazyjean/sla2/internal/infrastructure/cache/redis"
	"github.com/lazyjean/sla2/internal/infrastructure/listen"
	"github.com/lazyjean/sla2/internal/infrastructure/oauth"
	"github.com/lazyjean/sla2/internal/infrastructure/persistence/postgres"
	infrasecurity "github.com/lazyjean/sla2/internal/infrastructure/security"
	"github.com/lazyjean/sla2/internal/infrastructure/test"

	grpcserver "github.com/lazyjean/sla2/internal/transport/grpc"
	"github.com/lazyjean/sla2/internal/transport/http/ws/handler"

	adminconverter "github.com/lazyjean/sla2/internal/transport/grpc/admin/converter"
	courseconverter "github.com/lazyjean/sla2/internal/transport/grpc/course/converter"
	learningconverter "github.com/lazyjean/sla2/internal/transport/grpc/learning/converter"
	questionconverter "github.com/lazyjean/sla2/internal/transport/grpc/question/converter"
	vocabularyconverter "github.com/lazyjean/sla2/internal/transport/grpc/vocabulary/converter"
)

// 提供 Logger 实例
func ProvideLogger() *zap.Logger {
	// 使用项目中的标准 logger 实例
	return logger.Log
}

// 配置集
var configSet = wire.NewSet(
	config.GetConfig,
	wire.FieldsOf(new(*config.Config), "Database", "Redis", "JWT", "Apple", "RBAC", "Log"),
)

// 数据库集
var dbSet = wire.NewSet(
	postgres.NewDB,
)

var testDBSet = wire.NewSet(
	test.NewTestDB,
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
	postgres.NewQuestionRepository,
	postgres.NewHanCharRepository,
	postgres.NewMemoryUnitRepository,
)

// 转换器
var converterSet = wire.NewSet(
	courseconverter.NewCourseConverter,
	vocabularyconverter.NewVocabularyConverter,
	questionconverter.NewQuestionConverter,
	adminconverter.NewAdminConverter,
	learningconverter.NewLearningConverter,
)

// 服务集
var serviceSet = wire.NewSet(
	service.NewVocabularyService,
	service.NewLearningService,
	service.NewUserService,
	service.NewCourseService,
	provideAdminService,
	service.NewQuestionService,
	service.NewMemoryService,
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

var testAuthSet = wire.NewSet(
	test.NewMockAppleAuthService,
)

// 安全服务集
var securitySet = wire.NewSet(
	// 密码服务
	infrasecurity.NewBCryptPasswordService,

	// 令牌服务
	infrasecurity.NewJWTTokenService,
)

// WebSocket处理器集
var wsSet = wire.NewSet(
	handler.NewWebSocketHandler,
)

// gRPC服务器集
var grpcSet = wire.NewSet(
	listen.NewListener,
	grpcserver.NewGRPCServer,
)

var bufConnGrpcSet = wire.NewSet(
	test.NewTestListener,
	grpcserver.NewGRPCServer,
)

// 提供 PermissionHelper
func ProvidePermissionHelper(rbacConfig *config.RBACConfig) *domainsecurity.PermissionHelper {
	return &domainsecurity.PermissionHelper{}
}

// RBAC权限集
var rbacSet = wire.NewSet(
	domainsecurity.NewRBACProvider,
	wire.FieldsOf(new(*domainsecurity.RBACProvider), "PermissionHelper"),
)

// 测试安全服务集
var testSecuritySet = wire.NewSet(
	// 使用 mock token service
	test.NewMockTokenService,
	wire.Bind(new(domainsecurity.TokenService), new(*test.MockTokenService)),
	// 使用真实的密码服务
	infrasecurity.NewBCryptPasswordService,
)

// InitializeApp 初始化应用程序
func InitializeApp() (*Application, error) {
	wire.Build(
		NewApplication,
		logger.NewAppLogger,
		configSet,
		dbSet,
		repositorySet,
		serviceSet,
		authSet,
		securitySet,
		wsSet,
		grpcSet,
		rbacSet,
	)
	return nil, nil
}

func InitializeTestApp(t *testing.T) (*Application, error) {
	wire.Build(
		NewApplication,
		logger.NewAppLogger,
		configSet,
		testDBSet,
		repositorySet,
		serviceSet,
		testAuthSet,
		testSecuritySet,
		wsSet,
		bufConnGrpcSet,
		rbacSet,
	)
	return nil, nil
}
