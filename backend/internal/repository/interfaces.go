package repository

import (
	"crm/internal/domain"
)

// CustomerRepository 客户信息数据库接口
type CustomerRepository interface {
	GetByID(id uint64) (*domain.Customer, error)
	GetList(req domain.CustomerListRequest) ([]*domain.CustomerResponse, int64, error)
	Search(req domain.SearchRequest) ([]*domain.CustomerResponse, int64, error)
	Create(customer *domain.Customer) error
	Update(customer *domain.Customer) error
	Delete(id uint64) error
	UpdateFavors(id uint64, favors []map[string]interface{}) (*domain.Customer, error)
	UpdateRemark(id uint64, remark string) (*domain.Customer, error)
	UpdateSystemTags(id uint64, systemTags []int64) (*domain.Customer, error)
	GetSpecialCustomers(customerType string, page, pageSize int) ([]*domain.CustomerResponse, int64, error)
}

// TodoRepository 待办信息数据库接口
type TodoRepository interface {
	GetByID(id uint64) (*domain.Todo, error)
	GetList(req domain.TodoQueryRequest) ([]domain.TodoResponse, int64, error)
	Create(todo *domain.Todo) error
	Update(todo *domain.Todo) error
	Delete(id uint64) error
	UpdateStatus(id uint64, status domain.TodoStatus) error
	GetOverdueTodos() ([]domain.Todo, error)
	CreateLog(log *domain.TodoLog) error
}

// ActivityRepository 活动信息数据库接口
type ActivityRepository interface {
	GetByID(id uint64) (*domain.Activity, error)
	GetList(req domain.ActivityQueryRequest) ([]domain.ActivityResponse, int64, error)
	Create(activity *domain.Activity) error
	Update(activity *domain.Activity) error
	Delete(id uint64) error
	GetByCustomerID(customerID uint64, page, pageSize int) ([]domain.ActivityResponse, int64, error)
	GetStatistics(customerID uint64) (map[string]interface{}, error)
	GetNeedingFollowUp() ([]domain.ActivityResponse, error)
}

// UserRepository 用户信息数据库接口
type UserRepository interface {
	GetByID(id uint64) (*domain.User, error)
	GetList(req domain.UserQueryRequest) ([]domain.UserResponse, int64, error)
	Create(user *domain.User) error
	Update(user *domain.User) error
	Delete(id uint64) error
	GetByEmail(email string) (*domain.User, error)
	GetByPhone(phone string) (*domain.User, error)
	GetAllActiveUsers() ([]domain.UserResponse, error)
	GetTodayFollowUpsCount(userID uint64) (int, error)
}

// TagRepository 标签信息数据库接口
type TagRepository interface {
	GetByID(id uint64) (*domain.Tag, error)
	GetList(req domain.TagQueryRequest) ([]domain.TagResponse, int64, error)
	Create(tag *domain.Tag) error
	Update(tag *domain.Tag) error
	Delete(id uint64) error
	GetByDimensionID(dimensionID uint64) ([]domain.Tag, error)
	GetAllActiveTags() ([]domain.TagResponse, error)
}

// TagDimensionRepository 标签维度信息数据库接口
type TagDimensionRepository interface {
	GetByID(id uint64) (*domain.TagDimension, error)
	GetList(page, pageSize int) ([]domain.TagDimensionResponse, int64, error)
	Create(dimension *domain.TagDimension) error
	Update(dimension *domain.TagDimension) error
	Delete(id uint64) error
	GetWithTags() ([]domain.TagDimensionResponse, error)
}

// ReminderRepository 提醒信息数据库接口
type ReminderRepository interface {
	GetByID(id uint64) (*domain.Reminder, error)
	GetList(req domain.ReminderQueryRequest) ([]domain.ReminderResponse, int64, error)
	Create(reminder *domain.Reminder) error
	Update(reminder *domain.Reminder) error
	Delete(id uint64) error
	GetPendingReminders() ([]domain.Reminder, error)
	UpdateStatus(id uint64, status domain.ReminderStatus) error
	GetStats(userID *uint64) (map[string]int64, error)
	GetReminderConfig(userID uint64) (*domain.ReminderConfig, error)
	UpdateReminderConfig(config *domain.ReminderConfig) error
	GetReminderTemplate(reminderType domain.ReminderType, isDefault bool) (*domain.ReminderTemplate, error)
}

// Repository 聚合接口
type Repository struct {
	Customer     CustomerRepository
	Todo         TodoRepository
	Activity     ActivityRepository
	User         UserRepository
	Tag          TagRepository
	TagDimension TagDimensionRepository
	Reminder     ReminderRepository
}
