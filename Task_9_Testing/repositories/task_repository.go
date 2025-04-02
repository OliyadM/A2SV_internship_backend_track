package repositories

import (
	"context"
	"task_manager/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type TaskRepositoryImpl struct {
	collection *mongo.Collection
}

func NewTaskRepository(collection *mongo.Collection) domain.TaskRepository {
	return &TaskRepositoryImpl{collection: collection}
}

func (r *TaskRepositoryImpl) AddTask(ctx context.Context, task domain.Task) (string, error) {
	result, err := r.collection.InsertOne(ctx, task)
	if err != nil {
		return "", err
	}
	return result.InsertedID.(string), nil
}

func (r *TaskRepositoryImpl) GetAllTasks(ctx context.Context) ([]domain.Task, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var tasks []domain.Task
	err = cursor.All(ctx, &tasks)
	return tasks, err
}

func (r *TaskRepositoryImpl) GetTaskByID(ctx context.Context, id string) (*domain.Task, error) {
	var task domain.Task
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&task)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *TaskRepositoryImpl) UpdateTask(ctx context.Context, id string, task domain.Task) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": task})
	return err
}

func (r *TaskRepositoryImpl) DeleteTask(ctx context.Context, id string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
