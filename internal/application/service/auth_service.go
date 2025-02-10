package service

import (
	"context"
	"time"

	"github.com/lazyjean/sla2/internal/application/dto"
	"github.com/lazyjean/sla2/internal/domain/entity"
	"github.com/lazyjean/sla2/internal/domain/errors"
	"github.com/lazyjean/sla2/internal/domain/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(ctx context.Context, req *dto.RegisterDTO) (*dto.UserDTO, error)
	Login(ctx context.Context, req *dto.LoginDTO) (*dto.UserDTO, error)
	GetUserByID(ctx context.Context, id uint) (*dto.UserDTO, error)
}

type authService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{userRepo: userRepo}
}

func (s *authService) Register(ctx context.Context, req *dto.RegisterDTO) (*dto.UserDTO, error) {
	// 验证必填字段
	if req.Username == "" && req.Email == "" && req.Phone == "" {
		return nil, errors.NewError(errors.CodeInvalidArgument, "用户名、邮箱、手机号至少填写一项")
	}

	if req.Password == "" {
		return nil, errors.NewError(errors.CodeInvalidArgument, "密码不能为空")
	}

	// 如果提供了用户名，检查用户名是否已存在
	if req.Username != "" {
		if _, err := s.userRepo.FindByUsername(ctx, req.Username); err == nil {
			return nil, errors.NewError(errors.CodeUserAlreadyExists, "用户名已存在")
		}
	}

	// 如果提供了邮箱，检查邮箱是否已存在
	if req.Email != "" {
		if _, err := s.userRepo.FindByEmail(ctx, req.Email); err == nil {
			return nil, errors.NewError(errors.CodeUserAlreadyExists, "邮箱已存在")
		}
	}

	// 如果提供了手机号，检查手机号是否已存在
	if req.Phone != "" {
		if _, err := s.userRepo.FindByPhone(ctx, req.Phone); err == nil {
			return nil, errors.NewError(errors.CodeUserAlreadyExists, "手机号已存在")
		}
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.NewError(errors.CodeInternalError, "密码加密失败")
	}

	// 创建用户
	user := &entity.User{
		Username:  req.Username,
		Email:     req.Email,
		Phone:     req.Phone,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, errors.NewError(errors.CodeInternalError, "创建用户失败")
	}

	return &dto.UserDTO{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Phone:    user.Phone,
	}, nil
}

func (s *authService) Login(ctx context.Context, req *dto.LoginDTO) (*dto.UserDTO, error) {
	// 查找用户
	user, err := s.userRepo.FindByAccount(ctx, req.Account)
	if err != nil {
		return nil, errors.NewError(errors.CodeInvalidCredentials, "账号或密码错误")
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.NewError(errors.CodeInvalidCredentials, "账号或密码错误")
	}
	return &dto.UserDTO{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Phone:    user.Phone,
	}, nil
}

func (s *authService) GetUserByID(ctx context.Context, id uint) (*dto.UserDTO, error) {
	// 参数校验
	if id == 0 {
		return nil, errors.NewError(errors.CodeInvalidArgument, "用户ID不能为空")
	}

	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.NewError(errors.CodeUserNotFound, "用户不存在")
	}

	return &dto.UserDTO{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Phone:    user.Phone,
	}, nil
}
