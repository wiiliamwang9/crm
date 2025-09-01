package services

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"

	"crm/models"
)

// TodoService 待办服务
type TodoService struct {
	db *gorm.DB
}

// NewTodoService 创建待办服务
func NewTodoService(db *gorm.DB) *TodoService {
	return &TodoService{db: db}
}

// CreateTodo 创建待办
func (s *TodoService) CreateTodo(req models.TodoCreateRequest, creatorID uint64) (*models.Todo, error) {
	todo := models.Todo{
		CustomerID:     req.CustomerID,
		CreatorID:      creatorID,
		ExecutorID:     req.ExecutorID,
		Title:          req.Title,
		Content:        req.Content,
		PlannedTime:    req.PlannedTime,
		IsReminder:     req.IsReminder,
		ReminderType:   req.ReminderType,
		ReminderUserID: req.ReminderUserID,
		ReminderTime:   req.ReminderTime,
		Priority:       req.Priority,
		Tags:           req.Tags,
	}

	if err := s.db.Create(&todo).Error; err != nil {
		return nil, err
	}

	// 记录操作日志
	s.createTodoLog(todo.ID, creatorID, models.ActionCreate, nil, todo)

	return &todo, nil
}

// GetTodos 获取待办列表
func (s *TodoService) GetTodos(req models.TodoQueryRequest) ([]models.TodoResponse, int64, error) {
	// 构建查询
	query := s.db.Model(&models.Todo{}).Where("is_deleted = ?", false)

	// 添加筛选条件
	if req.CustomerID != nil {
		query = query.Where("customer_id = ?", *req.CustomerID)
	}
	if req.ExecutorID != nil {
		query = query.Where("executor_id = ?", *req.ExecutorID)
	}
	if req.CreatorID != nil {
		query = query.Where("creator_id = ?", *req.CreatorID)
	}
	if req.Status != nil {
		query = query.Where("status = ?", *req.Status)
	}
	if req.Priority != nil {
		query = query.Where("priority = ?", *req.Priority)
	}

	// 时间筛选
	if req.DateType != "" {
		now := time.Now()
		switch req.DateType {
		case "yesterday":
			yesterday := now.AddDate(0, 0, -1)
			startOfDay := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 0, 0, 0, 0, yesterday.Location())
			endOfDay := startOfDay.Add(24 * time.Hour)
			query = query.Where("planned_time >= ? AND planned_time < ?", startOfDay, endOfDay)
		case "today":
			startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
			endOfDay := startOfDay.Add(24 * time.Hour)
			query = query.Where("planned_time >= ? AND planned_time < ?", startOfDay, endOfDay)
		case "tomorrow":
			tomorrow := now.AddDate(0, 0, 1)
			startOfDay := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 0, 0, 0, 0, tomorrow.Location())
			endOfDay := startOfDay.Add(24 * time.Hour)
			query = query.Where("planned_time >= ? AND planned_time < ?", startOfDay, endOfDay)
		case "upcoming":
			// 明天及之后的数据
			tomorrow := now.AddDate(0, 0, 1)
			startOfTomorrow := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 0, 0, 0, 0, tomorrow.Location())
			query = query.Where("planned_time >= ?", startOfTomorrow)
		case "all":
			// 查询所有待办，不限制时间
			// 不添加时间筛选条件，返回所有待办
		case "overdue":
			query = query.Where("planned_time < ? AND status NOT IN (?, ?)", now, models.TodoStatusCompleted, models.TodoStatusCancelled)
		}
	}

	// 自定义时间范围
	if req.StartDate != nil {
		query = query.Where("planned_time >= ?", *req.StartDate)
	}
	if req.EndDate != nil {
		query = query.Where("planned_time <= ?", *req.EndDate)
	}

	// 关键词搜索
	if req.Keyword != "" {
		query = query.Where("title LIKE ? OR content LIKE ?", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}

	// 获取总数
	var total int64
	query.Count(&total)

	// 分页查询
	var todos []models.Todo
	offset := (req.Page - 1) * req.PageSize
	err := query.Preload("Customer").Preload("Creator").Preload("Executor").Preload("ReminderUser").
		Order("created_at DESC").Offset(offset).Limit(req.PageSize).Find(&todos).Error

	if err != nil {
		return nil, 0, err
	}

	// 转换为响应格式
	var responses []models.TodoResponse
	for _, todo := range todos {
		response := models.TodoResponse{
			Todo:         todo,
			CreatorName:  todo.Creator.Name,
			ExecutorName: todo.Executor.Name,
			CustomerName: todo.Customer.Name,
			IsOverdue:    todo.IsOverdue(),
			DaysLeft:     todo.GetDaysLeft(),
		}
		if todo.ReminderUser != nil {
			reminderUserName := todo.ReminderUser.Name
			response.ReminderUserName = &reminderUserName
		}
		responses = append(responses, response)
	}

	return responses, total, nil
}

