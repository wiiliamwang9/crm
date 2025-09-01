package repository

import (
	"crm/internal/domain"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type activityRepository struct {
	db *gorm.DB
}

// NewActivityRepository 创建活动仓储实例
func NewActivityRepository(db *gorm.DB) ActivityRepository {
	return &activityRepository{db: db}
}

// GetByID 根据ID获取活动
func (r *activityRepository) GetByID(id uint64) (*domain.Activity, error) {
	var activity domain.Activity
	err := r.db.Preload("User").Preload("Customer").
		Where("id = ? AND is_deleted = ?", id, false).
		First(&activity).Error
	return &activity, err
}

// GetList 获取活动列表
func (r *activityRepository) GetList(req domain.ActivityQueryRequest) ([]domain.ActivityResponse, int64, error) {
	var activities []domain.Activity
	query := r.db.Model(&domain.Activity{}).Where("is_deleted = ?", false)

	// 添加筛选条件
	if req.CustomerID != nil {
		query = query.Where("customer_id = ?", *req.CustomerID)
	}
	if req.UserID != nil {
		query = query.Where("user_id = ?", *req.UserID)
	}
	if req.Kind != nil {
		query = query.Where("kind = ?", *req.Kind)
	}

	// 时间筛选
	if req.StartDate != nil {
		query = query.Where("created_at >= ?", *req.StartDate)
	}
	if req.EndDate != nil {
		query = query.Where("created_at <= ?", *req.EndDate)
	}

	// 关键词搜索
	if req.Keyword != "" {
		query = query.Where("title LIKE ? OR remark LIKE ?", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}

	// 获取总数
	var total int64
	query.Count(&total)

	// 分页查询
	offset := (req.Page - 1) * req.PageSize
	err := query.Preload("User").Preload("Customer").
		Order("created_at DESC").Offset(offset).Limit(req.PageSize).Find(&activities).Error

	if err != nil {
		return nil, 0, err
	}

	// 转换为响应格式
	responses := make([]domain.ActivityResponse, len(activities))
	for i, activity := range activities {
		data := activity.GetDataStruct()
		responses[i] = domain.ActivityResponse{
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

// GetByCustomerID 根据客户ID获取活动列表
func (r *activityRepository) GetByCustomerID(customerID uint64, page, pageSize int) ([]domain.ActivityResponse, int64, error) {
	var activities []domain.Activity
	var total int64

	// 构建查询
	query := r.db.Model(&domain.Activity{}).
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
	responses := make([]domain.ActivityResponse, len(activities))
	for i, activity := range activities {
		data := activity.GetDataStruct()
		responses[i] = domain.ActivityResponse{
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

// Create 创建活动
func (r *activityRepository) Create(activity *domain.Activity) error {
	// 设置默认值
	if activity.Kind == "" {
		activity.Kind = domain.ActivityKindOther
	}
	return r.db.Create(activity).Error
}

// Update 更新活动
func (r *activityRepository) Update(activity *domain.Activity) error {
	return r.db.Save(activity).Error
}

// Delete 删除活动（软删除）
func (r *activityRepository) Delete(id uint64) error {
	result := r.db.Model(&domain.Activity{}).
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

// GetStatistics 获取跟进记录统计信息
func (r *activityRepository) GetStatistics(customerID uint64) (map[string]interface{}, error) {
	var results []struct {
		Kind   string  `json:"kind"`
		Count  int     `json:"count"`
		Amount float64 `json:"amount"`
		Cost   float64 `json:"cost"`
	}

	// 统计各类型记录数量和金额
	if err := r.db.Raw(`
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

// GetNeedingFollowUp 获取需要跟进的记录
func (r *activityRepository) GetNeedingFollowUp() ([]domain.ActivityResponse, error) {
	var activities []domain.Activity

	if err := r.db.Preload("User").Preload("Customer").
		Where("next_follow_time <= ? AND is_deleted = ?", time.Now(), false).
		Order("next_follow_time ASC").
		Find(&activities).Error; err != nil {
		return nil, fmt.Errorf("查询需要跟进的记录失败: %w", err)
	}

	// 转换为响应格式
	responses := make([]domain.ActivityResponse, len(activities))
	for i, activity := range activities {
		data := activity.GetDataStruct()
		responses[i] = domain.ActivityResponse{
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
