package wire

import (
	domainsecurity "github.com/lazyjean/sla2/internal/domain/security"
)

// RBAC权限集旧实现
// 现在使用 wire.go 中的新实现
// var rbacSet = wire.NewSet(
// 	domainsecurity.NewRBACProvider,
// 	providePermissionManager,
// 	ProvidePermissionHelper,
// 	middleware.NewRBACInterceptor,
// )

// providePermissionManager 提供权限管理器
func providePermissionManager(provider *domainsecurity.RBACProvider) domainsecurity.PermissionManager {
	return provider.GetPermissionManager()
}