// GetTodoByID 根据ID获取待办
func (s *TodoService) GetTodoByID(id uint64) (*models.TodoResponse, error) {
	var todo models.Todo
	err := s.db.Preload("Customer").Preload("Creator").Preload("Executor").Preload("ReminderUser").
		Where("id = ? AND is_deleted = ?", id, false).First(&todo).Error

	if err != nil {
		return nil, err
	}

	// 转换为响应格式
	response := models.TodoResponse{
		Todo:         todo,
		CreatorName:  todo.Creator.Name,
		ExecutorName: todo.Executor.Name,
		CustomerName: todo.Customer.Name,
		IsOverdue:    todo.IsOverdue(),
		DaysLeft:     todo.GetDaysLeft(),
	}
	if todo.ReminderUser != nil {
		reminderUserName := todo.ReminderUser.Name
		response.ReminderUserName = &reminderUserName
	}

	return &response, nil
}

// UpdateTodo 更新待办
func (s *TodoService) UpdateTodo(id uint64, req models.TodoUpdateRequest, operatorID uint64) (*models.Todo, error) {
	// 查找待办
	var todo models.Todo
	err := s.db.Where("id = ? AND is_deleted = ?", id, false).First(&todo).Error
	if err != nil {
		return nil, err
	}

	// 保存旧数据用于日志
	oldTodo := todo

	// 更新字段
	if req.Title != nil {
		todo.Title = *req.Title
	}
	if req.Content != nil {
		todo.Content = *req.Content
	}
	if req.Status != nil {
		todo.Status = *req.Status
	}
	if req.PlannedTime != nil {
		todo.PlannedTime = *req.PlannedTime
	}
	if req.ExecutorID != nil {
		todo.ExecutorID = *req.ExecutorID
	}
	if req.IsReminder != nil {
		todo.IsReminder = *req.IsReminder
	}
	if req.ReminderType != nil {
		todo.ReminderType = req.ReminderType
	}
	if req.ReminderUserID != nil {
		todo.ReminderUserID = req.ReminderUserID
	}
	if req.ReminderTime != nil {
		todo.ReminderTime = req.ReminderTime
	}
	if req.Priority != nil {
		todo.Priority = *req.Priority
	}
	if req.Tags != nil {
		todo.Tags = req.Tags
	}

	// 保存更新
	if err := s.db.Save(&todo).Error; err != nil {
		return nil, err
	}

	// 记录操作日志
	s.createTodoLog(todo.ID, operatorID, models.ActionUpdate, oldTodo, todo)

	return &todo, nil
}

// DeleteTodo 删除待办（软删除）
func (s *TodoService) DeleteTodo(id uint64, operatorID uint64) error {
	// 查找待办
	var todo models.Todo
	err := s.db.Where("id = ? AND is_deleted = ?", id, false).First(&todo).Error
	if err != nil {
		return err
	}

	// 软删除
	now := time.Now()
	todo.IsDeleted = true
	todo.DeletedAt = &now

	if err := s.db.Save(&todo).Error; err != nil {
		return err
	}

	// 记录操作日志
	s.createTodoLog(todo.ID, operatorID, models.ActionDelete, todo, nil)

	return nil
}

