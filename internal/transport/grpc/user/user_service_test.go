package user

import (
	"context"
	"net"
	"net/http"
	"testing"

	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/validator"
	pb "github.com/lazyjean/sla2/api/proto/v1"
	"github.com/lazyjean/sla2/config"
	"github.com/lazyjean/sla2/internal/application/dto"
	"github.com/lazyjean/sla2/internal/application/service"
	"github.com/lazyjean/sla2/internal/domain/entity"
	domainOauth "github.com/lazyjean/sla2/internal/domain/oauth"
	domainSecurity "github.com/lazyjean/sla2/internal/domain/security"
	"github.com/lazyjean/sla2/internal/infrastructure/persistence/postgres"
	"github.com/lazyjean/sla2/internal/infrastructure/test"
	"github.com/lazyjean/sla2/pkg/logger"
	"github.com/lazyjean/sla2/pkg/utils"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
	"gorm.io/gorm"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

type mockTokenService struct {
	domainSecurity.TokenService
}

func (m *mockTokenService) GenerateToken(userID entity.UID, roles []string) (string, error) {
	return "test-token", nil
}

func (m *mockTokenService) GenerateRefreshToken(userID entity.UID, roles []string) (string, error) {
	return "test-refresh-token", nil
}

func (m *mockTokenService) ValidateToken(token string) (entity.UID, []string, error) {
	return entity.UID(1), []string{"user"}, nil
}

func (m *mockTokenService) ValidateRefreshToken(token string) (entity.UID, []string, error) {
	return entity.UID(1), []string{"user"}, nil
}

func (m *mockTokenService) ValidateTokenFromContext(ctx context.Context) (entity.UID, []string, error) {
	return entity.UID(1), []string{"user"}, nil
}

func (m *mockTokenService) ValidateTokenFromRequest(r *http.Request) (entity.UID, []string, error) {
	return entity.UID(1), []string{"user"}, nil
}

type mockPasswordService struct {
	domainSecurity.PasswordService
}

func (m *mockPasswordService) HashPassword(password string) (string, error) {
	return "hashed-" + password, nil
}

func (m *mockPasswordService) VerifyPassword(password, hash string) bool {
	return hash == "hashed-"+password
}

type mockAppleAuthService struct{}

func (m *mockAppleAuthService) AuthCodeWithApple(ctx context.Context, code string) (*domainOauth.AppleIDToken, error) {
	return &domainOauth.AppleIDToken{
		Subject: "test_apple_id",
		Email:   "test@example.com",
	}, nil
}

var _ domainOauth.AppleAuthService = (*mockAppleAuthService)(nil)

// setupTestDB 设置测试数据库
func setupTestDB(t *testing.T) (*gorm.DB, func()) {
	db, cleanup := test.SetupTestDB(t)

	// 清理测试数据
	cleanupData := func() {
		if err := db.Exec("DELETE FROM users").Error; err != nil {
			t.Logf("清理用户数据失败: %v", err)
		}
		cleanup()
	}

	return db, cleanupData
}

// setupTestServer 设置测试服务器
func setupTestServer(t *testing.T, log *zap.Logger) (*grpc.Server, func()) {
	lis = bufconn.Listen(bufSize)

	// 创建 gRPC 服务器，添加验证中间件
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpc_validator.UnaryServerInterceptor(),
		),
	)

	// 创建依赖
	db, cleanup := test.SetupTestDB(t)
	userRepo := postgres.NewUserRepository(db)
	passwordService := &mockPasswordService{}
	tokenService := &mockTokenService{}
	appleAuth := &mockAppleAuthService{}

	// 初始化全局 logger
	_ = logger.NewAppLogger(&config.LogConfig{
		Production: false,
	})

	// 创建服务
	userService := service.NewUserService(
		userRepo,
		tokenService,
		passwordService,
		appleAuth,
	)

	// 注册服务
	pb.RegisterUserServiceServer(grpcServer, NewUserService(userService))

	// 启动服务器
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			t.Logf("gRPC server error: %v", err)
		}
	}()

	// 返回清理函数
	cleanupFunc := func() {
		grpcServer.GracefulStop()
		cleanup()
	}

	return grpcServer, cleanupFunc
}

