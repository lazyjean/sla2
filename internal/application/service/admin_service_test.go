package service

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/lazyjean/sla2/internal/application/dto"
	"github.com/lazyjean/sla2/internal/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAdminRepository 模拟管理员仓库
type MockAdminRepository struct {
	mock.Mock
}

func (m *MockAdminRepository) Create(ctx context.Context, admin *entity.Admin) error {
	args := m.Called(ctx, admin)
	return args.Error(0)
}

func (m *MockAdminRepository) FindByID(ctx context.Context, id entity.UID) (*entity.Admin, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Admin), args.Error(1)
}

func (m *MockAdminRepository) FindByUsername(ctx context.Context, username string) (*entity.Admin, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Admin), args.Error(1)
}

func (m *MockAdminRepository) Update(ctx context.Context, admin *entity.Admin) error {
	args := m.Called(ctx, admin)
	return args.Error(0)
}

func (m *MockAdminRepository) IsSystemInitialized(ctx context.Context) (bool, error) {
	args := m.Called(ctx)
	return args.Bool(0), args.Error(1)
}

func (m *MockAdminRepository) IsInitialized(ctx context.Context) (bool, error) {
	args := m.Called(ctx)
	return args.Bool(0), args.Error(1)
}

// MockPasswordService 模拟密码服务
type MockPasswordService struct {
	mock.Mock
}

func (m *MockPasswordService) HashPassword(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

func (m *MockPasswordService) VerifyPassword(password, hashedPassword string) bool {
	args := m.Called(password, hashedPassword)
	return args.Bool(0)
}

// MockTokenService 模拟令牌服务
type MockTokenService struct {
	mock.Mock
}

func (m *MockTokenService) GenerateToken(userID entity.UID, roles []string) (string, error) {
	args := m.Called(userID, roles)
	return args.String(0), args.Error(1)
}

func (m *MockTokenService) GenerateRefreshToken(userID entity.UID, roles []string) (string, error) {
	args := m.Called(userID, roles)
	return args.String(0), args.Error(1)
}

func (m *MockTokenService) ValidateRefreshToken(token string) (entity.UID, []string, error) {
	args := m.Called(token)
	return args.Get(0).(entity.UID), args.Get(1).([]string), args.Error(2)
}

func (m *MockTokenService) ValidateToken(token string) (entity.UID, []string, error) {
	args := m.Called(token)
	return args.Get(0).(entity.UID), args.Get(1).([]string), args.Error(2)
}

func (m *MockTokenService) ValidateTokenFromContext(ctx context.Context) (entity.UID, []string, error) {
	args := m.Called(ctx)
	return args.Get(0).(entity.UID), args.Get(1).([]string), args.Error(2)
}

func (m *MockTokenService) ValidateTokenFromRequest(r *http.Request) (entity.UID, []string, error) {
	args := m.Called(r)
	return args.Get(0).(entity.UID), args.Get(1).([]string), args.Error(2)
}

func TestAdminService_InitializeSystem(t *testing.T) {
	// 准备测试数据
	ctx := context.Background()
	req := &dto.InitializeSystemRequest{
		Username: "admin",
		Password: "password123",
	}

	// 创建模拟对象
	mockRepo := new(MockAdminRepository)
	mockPasswordService := new(MockPasswordService)
	mockTokenService := new(MockTokenService)

	// 设置模拟行为
	mockRepo.On("IsInitialized", ctx).Return(false, nil)
	mockPasswordService.On("HashPassword", req.Password).Return("hashed_password", nil)
	mockRepo.On("Create", ctx, mock.MatchedBy(func(admin *entity.Admin) bool {
		return admin.Username == "admin" &&
			admin.Password == "hashed_password" &&
			admin.Nickname == "admin" &&
			len(admin.Roles) == 1 &&
			admin.Roles[0] == "admin"
	})).Return(nil)
	mockTokenService.On("GenerateToken", mock.Anything, []string{"admin"}).Return("access_token", nil)
	mockTokenService.On("GenerateRefreshToken", mock.Anything, []string{"admin"}).Return("refresh_token", nil)

	// 创建服务实例
	service := NewAdminService(mockRepo, mockPasswordService, mockTokenService)

	// 执行测试
	resp, err := service.InitializeSystem(ctx, req)

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "admin", resp.Admin.Username)
	assert.Equal(t, "access_token", resp.AccessToken)
	assert.Equal(t, "refresh_token", resp.RefreshToken)

	// 验证模拟对象是否按预期被调用
	mockRepo.AssertExpectations(t)
	mockPasswordService.AssertExpectations(t)
	mockTokenService.AssertExpectations(t)
}