// CompleteTodo 完成待办
func (s *TodoService) CompleteTodo(id uint64, operatorID uint64) (*models.Todo, error) {
	// 查找待办
	var todo models.Todo
	err := s.db.Where("id = ? AND is_deleted = ?", id, false).First(&todo).Error
	if err != nil {
		return nil, err
	}

	if todo.Status == models.TodoStatusCompleted {
		return &todo, nil // 已完成，直接返回
	}

	// 保存旧数据
	oldTodo := todo

	// 更新状态
	todo.Status = models.TodoStatusCompleted
	now := time.Now()
	todo.CompletedTime = &now

	if err := s.db.Save(&todo).Error; err != nil {
		return nil, err
	}

	// 记录操作日志
	s.createTodoLog(todo.ID, operatorID, models.ActionComplete, oldTodo, todo)

	return &todo, nil
}

// CancelTodo 取消待办
func (s *TodoService) CancelTodo(id uint64, operatorID uint64) (*models.Todo, error) {
	// 查找待办
	var todo models.Todo
	err := s.db.Where("id = ? AND is_deleted = ?", id, false).First(&todo).Error
	if err != nil {
		return nil, err
	}

	if todo.Status == models.TodoStatusCancelled {
		return &todo, nil // 已取消，直接返回
	}

	// 保存旧数据
	oldTodo := todo

	// 更新状态
	todo.Status = models.TodoStatusCancelled

	if err := s.db.Save(&todo).Error; err != nil {
		return nil, err
	}

	// 记录操作日志
	s.createTodoLog(todo.ID, operatorID, models.ActionCancel, oldTodo, todo)

	return &todo, nil
}

// GetTodoStats 获取待办统计
func (s *TodoService) GetTodoStats(customerID, executorID *uint64) (map[string]int64, error) {
	query := s.db.Model(&models.Todo{}).Where("is_deleted = ?", false)

	if customerID != nil {
		query = query.Where("customer_id = ?", *customerID)
	}
	if executorID != nil {
		query = query.Where("executor_id = ?", *executorID)
	}

	// 统计各状态数量
	stats := make(map[string]int64)
	statuses := []models.TodoStatus{models.TodoStatusPending, models.TodoStatusCompleted, models.TodoStatusOverdue, models.TodoStatusCancelled}

	for _, status := range statuses {
		var count int64
		query.Where("status = ?", status).Count(&count)
		stats[string(status)] = count
	}

	// 统计延期数量
	var overdueCount int64
	now := time.Now()
	query.Where("planned_time < ? AND status NOT IN (?, ?)", now, models.TodoStatusCompleted, models.TodoStatusCancelled).Count(&overdueCount)
	stats["overdue"] = overdueCount

	// 统计今日待办
	var todayCount int64
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)
	query.Where("planned_time >= ? AND planned_time < ?", startOfDay, endOfDay).Count(&todayCount)
	stats["today"] = todayCount

	return stats, nil
}

// createTodoLog 创建待办操作日志
func (s *TodoService) createTodoLog(todoID, operatorID uint64, action models.ActionType, oldData, newData interface{}) {
	log := models.TodoLog{
		TodoID:     todoID,
		OperatorID: operatorID,
		Action:     action,
	}

	if oldData != nil {
		oldDataBytes, _ := json.Marshal(oldData)
		oldDataMap := make(map[string]interface{})
		json.Unmarshal(oldDataBytes, &oldDataMap)
		log.OldData = models.JSONB(oldDataMap)
	}

	if newData != nil {
		newDataBytes, _ := json.Marshal(newData)
		newDataMap := make(map[string]interface{})
		json.Unmarshal(newDataBytes, &newDataMap)
		log.NewData = models.JSONB(newDataMap)
	}

	s.db.Create(&log)
}
