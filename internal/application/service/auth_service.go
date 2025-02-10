package service

import (
	"context"

	"github.com/lazyjean/sla2/internal/domain/errors"
	"github.com/lazyjean/sla2/internal/domain/repository"
	"github.com/lazyjean/sla2/pkg/auth"
)

// AuthService 认证服务
type AuthService struct {
	userRepo repository.UserRepository
	jwtSvc   *auth.JWTService
}

func NewAuthService(userRepo repository.UserRepository, jwtSvc *auth.JWTService) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		jwtSvc:   jwtSvc,
	}
}

type RegisterRequest struct {
	Username string
	Password string
	Email    string
	Nickname string
}

type LoginRequest struct {
	Account  string
	Password string
}

type AuthResponse struct {
	UserID   uint
	Username string
	Email    string
	Nickname string
	Avatar   string
	Token    string
}

func (s *AuthService) Register(ctx context.Context, req *RegisterRequest) (*AuthResponse, error) {
	// 验证必填字段
	if req.Username == "" || req.Password == "" || req.Email == "" {
		return nil, errors.NewError(errors.CodeInvalidArgument, "用户名、密码、邮箱不能为空")
	}

	// 检查用户名是否已存在
	exists, err := s.userRepo.ExistsByUsername(ctx, req.Username)
	if err != nil {
		return nil, errors.NewError(errors.CodeInternalError, "检查用户名失败")
	}
	if exists {
		return nil, errors.NewError(errors.CodeUserAlreadyExists, "用户名已存在")
	}

	// 检查邮箱是否已存在
	exists, err = s.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.NewError(errors.CodeInternalError, "检查邮箱失败")
	}
	if exists {
		return nil, errors.NewError(errors.CodeUserAlreadyExists, "邮箱已存在")
	}

	// 加密密码
	hashedPassword, err := s.jwtSvc.HashPassword(req.Password)
	if err != nil {
		return nil, errors.NewError(errors.CodeInternalError, "密码加密失败")
	}

	// 创建用户
	user, err := s.userRepo.Create(ctx, req.Username, hashedPassword, req.Email, req.Nickname)
	if err != nil {
		return nil, errors.NewError(errors.CodeInternalError, "创建用户失败")
	}

	// 生成 token
	token, err := s.jwtSvc.GenerateToken(user.ID)
	if err != nil {
		return nil, errors.NewError(errors.CodeInternalError, "生成token失败")
	}

	return &AuthResponse{
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
		Nickname: user.Nickname,
		Avatar:   user.Avatar,
		Token:    token,
	}, nil
}

func (s *AuthService) Login(ctx context.Context, req *LoginRequest) (*AuthResponse, error) {
	// 尝试通过用户名查找用户
	user, err := s.userRepo.FindByUsername(ctx, req.Account)
	if err != nil {
		// 如果通过用户名找不到，尝试通过邮箱查找
		user, err = s.userRepo.FindByEmail(ctx, req.Account)
		if err != nil {
			return nil, errors.NewError(errors.CodeInvalidCredentials, "账号或密码错误")
		}
	}

	// 验证密码
	if !s.jwtSvc.ComparePasswords(user.Password, req.Password) {
		return nil, errors.NewError(errors.CodeInvalidCredentials, "账号或密码错误")
	}

	// 生成 token
	token, err := s.jwtSvc.GenerateToken(user.ID)
	if err != nil {
		return nil, errors.NewError(errors.CodeInternalError, "生成token失败")
	}

	return &AuthResponse{
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
		Nickname: user.Nickname,
		Avatar:   user.Avatar,
		Token:    token,
	}, nil
}

func (s *AuthService) GetUserByID(ctx context.Context, id uint) (*AuthResponse, error) {
	// 参数校验
	if id == 0 {
		return nil, errors.NewError(errors.CodeInvalidArgument, "用户ID不能为空")
	}

	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.NewError(errors.CodeUserNotFound, "用户不存在")
	}

	return &AuthResponse{
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
		Nickname: user.Nickname,
		Avatar:   user.Avatar,
	}, nil
}
