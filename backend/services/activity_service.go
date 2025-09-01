package services

import (
	"errors"
	"fmt"
	"time"
	"unicode/utf8"

	"gorm.io/gorm"

	"crm/models"
)

type ActivityService struct {
	db *gorm.DB
}

func NewActivityService(db *gorm.DB) *ActivityService {
	return &ActivityService{db: db}
}

// CreateActivity 创建跟进记录
func (s *ActivityService) CreateActivity(req *models.ActivityCreateRequest) (*models.Activity, error) {
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
	data := models.ActivityData{
		Content:      req.Content,
		Result:       req.Result,
		Amount:       req.Amount,
		Cost:         req.Cost,
		Feedback:     req.Feedback,
		Satisfaction: req.Satisfaction,
	}

	activity := &models.Activity{
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

	// 开始事务
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 创建活动记录
	if err := tx.Create(activity).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("创建跟进记录失败: %w", err)
	}

	// 如果需要创建待办事项
	if req.CreateTodo && req.TodoPlannedTime != nil {
		todo := &models.Todo{
			CustomerID:  req.CustomerID,
			CreatorID:   1, // TODO: 从上下文获取当前用户ID
			ExecutorID:  *req.TodoExecutorID,
			Title:       req.TodoContent,
			Content:     req.Content,
			PlannedTime: *req.TodoPlannedTime,
			Status:      models.TodoStatusPending,
			Priority:    models.PriorityMedium,
		}

		if req.TodoExecutorID == nil {
			todo.ExecutorID = 1 // 默认执行人
		}

		if err := tx.Create(todo).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("创建待办事项失败: %w", err)
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("事务提交失败: %w", err)
	}

	return activity, nil
}

