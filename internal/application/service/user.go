package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/lazyjean/sla2/internal/application/dto"
	"github.com/lazyjean/sla2/internal/domain/entity"
	"github.com/lazyjean/sla2/internal/domain/errors"
	"github.com/lazyjean/sla2/internal/domain/oauth"
	"github.com/lazyjean/sla2/internal/domain/repository"
	"github.com/lazyjean/sla2/internal/domain/security"
	"github.com/lazyjean/sla2/pkg/logger"
	"github.com/lazyjean/sla2/pkg/utils"
	"go.uber.org/zap"
)

type UserService struct {
	userRepo        repository.UserRepository
	tokenService    security.TokenService
	passwordService security.PasswordService
	appleAuth       oauth.AppleAuthService
}

func NewUserService(
	userRepo repository.UserRepository,
	tokenService security.TokenService,
	passwordService security.PasswordService,
	appleAuth oauth.AppleAuthService,
) *UserService {
	return &UserService{
		userRepo:        userRepo,
		tokenService:    tokenService,
		passwordService: passwordService,
		appleAuth:       appleAuth,
	}
}

func (s *UserService) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.RegisterResponse, error) {
	// 验证必填字段：至少需要用户名或邮箱
	if req.Username == "" || req.Email == "" {
		// todo: 这里需要更准确的错误描述, 以及要考虑是否使用统一的参数校验逻辑
		return nil, errors.ErrInvalidInput
	}

	// 密码不能为空
	if req.Password == "" {
		return nil, errors.ErrInvalidInput
	}

	// 加密密码
	hashedPassword, err := s.passwordService.HashPassword(req.Password)
	if err != nil {
		return nil, errors.ErrFailedToSave
	}

	// 创建用户
	user := &entity.User{
		Username:      req.Username,
		Email:         req.Email,
		Password:      hashedPassword,
		Nickname:      req.Nickname,
		AppleID:       "",
		EmailVerified: false,
	}
	err = s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, errors.ErrFailedToSave
	}

	// 生成 token，用户统一使用 "user" 角色
	accessToken, err := s.tokenService.GenerateToken(user.ID, []string{"user"})
	if err != nil {
		return nil, errors.ErrUnauthenticated
	}
	refreshToken, err := s.tokenService.GenerateRefreshToken(user.ID, []string{"user"})
	if err != nil {
		return nil, errors.ErrUnauthenticated
	}

	return &dto.RegisterResponse{
		UserID:       user.ID,
		Token:        accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *UserService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
	var user *entity.User
	var err error
	log := logger.GetLogger(ctx)

	// 1. 尝试通过邮箱查找
	if utils.IsValidEmail(req.Account) {
		log.Debug("通过邮箱查找用户",
			zap.String("email", req.Account),
		)
		user, err = s.userRepo.FindByEmail(ctx, req.Account)
		if err != nil {
			log.Error("通过邮箱查找用户失败",
				zap.String("email", req.Account),
				zap.Error(err),
			)
		}
	} else if utils.IsValidPhone(req.Account) {
		// 2. 尝试通过手机号查找
		log.Debug("通过手机号查找用户",
			zap.String("phone", req.Account),
		)
		user, err = s.userRepo.FindByPhone(ctx, req.Account)
		if err != nil {
			log.Error("通过手机号查找用户失败",
				zap.String("phone", req.Account),
				zap.Error(err),
			)
		}
	} else if utils.IsValidUsername(req.Account) {
		// 3. 尝试通过用户名查找
		log.Debug("通过用户名查找用户",
			zap.String("username", req.Account),
		)
		user, err = s.userRepo.FindByUsername(ctx, req.Account)
		if err != nil {
			log.Error("通过用户名查找用户失败",
				zap.String("username", req.Account),
				zap.Error(err),
			)
		}
	}

	if err != nil {
		log.Error("用户查找失败",
			zap.String("account", req.Account),
			zap.Error(err),
		)
		return nil, errors.ErrUserNotFound
	}

	if user == nil {
		log.Error("用户对象为空",
			zap.String("account", req.Account),
		)
		return nil, errors.ErrUserNotFound
	}

	if user.Password == "" {
		return nil, errors.ErrEmptyPassword
	}

	// 验证密码
	if !s.passwordService.VerifyPassword(req.Password, user.Password) {
		return nil, errors.ErrInvalidPassword
	}

	// 生成 token，用户统一使用 "user" 角色
	accessToken, err := s.tokenService.GenerateToken(user.ID, []string{"user"})
	if err != nil {
		return nil, errors.ErrUnauthenticated
	}

	refreshToken, err := s.tokenService.GenerateRefreshToken(user.ID, []string{"user"})
	if err != nil {
		return nil, errors.ErrUnauthenticated
	}

	return &dto.LoginResponse{
		UserID:        user.ID,
		Username:      user.Username,
		Email:         user.Email,
		EmailVerified: user.EmailVerified,
		Nickname:      user.Nickname,
		Avatar:        user.Avatar,
		Token:         accessToken,
		RefreshToken:  refreshToken,
		IsNewUser:     false,
	}, nil
}

