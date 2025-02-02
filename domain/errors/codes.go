package errors

const (
	// 通用错误码 (1000-1999)
	CodeInvalidInput    = 1000
	CodeNotFound        = 1001
	CodeInvalidArgument = 1002

	// 单词相关错误码 (2000-2999)
	CodeEmptyWordText     = 2000
	CodeEmptyTranslation  = 2001
	CodeWordNotFound      = 2002
	CodeWordAlreadyExists = 2003

	// 系统错误码 (9000-9999)
	CodeInternalError = 9000

	// 用户相关错误码 (3000-3999)
	CodeUserNotFound          = 3000
	CodeUserAlreadyExists     = 3001
	CodeInvalidCredentials    = 3002
	CodeInvalidPhoneFormat    = 3003
	CodeInvalidEmailFormat    = 3004
	CodeInvalidUsernameFormat = 3005
)
