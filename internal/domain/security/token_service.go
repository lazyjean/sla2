package security

// TokenService 令牌服务接口
/*
TokenService 目前处理用户认证和基本的角色权限管理。

未来会员系统设计建议：
1. 会员信息应独立于角色系统，在 Claims 中单独维护，例如：
   type MemberInfo struct {
       Level      string    // 会员等级
       ExpireAt   time.Time // 到期时间
       Privileges []string  // 特权列表
   }

2. 接口方法扩展建议：
   GenerateToken(userID string, roles []string, memberInfo *MemberInfo) (string, error)
   ValidateToken(token string) (userID string, roles []string, memberInfo *MemberInfo, err error)

3. Claims 结构建议：
   claims := jwt.MapClaims{
       "sub": userID,
       "roles": roles,
       "member": map[string]interface{}{
           "level": memberInfo.Level,
           "expireAt": memberInfo.ExpireAt.Unix(),
           "privileges": memberInfo.Privileges,
       },
   }

这样的设计可以：
- 保持角色和会员系统的职责分离
- 支持会员等级、到期时间、特权等信息的独立管理
- 避免会员状态变化导致频繁更新角色信息
*/
type TokenService interface {
	// GenerateToken 生成访问令牌
	GenerateToken(userID string, roles []string) (string, error)
	// ValidateToken 验证访问令牌
	ValidateToken(token string) (string, []string, error)
	// GenerateRefreshToken 生成刷新令牌
	GenerateRefreshToken(userID string, roles []string) (string, error)
	// ValidateRefreshToken 验证刷新令牌
	ValidateRefreshToken(token string) (string, []string, error)
}
