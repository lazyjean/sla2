package security

// PasswordService 密码服务接口
type PasswordService interface {
	// HashPassword 对密码进行哈希处理
	HashPassword(password string) (string, error)

	// VerifyPassword 验证密码是否匹配
	VerifyPassword(password, hashedPassword string) bool
}