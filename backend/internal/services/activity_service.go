package services

import (
	models2 "crm/internal/domain"
	"crm/internal/repository"
	"fmt"
	"unicode/utf8"
)

type ActivityService struct {
	activityRepo repository.ActivityRepository
	todoRepo     repository.TodoRepository
}

func NewActivityService(activityRepo repository.ActivityRepository, todoRepo repository.TodoRepository) *ActivityService {
	return &ActivityService{
		activityRepo: activityRepo,
		todoRepo:     todoRepo,
	}
}

// CreateActivity 创建跟进记录
func (s *ActivityService) CreateActivity(req *models2.ActivityCreateRequest) (*models2.Activity, error) {
	// 验证UTF-8字符编码
	if !utf8.ValidString(req.Content) {
		return nil, fmt.Errorf("内容包含无效的UTF-8字符")
	}
	if !utf8.ValidString(req.Title) {
		return nil, fmt.Errorf("标题包含无效的UTF-8字符")
	}
	if !utf8.ValidString(req.Feedback) {
		return nil, fmt.Errorf("反馈包含无效的UTF-8字符")
	}
	if !utf8.ValidString(req.Remark) {
		return nil, fmt.Errorf("备注包含无效的UTF-8字符")
	}

	// 构建活动数据
	data := models2.ActivityData{
		Content:      req.Content,
		Result:       req.Result,
		Amount:       req.Amount,
		Cost:         req.Cost,
		Feedback:     req.Feedback,
		Satisfaction: req.Satisfaction,
	}

	activity := &models2.Activity{
		CustomerID:     req.CustomerID,
		UserID:         1, // TODO: 从上下文获取当前用户ID
		Kind:           req.Kind,
		Title:          req.Title,
		Remark:         req.Remark,
		Duration:       req.Duration,
		Location:       req.Location,
		NextFollowTime: req.NextFollowTime,
		Attachments:    req.Attachments,
	}

	// 设置数据结构
	if err := activity.SetDataStruct(data); err != nil {
		return nil, fmt.Errorf("设置数据结构失败: %w", err)
	}

	// 创建活动记录
	if err := s.activityRepo.Create(activity); err != nil {
		return nil, fmt.Errorf("创建跟进记录失败: %w", err)
	}

	// 如果需要创建待办事项
	if req.CreateTodo && req.TodoPlannedTime != nil {
		todo := &models2.Todo{
			CustomerID:  req.CustomerID,
			CreatorID:   1, // TODO: 从上下文获取当前用户ID
			ExecutorID:  *req.TodoExecutorID,
			Title:       req.TodoContent,
			Content:     req.Content,
			PlannedTime: *req.TodoPlannedTime,
			Status:      models2.TodoStatusPending,
			Priority:    models2.PriorityMedium,
		}

		if req.TodoExecutorID == nil {
			todo.ExecutorID = 1 // 默认执行人
		}

		if err := s.todoRepo.Create(todo); err != nil {
			return nil, fmt.Errorf("创建待办事项失败: %w", err)
		}
	}

	return activity, nil
}

// GetActivitiesByCustomer 获取客户的跟进记录列表
func (s *ActivityService) GetActivitiesByCustomer(customerID uint64, page, pageSize int) ([]*models2.ActivityResponse, int64, error) {
	responses, total, err := s.activityRepo.GetByCustomerID(customerID, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	// 转换为指针切片
	resultList := make([]*models2.ActivityResponse, len(responses))
	for i := range responses {
		resultList[i] = &responses[i]
	}

	return resultList, total, nil
}

// GetActivityByID 根据ID获取跟进记录
func (s *ActivityService) GetActivityByID(id uint64) (*models2.ActivityResponse, error) {
	activity, err := s.activityRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	data := activity.GetDataStruct()
	response := &models2.ActivityResponse{
		Activity:     *activity,
		UserName:     activity.User.Name,
		CustomerName: activity.Customer.Name,
		Content:      data.Content,
		Result:       data.Result,
		Amount:       data.Amount,
		Cost:         data.Cost,
		Feedback:     data.Feedback,
		Satisfaction: data.Satisfaction,
		TimeAgo:      activity.GetTimeAgo(),
	}

	return response, nil
}

// UpdateActivity 更新跟进记录
func (s *ActivityService) UpdateActivity(id uint64, req *models2.ActivityUpdateRequest) error {
	activity, err := s.activityRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("查询跟进记录失败: %w", err)
	}

	// 更新基础字段
	if req.Title != nil {
		activity.Title = *req.Title
	}
	if req.Remark != nil {
		activity.Remark = *req.Remark
	}
	if req.Duration != nil {
		activity.Duration = req.Duration
	}
	if req.Location != nil {
		activity.Location = *req.Location
	}
	if req.NextFollowTime != nil {
		activity.NextFollowTime = req.NextFollowTime
	}
	if req.Attachments != nil {
		activity.Attachments = req.Attachments
	}

	// 更新Data字段中的数据
	data := activity.GetDataStruct()
	needUpdateData := false

	if req.Content != nil {
		data.Content = *req.Content
		needUpdateData = true
	}
	if req.Result != nil {
		data.Result = *req.Result
		needUpdateData = true
	}
	if req.Amount != nil {
		data.Amount = *req.Amount
		needUpdateData = true
	}
	if req.Cost != nil {
		data.Cost = *req.Cost
		needUpdateData = true
	}
	if req.Feedback != nil {
		data.Feedback = *req.Feedback
		needUpdateData = true
	}
	if req.Satisfaction != nil {
		data.Satisfaction = *req.Satisfaction
		needUpdateData = true
	}

	if needUpdateData {
		if err := activity.SetDataStruct(data); err != nil {
			return fmt.Errorf("设置数据结构失败: %w", err)
		}
	}

	// 执行更新
	if err := s.activityRepo.Update(activity); err != nil {
		return fmt.Errorf("更新跟进记录失败: %w", err)
	}

	return nil
}

// UpdateActivityFeedback 更新跟进记录反馈
func (s *ActivityService) UpdateActivityFeedback(id uint64, feedback string, satisfaction int) error {
	activity, err := s.activityRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("查询跟进记录失败: %w", err)
	}

	// 更新Data字段中的反馈信息
	data := activity.GetDataStruct()
	data.Feedback = feedback
	if satisfaction >= 1 && satisfaction <= 5 {
		data.Satisfaction = satisfaction
	}

	if err := activity.SetDataStruct(data); err != nil {
		return fmt.Errorf("设置数据结构失败: %w", err)
	}

	// 执行更新
	if err := s.activityRepo.Update(activity); err != nil {
		return fmt.Errorf("更新反馈失败: %w", err)
	}

	return nil
}

// DeleteActivity 删除跟进记录（软删除）
func (s *ActivityService) DeleteActivity(id uint64) error {
	return s.activityRepo.Delete(id)
}

// GetActivityStatistics 获取跟进记录统计信息
func (s *ActivityService) GetActivityStatistics(customerID uint64) (map[string]interface{}, error) {
	return s.activityRepo.GetStatistics(customerID)
}

// GetActivitiesNeedingFollowUp 获取需要跟进的记录
func (s *ActivityService) GetActivitiesNeedingFollowUp() ([]*models2.ActivityResponse, error) {
	responses, err := s.activityRepo.GetNeedingFollowUp()
	if err != nil {
		return nil, err
	}

	// 转换为指针切片
	resultList := make([]*models2.ActivityResponse, len(responses))
	for i := range responses {
		resultList[i] = &responses[i]
	}

	return resultList, nil
}
