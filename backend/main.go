package main

import (
	"fmt"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置文件
	_, err := LoadConfig("config.yml")
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// 设置Gin模式
	SetGinMode()

	// 连接数据库
	ConnectDatabase()

	// 自动迁移数据库表
	DB.AutoMigrate(&Customer{}, &Todo{}, &TodoLog{},
		&Reminder{}, &ReminderTemplate{}, &ReminderConfig{},
		&Activity{}, &User{}, &TagDimension{}, &Tag{})

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

	// 设置路由
	SetupRoutes(r)

	port := GetServerPort()
	fmt.Printf("Starting CRM server on %s\n", port)
	r.Run(port)
}
