package grpc

import (
	"context"
	"fmt"
	"net"

	pb "github.com/lazyjean/sla2/api/proto/v1"
	"github.com/lazyjean/sla2/internal/application/dto"
	"github.com/lazyjean/sla2/internal/application/service"
	"github.com/lazyjean/sla2/internal/domain/repository"
	"github.com/lazyjean/sla2/pkg/auth"
	"google.golang.org/grpc"
)

// Server gRPC 服务器
type Server struct {
	pb.UnimplementedUserServiceServer
	server   *grpc.Server
	userRepo repository.UserRepository
	authSvc  auth.JWTServicer
	userSvc  *service.UserService
}

// NewServer 创建 gRPC 服务器
func NewServer(userRepo repository.UserRepository, authSvc auth.JWTServicer) *Server {
	server := grpc.NewServer()

	// 创建服务实例
	userSvc := service.NewUserService(userRepo, authSvc)

	s := &Server{
		server:   server,
		userRepo: userRepo,
		authSvc:  authSvc,
		userSvc:  userSvc,
	}

	// 注册服务
	pb.RegisterUserServiceServer(server, s)

	return s
}

// Start 启动 gRPC 服务器
func (s *Server) Start(port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	if err := s.server.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}

	return nil
}

// Stop 停止 gRPC 服务器
func (s *Server) Stop() {
	s.server.GracefulStop()
}

// 实现 UserServiceServer 接口
func (s *Server) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	resp, err := s.userSvc.Register(ctx, &dto.RegisterRequest{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
		Nickname: req.Nickname,
	})
	if err != nil {
		return nil, err
	}

	return &pb.RegisterResponse{
		User: &pb.User{
			Id:       uint64(resp.UserID),
			Username: resp.Username,
			Email:    resp.Email,
			Nickname: resp.Nickname,
			Avatar:   resp.Avatar,
		},
		Token: resp.Token,
	}, nil
}

func (s *Server) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	resp, err := s.userSvc.Login(ctx, &dto.LoginRequest{
		Account:  req.Account,
		Password: req.Password,
	})
	if err != nil {
		return nil, err
	}

	return &pb.LoginResponse{
		User: &pb.User{
			Id:       uint64(resp.UserID),
			Username: resp.Username,
			Email:    resp.Email,
			Nickname: resp.Nickname,
			Avatar:   resp.Avatar,
		},
		Token: resp.Token,
	}, nil
}

func (s *Server) GetUserInfo(ctx context.Context, req *pb.GetUserInfoRequest) (*pb.GetUserInfoResponse, error) {
	resp, err := s.userSvc.GetUserByID(ctx)
	if err != nil {
		return nil, err
	}

	return &pb.GetUserInfoResponse{
		User: &pb.User{
			Id:       uint64(resp.UserID),
			Username: resp.Username,
			Email:    resp.Email,
			Nickname: resp.Nickname,
			Avatar:   resp.Avatar,
		},
	}, nil
}

func (s *Server) UpdateUserInfo(ctx context.Context, req *pb.UpdateUserInfoRequest) (*pb.UpdateUserInfoResponse, error) {
	err := s.userSvc.UpdateUser(ctx, &dto.UpdateUserRequest{
		Nickname: req.Nickname,
		Avatar:   req.Avatar,
	})
	if err != nil {
		return nil, err
	}

	return &pb.UpdateUserInfoResponse{}, nil
}

func (s *Server) ChangePassword(ctx context.Context, req *pb.ChangePasswordRequest) (*pb.ChangePasswordResponse, error) {
	err := s.userSvc.ChangePassword(ctx, &dto.ChangePasswordRequest{
		OldPassword: req.OldPassword,
		NewPassword: req.NewPassword,
	})
	if err != nil {
		return nil, err
	}

	return &pb.ChangePasswordResponse{}, nil
}

func (s *Server) ResetPassword(ctx context.Context, req *pb.ResetPasswordRequest) (*pb.ResetPasswordResponse, error) {
	err := s.userSvc.ResetPassword(ctx, &dto.ResetPasswordRequest{
		ResetType:        req.ResetType,
		NewPassword:      req.NewPassword,
		Phone:            req.Phone,
		VerificationCode: req.VerificationCode,
		AppleToken:       req.AppleToken,
	})
	if err != nil {
		return nil, err
	}

	return &pb.ResetPasswordResponse{}, nil
}
