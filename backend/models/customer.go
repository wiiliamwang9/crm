package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/lib/pq"
)

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
		*j = nil // 对于无法处理的类型，设为nil而不是报错
		return nil
	}

	// 处理空或无效数据
	if len(bytes) == 0 {
		*j = nil
		return nil
	}

	// 如果是单个对象，则包装成数组
	var raw json.RawMessage
	if err := json.Unmarshal(bytes, &raw); err != nil {
		// JSON解析失败，设为nil而不是报错
		*j = nil
		return nil
	}

	// 检查是否是数组
	var arrayCheck []json.RawMessage
	if err := json.Unmarshal(raw, &arrayCheck); err != nil {
		// 不是数组，尝试解析为单个对象
		var singleObj map[string]interface{}
		if err := json.Unmarshal(raw, &singleObj); err != nil {
			// 解析失败，设为nil
			*j = nil
			return nil
		}
		*j = JSONB(singleObj)
	} else {
		// 是数组，直接解析
		if err := json.Unmarshal(bytes, j); err != nil {
			// 解析失败，设为nil
			*j = nil
			return nil
		}
	}

	return nil
}

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

	// 从销售记录中提取的额外信息
	OriginalCustomerID string `json:"original_customer_id" gorm:"size:256"` // 原始客户ID
	ImportSource       string `json:"import_source" gorm:"size:256"`        // 导入来源
	SallerName         string `json:"saller_name" gorm:"size:256"`          // 销售员姓名

	// 系统标签字段
	SystemTags pq.Int64Array `json:"system_tags" gorm:"type:int4[]"` // 系统标签ID数组

	// 订单相关字段
	LastOrderDate *time.Time `json:"last_order_date" gorm:"comment:最后下单时间"` // 最后下单时间
}

func (Customer) TableName() string {
	return "customers"
}