// TestUserService_Register 测试注册功能
func TestUserService_Register(t *testing.T) {
	// 创建测试用的 logger
	testLogger, _ := zap.NewDevelopment()
	ctx := context.Background()
	ctx = logger.WithContext(ctx, testLogger)

	// 设置测试服务器
	_, cleanup := setupTestServer(t, testLogger)
	defer cleanup()

	// 创建 gRPC 客户端
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(t, err)
	defer conn.Close()

	client := pb.NewUserServiceClient(conn)

	// 测试用例
	t.Run("success", func(t *testing.T) {
		req := &pb.RegisterRequest{
			Username: "testuser",
			Password: "Aa1@test123",
			Email:    "test@example.com",
			Nickname: "Test User",
		}

		resp, err := client.Register(ctx, req)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotEmpty(t, resp.User.Id)
		assert.NotEmpty(t, resp.Token)
		assert.NotEmpty(t, resp.RefreshToken)
	})

	// 测试无效的用户名
	t.Run("invalid username", func(t *testing.T) {
		req := &pb.RegisterRequest{
			Username: "te", // 太短
			Password: "Aa1@test123",
			Email:    "test@example.com",
			Nickname: "Test User",
		}

		_, err := client.Register(ctx, req)
		assert.Error(t, err)
		assert.Equal(t, codes.InvalidArgument, status.Code(err))
	})

	// 测试无效的密码
	t.Run("invalid password", func(t *testing.T) {
		req := &pb.RegisterRequest{
			Username: "testuser",
			Password: "password", // 不符合密码规则
			Email:    "test@example.com",
			Nickname: "Test User",
		}

		_, err := client.Register(ctx, req)
		assert.Error(t, err)
		assert.Equal(t, codes.InvalidArgument, status.Code(err))
	})

	// 测试无效的邮箱
	t.Run("invalid email", func(t *testing.T) {
		req := &pb.RegisterRequest{
			Username: "testuser",
			Password: "Aa1@test123",   // 包含大写字母、小写字母、数字和特殊字符
			Email:    "invalid-email", // 无效的邮箱格式
			Nickname: "Test User",
		}

		_, err := client.Register(ctx, req)
		assert.Error(t, err)
		assert.Equal(t, codes.InvalidArgument, status.Code(err))
	})
}

// TestUserService_Login 测试登录功能
func TestUserService_Login(t *testing.T) {
	// 创建测试用的 logger
	testLogger, _ := zap.NewDevelopment()
	ctx := context.Background()
	ctx = logger.WithContext(ctx, testLogger)

	// 设置测试服务器
	_, cleanup := setupTestServer(t, testLogger)
	defer cleanup()

	// 创建 gRPC 客户端
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(t, err)
	defer conn.Close()

	client := pb.NewUserServiceClient(conn)

	// 先注册一个用户
	registerReq := &pb.RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "Aa1@test123",
		Nickname: "Test User",
	}
	registerResp, err := client.Register(ctx, registerReq)
	assert.NoError(t, err)

	// 测试用例
	t.Run("success", func(t *testing.T) {
		req := &pb.LoginRequest{
			Account:  "test@example.com",
			Password: "Aa1@test123",
		}

		resp, err := client.Login(ctx, req)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, registerResp.User.Id, resp.User.Id)
		assert.Equal(t, "testuser", resp.User.Username)
		assert.Equal(t, "test@example.com", resp.User.Email)
		assert.Equal(t, "Test User", resp.User.Nickname)
		assert.NotEmpty(t, resp.Token)
		assert.NotEmpty(t, resp.RefreshToken)
	})

	// 测试无效的账号
	t.Run("invalid account", func(t *testing.T) {
		req := &pb.LoginRequest{
			Account:  "te", // 太短
			Password: "Aa1@test123",
		}

		_, err := client.Login(ctx, req)
		assert.Error(t, err)
		assert.Equal(t, codes.InvalidArgument, status.Code(err))
	})

	// 测试无效的密码
	t.Run("invalid password", func(t *testing.T) {
		req := &pb.LoginRequest{
			Account:  "test@example.com",
			Password: "pass", // 太短
		}

		_, err := client.Login(ctx, req)
		assert.Error(t, err)
		assert.Equal(t, codes.InvalidArgument, status.Code(err))
	})
}

