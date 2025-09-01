package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"crm/models"
	"crm/services"
)

// ReminderHandler 提醒处理器
type ReminderHandler struct {
	service *services.ReminderService
}

// NewReminderHandler 创建提醒处理器
func NewReminderHandler(service *services.ReminderService) *ReminderHandler {
	return &ReminderHandler{service: service}
}

// CreateReminder 创建提醒
func (h *ReminderHandler) CreateReminder(c *gin.Context) {
	var req models.ReminderCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reminder, err := h.service.CreateReminder(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建提醒失败"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "提醒创建成功",
		"data":    reminder,
	})
}

// GetReminders 获取提醒列表
func (h *ReminderHandler) GetReminders(c *gin.Context) {
	var req models.ReminderQueryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 设置默认分页参数
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	responses, total, err := h.service.GetReminders(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询提醒失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": responses,
		"pagination": gin.H{
			"page":       req.Page,
			"page_size":  req.PageSize,
			"total":      total,
			"total_page": (total + int64(req.PageSize) - 1) / int64(req.PageSize),
		},
	})
}

// GetReminder 获取单个提醒
func (h *ReminderHandler) GetReminder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的提醒ID"})
		return
	}

	response, err := h.service.GetReminderByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "提醒不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "查询提醒失败"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}

// UpdateReminder 更新提醒
func (h *ReminderHandler) UpdateReminder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的提醒ID"})
		return
	}

	var req models.ReminderUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reminder, err := h.service.UpdateReminder(id, req)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "提醒不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "更新提醒失败"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "提醒更新成功",
		"data":    reminder,
	})
}

// DeleteReminder 删除提醒
func (h *ReminderHandler) DeleteReminder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的提醒ID"})
		return
	}

	err = h.service.DeleteReminder(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "提醒不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "删除提醒失败"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "提醒删除成功"})
}

// CancelReminder 取消提醒
func (h *ReminderHandler) CancelReminder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的提醒ID"})
		return
	}

	reminder, err := h.service.CancelReminder(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "提醒不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "取消提醒失败"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "提醒已取消",
		"data":    reminder,
	})
}

// GetReminderStats 获取提醒统计
func (h *ReminderHandler) GetReminderStats(c *gin.Context) {
	userIDStr := c.Query("user_id")
	var userID *uint64

	if userIDStr != "" {
		if id, err := strconv.ParseUint(userIDStr, 10, 64); err == nil {
			userID = &id
		}
	}

	stats, err := h.service.GetReminderStats(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取统计数据失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": stats})
}

// GetUserReminderConfig 获取用户提醒配置
func (h *ReminderHandler) GetUserReminderConfig(c *gin.Context) {
	// 获取当前用户ID（从JWT或session中获取）
	userID := uint64(1) // 临时硬编码，实际应从认证中获取

	config, err := h.service.GetUserReminderConfig(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取配置失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": config})
}

// UpdateUserReminderConfig 更新用户提醒配置
func (h *ReminderHandler) UpdateUserReminderConfig(c *gin.Context) {
	// 获取当前用户ID（从JWT或session中获取）
	userID := uint64(1) // 临时硬编码，实际应从认证中获取

	var config models.ReminderConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.UpdateUserReminderConfig(userID, &config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新配置失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "配置更新成功"})
}

// ProcessReminders 手动触发处理提醒
func (h *ReminderHandler) ProcessReminders(c *gin.Context) {
	err := h.service.ProcessPendingReminders()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "处理提醒失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "提醒处理完成"})
}