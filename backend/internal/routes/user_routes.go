package routes

import (
	"crm/internal/api"
	"crm/internal/repository"
	"crm/internal/services"

	"github.com/gin-gonic/gin"
)

// SetupUserRoutes 设置用户路由
func SetupUserRoutes(router *gin.RouterGroup, repo *repository.Repository) {
	// 创建服务和处理器
	userService := services.NewUserService(repo.User)
	userHandler := api.NewUserHandler(userService)

	// 用户路由组
	userGroup := router.Group("/v1/users")
	{
		userGroup.GET("", userHandler.GetUsers)                         // 获取用户列表
		userGroup.GET("/active", userHandler.GetAllActiveUsers)         // 获取所有活跃用户（下拉选择）
		userGroup.GET("/:id", userHandler.GetUser)                      // 获取单个用户
		userGroup.GET("/:id/homepage", userHandler.GetHomepageUserInfo) // 获取首页用户信息
	}
}
