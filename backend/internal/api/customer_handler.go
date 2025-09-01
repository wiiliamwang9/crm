package api

import (
	"crm/internal/domain"
	"crm/internal/middleware"
	"crm/internal/services"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CustomerHandler 客户处理器
type CustomerHandler struct {
	service *services.CustomerService
}

// NewCustomerHandler 创建客户处理器实例
func NewCustomerHandler(service *services.CustomerService) *CustomerHandler {
	return &CustomerHandler{service: service}
}

// findCustomerByID 通用的根据ID查找客户的辅助函数
func (h *CustomerHandler) findCustomerByID(c *gin.Context, id string) (*domain.Customer, bool) {
	customerID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		middleware.ErrorResponse(c, middleware.InvalidParams, "无效的客户ID")
		return nil, false
	}

	customer, err := h.service.GetCustomerByID(customerID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			middleware.ErrorResponse(c, middleware.NotFound, "客户不存在")
		} else {
			middleware.ErrorResponse(c, middleware.DatabaseError, "查询客户失败")
		}
		return nil, false
	}
	return customer, true
}

// CreateCustomer 创建客户
func (h *CustomerHandler) CreateCustomer(c *gin.Context) {
	var req domain.CustomerRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ErrorResponse(c, middleware.ValidationError, "数据验证失败: "+err.Error())
		return
	}

	customer, err := h.service.CreateCustomer(req)
	if err != nil {
		middleware.ErrorResponse(c, middleware.DatabaseError, "创建客户失败")
		return
	}

	middleware.SuccessResponse(c, domain.CustomerToResponse(customer))
}

// GetCustomer 获取单个客户
func (h *CustomerHandler) GetCustomer(c *gin.Context) {
	id := c.Param("id")
	customer, ok := h.findCustomerByID(c, id)
	if !ok {
		return
	}

	middleware.SuccessResponse(c, domain.CustomerToResponse(customer))
}

// UpdateCustomer 更新客户
func (h *CustomerHandler) UpdateCustomer(c *gin.Context) {
	id := c.Param("id")
	customerID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		middleware.ErrorResponse(c, middleware.InvalidParams, "无效的客户ID")
		return
	}

	var req domain.CustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ErrorResponse(c, middleware.ValidationError, "数据验证失败: "+err.Error())
		return
	}

	customer, err := h.service.UpdateCustomer(customerID, req)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			middleware.ErrorResponse(c, middleware.NotFound, "客户不存在")
		} else {
			middleware.ErrorResponse(c, middleware.DatabaseError, "更新客户失败")
		}
		return
	}

	middleware.SuccessResponse(c, domain.CustomerToResponse(customer))
}

// DeleteCustomer 删除客户
func (h *CustomerHandler) DeleteCustomer(c *gin.Context) {
	id := c.Param("id")
	customerID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		middleware.ErrorResponse(c, middleware.InvalidParams, "无效的客户ID")
		return
	}

	err = h.service.DeleteCustomer(customerID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			middleware.ErrorResponse(c, middleware.NotFound, "客户不存在")
		} else {
			middleware.ErrorResponse(c, middleware.DatabaseError, "删除客户失败")
		}
		return
	}

	middleware.SuccessResponse(c, gin.H{"message": "客户删除成功"})
}

// GetCustomers 获取客户列表
func (h *CustomerHandler) GetCustomers(c *gin.Context) {
	var req domain.CustomerListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		middleware.ErrorResponse(c, middleware.InvalidParams, "参数错误: "+err.Error())
		return
	}

	customers, total, err := h.service.GetCustomers(req)
	if err != nil {
		middleware.ErrorResponse(c, middleware.DatabaseError, "查询客户列表失败")
		return
	}

	response := gin.H{
		"customers": customers,
		"total":     total,
		"page":      req.Page,
		"limit":     req.Limit,
	}

	middleware.SuccessResponse(c, response)
}

// PostSearchCustomers 搜索客户
func (h *CustomerHandler) PostSearchCustomers(c *gin.Context) {
	var req domain.SearchRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ErrorResponse(c, middleware.InvalidParams, "参数错误: "+err.Error())
		return
	}

	if req.Limit <= 0 || req.Limit > 100 {
		req.Limit = 50
	}
	if req.Page <= 0 {
		req.Page = 1
	}

	c.Header("X-Request-ID", strconv.FormatInt(req.Timestamp, 10))

	customers, total, err := h.service.SearchCustomers(req)
	if err != nil {
		middleware.ErrorResponse(c, middleware.DatabaseError, "搜索客户失败")
		return
	}

	response := gin.H{
		"customers":   customers,
		"total":       total,
		"page":        req.Page,
		"limit":       req.Limit,
		"search_term": req.Search,
	}

	middleware.SuccessResponse(c, response)
}

