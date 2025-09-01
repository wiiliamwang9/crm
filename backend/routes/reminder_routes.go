package routes

import (
	"crm/handlers"
	"crm/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupReminderRoutes 设置提醒相关路由
func SetupReminderRoutes(router *gin.RouterGroup, db *gorm.DB) {
	reminderService := services.NewReminderService(db)
	reminderHandler := handlers.NewReminderHandler(reminderService)

	reminderGroup := router.Group("/reminders")
	{
		// 提醒管理
		reminderGroup.POST("", reminderHandler.CreateReminder)           // 创建提醒
		reminderGroup.GET("", reminderHandler.GetReminders)             // 获取提醒列表
		reminderGroup.GET("/:id", reminderHandler.GetReminder)          // 获取单个提醒
		reminderGroup.PUT("/:id", reminderHandler.UpdateReminder)       // 更新提醒
		reminderGroup.DELETE("/:id", reminderHandler.DeleteReminder)    // 删除提醒
		reminderGroup.POST("/:id/cancel", reminderHandler.CancelReminder) // 取消提醒
		
		// 提醒统计
		reminderGroup.GET("/stats", reminderHandler.GetReminderStats)
		
		// 手动触发处理提醒
		reminderGroup.POST("/process", reminderHandler.ProcessReminders)
		
		// 用户配置
		reminderGroup.GET("/config", reminderHandler.GetUserReminderConfig)       // 获取用户配置
		reminderGroup.PUT("/config", reminderHandler.UpdateUserReminderConfig)    // 更新用户配置
	}
}