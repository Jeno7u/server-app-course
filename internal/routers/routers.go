package routers
package routers

import (
	"github.com/Jeno7u/server-app-course/internal/handlers"
	"github.com/Jeno7u/server-app-course/internal/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterTaskRoutes(group *gin.RouterGroup) {
	tasks := group.Group("/tasks")
	tasks.Use(middleware.MockAuth())
	{
		tasks.POST("", handlers.CreateTask)
		tasks.GET("", handlers.GetTasks)
		tasks.GET("/:task_id", handlers.GetTaskByID)
		tasks.PATCH("/:task_id/status", handlers.UpdateTaskStatus)
		tasks.DELETE("/:task_id", handlers.DeleteTask)
	}
}

func RegisterUserRoutes(group *gin.RouterGroup) {
	users := group.Group("/users")
	users.Use(middleware.MockAuth())
	{
		users.GET("/me", handlers.GetMe)
		users.GET("/:user_id", handlers.GetUserByID)
	}
}

func RegisterAdminRoutes(group *gin.RouterGroup) {
	admin := group.Group("/admin")
	admin.Use(middleware.MockAuth(), middleware.RequireAdmin())
	{
		admin.GET("/stats", handlers.AdminStats)
		admin.DELETE("/tasks/:task_id", handlers.AdminDeleteTask)
	}
}