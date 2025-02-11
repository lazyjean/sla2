package service

import (
	"context"
	"fmt"

	"github.com/lazyjean/sla2/internal/application/dto"
	"github.com/lazyjean/sla2/internal/domain/entity"
	"github.com/lazyjean/sla2/internal/domain/errors"
	"github.com/lazyjean/sla2/internal/domain/repository"
	"github.com/lazyjean/sla2/pkg/auth"
	"github.com/lazyjean/sla2/pkg/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserService struct {
	userRepo repository.UserRepository
	authSvc  auth.JWTServicer
}

func NewUserService(userRepo repository.UserRepository, authSvc auth.JWTServicer) *UserService {
	return &UserService{
		userRepo: userRepo,
		authSvc:  authSvc,
	}
}

func (s *UserService) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.RegisterResponse, error) {
	// 验证必填字段
	if req.Password == "" || (req.Email == "" && req.Username == "") {
		return nil, errors.NewError(errors.CodeInvalidArgument, "密码不能为空, 邮箱或用户名不能为空")
	}

	// 检查用户名是否已存在
	var exists bool
	var err error
	if req.Username != "" {
		exists, err = s.userRepo.ExistsByUsername(ctx, req.Username)
	} else if req.Email != "" {
		exists, err = s.userRepo.ExistsByEmail(ctx, req.Email)
	}

	if err != nil && !errors.Is(err, errors.ErrUserNotFound) {
		return nil, errors.NewError(errors.CodeInternalError, "检查用户名失败")
	}
	if exists {
		return nil, errors.NewError(errors.CodeUserAlreadyExists, "用户名已存在")
	}
	// 加密密码
	hashedPassword, err := s.authSvc.HashPassword(req.Password)
	if err != nil {
		return nil, errors.NewError(errors.CodeInternalError, "密码加密失败")
	}

	// 创建用户
	user, err := s.userRepo.Create(ctx, req.Username, req.Email, hashedPassword, req.Nickname)
	if err != nil {
		return nil, errors.NewError(errors.CodeInternalError, "创建用户失败")
	}

	// 生成 token
	token, err := s.authSvc.GenerateToken(user.ID)
	if err != nil {
		return nil, errors.NewError(errors.CodeInternalError, "生成token失败")
	}
	refreshToken, err := s.authSvc.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, errors.NewError(errors.CodeInternalError, "生成refresh token失败")
	}

	return &dto.RegisterResponse{
		UserID:       uint32(user.ID),
		Token:        token,
		RefreshToken: refreshToken,
	}, nil
}

