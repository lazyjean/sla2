package user

import (
	"context"

	pb "github.com/lazyjean/sla2/api/proto/v1"
	"github.com/lazyjean/sla2/internal/domain/entity"
	"github.com/lazyjean/sla2/internal/domain/repository"
	"github.com/lazyjean/sla2/pkg/auth"
	"github.com/lazyjean/sla2/pkg/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserService struct {
	pb.UnimplementedUserServiceServer
	userRepo repository.UserRepository
	authSvc  auth.JWTServicer
}

func NewUserService(userRepo repository.UserRepository, authSvc auth.JWTServicer) *UserService {
	return &UserService{
		userRepo: userRepo,
		authSvc:  authSvc,
	}
}

func (s *UserService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	// 检查用户名是否已存在
	exists, err := s.userRepo.ExistsByUsername(ctx, req.Username)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to check username")
	}
	if exists {
		return nil, status.Error(codes.AlreadyExists, "username already exists")
	}

	// 检查邮箱是否已存在
	exists, err = s.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to check email")
	}
	if exists {
		return nil, status.Error(codes.AlreadyExists, "email already exists")
	}

	// 创建用户
	hashedPassword, err := s.authSvc.HashPassword(req.Password)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to hash password")
	}

	user, err := s.userRepo.Create(ctx, req.Username, hashedPassword, req.Email, req.Nickname)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create user")
	}

	// 生成 token
	token, err := s.authSvc.GenerateToken(user.ID)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to generate token")
	}
	refreshToken, err := s.authSvc.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to generate refresh token")
	}

	return &pb.RegisterResponse{
		User:         convertUserToPb(user),
		Token:        token,
		RefreshToken: refreshToken,
	}, nil
}

func (s *UserService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
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

	return &pb.LoginResponse{
		User:         convertUserToPb(user),
		Token:        accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *UserService) GetUserInfo(ctx context.Context, req *pb.GetUserInfoRequest) (*pb.GetUserInfoResponse, error) {
	userID, err := utils.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "user id not found")
	}

	user, err := s.userRepo.FindByID(ctx, uint(userID))
	if err != nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	return &pb.GetUserInfoResponse{
		User: convertUserToPb(user),
	}, nil
}

func (s *UserService) UpdateUserInfo(ctx context.Context, req *pb.UpdateUserInfoRequest) (*pb.UpdateUserInfoResponse, error) {
	userID, err := utils.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "user id not found")
	}

	user, err := s.userRepo.FindByID(ctx, uint(userID))
	if err != nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	// 更新用户信息
	user.Nickname = req.Nickname
	user.Avatar = req.Avatar

	err = s.userRepo.Update(ctx, user)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to update user")
	}

	return &pb.UpdateUserInfoResponse{}, nil
}

func (s *UserService) ChangePassword(ctx context.Context, req *pb.ChangePasswordRequest) (*pb.ChangePasswordResponse, error) {
	user, err := s.userRepo.FindByID(ctx, uint(req.UserId))
	if err != nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	// 验证旧密码
	if !s.authSvc.ComparePasswords(user.Password, req.OldPassword) {
		return nil, status.Error(codes.InvalidArgument, "invalid old password")
	}

	// 更新密码
	hashedPassword, err := s.authSvc.HashPassword(req.NewPassword)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to hash password")
	}

	user.Password = hashedPassword
	err = s.userRepo.Update(ctx, user)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to update password")
	}

	return &pb.ChangePasswordResponse{}, nil
}

