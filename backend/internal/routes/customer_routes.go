package routes

import (
	"crm/internal/api"
	"crm/internal/repository"
	"crm/internal/services"

	"github.com/gin-gonic/gin"
)

// SetupCustomerRoutes 设置客户相关路由
func SetupCustomerRoutes(router *gin.RouterGroup, repo *repository.Repository) {
	customerService := services.NewCustomerService(repo.Customer)
	customerHandler := api.NewCustomerHandler(customerService)

	v1 := router.Group("/v1")
	{
		v1.GET("/customers", customerHandler.GetCustomers)
		v1.POST("/customers", customerHandler.CreateCustomer)
		v1.GET("/customers/:id", customerHandler.GetCustomer)
		v1.PUT("/customers/:id", customerHandler.UpdateCustomer)
		v1.DELETE("/customers/:id", customerHandler.DeleteCustomer)
		v1.PUT("/customers/:id/favors", customerHandler.UpdateCustomerFavors)
		v1.PUT("/customers/:id/remark", customerHandler.UpdateCustomerRemark)
		v1.PUT("/customers/:id/system-tags", customerHandler.UpdateCustomerSystemTags)
		v1.POST("/customers/search", customerHandler.PostSearchCustomers)
		v1.GET("/customers/special", customerHandler.GetSpecialCustomers)

		// Excel导入导出路由
		v1.POST("/upload-excel", customerHandler.UploadCustomersExcel)
		v1.GET("/export-excel", customerHandler.ExportCustomersExcel)
	}
}
