package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"crm/models"
	"crm/services"
)

type ActivityHandler struct {
	activityService *services.ActivityService
}

func NewActivityHandler(activityService *services.ActivityService) *ActivityHandler {
	return &ActivityHandler{
		activityService: activityService,
	}
}

// CreateActivity 创建跟进记录
func (h *ActivityHandler) CreateActivity(c *gin.Context) {
	var req models.ActivityCreateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误: " + err.Error(),
		})
		return
	}

	activity, err := h.activityService.CreateActivity(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "创建跟进记录失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "创建跟进记录成功",
		"data":    activity,
	})
}

// GetActivitiesByCustomer 获取客户的跟进记录列表
func (h *ActivityHandler) GetActivitiesByCustomer(c *gin.Context) {
	customerIDStr := c.Param("customer_id")
	customerID, err := strconv.ParseUint(customerIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "客户ID格式错误",
		})
		return
	}

	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}

	activities, total, err := h.activityService.GetActivitiesByCustomer(customerID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取跟进记录失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取跟进记录成功",
		"data": gin.H{
			"list":      activities,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetActivityByID 根据ID获取跟进记录详情
func (h *ActivityHandler) GetActivityByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "记录ID格式错误",
		})
		return
	}

	activity, err := h.activityService.GetActivityByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取跟进记录失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取跟进记录成功",
		"data":    activity,
	})
}

// UpdateActivity 更新跟进记录
func (h *ActivityHandler) UpdateActivity(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "记录ID格式错误",
		})
		return
	}

	var req models.ActivityUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误: " + err.Error(),
		})
		return
	}

	if err := h.activityService.UpdateActivity(id, &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "更新跟进记录失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新跟进记录成功",
	})
}

// UpdateActivityFeedback 更新跟进记录反馈
func (h *ActivityHandler) UpdateActivityFeedback(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "记录ID格式错误",
		})
		return
	}

	var req struct {
		Feedback     string `json:"feedback" binding:"required"`
		Satisfaction int    `json:"satisfaction"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误: " + err.Error(),
		})
		return
	}

	if err := h.activityService.UpdateActivityFeedback(id, req.Feedback, req.Satisfaction); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "更新反馈失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新反馈成功",
	})
}

// DeleteActivity 删除跟进记录
func (h *ActivityHandler) DeleteActivity(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "记录ID格式错误",
		})
		return
	}

	if err := h.activityService.DeleteActivity(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "删除跟进记录失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除跟进记录成功",
	})
}

// GetActivityStatistics 获取跟进记录统计信息
func (h *ActivityHandler) GetActivityStatistics(c *gin.Context) {
	customerIDStr := c.Param("customer_id")
	customerID, err := strconv.ParseUint(customerIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "客户ID格式错误",
		})
		return
	}

	stats, err := h.activityService.GetActivityStatistics(customerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取统计信息失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取统计信息成功",
		"data":    stats,
	})
}

// GetActivitiesNeedingFollowUp 获取需要跟进的记录
func (h *ActivityHandler) GetActivitiesNeedingFollowUp(c *gin.Context) {
	activities, err := h.activityService.GetActivitiesNeedingFollowUp()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取需要跟进的记录失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取需要跟进的记录成功",
		"data":    activities,
	})
}
