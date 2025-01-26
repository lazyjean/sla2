package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lazyjean/sla2/application/dto"
	"github.com/lazyjean/sla2/application/service"
	"github.com/lazyjean/sla2/domain/errors"
	"github.com/lazyjean/sla2/pkg/logger"
	"go.uber.org/zap"
)

type WordHandler struct {
	wordService *service.WordService
}

func NewWordHandler(wordService *service.WordService) *WordHandler {
	return &WordHandler{
		wordService: wordService,
	}
}

// CreateWord 创建单词
// @Summary      创建新单词
// @Description  创建一个新的单词记录
// @Tags         words
// @Accept       json
// @Produce      json
// @Param        word  body      dto.WordCreateDTO  true  "单词信息"
// @Success      201   {object}  Response{data=dto.WordResponseDTO}
// @Failure      400   {object}  Response
// @Failure      500   {object}  Response
// @Router       /words [post]
func (h *WordHandler) CreateWord(c *gin.Context) {
	var createDTO dto.WordCreateDTO
	if err := c.ShouldBindJSON(&createDTO); err != nil {
		logger.Log.Error("Failed to bind JSON", zap.Error(err))
		c.JSON(http.StatusBadRequest, NewErrorResponse(400, errors.ErrInvalidInput.Error()))
		return
	}

	resp, err := h.wordService.CreateWord(c.Request.Context(), &createDTO)
	if err != nil {
		logger.Log.Error("Failed to create word", zap.Error(err))
		c.JSON(http.StatusInternalServerError, NewErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusCreated, NewResponse(0, "success", resp))
}

// GetWord 获取单词
// @Summary      获取单词详情
// @Description  通过ID获取单词的详细信息
// @Tags         words
// @Accept       json
// @Produce      json
// @Param        id   path      uint    true  "单词ID"
// @Success      200  {object}  Response{data=dto.WordResponseDTO}
// @Failure      400  {object}  Response
// @Failure      404  {object}  Response
// @Failure      500  {object}  Response
// @Router       /words/{id} [get]
func (h *WordHandler) GetWord(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, NewErrorResponse(400, "id is required"))
		return
	}

	resp, err := h.wordService.GetWord(c.Request.Context(), id)
	if err != nil {
		if err == errors.ErrWordNotFound {
			c.JSON(http.StatusNotFound, NewErrorResponse(404, err.Error()))
			return
		}
		logger.Log.Error("Failed to get word", zap.Error(err), zap.String("id", id))
		c.JSON(http.StatusInternalServerError, NewErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewResponse(0, "success", resp))
}

// ListWords 获取单词列表
// @Summary      获取单词列表
// @Description  分页获取单词列表
// @Tags         words
// @Accept       json
// @Produce      json
// @Param        page     query     int  false  "页码"  default(1)
// @Param        perPage  query     int  false  "每页数量"  default(10)
// @Success      200      {object}  Response{data=ListResponse{items=[]dto.WordResponseDTO}}
// @Failure      500      {object}  Response
// @Router       /words [get]
func (h *WordHandler) ListWords(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	words, total, err := h.wordService.ListWords(c.Request.Context(), page, pageSize)
	if err != nil {
		logger.Log.Error("Failed to list words",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("page_size", pageSize),
		)
		c.JSON(http.StatusInternalServerError, NewErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewListResponse(words, page, pageSize, total))
}

// UpdateWord 更新单词
// @Summary      更新单词信息
// @Description  更新指定ID的单词信息
// @Tags         words
// @Accept       json
// @Produce      json
// @Param        id    path      uint             true  "单词ID"
// @Param        word  body      dto.WordCreateDTO  true  "单词信息"
// @Success      200   {object}  Response{data=dto.WordResponseDTO}
// @Failure      400   {object}  Response
// @Failure      404   {object}  Response
// @Failure      500   {object}  Response
// @Router       /words/{id} [put]
func (h *WordHandler) UpdateWord(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(400, "invalid id format"))
		return
	}

	var updateDTO dto.WordCreateDTO
	if err := c.ShouldBindJSON(&updateDTO); err != nil {
		logger.Log.Error("Failed to bind JSON", zap.Error(err))
		c.JSON(http.StatusBadRequest, NewErrorResponse(400, errors.ErrInvalidInput.Error()))
		return
	}

	resp, err := h.wordService.UpdateWord(c.Request.Context(), uint(id), &updateDTO)
	if err != nil {
		if err == errors.ErrWordNotFound {
			c.JSON(http.StatusNotFound, NewErrorResponse(404, err.Error()))
			return
		}
		logger.Log.Error("Failed to update word", zap.Error(err), zap.Uint64("id", id))
		c.JSON(http.StatusInternalServerError, NewErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewResponse(0, "success", resp))
}

// DeleteWord 删除单词
// @Summary      删除单词
// @Description  删除指定ID的单词
// @Tags         words
// @Accept       json
// @Produce      json
// @Param        id   path      uint    true  "单词ID"
// @Success      204  {object}  nil
// @Failure      400  {object}  Response
// @Failure      404  {object}  Response
// @Failure      500  {object}  Response
// @Router       /words/{id} [delete]
func (h *WordHandler) DeleteWord(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(400, "invalid id format"))
		return
	}

	if err := h.wordService.DeleteWord(c.Request.Context(), uint(id)); err != nil {
		if err == errors.ErrWordNotFound {
			c.JSON(http.StatusNotFound, NewErrorResponse(404, err.Error()))
			return
		}
		logger.Log.Error("Failed to delete word", zap.Error(err), zap.Uint64("id", id))
		c.JSON(http.StatusInternalServerError, NewErrorResponse(500, err.Error()))
		return
	}

	c.Status(http.StatusNoContent)
}
