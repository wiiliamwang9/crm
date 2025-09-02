package main

import (
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupRoutes 设置所有路由
func SetupRoutes(r *gin.Engine) {
	// CORS中间件
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	}))

	// 字符编码中间件
	r.Use(func(c *gin.Context) {
		c.Header("Content-Type", "application/json; charset=utf-8")
		c.Next()
	})

	// API路由组
	api := r.Group("/api/v1")
	{
		// 客户相关路由
		api.GET("/customers", func(c *gin.Context) {
			page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
			limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
			search := c.Query("search")
			customers, total := getCustomers(page, limit, search)
			c.JSON(200, gin.H{"data": customers, "total": total})
		})

		api.POST("/customers", func(c *gin.Context) {
			var req CustomerRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}
			customer := createCustomer(req)
			c.JSON(200, gin.H{"data": customer})
		})

		api.GET("/customers/:id", func(c *gin.Context) {
			id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
			customer := getCustomer(id)
			c.JSON(200, gin.H{"data": customer})
		})

		api.PUT("/customers/:id", func(c *gin.Context) {
			id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
			var req CustomerRequest
			c.ShouldBindJSON(&req)
			customer := updateCustomer(id, req)
			c.JSON(200, gin.H{"data": customer})
		})

		api.DELETE("/customers/:id", func(c *gin.Context) {
			id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
			deleteCustomer(id)
			c.JSON(200, gin.H{"message": "删除成功"})
		})

		// 客户搜索路由
		api.GET("/customers/search", func(c *gin.Context) {
			keyword := c.Query("keyword")
			systemTagsStr := c.Query("system_tags") // 逗号分隔的标签ID，如"1,2,3"
			customers := searchCustomers(keyword, systemTagsStr)
			c.JSON(200, gin.H{"data": customers})
		})

		// 待办事项路由
		api.GET("/todos", func(c *gin.Context) {
			customerID, _ := strconv.ParseUint(c.Query("customer_id"), 10, 64)
			page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
			pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
			todos, total := getTodos(customerID, page, pageSize)
			c.JSON(200, gin.H{"data": todos, "total": total})
		})

		api.POST("/todos", func(c *gin.Context) {
			var req TodoCreateRequest
			c.ShouldBindJSON(&req)
			todo := createTodo(req)
			c.JSON(200, gin.H{"data": todo})
		})

		api.PUT("/todos/:id", func(c *gin.Context) {
			id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
			var req TodoUpdateRequest
			c.ShouldBindJSON(&req)
			todo := updateTodo(id, req)
			c.JSON(200, gin.H{"data": todo})
		})

		// 跟进记录路由
		api.GET("/activities", func(c *gin.Context) {
			customerID, _ := strconv.ParseUint(c.Query("customer_id"), 10, 64)
			page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
			pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
			activities, total := getActivities(customerID, page, pageSize)
			c.JSON(200, gin.H{"data": activities, "total": total})
		})

		api.POST("/activities", func(c *gin.Context) {
			var req ActivityCreateRequest
			c.ShouldBindJSON(&req)
			activity := createActivity(req)
			c.JSON(200, gin.H{"data": activity})
		})

		// 用户路由
		api.GET("/users", func(c *gin.Context) {
			users := getUsers()
			c.JSON(200, gin.H{"data": users})
		})

		api.GET("/users/:id", func(c *gin.Context) {
			id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
			user := getUserDetail(id)
			c.JSON(200, gin.H{"data": user})
		})

		// 仪表板搜索路由
		api.POST("/dashboard/search", func(c *gin.Context) {
			var req DashboardSearchRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}
			if req.Page <= 0 {
				req.Page = 1
			}
			if req.PageSize <= 0 {
				req.PageSize = 20
			}

			// 如果设置了查看全部标志，则设置大页大小展示全部数据
			if req.ShowAll {
				req.PageSize = 10000 // 设置足够大的页大小
				req.Page = 1
			}

			results, total := searchDashboardData(req)
			c.JSON(200, gin.H{"data": gin.H{"list": results, "total": total, "page": req.Page, "page_size": req.PageSize}})
		})

		// 标签路由
		api.GET("/tags", func(c *gin.Context) {
			tags := getTags()
			c.JSON(200, gin.H{"data": tags})
		})

		api.POST("/tags", func(c *gin.Context) {
			var req TagCreateRequest
			c.ShouldBindJSON(&req)
			tag := createTag(req)
			c.JSON(200, gin.H{"data": tag})
		})

		// 提醒路由
		api.GET("/reminders", func(c *gin.Context) {
			userID, _ := strconv.ParseUint(c.Query("user_id"), 10, 64)
			page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
			pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
			reminders, total := getReminders(userID, page, pageSize)
			c.JSON(200, gin.H{"data": reminders, "total": total})
		})

		api.POST("/reminders", func(c *gin.Context) {
			var req ReminderCreateRequest
			c.ShouldBindJSON(&req)
			reminder := createReminder(req)
			c.JSON(200, gin.H{"data": reminder})
		})
	}

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy", "service": "CRM API"})
	})
}
