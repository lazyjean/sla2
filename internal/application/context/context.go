package context

// ContextKey 上下文键类型
type ContextKey string

const (
	// AdminIDKey 管理员ID的上下文键
	AdminIDKey ContextKey = "admin_id"
	// UserIDKey 用户ID的上下文键
	UserIDKey ContextKey = "user_id"
)