func (s *UserService) GetLoginUser(ctx context.Context) (*entity.User, error) {
	userID, err := GetUserID(ctx)
	if err != nil {
		return nil, errors.ErrLoginUserIdIsMissingInCtx
	}
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.ErrUserNotFound
	}
	return user, nil
}

func (s *UserService) UpdateUser(ctx context.Context, req *dto.UpdateUserRequest) error {
	var userID entity.UID
	var err error

	if userID, err = GetUserID(ctx); err == nil {
		return errors.ErrLoginUserIdIsMissingInCtx
	}

	user := &entity.User{
		ID:        userID,
		Nickname:  req.Nickname,
		Avatar:    req.Avatar,
		UpdatedAt: time.Now(),
	}
	
	// 保存更新
	if err := s.userRepo.Update(ctx, user); err != nil {
		return errors.ErrFailedToUpdate
	}
	return nil
}

func (s *UserService) ChangePassword(ctx context.Context, req *dto.ChangePasswordRequest) error {
	userID, err := GetUserID(ctx)
	if err != nil {
		return errors.ErrLoginUserIdIsMissingInCtx
	}
	if req.OldPassword == "" || req.NewPassword == "" {
		return errors.ErrInvalidInput
	}

	// 获取用户信息
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return errors.ErrUserNotFound
	}

	// 验证旧密码
	if !s.passwordService.VerifyPassword(req.OldPassword, user.Password) {
		return errors.ErrInvalidPassword
	}

	// 加密新密码
	hashedPassword, err := s.passwordService.HashPassword(req.NewPassword)
	if err != nil {
		return errors.ErrFailedToSave
	}

	// 更新密码
	user.Password = hashedPassword
	if err := s.userRepo.Update(ctx, user); err != nil {
		return errors.ErrFailedToUpdate
	}

	return nil
}

func (s *UserService) ResetPassword(ctx context.Context, req *dto.ResetPasswordRequest) error {
	// 验证请求参数
	if req.NewPassword == "" {
		return errors.ErrInvalidInput
	}

	var user *entity.User
	var err error

	switch req.ResetType {
	case "phone":
		// 验证手机号和验证码
		if req.Phone == "" || req.VerificationCode == "" {
			return errors.ErrInvalidInput
		}
		// TODO: 调用验证码服务验证手机号和验证码
		// 验证通过后，根据手机号查找用户
		user, err = s.userRepo.FindByPhone(ctx, req.Phone)
		if err != nil {
			return errors.ErrUserNotFound
		}

	case "apple":
		// 验证苹果登录票据
		if req.AppleToken == "" {
			return errors.ErrInvalidInput
		}
		// TODO: 调用苹果登录服务验证票据
		// 验证通过后，根据苹果用户ID查找用户
		appleUserID := "TODO: 从AppleToken中获取用户ID"
		user, err = s.userRepo.FindByAppleID(ctx, appleUserID)
		if err != nil {
			return errors.ErrUserNotFound
		}

	default:
		return errors.ErrInvalidInput
	}

	// 加密新密码
	hashedPassword, err := s.passwordService.HashPassword(req.NewPassword)
	if err != nil {
		return errors.ErrFailedToSave
	}

	// 更新密码
	user.Password = hashedPassword
	if err := s.userRepo.Update(ctx, user); err != nil {
		return errors.ErrFailedToUpdate
	}

	return nil
}

