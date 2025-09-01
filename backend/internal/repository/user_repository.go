package repository

import (
	"crm/internal/domain"
	"time"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository 创建用户仓储实例
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// GetByID 根据ID获取用户
func (r *userRepository) GetByID(id uint64) (*domain.User, error) {
	var user domain.User
	err := r.db.Preload("Manager").
		Where("id = ? AND is_deleted = ?", id, false).First(&user).Error
	return &user, err
}

// GetList 获取用户列表
func (r *userRepository) GetList(req domain.UserQueryRequest) ([]domain.UserResponse, int64, error) {
	// 构建查询
	query := r.db.Model(&domain.User{}).Where("is_deleted = ?", false)

	// 添加筛选条件
	if req.Department != "" {
		query = query.Where("department = ?", req.Department)
	}
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}
	if req.Name != "" {
		query = query.Where("name LIKE ?", "%"+req.Name+"%")
	}

	// 获取总数
	var total int64
	query.Count(&total)

	// 分页查询
	var users []domain.User
	offset := (req.Page - 1) * req.PageSize
	err := query.Preload("Manager").
		Order("created_at DESC").Offset(offset).Limit(req.PageSize).Find(&users).Error

	if err != nil {
		return nil, 0, err
	}

	// 转换为响应格式
	var responses []domain.UserResponse
	for _, user := range users {
		response := domain.UserResponse{
			ID:         user.ID,
			Name:       user.Name,
			Department: user.Department,
			Position:   user.Position,
			Email:      user.Email,
			Phone:      user.Phone,
			Status:     user.Status,
			AvatarURL:  user.AvatarURL,
		}
		if user.Manager != nil {
			response.ManagerName = user.Manager.Name
		}
		responses = append(responses, response)
	}

	return responses, total, nil
}

// Create 创建用户
func (r *userRepository) Create(user *domain.User) error {
	// 设置默认值
	if user.Status == "" {
		user.Status = "active"
	}
	return r.db.Create(user).Error
}

// Update 更新用户
func (r *userRepository) Update(user *domain.User) error {
	return r.db.Save(user).Error
}

// Delete 删除用户（软删除）
func (r *userRepository) Delete(id uint64) error {
	result := r.db.Model(&domain.User{}).Where("id = ? AND is_deleted = ?", id, false).
		Updates(map[string]interface{}{
			"is_deleted": true,
		})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// GetByEmail 根据邮箱获取用户
func (r *userRepository) GetByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("email = ? AND is_deleted = ?", email, false).First(&user).Error
	return &user, err
}

// GetByPhone 根据手机号获取用户
func (r *userRepository) GetByPhone(phone string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("phone = ? AND is_deleted = ?", phone, false).First(&user).Error
	return &user, err
}

// GetAllActiveUsers 获取所有活跃用户
func (r *userRepository) GetAllActiveUsers() ([]domain.UserResponse, error) {
	var users []domain.User
	err := r.db.Where("is_deleted = ? AND (status = ? OR status = ?)", false, "active", "在职").
		Order("name ASC").Find(&users).Error

	if err != nil {
		return nil, err
	}

	// 转换为响应格式
	var responses []domain.UserResponse
	for _, user := range users {
		response := domain.UserResponse{
			ID:   user.ID,
			Name: user.Name,
		}
		responses = append(responses, response)
	}

	return responses, nil
}

// GetTodayFollowUpsCount 获取今日跟进数量
func (r *userRepository) GetTodayFollowUpsCount(userID uint64) (int, error) {
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	var count int64
	err := r.db.Model(&domain.Todo{}).
		Where("(creator_id = ? OR executor_id = ?) AND planned_time >= ? AND planned_time < ? AND is_deleted = ?",
			userID, userID, startOfDay, endOfDay, false).
		Count(&count).Error

	return int(count), err
}
