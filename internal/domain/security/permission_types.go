package security

// 定义资源类型常量
const (
	// ResourceUser 用户资源
	ResourceUser = "user"
	// ResourceRole 角色资源
	ResourceRole = "role"
	// ResourcePermission 权限资源
	ResourcePermission = "permission"
	// ResourceUserRole 用户角色关联资源
	ResourceUserRole = "user_role"
	// ResourceRolePermission 角色权限关联资源
	ResourceRolePermission = "role_permission"
	// ResourceCourse 课程资源
	ResourceCourse = "course"
	// ResourceQuestion 问题资源
	ResourceQuestion = "question"
	// ResourceQuestionTag 问题标签资源
	ResourceQuestionTag = "question_tag"
	// ResourceWord 单词资源
	ResourceWord = "word"
	// ResourceAny 任意资源，通配符
	ResourceAny = "*"
)

// 定义操作类型常量
const (
	// ActionCreate 创建操作
	ActionCreate = "create"
	// ActionRead 读取操作
	ActionRead = "read"
	// ActionUpdate 更新操作
	ActionUpdate = "update"
	// ActionDelete 删除操作
	ActionDelete = "delete"
	// ActionList 列表操作
	ActionList = "list"
	// ActionAssign 分配操作（如分配角色）
	ActionAssign = "assign"
	// ActionAny 任意操作，通配符
	ActionAny = "*"
)

// 定义预设角色
const (
	// RoleAdmin 管理员角色
	RoleAdmin = "admin"
	// RoleUser 普通用户角色
	RoleUser = "user"
	// RoleGuest 访客角色
	RoleGuest = "guest"
	// RoleContentManager 内容管理员角色
	RoleContentManager = "content_manager"
	// RoleUserManager 用户管理员角色
	RoleUserManager = "user_manager"
)