// GetActivitiesByCustomer 获取客户的跟进记录列表
func (s *ActivityService) GetActivitiesByCustomer(customerID uint64, page, pageSize int) ([]*models.ActivityResponse, int64, error) {
	var activities []models.Activity
	var total int64

	// 构建查询
	query := s.db.Model(&models.Activity{}).
		Where("customer_id = ? AND is_deleted = ?", customerID, false)

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("统计记录数失败: %w", err)
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Preload("User").Preload("Customer").
		Offset(offset).Limit(pageSize).
		Order("created_at DESC").
		Find(&activities).Error; err != nil {
		return nil, 0, fmt.Errorf("查询跟进记录失败: %w", err)
	}

	// 转换为响应格式
	responses := make([]*models.ActivityResponse, len(activities))
	for i, activity := range activities {
		data := activity.GetDataStruct()
		responses[i] = &models.ActivityResponse{
			Activity:     activity,
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
	}

	return responses, total, nil
}

// GetActivityByID 根据ID获取跟进记录
func (s *ActivityService) GetActivityByID(id uint64) (*models.ActivityResponse, error) {
	var activity models.Activity

	if err := s.db.Preload("User").Preload("Customer").
		Where("id = ? AND is_deleted = ?", id, false).
		First(&activity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("跟进记录不存在")
		}
		return nil, fmt.Errorf("查询跟进记录失败: %w", err)
	}

	data := activity.GetDataStruct()
	response := &models.ActivityResponse{
		Activity:     activity,
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
func (s *ActivityService) UpdateActivity(id uint64, req *models.ActivityUpdateRequest) error {
	var activity models.Activity

	if err := s.db.Where("id = ? AND is_deleted = ?", id, false).
		First(&activity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("跟进记录不存在")
		}
		return fmt.Errorf("查询跟进记录失败: %w", err)
	}

	// 更新基础字段
	updates := make(map[string]interface{})
	if req.Title != nil {
		updates["title"] = *req.Title
	}
	if req.Remark != nil {
		updates["remark"] = *req.Remark
	}
	if req.Duration != nil {
		updates["duration"] = *req.Duration
	}
	if req.Location != nil {
		updates["location"] = *req.Location
	}
	if req.NextFollowTime != nil {
		updates["next_follow_time"] = *req.NextFollowTime
	}
	if req.Attachments != nil {
		updates["attachments"] = req.Attachments
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
		updates["data"] = activity.Data
	}

	// 执行更新
	if len(updates) > 0 {
		if err := s.db.Model(&activity).Updates(updates).Error; err != nil {
			return fmt.Errorf("更新跟进记录失败: %w", err)
		}
	}

	return nil
}

// UpdateActivityFeedback 更新跟进记录反馈
func (s *ActivityService) UpdateActivityFeedback(id uint64, feedback string, satisfaction int) error {
	var activity models.Activity

	if err := s.db.Where("id = ? AND is_deleted = ?", id, false).
		First(&activity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("跟进记录不存在")
		}
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
	if err := s.db.Model(&activity).Update("data", activity.Data).Error; err != nil {
		return fmt.Errorf("更新反馈失败: %w", err)
	}

	return nil
}

// DeleteActivity 删除跟进记录（软删除）
func (s *ActivityService) DeleteActivity(id uint64) error {
	result := s.db.Model(&models.Activity{}).
		Where("id = ? AND is_deleted = ?", id, false).
		Updates(map[string]interface{}{
			"is_deleted": true,
			"deleted_at": time.Now(),
		})

	if result.Error != nil {
		return fmt.Errorf("删除跟进记录失败: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return errors.New("跟进记录不存在")
	}

	return nil
}

// GetActivityStatistics 获取跟进记录统计信息
func (s *ActivityService) GetActivityStatistics(customerID uint64) (map[string]interface{}, error) {
	var results []struct {
		Kind   string  `json:"kind"`
		Count  int     `json:"count"`
		Amount float64 `json:"amount"`
		Cost   float64 `json:"cost"`
	}

	// 统计各类型记录数量和金额
	if err := s.db.Raw(`
		SELECT 
			kind,
			COUNT(*) as count,
			COALESCE(SUM(CAST(data->>'amount' AS DECIMAL)), 0) as amount,
			COALESCE(SUM(CAST(data->>'cost' AS DECIMAL)), 0) as cost
		FROM activities 
		WHERE customer_id = ? AND is_deleted = false 
		GROUP BY kind
	`, customerID).Scan(&results).Error; err != nil {
		return nil, fmt.Errorf("统计跟进记录失败: %w", err)
	}

	// 构建统计结果
	stats := map[string]interface{}{
		"total_records":   0,
		"records_by_kind": make(map[string]int),
		"total_amount":    0.0,
		"total_cost":      0.0,
		"order_count":     0,
		"sample_count":    0,
	}

	totalRecords := 0
	totalAmount := 0.0
	totalCost := 0.0

	for _, result := range results {
		totalRecords += result.Count
		totalAmount += result.Amount
		totalCost += result.Cost

		stats["records_by_kind"].(map[string]int)[result.Kind] = result.Count

		if result.Kind == "order" {
			stats["order_count"] = result.Count
		} else if result.Kind == "sample" {
			stats["sample_count"] = result.Count
		}
	}

	stats["total_records"] = totalRecords
	stats["total_amount"] = totalAmount
	stats["total_cost"] = totalCost

	return stats, nil
}

// GetActivitiesNeedingFollowUp 获取需要跟进的记录
func (s *ActivityService) GetActivitiesNeedingFollowUp() ([]*models.ActivityResponse, error) {
	var activities []models.Activity

	if err := s.db.Preload("User").Preload("Customer").
		Where("next_follow_time <= ? AND is_deleted = ?", time.Now(), false).
		Order("next_follow_time ASC").
		Find(&activities).Error; err != nil {
		return nil, fmt.Errorf("查询需要跟进的记录失败: %w", err)
	}

	// 转换为响应格式
	responses := make([]*models.ActivityResponse, len(activities))
	for i, activity := range activities {
		data := activity.GetDataStruct()
		responses[i] = &models.ActivityResponse{
			Activity:     activity,
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
	}

	return responses, nil
}