func TestAdminService_Login(t *testing.T) {
	// 准备测试数据
	ctx := context.Background()
	req := &dto.AdminLoginRequest{
		Username: "admin",
		Password: "password123",
	}

	// 创建模拟对象
	mockRepo := new(MockAdminRepository)
	mockPasswordService := new(MockPasswordService)
	mockTokenService := new(MockTokenService)

	// 创建测试用的管理员实体
	admin := &entity.Admin{
		ID:        entity.UID(1),
		Username:  "admin",
		Password:  "hashed_password",
		Nickname:  "Admin",
		Roles:     []string{"admin"},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 设置模拟行为
	mockRepo.On("FindByUsername", ctx, req.Username).Return(admin, nil)
	mockPasswordService.On("VerifyPassword", req.Password, admin.Password).Return(true)
	mockTokenService.On("GenerateToken", admin.ID, admin.Roles).Return("access_token", nil)
	mockTokenService.On("GenerateRefreshToken", admin.ID, admin.Roles).Return("refresh_token", nil)

	// 创建服务实例
	service := NewAdminService(mockRepo, mockPasswordService, mockTokenService)

	// 执行测试
	resp, err := service.Login(ctx, req)

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, admin.ID, resp.Admin.ID)
	assert.Equal(t, "access_token", resp.AccessToken)
	assert.Equal(t, "refresh_token", resp.RefreshToken)

	// 验证模拟对象是否按预期被调用
	mockRepo.AssertExpectations(t)
	mockPasswordService.AssertExpectations(t)
	mockTokenService.AssertExpectations(t)
}

func TestAdminService_RefreshToken(t *testing.T) {
	// 准备测试数据
	ctx := context.Background()
	req := &dto.RefreshTokenRequest{
		RefreshToken: "refresh_token",
	}

	// 创建模拟对象
	mockRepo := new(MockAdminRepository)
	mockPasswordService := new(MockPasswordService)
	mockTokenService := new(MockTokenService)

	// 创建测试用的管理员实体
	admin := &entity.Admin{
		ID:        entity.UID(1),
		Username:  "admin",
		Password:  "hashed_password",
		Nickname:  "Admin",
		Roles:     []string{"admin"},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 添加用户ID和角色到上下文
	ctx = WithUserID(ctx, admin.ID)
	ctx = WithRoles(ctx, admin.Roles)

	// 设置模拟行为
	mockTokenService.On("ValidateRefreshToken", req.RefreshToken).Return(admin.ID, admin.Roles, nil)
	mockRepo.On("FindByID", ctx, admin.ID).Return(admin, nil)
	mockTokenService.On("GenerateToken", admin.ID, admin.Roles).Return("new_access_token", nil)
	mockTokenService.On("GenerateRefreshToken", admin.ID, admin.Roles).Return("new_refresh_token", nil)

	// 创建服务实例
	service := NewAdminService(mockRepo, mockPasswordService, mockTokenService)

	// 执行测试
	resp, err := service.RefreshToken(ctx, req)

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "new_access_token", resp.AccessToken)
	assert.Equal(t, "new_refresh_token", resp.RefreshToken)

	// 验证模拟对象是否按预期被调用
	mockRepo.AssertExpectations(t)
	mockTokenService.AssertExpectations(t)
}

func TestAdminService_GetCurrentAdminInfo(t *testing.T) {
	// 准备测试数据
	ctx := context.Background()

	// 创建模拟对象
	mockRepo := new(MockAdminRepository)
	mockPasswordService := new(MockPasswordService)
	mockTokenService := new(MockTokenService)

	// 创建测试用的管理员实体
	admin := &entity.Admin{
		ID:        entity.UID(1),
		Username:  "admin",
		Password:  "hashed_password",
		Nickname:  "Admin",
		Roles:     []string{"admin"},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 添加用户ID和角色到上下文
	ctx = WithUserID(ctx, admin.ID)
	ctx = WithRoles(ctx, admin.Roles)

	// 设置模拟行为
	mockRepo.On("FindByID", ctx, admin.ID).Return(admin, nil)

	// 创建服务实例
	service := NewAdminService(mockRepo, mockPasswordService, mockTokenService)

	// 执行测试
	resp, err := service.GetCurrentAdminInfo(ctx)

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, admin.ID, resp.ID)
	assert.Equal(t, admin.Username, resp.Username)
	assert.Equal(t, admin.Nickname, resp.Nickname)
	assert.Equal(t, admin.Roles, resp.Roles)

	// 验证模拟对象是否按预期被调用
	mockRepo.AssertExpectations(t)
}
