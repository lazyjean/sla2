package admin

import (
	"context"

	pb "github.com/lazyjean/sla2/api/proto/v1"
	"github.com/lazyjean/sla2/internal/application/dto"
	"github.com/lazyjean/sla2/internal/application/service"
	"github.com/lazyjean/sla2/internal/transport/grpc/admin/converter"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AdminService 管理员服务
type AdminService struct {
	pb.UnimplementedAdminServiceServer
	adminService *service.AdminService
	converter    *converter.AdminConverter
}

// NewAdminService 创建新的管理员服务
func NewAdminService(adminService *service.AdminService) *AdminService {
	return &AdminService{
		adminService: adminService,
		converter:    converter.NewAdminConverter(),
	}
}

// CheckSystemStatus 检查系统状态
func (s *AdminService) CheckSystemStatus(ctx context.Context, req *pb.AdminServiceCheckSystemStatusRequest) (*pb.AdminServiceCheckSystemStatusResponse, error) {
	resp, err := s.adminService.CheckSystemStatus(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.AdminServiceCheckSystemStatusResponse{
		Initialized: resp.Initialized,
	}, nil
}

// InitializeSystem 初始化系统
func (s *AdminService) InitializeSystem(ctx context.Context, req *pb.AdminServiceInitializeSystemRequest) (*pb.AdminServiceInitializeSystemResponse, error) {
	resp, err := s.adminService.InitializeSystem(ctx, &dto.InitializeSystemRequest{
		Username: req.Username,
		Nickname: req.Nickname,
		Password: req.Password,
		Email:    req.Email,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.AdminServiceInitializeSystemResponse{
		Admin:        s.converter.ToPB(resp.Admin),
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
	}, nil
}

// AdminLogin 管理员登录
func (s *AdminService) AdminLogin(ctx context.Context, req *pb.AdminServiceAdminLoginRequest) (*pb.AdminServiceAdminLoginResponse, error) {
	resp, err := s.adminService.Login(ctx, &dto.AdminLoginRequest{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.AdminServiceAdminLoginResponse{
		Admin:        s.converter.ToPB(resp.Admin),
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
	}, nil
}

// RefreshToken 刷新令牌
func (s *AdminService) RefreshToken(ctx context.Context, req *pb.AdminServiceRefreshTokenRequest) (*pb.AdminServiceRefreshTokenResponse, error) {
	resp, err := s.adminService.RefreshToken(ctx, &dto.RefreshTokenRequest{
		RefreshToken: req.RefreshToken,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.AdminServiceRefreshTokenResponse{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
	}, nil
}

// GetCurrentAdminInfo 获取当前管理员信息
func (s *AdminService) GetCurrentAdminInfo(ctx context.Context, req *pb.AdminServiceGetCurrentAdminInfoRequest) (*pb.AdminServiceGetCurrentAdminInfoResponse, error) {
	resp, err := s.adminService.GetCurrentAdminInfo(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.AdminServiceGetCurrentAdminInfoResponse{
		Admin: s.converter.ToPBFromResponse(resp),
	}, nil
}
