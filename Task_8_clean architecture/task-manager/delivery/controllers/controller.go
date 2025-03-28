package controllers

import (
	"net/http"
	"task_manager/domain"

	"github.com/gin-gonic/gin"
)

type TaskController struct {
	taskUsecase domain.TaskUsecase
}

func NewTaskController(taskUsecase domain.TaskUsecase) *TaskController {
	return &TaskController{taskUsecase: taskUsecase}
}

func (ctrl *TaskController) GetTasks(c *gin.Context) {
	tasks, err := ctrl.taskUsecase.GetAllTasks(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching tasks"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"tasks": tasks})
}

func (ctrl *TaskController) GetTask(c *gin.Context) {
	id := c.Param("id")
	task, err := ctrl.taskUsecase.GetTaskByID(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Task not found"})
		return
	}
	c.JSON(http.StatusOK, task)
}

func (ctrl *TaskController) AddTask(c *gin.Context) {
	var task domain.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id, err := ctrl.taskUsecase.AddTask(c, task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Task created", "id": id})
}

func (ctrl *TaskController) UpdateTask(c *gin.Context) {
	id := c.Param("id")
	var task domain.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := ctrl.taskUsecase.UpdateTask(c, id, task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating task"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Task updated"})
}

func (ctrl *TaskController) RemoveTask(c *gin.Context) {
	id := c.Param("id")
	if err := ctrl.taskUsecase.DeleteTask(c, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting task"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Task removed"})
}

type UserController struct {
	userUsecase domain.UserUsecase
}

func NewUserController(userUsecase domain.UserUsecase) *UserController {
	return &UserController{userUsecase: userUsecase}
}

func (ctrl *UserController) Register(c *gin.Context) {
	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := ctrl.userUsecase.Register(c, user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "User created"})
}

func (ctrl *UserController) Login(c *gin.Context) {
	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, err := ctrl.userUsecase.Login(c, creds.Username, creds.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (ctrl *UserController) PromoteUser(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := ctrl.userUsecase.PromoteUser(c, req.Username); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User promoted to admin"})
}

// In controllers/user_controller.go
func (ctrl *UserController) GetAllUsers(c *gin.Context) {
	users, err := ctrl.userUsecase.GetAllUsers(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}
