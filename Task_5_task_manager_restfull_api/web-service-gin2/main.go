package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"time"
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

func main(){
	router := gin.Default()
	router.GET("/Tasks",getTasks)
	router.GET("/Tasks/:id",getTaskBYID)
	router.POST("/Tasks",postTask)
	router.DELETE("/Tasks/:id",deleteTask)
	router.PUT("/Tasks/",updateTask)
}

func getTasks(c*gin.Context){
	c.IndentedJSON(http.StatusOK, tasks)
}

func getTaskBYID(c*gin.Context){
	var id  = c.Params("id")
	for _ ,a := range tasks{
		
	}
	 
}

func postTask(c*gin.Context){
	 
}

func updateTask(c*gin.Context){

}

func deleteTask(c*gin.Context){

}

