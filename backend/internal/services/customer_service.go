package services

import (
	"crm/internal/domain"
	"crm/internal/repository"
)

// CustomerService 客户服务
type CustomerService struct {
	repo repository.CustomerRepository
}

// NewCustomerService 创建客户服务
func NewCustomerService(repo repository.CustomerRepository) *CustomerService {
	return &CustomerService{repo: repo}
}

// GetCustomers 获取客户列表
func (s *CustomerService) GetCustomers(req domain.CustomerListRequest) ([]*domain.CustomerResponse, int64, error) {
	return s.repo.GetList(req)
}

// GetCustomerByID 根据ID获取客户
func (s *CustomerService) GetCustomerByID(id uint64) (*domain.Customer, error) {
	return s.repo.GetByID(id)
}

// CreateCustomer 创建客户
func (s *CustomerService) CreateCustomer(req domain.CustomerRequest) (*domain.Customer, error) {
	customer := req.ToModel()
	if err := s.repo.Create(customer); err != nil {
		return nil, err
	}
	return customer, nil
}

// UpdateCustomer 更新客户
func (s *CustomerService) UpdateCustomer(id uint64, req domain.CustomerRequest) (*domain.Customer, error) {
	customer, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	updateData := req.ToModel()
	updateData.ID = customer.ID
	updateData.CreatedAt = customer.CreatedAt

	if err := s.repo.Update(updateData); err != nil {
		return nil, err
	}

	return updateData, nil
}

// DeleteCustomer 删除客户
func (s *CustomerService) DeleteCustomer(id uint64) error {
	return s.repo.Delete(id)
}

// SearchCustomers 搜索客户（带系统标签筛选）
func (s *CustomerService) SearchCustomers(req domain.SearchRequest) ([]*domain.CustomerResponse, int64, error) {
	return s.repo.Search(req)
}

// GetSpecialCustomers 获取特殊客户（半年未下单、一直未下单）
func (s *CustomerService) GetSpecialCustomers(customerType string, page, pageSize int) ([]*domain.CustomerResponse, int64, error) {
	return s.repo.GetSpecialCustomers(customerType, page, pageSize)
}

// UpdateCustomerFavors 更新客户偏好
func (s *CustomerService) UpdateCustomerFavors(id uint64, favors []map[string]interface{}) (*domain.Customer, error) {
	return s.repo.UpdateFavors(id, favors)
}

// UpdateCustomerRemark 更新客户备注
func (s *CustomerService) UpdateCustomerRemark(id uint64, remark string) (*domain.Customer, error) {
	return s.repo.UpdateRemark(id, remark)
}

// UpdateCustomerSystemTags 更新客户系统标签
func (s *CustomerService) UpdateCustomerSystemTags(id uint64, systemTags []int64) (*domain.Customer, error) {
	return s.repo.UpdateSystemTags(id, systemTags)
}
