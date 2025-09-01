package services

import (
	"crm/internal/domain"
	"crm/internal/repository"
)

// UserService 用户服务
type UserService struct {
	repo repository.UserRepository
}

// NewUserService 创建用户服务
func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

// GetUsers 获取用户列表
func (s *UserService) GetUsers(req domain.UserQueryRequest) ([]domain.UserResponse, int64, error) {
	return s.repo.GetList(req)
}

// GetUserByID 根据ID获取用户
func (s *UserService) GetUserByID(id uint64) (*domain.UserResponse, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// 转换为响应格式
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

	return &response, nil
}

// GetAllActiveUsers 获取所有活跃用户（用于下拉选择）
func (s *UserService) GetAllActiveUsers() ([]domain.UserResponse, error) {
	return s.repo.GetAllActiveUsers()
}

// GetHomepageUserInfo 获取首页用户信息
func (s *UserService) GetHomepageUserInfo(id uint64) (*domain.HomepageUserResponse, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	todayFollowUps, err := s.repo.GetTodayFollowUpsCount(id)
	if err != nil {
		return nil, err
	}

	response := &domain.HomepageUserResponse{
		ID:             user.ID,
		Name:           user.Name,
		ShopName:       "四川一手货源",
		TodayRevenue:   0,
		TodayFollowUps: todayFollowUps,
	}

	return response, nil
}
