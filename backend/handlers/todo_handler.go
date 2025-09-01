package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"crm/models"
	"crm/services"
)

// TodoHandler 待办处理器
type TodoHandler struct {
	service *services.TodoService
}

// NewTodoHandler 创建待办处理器
func NewTodoHandler(service *services.TodoService) *TodoHandler {
	return &TodoHandler{service: service}
}

// CreateTodo 创建待办
func (h *TodoHandler) CreateTodo(c *gin.Context) {
	var req models.TodoCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 获取当前用户ID（从JWT或session中获取）
	creatorID := uint64(1) // 临时硬编码，实际应从认证中获取

	todo, err := h.service.CreateTodo(req, creatorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建待办失败"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "待办创建成功",
		"data":    todo,
	})
}

// GetTodos 获取待办列表
func (h *TodoHandler) GetTodos(c *gin.Context) {
	var req models.TodoQueryRequest
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

	responses, total, err := h.service.GetTodos(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询待办失败"})
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

// GetTodo 获取单个待办
func (h *TodoHandler) GetTodo(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的待办ID"})
		return
	}

	response, err := h.service.GetTodoByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "待办不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "查询待办失败"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}

// UpdateTodo 更新待办
func (h *TodoHandler) UpdateTodo(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的待办ID"})
		return
	}

	var req models.TodoUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 获取当前用户ID
	operatorID := uint64(1) // 临时硬编码

	todo, err := h.service.UpdateTodo(id, req, operatorID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "待办不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "更新待办失败"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "待办更新成功",
		"data":    todo,
	})
}

// DeleteTodo 删除待办（软删除）
func (h *TodoHandler) DeleteTodo(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的待办ID"})
		return
	}

	// 获取当前用户ID
	operatorID := uint64(1) // 临时硬编码

	err = h.service.DeleteTodo(id, operatorID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "待办不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "删除待办失败"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "待办删除成功"})
}

// CompleteTodo 完成待办
func (h *TodoHandler) CompleteTodo(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的待办ID"})
		return
	}

	// 获取当前用户ID
	operatorID := uint64(1) // 临时硬编码

	todo, err := h.service.CompleteTodo(id, operatorID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "待办不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "完成待办失败"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "待办已完成",
		"data":    todo,
	})
}

// CancelTodo 取消待办
func (h *TodoHandler) CancelTodo(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的待办ID"})
		return
	}

	// 获取当前用户ID
	operatorID := uint64(1) // 临时硬编码

	todo, err := h.service.CancelTodo(id, operatorID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "待办不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "取消待办失败"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "待办已取消",
		"data":    todo,
	})
}

// GetTodoStats 获取待办统计
func (h *TodoHandler) GetTodoStats(c *gin.Context) {
	customerIDStr := c.Query("customer_id")
	executorIDStr := c.Query("executor_id")

	var customerID, executorID *uint64

	if customerIDStr != "" {
		if id, err := strconv.ParseUint(customerIDStr, 10, 64); err == nil {
			customerID = &id
		}
	}

	if executorIDStr != "" {
		if id, err := strconv.ParseUint(executorIDStr, 10, 64); err == nil {
			executorID = &id
		}
	}

	stats, err := h.service.GetTodoStats(customerID, executorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取统计数据失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": stats})
}

