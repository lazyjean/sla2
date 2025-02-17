package user

import (
	"context"

	pb "github.com/lazyjean/sla2/api/proto/v1"
	"github.com/lazyjean/sla2/internal/application/dto"
	"github.com/lazyjean/sla2/internal/application/service"
	"github.com/lazyjean/sla2/pkg/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserService struct {
	pb.UnimplementedUserServiceServer
	userService *service.UserService
	authSvc     auth.JWTServicer
}

func NewUserService(userService *service.UserService, authSvc auth.JWTServicer) *UserService {
	return &UserService{
		userService: userService,
		authSvc:     authSvc,
	}
}

func (s *UserService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	user, err := s.userService.Register(ctx, &dto.RegisterRequest{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
	})
	if err != nil {
		return nil, err
	}

	return &pb.RegisterResponse{
		User: convertUserToPb(user),
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

	return &pb.LoginResponse{
		User:         convertUserToPb(result),
		Token:        result.Token,
		RefreshToken: result.RefreshToken,
	}, nil
}

func (s *UserService) GetUserInfo(ctx context.Context, req *pb.GetUserInfoRequest) (*pb.GetUserInfoResponse, error) {
	user, err := s.userService.GetLoginUser(ctx)
	if err != nil {
		return nil, err
	}

	return &pb.GetUserInfoResponse{
		User: convertUserToPb(user),
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
	// 调用应用层服务处理苹果登录
	resp, err := s.userService.AppleLogin(ctx, req.AuthorizationCode)
	if err != nil {
		return nil, err
	}

	return &pb.AppleLoginResponse{
		User: &pb.User{
			Id:       uint64(resp.UserID),
			Username: resp.Username,
			Email:    resp.Email,
			Nickname: resp.Nickname,
			Avatar:   resp.Avatar,
		},
		Token:        resp.Token,
		RefreshToken: resp.RefreshToken,
		IsNewUser:    resp.IsNewUser,
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

	return &pb.RefreshTokenResponse{
		Token:        resp.AccessToken,
		RefreshToken: resp.RefreshToken,
	}, nil
}

// 辅助函数：将领域对象转换为 pb 对象
func convertUserToPb(user interface{}) *pb.User {
	switch u := user.(type) {
	case *dto.LoginResponse:
		return &pb.User{
			Id:       uint64(u.UserID),
			Username: u.Username,
			Email:    u.Email,
			Nickname: u.Nickname,
			Avatar:   u.Avatar,
			Status:   pb.UserStatus_USER_STATUS_ACTIVE,
		}
	case *dto.RegisterResponse:
		return &pb.User{
			Id:     uint64(u.UserID),
			Status: pb.UserStatus_USER_STATUS_ACTIVE,
		}
	case *dto.UserDTO:
		return &pb.User{
			Id:        uint64(u.ID),
			Username:  u.Username,
			Email:     u.Email,
			Nickname:  u.Nickname,
			Avatar:    u.Avatar,
			Status:    convertUserStatusToPb(u.Status),
			CreatedAt: u.CreatedAt.Unix(),
			UpdatedAt: u.UpdatedAt.Unix(),
		}
	default:
		return nil
	}
}

func convertUserStatusToPb(status string) pb.UserStatus {
	switch status {
	case "active":
		return pb.UserStatus_USER_STATUS_ACTIVE
	case "inactive":
		return pb.UserStatus_USER_STATUS_INACTIVE
	case "suspended":
		return pb.UserStatus_USER_STATUS_SUSPENDED
	default:
		return pb.UserStatus_USER_STATUS_UNSPECIFIED
	}
}
