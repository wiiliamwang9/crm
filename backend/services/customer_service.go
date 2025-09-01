package services

import (
	"gorm.io/gorm"

	"crm/models"
)

// CustomerService 客户服务
type CustomerService struct {
	db *gorm.DB
}

// NewCustomerService 创建客户服务
func NewCustomerService(db *gorm.DB) *CustomerService {
	return &CustomerService{db: db}
}

// GetCustomers 获取客户列表
func (s *CustomerService) GetCustomers(page, limit int, search string) ([]models.Customer, int64, error) {
	var customers []models.Customer
	var total int64

	offset := (page - 1) * limit

	query := s.db.Model(&models.Customer{})

	if search != "" {
		query = query.Where("name LIKE ? OR contact_name LIKE ? OR remark LIKE ? OR phones::text LIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	// 获取总数
	query.Count(&total)

	// 分页查询
	err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&customers).Error

	return customers, total, err
}

// GetCustomerByID 根据ID获取客户
func (s *CustomerService) GetCustomerByID(id uint64) (*models.Customer, error) {
	var customer models.Customer
	err := s.db.First(&customer, id).Error
	return &customer, err
}

// CreateCustomer 创建客户
func (s *CustomerService) CreateCustomer(customer *models.Customer) error {
	return s.db.Create(customer).Error
}

// UpdateCustomer 更新客户
func (s *CustomerService) UpdateCustomer(customer *models.Customer) error {
	return s.db.Save(customer).Error
}

// DeleteCustomer 删除客户
func (s *CustomerService) DeleteCustomer(id uint64) error {
	return s.db.Delete(&models.Customer{}, id).Error
}