package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lazyjean/sla2/internal/application/dto"
	"github.com/lazyjean/sla2/internal/application/service"
	"github.com/lazyjean/sla2/internal/domain/errors"
	"github.com/lazyjean/sla2/pkg/jwt"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Register 用户注册
// @Summary      用户注册
// @Description  创建新用户账号
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        user  body      dto.RegisterDTO  true  "用户注册信息"
// @Success      201   {object}  Response{data=dto.UserDTO}
// @Failure      400   {object}  Response  "请求参数错误"
// @Failure      500   {object}  Response  "服务器内部错误"
// @Router       /register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(errors.CodeInvalidInput, "无效的请求参数"))
		return
	}

	user, err := h.authService.Register(c.Request.Context(), &service.RegisterRequest{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewErrorResponse(errors.CodeInternalError, err.Error()))
		return
	}

	c.JSON(http.StatusCreated, NewResponse(0, "注册成功", user))
}

// Login 用户登录
// @Summary      用户登录
// @Description  使用账号密码登录
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        credentials  body      dto.LoginDTO  true  "登录凭证"
// @Success      200         {object}  Response{data=dto.TokenDTO}
// @Failure      400         {object}  Response  "请求参数错误"
// @Failure      401         {object}  Response  "认证失败"
// @Failure      500         {object}  Response  "服务器内部错误"
// @Router       /login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(errors.CodeInvalidInput, "无效的请求参数"))
		return
	}

	userDTO, err := h.authService.Login(c.Request.Context(), &service.LoginRequest{
		Account:  req.Account,
		Password: req.Password,
	})
	if err != nil {
		c.JSON(http.StatusUnauthorized, NewErrorResponse(errors.CodeInvalidCredentials, err.Error()))
		return
	}

	// 生成JWT
	token, err := jwt.GenerateToken(int(userDTO.UserID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	// 设置 JWT token 到 cookie
	c.SetCookie(
		"jwt",     // cookie 名称
		token,     // cookie 值
		3600*24*7, // 过期时间（秒）：24小时
		"/",       // cookie 路径
		"",        // domain（留空表示当前域名）
		true,      // 仅限 HTTPS
		true,      // HTTP-only
	)

	c.JSON(http.StatusOK, NewResponse(0, "登录成功", token))
}
