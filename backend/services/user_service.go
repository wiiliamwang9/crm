package services

import (
	"gorm.io/gorm"

	"crm/models"
)

// UserService 用户服务
type UserService struct {
	db *gorm.DB
}

// NewUserService 创建用户服务
func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

// GetUsers 获取用户列表
func (s *UserService) GetUsers(req models.UserQueryRequest) ([]models.UserResponse, int64, error) {
	// 构建查询
	query := s.db.Model(&models.User{}).Where("is_deleted = ?", false)

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
	var users []models.User
	offset := (req.Page - 1) * req.PageSize
	err := query.Preload("Manager").
		Order("created_at DESC").Offset(offset).Limit(req.PageSize).Find(&users).Error

	if err != nil {
		return nil, 0, err
	}

	// 转换为响应格式
	var responses []models.UserResponse
	for _, user := range users {
		response := models.UserResponse{
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

// GetUserByID 根据ID获取用户
func (s *UserService) GetUserByID(id uint64) (*models.UserResponse, error) {
	var user models.User
	err := s.db.Preload("Manager").
		Where("id = ? AND is_deleted = ?", id, false).First(&user).Error

	if err != nil {
		return nil, err
	}

	// 转换为响应格式
	response := models.UserResponse{
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

	return &response, nil
}

// GetAllActiveUsers 获取所有活跃用户（用于下拉选择）
func (s *UserService) GetAllActiveUsers() ([]models.UserResponse, error) {
	var users []models.User
	err := s.db.Where("is_deleted = ? AND (status = ? OR status = ?)", false, "active", "在职").
		Order("name ASC").Find(&users).Error

	if err != nil {
		return nil, err
	}

	// 转换为响应格式
	var responses []models.UserResponse
	for _, user := range users {
		response := models.UserResponse{
			ID:   user.ID,
			Name: user.Name,
		}
		responses = append(responses, response)
	}

	return responses, nil
}
