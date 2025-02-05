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

	// 认证相关路由
	public := api.Group("/auth")
	{
		public.POST("/login", handlers.AuthHandler.Login)
		public.POST("/register", handlers.AuthHandler.Register)
	}

	// 需要认证的路由
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware())
	{
		// 单词相关
		protected.GET("/words", handlers.WordHandler.ListWords)
		protected.POST("/words", handlers.WordHandler.CreateWord)
		protected.GET("/words/:wordId", handlers.WordHandler.GetWord)
		protected.DELETE("/words/:wordId", handlers.WordHandler.DeleteWord)

		// 学习进度相关
		protected.GET("/learning/courses/progress", handlers.LearningHandler.ListCourseProgress)
		protected.POST("/learning/courses/:courseId/progress", handlers.LearningHandler.SaveCourseProgress)
		protected.GET("/learning/courses/:courseId/progress", handlers.LearningHandler.GetCourseProgress)
		protected.GET("/learning/courses/:courseId/sections/progress", handlers.LearningHandler.ListSectionProgress)
		protected.POST("/learning/sections/:sectionId/progress", handlers.LearningHandler.SaveSectionProgress)
		protected.GET("/learning/sections/:sectionId/progress", handlers.LearningHandler.GetSectionProgress)
		protected.POST("/learning/units/:unitId/progress", handlers.LearningHandler.SaveUnitProgress)
		protected.GET("/learning/units/:unitId/progress", handlers.LearningHandler.GetUnitProgress)
		protected.GET("/learning/sections/:sectionId/units/progress", handlers.LearningHandler.ListUnitProgress)
	}

	// 健康检查路由
	r.GET("/healthz", handlers.HealthHandler.HealthCheck)
}
