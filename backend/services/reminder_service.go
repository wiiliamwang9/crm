package services

import (
	"bytes"
	"fmt"
	"text/template"
	"time"

	"gorm.io/gorm"

	"crm/models"
)

// ReminderService 提醒服务
type ReminderService struct {
	db *gorm.DB
}

// NewReminderService 创建提醒服务
func NewReminderService(db *gorm.DB) *ReminderService {
	return &ReminderService{db: db}
}

// CreateReminder 创建提醒
func (s *ReminderService) CreateReminder(req models.ReminderCreateRequest) (*models.Reminder, error) {
	reminder := models.Reminder{
		TodoID:       req.TodoID,
		UserID:       req.UserID,
		Type:         req.Type,
		Title:        req.Title,
		Content:      req.Content,
		Frequency:    req.Frequency,
		ScheduleTime: req.ScheduleTime,
		MaxRetries:   req.MaxRetries,
	}

	if reminder.MaxRetries == 0 {
		reminder.MaxRetries = 3 // 默认最大重试3次
	}

	if err := s.db.Create(&reminder).Error; err != nil {
		return nil, err
	}

	return &reminder, nil
}

// GetReminders 获取提醒列表
func (s *ReminderService) GetReminders(req models.ReminderQueryRequest) ([]models.ReminderResponse, int64, error) {
	// 构建查询
	query := s.db.Model(&models.Reminder{})

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
	var reminders []models.Reminder
	offset := (req.Page - 1) * req.PageSize
	err := query.Preload("Todo").Preload("Todo.Customer").Preload("User").
		Order("schedule_time DESC").Offset(offset).Limit(req.PageSize).Find(&reminders).Error

	if err != nil {
		return nil, 0, err
	}

	// 转换为响应格式
	var responses []models.ReminderResponse
	for _, reminder := range reminders {
		response := models.ReminderResponse{
			Reminder: reminder,
			TodoTitle: reminder.Todo.Title,
			UserName:  reminder.User.Name,
			CustomerName: reminder.Todo.Customer.Name,
		}
		responses = append(responses, response)
	}

	return responses, total, nil
}

// GetReminderByID 根据ID获取提醒
func (s *ReminderService) GetReminderByID(id uint64) (*models.ReminderResponse, error) {
	var reminder models.Reminder
	err := s.db.Preload("Todo").Preload("Todo.Customer").Preload("User").
		First(&reminder, id).Error

	if err != nil {
		return nil, err
	}

	response := models.ReminderResponse{
		Reminder:     reminder,
		TodoTitle:    reminder.Todo.Title,
		UserName:     reminder.User.Name,
		CustomerName: reminder.Todo.Customer.Name,
	}

	return &response, nil
}

// UpdateReminder 更新提醒
func (s *ReminderService) UpdateReminder(id uint64, req models.ReminderUpdateRequest) (*models.Reminder, error) {
	var reminder models.Reminder
	if err := s.db.First(&reminder, id).Error; err != nil {
		return nil, err
	}

	// 更新字段
	if req.Title != nil {
		reminder.Title = *req.Title
	}
	if req.Content != nil {
		reminder.Content = *req.Content
	}
	if req.Status != nil {
		reminder.Status = *req.Status
	}
	if req.Frequency != nil {
		reminder.Frequency = *req.Frequency
	}
	if req.ScheduleTime != nil {
		reminder.ScheduleTime = *req.ScheduleTime
	}
	if req.MaxRetries != nil {
		reminder.MaxRetries = *req.MaxRetries
	}

	if err := s.db.Save(&reminder).Error; err != nil {
		return nil, err
	}

	return &reminder, nil
}

// DeleteReminder 删除提醒
func (s *ReminderService) DeleteReminder(id uint64) error {
	return s.db.Delete(&models.Reminder{}, id).Error
}

// CancelReminder 取消提醒
func (s *ReminderService) CancelReminder(id uint64) (*models.Reminder, error) {
	var reminder models.Reminder
	if err := s.db.First(&reminder, id).Error; err != nil {
		return nil, err
	}

	reminder.Status = models.ReminderStatusCancelled
	if err := s.db.Save(&reminder).Error; err != nil {
		return nil, err
	}

	return &reminder, nil
}

