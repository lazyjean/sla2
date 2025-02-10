package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lazyjean/sla2/internal/domain/entity"
)

type WordHandler struct {
}

func NewWordHandler() *WordHandler {
	return &WordHandler{}
}

func (h *WordHandler) CreateWord(c *gin.Context) {
	var word entity.Word
	if err := c.ShouldBindJSON(&word); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: 实现创建单词的逻辑

	c.JSON(http.StatusOK, word)
}

func (h *WordHandler) GetWords(c *gin.Context) {
	// TODO: 实现获取单词列表的逻辑
	c.JSON(http.StatusOK, gin.H{"message": "获取单词列表"})
}
