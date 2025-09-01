package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"crm/models"
	"crm/services"
)

// UserHandler 用户处理器
type UserHandler struct {
	service *services.UserService
}

// NewUserHandler 创建用户处理器
func NewUserHandler(service *services.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// GetUsers 获取用户列表
func (h *UserHandler) GetUsers(c *gin.Context) {
	var req models.UserQueryRequest
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

	responses, total, err := h.service.GetUsers(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询用户失败"})
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

// GetUser 获取单个用户
func (h *UserHandler) GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	response, err := h.service.GetUserByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "查询用户失败"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}

// GetAllActiveUsers 获取所有活跃用户（用于下拉选择）
func (h *UserHandler) GetAllActiveUsers(c *gin.Context) {
	responses, err := h.service.GetAllActiveUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询用户失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": responses})
}
