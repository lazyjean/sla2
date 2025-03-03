package service

// contextKey 上下文键类型
type contextKey string

const (
	// UserIDKey 用户ID的上下文键
	UserIDKey contextKey = "user_id"
	// UserRolesKey 用户角色的上下文键
	UserRolesKey contextKey = "user_roles"
	// AdminIDKey 管理员ID的上下文键
	AdminIDKey contextKey = "admin_id"
)
