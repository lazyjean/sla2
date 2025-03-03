package security

import (
	"golang.org/x/crypto/bcrypt"

	domainSecurity "github.com/lazyjean/sla2/internal/domain/security"
)

type Cost int

const (
	MinCost     Cost = Cost(bcrypt.MinCost)
	MaxCost     Cost = Cost(bcrypt.MaxCost)
	DefaultCost Cost = Cost(bcrypt.DefaultCost)
)

// BCryptPasswordService 使用 bcrypt 实现的密码服务
type BCryptPasswordService struct {
	cost Cost
}

// NewBCryptPasswordService 创建一个新的 BCrypt 密码服务
func NewBCryptPasswordService() domainSecurity.PasswordService {
	cost := DefaultCost
	return &BCryptPasswordService{
		cost: cost,
	}
}

// HashPassword 使用 bcrypt 对密码进行哈希处理
func (s *BCryptPasswordService) HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), int(s.cost))
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// VerifyPassword 验证密码是否匹配
func (s *BCryptPasswordService) VerifyPassword(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

var _ domainSecurity.PasswordService = (*BCryptPasswordService)(nil)
