package models

import (
	"time"

	"gorm.io/gorm"
)

// TodoStatus 待办状态枚举
type TodoStatus string

const (
	TodoStatusPending   TodoStatus = "pending"   // 未完成
	TodoStatusCompleted TodoStatus = "completed" // 已完成
	TodoStatusOverdue   TodoStatus = "overdue"   // 延期
	TodoStatusCancelled TodoStatus = "cancelled" // 取消
)

// ReminderType 提醒方式枚举
type ReminderType string

const (
	ReminderTypeWechat           ReminderType = "wechat"            // 微信
	ReminderTypeEnterpriseWechat ReminderType = "enterprise_wechat" // 企业微信
	ReminderTypeBoth             ReminderType = "both"              // 两者
	ReminderTypeSMS              ReminderType = "sms"               // 短信
)

// Priority 优先级枚举
type Priority string

const (
	PriorityLow    Priority = "low"    // 低
	PriorityMedium Priority = "medium" // 中
	PriorityHigh   Priority = "high"   // 高
	PriorityUrgent Priority = "urgent" // 紧急
)

// ActionType 操作类型枚举
type ActionType string

const (
	ActionCreate   ActionType = "create"   // 创建
	ActionUpdate   ActionType = "update"   // 更新
	ActionDelete   ActionType = "delete"   // 删除
	ActionComplete ActionType = "complete" // 完成
	ActionCancel   ActionType = "cancel"   // 取消
)

// Todo 待办事项模型
type Todo struct {
	ID             uint64        `json:"id" gorm:"primaryKey;autoIncrement;comment:待办ID"`
	CustomerID     uint64        `json:"customer_id" gorm:"not null;index;comment:关联客户ID"`
	CreatorID      uint64        `json:"creator_id" gorm:"not null;index;comment:创建人ID"`
	ExecutorID     uint64        `json:"executor_id" gorm:"not null;index;comment:执行人ID"`
	Title          string        `json:"title" gorm:"type:varchar(255);not null;comment:待办标题"`
	Content        string        `json:"content" gorm:"type:text;comment:待办内容详情"`
	Status         TodoStatus    `json:"status" gorm:"type:enum('pending','completed','overdue','cancelled');default:pending;index;comment:待办状态"`
	PlannedTime    time.Time     `json:"planned_time" gorm:"not null;index;comment:计划执行时间"`
	CompletedTime  *time.Time    `json:"completed_time" gorm:"comment:完成时间"`
	IsReminder     bool          `json:"is_reminder" gorm:"default:false;comment:是否提醒"`
	ReminderType   *ReminderType `json:"reminder_type" gorm:"type:enum('wechat','enterprise_wechat','both','sms');comment:提醒方式"`
	ReminderUserID *uint64       `json:"reminder_user_id" gorm:"comment:提醒人ID"`
	ReminderTime   *time.Time    `json:"reminder_time" gorm:"comment:提醒时间"`
	Priority       Priority      `json:"priority" gorm:"type:enum('low','medium','high','urgent');default:medium;comment:优先级"`
	Tags           JSONB         `json:"tags" gorm:"type:json;comment:标签"`
	Attachments    JSONB         `json:"attachments" gorm:"type:json;comment:附件信息"`
	CreatedAt      time.Time     `json:"created_at" gorm:"index;comment:创建时间"`
	UpdatedAt      time.Time     `json:"updated_at" gorm:"comment:更新时间"`
	DeletedAt      *time.Time    `json:"deleted_at" gorm:"comment:删除时间"`
	IsDeleted      bool          `json:"is_deleted" gorm:"default:false;index;comment:是否删除"`

	// 关联关系
	Customer     Customer `json:"customer" gorm:"foreignKey:CustomerID"`
	Creator      User     `json:"creator" gorm:"foreignKey:CreatorID"`
	Executor     User     `json:"executor" gorm:"foreignKey:ExecutorID"`
	ReminderUser *User    `json:"reminder_user" gorm:"foreignKey:ReminderUserID"`
}

// TodoLog 待办操作日志模型
type TodoLog struct {
	ID         uint64     `json:"id" gorm:"primaryKey;autoIncrement;comment:日志ID"`
	TodoID     uint64     `json:"todo_id" gorm:"not null;index;comment:待办ID"`
	OperatorID uint64     `json:"operator_id" gorm:"not null;index;comment:操作人ID"`
	Action     ActionType `json:"action" gorm:"type:enum('create','update','delete','complete','cancel');not null;index;comment:操作类型"`
	OldData    JSONB      `json:"old_data" gorm:"type:json;comment:变更前数据"`
	NewData    JSONB      `json:"new_data" gorm:"type:json;comment:变更后数据"`
	Remark     string     `json:"remark" gorm:"type:varchar(500);comment:操作备注"`
	CreatedAt  time.Time  `json:"created_at" gorm:"index;comment:操作时间"`

	// 关联关系
	Todo     Todo `json:"todo" gorm:"foreignKey:TodoID"`
	Operator User `json:"operator" gorm:"foreignKey:OperatorID"`
}

