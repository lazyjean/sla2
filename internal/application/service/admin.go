package service

import (
	"context"

	"github.com/lazyjean/sla2/internal/application/dto"
	"github.com/lazyjean/sla2/internal/domain/entity"
	errorsDomain "github.com/lazyjean/sla2/internal/domain/errors"
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

// NewAdminService 创建带权限助手的管理员服务
func NewAdminService(
	adminRepo repository.AdminRepository,
	passwordService security.PasswordService,
	tokenService security.TokenService,
	permissionHelper *security.PermissionHelper,
) *AdminService {
	return &AdminService{
		adminService:    domainService.NewAdminService(adminRepo, permissionHelper),
		passwordService: passwordService,
		tokenService:    tokenService,
	}
}

// IsSystemInitialized 检查系统状态
func (s *AdminService) IsSystemInitialized(ctx context.Context) (bool, error) {
	return s.adminService.IsSystemInitialized(ctx)
}

// InitializeSystem 初始化系统
func (s *AdminService) InitializeSystem(ctx context.Context, req *dto.InitializeSystemRequest) (*dto.InitializeSystemResponse, error) {
	// 检查系统是否已初始化
	initialized, err := s.adminService.IsSystemInitialized(ctx)
	if err != nil {
		return nil, err
	}
	if initialized {
		return nil, errorsDomain.ErrSystemAlreadyInitialized
	}

	// 哈希密码（应用层职责）
	hashedPassword, err := s.passwordService.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// 创建管理员实体（领域层职责）
	admin := entity.NewAdmin(req.Username, hashedPassword, req.Nickname, req.Email)

	// 调用领域服务初始化系统
	err = s.adminService.InitializeSystem(ctx, admin)
	if err != nil {
		return nil, err
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

	return &dto.InitializeSystemResponse{
		Admin:        admin,
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
		return nil, errorsDomain.ErrInvalidPassword
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
		Admin:        admin,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// RefreshToken 刷新访问令牌
func (s *AdminService) RefreshToken(ctx context.Context, token string) (string, string, error) {
	var uid entity.UID
	var err error
	if uid, err = GetUserID(ctx); err != nil {
		return "", "", err
	}

	if id, _, err := s.tokenService.ValidateToken(token); err != nil {
		return "", "", err
	} else if id != uid {
		return "", "", errorsDomain.ErrRefreshTokenMismatch
	}
	a, err := s.adminService.GetAdminByID(ctx, uid)
	if err != nil {
		return "", "", err
	}

	// 生成新的令牌（应用层职责）
	accessToken, err := s.tokenService.GenerateToken(uid, a.Roles)
	if err != nil {
		return "", "", err
	}
	refreshToken, err := s.tokenService.GenerateRefreshToken(uid, a.Roles)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// GetCurrentAdminInfo 获取当前管理员信息
func (s *AdminService) GetCurrentAdminInfo(ctx context.Context) (*entity.Admin, error) {
	userID, err := GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	// 获取管理员信息
	return s.adminService.GetAdminByID(ctx, userID)
}
