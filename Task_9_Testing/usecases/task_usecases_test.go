package usecases

import (
	"errors"
	"context"
	"testing"
	"time"
	"task_manager/domain"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
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

type TaskUsecaseTestSuite struct {
	suite.Suite
	mockRepo *MockTaskRepository
	usecase  domain.TaskUsecase
	ctx      context.Context
}

func (s *TaskUsecaseTestSuite) SetupTest() {
	s.mockRepo = &MockTaskRepository{}
	s.usecase = NewTaskUsecase(s.mockRepo)
	s.ctx = context.Background()
}

func (s *TaskUsecaseTestSuite) TearDownTest() {
	s.mockRepo.AssertExpectations(s.T())
}

func (s *TaskUsecaseTestSuite) TestAddTask() {
	s.Run("Success", func() {
		task := domain.Task{Title: "Test Task", DueDate: time.Now(), Status: "pending"}
		s.mockRepo.On("AddTask", s.ctx, mock.Anything).Return("1", nil).Once()

		id, err := s.usecase.AddTask(s.ctx, task)
		s.NoError(err)
		s.Equal("1", id)
	})
}

func (s *TaskUsecaseTestSuite) TestGetAllTasks() {
	s.Run("Success", func() {
		tasks := []domain.Task{
			{ID: "1", Title: "Task 1", DueDate: time.Now(), Status: "pending"},
			{ID: "2", Title: "Task 2", DueDate: time.Now(), Status: "done"},
		}
		s.mockRepo.On("GetAllTasks", s.ctx).Return(tasks, nil).Once()

		result, err := s.usecase.GetAllTasks(s.ctx)
		s.NoError(err)
		s.Len(result, 2)
		s.Equal("Task 1", result[0].Title)
	})
}

func (s *TaskUsecaseTestSuite) TestGetTaskByID() {
	s.Run("Success", func() {
		task := &domain.Task{ID: "1", Title: "Test Task", DueDate: time.Now(), Status: "pending"}
		s.mockRepo.On("GetTaskByID", s.ctx, "1").Return(task, nil).Once()

		result, err := s.usecase.GetTaskByID(s.ctx, "1")
		s.NoError(err)
		s.NotNil(result)
		s.Equal("Test Task", result.Title)
	})

	s.Run("NotFound", func() {
		s.mockRepo.On("GetTaskByID", s.ctx, "1").Return((*domain.Task)(nil), errors.New("not found")).Once()

		result, err := s.usecase.GetTaskByID(s.ctx, "1")
		s.Error(err)
		s.Nil(result)
	})
}

func (s *TaskUsecaseTestSuite) TestUpdateTask() {
	s.Run("Success", func() {
		task := domain.Task{Title: "Updated Task", DueDate: time.Now(), Status: "done"}
		s.mockRepo.On("UpdateTask", s.ctx, "1", task).Return(nil).Once()

		err := s.usecase.UpdateTask(s.ctx, "1", task)
		s.NoError(err)
	})
}

func (s *TaskUsecaseTestSuite) TestDeleteTask() {
	s.Run("Success", func() {
		s.mockRepo.On("DeleteTask", s.ctx, "1").Return(nil).Once()

		err := s.usecase.DeleteTask(s.ctx, "1")
		s.NoError(err)
	})
}

func TestTaskUsecaseSuite(t *testing.T) {
	suite.Run(t, new(TaskUsecaseTestSuite))
}