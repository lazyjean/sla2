package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lazyjean/sla2/internal/application/dto"
	"github.com/lazyjean/sla2/internal/application/service"
	"github.com/lazyjean/sla2/pkg/logger"
	"go.uber.org/zap"
)

type LearningHandler struct {
	learningService *service.LearningService
}

func NewLearningHandler(learningService *service.LearningService) *LearningHandler {
	return &LearningHandler{
		learningService: learningService,
	}
}

// SaveCourseProgress 保存课程进度
// @Summary      保存课程学习进度
// @Tags         learning
// @Security     Bearer
// @Param        courseId   path    int     true  "课程ID"
// @Param        status     query   string  true  "学习状态" Enums(not_started, in_progress, completed)
// @Param        score      query   int     false "得分"
// @Success      200  {object}  Response{data=dto.CourseProgressDTO}
// @Router       /learning/courses/{courseId}/progress [post]
func (h *LearningHandler) SaveCourseProgress(c *gin.Context) {
	userID := getUserIDFromContext(c)
	courseID, err := strconv.ParseUint(c.Param("courseId"), 10, 32)
	if err != nil || courseID == 0 {
		c.JSON(http.StatusBadRequest, NewErrorResponse(400, "无效的课程ID"))
		return
	}

	status := c.Query("status")
	if status == "" {
		c.JSON(http.StatusBadRequest, NewErrorResponse(400, "Status is required"))
		return
	}

	score, _ := strconv.Atoi(c.DefaultQuery("score", "0"))

	progress, err := h.learningService.SaveCourseProgress(c.Request.Context(), userID, uint(courseID), status, score)
	if err != nil {
		logger.Log.Error("Failed to save course progress",
			zap.Error(err),
			zap.Uint("userID", userID),
			zap.Uint64("courseID", courseID),
		)
		c.JSON(http.StatusInternalServerError, NewErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewResponse(0, "success", dto.CourseProgressToDTO(progress)))
}

// GetCourseProgress 获取课程进度
// @Summary      获取课程学习进度
// @Tags         learning
// @Security     Bearer
// @Param        courseId   path    int     true  "课程ID"
// @Success      200  {object}  Response{data=dto.CourseProgressDTO}
// @Router       /learning/courses/{courseId}/progress [get]
func (h *LearningHandler) GetCourseProgress(c *gin.Context) {
	userID := getUserIDFromContext(c)
	courseID, err := strconv.ParseUint(c.Param("courseId"), 10, 32)
	if err != nil || courseID == 0 {
		c.JSON(http.StatusBadRequest, NewErrorResponse(400, "无效的课程ID"))
		return
	}

	progress, err := h.learningService.GetCourseProgress(c.Request.Context(), userID, uint(courseID))
	if err != nil {
		logger.Log.Error("Failed to get course progress",
			zap.Error(err),
			zap.Uint("userID", userID),
			zap.Uint64("courseID", courseID),
		)
		c.JSON(http.StatusInternalServerError, NewErrorResponse(500, err.Error()))
		return
	}

	if progress == nil {
		c.JSON(http.StatusNotFound, NewErrorResponse(404, "Course progress not found"))
		return
	}

	c.JSON(http.StatusOK, NewResponse(0, "success", dto.CourseProgressToDTO(progress)))
}

// ListCourseProgress 获取课程进度列表
// @Summary      获取课程学习进度列表
// @Tags         learning
// @Security     Bearer
// @Param        page      query    int     false  "页码"  default(1)
// @Param        pageSize  query    int     false  "每页数量"  default(10)
// @Param        courseId  path    int     true  "课程ID"
// @Success      200  {object}  Response{data=[]dto.CourseProgressDTO}
// @Router       /learning/courses/{courseId}/progress [get]
func (h *LearningHandler) ListCourseProgress(c *gin.Context) {
	userID := getUserIDFromContext(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	progresses, total, err := h.learningService.ListCourseProgress(c.Request.Context(), userID, page, pageSize)
	if err != nil {
		logger.Log.Error("Failed to list course progress", zap.Error(err), zap.Uint("userID", userID))
		c.JSON(http.StatusInternalServerError, NewErrorResponse(500, err.Error()))
		return
	}

	items := make([]*dto.CourseProgressDTO, len(progresses))
	for i, p := range progresses {
		items[i] = dto.CourseProgressToDTO(p)
	}

	c.JSON(http.StatusOK, NewListResponse(items, page, pageSize, total))
}

// SaveSectionProgress 保存章节进度
// @Summary      保存章节学习进度
// @Tags         learning
// @Security     Bearer
// @Param        sectionId  path    int     true  "章节ID"
// @Param        courseId   query   int     true  "课程ID"
// @Param        status     query   string  true  "学习状态" Enums(not_started, in_progress, completed)
// @Param        progress   query   number  true  "进度" minimum(0) maximum(100)
// @Success      200  {object}  Response{data=dto.SectionProgressDTO}
// @Router       /learning/sections/{sectionId}/progress [post]
func (h *LearningHandler) SaveSectionProgress(c *gin.Context) {
	userID := getUserIDFromContext(c)
	sectionID, err := strconv.ParseUint(c.Param("sectionId"), 10, 32)
	if err != nil || sectionID == 0 {
		c.JSON(http.StatusBadRequest, NewErrorResponse(400, "无效的章节ID"))
		return
	}
	courseID, err := strconv.ParseUint(c.Query("courseId"), 10, 32)
	if err != nil || courseID == 0 {
		c.JSON(http.StatusBadRequest, NewErrorResponse(400, "无效的课程ID"))
		return
	}
	status := c.Query("status")
	progress, _ := strconv.ParseFloat(c.Query("progress"), 64)

	sectionProgress, err := h.learningService.SaveSectionProgress(c.Request.Context(), userID, uint(courseID), uint(sectionID), status, progress)
	if err != nil {
		logger.Log.Error("Failed to save section progress", zap.Error(err))
		c.JSON(http.StatusInternalServerError, NewErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewResponse(0, "success", dto.SectionProgressToDTO(sectionProgress)))
}

// GetSectionProgress 获取章节进度
// @Summary      获取章节学习进度
// @Tags         learning
// @Security     Bearer
// @Param        sectionId  path    int     true  "章节ID"
// @Success      200  {object}  Response{data=dto.SectionProgressDTO}
// @Router       /learning/sections/{sectionId}/progress [get]
func (h *LearningHandler) GetSectionProgress(c *gin.Context) {
	userID := getUserIDFromContext(c)
	sectionID, err := strconv.ParseUint(c.Param("sectionId"), 10, 32)
	if err != nil || sectionID == 0 {
		c.JSON(http.StatusBadRequest, NewErrorResponse(400, "无效的章节ID"))
		return
	}
	progress, err := h.learningService.GetSectionProgress(c.Request.Context(), userID, uint(sectionID))
	if err != nil {
		logger.Log.Error("Failed to get section progress", zap.Error(err))
		c.JSON(http.StatusInternalServerError, NewErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewResponse(0, "success", dto.SectionProgressToDTO(progress)))
}

// ListSectionProgress 获取章节进度列表
// @Summary      获取课程的章节学习进度列表
// @Tags         learning
// @Security     Bearer
// @Param        courseId   path    int     true  "课程ID"
// @Success      200  {object}  Response{data=[]dto.SectionProgressDTO}
// @Router       /learning/courses/{courseId}/sections/progress [get]
func (h *LearningHandler) ListSectionProgress(c *gin.Context) {
	userID := getUserIDFromContext(c)
	courseID, err := strconv.ParseUint(c.Param("courseId"), 10, 32)
	if err != nil || courseID == 0 {
		c.JSON(http.StatusBadRequest, NewErrorResponse(400, "无效的课程ID"))
		return
	}

	progresses, err := h.learningService.ListSectionProgress(c.Request.Context(), userID, uint(courseID))
	if err != nil {
		logger.Log.Error("Failed to list section progress", zap.Error(err))
		c.JSON(http.StatusInternalServerError, NewErrorResponse(500, err.Error()))
		return
	}

	items := make([]*dto.SectionProgressDTO, len(progresses))
	for i, p := range progresses {
		items[i] = dto.SectionProgressToDTO(p)
	}

	c.JSON(http.StatusOK, NewResponse(0, "success", items))
}

// SaveUnitProgress 保存单元进度
// @Summary      保存单元学习进度
// @Tags         learning
// @Security     Bearer
// @Param        unitId     path    int     true  "单元ID"
// @Param        sectionId  query   int     true  "章节ID"
// @Param        status     query   string  true  "学习状态" Enums(not_started, in_progress, completed)
// @Param        progress   query   number  true  "进度" minimum(0) maximum(100)
// @Success      200  {object}  Response{data=dto.UnitProgressDTO}
// @Router       /learning/units/{unitId}/progress [post]
func (h *LearningHandler) SaveUnitProgress(c *gin.Context) {
	userID := getUserIDFromContext(c)
	unitID, err := strconv.ParseUint(c.Param("unitId"), 10, 32)
	if err != nil || unitID == 0 {
		c.JSON(http.StatusBadRequest, NewErrorResponse(400, "无效的单元ID"))
		return
	}
	sectionID, err := strconv.ParseUint(c.Query("sectionId"), 10, 32)
	if err != nil || sectionID == 0 {
		c.JSON(http.StatusBadRequest, NewErrorResponse(400, "无效的章节ID"))
		return
	}
	status := c.Query("status")
	progress, _ := strconv.ParseFloat(c.Query("progress"), 64)
	var lastWordID *uint
	if id, err := strconv.ParseUint(c.Query("lastWordId"), 10, 32); err == nil {
		uid := uint(id)
		lastWordID = &uid
	}

	unitProgress, err := h.learningService.SaveUnitProgress(c.Request.Context(), userID, uint(sectionID), uint(unitID), status, progress, lastWordID)
	if err != nil {
		logger.Log.Error("Failed to save unit progress", zap.Error(err))
		c.JSON(http.StatusInternalServerError, NewErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewResponse(0, "success", dto.UnitProgressToDTO(unitProgress)))
}

// GetUnitProgress 获取单元进度
// @Summary      获取单元学习进度
// @Tags         learning
// @Security     Bearer
// @Param        unitId     path    int     true  "单元ID"
// @Success      200  {object}  Response{data=dto.UnitProgressDTO}
// @Router       /learning/units/{unitId}/progress [get]
func (h *LearningHandler) GetUnitProgress(c *gin.Context) {
	userID := getUserIDFromContext(c)
	unitID, err := strconv.ParseUint(c.Param("unitId"), 10, 32)
	if err != nil || unitID == 0 {
		c.JSON(http.StatusBadRequest, NewErrorResponse(400, "无效的单元ID"))
		return
	}

	progress, err := h.learningService.GetUnitProgress(c.Request.Context(), userID, uint(unitID))
	if err != nil {
		logger.Log.Error("Failed to get unit progress", zap.Error(err))
		c.JSON(http.StatusInternalServerError, NewErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewResponse(0, "success", dto.UnitProgressToDTO(progress)))
}

// ListUnitProgress 获取单元进度列表
// @Summary      获取章节的单元学习进度列表
// @Tags         learning
// @Security     Bearer
// @Param        sectionId  path    int     true  "章节ID"
// @Success      200  {object}  Response{data=[]dto.UnitProgressDTO}
// @Router       /learning/sections/{sectionId}/units/progress [get]
func (h *LearningHandler) ListUnitProgress(c *gin.Context) {
	userID := getUserIDFromContext(c)
	sectionID, err := strconv.ParseUint(c.Param("sectionId"), 10, 32)
	if err != nil || sectionID == 0 {
		c.JSON(http.StatusBadRequest, NewErrorResponse(400, "无效的章节ID"))
		return
	}

	progresses, err := h.learningService.ListUnitProgress(c.Request.Context(), userID, uint(sectionID))
	if err != nil {
		logger.Log.Error("Failed to list unit progress", zap.Error(err))
		c.JSON(http.StatusInternalServerError, NewErrorResponse(500, err.Error()))
		return
	}

	items := make([]*dto.UnitProgressDTO, len(progresses))
	for i, p := range progresses {
		items[i] = dto.UnitProgressToDTO(p)
	}

	c.JSON(http.StatusOK, NewResponse(0, "success", items))
}
