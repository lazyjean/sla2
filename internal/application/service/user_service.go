package service

import (
	"context"

	"github.com/lazyjean/sla2/internal/application/dto"
	"github.com/lazyjean/sla2/internal/domain/entity"
	"github.com/lazyjean/sla2/internal/domain/errors"
	"github.com/lazyjean/sla2/internal/domain/repository"
	"github.com/lazyjean/sla2/pkg/auth"
	"github.com/lazyjean/sla2/pkg/utils"
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

func (s *UserService) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.AuthResponse, error) {
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
	hashedPassword, err := s.authSvc.HashPassword(req.Password)
	if err != nil {
		return nil, errors.NewError(errors.CodeInternalError, "密码加密失败")
	}

	// 创建用户
	user, err := s.userRepo.Create(ctx, req.Username, hashedPassword, req.Email, req.Nickname)
	if err != nil {
		return nil, errors.NewError(errors.CodeInternalError, "创建用户失败")
	}

	// 生成 token
	token, err := s.authSvc.GenerateToken(user.ID)
	if err != nil {
		return nil, errors.NewError(errors.CodeInternalError, "生成token失败")
	}

	return &dto.AuthResponse{
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
		Nickname: user.Nickname,
		Avatar:   user.Avatar,
		Token:    token,
	}, nil
}

func (s *UserService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.AuthResponse, error) {
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
	if !s.authSvc.ComparePasswords(user.Password, req.Password) {
		return nil, errors.NewError(errors.CodeInvalidCredentials, "账号或密码错误")
	}

	// 生成 token
	token, err := s.authSvc.GenerateToken(user.ID)
	if err != nil {
		return nil, errors.NewError(errors.CodeInternalError, "生成token失败")
	}

	return &dto.AuthResponse{
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
		Nickname: user.Nickname,
		Avatar:   user.Avatar,
		Token:    token,
	}, nil
}

func (s *UserService) GetUserByID(ctx context.Context) (*dto.AuthResponse, error) {
	userID, err := utils.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, errors.NewError(errors.CodeUserNotFound, "用户不存在")
	}

	return &dto.AuthResponse{
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
		Nickname: user.Nickname,
		Avatar:   user.Avatar,
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
