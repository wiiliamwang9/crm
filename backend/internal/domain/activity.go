package domain

import (
	"encoding/json"
	"fmt"
	"time"
	"unicode/utf8"
)

// ActivityKind 跟进记录类型枚举
type ActivityKind string

const (
	ActivityKindCall      ActivityKind = "call"      // 电话沟通
	ActivityKindVisit     ActivityKind = "visit"     // 实地拜访
	ActivityKindEmail     ActivityKind = "email"     // 邮件
	ActivityKindWechat    ActivityKind = "wechat"    // 微信沟通
	ActivityKindMeeting   ActivityKind = "meeting"   // 会议洽谈
	ActivityKindOrder     ActivityKind = "order"     // 下单记录
	ActivityKindSample    ActivityKind = "sample"    // 发样记录
	ActivityKindFeedback  ActivityKind = "feedback"  // 客户反馈
	ActivityKindComplaint ActivityKind = "complaint" // 客户投诉
	ActivityKindPayment   ActivityKind = "payment"   // 付款记录
	ActivityKindOther     ActivityKind = "other"     // 其他
)

// Activity 跟进记录模型（基于现有activities表）
type Activity struct {
	ID             uint64       `json:"id" gorm:"primaryKey;autoIncrement;comment:记录ID"`
	CustomerID     uint64       `json:"customer_id" gorm:"not null;index;comment:客户ID"`
	UserID         uint64       `json:"user_id" gorm:"not null;index;comment:创建用户ID"`
	Kind           ActivityKind `json:"kind" gorm:"type:activity_kind;default:other;index;comment:跟进类型"`
	Title          string       `json:"title" gorm:"type:varchar(255);comment:跟进标题"`
	Data           JSONB        `json:"data" gorm:"type:jsonb;comment:结构化数据"`
	Remark         string       `json:"remark" gorm:"type:text;comment:备注"`
	Duration       *int         `json:"duration" gorm:"comment:持续时长（分钟）"`
	Location       string       `json:"location" gorm:"type:varchar(255);comment:地点"`
	NextFollowTime *time.Time   `json:"next_follow_time" gorm:"index;comment:下次跟进时间"`
	Attachments    JSONB        `json:"attachments" gorm:"type:jsonb;comment:附件信息"`
	CreatedAt      time.Time    `json:"created_at" gorm:"index;comment:创建时间"`
	UpdatedAt      time.Time    `json:"updated_at" gorm:"comment:更新时间"`
	DeletedAt      *time.Time   `json:"deleted_at" gorm:"comment:删除时间"`
	IsDeleted      bool         `json:"is_deleted" gorm:"default:false;index;comment:是否删除"`

	// 关联关系
	Customer Customer `json:"customer" gorm:"foreignKey:CustomerID"`
	User     User     `json:"user" gorm:"foreignKey:UserID"`
}

// 扩展的数据结构 - 存储在Data字段中

// ActivityData 基础活动数据
type ActivityData struct {
	Content      string  `json:"content"`      // 跟进内容
	Result       string  `json:"result"`       // 跟进结果
	Amount       float64 `json:"amount"`       // 金额
	Cost         float64 `json:"cost"`         // 成本
	Feedback     string  `json:"feedback"`     // 客户反馈
	Satisfaction int     `json:"satisfaction"` // 客户满意度(1-5)
}

// ProductItem 产品明细
type ProductItem struct {
	Name      string  `json:"name"`       // 产品名称
	Quantity  int     `json:"quantity"`   // 数量
	UnitPrice float64 `json:"unit_price"` // 单价
}

// SampleProduct 样品明细
type SampleProduct struct {
	Name     string  `json:"name"`     // 产品名称
	Quantity string  `json:"quantity"` // 数量(如100g)
	Cost     float64 `json:"cost"`     // 成本
}

// OrderActivityData 下单记录数据
type OrderActivityData struct {
	ActivityData
	OrderNumber     string        `json:"order_number"`     // 订单号
	ProductItems    []ProductItem `json:"product_items"`    // 产品明细
	PaymentMethod   string        `json:"payment_method"`   // 支付方式
	DeliveryAddress string        `json:"delivery_address"` // 配送地址
}

// SampleActivityData 发样记录数据
type SampleActivityData struct {
	ActivityData
	SampleProducts       []SampleProduct `json:"sample_products"`        // 样品明细
	ShippingMethod       string          `json:"shipping_method"`        // 配送方式
	TrackingNumber       string          `json:"tracking_number"`        // 快递单号
	ExpectedFeedbackDate string          `json:"expected_feedback_date"` // 预期反馈日期
}

