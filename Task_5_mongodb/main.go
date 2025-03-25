package main

import (
	"task_manager/data"
	"task_manager/router"
)

func main() {
	
	data.ConnectDB()

	r := router.SetupRouter()
	r.Run(":8080")
}