// GetPendingReminders 获取待发送的提醒
func (s *ReminderService) GetPendingReminders(limit int) ([]models.Reminder, error) {
	var reminders []models.Reminder
	now := time.Now()

	err := s.db.Preload("Todo").Preload("Todo.Customer").Preload("User").
		Where("status = ? AND schedule_time <= ?", models.ReminderStatusPending, now).
		Order("schedule_time ASC").
		Limit(limit).
		Find(&reminders).Error

	return reminders, err
}

// MarkReminderSent 标记提醒已发送
func (s *ReminderService) MarkReminderSent(id uint64) error {
	now := time.Now()
	return s.db.Model(&models.Reminder{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":    models.ReminderStatusSent,
			"sent_time": &now,
		}).Error
}

// MarkReminderFailed 标记提醒发送失败
func (s *ReminderService) MarkReminderFailed(id uint64, reason string) error {
	return s.db.Model(&models.Reminder{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":      models.ReminderStatusFailed,
			"fail_reason": reason,
			"retry_count": gorm.Expr("retry_count + 1"),
		}).Error
}

// CreateReminderFromTodo 从待办创建提醒
func (s *ReminderService) CreateReminderFromTodo(todo *models.Todo, userID uint64, reminderTime time.Time, reminderType models.ReminderType) (*models.Reminder, error) {
	// 获取默认模板
	template, err := s.GetDefaultTemplate(reminderType)
	if err != nil {
		return nil, err
	}

	// 渲染模板
	title, content, err := s.RenderTemplate(template, todo)
	if err != nil {
		return nil, err
	}

	req := models.ReminderCreateRequest{
		TodoID:       todo.ID,
		UserID:       userID,
		Type:         reminderType,
		Title:        title,
		Content:      content,
		Frequency:    models.ReminderFrequencyOnce,
		ScheduleTime: reminderTime,
		MaxRetries:   3,
	}

	return s.CreateReminder(req)
}

// GetDefaultTemplate 获取默认模板
func (s *ReminderService) GetDefaultTemplate(reminderType models.ReminderType) (*models.ReminderTemplate, error) {
	var template models.ReminderTemplate
	err := s.db.Where("type = ? AND is_default = ? AND is_active = ?", 
		reminderType, true, true).First(&template).Error
	
	if err == gorm.ErrRecordNotFound {
		// 如果没有找到默认模板，返回一个基本模板
		return &models.ReminderTemplate{
			Title:   "待办提醒：{{.Title}}",
			Content: "您有一个待办事项需要处理：\n\n标题：{{.Title}}\n内容：{{.Content}}\n客户：{{.CustomerName}}\n计划时间：{{.PlannedTime}}\n\n请及时处理！",
		}, nil
	}
	
	return &template, err
}

// RenderTemplate 渲染模板
func (s *ReminderService) RenderTemplate(tmpl *models.ReminderTemplate, todo *models.Todo) (string, string, error) {
	// 准备模板变量
	data := map[string]interface{}{
		"Title":        todo.Title,
		"Content":      todo.Content,
		"CustomerName": todo.Customer.Name,
		"PlannedTime":  todo.PlannedTime.Format("2006-01-02 15:04"),
		"Priority":     string(todo.Priority),
		"Status":       string(todo.Status),
	}

	// 渲染标题
	titleTmpl, err := template.New("title").Parse(tmpl.Title)
	if err != nil {
		return "", "", err
	}
	
	var titleBuf bytes.Buffer
	if err := titleTmpl.Execute(&titleBuf, data); err != nil {
		return "", "", err
	}

	// 渲染内容
	contentTmpl, err := template.New("content").Parse(tmpl.Content)
	if err != nil {
		return "", "", err
	}
	
	var contentBuf bytes.Buffer
	if err := contentTmpl.Execute(&contentBuf, data); err != nil {
		return "", "", err
	}

	return titleBuf.String(), contentBuf.String(), nil
}

// GetUserReminderConfig 获取用户提醒配置
func (s *ReminderService) GetUserReminderConfig(userID uint64) (*models.ReminderConfig, error) {
	var config models.ReminderConfig
	err := s.db.Where("user_id = ?", userID).First(&config).Error
	
	if err == gorm.ErrRecordNotFound {
		// 如果没有配置，创建默认配置
		config = models.ReminderConfig{
			UserID:                userID,
			EnableWechat:          true,
			EnableEnterpriseWechat: true,
			DefaultAdvanceMinutes: 30,
			QuietStartTime:       "22:00",
			QuietEndTime:         "08:00",
		}
		
		if err := s.db.Create(&config).Error; err != nil {
			return nil, err
		}
	}
	
	return &config, nil
}

