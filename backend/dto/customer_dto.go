package dto

import (
	"crm/models"
	"strconv"
)

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

type CustomerListRequest struct {
	Page   int    `form:"page,default=1" binding:"min=1"`
	Limit  int    `form:"limit,default=20" binding:"min=1,max=100"`
	Search string `form:"search"`
}

type SearchRequest struct {
	Search     string  `json:"search"`
	SystemTags []int64 `json:"system_tags"`
	Limit      int     `json:"limit" binding:"min=1,max=100"`
	Page       int     `json:"page" binding:"min=1"`
	Timestamp  int64   `json:"timestamp"`
}

type CustomerFavorsRequest struct {
	Favors []map[string]interface{} `json:"favors" binding:"required"`
}

func (req *CustomerRequest) ToModel() *models.Customer {
	return &models.Customer{
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

func CustomerToResponse(customer *models.Customer) *CustomerResponse {
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
