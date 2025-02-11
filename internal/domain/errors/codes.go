package errors

const (
	// 通用错误码 (1000-1999)
	CodeInvalidInput = 1000 + iota
	CodeNotFound
	CodeInvalidArgument

	// 单词相关错误码 (2000-2999)
	CodeEmptyWordText = 2000 + iota
	CodeWordNotFound
	CodeWordAlreadyExists
	CodeInvalidDifficulty
	CodeInvalidMasteryLevel
	CodeDuplicateTag

	// repository 错误码 (3000-3999)
	CodeFailedToSave = 3000 + iota
	CodeFailedToUpdate
	CodeFailedToDelete
	CodeFailedToQuery

	// 用户相关错误码 (4000-4999)
	CodeUserNotFound = 4000 + iota
	CodeUserAlreadyExists
	CodeInvalidPhoneFormat
	CodeInvalidEmailFormat
	CodeInvalidUsernameFormat
	CodeInvalidUserID
	CodeEmptyTranslation
	CodeEmptyExample
	CodeEmptyTag

	// 认证相关错误码 (5000-5999)
	CodeInvalidCredentials = 5000 + iota
	CodeProgressNotFound

	// 其他错误码 (6000-6999)
	CodeUnknown = 6000 + iota
	CodeAlreadyExists
	CodePermissionDenied
	CodeUnauthenticated

	// 系统错误码 (9000-9999)
	CodeInternalError = 9000
)
