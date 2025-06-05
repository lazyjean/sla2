//go:build wireinject
// +build wireinject

package wire

import (
	"testing"

	"github.com/lazyjean/sla2/internal/transport/grpc/admin"
	"github.com/lazyjean/sla2/internal/transport/grpc/course"
	"github.com/lazyjean/sla2/internal/transport/grpc/learning"
	"github.com/lazyjean/sla2/internal/transport/grpc/question"
	"github.com/lazyjean/sla2/internal/transport/grpc/user"
	"github.com/lazyjean/sla2/internal/transport/grpc/vocabulary"

	"github.com/lazyjean/sla2/internal/transport/grpc/middleware"

	"github.com/google/wire"
	"github.com/lazyjean/sla2/config"
	"github.com/lazyjean/sla2/pkg/logger"

	"github.com/lazyjean/sla2/internal/application/service"
	domainsecurity "github.com/lazyjean/sla2/internal/domain/security"
	domainservice "github.com/lazyjean/sla2/internal/domain/service"

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

	infravalidator "github.com/lazyjean/sla2/internal/infrastructure/validator"
)

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
	postgres.NewVocabularyRepository,
	postgres.NewCachedWordRepository,
	postgres.NewLearningRepository,
	postgres.NewUserRepository,
	postgres.NewCourseRepository,
	postgres.NewCourseSectionRepository,
	postgres.NewCourseSectionUnitRepository,
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
	domainservice.NewMemoryService,
	service.NewVocabularyService,
	service.NewLearningService,
	service.NewUserService,
	service.NewCourseService,
	service.NewQuestionService,
	service.NewAdminService,
)

// 认证集
var authSet = wire.NewSet(
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

// transport 服务集
var transportSet = wire.NewSet(
	user.NewUserService,
	question.NewQuestionService,
	course.NewCourseService,
	learning.NewLearningService,
	admin.NewAdminService,
	vocabulary.NewVocabularyService,
)

// gRPC服务器集
var grpcSet = wire.NewSet(
	listen.NewListener,
	grpcserver.NewGRPCServer,
	grpcserver.NewServer,
)

var bufConnGrpcSet = wire.NewSet(
	test.NewTestListener,
	grpcserver.NewGRPCServer,
	grpcserver.NewServer,
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

// 验证器集
var validatorSet = wire.NewSet(
	infravalidator.NewValidator,
)

// 测试验证器集
var testValidatorSet = wire.NewSet(
	infravalidator.NewValidator,
)

// InitializeApp 初始化应用程序
func InitializeApp() (*Application, error) {
	wire.Build(
		NewApplication,
		logger.NewAppLogger,
		middleware.NewMetrics,
		middleware.NewRegistry,
		configSet,
		dbSet,
		repositorySet,
		serviceSet,
		authSet,
		securitySet,
		wsSet,
		transportSet,
		converterSet,
		grpcSet,
		rbacSet,
		//validatorSet,
	)
	return nil, nil
}

func InitializeTestApp(t *testing.T) (*Application, error) {
	wire.Build(
		NewApplication,
		logger.NewAppLogger,
		middleware.NewMetrics,
		middleware.NewRegistry,
		configSet,
		testDBSet,
		repositorySet,
		serviceSet,
		testAuthSet,
		testSecuritySet,
		wsSet,
		transportSet,
		converterSet,
		bufConnGrpcSet,
		rbacSet,
		//testValidatorSet,
	)
	return nil, nil
}
