package main

import (
	"time"

	"github.com/lib/pq"
)

// ============================================================================
// 请求和响应结构体定义
// ============================================================================

// Customer 相关请求响应
type CustomerRequest struct {
	Name         string   `json:"name" binding:"required,min=1,max=100"`
	ContactName  string   `json:"contact_name" binding:"max=50"`
	Phones       []string `json:"phones"`
	Wechats      []string `json:"wechats"`
	Province     string   `json:"province" binding:"max=20"`
	City         string   `json:"city" binding:"max=20"`
	District     string   `json:"district" binding:"max=20"`
	Address      string   `json:"address" binding:"max=200"`
	Products     []string `json:"products"`
	Category     string   `json:"category" binding:"max=50"`
	Tags         []string `json:"tags"`
	State        int      `json:"state" binding:"min=0,max=10"`
	Level        int      `json:"level" binding:"min=0,max=10"`
	Source       string   `json:"source" binding:"max=50"`
	ImportSource string   `json:"import_source" binding:"max=50"`
	Remark       string   `json:"remark" binding:"max=500"`
	SallerName   string   `json:"saller_name" binding:"max=50"`
	Sellers      []int64  `json:"sellers"`
}

type CustomerResponse struct {
	ID           uint     `json:"id"`
	Name         string   `json:"name"`
	Phone        string   `json:"phone"`
	Phones       []string `json:"phones"`
	Wechat       string   `json:"wechat"`
	Wechats      []string `json:"wechats"`
	Seller       string   `json:"seller"`
	Sellers      []int64  `json:"sellers"`
	SallerName   string   `json:"saller_name"`
	Address      string   `json:"address"`
	Province     string   `json:"province"`
	City         string   `json:"city"`
	District     string   `json:"district"`
	Company      string   `json:"company"`
	Products     []string `json:"products"`
	Category     string   `json:"category"`
	Tags         []string `json:"tags"`
	State        int      `json:"state"`
	Level        int      `json:"level"`
	ContactName  string   `json:"contact_name"`
	Source       string   `json:"source"`
	ImportSource string   `json:"import_source"`
	Remark       string   `json:"remark"`
	Organization string   `json:"organization,omitempty"`
	CreatedAt    string   `json:"created_at"`
}

