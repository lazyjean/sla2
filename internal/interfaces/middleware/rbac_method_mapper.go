package middleware

import (
	"github.com/lazyjean/sla2/internal/domain/security"
	"github.com/lazyjean/sla2/pkg/logger"
)

// 定义额外的资源类型常量
const (
	// ResourceSystem 系统资源
	ResourceSystem = "system"
)

// RegisterRBACMethodMappings 注册gRPC方法与权限的映射
func RegisterRBACMethodMappings(interceptor *RBACInterceptor) {
	logger.Log.Info("Registering RBAC method mappings...")

	// 注册认证和登录相关API到白名单
	registerAuthWhitelist(interceptor)

	// 注册用户相关API权限
	registerUserMethodPermissions(interceptor)

	// 注册角色相关API权限
	registerRoleMethodPermissions(interceptor)

	// 注册课程相关API权限
	registerCourseMethodPermissions(interceptor)

	// 注册问题相关API权限
	registerQuestionMethodPermissions(interceptor)

	// 注册单词相关API权限
	registerWordMethodPermissions(interceptor)

	// 注册健康检查API到白名单
	registerHealthCheckWhitelist(interceptor)

	// 注册WebSocket相关API权限
	registerWebSocketPermissions(interceptor)

	logger.Log.Info("RBAC method mappings registered successfully")
}

// registerAuthWhitelist 注册认证和登录相关API到白名单
func registerAuthWhitelist(interceptor *RBACInterceptor) {
	// 添加认证相关方法到白名单
	authMethods := []string{
		"/auth.AuthService/Login",
		"/auth.AuthService/Register",
		"/auth.AuthService/RefreshToken",
		"/auth.AuthService/LoginWithApple",
		"/health.HealthService/Check",
	}

	for _, method := range authMethods {
		interceptor.AddToWhitelist(method)
	}
}

// registerUserMethodPermissions 注册用户相关API权限
func registerUserMethodPermissions(interceptor *RBACInterceptor) {
	// 用户相关方法权限映射
	userMethods := map[string]map[string]string{
		"/user.UserService/GetUser": {
			"resource": security.ResourceUser,
			"action":   security.ActionRead,
		},
		"/user.UserService/GetUserProfile": {
			"resource": security.ResourceUser,
			"action":   security.ActionRead,
		},
		"/user.UserService/UpdateUserProfile": {
			"resource": security.ResourceUser,
			"action":   security.ActionUpdate,
		},
		"/user.UserService/ListUsers": {
			"resource": security.ResourceUser,
			"action":   security.ActionList,
		},
		"/user.UserService/DeleteUser": {
			"resource": security.ResourceUser,
			"action":   security.ActionDelete,
		},
	}

	for method, permission := range userMethods {
		interceptor.RegisterMethodPermission(method, permission["resource"], permission["action"])
	}
}

// registerRoleMethodPermissions 注册角色相关API权限
func registerRoleMethodPermissions(interceptor *RBACInterceptor) {
	// 角色相关方法权限映射
	roleMethods := map[string]map[string]string{
		"/role.RoleService/GetRole": {
			"resource": security.ResourceRole,
			"action":   security.ActionRead,
		},
		"/role.RoleService/ListRoles": {
			"resource": security.ResourceRole,
			"action":   security.ActionList,
		},
		"/role.RoleService/CreateRole": {
			"resource": security.ResourceRole,
			"action":   security.ActionCreate,
		},
		"/role.RoleService/UpdateRole": {
			"resource": security.ResourceRole,
			"action":   security.ActionUpdate,
		},
		"/role.RoleService/DeleteRole": {
			"resource": security.ResourceRole,
			"action":   security.ActionDelete,
		},
		"/role.RoleService/AssignRoleToUser": {
			"resource": security.ResourceRole,
			"action":   security.ActionAssign,
		},
		"/role.RoleService/RevokeRoleFromUser": {
			"resource": security.ResourceRole,
			"action":   security.ActionDelete,
		},
	}

	for method, permission := range roleMethods {
		interceptor.RegisterMethodPermission(method, permission["resource"], permission["action"])
	}
}