// AppleLogin 处理苹果登录
func (s *UserService) AppleLogin(ctx context.Context, req *dto.AppleLoginRequest) (*dto.AppleLoginResponse, error) {
	// 验证 Apple Authorization Code
	appleIDToken, err := s.appleAuth.AuthCodeWithApple(ctx, req.AuthorizationCode)
	if err != nil || appleIDToken.Subject != req.UserIdentifier {
		return nil, errors.ErrUnauthenticated
	}

	// 查找用户是否已存在
	user, err := s.userRepo.FindByAppleID(ctx, appleIDToken.Subject)
	isNewUser := false

	// 这里需要明确是用户不存在的错误还是其他错误，只有用户不存在时才创建用户
	var domainErr *errors.Error
	if err != nil && errors.As(err, &domainErr) && domainErr.Code == errors.CodeUserNotFound {
		isNewUser = true

		// 生成昵称
		nickname := fmt.Sprintf("apple_%s", strings.Split(appleIDToken.Subject, ".")[0])

		// 创建用户，用户名留空
		emailVerified := appleIDToken.Email != ""
		user := &entity.User{
			Username:      "",
			Email:         appleIDToken.Email,
			Password:      "",
			Nickname:      nickname,
			AppleID:       appleIDToken.Subject,
			EmailVerified: emailVerified,
		}
		err = s.userRepo.Create(ctx, user)
		if err != nil {
			return nil, errors.ErrFailedToSave
		}
	}

	// 生成 token，用户统一使用 "user" 角色
	token, err := s.tokenService.GenerateToken(user.ID, []string{"user"})
	if err != nil {
		return nil, errors.ErrUnauthenticated
	}

	refreshToken, err := s.tokenService.GenerateRefreshToken(user.ID, []string{"user"})
	if err != nil {
		return nil, errors.ErrUnauthenticated
	}

	return &dto.AppleLoginResponse{
		UserID:       user.ID,
		Username:     user.Username,
		Email:        user.Email,
		Nickname:     user.Nickname,
		Avatar:       user.Avatar,
		Token:        token,
		RefreshToken: refreshToken,
		IsNewUser:    isNewUser,
	}, nil
}

// RefreshToken 刷新访问令牌
func (s *UserService) RefreshToken(ctx context.Context, token string) (string, string, error) {
	var uid entity.UID
	var err error
	if uid, err = GetUserID(ctx); err != nil {
		return "", "", errors.ErrLoginUserIdIsMissingInCtx
	}

	// 验证令牌并检查用户身份
	if id, roles, err := s.tokenService.ValidateToken(token); err != nil {
		return "", "", errors.ErrUnauthenticated
	} else if id != uid {
		return "", "", errors.ErrRefreshTokenMismatch
	} else {
		// 生成新的令牌
		accessToken, err := s.tokenService.GenerateToken(uid, roles)
		if err != nil {
			return "", "", errors.ErrUnauthenticated
		}

		refreshToken, err := s.tokenService.GenerateRefreshToken(uid, roles)
		if err != nil {
			return "", "", errors.ErrUnauthenticated
		}

		return accessToken, refreshToken, nil
	}
}

func (s *UserService) LoginWithApple(ctx context.Context, appleToken string) (string, string, error) {
	// 验证 Apple token
	appleUserID, err := s.appleAuth.AuthCodeWithApple(ctx, appleToken)
	if err != nil {
		return "", "", err
	}

	// 查找或创建用户
	user, err := s.userRepo.FindByAppleID(ctx, appleUserID.Subject)
	if err != nil {
		// 检查是否是用户不存在的错误
		if errors.Is(err, errors.ErrUserNotFound) {
			// 创建新用户
			username := fmt.Sprintf("apple_%s", appleUserID.Subject)
			nickname := fmt.Sprintf("Apple User %s", appleUserID.Subject)
			user := &entity.User{
				Username:      username,
				Email:         appleUserID.Email,
				Password:      "",
				Nickname:      nickname,
				AppleID:       appleUserID.Subject,
				EmailVerified: true,
			}
			err = s.userRepo.Create(ctx, user)
			if err != nil {
				return "", "", err
			}
		} else {
			return "", "", err
		}
	}

	// 生成 token，用户统一使用 "user" 角色
	token, err := s.tokenService.GenerateToken(user.ID, []string{"user"})
	if err != nil {
		return "", "", err
	}

	refreshToken, err := s.tokenService.GenerateRefreshToken(user.ID, []string{"user"})
	if err != nil {
		return "", "", err
	}

	return token, refreshToken, nil
}

// Logout handles user logout
func (s *UserService) Logout(ctx context.Context) error {
	_, err := GetUserID(ctx)
	if err != nil {
		return err
	}
	return nil
}
