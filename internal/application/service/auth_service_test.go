package service

import (
	"context"
	"testing"

	"github.com/lazyjean/sla2/internal/domain/entity"
	"github.com/lazyjean/sla2/internal/domain/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository 模拟用户仓储
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, username, password, email, nickname string) (*entity.User, error) {
	args := m.Called(ctx, username, password, email, nickname)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) FindByID(ctx context.Context, id uint) (*entity.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) FindByUsername(ctx context.Context, username string) (*entity.User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	args := m.Called(ctx, username)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	args := m.Called(ctx, email)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) Update(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

// MockJWTService 模拟 JWT 服务
type MockJWTService struct {
	mock.Mock
}

func (m *MockJWTService) GenerateToken(userID uint) (string, error) {
	args := m.Called(userID)
	return args.String(0), args.Error(1)
}

func (m *MockJWTService) ValidateToken(token string) (uint, error) {
	args := m.Called(token)
	return args.Get(0).(uint), args.Error(1)
}

func (m *MockJWTService) HashPassword(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

func (m *MockJWTService) ComparePasswords(hashedPassword, password string) bool {
	args := m.Called(hashedPassword, password)
	return args.Bool(0)
}

func (m *MockJWTService) GenerateRandomPassword() string {
	args := m.Called()
	return args.String(0)
}

func TestAuthService_Register(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockJWT := new(MockJWTService)
	service := NewAuthService(mockRepo, mockJWT)
	ctx := context.Background()

	t.Run("成功注册", func(t *testing.T) {
		req := &RegisterRequest{
			Username: "testuser",
			Password: "password123",
			Email:    "test@example.com",
			Nickname: "Test User",
		}

		hashedPassword := "hashed_password"
		mockJWT.On("HashPassword", req.Password).Return(hashedPassword, nil)

		mockRepo.On("ExistsByUsername", ctx, req.Username).Return(false, nil)
		mockRepo.On("ExistsByEmail", ctx, req.Email).Return(false, nil)

		newUser := &entity.User{
			ID:       1,
			Username: req.Username,
			Email:    req.Email,
			Nickname: req.Nickname,
		}

		mockRepo.On("Create", ctx, req.Username, hashedPassword, req.Email, req.Nickname).Return(newUser, nil)
		mockJWT.On("GenerateToken", newUser.ID).Return("test_token", nil)

		response, err := service.Register(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, req.Username, response.Username)
		assert.Equal(t, req.Email, response.Email)
		assert.Equal(t, req.Nickname, response.Nickname)
		assert.Equal(t, "test_token", response.Token)
	})

	t.Run("用户名已存在", func(t *testing.T) {
		req := &RegisterRequest{
			Username: "existinguser",
			Password: "password123",
			Email:    "test@example.com",
			Nickname: "Test User",
		}

		mockRepo.On("ExistsByUsername", ctx, req.Username).Return(true, nil)

		response, err := service.Register(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, errors.CodeUserAlreadyExists, err.(*errors.Error).Code)
	})

	t.Run("邮箱已存在", func(t *testing.T) {
		req := &RegisterRequest{
			Username: "newuser",
			Password: "password123",
			Email:    "existing@example.com",
			Nickname: "Test User",
		}

		mockRepo.On("ExistsByUsername", ctx, req.Username).Return(false, nil)
		mockRepo.On("ExistsByEmail", ctx, req.Email).Return(true, nil)

		response, err := service.Register(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, errors.CodeUserAlreadyExists, err.(*errors.Error).Code)
	})

	t.Run("必填字段为空", func(t *testing.T) {
		req := &RegisterRequest{
			Username: "",
			Password: "password123",
			Email:    "",
			Nickname: "Test User",
		}

		response, err := service.Register(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, errors.CodeInvalidArgument, err.(*errors.Error).Code)
	})
}

func TestAuthService_Login(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockJWT := new(MockJWTService)
	service := NewAuthService(mockRepo, mockJWT)
	ctx := context.Background()

	t.Run("成功登录", func(t *testing.T) {
		req := &LoginRequest{
			Account:  "testuser",
			Password: "password123",
		}

		user := &entity.User{
			ID:       1,
			Username: "testuser",
			Password: "hashed_password",
			Email:    "test@example.com",
			Nickname: "Test User",
		}

		mockRepo.On("FindByUsername", ctx, req.Account).Return(user, nil)
		mockJWT.On("ComparePasswords", user.Password, req.Password).Return(true)
		mockJWT.On("GenerateToken", user.ID).Return("test_token", nil)

		response, err := service.Login(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, user.Username, response.Username)
		assert.Equal(t, "test_token", response.Token)
	})

	t.Run("用户不存在", func(t *testing.T) {
		req := &LoginRequest{
			Account:  "nonexistent",
			Password: "password123",
		}

		mockRepo.On("FindByUsername", ctx, req.Account).Return(nil, errors.ErrNotFound)
		mockRepo.On("FindByEmail", ctx, req.Account).Return(nil, errors.ErrNotFound)

		response, err := service.Login(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, errors.CodeInvalidCredentials, err.(*errors.Error).Code)
	})

	t.Run("密码错误", func(t *testing.T) {
		req := &LoginRequest{
			Account:  "testuser",
			Password: "wrongpassword",
		}

		user := &entity.User{
			ID:       1,
			Username: "testuser",
			Password: "hashed_password",
		}

		mockRepo.On("FindByUsername", ctx, req.Account).Return(user, nil)
		mockJWT.On("ComparePasswords", user.Password, req.Password).Return(false)

		response, err := service.Login(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, errors.CodeInvalidCredentials, err.(*errors.Error).Code)
	})
}

func TestAuthService_GetUserByID(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockJWT := new(MockJWTService)
	service := NewAuthService(mockRepo, mockJWT)
	ctx := context.Background()

	t.Run("成功获取用户", func(t *testing.T) {
		userID := uint(1)
		user := &entity.User{
			ID:       userID,
			Username: "testuser",
			Email:    "test@example.com",
			Nickname: "Test User",
		}

		mockRepo.On("FindByID", ctx, userID).Return(user, nil)

		response, err := service.GetUserByID(ctx, userID)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, user.Username, response.Username)
		assert.Equal(t, user.Email, response.Email)
		assert.Equal(t, user.Nickname, response.Nickname)
	})

	t.Run("用户不存在", func(t *testing.T) {
		userID := uint(999)

		mockRepo.On("FindByID", ctx, userID).Return(nil, errors.ErrNotFound)

		response, err := service.GetUserByID(ctx, userID)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, errors.CodeUserNotFound, err.(*errors.Error).Code)
	})

	t.Run("无效用户ID", func(t *testing.T) {
		var userID uint = 0

		response, err := service.GetUserByID(ctx, userID)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, errors.CodeInvalidArgument, err.(*errors.Error).Code)
	})
}
