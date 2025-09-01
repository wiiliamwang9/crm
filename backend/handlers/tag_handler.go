package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"crm/models"
	"crm/services"
)

// TagHandler 标签处理器
type TagHandler struct {
	service *services.TagService
}

// NewTagHandler 创建标签处理器
func NewTagHandler(service *services.TagService) *TagHandler {
	return &TagHandler{service: service}
}

// GetDimensions 获取所有维度及其标签
func (h *TagHandler) GetDimensions(c *gin.Context) {
	responses, err := h.service.GetDimensions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询维度失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": responses})
}

// GetDimension 获取单个维度
func (h *TagHandler) GetDimension(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的维度ID"})
		return
	}

	dimension, err := h.service.GetDimensionByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "维度不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "查询维度失败"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": dimension})
}

// CreateDimension 创建维度
func (h *TagHandler) CreateDimension(c *gin.Context) {
	var req models.TagDimensionCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dimension, err := h.service.CreateDimension(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建维度失败"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": dimension})
}

// UpdateDimension 更新维度
func (h *TagHandler) UpdateDimension(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的维度ID"})
		return
	}

	var req models.TagDimensionUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.service.UpdateDimension(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新维度失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

// DeleteDimension 删除维度
func (h *TagHandler) DeleteDimension(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的维度ID"})
		return
	}

	err = h.service.DeleteDimension(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除维度失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// GetTags 获取标签列表
func (h *TagHandler) GetTags(c *gin.Context) {
	var req models.TagQueryRequest
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

	responses, total, err := h.service.GetTags(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询标签失败"})
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

// GetTagsByDimensionID 根据维度ID获取标签
func (h *TagHandler) GetTagsByDimensionID(c *gin.Context) {
	idStr := c.Param("dimension_id")
	dimensionID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的维度ID"})
		return
	}

	tags, err := h.service.GetTagsByDimensionID(dimensionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询标签失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": tags})
}

// GetTag 获取单个标签
func (h *TagHandler) GetTag(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的标签ID"})
		return
	}

	tag, err := h.service.GetTagByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "标签不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "查询标签失败"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": tag})
}

// CreateTag 创建标签
func (h *TagHandler) CreateTag(c *gin.Context) {
	var req models.TagCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tag, err := h.service.CreateTag(req)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusBadRequest, gin.H{"error": "维度不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "创建标签失败"})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": tag})
}

// UpdateTag 更新标签
func (h *TagHandler) UpdateTag(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的标签ID"})
		return
	}

	var req models.TagUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.service.UpdateTag(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新标签失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

// DeleteTag 删除标签
func (h *TagHandler) DeleteTag(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的标签ID"})
		return
	}

	err = h.service.DeleteTag(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除标签失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// GetAllActiveTags 获取所有活跃标签（用于下拉选择）
func (h *TagHandler) GetAllActiveTags(c *gin.Context) {
	tags, err := h.service.GetAllActiveTags()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询标签失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": tags})
}
