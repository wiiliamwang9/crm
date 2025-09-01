package models

import (
	"time"
)

// ReminderStatus 提醒状态枚举
type ReminderStatus string

const (
	ReminderStatusPending   ReminderStatus = "pending"   // 待发送
	ReminderStatusSent      ReminderStatus = "sent"      // 已发送
	ReminderStatusFailed    ReminderStatus = "failed"    // 发送失败
	ReminderStatusCancelled ReminderStatus = "cancelled" // 已取消
)

// ReminderFrequency 提醒频率枚举
type ReminderFrequency string

const (
	ReminderFrequencyOnce    ReminderFrequency = "once"    // 单次提醒
	ReminderFrequencyDaily   ReminderFrequency = "daily"   // 每日提醒
	ReminderFrequencyWeekly  ReminderFrequency = "weekly"  // 每周提醒
	ReminderFrequencyMonthly ReminderFrequency = "monthly" // 每月提醒
)

// Reminder 提醒记录模型
type Reminder struct {
	ID           uint64            `json:"id" gorm:"primaryKey;autoIncrement;comment:提醒ID"`
	TodoID       uint64            `json:"todo_id" gorm:"not null;index;comment:关联待办ID"`
	UserID       uint64            `json:"user_id" gorm:"not null;index;comment:提醒用户ID"`
	Type         ReminderType      `json:"type" gorm:"type:enum('wechat','enterprise_wechat','both','sms');not null;comment:提醒方式"`
	Title        string            `json:"title" gorm:"type:varchar(255);not null;comment:提醒标题"`
	Content      string            `json:"content" gorm:"type:text;comment:提醒内容"`
	Status       ReminderStatus    `json:"status" gorm:"type:enum('pending','sent','failed','cancelled');default:pending;index;comment:提醒状态"`
	Frequency    ReminderFrequency `json:"frequency" gorm:"type:enum('once','daily','weekly','monthly');default:once;comment:提醒频率"`
	ScheduleTime time.Time         `json:"schedule_time" gorm:"not null;index;comment:计划提醒时间"`
	SentTime     *time.Time        `json:"sent_time" gorm:"comment:实际发送时间"`
	FailReason   string            `json:"fail_reason" gorm:"type:varchar(500);comment:失败原因"`
	RetryCount   int               `json:"retry_count" gorm:"default:0;comment:重试次数"`
	MaxRetries   int               `json:"max_retries" gorm:"default:3;comment:最大重试次数"`
	CreatedAt    time.Time         `json:"created_at" gorm:"index;comment:创建时间"`
	UpdatedAt    time.Time         `json:"updated_at" gorm:"comment:更新时间"`

	// 关联关系
	Todo Todo `json:"todo" gorm:"foreignKey:TodoID"`
	User User `json:"user" gorm:"foreignKey:UserID"`
}

// ReminderTemplate 提醒模板
type ReminderTemplate struct {
	ID        uint64       `json:"id" gorm:"primaryKey;autoIncrement;comment:模板ID"`
	Name      string       `json:"name" gorm:"type:varchar(100);not null;comment:模板名称"`
	Type      ReminderType `json:"type" gorm:"type:enum('wechat','enterprise_wechat','both','sms');not null;comment:适用的提醒方式"`
	Title     string       `json:"title" gorm:"type:varchar(255);not null;comment:标题模板"`
	Content   string       `json:"content" gorm:"type:text;not null;comment:内容模板"`
	Variables JSONB        `json:"variables" gorm:"type:json;comment:可用变量说明"`
	IsActive  bool         `json:"is_active" gorm:"default:true;comment:是否启用"`
	IsDefault bool         `json:"is_default" gorm:"default:false;comment:是否为默认模板"`
	CreatedBy uint64       `json:"created_by" gorm:"not null;comment:创建人ID"`
	CreatedAt time.Time    `json:"created_at" gorm:"comment:创建时间"`
	UpdatedAt time.Time    `json:"updated_at" gorm:"comment:更新时间"`
}

