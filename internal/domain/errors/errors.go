package errors

import "errors"

// Common errors
var (
	ErrInvalidInput = NewError(CodeInvalidInput, "输入错误")
	ErrNotFound     = NewError(CodeNotFound, "资源不存在")
	ErrInvalidWord  = NewError(CodeInvalidWord, "无效的单词")
)

// Word related errors
var (
	ErrEmptyWordText       = NewError(CodeEmptyWordText, "单词不能为空")
	ErrWordNotFound        = NewError(CodeWordNotFound, "单词不存在")
	ErrWordAlreadyExists   = NewError(CodeWordAlreadyExists, "单词已存在")
	ErrInvalidDifficulty   = NewError(CodeInvalidDifficulty, "难度必须在1到5之间")
	ErrInvalidMasteryLevel = NewError(CodeInvalidMasteryLevel, "熟练度必须在0到5之间")
	ErrDuplicateTag        = NewError(CodeDuplicateTag, "标签已存在")
	ErrEmptyDefinition     = NewError(CodeEmptyDefinition, "释义不能为空")
)

// Repository related errors
var (
	ErrFailedToSave   = NewError(CodeFailedToSave, "保存失败")
	ErrFailedToUpdate = NewError(CodeFailedToUpdate, "更新失败")
	ErrFailedToDelete = NewError(CodeFailedToDelete, "删除失败")
	ErrFailedToQuery  = NewError(CodeFailedToQuery, "查询失败")
)

// User related errors
var (
	ErrUserNotFound          = NewError(CodeUserNotFound, "用户不存在")
	ErrUserAlreadyExists     = NewError(CodeUserAlreadyExists, "用户已存在")
	ErrInvalidPhoneFormat    = NewError(CodeInvalidPhoneFormat, "手机号格式错误")
	ErrInvalidEmailFormat    = NewError(CodeInvalidEmailFormat, "邮箱格式错误")
	ErrInvalidUsernameFormat = NewError(CodeInvalidUsernameFormat, "用户名格式错误")
	ErrInvalidUserID         = NewError(CodeInvalidUserID, "无效的用户ID")
	ErrEmptyTranslation      = NewError(CodeEmptyTranslation, "翻译不能为空")
	ErrEmptyExample          = NewError(CodeEmptyExample, "例句不能为空")
	ErrEmptyTag              = NewError(CodeEmptyTag, "标签不能为空")
	ErrUnauthenticated       = NewError(CodeUnauthenticated, "用户未认证或认证无效")
	ErrInvalidPassword       = NewError(CodeInvalidPassword, "密码错误")
	ErrEmptyPassword         = NewError(CodeEmptyPassword, "密码为空")
)

var (
	ErrInvalidCredentials        = NewError(CodeInvalidCredentials, "用户名或密码错误")
	ErrProgressNotFound          = NewError(CodeProgressNotFound, "学习进度不存在")
	ErrUnauthorized              = NewError(CodeUnauthorized, "未授权")
	ErrRoleNotFound              = NewError(CodeRoleNotFound, "角色不存在")
	ErrLoginUserIdIsMissingInCtx = NewError(CodeLoginUserIdIsMissingInCtx, "上下文缺失登陆用户 ID")
)

// HanChar related errors
var (
	ErrHanCharAlreadyExists   = NewError(CodeHanCharAlreadyExists, "汉字已存在")
	ErrNotImplemented         = NewError(CodeNotImplemented, "功能未实现")
	ErrInvalidDifficultyLevel = NewError(CodeInvalidDifficultyLevel, "无效的难度等级")
)

// admin system
var (
	ErrSystemAlreadyInitialized = NewError(CodeSystemAlreadyInitialized, "系统已初始化")
	ErrSystemNotInitialized     = NewError(CodeSystemNotInitialized, "系统未初始化")
	ErrAdminNotFound            = NewError(CodeAdminNotFound, "管理员不存在")
)

// Memory unit related errors
var (
	ErrMemoryUnitNotFound = NewError(CodeMemoryUnitNotFound, "记忆单元不存在")
	ErrInvalidUnitType    = NewError(CodeInvalidUnitType, "无效的记忆单元类型")
)

// Review related errors
var (
	ErrReviewNotFound = NewError(CodeReviewNotFound, "复习记录不存在")
	ErrInvalidReview  = NewError(CodeInvalidReview, "无效的复习记录")
)

// Content related errors
var (
	ErrContentNotFound = NewError(CodeContentNotFound, "内容不存在")
	ErrInvalidContent  = NewError(CodeInvalidContent, "无效的内容")
)

// Stats related errors
var (
	ErrStatsNotFound = NewError(CodeStatsNotFound, "统计记录不存在")
	ErrInvalidStats  = NewError(CodeInvalidStats, "无效的统计记录")
)

var (
	ErrRefreshTokenMismatch = NewError(CodeRefreshTokenMismatch, "刷新 Token 与登陆态不匹配")
)

// Error 自定义错误类型
type Error struct {
	Code    int    // 错误码
	Message string // 错误信息
}

// Error 实现 error 接口
func (e *Error) Error() string {
	return e.Message
}

// New 创建新的错误
func NewError(code int, message string) error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

func Is(err error, target error) bool {
	return errors.Is(err, target)
}

func As(err error, target any) bool {
	return errors.As(err, target)
}

// IsNotFound 判断是否为未找到错误
func IsNotFound(err error) bool {
	return err == ErrUserNotFound ||
		err == ErrMemoryUnitNotFound ||
		err == ErrReviewNotFound ||
		err == ErrContentNotFound ||
		err == ErrStatsNotFound
}

// IsInvalid 判断是否为无效错误
func IsInvalid(err error) bool {
	return err == ErrInvalidPassword ||
		err == ErrInvalidUnitType ||
		err == ErrInvalidReview ||
		err == ErrInvalidContent ||
		err == ErrInvalidStats
}
