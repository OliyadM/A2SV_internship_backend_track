package repositories

import (
	"context"
	"testing"
	"time"
	"task_manager/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTaskRepository struct {
	mock.Mock
}

func (m *MockTaskRepository) AddTask(ctx context.Context, task domain.Task) (string, error) {
	args := m.Called(ctx, task)
	return args.String(0), args.Error(1)
}

func (m *MockTaskRepository) GetAllTasks(ctx context.Context) ([]domain.Task, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.Task), args.Error(1)
}

func (m *MockTaskRepository) GetTaskByID(ctx context.Context, id string) (*domain.Task, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*domain.Task), args.Error(1)
}

func (m *MockTaskRepository) UpdateTask(ctx context.Context, id string, task domain.Task) error {
	return m.Called(ctx, id, task).Error(0)
}

func (m *MockTaskRepository) DeleteTask(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}

func TestAddTask(t *testing.T) {
	mockRepo := &MockTaskRepository{}
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		task := domain.Task{
			ID:      "1",
			Title:   "Test Task",
			DueDate: time.Now(),
			Status:  "pending",
		}
		mockRepo.On("AddTask", ctx, task).Return("1", nil).Once()

		id, err := mockRepo.AddTask(ctx, task)
		assert.NoError(t, err, "AddTask should succeed")
		assert.Equal(t, "1", id, "AddTask should return the correct ID")

		mockRepo.AssertExpectations(t)
	})
}

func TestGetAllTasks(t *testing.T) {
	mockRepo := &MockTaskRepository{}
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		tasks := []domain.Task{
			{ID: "1", Title: "Task 1", DueDate: time.Now(), Status: "pending"},
			{ID: "2", Title: "Task 2", DueDate: time.Now(), Status: "done"},
		}
		mockRepo.On("GetAllTasks", ctx).Return(tasks, nil).Once()

		result, err := mockRepo.GetAllTasks(ctx)
		assert.NoError(t, err, "GetAllTasks should succeed")
		assert.Len(t, result, 2, "Should return two tasks")
		assert.Equal(t, "Task 1", result[0].Title, "First task title should match")
		assert.Equal(t, "Task 2", result[1].Title, "Second task title should match")

		mockRepo.AssertExpectations(t)
	})
}

func TestGetTaskByID(t *testing.T) {
	mockRepo := &MockTaskRepository{}
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		task := &domain.Task{
			ID:      "1",
			Title:   "Test Task",
			DueDate: time.Now(),
			Status:  "pending",
		}
		mockRepo.On("GetTaskByID", ctx, "1").Return(task, nil).Once()

		result, err := mockRepo.GetTaskByID(ctx, "1")
		assert.NoError(t, err, "GetTaskByID should succeed")
		assert.NotNil(t, result, "Result should not be nil")
		assert.Equal(t, "Test Task", result.Title, "Task title should match")

		mockRepo.AssertExpectations(t)
	})
}

func TestUpdateTask(t *testing.T) {
	mockRepo := &MockTaskRepository{}
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		task := domain.Task{
			ID:      "1",
			Title:   "Updated Task",
			DueDate: time.Now(),
			Status:  "done",
		}
		mockRepo.On("UpdateTask", ctx, "1", task).Return(nil).Once()

		err := mockRepo.UpdateTask(ctx, "1", task)
		assert.NoError(t, err, "UpdateTask should succeed")

		mockRepo.AssertExpectations(t)
	})
}

func TestDeleteTask(t *testing.T) {
	mockRepo := &MockTaskRepository{}
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("DeleteTask", ctx, "1").Return(nil).Once()

		err := mockRepo.DeleteTask(ctx, "1")
		assert.NoError(t, err, "DeleteTask should succeed")

		mockRepo.AssertExpectations(t)
	})
}