// 移除内联的User模型定义，现在使用独立的user.go文件

// TodoCreateRequest 创建待办请求
type TodoCreateRequest struct {
	CustomerID     uint64        `json:"customer_id" binding:"required"`
	ExecutorID     uint64        `json:"executor_id" binding:"required"`
	Title          string        `json:"title" binding:"required,max=255"`
	Content        string        `json:"content"`
	PlannedTime    time.Time     `json:"planned_time" binding:"required"`
	IsReminder     bool          `json:"is_reminder"`
	ReminderType   *ReminderType `json:"reminder_type"`
	ReminderUserID *uint64       `json:"reminder_user_id"`
	ReminderTime   *time.Time    `json:"reminder_time"`
	Priority       Priority      `json:"priority"`
	Tags           JSONB         `json:"tags"`
}

// TodoUpdateRequest 更新待办请求
type TodoUpdateRequest struct {
	Title          *string       `json:"title"`
	Content        *string       `json:"content"`
	Status         *TodoStatus   `json:"status"`
	PlannedTime    *time.Time    `json:"planned_time"`
	ExecutorID     *uint64       `json:"executor_id"`
	IsReminder     *bool         `json:"is_reminder"`
	ReminderType   *ReminderType `json:"reminder_type"`
	ReminderUserID *uint64       `json:"reminder_user_id"`
	ReminderTime   *time.Time    `json:"reminder_time"`
	Priority       *Priority     `json:"priority"`
	Tags           JSONB         `json:"tags"`
}

// TodoQueryRequest 查询待办请求
type TodoQueryRequest struct {
	CustomerID *uint64     `json:"customer_id" form:"customer_id"`
	ExecutorID *uint64     `json:"executor_id" form:"executor_id"`
	CreatorID  *uint64     `json:"creator_id" form:"creator_id"`
	Status     *TodoStatus `json:"status" form:"status"`
	Priority   *Priority   `json:"priority" form:"priority"`
	DateType   string      `json:"date_type" form:"date_type"` // yesterday, today, tomorrow, overdue
	StartDate  *time.Time  `json:"start_date" form:"start_date"`
	EndDate    *time.Time  `json:"end_date" form:"end_date"`
	Keyword    string      `json:"keyword" form:"keyword"`
	Page       int         `json:"page" form:"page"`
	PageSize   int         `json:"page_size" form:"page_size"`
}

// TodoResponse 待办响应
type TodoResponse struct {
	Todo
	CreatorName      string  `json:"creator_name"`
	ExecutorName     string  `json:"executor_name"`
	CustomerName     string  `json:"customer_name"`
	ReminderUserName *string `json:"reminder_user_name"`
	IsOverdue        bool    `json:"is_overdue"`
	DaysLeft         int     `json:"days_left"`
}

// TableName 指定表名
func (Todo) TableName() string {
	return "todos"
}

// TableName 指定表名
func (TodoLog) TableName() string {
	return "todo_logs"
}

// BeforeCreate GORM钩子：创建前
func (t *Todo) BeforeCreate(tx *gorm.DB) error {
	if t.Status == "" {
		t.Status = TodoStatusPending
	}
	if t.Priority == "" {
		t.Priority = PriorityMedium
	}
	return nil
}

// BeforeUpdate GORM钩子：更新前
func (t *Todo) BeforeUpdate(tx *gorm.DB) error {
	// 如果状态变为已完成，设置完成时间
	if t.Status == TodoStatusCompleted && t.CompletedTime == nil {
		now := time.Now()
		t.CompletedTime = &now
	}
	return nil
}

// IsOverdue 判断是否延期
func (t *Todo) IsOverdue() bool {
	if t.Status == TodoStatusCompleted || t.Status == TodoStatusCancelled {
		return false
	}
	return time.Now().After(t.PlannedTime)
}

// GetDaysLeft 获取剩余天数
func (t *Todo) GetDaysLeft() int {
	if t.Status == TodoStatusCompleted || t.Status == TodoStatusCancelled {
		return 0
	}
	duration := time.Until(t.PlannedTime)
	return int(duration.Hours() / 24)
}
