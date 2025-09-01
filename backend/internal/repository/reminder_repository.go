package repository

import (
	"crm/internal/domain"
	"time"

	"gorm.io/gorm"
)

type reminderRepository struct {
	db *gorm.DB
}

// NewReminderRepository 创建提醒仓储实例
func NewReminderRepository(db *gorm.DB) ReminderRepository {
	return &reminderRepository{db: db}
}

// GetByID 根据ID获取提醒
func (r *reminderRepository) GetByID(id uint64) (*domain.Reminder, error) {
	var reminder domain.Reminder
	err := r.db.Preload("Todo").Preload("Todo.Customer").Preload("User").
		First(&reminder, id).Error
	return &reminder, err
}

// GetList 获取提醒列表
func (r *reminderRepository) GetList(req domain.ReminderQueryRequest) ([]domain.ReminderResponse, int64, error) {
	// 构建查询
	query := r.db.Model(&domain.Reminder{})

	// 添加筛选条件
	if req.TodoID != nil {
		query = query.Where("todo_id = ?", *req.TodoID)
	}
	if req.UserID != nil {
		query = query.Where("user_id = ?", *req.UserID)
	}
	if req.Status != nil {
		query = query.Where("status = ?", *req.Status)
	}
	if req.Type != nil {
		query = query.Where("type = ?", *req.Type)
	}
	if req.StartDate != nil {
		query = query.Where("schedule_time >= ?", *req.StartDate)
	}
	if req.EndDate != nil {
		query = query.Where("schedule_time <= ?", *req.EndDate)
	}

	// 获取总数
	var total int64
	query.Count(&total)

	// 分页查询
	var reminders []domain.Reminder
	offset := (req.Page - 1) * req.PageSize
	err := query.Preload("Todo").Preload("Todo.Customer").Preload("User").
		Order("schedule_time DESC").Offset(offset).Limit(req.PageSize).Find(&reminders).Error

	if err != nil {
		return nil, 0, err
	}

	// 转换为响应格式
	var responses []domain.ReminderResponse
	for _, reminder := range reminders {
		response := domain.ReminderResponse{
			Reminder:     reminder,
			TodoTitle:    reminder.Todo.Title,
			UserName:     reminder.User.Name,
			CustomerName: reminder.Todo.Customer.Name,
		}
		responses = append(responses, response)
	}

	return responses, total, nil
}

// Create 创建提醒
func (r *reminderRepository) Create(reminder *domain.Reminder) error {
	if reminder.MaxRetries == 0 {
		reminder.MaxRetries = 3 // 默认最大重试3次
	}
	return r.db.Create(reminder).Error
}

// Update 更新提醒
func (r *reminderRepository) Update(reminder *domain.Reminder) error {
	return r.db.Save(reminder).Error
}

// Delete 删除提醒
func (r *reminderRepository) Delete(id uint64) error {
	return r.db.Delete(&domain.Reminder{}, id).Error
}

// GetPendingReminders 获取待发送的提醒
func (r *reminderRepository) GetPendingReminders() ([]domain.Reminder, error) {
	var reminders []domain.Reminder
	now := time.Now()

	err := r.db.Preload("Todo").Preload("Todo.Customer").Preload("User").
		Where("status = ? AND schedule_time <= ?", domain.ReminderStatusPending, now).
		Order("schedule_time ASC").
		Limit(100).
		Find(&reminders).Error

	return reminders, err
}

// UpdateStatus 更新提醒状态
func (r *reminderRepository) UpdateStatus(id uint64, status domain.ReminderStatus) error {
	updates := map[string]interface{}{
		"status": status,
	}

	if status == domain.ReminderStatusSent {
		now := time.Now()
		updates["sent_time"] = &now
	}

	result := r.db.Model(&domain.Reminder{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// MarkSent 标记提醒已发送
func (r *reminderRepository) MarkSent(id uint64) error {
	now := time.Now()
	return r.db.Model(&domain.Reminder{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":    domain.ReminderStatusSent,
			"sent_time": &now,
		}).Error
}

// MarkFailed 标记提醒发送失败
func (r *reminderRepository) MarkFailed(id uint64, reason string) error {
	return r.db.Model(&domain.Reminder{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":      domain.ReminderStatusFailed,
			"fail_reason": reason,
			"retry_count": gorm.Expr("retry_count + 1"),
		}).Error
}

// GetStats 获取提醒统计
func (r *reminderRepository) GetStats(userID *uint64) (map[string]int64, error) {
	query := r.db.Model(&domain.Reminder{})

	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}

	stats := make(map[string]int64)

	// 统计各状态数量
	statuses := []domain.ReminderStatus{
		domain.ReminderStatusPending,
		domain.ReminderStatusSent,
		domain.ReminderStatusFailed,
		domain.ReminderStatusCancelled,
	}

	for _, status := range statuses {
		var count int64
		query.Where("status = ?", status).Count(&count)
		stats[string(status)] = count
	}

	// 统计今日待发送
	var todayCount int64
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)
	query.Where("status = ? AND schedule_time >= ? AND schedule_time < ?",
		domain.ReminderStatusPending, startOfDay, endOfDay).Count(&todayCount)
	stats["today_pending"] = todayCount

	return stats, nil
}

// GetReminderConfig 获取用户提醒配置
func (r *reminderRepository) GetReminderConfig(userID uint64) (*domain.ReminderConfig, error) {
	var config domain.ReminderConfig
	err := r.db.Where("user_id = ?", userID).First(&config).Error

	if err == gorm.ErrRecordNotFound {
		// 如果没有配置，创建默认配置
		config = domain.ReminderConfig{
			UserID:                 userID,
			EnableWechat:           true,
			EnableEnterpriseWechat: true,
			DefaultAdvanceMinutes:  30,
			QuietStartTime:         "22:00",
			QuietEndTime:           "08:00",
		}

		if err := r.db.Create(&config).Error; err != nil {
			return nil, err
		}
	}

	return &config, nil
}

// UpdateReminderConfig 更新用户提醒配置
func (r *reminderRepository) UpdateReminderConfig(config *domain.ReminderConfig) error {
	return r.db.Save(config).Error
}

// GetReminderTemplate 获取提醒模板
func (r *reminderRepository) GetReminderTemplate(reminderType domain.ReminderType, isDefault bool) (*domain.ReminderTemplate, error) {
	var template domain.ReminderTemplate
	err := r.db.Where("type = ? AND is_default = ? AND is_active = ?",
		reminderType, isDefault, true).First(&template).Error

	if err == gorm.ErrRecordNotFound && isDefault {
		// 如果没有找到默认模板，返回一个基本模板
		return &domain.ReminderTemplate{
			Title:   "待办提醒：{{.Title}}",
			Content: "您有一个待办事项需要处理：\n\n标题：{{.Title}}\n内容：{{.Content}}\n客户：{{.CustomerName}}\n计划时间：{{.PlannedTime}}\n\n请及时处理！",
		}, nil
	}

	return &template, err
}
