package errors

import "errors"

// Common errors
var (
	ErrInvalidInput = NewError(CodeInvalidInput, "输入错误")
	ErrNotFound     = NewError(CodeNotFound, "资源不存在")
)

// Word related errors
var (
	ErrEmptyWordText       = NewError(CodeEmptyWordText, "单词不能为空")
	ErrWordNotFound        = NewError(CodeWordNotFound, "单词不存在")
	ErrWordAlreadyExists   = NewError(CodeWordAlreadyExists, "单词已存在")
	ErrInvalidDifficulty   = NewError(CodeInvalidDifficulty, "难度必须在1到5之间")
	ErrInvalidMasteryLevel = NewError(CodeInvalidMasteryLevel, "熟练度必须在0到5之间")
	ErrDuplicateTag        = NewError(CodeDuplicateTag, "标签已存在")
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
)

var (
	ErrInvalidCredentials = NewError(CodeInvalidCredentials, "用户名或密码错误")
	ErrProgressNotFound   = NewError(CodeProgressNotFound, "学习进度不存在")
)

type Error struct {
	Code    int
	Message string
}

func (e *Error) Error() string {
	return e.Message
}

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
