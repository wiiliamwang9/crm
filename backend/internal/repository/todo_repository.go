package repository

import (
	"crm/internal/domain"
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

type todoRepository struct {
	db *gorm.DB
}

// NewTodoRepository 创建待办仓储实例
func NewTodoRepository(db *gorm.DB) TodoRepository {
	return &todoRepository{db: db}
}

// GetByID 根据ID获取待办
func (r *todoRepository) GetByID(id uint64) (*domain.Todo, error) {
	var todo domain.Todo
	err := r.db.Preload("Customer").Preload("Creator").Preload("Executor").Preload("ReminderUser").
		Where("id = ? AND is_deleted = ?", id, false).First(&todo).Error
	return &todo, err
}

// GetList 获取待办列表
func (r *todoRepository) GetList(req domain.TodoQueryRequest) ([]domain.TodoResponse, int64, error) {
	// 构建查询
	query := r.db.Model(&domain.Todo{}).Where("is_deleted = ?", false)

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
		case "overdue":
			query = query.Where("planned_time < ? AND status NOT IN (?, ?)", now, domain.TodoStatusCompleted, domain.TodoStatusCancelled)
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
	var todos []domain.Todo
	offset := (req.Page - 1) * req.PageSize
	err := query.Preload("Customer").Preload("Creator").Preload("Executor").Preload("ReminderUser").
		Order("created_at DESC").Offset(offset).Limit(req.PageSize).Find(&todos).Error

	if err != nil {
		return nil, 0, err
	}

	// 转换为响应格式
	var responses []domain.TodoResponse
	for _, todo := range todos {
		response := domain.TodoResponse{
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

// Create 创建待办
func (r *todoRepository) Create(todo *domain.Todo) error {
	// 设置默认值
	if todo.Status == "" {
		todo.Status = domain.TodoStatusPending
	}
	if todo.Priority == "" {
		todo.Priority = domain.PriorityMedium
	}
	return r.db.Create(todo).Error
}

// Update 更新待办
func (r *todoRepository) Update(todo *domain.Todo) error {
	// 如果状态变为已完成，设置完成时间
	if todo.Status == domain.TodoStatusCompleted && todo.CompletedTime == nil {
		now := time.Now()
		todo.CompletedTime = &now
	}
	return r.db.Save(todo).Error
}

// Delete 删除待办（软删除）
func (r *todoRepository) Delete(id uint64) error {
	now := time.Now()
	result := r.db.Model(&domain.Todo{}).Where("id = ? AND is_deleted = ?", id, false).
		Updates(map[string]interface{}{
			"is_deleted": true,
			"deleted_at": &now,
		})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// UpdateStatus 更新待办状态
func (r *todoRepository) UpdateStatus(id uint64, status domain.TodoStatus) error {
	updates := map[string]interface{}{
		"status": status,
	}

	if status == domain.TodoStatusCompleted {
		now := time.Now()
		updates["completed_time"] = &now
	}

	result := r.db.Model(&domain.Todo{}).Where("id = ? AND is_deleted = ?", id, false).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// GetOverdueTodos 获取延期待办
func (r *todoRepository) GetOverdueTodos() ([]domain.Todo, error) {
	var todos []domain.Todo
	now := time.Now()
	err := r.db.Where("planned_time < ? AND status NOT IN (?, ?) AND is_deleted = ?",
		now, domain.TodoStatusCompleted, domain.TodoStatusCancelled, false).Find(&todos).Error
	return todos, err
}

// CreateLog 创建待办操作日志
func (r *todoRepository) CreateLog(log *domain.TodoLog) error {
	return r.db.Create(log).Error
}

// createTodoLog 创建待办操作日志的辅助方法
func (r *todoRepository) createTodoLog(todoID, operatorID uint64, action domain.ActionType, oldData, newData interface{}) {
	log := domain.TodoLog{
		TodoID:     todoID,
		OperatorID: operatorID,
		Action:     action,
	}

	if oldData != nil {
		oldDataBytes, _ := json.Marshal(oldData)
		oldDataMap := make(map[string]interface{})
		json.Unmarshal(oldDataBytes, &oldDataMap)
		log.OldData = domain.JSONB(oldDataMap)
	}

	if newData != nil {
		newDataBytes, _ := json.Marshal(newData)
		newDataMap := make(map[string]interface{})
		json.Unmarshal(newDataBytes, &newDataMap)
		log.NewData = domain.JSONB(newDataMap)
	}

	r.db.Create(&log)
}
