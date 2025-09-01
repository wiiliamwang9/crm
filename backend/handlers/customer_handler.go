package handlers

import (
	"crm/config"
	"crm/dto"
	"crm/middleware"
	"crm/models"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// CustomerHandler 客户处理器
type CustomerHandler struct{}

// NewCustomerHandler 创建客户处理器实例
func NewCustomerHandler() *CustomerHandler {
	return &CustomerHandler{}
}

// findCustomerByID 通用的根据ID查找客户的辅助函数
func (h *CustomerHandler) findCustomerByID(c *gin.Context, id string) (*models.Customer, bool) {
	var customer models.Customer
	if err := config.DB.First(&customer, id).Error; err != nil {
		middleware.ErrorResponse(c, middleware.NotFound, "客户不存在")
		return nil, false
	}
	return &customer, true
}

// CreateCustomer 创建客户
func (h *CustomerHandler) CreateCustomer(c *gin.Context) {
	var req dto.CustomerRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ErrorResponse(c, middleware.ValidationError, "数据验证失败: "+err.Error())
		return
	}

	customer := req.ToModel()

	if err := config.DB.Create(customer).Error; err != nil {
		middleware.ErrorResponse(c, middleware.DatabaseError, "创建客户失败")
		return
	}

	middleware.SuccessResponse(c, dto.CustomerToResponse(customer))
}

// GetCustomer 获取单个客户
func (h *CustomerHandler) GetCustomer(c *gin.Context) {
	id := c.Param("id")
	customer, ok := h.findCustomerByID(c, id)
	if !ok {
		return
	}

	middleware.SuccessResponse(c, dto.CustomerToResponse(customer))
}

// UpdateCustomer 更新客户
func (h *CustomerHandler) UpdateCustomer(c *gin.Context) {
	id := c.Param("id")
	customer, ok := h.findCustomerByID(c, id)
	if !ok {
		return
	}

	var req dto.CustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ErrorResponse(c, middleware.ValidationError, "数据验证失败: "+err.Error())
		return
	}

	updateData := req.ToModel()
	updateData.ID = customer.ID
	updateData.CreatedAt = customer.CreatedAt

	if err := config.DB.Save(updateData).Error; err != nil {
		middleware.ErrorResponse(c, middleware.DatabaseError, "更新客户失败")
		return
	}

	middleware.SuccessResponse(c, dto.CustomerToResponse(updateData))
}

// DeleteCustomer 删除客户
func (h *CustomerHandler) DeleteCustomer(c *gin.Context) {
	id := c.Param("id")

	result := config.DB.Delete(&models.Customer{}, id)
	if result.Error != nil {
		middleware.ErrorResponse(c, middleware.DatabaseError, "删除客户失败")
		return
	}

	if result.RowsAffected == 0 {
		middleware.ErrorResponse(c, middleware.NotFound, "客户不存在")
		return
	}

	middleware.SuccessResponse(c, gin.H{"message": "客户删除成功"})
}