func (s *UserService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
	var user *entity.User
	var err error

	// 1. 尝试通过邮箱查找
	if utils.IsValidEmail(req.Account) {
		user, err = s.userRepo.FindByEmail(ctx, req.Account)
	}

	// 2. 尝试通过手机号查找
	if err != nil && utils.IsValidPhone(req.Account) {
		user, err = s.userRepo.FindByPhone(ctx, req.Account)
	}

	// 3. 尝试通过用户名查找
	if err != nil && utils.IsValidUsername(req.Account) {
		user, err = s.userRepo.FindByUsername(ctx, req.Account)
	}

	if err != nil {
		return nil, status.Error(codes.NotFound, "用户不存在")
	}

	if user.Password == "" {
		return nil, status.Error(codes.Aborted, "用户未设置密码")
	}

	// validate password
	if !s.authSvc.ComparePasswords(user.Password, req.Password) {
		return nil, status.Error(codes.Unauthenticated, "密码错误")
	}

	// 生成 access token 和 refresh token
	accessToken, err := s.authSvc.GenerateToken(user.ID)
	if err != nil {
		return nil, status.Error(codes.Internal, "生成token失败")
	}

	refreshToken, err := s.authSvc.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, status.Error(codes.Internal, "生成refresh token失败")
	}

	return &dto.LoginResponse{
		UserID:        uint32(user.ID),
		Username:      user.Username,
		Email:         user.Email,
		EmailVerified: user.EmailVerified,
		Nickname:      user.Nickname,
		Avatar:        user.Avatar,
		Token:         accessToken,
		RefreshToken:  refreshToken,
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
		ID:        uint32(user.ID),
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
	if !s.authSvc.ComparePasswords(user.Password, req.OldPassword) {
		return errors.NewError(errors.CodeInvalidCredentials, "旧密码错误")
	}

	// 加密新密码
	hashedPassword, err := s.authSvc.HashPassword(req.NewPassword)
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
	hashedPassword, err := s.authSvc.HashPassword(req.NewPassword)
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
func (s *UserService) AppleLogin(ctx context.Context, idToken string) (*dto.LoginResponse, error) {
	if idToken == "" {
		return nil, errors.NewError(errors.CodeInvalidInput, "Apple ID Token不能为空")
	}

	// 验证 Apple ID Token
	appleIDToken, err := s.authSvc.VerifyAppleIDToken(ctx, idToken)
	if err != nil {
		return nil, errors.NewError(errors.CodeInvalidCredentials, "Apple ID Token验证失败")
	}

	// 查找用户是否已存在
	user, err := s.userRepo.FindByAppleID(ctx, appleIDToken.Subject)
	if err == nil {
		// 用户已存在，生成 token 并返回
		accessToken, err := s.authSvc.GenerateToken(user.ID)
		if err != nil {
			return nil, errors.NewError(errors.CodeInternalError, "生成token失败")
		}

		refreshToken, err := s.authSvc.GenerateRefreshToken(user.ID)
		if err != nil {
			return nil, errors.NewError(errors.CodeInternalError, "生成refresh token失败")
		}

		return &dto.LoginResponse{
			UserID:        uint32(user.ID),
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

	// 创建新用户
	username := fmt.Sprintf("apple_%s", utils.GenerateRandomString(8))
	user = &entity.User{
		Username:      username,
		Email:         appleIDToken.Email,
		EmailVerified: true, // Apple 登录的邮箱默认已验证
		Nickname:      appleIDToken.Name,
		AppleID:       appleIDToken.Subject,
		Status:        entity.UserStatusActive,
	}

	if err := s.userRepo.CreateWithAppleID(ctx, user); err != nil {
		return nil, errors.NewError(errors.CodeInternalError, "创建用户失败")
	}

	// 生成 token
	accessToken, err := s.authSvc.GenerateToken(user.ID)
	if err != nil {
		return nil, errors.NewError(errors.CodeInternalError, "生成token失败")
	}

	refreshToken, err := s.authSvc.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, errors.NewError(errors.CodeInternalError, "生成refresh token失败")
	}

	return &dto.LoginResponse{
		UserID:        uint32(user.ID),
		Username:      user.Username,
		Email:         user.Email,
		EmailVerified: user.EmailVerified,
		Nickname:      user.Nickname,
		Avatar:        user.Avatar,
		Token:         accessToken,
		RefreshToken:  refreshToken,
		IsNewUser:     true,
	}, nil
}

// RefreshToken 刷新token
func (s *UserService) RefreshToken(ctx context.Context, req *dto.RefreshTokenRequest) (*dto.RefreshTokenResponse, error) {
	if req.RefreshToken == "" {
		return nil, errors.NewError(errors.CodeInvalidInput, "refresh token不能为空")
	}

	// 验证refresh token
	userID, err := s.authSvc.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, errors.NewError(errors.CodeInvalidCredentials, "无效的refresh token")
	}

	// 生成新的token
	accessToken, err := s.authSvc.GenerateToken(userID)
	if err != nil {
		return nil, errors.NewError(errors.CodeInternalError, "生成token失败")
	}

	refreshToken, err := s.authSvc.GenerateRefreshToken(userID)
	if err != nil {
		return nil, errors.NewError(errors.CodeInternalError, "生成refresh token失败")
	}

	return &dto.RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
