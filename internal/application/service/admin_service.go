package service

import (
	"context"
	"errors"

	"github.com/lazyjean/sla2/internal/application/dto"
	"github.com/lazyjean/sla2/internal/domain/entity"
	"github.com/lazyjean/sla2/internal/domain/repository"
	"github.com/lazyjean/sla2/internal/domain/security"
	domainService "github.com/lazyjean/sla2/internal/domain/service"
)

// AdminService 管理员服务
type AdminService struct {
	adminService    domainService.AdminService
	passwordService security.PasswordService
	tokenService    security.TokenService
}

// NewAdminService 创建管理员服务
func NewAdminService(
	adminRepo repository.AdminRepository,
	passwordService security.PasswordService,
	tokenService security.TokenService,
) *AdminService {
	return &AdminService{
		adminService:    domainService.NewAdminService(adminRepo),
		passwordService: passwordService,
		tokenService:    tokenService,
	}
}

// CheckSystemStatus 检查系统状态
func (s *AdminService) CheckSystemStatus(ctx context.Context) (*dto.SystemStatusResponse, error) {
	initialized, err := s.adminService.IsSystemInitialized(ctx)
	if err != nil {
		return nil, err
	}

	return &dto.SystemStatusResponse{
		Initialized: initialized,
	}, nil
}

// InitializeSystem 初始化系统
func (s *AdminService) InitializeSystem(ctx context.Context, req *dto.InitializeSystemRequest) (*dto.InitializeSystemResponse, error) {
	// 检查系统是否已初始化
	initialized, err := s.adminService.IsSystemInitialized(ctx)
	if err != nil {
		return nil, err
	}
	if initialized {
		return nil, errors.New("system already initialized")
	}

	// 哈希密码（应用层职责）
	hashedPassword, err := s.passwordService.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// 创建管理员
	admin := entity.NewAdmin(req.Username, hashedPassword, req.Nickname)
	// 设置管理员权限
	admin.Roles = []string{"admin"}

	// 调用领域服务初始化系统
	err = s.adminService.InitializeSystem(ctx, admin)
	if err != nil {
		return nil, err
	}

	// 生成令牌（应用层职责）
	accessToken, err := s.tokenService.GenerateToken(admin.ID, []string{"admin"})
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.tokenService.GenerateRefreshToken(admin.ID, []string{"admin"})
	if err != nil {
		return nil, err
	}

	return &dto.InitializeSystemResponse{
		Admin: &dto.AdminInfo{
			ID:        admin.ID,
			Username:  admin.Username,
			Nickname:  admin.Nickname,
			Roles:     admin.Roles,
			CreatedAt: admin.CreatedAt,
			UpdatedAt: admin.UpdatedAt,
		},
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// Login 管理员登录
func (s *AdminService) Login(ctx context.Context, req *dto.AdminLoginRequest) (*dto.AdminLoginResponse, error) {
	// 获取管理员信息
	admin, err := s.adminService.GetAdminByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}

	// 验证密码（应用层职责）
	if !s.passwordService.VerifyPassword(req.Password, admin.Password) {
		return nil, errors.New("invalid credentials")
	}

	// 生成令牌（应用层职责）
	accessToken, err := s.tokenService.GenerateToken(admin.ID, admin.Roles)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.tokenService.GenerateRefreshToken(admin.ID, admin.Roles)
	if err != nil {
		return nil, err
	}

	return &dto.AdminLoginResponse{
		Admin: &dto.AdminInfo{
			ID:        admin.ID,
			Username:  admin.Username,
			Nickname:  admin.Nickname,
			Roles:     admin.Roles,
			CreatedAt: admin.CreatedAt,
			UpdatedAt: admin.UpdatedAt,
		},
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// RefreshToken 刷新访问令牌
func (s *AdminService) RefreshToken(ctx context.Context, req *dto.RefreshTokenRequest) (*dto.AdminRefreshTokenResponse, error) {
	// 验证刷新令牌（应用层职责）
	adminID, roles, err := s.tokenService.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, err
	}

	// 检查管理员是否存在
	_, err = s.adminService.GetAdminByID(ctx, adminID)
	if err != nil {
		return nil, err
	}

	// 生成新的令牌（应用层职责）
	accessToken, err := s.tokenService.GenerateToken(adminID, roles)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.tokenService.GenerateRefreshToken(adminID, roles)
	if err != nil {
		return nil, err
	}

	return &dto.AdminRefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// GetCurrentAdminInfo 获取当前管理员信息
func (s *AdminService) GetCurrentAdminInfo(ctx context.Context) (*dto.AdminInfoResponse, error) {
	userID, err := GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	// 获取管理员信息
	admin, err := s.adminService.GetAdminByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &dto.AdminInfoResponse{
		ID:       admin.ID,
		Username: admin.Username,
		Nickname: admin.Nickname,
		Roles:    admin.Roles,
	}, nil
}

// 错误定义
var (
	ErrSystemAlreadyInitialized = errors.New("system already initialized")
	ErrInvalidCredentials       = errors.New("invalid credentials")
	ErrAdminNotFound            = errors.New("admin not found")
	ErrUnauthorized             = errors.New("unauthorized")
	ErrRoleNotFound             = errors.New("role not found")
)

// GetAdminIDFromContext 从上下文中获取管理员ID
func GetAdminIDFromContext(ctx context.Context) (entity.UID, error) {
	userID, err := GetUserID(ctx)
	if err != nil {
		return 0, err
	}
	// todo: 当前只有admin角色，可以访问管理端接口
	if !HasAnyRole(ctx, "admin") {
		return 0, ErrRoleNotFound
	}
	return userID, nil
}
