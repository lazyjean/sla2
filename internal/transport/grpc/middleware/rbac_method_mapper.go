package middleware

import (
	"context"

	"github.com/lazyjean/sla2/internal/domain/security"
	"github.com/lazyjean/sla2/pkg/logger"
)

// 定义额外的资源类型常量
const (
	// ResourceSystem 系统资源
	ResourceSystem = "system"
)

// RegisterRBACMethodMappings 注册gRPC方法与权限的映射
func RegisterRBACMethodMappings(ctx context.Context, interceptor *RBACInterceptor) {
	logger.GetLogger(ctx).Info("Registering RBAC method mappings...")

	// 注册认证和登录相关API到白名单
	registerAuthWhitelist(ctx, interceptor)

	// 注册用户相关API权限
	registerUserMethodPermissions(ctx, interceptor)

	// 注册角色相关API权限
	registerRoleMethodPermissions(ctx, interceptor)

	// 注册课程相关API权限
	registerCourseMethodPermissions(ctx, interceptor)

	// 注册问题相关API权限
	registerQuestionMethodPermissions(ctx, interceptor)

	// 注册单词相关API权限
	registerWordMethodPermissions(ctx, interceptor)

	// 注册健康检查API到白名单
	registerHealthCheckWhitelist(ctx, interceptor)

	// 注册WebSocket相关API权限
	registerWebSocketPermissions(ctx, interceptor)
}

// registerAuthWhitelist 注册认证相关API到白名单（不需要检查权限）
func registerAuthWhitelist(ctx context.Context, interceptor *RBACInterceptor) {
	// 系统管理接口
	interceptor.AddToWhitelist(ctx, "/proto.v1.AdminService/IsSystemInitialized")
	interceptor.AddToWhitelist(ctx, "/proto.v1.AdminService/InitializeSystem")

	// 管理员认证接口
	interceptor.AddToWhitelist(ctx, "/proto.v1.AdminService/AdminLogin")
	interceptor.AddToWhitelist(ctx, "/proto.v1.AdminService/RefreshToken")

	// 用户认证接口
	interceptor.AddToWhitelist(ctx, "/proto.v1.UserService/Register")
	interceptor.AddToWhitelist(ctx, "/proto.v1.UserService/Login")
	interceptor.AddToWhitelist(ctx, "/proto.v1.UserService/RefreshToken")
	interceptor.AddToWhitelist(ctx, "/proto.v1.UserService/AppleLogin")
	interceptor.AddToWhitelist(ctx, "/proto.v1.UserService/ResetPassword")

	// 健康检查接口
	interceptor.AddToWhitelist(ctx, "/grpc.health.v1.Health/Check")
	interceptor.AddToWhitelist(ctx, "/grpc.health.v1.Health/Watch")
}

// registerUserMethodPermissions 注册用户相关API权限
func registerUserMethodPermissions(ctx context.Context, interceptor *RBACInterceptor) {
	// 用户管理相关API
	// 1. 用户资源的读取权限
	interceptor.RegisterMethodPermission(
		ctx,
		"/proto.v1.UserService/GetUser",
		security.ResourceUser,
		security.ActionRead,
	)
	interceptor.RegisterMethodPermission(
		ctx,
		"/proto.v1.UserService/ListUsers",
		security.ResourceUser,
		security.ActionList,
	)

	// 2. 用户资源的写入权限
	interceptor.RegisterMethodPermission(
		ctx,
		"/proto.v1.UserService/UpdateUser",
		security.ResourceUser,
		security.ActionUpdate,
	)
	interceptor.RegisterMethodPermission(
		ctx,
		"/proto.v1.UserService/DeleteUser",
		security.ResourceUser,
		security.ActionDelete,
	)

	// 3. 用户角色相关权限
	interceptor.RegisterMethodPermission(
		ctx,
		"/proto.v1.UserService/AssignRole",
		security.ResourceUserRole,
		security.ActionCreate,
	)
	interceptor.RegisterMethodPermission(
		ctx,
		"/proto.v1.UserService/RevokeRole",
		security.ResourceUserRole,
		security.ActionDelete,
	)
}

// registerRoleMethodPermissions 注册角色相关API权限
func registerRoleMethodPermissions(ctx context.Context, interceptor *RBACInterceptor) {
	// 1. 角色资源的读取权限
	interceptor.RegisterMethodPermission(
		ctx,
		"/proto.v1.AdminService/GetRole",
		security.ResourceRole,
		security.ActionRead,
	)
	interceptor.RegisterMethodPermission(
		ctx,
		"/proto.v1.AdminService/ListRoles",
		security.ResourceRole,
		security.ActionList,
	)

	// 2. 角色资源的写入权限
	interceptor.RegisterMethodPermission(
		ctx,
		"/proto.v1.AdminService/CreateRole",
		security.ResourceRole,
		security.ActionCreate,
	)
	interceptor.RegisterMethodPermission(
		ctx,
		"/proto.v1.AdminService/UpdateRole",
		security.ResourceRole,
		security.ActionUpdate,
	)
	interceptor.RegisterMethodPermission(
		ctx,
		"/proto.v1.AdminService/DeleteRole",
		security.ResourceRole,
		security.ActionDelete,
	)

	// 3. 角色权限相关API
	interceptor.RegisterMethodPermission(
		ctx,
		"/proto.v1.AdminService/AssignPermission",
		security.ResourceRolePermission,
		security.ActionCreate,
	)
	interceptor.RegisterMethodPermission(
		ctx,
		"/proto.v1.AdminService/RevokePermission",
		security.ResourceRolePermission,
		security.ActionDelete,
	)
	interceptor.RegisterMethodPermission(
		ctx,
		"/proto.v1.AdminService/ListRolePermissions",
		security.ResourceRolePermission,
		security.ActionList,
	)
}

