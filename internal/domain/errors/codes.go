package errors

const (
	// 通用错误码 (1000-1999)
	CodeInvalidInput = 1000 + iota
	CodeNotFound
	CodeInvalidArgument
	CodeInvalidWord

	// 单词相关错误码 (2000-2999)
	CodeEmptyWordText = 2000 + iota
	CodeWordNotFound
	CodeWordAlreadyExists
	CodeInvalidDifficulty
	CodeInvalidMasteryLevel
	CodeDuplicateTag
	CodeEmptyDefinition

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
	CodeUnauthenticated
	CodeInvalidPassword
	CodeUnauthorized
	CodeRoleNotFound
	CodeEmptyPassword

	// 认证相关错误码 (5000-5999)
	CodeInvalidCredentials = 5000 + iota
	CodeProgressNotFound
	CodeRefreshTokenMismatch
	CodeLoginUserIdIsMissingInCtx

	// 其他错误码 (6000-6999)
	CodeUnknown = 6000 + iota
	CodeAlreadyExists
	CodePermissionDenied

	// HanChar related error codes (7000-7999)
	CodeHanCharAlreadyExists = 7000 + iota
	CodeNotImplemented
	CodeInvalidDifficultyLevel

	// 系统相关错误码 (8000-8999)
	CodeSystemAlreadyInitialized = 8000 + iota
	CodeSystemNotInitialized
	CodeAdminNotFound

	// 记忆单元相关错误码 (9000-9999)
	CodeMemoryUnitNotFound = 9000 + iota
	CodeInvalidUnitType

	// 复习相关错误码 (10000-10999)
	CodeReviewNotFound = 10000 + iota
	CodeInvalidReview

	// 内容相关错误码 (11000-11999)
	CodeContentNotFound = 11000 + iota
	CodeInvalidContent

	// 统计相关错误码 (12000-12999)
	CodeStatsNotFound = 12000 + iota
	CodeInvalidStats

	// 系统错误码 (13000-13999)
	CodeInternalError = 13000
)
