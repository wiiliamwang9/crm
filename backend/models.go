package main

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/lib/pq"
)

// ============================================================================
// 公共基础结构和类型定义
// ============================================================================

// BaseModel 基础模型，包含公共字段
type BaseModel struct {
	CreatedAt time.Time  `json:"created_at" gorm:"index;comment:创建时间"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"comment:更新时间"`
	DeletedAt *time.Time `json:"deleted_at" gorm:"comment:删除时间"`
	IsDeleted bool       `json:"is_deleted" gorm:"default:false;index;comment:是否删除"`
}

// JSONB 自定义JSONB类型
type JSONB map[string]interface{}

func (j JSONB) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return json.Marshal(j)
}

func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		*j = nil
		return nil
	}

	if len(bytes) == 0 {
		*j = nil
		return nil
	}

	// 先尝试解析为对象
	if err := json.Unmarshal(bytes, j); err == nil {
		return nil
	}

	// 如果解析失败，可能是数组或其他类型，设置为nil
	*j = nil
	return nil
}

// ============================================================================
// 枚举定义
// ============================================================================

// TodoStatus 待办状态枚举
type TodoStatus string

const (
	TodoStatusPending   TodoStatus = "pending"
	TodoStatusCompleted TodoStatus = "completed"
	TodoStatusOverdue   TodoStatus = "overdue"
	TodoStatusCancelled TodoStatus = "cancelled"
)

// ReminderType 提醒方式枚举
type ReminderType string

const (
	ReminderTypeWechat           ReminderType = "wechat"
	ReminderTypeEnterpriseWechat ReminderType = "enterprise_wechat"
	ReminderTypeBoth             ReminderType = "both"
	ReminderTypeSMS              ReminderType = "sms"
)

// Priority 优先级枚举
type Priority string

const (
	PriorityLow    Priority = "low"
	PriorityMedium Priority = "medium"
	PriorityHigh   Priority = "high"
	PriorityUrgent Priority = "urgent"
)

// ActionType 操作类型枚举
type ActionType string

const (
	ActionCreate   ActionType = "create"
	ActionUpdate   ActionType = "update"
	ActionDelete   ActionType = "delete"
	ActionComplete ActionType = "complete"
	ActionCancel   ActionType = "cancel"
)

// ReminderStatus 提醒状态枚举
type ReminderStatus string

const (
	ReminderStatusPending   ReminderStatus = "pending"
	ReminderStatusSent      ReminderStatus = "sent"
	ReminderStatusFailed    ReminderStatus = "failed"
	ReminderStatusCancelled ReminderStatus = "cancelled"
)

// ReminderFrequency 提醒频率枚举
type ReminderFrequency string

const (
	ReminderFrequencyOnce    ReminderFrequency = "once"
	ReminderFrequencyDaily   ReminderFrequency = "daily"
	ReminderFrequencyWeekly  ReminderFrequency = "weekly"
	ReminderFrequencyMonthly ReminderFrequency = "monthly"
)

// ============================================================================
// 数据模型定义
// ============================================================================

