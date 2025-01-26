package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lazyjean/sla2/interfaces/api/handler"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func RegisterWordRoutes(r *gin.Engine, wordHandler *handler.WordHandler) {
	// Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api")
	{
		words := api.Group("/words")
		{
			words.POST("", wordHandler.CreateWord)
			words.GET("/:id", wordHandler.GetWord)
			words.GET("", wordHandler.ListWords)
			words.PUT("/:id", wordHandler.UpdateWord)
			words.DELETE("/:id", wordHandler.DeleteWord)
		}
	}
}
