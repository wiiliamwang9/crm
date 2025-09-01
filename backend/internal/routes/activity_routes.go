package routes

import (
	"crm/internal/api"
	"crm/internal/repository"
	"crm/internal/services"

	"github.com/gin-gonic/gin"
)

// SetupActivityRoutes 设置跟进记录相关路由
func SetupActivityRoutes(router *gin.RouterGroup, repo *repository.Repository) {
	activityService := services.NewActivityService(repo.Activity, repo.Todo)
	activityHandler := api.NewActivityHandler(activityService)
	// 跟进记录相关路由
	v1 := router.Group("/v1")
	activities := v1.Group("/activities")
	{
		// 创建跟进记录
		activities.POST("", activityHandler.CreateActivity)

		// 获取客户的跟进记录列表
		activities.GET("/customer/:customer_id", activityHandler.GetActivitiesByCustomer)

		// 获取跟进记录详情
		activities.GET("/:id", activityHandler.GetActivityByID)

		// 更新跟进记录
		activities.PUT("/:id", activityHandler.UpdateActivity)

		// 更新跟进记录反馈
		activities.PUT("/:id/feedback", activityHandler.UpdateActivityFeedback)

		// 删除跟进记录
		activities.DELETE("/:id", activityHandler.DeleteActivity)

		// 获取客户的跟进记录统计信息
		activities.GET("/customer/:customer_id/statistics", activityHandler.GetActivityStatistics)

		// 获取需要跟进的记录
		activities.GET("/need-follow-up", activityHandler.GetActivitiesNeedingFollowUp)
	}
}
