package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// 测试数据库设置
func setupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// 自动迁移
	db.AutoMigrate(&Customer{})

	return db
}

// 创建测试客户
func createTestCustomer(db *gorm.DB) *Customer {
	customer := &Customer{
		BaseModel: BaseModel{ID: 1},
		Name:      "测试客户",
		Phone:     "13800138000",
		Email:     "test@example.com",
		Favors:    JSONB{},
	}
	db.Create(customer)
	return customer
}

// 设置测试路由
func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	SetupRoutes(router)
	return router
}

// TestGetCustomerPreferences 测试获取客户偏好列表
func TestGetCustomerPreferences(t *testing.T) {
	// 设置测试数据库
	originalDB := db
	db = setupTestDB()
	defer func() { db = originalDB }()

	// 创建测试客户
	customer := createTestCustomer(db)

	// 添加测试偏好数据
	preferences := []map[string]interface{}{
		{
			"id":          "pref_1",
			"type":        "product",
			"content":     "喜欢高端产品",
			"created_at":  time.Now().Format(time.RFC3339),
		},
		{
			"id":          "pref_2",
			"type":        "service",
			"content":     "偏好上门服务",
			"created_at":  time.Now().Format(time.RFC3339),
		},
	}
	preferencesJSON, _ := json.Marshal(preferences)
	customer.Favors = JSONB(preferencesJSON)
	db.Save(customer)

	// 设置路由
	router := setupTestRouter()

	t.Run("成功获取偏好列表", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/customers/1/preferences", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response CustomerPreferenceListResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(response.Preferences))
		assert.Equal(t, "pref_1", response.Preferences[0].ID)
		assert.Equal(t, "product", response.Preferences[0].Type)
	})

	t.Run("客户不存在", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/customers/999/preferences", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("无效的客户ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/customers/invalid/preferences", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

// TestCreateCustomerPreference 测试创建客户偏好
func TestCreateCustomerPreference(t *testing.T) {
	// 设置测试数据库
	originalDB := db
	db = setupTestDB()
	defer func() { db = originalDB }()

	// 创建测试客户
	customer := createTestCustomer(db)

	// 设置路由
	router := setupTestRouter()

	t.Run("成功创建偏好", func(t *testing.T) {
		request := CustomerPreferenceCreateRequest{
			Type:    "product",
			Content: "喜欢智能家居产品",
		}
		requestJSON, _ := json.Marshal(request)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/customers/1/preferences", bytes.NewBuffer(requestJSON))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response CustomerPreferenceResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotEmpty(t, response.Preference.ID)
		assert.Equal(t, "product", response.Preference.Type)
		assert.Equal(t, "喜欢智能家居产品", response.Preference.Content)
		assert.NotEmpty(t, response.Preference.CreatedAt)

		// 验证数据库中的数据
		var updatedCustomer Customer
		db.First(&updatedCustomer, customer.ID)
		var preferences []map[string]interface{}
		json.Unmarshal([]byte(updatedCustomer.Favors), &preferences)
		assert.Equal(t, 1, len(preferences))
	})

	t.Run("缺少必填字段", func(t *testing.T) {
		request := CustomerPreferenceCreateRequest{
			Type: "product",
			// Content 缺失
		}
		requestJSON, _ := json.Marshal(request)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/customers/1/preferences", bytes.NewBuffer(requestJSON))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("客户不存在", func(t *testing.T) {
		request := CustomerPreferenceCreateRequest{
			Type:    "product",
			Content: "喜欢智能家居产品",
		}
		requestJSON, _ := json.Marshal(request)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/customers/999/preferences", bytes.NewBuffer(requestJSON))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("无效的JSON格式", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/customers/1/preferences", strings.NewReader("invalid json"))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

// TestUpdateCustomerPreference 测试更新客户偏好
func TestUpdateCustomerPreference(t *testing.T) {
	// 设置测试数据库
	originalDB := db
	db = setupTestDB()
	defer func() { db = originalDB }()

	// 创建测试客户
	customer := createTestCustomer(db)

	// 添加测试偏好数据
	preferences := []map[string]interface{}{
		{
			"id":          "pref_1",
			"type":        "product",
			"content":     "喜欢高端产品",
			"created_at":  time.Now().Format(time.RFC3339),
		},
	}
	preferencesJSON, _ := json.Marshal(preferences)
	customer.Favors = JSONB(preferencesJSON)
	db.Save(customer)

	// 设置路由
	router := setupTestRouter()

	t.Run("成功更新偏好", func(t *testing.T) {
		request := CustomerPreferenceUpdateRequest{
			Type:    "product",
			Content: "更喜欢超高端产品",
		}
		requestJSON, _ := json.Marshal(request)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/api/v1/customers/1/preferences/pref_1", bytes.NewBuffer(requestJSON))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response CustomerPreferenceResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "pref_1", response.Preference.ID)
		assert.Equal(t, "product", response.Preference.Type)
		assert.Equal(t, "更喜欢超高端产品", response.Preference.Content)
	})

	t.Run("偏好不存在", func(t *testing.T) {
		request := CustomerPreferenceUpdateRequest{
			Type:    "product",
			Content: "更喜欢超高端产品",
		}
		requestJSON, _ := json.Marshal(request)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/api/v1/customers/1/preferences/nonexistent", bytes.NewBuffer(requestJSON))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("客户不存在", func(t *testing.T) {
		request := CustomerPreferenceUpdateRequest{
			Type:    "product",
			Content: "更喜欢超高端产品",
		}
		requestJSON, _ := json.Marshal(request)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/api/v1/customers/999/preferences/pref_1", bytes.NewBuffer(requestJSON))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

// TestDeleteCustomerPreference 测试删除客户偏好
func TestDeleteCustomerPreference(t *testing.T) {
	// 设置测试数据库
	originalDB := db
	db = setupTestDB()
	defer func() { db = originalDB }()

	// 创建测试客户
	customer := createTestCustomer(db)

	// 添加测试偏好数据
	preferences := []map[string]interface{}{
		{
			"id":          "pref_1",
			"type":        "product",
			"content":     "喜欢高端产品",
			"created_at":  time.Now().Format(time.RFC3339),
		},
		{
			"id":          "pref_2",
			"type":        "service",
			"content":     "偏好上门服务",
			"created_at":  time.Now().Format(time.RFC3339),
		},
	}
	preferencesJSON, _ := json.Marshal(preferences)
	customer.Favors = JSONB(preferencesJSON)
	db.Save(customer)

	// 设置路由
	router := setupTestRouter()

	t.Run("成功删除偏好", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/api/v1/customers/1/preferences/pref_1", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// 验证数据库中的数据
		var updatedCustomer Customer
		db.First(&updatedCustomer, customer.ID)
		var remainingPreferences []map[string]interface{}
		json.Unmarshal([]byte(updatedCustomer.Favors), &remainingPreferences)
		assert.Equal(t, 1, len(remainingPreferences))
		assert.Equal(t, "pref_2", remainingPreferences[0]["id"])
	})

	t.Run("偏好不存在", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/api/v1/customers/1/preferences/nonexistent", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("客户不存在", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/api/v1/customers/999/preferences/pref_1", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

// TestPreferenceIntegration 集成测试：完整的偏好管理流程
func TestPreferenceIntegration(t *testing.T) {
	// 设置测试数据库
	originalDB := db
	db = setupTestDB()
	defer func() { db = originalDB }()

	// 创建测试客户
	customer := createTestCustomer(db)

	// 设置路由
	router := setupTestRouter()

	t.Run("完整的偏好管理流程", func(t *testing.T) {
		// 1. 初始状态：获取空的偏好列表
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/customers/1/preferences", nil)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		var listResponse CustomerPreferenceListResponse
		json.Unmarshal(w.Body.Bytes(), &listResponse)
		assert.Equal(t, 0, len(listResponse.Preferences))

		// 2. 创建第一个偏好
		createRequest1 := CustomerPreferenceCreateRequest{
			Type:    "product",
			Content: "喜欢智能手机",
		}
		createJSON1, _ := json.Marshal(createRequest1)

		w = httptest.NewRecorder()
		req, _ = http.NewRequest("POST", "/api/v1/customers/1/preferences", bytes.NewBuffer(createJSON1))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)

		var createResponse1 CustomerPreferenceResponse
		json.Unmarshal(w.Body.Bytes(), &createResponse1)
		preferenceID1 := createResponse1.Preference.ID

		// 3. 创建第二个偏好
		createRequest2 := CustomerPreferenceCreateRequest{
			Type:    "service",
			Content: "偏好线上咨询",
		}
		createJSON2, _ := json.Marshal(createRequest2)

		w = httptest.NewRecorder()
		req, _ = http.NewRequest("POST", "/api/v1/customers/1/preferences", bytes.NewBuffer(createJSON2))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)

		// 4. 获取偏好列表，应该有2个偏好
		w = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", "/api/v1/customers/1/preferences", nil)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		json.Unmarshal(w.Body.Bytes(), &listResponse)
		assert.Equal(t, 2, len(listResponse.Preferences))

		// 5. 更新第一个偏好
		updateRequest := CustomerPreferenceUpdateRequest{
			Type:    "product",
			Content: "更喜欢iPhone",
		}
		updateJSON, _ := json.Marshal(updateRequest)

		w = httptest.NewRecorder()
		req, _ = http.NewRequest("PUT", fmt.Sprintf("/api/v1/customers/1/preferences/%s", preferenceID1), bytes.NewBuffer(updateJSON))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		var updateResponse CustomerPreferenceResponse
		json.Unmarshal(w.Body.Bytes(), &updateResponse)
		assert.Equal(t, "更喜欢iPhone", updateResponse.Preference.Content)

		// 6. 删除第一个偏好
		w = httptest.NewRecorder()
		req, _ = http.NewRequest("DELETE", fmt.Sprintf("/api/v1/customers/1/preferences/%s", preferenceID1), nil)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		// 7. 最终检查：应该只剩1个偏好
		w = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", "/api/v1/customers/1/preferences", nil)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		json.Unmarshal(w.Body.Bytes(), &listResponse)
		assert.Equal(t, 1, len(listResponse.Preferences))
		assert.Equal(t, "service", listResponse.Preferences[0].Type)
		assert.Equal(t, "偏好线上咨询", listResponse.Preferences[0].Content)
	})
}