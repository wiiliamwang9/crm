package services

import (
	"bytes"
	models2 "crm/internal/domain"
	"crm/internal/repository"
	"fmt"
	"text/template"
	"time"
)

// ReminderService 提醒服务
type ReminderService struct {
	repo repository.ReminderRepository
}

// NewReminderService 创建提醒服务
func NewReminderService(repo repository.ReminderRepository) *ReminderService {
	return &ReminderService{repo: repo}
}

// CreateReminder 创建提醒
func (s *ReminderService) CreateReminder(req models2.ReminderCreateRequest) (*models2.Reminder, error) {
	reminder := models2.Reminder{
		TodoID:       req.TodoID,
		UserID:       req.UserID,
		Type:         req.Type,
		Title:        req.Title,
		Content:      req.Content,
		Frequency:    req.Frequency,
		ScheduleTime: req.ScheduleTime,
		MaxRetries:   req.MaxRetries,
	}

	if err := s.repo.Create(&reminder); err != nil {
		return nil, err
	}

	return &reminder, nil
}

// GetReminders 获取提醒列表
func (s *ReminderService) GetReminders(req models2.ReminderQueryRequest) ([]models2.ReminderResponse, int64, error) {
	return s.repo.GetList(req)
}

// GetReminderByID 根据ID获取提醒
func (s *ReminderService) GetReminderByID(id uint64) (*models2.ReminderResponse, error) {
	reminder, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	response := models2.ReminderResponse{
		Reminder:     *reminder,
		TodoTitle:    reminder.Todo.Title,
		UserName:     reminder.User.Name,
		CustomerName: reminder.Todo.Customer.Name,
	}

	return &response, nil
}

// UpdateReminder 更新提醒
func (s *ReminderService) UpdateReminder(id uint64, req models2.ReminderUpdateRequest) (*models2.Reminder, error) {
	reminder, err := s.repo.GetByID(id)
	if err != nil {
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

	if err := s.repo.Update(reminder); err != nil {
		return nil, err
	}

	return reminder, nil
}

// DeleteReminder 删除提醒
func (s *ReminderService) DeleteReminder(id uint64) error {
	return s.repo.Delete(id)
}

// CancelReminder 取消提醒
func (s *ReminderService) CancelReminder(id uint64) (*models2.Reminder, error) {
	reminder, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	reminder.Status = models2.ReminderStatusCancelled
	if err := s.repo.Update(reminder); err != nil {
		return nil, err
	}

	return reminder, nil
}

// GetPendingReminders 获取待发送的提醒
func (s *ReminderService) GetPendingReminders(limit int) ([]models2.Reminder, error) {
	return s.repo.GetPendingReminders()
}

// MarkReminderSent 标记提醒已发送
func (s *ReminderService) MarkReminderSent(id uint64) error {
	return s.repo.UpdateStatus(id, models2.ReminderStatusSent)
}

// MarkReminderFailed 标记提醒发送失败
func (s *ReminderService) MarkReminderFailed(id uint64, reason string) error {
	// TODO: 添加MarkFailed方法到repository接口
	return s.repo.UpdateStatus(id, models2.ReminderStatusFailed)
}

// CreateReminderFromTodo 从待办创建提醒
func (s *ReminderService) CreateReminderFromTodo(todo *models2.Todo, userID uint64, reminderTime time.Time, reminderType models2.ReminderType) (*models2.Reminder, error) {
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

	req := models2.ReminderCreateRequest{
		TodoID:       todo.ID,
		UserID:       userID,
		Type:         reminderType,
		Title:        title,
		Content:      content,
		Frequency:    models2.ReminderFrequencyOnce,
		ScheduleTime: reminderTime,
		MaxRetries:   3,
	}

	return s.CreateReminder(req)
}

// GetDefaultTemplate 获取默认模板
func (s *ReminderService) GetDefaultTemplate(reminderType models2.ReminderType) (*models2.ReminderTemplate, error) {
	return s.repo.GetReminderTemplate(reminderType, true)
}

// RenderTemplate 渲染模板
func (s *ReminderService) RenderTemplate(tmpl *models2.ReminderTemplate, todo *models2.Todo) (string, string, error) {
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
func (s *ReminderService) GetUserReminderConfig(userID uint64) (*models2.ReminderConfig, error) {
	return s.repo.GetReminderConfig(userID)
}

// UpdateUserReminderConfig 更新用户提醒配置
func (s *ReminderService) UpdateUserReminderConfig(userID uint64, config *models2.ReminderConfig) error {
	config.UserID = userID
	return s.repo.UpdateReminderConfig(config)
}

// SendReminder 发送提醒（这里是模拟实现）
func (s *ReminderService) SendReminder(reminder *models2.Reminder) error {
	fmt.Printf("发送提醒 - ID: %d, 类型: %s, 标题: %s\n",
		reminder.ID, reminder.Type, reminder.Title)

	// 在实际实现中，这里应该调用具体的推送服务
	switch reminder.Type {
	case models2.ReminderTypeWechat:
		return s.sendWechatReminder(reminder)
	case models2.ReminderTypeEnterpriseWechat:
		return s.sendEnterpriseWechatReminder(reminder)
	case models2.ReminderTypeBoth:
		if err := s.sendWechatReminder(reminder); err != nil {
			return err
		}
		return s.sendEnterpriseWechatReminder(reminder)
	}

	return nil
}

// sendWechatReminder 发送微信提醒
func (s *ReminderService) sendWechatReminder(reminder *models2.Reminder) error {
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
func (s *ReminderService) sendEnterpriseWechatReminder(reminder *models2.Reminder) error {
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
				nextReminder := models2.ReminderCreateRequest{
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
	return s.repo.GetStats(userID)
}
