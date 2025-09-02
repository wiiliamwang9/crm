package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// 设置测试环境
func setupTestServer() *gin.Engine {
	// 加载测试配置
	gin.SetMode(gin.TestMode)

	// 连接数据库（使用相同的配置）
	LoadConfig("config.yml")
	ConnectDatabase()

	// 创建路由
	r := gin.New()
	SetupRoutes(r)

	return r
}

// 测试健康检查接口
func TestHealthCheck(t *testing.T) {
	router := setupTestServer()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "healthy", response["status"])
	assert.Equal(t, "CRM API", response["service"])
}

// 测试客户列表接口
func TestGetCustomers(t *testing.T) {
	router := setupTestServer()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/customers?page=1&limit=10", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response, "data")
	assert.Contains(t, response, "total")
}

// 测试客户搜索接口
func TestSearchCustomers(t *testing.T) {
	router := setupTestServer()

	// 测试关键词搜索
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/customers/search?keyword=测试", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response, "data")

	// 验证返回的数据结构
	data, ok := response["data"].([]interface{})
	assert.True(t, ok)

	if len(data) > 0 {
		customer := data[0].(map[string]interface{})
		assert.Contains(t, customer, "id")
		assert.Contains(t, customer, "name")
		assert.Contains(t, customer, "contact_name")
	}
}

// 测试创建客户接口
func TestCreateCustomer(t *testing.T) {
	router := setupTestServer()

	customerReq := CustomerRequest{
		Name:        "测试客户" + fmt.Sprintf("%d", time.Now().Unix()),
		ContactName: "测试联系人",
		Phones:      []string{"13800138000"},
		Province:    "广东省",
		City:        "深圳市",
		District:    "南山区",
		Address:     "测试地址",
		Category:    "零售店",
		State:       1,
		Level:       2,
		Source:      "测试来源",
	}

	jsonData, _ := json.Marshal(customerReq)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/customers", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response, "data")

	// 验证创建的客户数据
	data := response["data"].(map[string]interface{})
	assert.Equal(t, customerReq.Name, data["name"])
	assert.Equal(t, customerReq.ContactName, data["contact_name"])
}

// 测试获取单个客户接口
func TestGetCustomer(t *testing.T) {
	router := setupTestServer()

	// 先创建一个测试客户
	customerReq := CustomerRequest{
		Name:        "获取测试客户" + fmt.Sprintf("%d", time.Now().Unix()),
		ContactName: "测试联系人",
		Phones:      []string{"13900139000"},
	}

	jsonData, _ := json.Marshal(customerReq)

	// 创建客户
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("POST", "/api/v1/customers", bytes.NewBuffer(jsonData))
	req1.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w1, req1)

	var createResponse map[string]interface{}
	json.Unmarshal(w1.Body.Bytes(), &createResponse)
	customerData := createResponse["data"].(map[string]interface{})
	customerID := int(customerData["id"].(float64))

	// 获取客户详情
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/customers/%d", customerID), nil)
	router.ServeHTTP(w2, req2)

	assert.Equal(t, 200, w2.Code)

	var getResponse map[string]interface{}
	json.Unmarshal(w2.Body.Bytes(), &getResponse)
	assert.Contains(t, getResponse, "data")

	data := getResponse["data"].(map[string]interface{})
	assert.Equal(t, customerReq.Name, data["name"])
	assert.Equal(t, float64(customerID), data["id"])
}

// 测试用户列表接口
func TestGetUsers(t *testing.T) {
	router := setupTestServer()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/users", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response, "data")
}

// 测试待办事项列表接口
func TestGetTodos(t *testing.T) {
	router := setupTestServer()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/todos?page=1&page_size=10", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response, "data")
	assert.Contains(t, response, "total")
}

// 测试跟进记录列表接口
func TestGetActivities(t *testing.T) {
	router := setupTestServer()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/activities?page=1&page_size=10", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response, "data")
	assert.Contains(t, response, "total")
}

// 测试仪表板搜索接口
func TestDashboardSearch(t *testing.T) {
	router := setupTestServer()

	searchReq := DashboardSearchRequest{
		UserID:       1,
		TimeFilter:   "今日待跟进",
		StatusFilter: "全部",
		Page:         1,
		PageSize:     10,
	}

	jsonData, _ := json.Marshal(searchReq)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/dashboard/search", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response, "data")

	data := response["data"].(map[string]interface{})
	assert.Contains(t, data, "list")
	assert.Contains(t, data, "total")
	assert.Contains(t, data, "page")
	assert.Contains(t, data, "page_size")
}

// 测试标签列表接口
func TestGetTags(t *testing.T) {
	router := setupTestServer()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/tags", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response, "data")
}

// 测试系统标签搜索
func TestSearchCustomersWithSystemTags(t *testing.T) {
	router := setupTestServer()

	// 测试系统标签搜索
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/customers/search?system_tags=1,2,3", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response, "data")
}

// 集成测试：完整的客户操作流程
func TestCustomerCRUDFlow(t *testing.T) {
	router := setupTestServer()

	// 1. 创建客户
	customerReq := CustomerRequest{
		Name:        "CRUD测试客户" + fmt.Sprintf("%d", time.Now().Unix()),
		ContactName: "CRUD联系人",
		Phones:      []string{"13700137000"},
		Category:    "连锁店",
		State:       1,
		Level:       3,
	}

	jsonData, _ := json.Marshal(customerReq)

	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("POST", "/api/v1/customers", bytes.NewBuffer(jsonData))
	req1.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w1, req1)

	assert.Equal(t, 200, w1.Code)

	var createResponse map[string]interface{}
	json.Unmarshal(w1.Body.Bytes(), &createResponse)
	customerData := createResponse["data"].(map[string]interface{})
	customerID := int(customerData["id"].(float64))

	// 2. 获取客户详情
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/customers/%d", customerID), nil)
	router.ServeHTTP(w2, req2)

	assert.Equal(t, 200, w2.Code)

	// 3. 更新客户
	updateReq := CustomerRequest{
		Name:        customerReq.Name + "_更新",
		ContactName: customerReq.ContactName + "_更新",
		Phones:      customerReq.Phones,
		Category:    "更新分类",
		State:       2,
		Level:       4,
	}

	updateData, _ := json.Marshal(updateReq)

	w3 := httptest.NewRecorder()
	req3, _ := http.NewRequest("PUT", fmt.Sprintf("/api/v1/customers/%d", customerID), bytes.NewBuffer(updateData))
	req3.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w3, req3)

	assert.Equal(t, 200, w3.Code)

	var updateResponse map[string]interface{}
	json.Unmarshal(w3.Body.Bytes(), &updateResponse)
	updatedData := updateResponse["data"].(map[string]interface{})
	assert.Equal(t, updateReq.Name, updatedData["name"])
	assert.Equal(t, updateReq.Category, updatedData["category"])

	// 4. 删除客户
	w4 := httptest.NewRecorder()
	req4, _ := http.NewRequest("DELETE", fmt.Sprintf("/api/v1/customers/%d", customerID), nil)
	router.ServeHTTP(w4, req4)

	assert.Equal(t, 200, w4.Code)
}
