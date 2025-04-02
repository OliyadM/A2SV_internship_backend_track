package usecases

import (
	"context"
	"task_manager/domain"

	"github.com/google/uuid"
)

type TaskUsecaseImpl struct {
	taskRepo domain.TaskRepository
}

func NewTaskUsecase(taskRepo domain.TaskRepository) domain.TaskUsecase {
	return &TaskUsecaseImpl{taskRepo: taskRepo}
}

func (u *TaskUsecaseImpl) AddTask(ctx context.Context, task domain.Task) (string, error) {
	
	task.ID = uuid.New().String()

	return u.taskRepo.AddTask(ctx, task)
}

func (u *TaskUsecaseImpl) GetAllTasks(ctx context.Context) ([]domain.Task, error) {
	return u.taskRepo.GetAllTasks(ctx)
}

func (u *TaskUsecaseImpl) GetTaskByID(ctx context.Context, id string) (*domain.Task, error) {
	return u.taskRepo.GetTaskByID(ctx, id)
}

func (u *TaskUsecaseImpl) UpdateTask(ctx context.Context, id string, task domain.Task) error {
	return u.taskRepo.UpdateTask(ctx, id, task)
}

func (u *TaskUsecaseImpl) DeleteTask(ctx context.Context, id string) error {
	return u.taskRepo.DeleteTask(ctx, id)
}
