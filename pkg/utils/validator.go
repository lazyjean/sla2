package utils

import "regexp"

// IsValidEmail 验证邮箱格式
func IsValidEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}

// IsValidPhone 验证手机号格式（中国大陆手机号）
func IsValidPhone(phone string) bool {
	pattern := `^1[3-9]\d{9}$`
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(phone)
}

// IsValidUsername 验证用户名格式
func IsValidUsername(username string) bool {
	pattern := `^[a-zA-Z0-9_-]{3,20}$`
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(username)
}
