package main

import (
	"context"
	"log"
	"task_manager/delivery/controllers"
	"task_manager/delivery/routers"
	"task_manager/infrastructure"
	"task_manager/repositories"
	"task_manager/usecases"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal("MongoDB connection failed:", err)
	}
	db := client.Database("task_maanager")
	taskCollection := db.Collection("tasks")
	userCollection := db.Collection("users")

	
	taskRepo := repositories.NewTaskRepository(taskCollection)
	userRepo := repositories.NewUserRepository(userCollection)
	passwordSvc := infrastructure.NewPasswordService()
	jwtSvc := infrastructure.NewJWTService("oliyads-secrete-jwt")
	taskUsecase := usecases.NewTaskUsecase(taskRepo)
	userUsecase := usecases.NewUserUsecase(userRepo, passwordSvc, jwtSvc)


	taskCtrl := controllers.NewTaskController(taskUsecase)
	userCtrl := controllers.NewUserController(userUsecase)

	router := routers.SetupRouter(taskCtrl, userCtrl, jwtSvc)

	if err := router.Run(":8080"); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