// Customer 客户模型
type Customer struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"size:256"`
	ContactName string         `json:"contact_name" gorm:"size:256"`
	Gender      int            `json:"gender"`
	Avatar      string         `json:"avatar" gorm:"size:2048"`
	Photos      pq.StringArray `json:"photos" gorm:"type:varchar(2048)[]"`
	Remark      string         `json:"remark" gorm:"type:text"`

	Source    string    `json:"source" gorm:"size:256"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy uint      `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy uint      `json:"updated_by"`

	Phones        pq.StringArray `json:"phones" gorm:"type:varchar(128)[];default:null"`
	Wechats       pq.StringArray `json:"wechats" gorm:"type:varchar(128)[]"`
	Douyins       pq.StringArray `json:"douyins" gorm:"type:varchar(128)[]"`
	Kwais         pq.StringArray `json:"kwais" gorm:"type:varchar(128)[]"`
	Redbooks      pq.StringArray `json:"redbooks" gorm:"type:varchar(128)[]"`
	WeworkOpenids pq.StringArray `json:"wework_openids" gorm:"type:varchar(128)[]"`

	Province   string  `json:"province" gorm:"size:256"`
	City       string  `json:"city" gorm:"size:256"`
	District   string  `json:"district" gorm:"size:256"`
	DistrictID int     `json:"district_id"`
	Street     string  `json:"street" gorm:"size:256"`
	Address    string  `json:"address" gorm:"size:2048"`
	Lat        float64 `json:"lat"`
	Lon        float64 `json:"lon"`

	Category    string         `json:"category" gorm:"size:256"`
	Flags       int            `json:"flags"`
	Tags        pq.StringArray `json:"tags" gorm:"type:varchar(128)[]"`
	Level       int            `json:"level"`
	State       int            `json:"state"`
	Kind        int            `json:"kind"`
	AddedWechat bool           `json:"added_wechat"`

	WorkPhone   pq.StringArray `json:"work_phone" gorm:"type:varchar(256)[]"`
	WorkWechat  pq.StringArray `json:"work_wechat" gorm:"type:varchar(256)[]"`
	CreditSale  float64        `json:"credit_sale" gorm:"type:decimal"`
	Sellers     pq.Int64Array  `json:"sellers" gorm:"type:int4[]"`
	LastVisited *time.Time     `json:"last_visited"`
	LastCalled  *time.Time     `json:"last_called"`

	GroupID        pq.Int64Array  `json:"group_id" gorm:"type:int4[]"`
	BirthPlace     string         `json:"birth_place" gorm:"size:256"`
	BirthYear      int            `json:"birth_year"`
	BirthMonth     int            `json:"birth_month"`
	BirthDate      int            `json:"birth_date"`
	Favors         JSONB          `json:"favors" gorm:"type:jsonb[]"`
	Products       pq.StringArray `json:"products" gorm:"type:varchar(512)[]"`
	AnnualTurnover string         `json:"annual_turnover" gorm:"size:512"`
	ShippingInfos  JSONB          `json:"shipping_infos" gorm:"type:jsonb[]"`
	ExtraInfo      JSONB          `json:"extra_info" gorm:"type:jsonb"`

	OriginalCustomerID string `json:"original_customer_id" gorm:"size:256"`
	ImportSource       string `json:"import_source" gorm:"size:256"`
	SallerName         string `json:"saller_name" gorm:"size:256"`

	SystemTags pq.Int64Array `json:"system_tags" gorm:"type:int4[]"`

	LastOrderDate           *time.Time `json:"last_order_date" gorm:"comment:最后下单时间"`
	OrderCount              int        `json:"order_count" gorm:"default:0;comment:订单数量"`
	AvgOrderValue           *float64   `json:"avg_order_value" gorm:"type:decimal(15,2);comment:平均订单金额"`
	PreferredDeliveryMethod string     `json:"preferred_delivery_method" gorm:"type:varchar(128);comment:偏好配送方式"` // 修正字段长度以匹配数据库
}

func (Customer) TableName() string {
	return "customers"
}

// User 用户模型
type User struct {
	ID                 uint64     `json:"id" gorm:"primaryKey;autoIncrement;comment:用户ID"`
	Name               string     `json:"name" gorm:"type:varchar(256);not null;comment:用户名"`
	Username           string     `json:"username" gorm:"type:text;comment:用户名（登录用）"`
	ManagerID          *uint64    `json:"manager_id" gorm:"comment:主管ID"`
	Email              string     `json:"email" gorm:"type:varchar(256);comment:邮箱"`
	Phone              string     `json:"phone" gorm:"type:varchar(32);comment:手机号"`
	Department         string     `json:"department" gorm:"type:varchar(128);comment:部门"`
	DepartmentLeaderID *uint64    `json:"department_leader_id" gorm:"comment:部门领导ID"`
	Position           string     `json:"position" gorm:"type:varchar(128);comment:职位"`
	WechatWorkID       string     `json:"wechat_work_id" gorm:"type:varchar(128);comment:企业微信ID"`
	WechatID           string     `json:"wechat_id" gorm:"type:varchar(128);comment:微信ID"`
	Status             string     `json:"status" gorm:"type:varchar(32);default:active;comment:状态"`
	AvatarURL          string     `json:"avatar_url" gorm:"type:varchar(512);comment:头像URL"`
	LastLoginAt        *time.Time `json:"last_login_at" gorm:"comment:最后登录时间"`
	BaseModel

	Manager          *User `json:"manager,omitempty" gorm:"foreignKey:ManagerID"`
	DepartmentLeader *User `json:"department_leader,omitempty" gorm:"foreignKey:DepartmentLeaderID"`
}

func (User) TableName() string {
	return "users"
}

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
	BaseModel

	Customer     Customer `json:"customer" gorm:"foreignKey:CustomerID"`
	Creator      User     `json:"creator" gorm:"foreignKey:CreatorID"`
	Executor     User     `json:"executor" gorm:"foreignKey:ExecutorID"`
	ReminderUser *User    `json:"reminder_user" gorm:"foreignKey:ReminderUserID"`
}

func (Todo) TableName() string {
	return "todos"
}

// IsOverdue 检查待办是否已过期
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

// TodoLog 待办操作日志
type TodoLog struct {
	ID         uint64     `json:"id" gorm:"primaryKey;autoIncrement;comment:日志ID"`
	TodoID     uint64     `json:"todo_id" gorm:"not null;index;comment:待办ID"`
	OperatorID uint64     `json:"operator_id" gorm:"not null;index;comment:操作人ID"`
	Action     ActionType `json:"action" gorm:"type:varchar(32);not null;index;comment:操作类型"`
	OldData    JSONB      `json:"old_data" gorm:"type:json;comment:变更前数据"`
	NewData    JSONB      `json:"new_data" gorm:"type:json;comment:变更后数据"`
	Remark     string     `json:"remark" gorm:"type:varchar(500);comment:操作备注"`
	CreatedAt  time.Time  `json:"created_at" gorm:"index;comment:操作时间"`

	Todo     Todo `json:"todo" gorm:"foreignKey:TodoID"`
	Operator User `json:"operator" gorm:"foreignKey:OperatorID"`
}

func (TodoLog) TableName() string {
	return "todo_logs"
}

// Reminder 提醒记录模型
type Reminder struct {
	ID           uint64            `json:"id" gorm:"primaryKey;autoIncrement;comment:提醒ID"`
	TodoID       uint64            `json:"todo_id" gorm:"not null;index;comment:关联待办ID"`
	UserID       uint64            `json:"user_id" gorm:"not null;index;comment:提醒用户ID"`
	Type         ReminderType      `json:"type" gorm:"type:varchar(32);not null;comment:提醒方式"`
	Title        string            `json:"title" gorm:"type:varchar(255);not null;comment:提醒标题"`
	Content      string            `json:"content" gorm:"type:text;comment:提醒内容"`
	Status       ReminderStatus    `json:"status" gorm:"type:varchar(32);default:pending;index;comment:提醒状态"`
	Frequency    ReminderFrequency `json:"frequency" gorm:"type:varchar(32);default:once;comment:提醒频率"`
	ScheduleTime time.Time         `json:"schedule_time" gorm:"not null;index;comment:计划提醒时间"`
	SentTime     *time.Time        `json:"sent_time" gorm:"comment:实际发送时间"`
	FailReason   string            `json:"fail_reason" gorm:"type:varchar(500);comment:失败原因"`
	RetryCount   int               `json:"retry_count" gorm:"default:0;comment:重试次数"`
	MaxRetries   int               `json:"max_retries" gorm:"default:3;comment:最大重试次数"`
	CreatedAt    time.Time         `json:"created_at" gorm:"index;comment:创建时间"`
	UpdatedAt    time.Time         `json:"updated_at" gorm:"comment:更新时间"`

	Todo Todo `json:"todo" gorm:"foreignKey:TodoID"`
	User User `json:"user" gorm:"foreignKey:UserID"`
}

func (Reminder) TableName() string {
	return "reminders"
}

// ReminderTemplate 提醒模板
type ReminderTemplate struct {
	ID        uint64       `json:"id" gorm:"primaryKey;autoIncrement;comment:模板ID"`
	Name      string       `json:"name" gorm:"type:varchar(100);not null;comment:模板名称"`
	Type      ReminderType `json:"type" gorm:"type:varchar(32);not null;comment:适用的提醒方式"`
	Title     string       `json:"title" gorm:"type:varchar(255);not null;comment:标题模板"`
	Content   string       `json:"content" gorm:"type:text;not null;comment:内容模板"`
	Variables JSONB        `json:"variables" gorm:"type:json;comment:可用变量说明"`
	IsActive  bool         `json:"is_active" gorm:"default:true;comment:是否启用"`
	IsDefault bool         `json:"is_default" gorm:"default:false;comment:是否为默认模板"`
	CreatedBy uint64       `json:"created_by" gorm:"not null;comment:创建人ID"`
	CreatedAt time.Time    `json:"created_at" gorm:"comment:创建时间"`
	UpdatedAt time.Time    `json:"updated_at" gorm:"comment:更新时间"`
}

func (ReminderTemplate) TableName() string {
	return "reminder_templates"
}

// ReminderConfig 提醒配置
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

	User User `json:"user" gorm:"foreignKey:UserID"`
}

func (ReminderConfig) TableName() string {
	return "reminder_configs"
}

// TagDimension 标签维度模型
type TagDimension struct {
	ID          uint64 `json:"id" gorm:"primaryKey;autoIncrement;comment:维度ID"`
	Name        string `json:"name" gorm:"type:varchar(128);not null;uniqueIndex;comment:维度名称"`
	Description string `json:"description" gorm:"type:text;comment:维度描述"`
	SortOrder   int    `json:"sort_order" gorm:"default:0;index;comment:排序顺序"`
	BaseModel

	Tags []Tag `json:"tags,omitempty" gorm:"foreignKey:DimensionID"`
}

func (TagDimension) TableName() string {
	return "tag_dimensions"
}

// Tag 标签模型
type Tag struct {
	ID          uint64 `json:"id" gorm:"primaryKey;autoIncrement;comment:标签ID"`
	DimensionID uint64 `json:"dimension_id" gorm:"not null;index;comment:维度ID"`
	Name        string `json:"name" gorm:"type:varchar(128);not null;comment:标签名称"`
	Color       string `json:"color" gorm:"type:varchar(32);comment:标签颜色"`
	Description string `json:"description" gorm:"type:text;comment:标签描述"`
	SortOrder   int    `json:"sort_order" gorm:"default:0;index;comment:排序顺序"`
	BaseModel

	Dimension TagDimension `json:"dimension,omitempty" gorm:"foreignKey:DimensionID"`
}

func (Tag) TableName() string {
	return "tags"
}

// FollowUpRecord 跟进记录模型
// 注意：该模型对应的 follow_up_records 表是从原 activities 表迁移而来
type FollowUpRecord struct {
	ID                   uint64     `json:"id" gorm:"primaryKey;autoIncrement;comment:跟进记录ID"`
	CustomerID           uint64     `json:"customer_id" gorm:"not null;index;comment:客户ID"`
	UserID               uint64     `json:"user_id" gorm:"not null;index;comment:创建用户ID"`
	Kind                 string     `json:"kind" gorm:"type:varchar(32);default:other;index;comment:跟进记录类型"`
	Title                string     `json:"title" gorm:"type:varchar(255);comment:跟进标题"`
	Content              string     `json:"content" gorm:"type:text;comment:跟进内容"`
	Data                 JSONB      `json:"data" gorm:"type:jsonb;comment:结构化数据"`
	Remark               string     `json:"remark" gorm:"type:text;comment:备注"`
	Duration             *int       `json:"duration" gorm:"comment:持续时长（分钟）"`
	Location             string     `json:"location" gorm:"type:varchar(255);comment:地点"`
	Amount               *float64   `json:"amount" gorm:"type:decimal(15,2);comment:金额"`
	Cost                 *float64   `json:"cost" gorm:"type:decimal(15,2);comment:成本"`
	RelatedTodoID        *uint64    `json:"related_todo_id" gorm:"index;comment:关联的待办ID"`
	ParentRecordID       *uint64    `json:"parent_record_id" gorm:"index;comment:父记录ID"`
	NextFollowTime       *time.Time `json:"next_follow_time" gorm:"index;comment:下次跟进时间"`
	NextFollowContent    string     `json:"next_follow_content" gorm:"type:text;comment:下次跟进内容"`
	Attachments          JSONB      `json:"attachments" gorm:"type:jsonb;comment:附件信息"`
	Photos               JSONB      `json:"photos" gorm:"type:jsonb;comment:照片信息"`
	CustomerSatisfaction *int       `json:"customer_satisfaction" gorm:"comment:客户满意度(1-5)"`
	CustomerFeedback     string     `json:"customer_feedback" gorm:"type:text;comment:客户反馈"`
	BaseModel

	Customer     Customer        `json:"customer" gorm:"foreignKey:CustomerID"`
	User         User            `json:"user" gorm:"foreignKey:UserID"`
	RelatedTodo  *Todo           `json:"related_todo,omitempty" gorm:"foreignKey:RelatedTodoID"`
	ParentRecord *FollowUpRecord `json:"parent_record,omitempty" gorm:"foreignKey:ParentRecordID"`
}

func (FollowUpRecord) TableName() string {
	return "follow_up_records"
}

// GetTimeAgo 获取时间差描述
func (f *FollowUpRecord) GetTimeAgo() string {
	now := time.Now()
	diff := now.Sub(f.CreatedAt)

	if diff < time.Minute {
		return "刚刚"
	} else if diff < time.Hour {
		return fmt.Sprintf("%d分钟前", int(diff.Minutes()))
	} else if diff < 24*time.Hour {
		return fmt.Sprintf("%d小时前", int(diff.Hours()))
	} else {
		return fmt.Sprintf("%d天前", int(diff.Hours()/24))
	}
}

// Group 分组模型
type Group struct {
	ID          uint64 `json:"id" gorm:"primaryKey;autoIncrement;comment:分组ID"`
	Name        string `json:"name" gorm:"type:varchar(128);not null;comment:分组名称"`
	Description string `json:"description" gorm:"type:text;comment:分组描述"`
	CreatedBy   uint64 `json:"created_by" gorm:"not null;index;comment:创建人ID"`
	SortOrder   int    `json:"sort_order" gorm:"default:0;index;comment:排序顺序"`
	BaseModel

	Creator User `json:"creator" gorm:"foreignKey:CreatedBy"`
}

func (Group) TableName() string {
	return "groups"
}