// ActivityCreateRequest 创建跟进记录请求
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
	// 待办事项相关
	CreateTodo      bool       `json:"create_todo"`       // 是否创建待办事项
	TodoPlannedTime *time.Time `json:"todo_planned_time"` // 待办计划时间
	TodoContent     string     `json:"todo_content"`      // 待办内容
	TodoExecutorID  *uint64    `json:"todo_executor_id"`  // 待办执行人ID
}

// ActivityUpdateRequest 更新跟进记录请求
type ActivityUpdateRequest struct {
	Title          *string    `json:"title"`
	Content        *string    `json:"content"`
	Result         *string    `json:"result"`
	Amount         *float64   `json:"amount"`
	Cost           *float64   `json:"cost"`
	Feedback       *string    `json:"feedback"`
	Satisfaction   *int       `json:"satisfaction"`
	Remark         *string    `json:"remark"`
	Duration       *int       `json:"duration"`
	Location       *string    `json:"location"`
	NextFollowTime *time.Time `json:"next_follow_time"`
	Attachments    JSONB      `json:"attachments"`
}

// ActivityQueryRequest 查询跟进记录请求
type ActivityQueryRequest struct {
	CustomerID *uint64       `json:"customer_id" form:"customer_id"`
	UserID     *uint64       `json:"user_id" form:"user_id"`
	Kind       *ActivityKind `json:"kind" form:"kind"`
	StartDate  *time.Time    `json:"start_date" form:"start_date"`
	EndDate    *time.Time    `json:"end_date" form:"end_date"`
	Keyword    string        `json:"keyword" form:"keyword"`
	Page       int           `json:"page" form:"page"`
	PageSize   int           `json:"page_size" form:"page_size"`
}

// ActivityResponse 跟进记录响应
type ActivityResponse struct {
	Activity
	UserName     string  `json:"user_name"`
	CustomerName string  `json:"customer_name"`
	Content      string  `json:"content"`      // 从Data中提取
	Result       string  `json:"result"`       // 从Data中提取
	Amount       float64 `json:"amount"`       // 从Data中提取
	Cost         float64 `json:"cost"`         // 从Data中提取
	Feedback     string  `json:"feedback"`     // 从Data中提取
	Satisfaction int     `json:"satisfaction"` // 从Data中提取
	TimeAgo      string  `json:"time_ago"`     // 时间显示
}

// TableName 指定表名
func (Activity) TableName() string {
	return "activities"
}

// GetKindDisplayName 获取跟进类型显示名称
func (a *Activity) GetKindDisplayName() string {
	kindNames := map[ActivityKind]string{
		ActivityKindCall:      "电话沟通",
		ActivityKindVisit:     "实地拜访",
		ActivityKindEmail:     "邮件",
		ActivityKindWechat:    "微信沟通",
		ActivityKindMeeting:   "会议洽谈",
		ActivityKindOrder:     "下单记录",
		ActivityKindSample:    "发样记录",
		ActivityKindFeedback:  "客户反馈",
		ActivityKindComplaint: "客户投诉",
		ActivityKindPayment:   "付款记录",
		ActivityKindOther:     "其他",
	}
	if name, ok := kindNames[a.Kind]; ok {
		return name
	}
	return string(a.Kind)
}

// GetDataStruct 获取解析后的数据结构
func (a *Activity) GetDataStruct() *ActivityData {
	if a.Data == nil {
		return &ActivityData{}
	}

	dataBytes, err := json.Marshal(a.Data)
	if err != nil {
		return &ActivityData{}
	}

	var data ActivityData
	if err := json.Unmarshal(dataBytes, &data); err != nil {
		return &ActivityData{}
	}

	return &data
}

// SetDataStruct 设置数据结构到Data字段
func (a *Activity) SetDataStruct(data interface{}) error {
	// 使用自定义的JSON编码器，确保UTF-8字符正确处理
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// 验证生成的JSON是否为有效的UTF-8
	if !utf8.Valid(jsonData) {
		return fmt.Errorf("生成的JSON包含无效的UTF-8字符")
	}

	var dataMap map[string]interface{}
	if err := json.Unmarshal(jsonData, &dataMap); err != nil {
		return err
	}

	a.Data = JSONB(dataMap)
	return nil
}

// GetTimeAgo 获取时间显示
func (a *Activity) GetTimeAgo() string {
	now := time.Now()
	duration := now.Sub(a.CreatedAt)

	if duration.Hours() < 1 {
		return "刚刚"
	} else if duration.Hours() < 24 {
		return "今天"
	} else if duration.Hours() < 48 {
		return "昨天"
	} else if duration.Hours() < 72 {
		return "2天前"
	} else {
		days := int(duration.Hours() / 24)
		return string(rune(days)) + "天前"
	}
}

// HasFeedback 检查是否有反馈
func (a *Activity) HasFeedback() bool {
	data := a.GetDataStruct()
	return data.Feedback != ""
}