// UpdateUserReminderConfig 更新用户提醒配置
func (s *ReminderService) UpdateUserReminderConfig(userID uint64, config *models.ReminderConfig) error {
	config.UserID = userID
	return s.db.Save(config).Error
}

// SendReminder 发送提醒（这里是模拟实现）
func (s *ReminderService) SendReminder(reminder *models.Reminder) error {
	fmt.Printf("发送提醒 - ID: %d, 类型: %s, 标题: %s\n", 
		reminder.ID, reminder.Type, reminder.Title)
	
	// 在实际实现中，这里应该调用具体的推送服务
	switch reminder.Type {
	case models.ReminderTypeWechat:
		return s.sendWechatReminder(reminder)
	case models.ReminderTypeEnterpriseWechat:
		return s.sendEnterpriseWechatReminder(reminder)
	case models.ReminderTypeBoth:
		if err := s.sendWechatReminder(reminder); err != nil {
			return err
		}
		return s.sendEnterpriseWechatReminder(reminder)
	}
	
	return nil
}

// sendWechatReminder 发送微信提醒
func (s *ReminderService) sendWechatReminder(reminder *models.Reminder) error {
	// 获取用户配置
	config, err := s.GetUserReminderConfig(reminder.UserID)
	if err != nil || !config.EnableWechat {
		return fmt.Errorf("用户未启用微信提醒")
	}

	// 检查免打扰时间
	if config.IsInQuietTime(time.Now()) {
		// 延后到免打扰时间结束后
		return fmt.Errorf("当前在免打扰时间内")
	}

	// 这里应该调用微信API发送消息
	fmt.Printf("发送微信提醒给用户 %s: %s\n", config.WechatUserID, reminder.Title)
	
	return nil
}

// sendEnterpriseWechatReminder 发送企业微信提醒
func (s *ReminderService) sendEnterpriseWechatReminder(reminder *models.Reminder) error {
	// 获取用户配置
	config, err := s.GetUserReminderConfig(reminder.UserID)
	if err != nil || !config.EnableEnterpriseWechat {
		return fmt.Errorf("用户未启用企业微信提醒")
	}

	// 检查免打扰时间
	if config.IsInQuietTime(time.Now()) {
		// 延后到免打扰时间结束后
		return fmt.Errorf("当前在免打扰时间内")
	}

	// 这里应该调用企业微信API发送消息
	fmt.Printf("发送企业微信提醒给用户 %s: %s\n", config.EnterpriseWechatUserID, reminder.Title)
	
	return nil
}

// ProcessPendingReminders 处理待发送的提醒
func (s *ReminderService) ProcessPendingReminders() error {
	reminders, err := s.GetPendingReminders(100) // 每次处理100条
	if err != nil {
		return err
	}

	for _, reminder := range reminders {
		if err := s.SendReminder(&reminder); err != nil {
			// 标记为失败
			s.MarkReminderFailed(reminder.ID, err.Error())
		} else {
			// 标记为已发送
			s.MarkReminderSent(reminder.ID)
			
			// 如果是重复提醒，创建下次提醒
			if nextTime := reminder.GetNextScheduleTime(); nextTime != nil {
				nextReminder := models.ReminderCreateRequest{
					TodoID:       reminder.TodoID,
					UserID:       reminder.UserID,
					Type:         reminder.Type,
					Title:        reminder.Title,
					Content:      reminder.Content,
					Frequency:    reminder.Frequency,
					ScheduleTime: *nextTime,
					MaxRetries:   reminder.MaxRetries,
				}
				s.CreateReminder(nextReminder)
			}
		}
	}

	return nil
}

// GetReminderStats 获取提醒统计
func (s *ReminderService) GetReminderStats(userID *uint64) (map[string]int64, error) {
	query := s.db.Model(&models.Reminder{})
	
	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}

	stats := make(map[string]int64)
	
	// 统计各状态数量
	statuses := []models.ReminderStatus{
		models.ReminderStatusPending, 
		models.ReminderStatusSent, 
		models.ReminderStatusFailed, 
		models.ReminderStatusCancelled,
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
		models.ReminderStatusPending, startOfDay, endOfDay).Count(&todayCount)
	stats["today_pending"] = todayCount

	return stats, nil
}