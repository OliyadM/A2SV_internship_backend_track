package router

import (
    "task_manager/controllers"
    "task_manager/middleware"

    "github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
    router := gin.Default()

    // Public routes
    router.POST("/register", controllers.Register)
    router.POST("/login", controllers.Login)

    // Authenticated routes
    auth := router.Group("/")
    auth.Use(middleware.AuthMiddleware())
    {
        auth.GET("/tasks", controllers.GetTasks)
        auth.GET("/tasks/:id", controllers.GetTask)
    }

    // Admin-only routes
    admin := router.Group("/")
    admin.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
    {
        admin.POST("/tasks", controllers.AddTask)
        admin.PUT("/tasks/:id", controllers.UpdateTask)
        admin.DELETE("/tasks/:id", controllers.RemoveTask)
        admin.POST("/promote", controllers.PromoteUser)
    }

    return router
}