package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lazyjean/sla2/config"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// HealthCheck 健康检查接口
// @Summary      服务健康状态
// @Description  获取服务运行状态
// @Tags         system
// @Accept       json
// @Produce      json
// @Success      200  {object}  Response{data=string}  "服务状态"
// @Router       /healthz [get]
func (h *HealthHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, NewResponse(0, "服务运行正常", gin.H{
		"status":  "ok",
		"version": config.GetConfig().Server.Version,
	}))
}
