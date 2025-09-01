package services

import (
	models2 "crm/internal/domain"
	"crm/internal/repository"
	"encoding/json"
	"time"
)

// TodoService 待办服务
type TodoService struct {
	repo repository.TodoRepository
}

// NewTodoService 创建待办服务
func NewTodoService(repo repository.TodoRepository) *TodoService {
	return &TodoService{repo: repo}
}

// CreateTodo 创建待办
func (s *TodoService) CreateTodo(req models2.TodoCreateRequest, creatorID uint64) (*models2.Todo, error) {
	todo := models2.Todo{
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

	if err := s.repo.Create(&todo); err != nil {
		return nil, err
	}

	// 记录操作日志
	s.createTodoLog(todo.ID, creatorID, models2.ActionCreate, nil, todo)

	return &todo, nil
}

// GetTodos 获取待办列表
func (s *TodoService) GetTodos(req models2.TodoQueryRequest) ([]models2.TodoResponse, int64, error) {
	return s.repo.GetList(req)
}

// GetTodoByID 根据ID获取待办
func (s *TodoService) GetTodoByID(id uint64) (*models2.TodoResponse, error) {
	todo, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// 转换为响应格式
	response := models2.TodoResponse{
		Todo:         *todo,
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

// CompleteTodo 完成待办
func (s *TodoService) CompleteTodo(id uint64, operatorID uint64) (*models2.Todo, error) {
	todo, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if todo.Status == models2.TodoStatusCompleted {
		return todo, nil // 已完成，直接返回
	}

	// 保存旧数据
	oldTodo := *todo

	// 更新状态
	todo.Status = models2.TodoStatusCompleted
	now := time.Now()
	todo.CompletedTime = &now

	if err := s.repo.Update(todo); err != nil {
		return nil, err
	}

	// 记录操作日志
	s.createTodoLog(todo.ID, operatorID, models2.ActionComplete, oldTodo, *todo)

	return todo, nil
}

// CancelTodo 取消待办
func (s *TodoService) CancelTodo(id uint64, operatorID uint64) (*models2.Todo, error) {
	todo, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if todo.Status == models2.TodoStatusCancelled {
		return todo, nil // 已取消，直接返回
	}

	// 保存旧数据
	oldTodo := *todo

	// 更新状态
	todo.Status = models2.TodoStatusCancelled

	if err := s.repo.Update(todo); err != nil {
		return nil, err
	}

	// 记录操作日志
	s.createTodoLog(todo.ID, operatorID, models2.ActionCancel, oldTodo, *todo)

	return todo, nil
}

// UpdateTodo 更新待办
func (s *TodoService) UpdateTodo(id uint64, req models2.TodoUpdateRequest, operatorID uint64) (*models2.Todo, error) {
	// 查找待办
	todo, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// 保存旧数据用于日志
	oldTodo := *todo

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
	if err := s.repo.Update(todo); err != nil {
		return nil, err
	}

	// 记录操作日志
	s.createTodoLog(todo.ID, operatorID, models2.ActionUpdate, oldTodo, *todo)

	return todo, nil
}

// DeleteTodo 删除待办（软删除）
func (s *TodoService) DeleteTodo(id uint64, operatorID uint64) error {
	// 查找待办
	todo, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	// 软删除
	if err := s.repo.Delete(id); err != nil {
		return err
	}

	// 记录操作日志
	s.createTodoLog(todo.ID, operatorID, models2.ActionDelete, *todo, nil)

	return nil
}

// GetTodoStats 获取待办统计
func (s *TodoService) GetTodoStats(customerID, executorID *uint64) (map[string]int64, error) {
	// TODO: 将统计逻辑移动到repository层
	return make(map[string]int64), nil
}

// createTodoLog 创建待办操作日志
func (s *TodoService) createTodoLog(todoID, operatorID uint64, action models2.ActionType, oldData, newData interface{}) {
	log := models2.TodoLog{
		TodoID:     todoID,
		OperatorID: operatorID,
		Action:     action,
	}

	if oldData != nil {
		oldDataBytes, _ := json.Marshal(oldData)
		oldDataMap := make(map[string]interface{})
		json.Unmarshal(oldDataBytes, &oldDataMap)
		log.OldData = models2.JSONB(oldDataMap)
	}

	if newData != nil {
		newDataBytes, _ := json.Marshal(newData)
		newDataMap := make(map[string]interface{})
		json.Unmarshal(newDataBytes, &newDataMap)
		log.NewData = models2.JSONB(newDataMap)
	}

	s.repo.CreateLog(&log)
}