// GetCustomers 获取客户列表
func (h *CustomerHandler) GetCustomers(c *gin.Context) {
	var req dto.CustomerListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		middleware.ErrorResponse(c, middleware.InvalidParams, "参数错误: "+err.Error())
		return
	}

	customers, total, err := h.queryCustomersWithSystemTags(req.Search, nil, req.Page, req.Limit)
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
	var req dto.SearchRequest

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

	customers, total, err := h.queryCustomersWithSystemTags(req.Search, req.SystemTags, req.Page, req.Limit)
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

	customers, total, err := h.querySpecialCustomers(customerType, page, pageSize)
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
	customer, ok := h.findCustomerByID(c, id)
	if !ok {
		return
	}

	var req dto.CustomerFavorsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ErrorResponse(c, middleware.ValidationError, "数据验证失败: "+err.Error())
		return
	}

	customer.Favors = models.JSONB(map[string]interface{}{"favors": req.Favors})

	if err := config.DB.Save(customer).Error; err != nil {
		middleware.ErrorResponse(c, middleware.DatabaseError, "更新客户偏好失败")
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
	customer, ok := h.findCustomerByID(c, id)
	if !ok {
		return
	}

	var req struct {
		Remark string `json:"remark" binding:"max=1000"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ErrorResponse(c, middleware.ValidationError, "数据验证失败: "+err.Error())
		return
	}

	customer.Remark = req.Remark

	if err := config.DB.Save(customer).Error; err != nil {
		middleware.ErrorResponse(c, middleware.DatabaseError, "更新客户备注失败")
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
	customer, ok := h.findCustomerByID(c, id)
	if !ok {
		return
	}

	var req struct {
		SystemTags []int64 `json:"system_tags" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ErrorResponse(c, middleware.ValidationError, "数据验证失败: "+err.Error())
		return
	}

	customer.SystemTags = req.SystemTags

	if err := config.DB.Save(customer).Error; err != nil {
		middleware.ErrorResponse(c, middleware.DatabaseError, "更新客户系统标签失败")
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
	var customers []models.Customer

	if err := config.DB.Order("created_at DESC").Find(&customers).Error; err != nil {
		middleware.ErrorResponse(c, middleware.DatabaseError, "查询客户数据失败")
		return
	}

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

// 辅助函数

// queryCustomersWithSystemTags 通用客户查询函数（支持系统标签筛选）
func (h *CustomerHandler) queryCustomersWithSystemTags(search string, systemTags []int64, page, limit int) ([]*dto.CustomerResponse, int64, error) {
	var customers []models.Customer
	offset := (page - 1) * limit

	query := config.DB.Model(&models.Customer{})

	if search != "" {
		searchTerm := "%" + search + "%"
		query = query.Where("name LIKE ? OR contact_name LIKE ? OR remark LIKE ? OR phones::text LIKE ? OR wechats::text LIKE ? OR address LIKE ? OR province LIKE ? OR city LIKE ? OR district LIKE ?",
			searchTerm, searchTerm, searchTerm, searchTerm, searchTerm, searchTerm, searchTerm, searchTerm, searchTerm)
	}

	if len(systemTags) > 0 {
		query = query.Where("system_tags && ARRAY[" + formatIntArray(systemTags) + "]")
	}

	var total int64
	countQuery := *query
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&customers).Error; err != nil {
		return nil, 0, err
	}

	customerList := make([]*dto.CustomerResponse, len(customers))
	for i, customer := range customers {
		customerList[i] = dto.CustomerToResponse(&customer)
	}

	return customerList, total, nil
}

// querySpecialCustomers 查询特殊客户
func (h *CustomerHandler) querySpecialCustomers(customerType string, page, pageSize int) ([]*dto.CustomerResponse, int64, error) {
	var customers []models.Customer
	offset := (page - 1) * pageSize
	query := config.DB.Model(&models.Customer{})

	now := time.Now()
	switch customerType {
	case "no_order_half_year":
		// 半年未下单：last_order_date不为空且距今超过6个月
		sixMonthsAgo := now.AddDate(0, -6, 0)
		query = query.Where("last_order_date IS NOT NULL AND last_order_date < ?", sixMonthsAgo)
	case "never_ordered":
		// 一直未下单：last_order_date为空
		query = query.Where("last_order_date IS NULL")
	default:
		return nil, 0, errors.New("无效的客户类型")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&customers).Error; err != nil {
		return nil, 0, err
	}

	// 即使没有找到客户，也返回空列表而不是错误
	customerList := make([]*dto.CustomerResponse, len(customers))
	for i, customer := range customers {
		customerList[i] = dto.CustomerToResponse(&customer)
	}

	return customerList, total, nil
}

// formatIntArray 将int64数组格式化为PostgreSQL数组字符串
func formatIntArray(arr []int64) string {
	if len(arr) == 0 {
		return ""
	}

	strArr := make([]string, len(arr))
	for i, v := range arr {
		strArr[i] = strconv.FormatInt(v, 10)
	}

	return strings.Join(strArr, ",")
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
