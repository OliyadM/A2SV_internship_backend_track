package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Task struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"due_date"`
	Status      string    `json:"status"`
}

var tasks = []Task{
	{ID: "1", Title: "Task 1", Description: "First task", DueDate: time.Now(), Status: "Pending"},
	{ID: "2", Title: "Task 2", Description: "Second task", DueDate: time.Now().AddDate(0, 0, 1), Status: "In Progress"},
	{ID: "3", Title: "Task 3", Description: "Third task", DueDate: time.Now().AddDate(0, 0, 2), Status: "Completed"},
}

func main() {
	router := gin.Default()
	router.GET("/tasks", getTasks)
	router.GET("/tasks/:id", getTask)
	router.POST("/tasks", addTask)
	router.PUT("/tasks/:id", updateTask)
	router.DELETE("/tasks/:id", removeTask)
	router.Run(":8080")
}

func getTasks(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"tasks": tasks})
}

func getTask(c *gin.Context) {
	id := c.Param("id")
	for _, val := range tasks {
		if val.ID == id {
			c.JSON(http.StatusOK, val)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "Task not found"})
}

func addTask(c *gin.Context) {
	var newTask Task
	if err := c.ShouldBindJSON(&newTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newTask.ID = time.Now().Format("20060102150405")
	tasks = append(tasks, newTask)
	c.JSON(http.StatusCreated, gin.H{"message": "Task created", "task": newTask})
}

func updateTask(c *gin.Context) {
	id := c.Param("id")
	var updatedTask Task
	if err := c.ShouldBindJSON(&updatedTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	for i, task := range tasks {
		if task.ID == id {
			if updatedTask.Title != "" {
				tasks[i].Title = updatedTask.Title
			}
			if updatedTask.Description != "" {
				tasks[i].Description = updatedTask.Description
			}
			if !updatedTask.DueDate.IsZero() {
				tasks[i].DueDate = updatedTask.DueDate
			}
			if updatedTask.Status != "" {
				tasks[i].Status = updatedTask.Status
			}
			c.JSON(http.StatusOK, gin.H{"message": "Task updated", "task": tasks[i]})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "Task not found"})
}

func removeTask(c *gin.Context) {
	id := c.Param("id")
	for i, val := range tasks {
		if val.ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			c.JSON(http.StatusOK, gin.H{"message": "Task removed"})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "Task not found"})
}
