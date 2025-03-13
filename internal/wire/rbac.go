package wire

import (
	"github.com/google/wire"
	domainsecurity "github.com/lazyjean/sla2/internal/domain/security"
	"github.com/lazyjean/sla2/internal/interfaces/grpc/middleware"
)

// RBAC权限集
var rbacSet = wire.NewSet(
	domainsecurity.NewRBACProvider,
	providePermissionManager,
	providePermissionHelper,
	middleware.NewRBACInterceptor,
)

// providePermissionManager 提供权限管理器
func providePermissionManager(provider *domainsecurity.RBACProvider) domainsecurity.PermissionManager {
	return provider.GetPermissionManager()
}

// providePermissionHelper 提供权限辅助工具
func providePermissionHelper(provider *domainsecurity.RBACProvider) *domainsecurity.PermissionHelper {
	return provider.GetPermissionHelper()
}
