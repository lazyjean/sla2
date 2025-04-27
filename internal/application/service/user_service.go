package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/lazyjean/sla2/internal/application/dto"
	"github.com/lazyjean/sla2/internal/domain/entity"
	"github.com/lazyjean/sla2/internal/domain/errors"
	"github.com/lazyjean/sla2/internal/domain/oauth"
	"github.com/lazyjean/sla2/internal/domain/repository"
	"github.com/lazyjean/sla2/internal/domain/security"
	"github.com/lazyjean/sla2/pkg/logger"
	"github.com/lazyjean/sla2/pkg/utils"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	if req.Username == "" && req.Email == "" {
		return nil, errors.NewError(
			errors.CodeInvalidArgument,
			"用户名或邮箱至少需要填写一个",
		)
	}

	// 密码不能为空
	if req.Password == "" {
		return nil, errors.NewError(
			errors.CodeInvalidArgument,
			"密码不能为空",
		)
	}

	// 加密密码
	hashedPassword, err := s.passwordService.HashPassword(req.Password)
	if err != nil {
		return nil, errors.NewError(errors.CodeInternalError, "密码加密失败")
	}

	// 创建用户
	user, err := s.userRepo.Create(ctx, req.Username, req.Email, hashedPassword, req.Nickname, "", false)
	if err != nil {
		return nil, errors.NewError(errors.CodeInternalError, "创建用户失败")
	}

	// 生成 token，用户统一使用 "user" 角色
	accessToken, err := s.tokenService.GenerateToken(user.ID, []string{"user"})
	if err != nil {
		return nil, errors.NewError(errors.CodeInternalError, "生成token失败")
	}
	refreshToken, err := s.tokenService.GenerateRefreshToken(user.ID, []string{"user"})
	if err != nil {
		return nil, errors.NewError(errors.CodeInternalError, "生成refresh token失败")
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
	log.Info("登录请求",
		zap.String("account", req.Account),
	)

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
		return nil, status.Error(codes.NotFound, "用户不存在")
	}

	if user == nil {
		log.Error("用户对象为空",
			zap.String("account", req.Account),
		)
		return nil, status.Error(codes.NotFound, "用户不存在")
	}

	if user.Password == "" {
		return nil, status.Error(codes.Aborted, "用户未设置密码")
	}

	// 验证密码
	if !s.passwordService.VerifyPassword(req.Password, user.Password) {
		return nil, status.Error(codes.Unauthenticated, "密码错误")
	}

	// 生成 token，用户统一使用 "user" 角色
	accessToken, err := s.tokenService.GenerateToken(user.ID, []string{"user"})
	if err != nil {
		return nil, status.Error(codes.Internal, "生成token失败")
	}

	refreshToken, err := s.tokenService.GenerateRefreshToken(user.ID, []string{"user"})
	if err != nil {
		return nil, status.Error(codes.Internal, "生成refresh token失败")
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

func (s *UserService) GetLoginUser(ctx context.Context) (*dto.UserDTO, error) {
	userID, err := utils.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, errors.NewError(errors.CodeUserNotFound, "用户不存在")
	}

	status := "active"
	switch user.Status {
	case entity.UserStatusInactive:
		status = "inactive"
	case entity.UserStatusSuspended:
		status = "suspended"
	}

	return &dto.UserDTO{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Nickname:  user.Nickname,
		Avatar:    user.Avatar,
		Status:    status,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (s *UserService) UpdateUser(ctx context.Context, req *dto.UpdateUserRequest) error {
	userID, err := utils.GetUserIDFromContext(ctx)
	if err != nil {
		return err
	}

	// 获取用户信息
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return errors.NewError(errors.CodeUserNotFound, "用户不存在")
	}

	// 更新用户信息
	user.Nickname = req.Nickname
	user.Avatar = req.Avatar

	// 保存更新
	if err := s.userRepo.Update(ctx, user); err != nil {
		return errors.NewError(errors.CodeInternalError, "更新用户信息失败")
	}

	return nil
}

func (s *UserService) ChangePassword(ctx context.Context, req *dto.ChangePasswordRequest) error {
	userID, err := utils.GetUserIDFromContext(ctx)
	if err != nil {
		return err
	}
	if req.OldPassword == "" || req.NewPassword == "" {
		return errors.NewError(errors.CodeInvalidArgument, "旧密码和新密码不能为空")
	}

	// 获取用户信息
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return errors.NewError(errors.CodeUserNotFound, "用户不存在")
	}

	// 验证旧密码
	if !s.passwordService.VerifyPassword(req.OldPassword, user.Password) {
		return errors.NewError(errors.CodeInvalidCredentials, "旧密码错误")
	}

	// 加密新密码
	hashedPassword, err := s.passwordService.HashPassword(req.NewPassword)
	if err != nil {
		return errors.NewError(errors.CodeInternalError, "密码加密失败")
	}

	// 更新密码
	user.Password = hashedPassword
	if err := s.userRepo.Update(ctx, user); err != nil {
		return errors.NewError(errors.CodeInternalError, "更新密码失败")
	}

	return nil
}

