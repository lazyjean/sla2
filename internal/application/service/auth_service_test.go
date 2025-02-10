package service

import (
	"context"
	"testing"
	"time"

	"github.com/lazyjean/sla2/internal/application/dto"
	"github.com/lazyjean/sla2/internal/domain/entity"
	"github.com/lazyjean/sla2/internal/domain/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

// MockUserRepository 模拟用户仓储
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
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

func (m *MockUserRepository) FindByPhone(ctx context.Context, phone string) (*entity.User, error) {
	args := m.Called(ctx, phone)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) FindByAccount(ctx context.Context, account string) (*entity.User, error) {
	args := m.Called(ctx, account)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func TestAuthService_Register(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewAuthService(mockRepo)
	ctx := context.Background()

	t.Run("成功注册", func(t *testing.T) {
		req := &dto.RegisterDTO{
			Username: "testuser",
			Password: "password123",
			Email:    "test@example.com",
			Phone:    "13800138000",
		}

		// 模拟用户名不存在
		mockRepo.On("FindByUsername", ctx, req.Username).Return(nil, errors.ErrNotFound)
		// 模拟邮箱不存在
		mockRepo.On("FindByEmail", ctx, req.Email).Return(nil, errors.ErrNotFound)
		// 模拟手机号不存在
		mockRepo.On("FindByPhone", ctx, req.Phone).Return(nil, errors.ErrNotFound)
		// 模拟创建用户成功
		mockRepo.On("Create", ctx, mock.AnythingOfType("*entity.User")).Return(nil)

		user, err := service.Register(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, req.Username, user.Username)
		assert.Equal(t, req.Email, user.Email)
		assert.Equal(t, req.Phone, user.Phone)
	})

	t.Run("用户名已存在", func(t *testing.T) {
		req := &dto.RegisterDTO{
			Username: "existinguser",
			Password: "password123",
		}

		existingUser := &entity.User{
			Username: req.Username,
		}

		mockRepo.On("FindByUsername", ctx, req.Username).Return(existingUser, nil)

		user, err := service.Register(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, errors.CodeUserAlreadyExists, err.(*errors.Error).Code)
	})

	t.Run("只有用户名有效", func(t *testing.T) {
		req := &dto.RegisterDTO{
			Username: "testuser",
			Password: "password123",
			Email:    "",
			Phone:    "",
		}

		// 模拟用户名不存在
		mockRepo.On("FindByUsername", ctx, req.Username).Return(nil, errors.ErrNotFound)
		// 模拟创建用户成功
		mockRepo.On("Create", ctx, mock.AnythingOfType("*entity.User")).Return(nil)

		user, err := service.Register(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, req.Username, user.Username)
	})

	t.Run("只有邮箱有效", func(t *testing.T) {
		req := &dto.RegisterDTO{
			Username: "",
			Password: "password123",
			Email:    "test@example.com",
			Phone:    "",
		}

		// 模拟空用户名检查
		mockRepo.On("FindByUsername", ctx, req.Username).Return(nil, errors.ErrNotFound)
		// 模拟邮箱不存在
		mockRepo.On("FindByEmail", ctx, req.Email).Return(nil, errors.ErrNotFound)
		// 模拟创建用户成功
		mockRepo.On("Create", ctx, mock.AnythingOfType("*entity.User")).Return(nil)

		user, err := service.Register(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, req.Email, user.Email)
	})

	t.Run("只有手机号有效", func(t *testing.T) {
		req := &dto.RegisterDTO{
			Username: "",
			Password: "password123",
			Email:    "",
			Phone:    "13800138000",
		}

		// 模拟空用户名检查
		mockRepo.On("FindByUsername", ctx, req.Username).Return(nil, errors.ErrNotFound)
		// 模拟空邮箱检查
		mockRepo.On("FindByEmail", ctx, req.Email).Return(nil, errors.ErrNotFound)
		// 模拟手机号不存在
		mockRepo.On("FindByPhone", ctx, req.Phone).Return(nil, errors.ErrNotFound)
		// 模拟创建用户成功
		mockRepo.On("Create", ctx, mock.AnythingOfType("*entity.User")).Return(nil)

		user, err := service.Register(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, req.Phone, user.Phone)
	})

	t.Run("三个字段都为空", func(t *testing.T) {
		req := &dto.RegisterDTO{
			Username: "",
			Password: "password123",
			Email:    "",
			Phone:    "",
		}

		user, err := service.Register(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, errors.CodeInvalidArgument, err.(*errors.Error).Code)

		// 验证 mock 方法没有被调用
		mockRepo.AssertNotCalled(t, "FindByUsername")
		mockRepo.AssertNotCalled(t, "FindByEmail")
		mockRepo.AssertNotCalled(t, "FindByPhone")
		mockRepo.AssertNotCalled(t, "Create")
	})

	t.Run("密码为空", func(t *testing.T) {
		req := &dto.RegisterDTO{
			Username: "testuser",
			Password: "",
			Email:    "test@example.com",
			Phone:    "13800138000",
		}

		user, err := service.Register(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, errors.CodeInvalidArgument, err.(*errors.Error).Code)
	})
}

func TestAuthService_Login(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewAuthService(mockRepo)
	ctx := context.Background()

	t.Run("成功登录", func(t *testing.T) {
		req := &dto.LoginDTO{
			Account:  "testuser",
			Password: "password123",
		}

		// 使用 bcrypt 生成正确的密码哈希
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
		require.NoError(t, err)

		user := &entity.User{
			ID:        1,
			Username:  "testuser",
			Password:  string(hashedPassword),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mockRepo.On("FindByAccount", ctx, req.Account).Return(user, nil)

		userDTO, err := service.Login(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, userDTO)
		assert.Equal(t, user.ID, userDTO.ID)
		assert.Equal(t, user.Username, userDTO.Username)
	})

	t.Run("账号不存在", func(t *testing.T) {
		req := &dto.LoginDTO{
			Account:  "nonexistent",
			Password: "password123",
		}

		mockRepo.On("FindByAccount", ctx, req.Account).Return(nil, errors.ErrNotFound)

		userDTO, err := service.Login(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, userDTO)
		assert.Equal(t, errors.CodeInvalidCredentials, err.(*errors.Error).Code)
	})

	t.Run("密码错误", func(t *testing.T) {
		req := &dto.LoginDTO{
			Account:  "testuser",
			Password: "wrongpassword",
		}

		hashedPassword := "$2a$10$ZWqFCxx0mZF7XL4/ZW5zPuqz9K.xF1C1YW5r9q5q5q5q5q5q5q5q5q"
		user := &entity.User{
			ID:        1,
			Username:  "testuser",
			Password:  hashedPassword,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mockRepo.On("FindByAccount", ctx, req.Account).Return(user, nil)

		userDTO, err := service.Login(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, userDTO)
		assert.Equal(t, errors.CodeInvalidCredentials, err.(*errors.Error).Code)
	})
}

func TestAuthService_GetUserByID(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewAuthService(mockRepo)
	ctx := context.Background()

	t.Run("成功获取用户信息", func(t *testing.T) {
		userID := uint(1)
		user := &entity.User{
			ID:        userID,
			Username:  "testuser",
			Email:     "test@example.com",
			Phone:     "13800138000",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mockRepo.On("FindByID", ctx, userID).Return(user, nil)

		userDTO, err := service.GetUserByID(ctx, userID)

		assert.NoError(t, err)
		assert.NotNil(t, userDTO)
		assert.Equal(t, user.ID, userDTO.ID)
		assert.Equal(t, user.Username, userDTO.Username)
		assert.Equal(t, user.Email, userDTO.Email)
		assert.Equal(t, user.Phone, userDTO.Phone)
	})

	t.Run("用户不存在", func(t *testing.T) {
		userID := uint(999)

		mockRepo.On("FindByID", ctx, userID).Return(nil, errors.ErrNotFound)

		userDTO, err := service.GetUserByID(ctx, userID)

		assert.Error(t, err)
		assert.Nil(t, userDTO)
		assert.Equal(t, errors.CodeUserNotFound, err.(*errors.Error).Code)
	})

	t.Run("无效的用户ID", func(t *testing.T) {
		userID := uint(0)

		userDTO, err := service.GetUserByID(ctx, userID)

		assert.Error(t, err)
		assert.Nil(t, userDTO)
		assert.Equal(t, errors.CodeInvalidArgument, err.(*errors.Error).Code)
	})
}
