package main

import (
	"crm/config"
	models2 "crm/internal/domain"
	"crm/internal/middleware"
	"crm/internal/routes"
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// 连接数据库
	config.ConnectDatabase()

	// 自动迁移数据库表
	config.DB.AutoMigrate(&models2.Customer{}, &models2.Todo{}, &models2.TodoLog{},
		&models2.Reminder{}, &models2.ReminderTemplate{}, &models2.ReminderConfig{},
		&models2.Activity{}, &models2.User{}, &models2.TagDimension{}, &models2.Tag{})

	// 创建Gin引擎
	r := gin.Default()

	// 添加CORS中间件
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	}))

	// 添加字符编码中间件
	r.Use(func(c *gin.Context) {
		c.Header("Content-Type", "application/json; charset=utf-8")
		c.Next()
	})

	// 添加中间件
	r.Use(middleware.ErrorHandler())

	// 设置路由
	routes.SetupRoutes(r, config.DB)

	// 添加静态文件服务和其他路由
	config.SetupStaticRoutes(r)

	fmt.Println("Starting CRM server on :8081")
	r.Run(":8081")
}