func (s *UserService) ResetPassword(ctx context.Context, req *dto.ResetPasswordRequest) error {
	// 验证请求参数
	if req.NewPassword == "" {
		return errors.NewError(errors.CodeInvalidArgument, "新密码不能为空")
	}

	var user *entity.User
	var err error

	switch req.ResetType {
	case "phone":
		// 验证手机号和验证码
		if req.Phone == "" || req.VerificationCode == "" {
			return errors.NewError(errors.CodeInvalidArgument, "手机号和验证码不能为空")
		}
		// TODO: 调用验证码服务验证手机号和验证码
		// 验证通过后，根据手机号查找用户
		user, err = s.userRepo.FindByPhone(ctx, req.Phone)
		if err != nil {
			return errors.NewError(errors.CodeUserNotFound, "用户不存在")
		}

	case "apple":
		// 验证苹果登录票据
		if req.AppleToken == "" {
			return errors.NewError(errors.CodeInvalidArgument, "苹果登录票据不能为空")
		}
		// TODO: 调用苹果登录服务验证票据
		// 验证通过后，根据苹果用户ID查找用户
		appleUserID := "TODO: 从AppleToken中获取用户ID"
		user, err = s.userRepo.FindByAppleID(ctx, appleUserID)
		if err != nil {
			return errors.NewError(errors.CodeUserNotFound, "用户不存在")
		}

	default:
		return errors.NewError(errors.CodeInvalidArgument, "不支持的重置方式")
	}

	// 加密新密码
	hashedPassword, err := s.passwordService.HashPassword(req.NewPassword)
	if err != nil {
		return errors.NewError(errors.CodeInternalError, "密码加密失败")
	}

	// 更新密码
	user.Password = hashedPassword
	if err := s.userRepo.Update(ctx, user); err != nil {
		return errors.NewError(errors.CodeInternalError, "重置密码失败")
	}

	return nil
}

// AppleLogin 处理苹果登录
func (s *UserService) AppleLogin(ctx context.Context, req *dto.AppleLoginRequest) (*dto.AppleLoginResponse, error) {
	// 验证 Apple Authorization Code
	appleIDToken, err := s.appleAuth.AuthCodeWithApple(ctx, req.AuthorizationCode)
	if err != nil || appleIDToken.Subject != req.UserIdentifier {
		return nil, errors.NewError(errors.CodeInvalidCredentials, "Apple Authorization Code 验证失败")
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
		user, err = s.userRepo.Create(ctx, "", appleIDToken.Email, "", nickname, appleIDToken.Subject, emailVerified)
		if err != nil {
			return nil, errors.NewError(errors.CodeInternalError, "创建用户失败")
		}
	}

	// 生成 token，用户统一使用 "user" 角色
	token, err := s.tokenService.GenerateToken(user.ID, []string{"user"})
	if err != nil {
		return nil, errors.NewError(errors.CodeInternalError, "生成 token 失败")
	}

	refreshToken, err := s.tokenService.GenerateRefreshToken(user.ID, []string{"user"})
	if err != nil {
		return nil, errors.NewError(errors.CodeInternalError, "生成 refresh token 失败")
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
func (s *UserService) RefreshToken(ctx context.Context, req *dto.RefreshTokenRequest) (*dto.UserRefreshTokenResponse, error) {
	// 验证刷新令牌
	userID, roles, err := s.tokenService.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, err
	}

	// 生成新的令牌，使用相同的角色信息
	accessToken, err := s.tokenService.GenerateToken(userID, roles)
	if err != nil {
		return nil, errors.NewError(errors.CodeInternalError, "生成 token 失败")
	}

	refreshToken, err := s.tokenService.GenerateRefreshToken(userID, roles)
	if err != nil {
		return nil, errors.NewError(errors.CodeInternalError, "生成 refresh token 失败")
	}

	return &dto.UserRefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
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
		exists, err := s.userRepo.ExistsByAppleID(ctx, appleUserID.Subject)
		if err != nil {
			return "", "", err
		}
		if !exists {
			// 创建新用户
			username := fmt.Sprintf("apple_%s", appleUserID.Subject)
			nickname := fmt.Sprintf("Apple User %s", appleUserID.Subject)
			user, err = s.userRepo.Create(ctx, username, appleUserID.Email, "", nickname, appleUserID.Subject, true)
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
func (s *UserService) Logout(ctx context.Context, req *dto.LogoutRequest) (*dto.LogoutResponse, error) {
	// Get user ID from context
	userID, err := utils.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// Log the logout action
	log := logger.GetLogger(ctx)
	log.Info("User logged out",
		zap.Int64("user_id", int64(userID)),
	)

	return &dto.LogoutResponse{}, nil
}