// GetSpecialCustomers 获取特殊客户（半年未下单、一直未下单）
func (h *CustomerHandler) GetSpecialCustomers(c *gin.Context) {
	customerType := c.Query("type") // "no_order_half_year" 或 "never_ordered"
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}

	customers, total, err := h.service.GetSpecialCustomers(customerType, page, pageSize)
	if err != nil {
		middleware.ErrorResponse(c, middleware.InvalidParams, err.Error())
		return
	}

	response := gin.H{
		"customers":  customers,
		"total":      total,
		"page":       page,
		"page_size":  pageSize,
		"total_page": (total + int64(pageSize) - 1) / int64(pageSize),
	}

	middleware.SuccessResponse(c, response)
}

// UpdateCustomerFavors 更新客户偏好
func (h *CustomerHandler) UpdateCustomerFavors(c *gin.Context) {
	id := c.Param("id")
	customerID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		middleware.ErrorResponse(c, middleware.InvalidParams, "无效的客户ID")
		return
	}

	var req domain.CustomerFavorsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ErrorResponse(c, middleware.ValidationError, "数据验证失败: "+err.Error())
		return
	}

	customer, err := h.service.UpdateCustomerFavors(customerID, req.Favors)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			middleware.ErrorResponse(c, middleware.NotFound, "客户不存在")
		} else {
			middleware.ErrorResponse(c, middleware.DatabaseError, "更新客户偏好失败")
		}
		return
	}

	response := gin.H{
		"message": "客户偏好更新成功",
		"favors":  customer.Favors,
	}

	middleware.SuccessResponse(c, response)
}

// UpdateCustomerRemark 更新客户备注
func (h *CustomerHandler) UpdateCustomerRemark(c *gin.Context) {
	id := c.Param("id")
	customerID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		middleware.ErrorResponse(c, middleware.InvalidParams, "无效的客户ID")
		return
	}

	var req struct {
		Remark string `json:"remark" binding:"max=1000"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ErrorResponse(c, middleware.ValidationError, "数据验证失败: "+err.Error())
		return
	}

	customer, err := h.service.UpdateCustomerRemark(customerID, req.Remark)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			middleware.ErrorResponse(c, middleware.NotFound, "客户不存在")
		} else {
			middleware.ErrorResponse(c, middleware.DatabaseError, "更新客户备注失败")
		}
		return
	}

	response := gin.H{
		"message": "客户备注更新成功",
		"remark":  customer.Remark,
	}

	middleware.SuccessResponse(c, response)
}

// UpdateCustomerSystemTags 更新客户系统标签
func (h *CustomerHandler) UpdateCustomerSystemTags(c *gin.Context) {
	id := c.Param("id")
	customerID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		middleware.ErrorResponse(c, middleware.InvalidParams, "无效的客户ID")
		return
	}

	var req struct {
		SystemTags []int64 `json:"system_tags" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ErrorResponse(c, middleware.ValidationError, "数据验证失败: "+err.Error())
		return
	}

	customer, err := h.service.UpdateCustomerSystemTags(customerID, req.SystemTags)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			middleware.ErrorResponse(c, middleware.NotFound, "客户不存在")
		} else {
			middleware.ErrorResponse(c, middleware.DatabaseError, "更新客户系统标签失败")
		}
		return
	}

	response := gin.H{
		"message":     "客户系统标签更新成功",
		"system_tags": customer.SystemTags,
	}

	middleware.SuccessResponse(c, response)
}

// ExportCustomersExcel 导出客户Excel
func (h *CustomerHandler) ExportCustomersExcel(c *gin.Context) {
	middleware.ErrorResponse(c, middleware.InternalError, "Excel导出功能暂未实现，请使用原有功能")
}

// UploadCustomersExcel 上传客户Excel
func (h *CustomerHandler) UploadCustomersExcel(c *gin.Context) {
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		middleware.ErrorResponse(c, middleware.InvalidParams, "上传文件失败")
		return
	}
	defer file.Close()

	middleware.ErrorResponse(c, middleware.InternalError, "Excel导入功能暂未实现，请使用原有功能")
}

// GetCustomerStateDesc 获取客户状态描述
func GetCustomerStateDesc(state int) string {
	switch state {
	case 0:
		return "未开发"
	case 1:
		return "意向客户"
	case 2:
		return "跟进中"
	case 3:
		return "已开发"
	case 4:
		return "成交客户"
	case 5:
		return "流失客户"
	default:
		return "未知"
	}
}

// GetCustomerLevelDesc 获取客户级别描述
func GetCustomerLevelDesc(level int) string {
	switch level {
	case 0:
		return "未分级"
	case 1:
		return "A级客户"
	case 2:
		return "B级客户"
	case 3:
		return "C级客户"
	case 4:
		return "D级客户"
	default:
		return "未知"
	}
}