// TestUserService_GetUserInfo 测试获取用户信息功能
func TestUserService_GetUserInfo(t *testing.T) {
	ctx := context.Background()

	// 初始化日志
	log := zap.NewExample()
	ctx = logger.WithContext(ctx, log)

	// 创建真实的依赖
	db, cleanup := setupTestDB(t)
	defer cleanup()

	userRepo := postgres.NewUserRepository(db)
	passwordService := &mockPasswordService{}
	tokenService := &mockTokenService{}
	appleAuth := &mockAppleAuthService{}

	// 创建 UserService 实例
	userService := service.NewUserService(
		userRepo,
		tokenService,
		passwordService,
		appleAuth,
	)

	// 先注册一个用户
	registerReq := &dto.RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
		Nickname: "Test User",
	}
	registerResp, err := userService.Register(ctx, registerReq)
	assert.NoError(t, err)

	// 设置用户ID到上下文
	ctx = utils.SetUserIDToContext(ctx, registerResp.UserID)

	// 调用获取用户信息方法
	user, err := userService.GetLoginUser(ctx)

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, registerResp.UserID, user.ID)
	assert.Equal(t, "testuser", user.Username)
	assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, "Test User", user.Nickname)
	assert.Equal(t, "active", user.Status)
}

// TestUserService_UpdateUserInfo 测试更新用户信息功能
func TestUserService_UpdateUserInfo(t *testing.T) {
	ctx := context.Background()

	// 初始化日志
	log := zap.NewExample()
	ctx = logger.WithContext(ctx, log)

	// 创建真实的依赖
	db, cleanup := setupTestDB(t)
	defer cleanup()

	userRepo := postgres.NewUserRepository(db)
	passwordService := &mockPasswordService{}
	tokenService := &mockTokenService{}
	appleAuth := &mockAppleAuthService{}

	// 创建 UserService 实例
	userService := service.NewUserService(
		userRepo,
		tokenService,
		passwordService,
		appleAuth,
	)

	// 先注册一个用户
	registerReq := &dto.RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
		Nickname: "Test User",
	}
	registerResp, err := userService.Register(ctx, registerReq)
	assert.NoError(t, err)

	// 设置用户ID到上下文
	ctx = utils.SetUserIDToContext(ctx, registerResp.UserID)

	// 创建更新请求
	req := &dto.UpdateUserRequest{
		Nickname: "New Nickname",
		Avatar:   "new_avatar.jpg",
	}

	// 调用更新用户信息方法
	err = userService.UpdateUser(ctx, req)
	assert.NoError(t, err)

	// 验证更新结果
	user, err := userService.GetLoginUser(ctx)
	assert.NoError(t, err)
	assert.Equal(t, "New Nickname", user.Nickname)
	assert.Equal(t, "new_avatar.jpg", user.Avatar)
}

// TestUserService_ChangePassword 测试修改密码功能
func TestUserService_ChangePassword(t *testing.T) {
	ctx := context.Background()

	// 初始化日志
	log := zap.NewExample()
	ctx = logger.WithContext(ctx, log)

	// 创建真实的依赖
	db, cleanup := setupTestDB(t)
	defer cleanup()

	userRepo := postgres.NewUserRepository(db)
	passwordService := &mockPasswordService{}
	tokenService := &mockTokenService{}
	appleAuth := &mockAppleAuthService{}

	// 创建 UserService 实例
	userService := service.NewUserService(
		userRepo,
		tokenService,
		passwordService,
		appleAuth,
	)

	// 先注册一个用户
	registerReq := &dto.RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "Aa1@test123",
		Nickname: "Test User",
	}
	registerResp, err := userService.Register(ctx, registerReq)
	assert.NoError(t, err)

	// 设置用户ID到上下文
	ctx = utils.SetUserIDToContext(ctx, registerResp.UserID)

	// 创建修改密码请求
	req := &dto.ChangePasswordRequest{
		OldPassword: "Aa1@test123",
		NewPassword: "Bb2@test456",
	}

	// 调用修改密码方法
	err = userService.ChangePassword(ctx, req)
	assert.NoError(t, err)

	// 验证新密码
	loginReq := &dto.LoginRequest{
		Account:  "test@example.com",
		Password: "Bb2@test456",
	}
	_, err = userService.Login(ctx, loginReq)
	assert.NoError(t, err)
}

