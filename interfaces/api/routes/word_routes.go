package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lazyjean/sla2/interfaces/api/handler"
)

func RegisterWordRoutes(r *gin.Engine, wordHandler *handler.WordHandler) {
	api := r.Group("/api")
	{
		words := api.Group("/words")
		{
			words.POST("", wordHandler.CreateWord)
			words.GET("/:id", wordHandler.GetWord)
			words.GET("", wordHandler.ListWords)
			words.PUT("/:id", wordHandler.UpdateWord)
			words.DELETE("/:id", wordHandler.DeleteWord)
			words.GET("/search", wordHandler.SearchWords)
		}
	}
}
