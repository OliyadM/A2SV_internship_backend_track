package routers

import (
	"task_manager/delivery/controllers"
	"task_manager/domain" 

	"task_manager/infrastructure"

	"github.com/gin-gonic/gin"
)

func SetupRouter(taskCtrl *controllers.TaskController, userCtrl *controllers.UserController, jwtSvc domain.JWTService) *gin.Engine {
	router := gin.Default()


	router.POST("/register", userCtrl.Register)
	router.POST("/login", userCtrl.Login)
	router.GET("/users", userCtrl.GetAllUsers)


	auth := router.Group("/").Use(infrastructure.AuthMiddleware(jwtSvc))
	{
		auth.GET("/tasks", taskCtrl.GetTasks)
		auth.GET("/tasks/:id", taskCtrl.GetTask)
		auth.DELETE("/tasks/:id", taskCtrl.RemoveTask)
		auth.POST("/tasks", taskCtrl.AddTask)

	}


	admin := router.Group("/").Use(infrastructure.AuthMiddleware(jwtSvc), infrastructure.AdminMiddleware())
	{

		admin.PUT("/tasks/:id", taskCtrl.UpdateTask)

		admin.POST("/promote", userCtrl.PromoteUser)
	}

	return router
}
