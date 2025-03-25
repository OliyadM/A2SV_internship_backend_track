package data

import (
	"context"
	"log"
	"time"

	"task_manager/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


var taskCollection *mongo.Collection


func ConnectDB() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal("MongoDB connection failed:", err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal("MongoDB ping failed:", err)
	}

	log.Println("Connected to MongoDB!")

	
	taskCollection = client.Database("task_manager").Collection("tasks")
}


func GetTaskCollection() *mongo.Collection {
	return taskCollection
}


func AddTask(task models.Task) (*mongo.InsertOneResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return taskCollection.InsertOne(ctx, task)
}


func GetAllTasks() ([]models.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := taskCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tasks []models.Task
	err = cursor.All(ctx, &tasks)
	return tasks, err
}


func GetTaskByID(taskID string) (*models.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var task models.Task
	err := taskCollection.FindOne(ctx, bson.M{"_id": taskID}).Decode(&task)
	if err != nil {
		return nil, err
	}

	return &task, nil
}

func UpdateTask(taskID string, updatedTask models.Task) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{"$set": updatedTask}

	_, err := taskCollection.UpdateOne(ctx, bson.M{"_id": taskID}, update)
	return err
}


func DeleteTask(taskID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := taskCollection.DeleteOne(ctx, bson.M{"_id": taskID})
	return err
}