// ReminderConfig 用户提醒配置
type ReminderConfig struct {
	ID                     uint64    `json:"id" gorm:"primaryKey;autoIncrement;comment:配置ID"`
	UserID                 uint64    `json:"user_id" gorm:"not null;unique;comment:用户ID"`
	EnableWechat           bool      `json:"enable_wechat" gorm:"default:true;comment:启用微信提醒"`
	EnableEnterpriseWechat bool      `json:"enable_enterprise_wechat" gorm:"default:true;comment:启用企业微信提醒"`
	WechatUserID           string    `json:"wechat_user_id" gorm:"type:varchar(100);comment:微信用户ID"`
	EnterpriseWechatUserID string    `json:"enterprise_wechat_user_id" gorm:"type:varchar(100);comment:企业微信用户ID"`
	DefaultAdvanceMinutes  int       `json:"default_advance_minutes" gorm:"default:30;comment:默认提前提醒分钟数"`
	QuietStartTime         string    `json:"quiet_start_time" gorm:"type:varchar(5);default:22:00;comment:免打扰开始时间"`
	QuietEndTime           string    `json:"quiet_end_time" gorm:"type:varchar(5);default:08:00;comment:免打扰结束时间"`
	CreatedAt              time.Time `json:"created_at" gorm:"comment:创建时间"`
	UpdatedAt              time.Time `json:"updated_at" gorm:"comment:更新时间"`

	// 关联关系
	User User `json:"user" gorm:"foreignKey:UserID"`
}

// ReminderCreateRequest 创建提醒请求
type ReminderCreateRequest struct {
	TodoID       uint64            `json:"todo_id" binding:"required"`
	UserID       uint64            `json:"user_id" binding:"required"`
	Type         ReminderType      `json:"type" binding:"required"`
	Title        string            `json:"title" binding:"required,max=255"`
	Content      string            `json:"content"`
	Frequency    ReminderFrequency `json:"frequency"`
	ScheduleTime time.Time         `json:"schedule_time" binding:"required"`
	MaxRetries   int               `json:"max_retries"`
}

// ReminderUpdateRequest 更新提醒请求
type ReminderUpdateRequest struct {
	Title        *string            `json:"title"`
	Content      *string            `json:"content"`
	Status       *ReminderStatus    `json:"status"`
	Frequency    *ReminderFrequency `json:"frequency"`
	ScheduleTime *time.Time         `json:"schedule_time"`
	MaxRetries   *int               `json:"max_retries"`
}

// ReminderQueryRequest 查询提醒请求
type ReminderQueryRequest struct {
	TodoID    *uint64         `json:"todo_id" form:"todo_id"`
	UserID    *uint64         `json:"user_id" form:"user_id"`
	Status    *ReminderStatus `json:"status" form:"status"`
	Type      *ReminderType   `json:"type" form:"type"`
	StartDate *time.Time      `json:"start_date" form:"start_date"`
	EndDate   *time.Time      `json:"end_date" form:"end_date"`
	Page      int             `json:"page" form:"page"`
	PageSize  int             `json:"page_size" form:"page_size"`
}

// ReminderResponse 提醒响应
type ReminderResponse struct {
	Reminder
	TodoTitle    string `json:"todo_title"`
	UserName     string `json:"user_name"`
	CustomerName string `json:"customer_name"`
}

// TableName 指定表名
func (Reminder) TableName() string {
	return "reminders"
}

// TableName 指定表名
func (ReminderTemplate) TableName() string {
	return "reminder_templates"
}

// TableName 指定表名
func (ReminderConfig) TableName() string {
	return "reminder_configs"
}

// CanRetry 判断是否可以重试
func (r *Reminder) CanRetry() bool {
	return r.Status == ReminderStatusFailed && r.RetryCount < r.MaxRetries
}

// IsInQuietTime 判断是否在免打扰时间内
func (rc *ReminderConfig) IsInQuietTime(t time.Time) bool {
	if rc.QuietStartTime == "" || rc.QuietEndTime == "" {
		return false
	}

	timeStr := t.Format("15:04")

	// 跨天的情况（如22:00-08:00）
	if rc.QuietStartTime > rc.QuietEndTime {
		return timeStr >= rc.QuietStartTime || timeStr <= rc.QuietEndTime
	} else {
		// 同天的情况（如12:00-14:00）
		return timeStr >= rc.QuietStartTime && timeStr <= rc.QuietEndTime
	}
}

// GetNextScheduleTime 获取下次提醒时间
func (r *Reminder) GetNextScheduleTime() *time.Time {
	if r.Frequency == ReminderFrequencyOnce {
		return nil // 单次提醒不需要下次时间
	}

	var nextTime time.Time
	switch r.Frequency {
	case ReminderFrequencyDaily:
		nextTime = r.ScheduleTime.AddDate(0, 0, 1)
	case ReminderFrequencyWeekly:
		nextTime = r.ScheduleTime.AddDate(0, 0, 7)
	case ReminderFrequencyMonthly:
		nextTime = r.ScheduleTime.AddDate(0, 1, 0)
	default:
		return nil
	}

	return &nextTime
}
