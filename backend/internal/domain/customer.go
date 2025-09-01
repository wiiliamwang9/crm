package domain

import (
	"database/sql/driver"
	"encoding/json"
	"strconv"
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

// DTO Structures and Methods (merged from customer_dto.go)

// CustomerRequest 客户请求结构体
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

// CustomerResponse 客户响应结构体
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

// CustomerListRequest 客户列表请求结构体
type CustomerListRequest struct {
	Page   int    `form:"page,default=1" binding:"min=1"`
	Limit  int    `form:"limit,default=20" binding:"min=1,max=100"`
	Search string `form:"search"`
}

// SearchRequest 搜索请求结构体
type SearchRequest struct {
	Search     string  `json:"search"`
	SystemTags []int64 `json:"system_tags"`
	Limit      int     `json:"limit" binding:"min=1,max=100"`
	Page       int     `json:"page" binding:"min=1"`
	Timestamp  int64   `json:"timestamp"`
}

// CustomerFavorsRequest 客户偏好请求结构体
type CustomerFavorsRequest struct {
	Favors []map[string]interface{} `json:"favors" binding:"required"`
}

// ToModel 将CustomerRequest转换为Customer模型
func (req *CustomerRequest) ToModel() *Customer {
	return &Customer{
		Name:         req.Name,
		ContactName:  req.ContactName,
		Phones:       req.Phones,
		Wechats:      req.Wechats,
		Province:     req.Province,
		City:         req.City,
		District:     req.District,
		Address:      req.Address,
		Products:     req.Products,
		Category:     req.Category,
		Tags:         req.Tags,
		State:        req.State,
		Level:        req.Level,
		Source:       req.Source,
		ImportSource: req.ImportSource,
		Remark:       req.Remark,
		SallerName:   req.SallerName,
		Sellers:      req.Sellers,
	}
}

// CustomerToResponse 将Customer模型转换为CustomerResponse
func CustomerToResponse(customer *Customer) *CustomerResponse {
	mainPhone := ""
	if len(customer.Phones) > 0 {
		mainPhone = customer.Phones[0]
	}

	mainWechat := ""
	if len(customer.Wechats) > 0 {
		mainWechat = customer.Wechats[0]
	}

	sellerInfo := ""
	if len(customer.Sellers) > 0 {
		sellerInfo = strconv.FormatInt(customer.Sellers[0], 10)
	}

	fullAddress := customer.Address
	if customer.Province != "" || customer.City != "" || customer.District != "" {
		fullAddress = customer.Province + customer.City + customer.District + customer.Address
	}

	organization := ""

	return &CustomerResponse{
		ID:           customer.ID,
		Name:         customer.Name,
		Phone:        mainPhone,
		Phones:       customer.Phones,
		Wechat:       mainWechat,
		Wechats:      customer.Wechats,
		Seller:       sellerInfo,
		Sellers:      customer.Sellers,
		SallerName:   customer.SallerName,
		Address:      fullAddress,
		Province:     customer.Province,
		City:         customer.City,
		District:     customer.District,
		Company:      customer.Name,
		Products:     customer.Products,
		Category:     customer.Category,
		Tags:         customer.Tags,
		State:        customer.State,
		Level:        customer.Level,
		ContactName:  customer.ContactName,
		Source:       customer.Source,
		ImportSource: customer.ImportSource,
		Remark:       customer.Remark,
		Organization: organization,
		CreatedAt:    customer.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}