// CustomerToResponse 将客户模型转换为响应
func CustomerToResponse(customer *Customer) *CustomerResponse {
	response := &CustomerResponse{
		ID:           customer.ID,
		Name:         customer.Name,
		Phones:       []string(customer.Phones),
		Wechats:      []string(customer.Wechats),
		Sellers:      []int64(customer.Sellers),
		SallerName:   customer.SallerName,
		Address:      customer.Address,
		Province:     customer.Province,
		City:         customer.City,
		District:     customer.District,
		Products:     []string(customer.Products),
		Category:     customer.Category,
		Tags:         []string(customer.Tags),
		State:        customer.State,
		Level:        customer.Level,
		ContactName:  customer.ContactName,
		Source:       customer.Source,
		ImportSource: customer.ImportSource,
		Remark:       customer.Remark,
		CreatedAt:    customer.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	if len(customer.Phones) > 0 {
		response.Phone = customer.Phones[0]
	}
	if len(customer.Wechats) > 0 {
		response.Wechat = customer.Wechats[0]
	}

	return response
}

// Todo 相关请求响应
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

type TodoResponse struct {
	Todo
	CreatorName      string  `json:"creator_name"`
	ExecutorName     string  `json:"executor_name"`
	CustomerName     string  `json:"customer_name"`
	ReminderUserName *string `json:"reminder_user_name"`
	IsOverdue        bool    `json:"is_overdue"`
	DaysLeft         int     `json:"days_left"`
}

// Activity 相关请求响应
type ActivityCreateRequest struct {
	CustomerID     uint64       `json:"customer_id" binding:"required"`
	Kind           ActivityKind `json:"kind" binding:"required"`
	Title          string       `json:"title" binding:"required,max=255"`
	Content        string       `json:"content" binding:"required"`
	Result         string       `json:"result"`
	Amount         float64      `json:"amount"`
	Cost           float64      `json:"cost"`
	Feedback       string       `json:"feedback"`
	Satisfaction   int          `json:"satisfaction"`
	Remark         string       `json:"remark"`
	Duration       *int         `json:"duration"`
	Location       string       `json:"location"`
	NextFollowTime *time.Time   `json:"next_follow_time"`
	Attachments    JSONB        `json:"attachments"`
}

type ActivityResponse struct {
	Activity
	UserName     string  `json:"user_name"`
	CustomerName string  `json:"customer_name"`
	Content      string  `json:"content"`
	Result       string  `json:"result"`
	Amount       float64 `json:"amount"`
	Cost         float64 `json:"cost"`
	Feedback     string  `json:"feedback"`
	Satisfaction int     `json:"satisfaction"`
	TimeAgo      string  `json:"time_ago"`
}

// User 相关请求响应
type UserResponse struct {
	ID          uint64 `json:"id"`
	Name        string `json:"name"`
	Department  string `json:"department"`
	Position    string `json:"position"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Status      string `json:"status"`
	ManagerName string `json:"manager_name,omitempty"`
	AvatarURL   string `json:"avatar_url"`
}

// UserDetailResponse 用户详情响应（智能判断员工/客户身份）
type UserDetailResponse struct {
	ID           uint64 `json:"id"`
	Name         string `json:"name"`
	DisplayInfo  string `json:"display_info"`  // 显示信息（department+position 或 客户公司名）
	IsEmployee   bool   `json:"is_employee"`   // 是否为员工
	TodayRevenue int    `json:"today_revenue"` // 今日营业额（暂时默认0）
	TodayFollows int    `json:"today_follows"` // 今日跟进数量
	AvatarURL    string `json:"avatar_url"`
}

// DashboardSearchRequest 仪表板搜索请求
type DashboardSearchRequest struct {
	UserID       uint64 `json:"user_id" binding:"required"`       // 当前用户ID
	TimeFilter   string `json:"time_filter" binding:"required"`   // 时间筛选条件
	StatusFilter string `json:"status_filter" binding:"required"` // 状态筛选条件
	ShowAll      bool   `json:"show_all"`                         // 是否查看全部（忽略分页）
	Page         int    `json:"page"`                             // 页码，默认1
	PageSize     int    `json:"page_size"`                        // 页大小，默认20
}

// DashboardSearchResponse 仪表板搜索响应
type DashboardSearchResponse struct {
	CustomerID    uint64   `json:"customer_id"`     // 客户ID
	ContactName   string   `json:"contact_name"`    // 联系人姓名
	CustomerName  string   `json:"customer_name"`   // 客户店铺名
	Tags          []string `json:"tags"`            // 客户标签
	TodoContents  string   `json:"todo_contents"`   // 待办内容（多个用逗号分隔）
	TodoCount     int      `json:"todo_count"`      // 待办数量
	PlannedTime   string   `json:"planned_time"`    // 计划时间
	LastCallTime  string   `json:"last_call_time"`  // 最后联系时间
	LastOrderTime string   `json:"last_order_time"` // 最后下单时间
}

// CustomerSearchResponse 客户搜索响应
type CustomerSearchResponse struct {
	ID          uint64   `json:"id"`           // 客户ID
	Name        string   `json:"name"`         // 客户店铺名
	ContactName string   `json:"contact_name"` // 联系人姓名
	Phone       string   `json:"phone"`        // 主要电话
	Category    string   `json:"category"`     // 客户分类
	Tags        []string `json:"tags"`         // 客户标签
	SystemTags  []int64  `json:"system_tags"`  // 系统标签ID
	Province    string   `json:"province"`     // 省份
	City        string   `json:"city"`         // 城市
	State       int      `json:"state"`        // 客户状态
	Level       int      `json:"level"`        // 客户级别
}

// Tag 相关请求响应
type TagCreateRequest struct {
	DimensionID uint64 `json:"dimension_id" binding:"required"`
	Name        string `json:"name" binding:"required,max=128"`
	Color       string `json:"color"`
	Description string `json:"description"`
	SortOrder   int    `json:"sort_order"`
}

type TagResponse struct {
	ID            uint64 `json:"id"`
	DimensionID   uint64 `json:"dimension_id"`
	DimensionName string `json:"dimension_name"`
	Name          string `json:"name"`
	Color         string `json:"color"`
	Description   string `json:"description"`
	SortOrder     int    `json:"sort_order"`
}

// Reminder 相关请求响应
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

type ReminderResponse struct {
	Reminder
	TodoTitle    string `json:"todo_title"`
	UserName     string `json:"user_name"`
	CustomerName string `json:"customer_name"`
}

// 类型转换辅助函数
func convertJSONBToStringArray(jsonb JSONB) pq.StringArray {
	if jsonb == nil {
		return pq.StringArray{}
	}
	var result pq.StringArray
	for _, v := range jsonb {
		if str, ok := v.(string); ok {
			result = append(result, str)
		}
	}
	return result
}

func convertStringArrayToJSONB(arr pq.StringArray) JSONB {
	result := make(JSONB)
	for i, v := range arr {
		result[string(rune(i))] = v
	}
	return result
}

func convertJSONBToInt64Array(jsonb JSONB) pq.Int64Array {
	if jsonb == nil {
		return pq.Int64Array{}
	}
	var result pq.Int64Array
	for _, v := range jsonb {
		if num, ok := v.(int64); ok {
			result = append(result, num)
		} else if num, ok := v.(float64); ok {
			result = append(result, int64(num))
		}
	}
	return result
}

func convertInt64ArrayToJSONB(arr pq.Int64Array) JSONB {
	result := make(JSONB)
	for i, v := range arr {
		result[string(rune(i))] = v
	}
	return result
}

// BaseDTO 基础 DTO 结构
type BaseDTO struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// AuditDTO 审计 DTO 结构
type AuditDTO struct {
	BaseDTO
	IsDeleted bool       `gorm:"default:false;index" json:"is_deleted"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}

// CustomerDTO 客户数据传输对象
type CustomerDTO struct {
	AuditDTO
	Name        string `gorm:"size:100;not null;index" json:"name"`
	ContactName string `gorm:"size:50" json:"contact_name"`
	Phones      JSONB  `gorm:"type:jsonb" json:"phones"`
	Wechats     JSONB  `gorm:"type:jsonb" json:"wechats"`
	Address     string `gorm:"size:200" json:"address"`
	Province    string `gorm:"size:50" json:"province"`
	City        string `gorm:"size:50" json:"city"`
	District    string `gorm:"size:50" json:"district"`
	Favors      JSONB  `gorm:"type:jsonb" json:"favors"`
	Remark      string `gorm:"type:text" json:"remark"`
	SystemTags  JSONB  `gorm:"type:jsonb" json:"system_tags"`
}

// TableName 指定表名
func (CustomerDTO) TableName() string {
	return "customers"
}

// ToModel 转换为业务模型
func (dto *CustomerDTO) ToModel() *Customer {
	return &Customer{
		ID:          uint(dto.ID),
		CreatedAt:   dto.CreatedAt,
		UpdatedAt:   dto.UpdatedAt,
		Name:        dto.Name,
		ContactName: dto.ContactName,
		Phones:      convertJSONBToStringArray(dto.Phones),
		Wechats:     convertJSONBToStringArray(dto.Wechats),
		Address:     dto.Address,
		Province:    dto.Province,
		City:        dto.City,
		District:    dto.District,
		Favors:      dto.Favors,
		Remark:      dto.Remark,
		SystemTags:  convertJSONBToInt64Array(dto.SystemTags),
	}
}

// FromModel 从业务模型转换
func (dto *CustomerDTO) FromModel(customer *Customer) {
	dto.ID = uint64(customer.ID)
	dto.CreatedAt = customer.CreatedAt
	dto.UpdatedAt = customer.UpdatedAt
	dto.Name = customer.Name
	dto.ContactName = customer.ContactName
	dto.Phones = convertStringArrayToJSONB(customer.Phones)
	dto.Wechats = convertStringArrayToJSONB(customer.Wechats)
	dto.Address = customer.Address
	dto.Province = customer.Province
	dto.City = customer.City
	dto.District = customer.District
	dto.Favors = customer.Favors
	dto.Remark = customer.Remark
	dto.SystemTags = convertInt64ArrayToJSONB(customer.SystemTags)
}

// UserDTO 用户数据传输对象
type UserDTO struct {
	AuditDTO
	Name               string     `gorm:"size:256;not null;index" json:"name"`
	ManagerID          *uint64    `gorm:"index" json:"manager_id"`
	Email              string     `gorm:"size:256" json:"email"`
	Phone              string     `gorm:"size:32" json:"phone"`
	Department         string     `gorm:"size:128" json:"department"`
	DepartmentLeaderID *uint64    `gorm:"index" json:"department_leader_id"`
	Position           string     `gorm:"size:128" json:"position"`
	WechatWorkID       string     `gorm:"size:128" json:"wechat_work_id"`
	WechatID           string     `gorm:"size:128" json:"wechat_id"`
	Status             string     `gorm:"size:32;default:'active'" json:"status"`
	AvatarURL          string     `gorm:"size:512" json:"avatar_url"`
	LastLoginAt        *time.Time `json:"last_login_at"`
	Manager            *UserDTO   `gorm:"foreignKey:ManagerID" json:"manager,omitempty"`
	DepartmentLeader   *UserDTO   `gorm:"foreignKey:DepartmentLeaderID" json:"department_leader,omitempty"`
}

// TableName 指定表名
func (UserDTO) TableName() string {
	return "users"
}

// ToModel 转换为业务模型
func (dto *UserDTO) ToModel() *User {
	user := &User{
		ID:                 dto.ID,
		Name:               dto.Name,
		ManagerID:          dto.ManagerID,
		Email:              dto.Email,
		Phone:              dto.Phone,
		Department:         dto.Department,
		DepartmentLeaderID: dto.DepartmentLeaderID,
		Position:           dto.Position,
		WechatWorkID:       dto.WechatWorkID,
		WechatID:           dto.WechatID,
		Status:             dto.Status,
		AvatarURL:          dto.AvatarURL,
		LastLoginAt:        dto.LastLoginAt,
		BaseModel: BaseModel{
			CreatedAt: dto.CreatedAt,
			UpdatedAt: dto.UpdatedAt,
			DeletedAt: dto.DeletedAt,
			IsDeleted: dto.IsDeleted,
		},
	}
	if dto.Manager != nil {
		user.Manager = dto.Manager.ToModel()
	}
	if dto.DepartmentLeader != nil {
		user.DepartmentLeader = dto.DepartmentLeader.ToModel()
	}
	return user
}

// FromModel 从业务模型转换
func (dto *UserDTO) FromModel(user *User) {
	dto.ID = user.ID
	dto.Name = user.Name
	dto.ManagerID = user.ManagerID
	dto.Email = user.Email
	dto.Phone = user.Phone
	dto.Department = user.Department
	dto.DepartmentLeaderID = user.DepartmentLeaderID
	dto.Position = user.Position
	dto.WechatWorkID = user.WechatWorkID
	dto.WechatID = user.WechatID
	dto.Status = user.Status
	dto.AvatarURL = user.AvatarURL
	dto.LastLoginAt = user.LastLoginAt
	dto.CreatedAt = user.BaseModel.CreatedAt
	dto.UpdatedAt = user.BaseModel.UpdatedAt
	dto.DeletedAt = user.BaseModel.DeletedAt
	dto.IsDeleted = user.BaseModel.IsDeleted
	if user.Manager != nil {
		dto.Manager = &UserDTO{}
		dto.Manager.FromModel(user.Manager)
	}
	if user.DepartmentLeader != nil {
		dto.DepartmentLeader = &UserDTO{}
		dto.DepartmentLeader.FromModel(user.DepartmentLeader)
	}
}

// TodoDTO 待办数据传输对象
type TodoDTO struct {
	AuditDTO
	Title          string       `gorm:"size:200;not null" json:"title"`
	Description    string       `gorm:"type:text" json:"description"`
	Status         TodoStatus   `gorm:"size:20;default:'pending'" json:"status"`
	Priority       Priority     `gorm:"size:20;default:'medium'" json:"priority"`
	DueDate        *time.Time   `json:"due_date"`
	CompletedAt    *time.Time   `json:"completed_at"`
	CustomerID     *uint64      `gorm:"index" json:"customer_id"`
	Customer       *CustomerDTO `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	CreatorID      uint64       `gorm:"not null;index" json:"creator_id"`
	Creator        *UserDTO     `gorm:"foreignKey:CreatorID" json:"creator,omitempty"`
	ExecutorID     *uint64      `gorm:"index" json:"executor_id"`
	Executor       *UserDTO     `gorm:"foreignKey:ExecutorID" json:"executor,omitempty"`
	ReminderUserID *uint64      `gorm:"index" json:"reminder_user_id"`
	ReminderUser   *UserDTO     `gorm:"foreignKey:ReminderUserID" json:"reminder_user,omitempty"`
	ReminderTime   *time.Time   `json:"reminder_time"`
}

// TableName 指定表名
func (TodoDTO) TableName() string {
	return "todos"
}

// ToModel 转换为业务模型
func (dto *TodoDTO) ToModel() *Todo {
	todo := &Todo{
		ID:             uint64(dto.ID),
		Title:          dto.Title,
		Content:        dto.Description,
		Status:         dto.Status,
		Priority:       dto.Priority,
		CreatorID:      uint64(dto.CreatorID),
		ReminderUserID: dto.ReminderUserID,
		ReminderTime:   dto.ReminderTime,
		BaseModel: BaseModel{
			CreatedAt: dto.CreatedAt,
			UpdatedAt: dto.UpdatedAt,
			DeletedAt: dto.DeletedAt,
			IsDeleted: dto.IsDeleted,
		},
	}
	if dto.DueDate != nil {
		todo.PlannedTime = *dto.DueDate
	}
	if dto.CompletedAt != nil {
		todo.CompletedTime = dto.CompletedAt
	}
	if dto.CustomerID != nil {
		todo.CustomerID = uint64(*dto.CustomerID)
	}
	if dto.ExecutorID != nil {
		todo.ExecutorID = uint64(*dto.ExecutorID)
	}
	// 关联对象在查询时通过 GORM 自动填充，这里不需要手动设置
	return todo
}

// FromModel 从业务模型转换
func (dto *TodoDTO) FromModel(todo *Todo) {
	dto.ID = uint64(todo.ID)
	dto.CreatedAt = todo.BaseModel.CreatedAt
	dto.UpdatedAt = todo.BaseModel.UpdatedAt
	dto.IsDeleted = todo.BaseModel.IsDeleted
	dto.DeletedAt = todo.BaseModel.DeletedAt
	dto.Title = todo.Title
	dto.Description = todo.Content
	dto.Status = todo.Status
	dto.Priority = todo.Priority
	dto.DueDate = &todo.PlannedTime
	dto.CompletedAt = todo.CompletedTime
	dto.CreatorID = uint64(todo.CreatorID)
	dto.ReminderUserID = todo.ReminderUserID
	dto.ReminderTime = todo.ReminderTime
	if todo.CustomerID != 0 {
		customerID := uint64(todo.CustomerID)
		dto.CustomerID = &customerID
	}
	if todo.ExecutorID != 0 {
		executorID := uint64(todo.ExecutorID)
		dto.ExecutorID = &executorID
	}
	// 处理关联对象
	dto.Customer = &CustomerDTO{}
	dto.Customer.FromModel(&todo.Customer)
	dto.Creator = &UserDTO{}
	dto.Creator.FromModel(&todo.Creator)
	dto.Executor = &UserDTO{}
	dto.Executor.FromModel(&todo.Executor)
	if todo.ReminderUser != nil {
		dto.ReminderUser = &UserDTO{}
		dto.ReminderUser.FromModel(todo.ReminderUser)
	}
}

// TodoLogDTO 待办日志数据传输对象
type TodoLogDTO struct {
	BaseDTO
	TodoID     uint64     `gorm:"not null;index" json:"todo_id"`
	OperatorID uint64     `gorm:"not null;index" json:"operator_id"`
	Action     ActionType `gorm:"size:50;not null" json:"action"`
	OldData    JSONB      `gorm:"type:jsonb" json:"old_data"`
	NewData    JSONB      `gorm:"type:jsonb" json:"new_data"`
}

// TableName 指定表名
func (TodoLogDTO) TableName() string {
	return "todo_logs"
}

// ToModel 转换为业务模型
func (dto *TodoLogDTO) ToModel() *TodoLog {
	return &TodoLog{
		ID:         uint64(dto.ID),
		TodoID:     uint64(dto.TodoID),
		OperatorID: uint64(dto.OperatorID),
		Action:     dto.Action,
		OldData:    dto.OldData,
		NewData:    dto.NewData,
		CreatedAt:  dto.CreatedAt,
	}
}

// FromModel 从业务模型转换
func (dto *TodoLogDTO) FromModel(log *TodoLog) {
	dto.ID = uint64(log.ID)
	dto.CreatedAt = log.CreatedAt
	// TodoLog 没有 UpdatedAt 字段
	dto.TodoID = uint64(log.TodoID)
	dto.OperatorID = uint64(log.OperatorID)
	dto.Action = log.Action
	dto.OldData = log.OldData
	dto.NewData = log.NewData
}

// ActivityDTO 跟进记录数据传输对象
type ActivityDTO struct {
	AuditDTO
	Kind         ActivityKind `gorm:"size:50;not null;index" json:"kind"`
	Title        string       `gorm:"size:200;not null" json:"title"`
	Content      string       `gorm:"type:text" json:"content"`
	CustomerID   uint64       `gorm:"not null;index" json:"customer_id"`
	Customer     *CustomerDTO `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	UserID       uint64       `gorm:"not null;index" json:"user_id"`
	User         *UserDTO     `gorm:"foreignKey:UserID" json:"user,omitempty"`
	FollowUpDate *time.Time   `gorm:"index" json:"follow_up_date"`
	Metadata     JSONB        `gorm:"type:jsonb" json:"metadata"`
}

// TableName 指定表名
func (ActivityDTO) TableName() string {
	return "activities"
}

// ToModel 转换为业务模型
func (dto *ActivityDTO) ToModel() *Activity {
	activity := &Activity{
		ID:             uint64(dto.ID),
		CustomerID:     uint64(dto.CustomerID),
		UserID:         uint64(dto.UserID),
		Kind:           dto.Kind,
		Title:          dto.Title,
		Remark:         dto.Content,
		Data:           dto.Metadata,
		NextFollowTime: dto.FollowUpDate,
		BaseModel: BaseModel{
			CreatedAt: dto.CreatedAt,
			UpdatedAt: dto.UpdatedAt,
			DeletedAt: dto.DeletedAt,
			IsDeleted: dto.IsDeleted,
		},
	}
	return activity
}

// FromModel 从业务模型转换
func (dto *ActivityDTO) FromModel(activity *Activity) {
	dto.ID = uint64(activity.ID)
	dto.CreatedAt = activity.CreatedAt
	dto.UpdatedAt = activity.UpdatedAt
	dto.IsDeleted = activity.IsDeleted
	dto.DeletedAt = activity.DeletedAt
	dto.Kind = activity.Kind
	dto.Title = activity.Title
	dto.Content = activity.Remark
	dto.CustomerID = uint64(activity.CustomerID)
	dto.UserID = uint64(activity.UserID)
	dto.FollowUpDate = activity.NextFollowTime
	dto.Metadata = activity.Data
	dto.Customer = &CustomerDTO{}
	dto.Customer.FromModel(&activity.Customer)
	dto.User = &UserDTO{}
	dto.User.FromModel(&activity.User)
}

// ReminderDTO 提醒数据传输对象
type ReminderDTO struct {
	BaseDTO
	TodoID       uint64            `gorm:"not null;index" json:"todo_id"`
	Todo         *TodoDTO          `gorm:"foreignKey:TodoID" json:"todo,omitempty"`
	UserID       uint64            `gorm:"not null;index" json:"user_id"`
	User         *UserDTO          `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Type         ReminderType      `gorm:"size:50;not null" json:"type"`
	Title        string            `gorm:"size:200;not null" json:"title"`
	Content      string            `gorm:"type:text" json:"content"`
	Status       ReminderStatus    `gorm:"size:20;default:'pending'" json:"status"`
	Frequency    ReminderFrequency `gorm:"size:20;default:'once'" json:"frequency"`
	ScheduleTime time.Time         `gorm:"not null;index" json:"schedule_time"`
	SentTime     *time.Time        `json:"sent_time"`
	FailReason   string            `gorm:"type:text" json:"fail_reason"`
	RetryCount   int               `gorm:"default:0" json:"retry_count"`
	MaxRetries   int               `gorm:"default:3" json:"max_retries"`
}

// TableName 指定表名
func (ReminderDTO) TableName() string {
	return "reminders"
}

// ToModel 转换为业务模型
func (dto *ReminderDTO) ToModel() *Reminder {
	reminder := &Reminder{
		ID:           uint64(dto.ID),
		TodoID:       uint64(dto.TodoID),
		UserID:       uint64(dto.UserID),
		Type:         dto.Type,
		Title:        dto.Title,
		Content:      dto.Content,
		Status:       dto.Status,
		Frequency:    dto.Frequency,
		ScheduleTime: dto.ScheduleTime,
		SentTime:     dto.SentTime,
		FailReason:   dto.FailReason,
		RetryCount:   dto.RetryCount,
		MaxRetries:   dto.MaxRetries,
		CreatedAt:    dto.CreatedAt,
		UpdatedAt:    dto.UpdatedAt,
	}
	return reminder
}

// FromModel 从业务模型转换
func (dto *ReminderDTO) FromModel(reminder *Reminder) {
	dto.ID = uint64(reminder.ID)
	dto.TodoID = uint64(reminder.TodoID)
	dto.UserID = uint64(reminder.UserID)
	dto.Type = reminder.Type
	dto.Title = reminder.Title
	dto.Content = reminder.Content
	dto.Status = reminder.Status
	dto.Frequency = reminder.Frequency
	dto.ScheduleTime = reminder.ScheduleTime
	dto.SentTime = reminder.SentTime
	dto.FailReason = reminder.FailReason
	dto.RetryCount = reminder.RetryCount
	dto.MaxRetries = reminder.MaxRetries
	dto.CreatedAt = reminder.CreatedAt
	dto.UpdatedAt = reminder.UpdatedAt
}

// ReminderTemplateDTO 提醒模板数据传输对象
type ReminderTemplateDTO struct {
	BaseDTO
	Type      ReminderType `gorm:"size:50;not null;index" json:"type"`
	Title     string       `gorm:"size:200;not null" json:"title"`
	Content   string       `gorm:"type:text;not null" json:"content"`
	IsDefault bool         `gorm:"default:false" json:"is_default"`
	IsActive  bool         `gorm:"default:true" json:"is_active"`
}

// TableName 指定表名
func (ReminderTemplateDTO) TableName() string {
	return "reminder_templates"
}

// ToModel 转换为业务模型
func (dto *ReminderTemplateDTO) ToModel() *ReminderTemplate {
	return &ReminderTemplate{
		ID:        uint64(dto.ID),
		CreatedAt: dto.CreatedAt,
		UpdatedAt: dto.UpdatedAt,
		Type:      dto.Type,
		Title:     dto.Title,
		Content:   dto.Content,
		IsDefault: dto.IsDefault,
		IsActive:  dto.IsActive,
	}
}

// FromModel 从业务模型转换
func (dto *ReminderTemplateDTO) FromModel(template *ReminderTemplate) {
	dto.ID = uint64(template.ID)
	dto.CreatedAt = template.CreatedAt
	dto.UpdatedAt = template.UpdatedAt
	dto.Type = template.Type
	dto.Title = template.Title
	dto.Content = template.Content
	dto.IsDefault = template.IsDefault
	dto.IsActive = template.IsActive
}

// ReminderConfigDTO 提醒配置数据传输对象
type ReminderConfigDTO struct {
	BaseDTO
	UserID          uint64            `gorm:"not null;uniqueIndex" json:"user_id"`
	User            *UserDTO          `gorm:"foreignKey:UserID" json:"user,omitempty"`
	EmailEnabled    bool              `gorm:"default:true" json:"email_enabled"`
	SmsEnabled      bool              `gorm:"default:false" json:"sms_enabled"`
	PushEnabled     bool              `gorm:"default:true" json:"push_enabled"`
	AdvanceMinutes  int               `gorm:"default:15" json:"advance_minutes"`
	Frequency       ReminderFrequency `gorm:"size:20;default:'once'" json:"frequency"`
	QuietHoursStart string            `gorm:"size:5;default:'22:00'" json:"quiet_hours_start"`
	QuietHoursEnd   string            `gorm:"size:5;default:'08:00'" json:"quiet_hours_end"`
	WeekendEnabled  bool              `gorm:"default:false" json:"weekend_enabled"`
}

// TableName 指定表名
func (ReminderConfigDTO) TableName() string {
	return "reminder_configs"
}

// ToModel 转换为业务模型
func (dto *ReminderConfigDTO) ToModel() *ReminderConfig {
	config := &ReminderConfig{
		ID:                     uint64(dto.ID),
		CreatedAt:              dto.CreatedAt,
		UpdatedAt:              dto.UpdatedAt,
		UserID:                 uint64(dto.UserID),
		EnableWechat:           dto.EmailEnabled,
		EnableEnterpriseWechat: dto.PushEnabled,
		DefaultAdvanceMinutes:  dto.AdvanceMinutes,
		QuietStartTime:         dto.QuietHoursStart,
		QuietEndTime:           dto.QuietHoursEnd,
	}
	return config
}

// FromModel 从业务模型转换
func (dto *ReminderConfigDTO) FromModel(config *ReminderConfig) {
	dto.ID = uint64(config.ID)
	dto.CreatedAt = config.CreatedAt
	dto.UpdatedAt = config.UpdatedAt
	dto.UserID = uint64(config.UserID)
	dto.EmailEnabled = config.EnableWechat
	dto.PushEnabled = config.EnableEnterpriseWechat
	dto.AdvanceMinutes = config.DefaultAdvanceMinutes
	dto.QuietHoursStart = config.QuietStartTime
	dto.QuietHoursEnd = config.QuietEndTime
}

// TagDimensionDTO 标签维度数据传输对象
type TagDimensionDTO struct {
	AuditDTO
	Name        string   `gorm:"size:50;not null;uniqueIndex" json:"name"`
	Description string   `gorm:"type:text" json:"description"`
	SortOrder   int      `gorm:"default:0" json:"sort_order"`
	Tags        []TagDTO `gorm:"foreignKey:DimensionID" json:"tags,omitempty"`
}

// TableName 指定表名
func (TagDimensionDTO) TableName() string {
	return "tag_dimensions"
}

// ToModel 转换为业务模型
func (dto *TagDimensionDTO) ToModel() *TagDimension {
	dimension := &TagDimension{
		ID:          uint64(dto.ID),
		Name:        dto.Name,
		Description: dto.Description,
		SortOrder:   dto.SortOrder,
	}
	if len(dto.Tags) > 0 {
		dimension.Tags = make([]Tag, len(dto.Tags))
		for i, tagDTO := range dto.Tags {
			dimension.Tags[i] = *tagDTO.ToModel()
		}
	}
	return dimension
}

// FromModel 从业务模型转换
func (dto *TagDimensionDTO) FromModel(dimension *TagDimension) {
	dto.ID = uint64(dimension.ID)
	dto.CreatedAt = dimension.CreatedAt
	dto.UpdatedAt = dimension.UpdatedAt
	dto.IsDeleted = dimension.IsDeleted
	dto.DeletedAt = dimension.DeletedAt
	dto.Name = dimension.Name
	dto.Description = dimension.Description
	dto.SortOrder = dimension.SortOrder
	if len(dimension.Tags) > 0 {
		dto.Tags = make([]TagDTO, len(dimension.Tags))
		for i, tag := range dimension.Tags {
			dto.Tags[i].FromModel(&tag)
		}
	}
}

// TagDTO 标签数据传输对象
type TagDTO struct {
	AuditDTO
	DimensionID uint64           `gorm:"not null;index" json:"dimension_id"`
	Dimension   *TagDimensionDTO `gorm:"foreignKey:DimensionID" json:"dimension,omitempty"`
	Name        string           `gorm:"size:50;not null" json:"name"`
	Color       string           `gorm:"size:7;default:'#007bff'" json:"color"`
	Description string           `gorm:"type:text" json:"description"`
	SortOrder   int              `gorm:"default:0" json:"sort_order"`
}

// TableName 指定表名
func (TagDTO) TableName() string {
	return "tags"
}

// ToModel 转换为业务模型
func (dto *TagDTO) ToModel() *Tag {
	tag := &Tag{
		ID:          uint64(dto.ID),
		DimensionID: uint64(dto.DimensionID),
		Name:        dto.Name,
		Color:       dto.Color,
		Description: dto.Description,
		SortOrder:   dto.SortOrder,
		BaseModel: BaseModel{
			CreatedAt: dto.CreatedAt,
			UpdatedAt: dto.UpdatedAt,
			DeletedAt: dto.DeletedAt,
			IsDeleted: dto.IsDeleted,
		},
	}
	return tag
}

// FromModel 从业务模型转换
func (dto *TagDTO) FromModel(tag *Tag) {
	dto.ID = uint64(tag.ID)
	dto.CreatedAt = tag.CreatedAt
	dto.UpdatedAt = tag.UpdatedAt
	dto.IsDeleted = tag.IsDeleted
	dto.DeletedAt = tag.DeletedAt
	dto.DimensionID = uint64(tag.DimensionID)
	dto.Name = tag.Name
	dto.Color = tag.Color
	dto.Description = tag.Description
	dto.SortOrder = tag.SortOrder
}