func (s *UserService) ResetPassword(ctx context.Context, req *pb.ResetPasswordRequest) (*pb.ResetPasswordResponse, error) {
	var user *entity.User
	var err error

	// 验证请求参数
	if req.NewPassword == "" {
		return nil, status.Error(codes.InvalidArgument, "new password is required")
	}

	switch req.ResetType {
	case "phone":
		// 验证手机号和验证码
		if req.Phone == "" || req.VerificationCode == "" {
			return nil, status.Error(codes.InvalidArgument, "phone and verification code are required")
		}
		// TODO: 验证验证码
		// if !s.verifySvc.VerifyCode(ctx, req.Phone, req.VerificationCode) {
		//     return nil, status.Error(codes.InvalidArgument, "invalid verification code")
		// }

		// 通过手机号查找用户
		user, err = s.userRepo.FindByPhone(ctx, req.Phone)
		if err != nil {
			return nil, status.Error(codes.NotFound, "user not found")
		}

	case "apple":
		// 验证苹果登录票据
		if req.AppleToken == "" {
			return nil, status.Error(codes.InvalidArgument, "apple token is required")
		}

		// 验证 Apple ID Token
		appleIDToken, err := s.authSvc.VerifyAppleIDToken(ctx, req.AppleToken)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid apple token")
		}

		// 通过 Apple ID 查找用户
		user, err = s.userRepo.FindByAppleID(ctx, appleIDToken.Subject)
		if err != nil {
			return nil, status.Error(codes.NotFound, "user not found")
		}

	default:
		return nil, status.Error(codes.InvalidArgument, "unsupported reset type")
	}

	// 更新密码
	hashedPassword, err := s.authSvc.HashPassword(req.NewPassword)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to hash password")
	}

	user.Password = hashedPassword
	err = s.userRepo.Update(ctx, user)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to update password")
	}

	return &pb.ResetPasswordResponse{}, nil
}

func (s *UserService) AppleLogin(ctx context.Context, req *pb.AppleLoginRequest) (*pb.AppleLoginResponse, error) {
	// 验证 Apple ID Token
	appleIDToken, err := s.authSvc.VerifyAppleIDToken(ctx, req.IdToken)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid apple id token")
	}

	// 查找用户是否已存在
	user, err := s.userRepo.FindByAppleID(ctx, appleIDToken.Subject)
	if err == nil {
		// 用户已存在，生成 token 并返回
		token, err := s.authSvc.GenerateToken(user.ID)
		if err != nil {
			return nil, status.Error(codes.Internal, "failed to generate token")
		}

		return &pb.AppleLoginResponse{
			User:      convertUserToPb(user),
			Token:     token,
			IsNewUser: false,
		}, nil
	}

	// 创建用户
	user = &entity.User{
		Email:    appleIDToken.Email,
		Nickname: appleIDToken.Name,
		AppleID:  appleIDToken.Subject,
		Status:   entity.UserStatusActive,
	}

	err = s.userRepo.CreateWithAppleID(ctx, user)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create user")
	}

	// 生成 token
	token, err := s.authSvc.GenerateToken(user.ID)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to generate token")
	}

	return &pb.AppleLoginResponse{
		User:      convertUserToPb(user),
		Token:     token,
		IsNewUser: true,
	}, nil
}

// 添加刷新令牌接口
func (s *UserService) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	// 验证 refresh token
	userID, err := s.authSvc.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "无效的刷新令牌")
	}

	// 生成新的令牌对
	accessToken, refreshToken, err := s.authSvc.GenerateTokenPair(userID)
	if err != nil {
		return nil, status.Error(codes.Internal, "生成token失败")
	}

	return &pb.RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// 辅助函数：将 domain.User 转换为 pb.User
func convertUserToPb(user *entity.User) *pb.User {
	return &pb.User{
		Id:        uint64(user.ID),
		Username:  user.Username,
		Email:     user.Email,
		Nickname:  user.Nickname,
		Avatar:    user.Avatar,
		Status:    convertUserStatusToPb(user.Status),
		CreatedAt: user.CreatedAt.Unix(),
		UpdatedAt: user.UpdatedAt.Unix(),
	}
}

// 辅助函数：将 domain.UserStatus 转换为 pb.UserStatus
func convertUserStatusToPb(status entity.UserStatus) pb.UserStatus {
	switch status {
	case entity.UserStatusActive:
		return pb.UserStatus_USER_STATUS_ACTIVE
	case entity.UserStatusInactive:
		return pb.UserStatus_USER_STATUS_INACTIVE
	case entity.UserStatusSuspended:
		return pb.UserStatus_USER_STATUS_SUSPENDED
	default:
		return pb.UserStatus_USER_STATUS_UNSPECIFIED
	}
}
