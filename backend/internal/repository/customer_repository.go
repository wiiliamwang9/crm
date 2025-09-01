package repository

import (
	"crm/internal/domain"
	"errors"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

type customerRepository struct {
	db *gorm.DB
}

// NewCustomerRepository 创建客户仓储实例
func NewCustomerRepository(db *gorm.DB) CustomerRepository {
	return &customerRepository{db: db}
}

// GetByID 根据ID获取客户
func (r *customerRepository) GetByID(id uint64) (*domain.Customer, error) {
	var customer domain.Customer
	err := r.db.First(&customer, id).Error
	return &customer, err
}

// GetList 获取客户列表
func (r *customerRepository) GetList(req domain.CustomerListRequest) ([]*domain.CustomerResponse, int64, error) {
	return r.Search(domain.SearchRequest{
		Search: req.Search,
		Page:   req.Page,
		Limit:  req.Limit,
	})
}

// Search 搜索客户（带系统标签筛选）
func (r *customerRepository) Search(req domain.SearchRequest) ([]*domain.CustomerResponse, int64, error) {
	var customers []domain.Customer
	offset := (req.Page - 1) * req.Limit

	query := r.db.Model(&domain.Customer{})

	if req.Search != "" {
		searchTerm := "%" + req.Search + "%"
		query = query.Where("name LIKE ? OR contact_name LIKE ? OR remark LIKE ? OR phones::text LIKE ? OR wechats::text LIKE ? OR address LIKE ? OR province LIKE ? OR city LIKE ? OR district LIKE ?",
			searchTerm, searchTerm, searchTerm, searchTerm, searchTerm, searchTerm, searchTerm, searchTerm, searchTerm)
	}

	if len(req.SystemTags) > 0 {
		query = query.Where("system_tags && ARRAY[" + formatIntArray(req.SystemTags) + "]")
	}

	var total int64
	countQuery := *query
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(req.Limit).Order("created_at DESC").Find(&customers).Error; err != nil {
		return nil, 0, err
	}

	customerList := make([]*domain.CustomerResponse, len(customers))
	for i, customer := range customers {
		customerList[i] = domain.CustomerToResponse(&customer)
	}

	return customerList, total, nil
}

// Create 创建客户
func (r *customerRepository) Create(customer *domain.Customer) error {
	return r.db.Create(customer).Error
}

// Update 更新客户
func (r *customerRepository) Update(customer *domain.Customer) error {
	return r.db.Save(customer).Error
}

// Delete 删除客户
func (r *customerRepository) Delete(id uint64) error {
	result := r.db.Delete(&domain.Customer{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// UpdateFavors 更新客户偏好
func (r *customerRepository) UpdateFavors(id uint64, favors []map[string]interface{}) (*domain.Customer, error) {
	var customer domain.Customer
	if err := r.db.First(&customer, id).Error; err != nil {
		return nil, err
	}

	customer.Favors = domain.JSONB(map[string]interface{}{"favors": favors})

	if err := r.db.Save(&customer).Error; err != nil {
		return nil, err
	}

	return &customer, nil
}

// UpdateRemark 更新客户备注
func (r *customerRepository) UpdateRemark(id uint64, remark string) (*domain.Customer, error) {
	var customer domain.Customer
	if err := r.db.First(&customer, id).Error; err != nil {
		return nil, err
	}

	customer.Remark = remark

	if err := r.db.Save(&customer).Error; err != nil {
		return nil, err
	}

	return &customer, nil
}

// UpdateSystemTags 更新客户系统标签
func (r *customerRepository) UpdateSystemTags(id uint64, systemTags []int64) (*domain.Customer, error) {
	var customer domain.Customer
	if err := r.db.First(&customer, id).Error; err != nil {
		return nil, err
	}

	customer.SystemTags = systemTags

	if err := r.db.Save(&customer).Error; err != nil {
		return nil, err
	}

	return &customer, nil
}

// GetSpecialCustomers 获取特殊客户（半年未下单、一直未下单）
func (r *customerRepository) GetSpecialCustomers(customerType string, page, pageSize int) ([]*domain.CustomerResponse, int64, error) {
	var customers []domain.Customer
	offset := (page - 1) * pageSize
	query := r.db.Model(&domain.Customer{})

	now := time.Now()
	switch customerType {
	case "no_order_half_year":
		sixMonthsAgo := now.AddDate(0, -6, 0)
		query = query.Where("last_order_date IS NOT NULL AND last_order_date < ?", sixMonthsAgo)
	case "never_ordered":
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

	customerList := make([]*domain.CustomerResponse, len(customers))
	for i, customer := range customers {
		customerList[i] = domain.CustomerToResponse(&customer)
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
