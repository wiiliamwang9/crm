package routes

import (
	"crm/internal/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupRoutes 设置所有路由
func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	// 初始化仓储层
	repo := repository.NewRepository(db)
	// 启用CORS中间件
	r.Use(func(c *gin.Context) {
		// 只设置一次CORS头，避免重复
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", origin)
		} else {
			c.Header("Access-Control-Allow-Origin", "*")
		}
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Requested-With")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// API路由组
	api := r.Group("/api")
	{
		// 设置客户路由
		SetupCustomerRoutes(api, repo)

		// 设置待办路由
		SetupTodoRoutes(api, repo)

		// 设置提醒路由
		SetupReminderRoutes(api, repo)

		// 设置用户路由
		SetupUserRoutes(api, repo)

		// 设置标签路由
		SetupTagRoutes(api, repo)

		// 设置跟进记录路由
		SetupActivityRoutes(api, repo)
	}
}
