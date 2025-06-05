package admin

import (
	"context"

	"github.com/lazyjean/sla2/internal/transport/grpc/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	pb "github.com/lazyjean/sla2/api/proto/v1"
	"github.com/lazyjean/sla2/internal/application/dto"
	"github.com/lazyjean/sla2/internal/application/service"
	"github.com/lazyjean/sla2/internal/transport/grpc/admin/converter"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Service 管理员服务
type Service struct {
	pb.UnimplementedAdminServiceServer
	adminService *service.AdminService
	converter    *converter.AdminConverter
}

// NewAdminService 创建新的管理员服务
func NewAdminService(adminService *service.AdminService) *Service {
	return &Service{
		adminService: adminService,
		converter:    converter.NewAdminConverter(),
	}
}

// CheckSystemStatus 检查系统状态
func (s *Service) CheckSystemStatus(ctx context.Context, req *pb.AdminServiceCheckSystemStatusRequest) (*pb.AdminServiceCheckSystemStatusResponse, error) {
	initialized, err := s.adminService.IsSystemInitialized(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.AdminServiceCheckSystemStatusResponse{
		Initialized: initialized,
	}, nil
}

// InitializeSystem 初始化系统
func (s *Service) InitializeSystem(ctx context.Context, req *pb.AdminServiceInitializeSystemRequest) (*pb.AdminServiceInitializeSystemResponse, error) {
	resp, err := s.adminService.InitializeSystem(ctx, &dto.InitializeSystemRequest{
		Username: req.Username,
		Nickname: req.Nickname,
		Password: req.Password,
		Email:    req.Email,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	md := metadata.Pairs(
		middleware.MDHeaderAccessToken, resp.AccessToken,
		middleware.MDHeaderRefreshToken, resp.RefreshToken,
	)
	if err := grpc.SetHeader(ctx, md); err != nil {
		return nil, err
	}
	return &pb.AdminServiceInitializeSystemResponse{
		Admin:        s.converter.ToPB(resp.Admin),
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
	}, nil
}

// AdminLogin 管理员登录
func (s *Service) AdminLogin(ctx context.Context, req *pb.AdminServiceAdminLoginRequest) (*pb.AdminServiceAdminLoginResponse, error) {
	resp, err := s.adminService.Login(ctx, &dto.AdminLoginRequest{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	md := metadata.Pairs(
		middleware.MDHeaderAccessToken, resp.AccessToken,
		middleware.MDHeaderRefreshToken, resp.RefreshToken,
	)
	if err := grpc.SetHeader(ctx, md); err != nil {
		return nil, err
	}
	return &pb.AdminServiceAdminLoginResponse{
		Admin:        s.converter.ToPB(resp.Admin),
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
	}, nil
}

// RefreshToken 刷新令牌
func (s *Service) RefreshToken(ctx context.Context, req *pb.AdminServiceRefreshTokenRequest) (*pb.AdminServiceRefreshTokenResponse, error) {
	access, refresh, err := s.adminService.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	md := metadata.Pairs(
		middleware.MDHeaderAccessToken, access,
		middleware.MDHeaderRefreshToken, refresh,
	)
	if err := grpc.SetHeader(ctx, md); err != nil {
		return nil, err
	}
	return &pb.AdminServiceRefreshTokenResponse{
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil
}

// GetCurrentAdminInfo 获取当前管理员信息
func (s *Service) GetCurrentAdminInfo(ctx context.Context, req *pb.AdminServiceGetCurrentAdminInfoRequest) (*pb.AdminServiceGetCurrentAdminInfoResponse, error) {
	a, err := s.adminService.GetCurrentAdminInfo(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.AdminServiceGetCurrentAdminInfoResponse{
		Admin: s.converter.ToPB(a),
	}, nil
}