// registerCourseMethodPermissions 注册课程相关API权限
func registerCourseMethodPermissions(ctx context.Context, interceptor *RBACInterceptor) {
	// 1. 课程资源的读取权限
	interceptor.RegisterMethodPermission(
		ctx,
		"/proto.v1.CourseService/Get",
		security.ResourceCourse,
		security.ActionRead,
	)
	interceptor.RegisterMethodPermission(
		ctx,
		"/proto.v1.CourseService/List",
		security.ResourceCourse,
		security.ActionList,
	)
	interceptor.RegisterMethodPermission(
		ctx,
		"/proto.v1.CourseService/Search",
		security.ResourceCourse,
		security.ActionList,
	)

	// 2. 课程资源的写入权限
	interceptor.RegisterMethodPermission(
		ctx,
		"/proto.v1.CourseService/Create",
		security.ResourceCourse,
		security.ActionCreate,
	)
	interceptor.RegisterMethodPermission(
		ctx,
		"/proto.v1.CourseService/BatchCreate",
		security.ResourceCourse,
		security.ActionCreate,
	)
	interceptor.RegisterMethodPermission(
		ctx,
		"/proto.v1.CourseService/Update",
		security.ResourceCourse,
		security.ActionUpdate,
	)
	interceptor.RegisterMethodPermission(
		ctx,
		"/proto.v1.CourseService/Delete",
		security.ResourceCourse,
		security.ActionDelete,
	)
}

// registerQuestionMethodPermissions 注册问题相关API权限
func registerQuestionMethodPermissions(ctx context.Context, interceptor *RBACInterceptor) {
	// 1. 问题资源的读取权限
	interceptor.RegisterMethodPermission(
		ctx,
		"/proto.v1.QuestionService/Get",
		security.ResourceQuestion,
		security.ActionRead,
	)
	interceptor.RegisterMethodPermission(
		ctx,
		"/proto.v1.QuestionService/List",
		security.ResourceQuestion,
		security.ActionList,
	)
	interceptor.RegisterMethodPermission(
		ctx,
		"/proto.v1.QuestionService/Search",
		security.ResourceQuestion,
		security.ActionList,
	)

	// 2. 问题资源的写入权限
	interceptor.RegisterMethodPermission(
		ctx,
		"/proto.v1.QuestionService/Create",
		security.ResourceQuestion,
		security.ActionCreate,
	)
	interceptor.RegisterMethodPermission(
		ctx,
		"/proto.v1.QuestionService/Update",
		security.ResourceQuestion,
		security.ActionUpdate,
	)
	interceptor.RegisterMethodPermission(
		ctx,
		"/proto.v1.QuestionService/Delete",
		security.ResourceQuestion,
		security.ActionDelete,
	)

	// 3. 问题标签相关API
	interceptor.RegisterMethodPermission(
		ctx,
		"/proto.v1.QuestionService/ListTags",
		security.ResourceQuestion,
		security.ActionList,
	)
}

// registerWordMethodPermissions 注册单词相关API权限
func registerWordMethodPermissions(ctx context.Context, interceptor *RBACInterceptor) {
	// 1. 单词资源的读取权限
	interceptor.RegisterMethodPermission(
		ctx,
		"/proto.v1.WordService/Get",
		security.ResourceWord,
		security.ActionRead,
	)
	interceptor.RegisterMethodPermission(
		ctx,
		"/proto.v1.WordService/List",
		security.ResourceWord,
		security.ActionList,
	)
	interceptor.RegisterMethodPermission(
		ctx,
		"/proto.v1.WordService/Search",
		security.ResourceWord,
		security.ActionList,
	)

	// 2. 单词资源的写入权限
	interceptor.RegisterMethodPermission(
		ctx,
		"/proto.v1.WordService/Create",
		security.ResourceWord,
		security.ActionCreate,
	)
	interceptor.RegisterMethodPermission(
		ctx,
		"/proto.v1.WordService/Update",
		security.ResourceWord,
		security.ActionUpdate,
	)
	interceptor.RegisterMethodPermission(
		ctx,
		"/proto.v1.WordService/Delete",
		security.ResourceWord,
		security.ActionDelete,
	)
}

// registerHealthCheckWhitelist 注册健康检查API到白名单
func registerHealthCheckWhitelist(ctx context.Context, interceptor *RBACInterceptor) {
	// 健康检查接口不需要权限检查
	interceptor.AddToWhitelist(ctx, "/grpc.health.v1.Health/Check")
	interceptor.AddToWhitelist(ctx, "/grpc.health.v1.Health/Watch")
}

// registerWebSocketPermissions 注册WebSocket相关API权限
func registerWebSocketPermissions(ctx context.Context, interceptor *RBACInterceptor) {
	// WebSocket相关API暂时加入白名单，不需要权限验证
	interceptor.AddToWhitelist(ctx, "/proto.v1.ChatService/Chat")
}
