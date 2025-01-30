package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lazyjean/sla2/application/dto"
	"github.com/lazyjean/sla2/application/service"
	"github.com/lazyjean/sla2/domain/errors"
	"github.com/lazyjean/sla2/domain/repository"
	"github.com/lazyjean/sla2/pkg/logger"
	"go.uber.org/zap"
)

// @title           单词本 API
// @version         1.0
// @description     单词本服务 API 文档
// @termsOfService  http://swagger.io/terms/

// @host      localhost:9000
// @BasePath  /api
// @schemes   http

// @servers   url=http://localhost:9000/api   description=本地开发环境
// @servers   url=http://api.example.com/api  description=生产环境

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description 请在此输入 Bearer token

type WordHandler struct {
	wordService *service.WordService
}

func NewWordHandler(wordService *service.WordService) *WordHandler {
	return &WordHandler{
		wordService: wordService,
	}
}

// CreateWord 创建生词
// @Summary      创建新单词
// @Description  创建一个新的单词记录
// @Tags         words
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        Authorization  header    string           true   "Bearer 令牌"
// @Param        word          body      dto.WordCreateDTO  true  "单词信息"
// @Success      201   {object}  Response{data=dto.WordResponseDTO}
// @Failure      400   {object}  Response          "请求参数错误"
// @Failure      401   {object}  Response          "未授权"
// @Failure      500   {object}  Response          "服务器内部错误"
// @Router       /v1/words [post]
func (h *WordHandler) CreateWord(c *gin.Context) {
	var createDTO dto.WordCreateDTO
	if err := c.ShouldBindJSON(&createDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 从上下文获取用户ID
	userID := getUserIDFromContext(c)

	word, err := h.wordService.CreateWord(c.Request.Context(), &createDTO, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dto.WordResponseDTOFromEntity(word))
}

// GetWord 获取单词
// @Summary      获取单词详情
// @Description  通过ID获取单词的详细信息
// @Tags         words
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        Authorization  header    string  true   "Bearer 令牌"
// @Param        id   path      uint    true  "单词ID"
// @Success      200  {object}  Response{data=dto.WordResponseDTO}
// @Failure      400  {object}  Response          "请求参数错误"
// @Failure      401  {object}  Response          "未授权"
// @Failure      404  {object}  Response          "单词不存在"
// @Failure      500  {object}  Response          "服务器内部错误"
// @Router       /v1/words/{id} [get]
func (h *WordHandler) GetWord(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(400, "invalid id format"))
		return
	}

	resp, err := h.wordService.GetWord(c.Request.Context(), uint(id))
	if err != nil {
		if err == errors.ErrWordNotFound {
			c.JSON(http.StatusNotFound, NewErrorResponse(404, err.Error()))
			return
		}
		logger.Log.Error("Failed to get word", zap.Error(err), zap.Uint64("id", id))
		c.JSON(http.StatusInternalServerError, NewErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewResponse(0, "success", resp))
}

// ListWords 获取生词列表
// @Summary      获取单词列表
// @Description  分页获取单词列表
// @Tags         words
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        Authorization  header    string  true   "Bearer 令牌"
// @Param        page     query     int  false  "页码"  default(1)
// @Param        perPage  query     int  false  "每页数量"  default(10)
// @Success      200      {object}  Response{data=ListResponse{items=[]dto.WordResponseDTO}}
// @Failure      401      {object}  Response          "未授权"
// @Failure      500      {object}  Response          "服务器内部错误"
// @Router       /v1/words [get]
func (h *WordHandler) ListWords(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	userID := getUserIDFromContext(c)
	offset := (page - 1) * pageSize

	words, total, err := h.wordService.ListWords(c.Request.Context(), userID, offset, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 转换为响应 DTO
	dtos := make([]*dto.WordResponseDTO, len(words))
	for i, word := range words {
		dtos[i] = dto.WordResponseDTOFromEntity(word)
	}

	c.JSON(http.StatusOK, gin.H{
		"total": total,
		"items": dtos,
	})
}

// DeleteWord 删除单词
// @Summary      删除单词
// @Description  删除指定ID的单词
// @Tags         words
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        Authorization  header    string  true   "Bearer 令牌"
// @Param        id   path      uint    true  "单词ID"
// @Success      200  {object}  Response
// @Failure      400  {object}  Response          "请求参数错误"
// @Failure      401  {object}  Response          "未授权"
// @Failure      404  {object}  Response          "单词不存在"
// @Failure      500  {object}  Response          "服务器内部错误"
// @Router       /v1/words/{id} [delete]
func (h *WordHandler) DeleteWord(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(400, "invalid id format"))
		return
	}

	err = h.wordService.DeleteWord(c.Request.Context(), uint(id))
	if err != nil {
		switch err {
		case errors.ErrWordNotFound:
			// 返回 200 但在响应中包含错误信息
			c.JSON(http.StatusOK, NewErrorResponse(404, err.Error()))
		default:
			logger.Log.Error("Failed to delete word", zap.Error(err), zap.Uint64("id", id))
			c.JSON(http.StatusInternalServerError, NewErrorResponse(500, err.Error()))
		}
		return
	}

	c.JSON(http.StatusOK, NewResponse(0, "success", nil))
}

// SearchWords 搜索单词
// @Summary      搜索单词
// @Description  根据条件搜索单词
// @Tags         words
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        Authorization  header    string    true   "Bearer 令牌"
// @Param        text           query    string    false  "单词文本"
// @Param        tags           query    []string  false  "标签列表"
// @Param        minDifficulty  query    int       false  "最小难度"
// @Param        maxDifficulty  query    int       false  "最大难度"
// @Param        page           query    int       false  "页码"  default(1)
// @Param        pageSize       query    int       false  "每页数量"  default(10)
// @Success      200            {object}  Response{data=ListResponse{items=[]dto.WordResponseDTO}}
// @Failure      400            {object}  Response          "请求参数错误"
// @Failure      401            {object}  Response          "未授权"
// @Failure      500            {object}  Response          "服务器内部错误"
// @Router       /v1/words/search [get]
func (h *WordHandler) SearchWords(c *gin.Context) {
	// 解析查询参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	minDifficulty, _ := strconv.Atoi(c.Query("min_difficulty"))
	maxDifficulty, _ := strconv.Atoi(c.Query("max_difficulty"))

	query := &repository.WordQuery{
		Text:          c.Query("text"),
		Tags:          c.QueryArray("tags"),
		MinDifficulty: minDifficulty,
		MaxDifficulty: maxDifficulty,
		Offset:        (page - 1) * pageSize,
		Limit:         pageSize,
	}

	words, total, err := h.wordService.SearchWords(c.Request.Context(), query)
	if err != nil {
		logger.Log.Error("Failed to search words",
			zap.Error(err),
			zap.Any("query", query),
		)
		c.JSON(http.StatusInternalServerError, NewErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewListResponse(words, page, pageSize, total))
}
