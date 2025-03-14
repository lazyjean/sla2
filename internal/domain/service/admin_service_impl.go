package service

import (
	"context"
	"errors"

	"github.com/lazyjean/sla2/internal/domain/entity"
	"github.com/lazyjean/sla2/internal/domain/repository"
	"github.com/lazyjean/sla2/internal/domain/security"
)

// Domain errors
var (
	ErrSystemAlreadyInitialized = errors.New("system already initialized")
	ErrSystemNotInitialized     = errors.New("system not initialized")
	ErrAdminNotFound            = errors.New("admin not found")
)

// adminService 管理员领域服务实现
type adminService struct {
	adminRepo        repository.AdminRepository
	permissionHelper *security.PermissionHelper

	// 系统初始化状态缓存（一旦初始化就永远为 true）
	isInitialized bool
}

// NewAdminService 创建管理员领域服务
func NewAdminService(
	adminRepo repository.AdminRepository,
	permissionHelper *security.PermissionHelper,
) AdminService {
	svc := &adminService{
		adminRepo:        adminRepo,
		permissionHelper: permissionHelper,
	}

	// 在服务创建时检查系统初始化状态
	// 使用 context.Background() 作为初始检查的上下文
	initialized, err := adminRepo.IsInitialized(context.Background())
	if err == nil && initialized {
		svc.isInitialized = true
	}

	return svc
}

// IsSystemInitialized 检查系统是否已初始化
func (s *adminService) IsSystemInitialized(ctx context.Context) (bool, error) {
	// 如果内存标记为已初始化，直接返回 true
	if s.isInitialized {
		return true, nil
	}

	// 查询数据库
	initialized, err := s.adminRepo.IsInitialized(ctx)
	if err != nil {
		return false, err
	}

	// 如果已初始化，更新内存状态
	if initialized {
		s.isInitialized = true
	}

	return initialized, nil
}

// InitializeSystem 初始化系统并创建初始管理员
func (s *adminService) InitializeSystem(ctx context.Context, admin *entity.Admin) error {
	// 检查系统是否已初始化
	initialized, err := s.IsSystemInitialized(ctx)
	if err != nil {
		return err
	}
	if initialized {
		return ErrSystemAlreadyInitialized
	}

	// 开始事务
	// 注意：实际实现可能需要根据具体的事务管理方式进行调整
	// 这里假设仓储层已实现事务支持

	// 1. 保存管理员
	if err := s.adminRepo.Create(ctx, admin); err != nil {
		return err
	}

	// 2. 为新创建的管理员分配admin角色
	// PermissionHelper.AssignRoleToUser 返回 (bool, error)
	assigned, err := s.permissionHelper.AssignRoleToUser(ctx, admin.ID, security.RoleAdmin)
	if err != nil {
		// 回滚管理员创建
		_ = s.adminRepo.Delete(ctx, admin.ID)
		return err
	}

	if !assigned {
		// 如果角色分配失败但没有错误，记录一个域错误
		_ = s.adminRepo.Delete(ctx, admin.ID)
		return NewDomainError("failed to assign admin role to admin user")
	}

	// 3. 为admin角色设置管理员权限已在PermissionInitializer中配置
	// 根据permission_initializer.go的代码，admin角色已经被赋予所有权限
	// fmt.Sprintf("r:%s", RoleAdmin), ResourceAny, ActionAny

	// 更新初始化状态缓存
	s.isInitialized = true

	return nil
}

// GetAdminByID 根据ID获取管理员信息
func (s *adminService) GetAdminByID(ctx context.Context, adminID entity.UID) (*entity.Admin, error) {
	if !s.isInitialized {
		return nil, ErrSystemNotInitialized
	}
	admin, err := s.adminRepo.FindByID(ctx, adminID)
	if err != nil {
		return nil, err
	}
	if admin == nil {
		return nil, ErrAdminNotFound
	}
	return admin, nil
}

// GetAdminByUsername 根据用户名获取管理员信息
func (s *adminService) GetAdminByUsername(ctx context.Context, username string) (*entity.Admin, error) {
	if !s.isInitialized {
		return nil, ErrSystemNotInitialized
	}
	admin, err := s.adminRepo.FindByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	if admin == nil {
		return nil, ErrAdminNotFound
	}
	return admin, nil
}

// DomainError 领域错误
type DomainError struct {
	Message string
}

// NewDomainError 创建领域错误
func NewDomainError(message string) *DomainError {
	return &DomainError{
		Message: message,
	}
}

// Error 实现 error 接口
func (e *DomainError) Error() string {
	return e.Message
}
