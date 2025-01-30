package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lazyjean/sla2/interfaces/api/handler"
	"github.com/lazyjean/sla2/interfaces/api/middleware"
)

// @title SLA2 API
// @version 1.0
// @description SLA2 单词学习助手 API 服务
// @BasePath /api/v1
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
func SetupRoutes(r *gin.Engine, handlers *handler.Handlers) {
	// API 路由组
	api := r.Group("/api/v1")

	// 公开路由
	public := api.Group("")
	{
		public.POST("/login", handlers.AuthHandler.Login)
		public.POST("/register", handlers.AuthHandler.Register)
	}

	// 需要认证的路由
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware())
	{
		// 单词相关路由
		protected.POST("/words", handlers.WordHandler.CreateWord)
		protected.GET("/words", handlers.WordHandler.ListWords)
		protected.GET("/words/:id", handlers.WordHandler.GetWord)
		protected.DELETE("/words/:id", handlers.WordHandler.DeleteWord)
	}
}
