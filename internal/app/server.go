package app

import (
	"github.com/Jeno7u/server-app-course/internal/errors"
	"github.com/Jeno7u/server-app-course/internal/handlers"
	"github.com/Jeno7u/server-app-course/internal/middleware"
	"github.com/Jeno7u/server-app-course/internal/routers"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	// Assignment 4: Custom Error Handling
	router.Use(errors.ErrorHandler())

	// Health check (Assignment 5)
	router.GET("/health", handlers.HealthCheck)

	// WebSockets (Assignment 5)
	router.GET("/ws/rooms/:room_id", handlers.ConnectRoom)
	router.GET("/rooms/:room_id/users", handlers.GetRoomUsers)

	// API Routing (Assignment 5)
	api := router.Group("/api")
	routers.RegisterTaskRoutes(api)
	routers.RegisterUserRoutes(api)
	routers.RegisterAdminRoutes(api)

	// Default routes from previous assignments
	router.GET("/", handlers.Index)
	router.POST("/calculate", handlers.Calculate)
	router.GET("/users", handlers.GetUsers)
	router.POST("/user", handlers.CreateUser)
	router.POST("/create_user", handlers.CreateUser2)
	router.POST("/feedback", handlers.Feedback)
	router.GET("/headers", handlers.Headers)
	router.GET("/info", handlers.Info)
	router.GET("/product/:productID", handlers.GetProduct)
	router.GET("/products/search", handlers.SearchProducts)

	// Assignment 3: Docs control with environment
	router.GET("/docs", middleware.DocsAuth(), handlers.Docs)

	// Assignment 3: Auth & Register with Rate Limiting and SQLite
	router.POST("/register", middleware.RateLimit(1, 60), handlers.Register)
	router.POST("/login", middleware.RateLimit(5, 60), handlers.Login)

	// Assignment 3: Todos CRUD with SQLite
	todos := router.Group("/todos")
	{
		todos.POST("", handlers.CreateTodo)
		todos.GET("/:id", handlers.GetTodo)
		todos.PUT("/:id", handlers.UpdateTodo)
		todos.DELETE("/:id", handlers.DeleteTodo)
	}

	// Trigger manual errors (Task 10.1)
	router.GET("/errorA", handlers.TriggerErrorA)
	router.GET("/errorB", handlers.TriggerErrorB)

	// Trigger validation errors (Task 10.2)
	router.POST("/validate_user", handlers.ValidateUser)

	// Assignment 4: Unit Test Mem Users
	memUsers := router.Group("/mem_users")
	{
		memUsers.POST("", handlers.CreateMemUser)
		memUsers.GET("/:id", handlers.GetMemUser)
		memUsers.DELETE("/:id", handlers.DeleteMemUser)
	}

	// Assignment 3: RBAC (Role-Based Access Control)
	protected := router.Group("/")
	protected.Use(middleware.JWTAuth())
	{
		// /protected_resource
		protected.GET("/protected_resource", middleware.RoleAuth("admin", "user"), handlers.ProtectedResource)

		// Admin only
		protected.GET("/admin_resource", middleware.RoleAuth("admin"), handlers.AdminResource)

		// User or Admin
		protected.GET("/user_resource", middleware.RoleAuth("admin", "user"), handlers.UserResource)
	}

	return router
}
