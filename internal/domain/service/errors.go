package service

import "errors"

// 通用错误定义
var (
	ErrUnauthorized = errors.New("unauthorized")
)