// registerCourseMethodPermissions 注册课程相关API权限
func registerCourseMethodPermissions(interceptor *RBACInterceptor) {
	// 课程相关方法权限映射
	courseMethods := map[string]map[string]string{
		"/course.CourseService/GetCourse": {
			"resource": security.ResourceCourse,
			"action":   security.ActionRead,
		},
		"/course.CourseService/ListCourses": {
			"resource": security.ResourceCourse,
			"action":   security.ActionList,
		},
		"/course.CourseService/CreateCourse": {
			"resource": security.ResourceCourse,
			"action":   security.ActionCreate,
		},
		"/course.CourseService/UpdateCourse": {
			"resource": security.ResourceCourse,
			"action":   security.ActionUpdate,
		},
		"/course.CourseService/DeleteCourse": {
			"resource": security.ResourceCourse,
			"action":   security.ActionDelete,
		},
	}

	for method, permission := range courseMethods {
		interceptor.RegisterMethodPermission(method, permission["resource"], permission["action"])
	}
}

// registerQuestionMethodPermissions 注册问题相关API权限
func registerQuestionMethodPermissions(interceptor *RBACInterceptor) {
	// 问题相关方法权限映射
	questionMethods := map[string]map[string]string{
		"/question.QuestionService/GetQuestion": {
			"resource": security.ResourceQuestion,
			"action":   security.ActionRead,
		},
		"/question.QuestionService/ListQuestions": {
			"resource": security.ResourceQuestion,
			"action":   security.ActionList,
		},
		"/question.QuestionService/CreateQuestion": {
			"resource": security.ResourceQuestion,
			"action":   security.ActionCreate,
		},
		"/question.QuestionService/UpdateQuestion": {
			"resource": security.ResourceQuestion,
			"action":   security.ActionUpdate,
		},
		"/question.QuestionService/DeleteQuestion": {
			"resource": security.ResourceQuestion,
			"action":   security.ActionDelete,
		},
	}

	for method, permission := range questionMethods {
		interceptor.RegisterMethodPermission(method, permission["resource"], permission["action"])
	}
}

// registerWordMethodPermissions 注册单词相关API权限
func registerWordMethodPermissions(interceptor *RBACInterceptor) {
	// 单词相关方法权限映射
	wordMethods := map[string]map[string]string{
		"/word.WordService/GetWord": {
			"resource": security.ResourceWord,
			"action":   security.ActionRead,
		},
		"/word.WordService/ListWords": {
			"resource": security.ResourceWord,
			"action":   security.ActionList,
		},
		"/word.WordService/CreateWord": {
			"resource": security.ResourceWord,
			"action":   security.ActionCreate,
		},
		"/word.WordService/UpdateWord": {
			"resource": security.ResourceWord,
			"action":   security.ActionUpdate,
		},
		"/word.WordService/DeleteWord": {
			"resource": security.ResourceWord,
			"action":   security.ActionDelete,
		},
	}

	for method, permission := range wordMethods {
		interceptor.RegisterMethodPermission(method, permission["resource"], permission["action"])
	}
}

// registerHealthCheckWhitelist 注册健康检查相关API到白名单
func registerHealthCheckWhitelist(interceptor *RBACInterceptor) {
	// 添加健康检查相关方法到白名单
	healthMethods := []string{
		"/grpc.health.v1.Health/Check",
		"/health.HealthService/Check",
	}

	for _, method := range healthMethods {
		interceptor.AddToWhitelist(method)
	}
}

// registerWebSocketPermissions 注册WebSocket相关API权限
func registerWebSocketPermissions(interceptor *RBACInterceptor) {
	// WebSocket相关方法权限映射
	wsMethods := map[string]map[string]string{
		"/ws.WebSocketService/Connect": {
			"resource": ResourceSystem,
			"action":   security.ActionRead,
		},
	}

	for method, permission := range wsMethods {
		interceptor.RegisterMethodPermission(method, permission["resource"], permission["action"])
	}
}
