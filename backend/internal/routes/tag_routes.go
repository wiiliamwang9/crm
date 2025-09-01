package routes

import (
	"crm/internal/api"
	"crm/internal/repository"
	"crm/internal/services"

	"github.com/gin-gonic/gin"
)

// SetupTagRoutes 设置标签路由
func SetupTagRoutes(router *gin.RouterGroup, repo *repository.Repository) {
	// 创建服务和处理器
	tagService := services.NewTagService(repo.Tag, repo.TagDimension)
	tagHandler := api.NewTagHandler(tagService)

	// 标签维度路由组
	dimensionGroup := router.Group("/v1/tag-dimensions")
	{
		dimensionGroup.GET("", tagHandler.GetDimensions)          // 获取所有维度及其标签
		dimensionGroup.GET("/:id", tagHandler.GetDimension)       // 获取单个维度
		dimensionGroup.POST("", tagHandler.CreateDimension)       // 创建维度
		dimensionGroup.PUT("/:id", tagHandler.UpdateDimension)    // 更新维度
		dimensionGroup.DELETE("/:id", tagHandler.DeleteDimension) // 删除维度
	}

	// 标签路由组
	tagGroup := router.Group("/v1/tags")
	{
		tagGroup.GET("", tagHandler.GetTags)                                      // 获取标签列表（带分页和筛选）
		tagGroup.GET("/active", tagHandler.GetAllActiveTags)                      // 获取所有活跃标签（下拉选择）
		tagGroup.GET("/:id", tagHandler.GetTag)                                   // 获取单个标签
		tagGroup.POST("", tagHandler.CreateTag)                                   // 创建标签
		tagGroup.PUT("/:id", tagHandler.UpdateTag)                                // 更新标签
		tagGroup.DELETE("/:id", tagHandler.DeleteTag)                             // 删除标签
		tagGroup.GET("/dimension/:dimension_id", tagHandler.GetTagsByDimensionID) // 根据维度ID获取标签
	}
}
