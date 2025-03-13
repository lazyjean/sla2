package user

import (
	"context"

	pb "github.com/lazyjean/sla2/api/proto/v1"
	"github.com/lazyjean/sla2/internal/application/dto"
	"github.com/lazyjean/sla2/internal/application/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type UserService struct {
	pb.UnimplementedUserServiceServer
	userService *service.UserService
}

func NewUserService(userService *service.UserService) *UserService {
	return &UserService{
		userService: userService,
	}
}

// Register 用户注册
func (s *UserService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	result, err := s.userService.Register(ctx, &dto.RegisterRequest{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
		Nickname: req.Nickname,
	})
	if err != nil {
		return nil, err
	}

	// 设置 token 到响应头，用于设置 cookie
	if err := grpc.SetHeader(ctx, metadata.Pairs("set-cookie-token", result.Token)); err != nil {
		return nil, err
	}

	return &pb.RegisterResponse{
		Token:        result.Token,
		RefreshToken: result.RefreshToken,
		User: &pb.User{
			Id:       uint64(result.UserID),
			Username: req.Username,
			Email:    req.Email,
			Nickname: req.Nickname,
			Status:   pb.UserStatus_USER_STATUS_ACTIVE,
		},
	}, nil
}

func (s *UserService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	result, err := s.userService.Login(ctx, &dto.LoginRequest{
		Account:  req.Account,
		Password: req.Password,
	})
	if err != nil {
		return nil, err
	}

	if err := grpc.SetHeader(ctx, metadata.Pairs("set-cookie-token", result.Token)); err != nil {
		return nil, err
	}

	return &pb.LoginResponse{
		Token:        result.Token,
		RefreshToken: result.RefreshToken,
		User: &pb.User{
			Id:       uint64(result.UserID),
			Username: result.Username,
			Email:    result.Email,
			Nickname: result.Nickname,
			Avatar:   result.Avatar,
			Status:   pb.UserStatus_USER_STATUS_ACTIVE,
		},
	}, nil
}

func (s *UserService) GetUserInfo(ctx context.Context, req *pb.GetUserInfoRequest) (*pb.GetUserInfoResponse, error) {
	user, err := s.userService.GetLoginUser(ctx)
	if err != nil {
		return nil, err
	}

	return &pb.GetUserInfoResponse{
		User: &pb.User{
			Id:       uint64(user.ID),
			Username: user.Username,
			Email:    user.Email,
			Nickname: user.Nickname,
			Avatar:   user.Avatar,
			Status:   pb.UserStatus_USER_STATUS_ACTIVE,
		},
	}, nil
}

func (s *UserService) UpdateUserInfo(ctx context.Context, req *pb.UpdateUserInfoRequest) (*pb.UpdateUserInfoResponse, error) {
	err := s.userService.UpdateUser(ctx, &dto.UpdateUserRequest{
		Nickname: req.Nickname,
		Avatar:   req.Avatar,
	})
	if err != nil {
		return nil, err
	}

	return &pb.UpdateUserInfoResponse{}, nil
}

func (s *UserService) ChangePassword(ctx context.Context, req *pb.ChangePasswordRequest) (*pb.ChangePasswordResponse, error) {
	err := s.userService.ChangePassword(ctx, &dto.ChangePasswordRequest{
		OldPassword: req.OldPassword,
		NewPassword: req.NewPassword,
	})
	if err != nil {
		return nil, err
	}

	return &pb.ChangePasswordResponse{}, nil
}

func (s *UserService) ResetPassword(ctx context.Context, req *pb.ResetPasswordRequest) (*pb.ResetPasswordResponse, error) {
	err := s.userService.ResetPassword(ctx, &dto.ResetPasswordRequest{
		ResetType:        req.ResetType,
		Phone:            req.Phone,
		VerificationCode: req.VerificationCode,
		NewPassword:      req.NewPassword,
		AppleToken:       req.AppleToken,
	})
	if err != nil {
		return nil, err
	}

	return &pb.ResetPasswordResponse{}, nil
}

// AppleLogin 处理苹果登录请求
func (s *UserService) AppleLogin(ctx context.Context, req *pb.AppleLoginRequest) (*pb.AppleLoginResponse, error) {
	resp, err := s.userService.AppleLogin(ctx, &dto.AppleLoginRequest{
		AuthorizationCode: req.AuthorizationCode,
		UserIdentifier:    req.UserIdentifier,
	})
	if err != nil {
		return nil, err
	}

	// 设置 token 到响应头，用于设置 cookie
	if err := grpc.SetHeader(ctx, metadata.Pairs("set-cookie-token", resp.Token)); err != nil {
		return nil, err
	}

	return &pb.AppleLoginResponse{
		Token:        resp.Token,
		RefreshToken: resp.RefreshToken,
		IsNewUser:    resp.IsNewUser,
		User: &pb.User{
			Id:       uint64(resp.UserID),
			Username: resp.Username,
			Email:    resp.Email,
			Nickname: resp.Nickname,
			Avatar:   resp.Avatar,
			Status:   pb.UserStatus_USER_STATUS_ACTIVE,
		},
	}, nil
}

// RefreshToken 刷新token
func (s *UserService) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	if req.RefreshToken == "" {
		return nil, status.Error(codes.InvalidArgument, "refresh token不能为空")
	}

	resp, err := s.userService.RefreshToken(ctx, &dto.RefreshTokenRequest{
		RefreshToken: req.RefreshToken,
	})
	if err != nil {
		return nil, err
	}

	// 设置新的 token 到响应头，用于设置 cookie
	if err := grpc.SetHeader(ctx, metadata.Pairs("set-cookie-token", resp.AccessToken)); err != nil {
		return nil, err
	}

	return &pb.RefreshTokenResponse{
		Token:        resp.AccessToken,
		RefreshToken: resp.RefreshToken,
	}, nil
}

// Logout 登出
func (s *UserService) Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	// 设置空的 token 到响应头，用于清除 cookie
	if err := grpc.SetHeader(ctx, metadata.Pairs("set-cookie-token", "")); err != nil {
		return nil, err
	}

	return &pb.LogoutResponse{}, nil
}