// TestUserService_AppleLogin 测试苹果登录功能
func TestUserService_AppleLogin(t *testing.T) {
	ctx := context.Background()

	// 初始化日志
	log := zap.NewExample()
	ctx = logger.WithContext(ctx, log)

	// 创建真实的依赖
	db, cleanup := setupTestDB(t)
	defer cleanup()

	userRepo := postgres.NewUserRepository(db)
	passwordService := &mockPasswordService{}
	tokenService := &mockTokenService{}
	appleAuth := &mockAppleAuthService{}

	// 创建 UserService 实例
	userService := service.NewUserService(
		userRepo,
		tokenService,
		passwordService,
		appleAuth,
	)

	// 创建苹果登录请求
	req := &dto.AppleLoginRequest{
		AuthorizationCode: "test_auth_code",
		UserIdentifier:    "test_apple_id",
	}

	// 调用苹果登录方法
	resp, err := userService.AppleLogin(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotZero(t, resp.UserID)
	assert.NotEmpty(t, resp.Token)
	assert.NotEmpty(t, resp.RefreshToken)
	assert.True(t, resp.IsNewUser)

	// 再次登录，应该返回已存在的用户
	resp2, err := userService.AppleLogin(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, resp2)
	assert.Equal(t, resp.UserID, resp2.UserID)
	assert.NotEmpty(t, resp2.Token)
	assert.NotEmpty(t, resp2.RefreshToken)
	assert.False(t, resp2.IsNewUser)
}

// TestUserService_RefreshToken 测试刷新令牌功能
func TestUserService_RefreshToken(t *testing.T) {
	ctx := context.Background()

	// 初始化日志
	log := zap.NewExample()
	ctx = logger.WithContext(ctx, log)

	// 创建真实的依赖
	db, cleanup := setupTestDB(t)
	defer cleanup()

	userRepo := postgres.NewUserRepository(db)
	passwordService := &mockPasswordService{}
	tokenService := &mockTokenService{}
	appleAuth := &mockAppleAuthService{}

	// 创建 UserService 实例
	userService := service.NewUserService(
		userRepo,
		tokenService,
		passwordService,
		appleAuth,
	)

	// 先注册一个用户
	registerReq := &dto.RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
		Nickname: "Test User",
	}
	registerResp, err := userService.Register(ctx, registerReq)
	assert.NoError(t, err)
	assert.NotNil(t, registerResp)

	// 登录获取 refresh token
	loginReq := &dto.LoginRequest{
		Account:  "test@example.com",
		Password: "password123",
	}
	loginResp, err := userService.Login(ctx, loginReq)
	assert.NoError(t, err)

	// 创建刷新令牌请求
	req := &dto.RefreshTokenRequest{
		RefreshToken: loginResp.RefreshToken,
	}

	// 调用刷新令牌方法
	resp, err := userService.RefreshToken(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.AccessToken)
	assert.NotEmpty(t, resp.RefreshToken)
}

// TestUserService_Logout 测试登出功能
func TestUserService_Logout(t *testing.T) {
	ctx := context.Background()

	// 初始化日志
	log := zap.NewExample()
	ctx = logger.WithContext(ctx, log)

	// 创建真实的依赖
	db, cleanup := setupTestDB(t)
	defer cleanup()

	userRepo := postgres.NewUserRepository(db)
	passwordService := &mockPasswordService{}
	tokenService := &mockTokenService{}
	appleAuth := &mockAppleAuthService{}

	// 创建 UserService 实例
	userService := service.NewUserService(
		userRepo,
		tokenService,
		passwordService,
		appleAuth,
	)

	// 先注册一个用户
	registerReq := &dto.RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
		Nickname: "Test User",
	}
	registerResp, err := userService.Register(ctx, registerReq)
	assert.NoError(t, err)

	// 设置用户ID到上下文
	ctx = utils.SetUserIDToContext(ctx, registerResp.UserID)

	// 创建登出请求
	req := &dto.LogoutRequest{}

	// 调用登出方法
	resp, err := userService.Logout(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.IsType(t, &dto.LogoutResponse{}, resp)
}
