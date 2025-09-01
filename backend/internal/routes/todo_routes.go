package routes

import (
	"crm/internal/api"
	"crm/internal/repository"
	"crm/internal/services"

	"github.com/gin-gonic/gin"
)

// SetupTodoRoutes 设置待办相关路由
func SetupTodoRoutes(router *gin.RouterGroup, repo *repository.Repository) {
	todoService := services.NewTodoService(repo.Todo)
	todoHandler := api.NewTodoHandler(todoService)

	v1 := router.Group("/v1")
	{
		v1.POST("/todos", todoHandler.CreateTodo)
		v1.GET("/todos", todoHandler.GetTodos)
		v1.GET("/todos/:id", todoHandler.GetTodo)
		v1.PUT("/todos/:id", todoHandler.UpdateTodo)
		v1.DELETE("/todos/:id", todoHandler.DeleteTodo)
		v1.POST("/todos/:id/complete", todoHandler.CompleteTodo)
		v1.POST("/todos/:id/cancel", todoHandler.CancelTodo)
		v1.GET("/todos/stats", todoHandler.GetTodoStats)
	}
}
