package admin

import (
	"context"

	pb "github.com/lazyjean/sla2/api/proto/v1"
	"github.com/lazyjean/sla2/internal/application/dto"
	"github.com/lazyjean/sla2/internal/application/service"
)

// AdminService 管理员 gRPC 服务
type AdminService struct {
	pb.UnimplementedAdminServiceServer
	adminService *service.AdminService
}

// NewAdminService 创建管理员 gRPC 服务
func NewAdminService(adminService *service.AdminService) *AdminService {
	return &AdminService{
		adminService: adminService,
	}
}

// CheckSystemStatus 检查系统状态
func (s *AdminService) CheckSystemStatus(ctx context.Context, req *pb.AdminServiceCheckSystemStatusRequest) (*pb.AdminServiceCheckSystemStatusResponse, error) {
	resp, err := s.adminService.CheckSystemStatus(ctx)
	if err != nil {
		return nil, err
	}

	return &pb.AdminServiceCheckSystemStatusResponse{
		Initialized: resp.Initialized,
	}, nil
}

// InitializeSystem 初始化系统
func (s *AdminService) InitializeSystem(ctx context.Context, req *pb.AdminServiceInitializeSystemRequest) (*pb.AdminServiceInitializeSystemResponse, error) {
	resp, err := s.adminService.InitializeSystem(ctx, &dto.InitializeSystemRequest{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		return nil, err
	}

	return &pb.AdminServiceInitializeSystemResponse{
		Admin: &pb.AdminInfo{
			Id:        uint64(resp.Admin.ID),
			Username:  resp.Admin.Username,
			Nickname:  resp.Admin.Nickname,
			Roles:     resp.Admin.Roles,
			CreatedAt: resp.Admin.CreatedAt.Unix(),
			UpdatedAt: resp.Admin.UpdatedAt.Unix(),
		},
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
		return nil, err
	}

	return &pb.AdminServiceAdminLoginResponse{
		Admin: &pb.AdminInfo{
			Id:        uint64(resp.Admin.ID),
			Username:  resp.Admin.Username,
			Nickname:  resp.Admin.Nickname,
			Roles:     resp.Admin.Roles,
			CreatedAt: resp.Admin.CreatedAt.Unix(),
			UpdatedAt: resp.Admin.UpdatedAt.Unix(),
		},
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
	}, nil
}

// RefreshToken 刷新访问令牌
func (s *AdminService) RefreshToken(ctx context.Context, req *pb.AdminServiceRefreshTokenRequest) (*pb.AdminServiceRefreshTokenResponse, error) {
	resp, err := s.adminService.RefreshToken(ctx, &dto.RefreshTokenRequest{
		RefreshToken: req.RefreshToken,
	})
	if err != nil {
		return nil, err
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
		return nil, err
	}

	return &pb.AdminServiceGetCurrentAdminInfoResponse{
		Admin: &pb.AdminInfo{
			Id:       uint64(resp.ID),
			Username: resp.Username,
			Nickname: resp.Nickname,
			Roles:    resp.Roles,
		},
	}, nil
}
