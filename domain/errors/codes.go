package errors

const (
	// 通用错误码 (1000-1999)
	CodeInvalidInput = 1000
	CodeNotFound     = 1001

	// 单词相关错误码 (2000-2999)
	CodeEmptyWordText     = 2000
	CodeEmptyTranslation  = 2001
	CodeWordNotFound      = 2002
	CodeWordAlreadyExists = 2003

	// 系统错误码 (9000-9999)
	CodeInternalError = 9000
)
