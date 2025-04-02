package controllers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"task_manager/domain"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MockTaskUsecase struct {
	mock.Mock
}

func (m *MockTaskUsecase) AddTask(ctx context.Context, task domain.Task) (string, error) {
	args := m.Called(ctx, task)
	return args.String(0), args.Error(1)
}

func (m *MockTaskUsecase) GetAllTasks(ctx context.Context) ([]domain.Task, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.Task), args.Error(1)
}

func (m *MockTaskUsecase) GetTaskByID(ctx context.Context, id string) (*domain.Task, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*domain.Task), args.Error(1)
}

func (m *MockTaskUsecase) UpdateTask(ctx context.Context, id string, task domain.Task) error {
	return m.Called(ctx, id, task).Error(0)
}

func (m *MockTaskUsecase) DeleteTask(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}

type MockUserUsecase struct {
	mock.Mock
}

func (m *MockUserUsecase) Register(ctx context.Context, user domain.User) error {
	return m.Called(ctx, user).Error(0)
}

func (m *MockUserUsecase) Login(ctx context.Context, username, password string) (string, error) {
	args := m.Called(ctx, username, password)
	return args.String(0), args.Error(1)
}

func (m *MockUserUsecase) PromoteUser(ctx context.Context, username string) error {
	return m.Called(ctx, username).Error(0)
}

func (m *MockUserUsecase) GetAllUsers(ctx context.Context) ([]*domain.User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*domain.User), args.Error(1)
}

type ControllerTestSuite struct {
	suite.Suite
	mockTaskUsecase *MockTaskUsecase
	mockUserUsecase *MockUserUsecase
	taskController  *TaskController
	userController  *UserController
	router          *gin.Engine
}

func (s *ControllerTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
	s.mockTaskUsecase = &MockTaskUsecase{}
	s.mockUserUsecase = &MockUserUsecase{}
	s.taskController = NewTaskController(s.mockTaskUsecase)
	s.userController = NewUserController(s.mockUserUsecase)
	s.router = gin.New()
	// Match routes to your controller.go
	s.router.POST("/tasks", s.taskController.AddTask)
	s.router.GET("/tasks", s.taskController.GetTasks)
	s.router.GET("/tasks/:id", s.taskController.GetTask)
	s.router.PUT("/tasks/:id", s.taskController.UpdateTask)
	s.router.DELETE("/tasks/:id", s.taskController.RemoveTask)
	s.router.POST("/register", s.userController.Register)
	s.router.POST("/login", s.userController.Login)
	s.router.PUT("/promote", s.userController.PromoteUser) // Uses JSON body
	s.router.GET("/users", s.userController.GetAllUsers)
}

func (s *ControllerTestSuite) TearDownTest() {
	s.mockTaskUsecase.AssertExpectations(s.T())
	s.mockUserUsecase.AssertExpectations(s.T())
}

func (s *ControllerTestSuite) TestAddTask() {
	s.Run("Success", func() {
		taskJSON := `{"title":"Test Task","due_date":"2025-04-03T00:00:00Z","status":"pending"}`
		s.mockTaskUsecase.On("AddTask", mock.Anything, mock.AnythingOfType("domain.Task")).Return("1", nil).Once()

		req, _ := http.NewRequest("POST", "/tasks", strings.NewReader(taskJSON))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		s.router.ServeHTTP(w, req)

		s.Equal(http.StatusCreated, w.Code)
		s.Contains(w.Body.String(), `"id":"1"`)
		s.Contains(w.Body.String(), `"message":"Task created"`)
	})
}

func (s *ControllerTestSuite) TestGetTasks() {
	s.Run("Success", func() {
		tasks := []domain.Task{{ID: "1", Title: "Task 1"}}
		s.mockTaskUsecase.On("GetAllTasks", mock.Anything).Return(tasks, nil).Once()

		req, _ := http.NewRequest("GET", "/tasks", nil)
		w := httptest.NewRecorder()
		s.router.ServeHTTP(w, req)

		s.Equal(http.StatusOK, w.Code)
		s.Contains(w.Body.String(), `"id":"1"`)
		s.Contains(w.Body.String(), `"title":"Task 1"`)
	})
}

func (s *ControllerTestSuite) TestGetTask() {
	s.Run("Success", func() {
		task := &domain.Task{ID: "1", Title: "Test Task"}
		s.mockTaskUsecase.On("GetTaskByID", mock.Anything, "1").Return(task, nil).Once()

		req, _ := http.NewRequest("GET", "/tasks/1", nil)
		w := httptest.NewRecorder()
		s.router.ServeHTTP(w, req)

		s.Equal(http.StatusOK, w.Code)
		s.Contains(w.Body.String(), `"id":"1"`)
		s.Contains(w.Body.String(), `"title":"Test Task"`)
	})
}

func (s *ControllerTestSuite) TestUpdateTask() {
	s.Run("Success", func() {
		taskJSON := `{"title":"Updated Task","due_date":"2025-04-03T00:00:00Z","status":"done"}`
		s.mockTaskUsecase.On("UpdateTask", mock.Anything, "1", mock.AnythingOfType("domain.Task")).Return(nil).Once()

		req, _ := http.NewRequest("PUT", "/tasks/1", strings.NewReader(taskJSON))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		s.router.ServeHTTP(w, req)

		s.Equal(http.StatusOK, w.Code)
		s.Contains(w.Body.String(), `"message":"Task updated"`)
	})
}

func (s *ControllerTestSuite) TestRemoveTask() {
	s.Run("Success", func() {
		s.mockTaskUsecase.On("DeleteTask", mock.Anything, "1").Return(nil).Once()

		req, _ := http.NewRequest("DELETE", "/tasks/1", nil)
		w := httptest.NewRecorder()
		s.router.ServeHTTP(w, req)

		s.Equal(http.StatusOK, w.Code)
		s.Contains(w.Body.String(), `"message":"Task removed"`)
	})
}

func (s *ControllerTestSuite) TestRegister() {
	s.Run("Success", func() {
		userJSON := `{"username":"testuser","password":"pass"}`
		s.mockUserUsecase.On("Register", mock.Anything, mock.AnythingOfType("domain.User")).Return(nil).Once()

		req, _ := http.NewRequest("POST", "/register", strings.NewReader(userJSON))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		s.router.ServeHTTP(w, req)

		s.Equal(http.StatusCreated, w.Code)
		s.Contains(w.Body.String(), `"message":"User created"`)
	})
}

func (s *ControllerTestSuite) TestLogin() {
	s.Run("Success", func() {
		credsJSON := `{"username":"testuser","password":"pass"}`
		s.mockUserUsecase.On("Login", mock.Anything, "testuser", "pass").Return("token", nil).Once()

		req, _ := http.NewRequest("POST", "/login", strings.NewReader(credsJSON))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		s.router.ServeHTTP(w, req)

		s.Equal(http.StatusOK, w.Code)
		s.Contains(w.Body.String(), `"token":"token"`)
	})
}

func (s *ControllerTestSuite) TestPromoteUser() {
	s.Run("Success", func() {
		promoteJSON := `{"username":"testuser"}`
		s.mockUserUsecase.On("PromoteUser", mock.Anything, "testuser").Return(nil).Once()

		req, _ := http.NewRequest("PUT", "/promote", strings.NewReader(promoteJSON))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		s.router.ServeHTTP(w, req)

		s.Equal(http.StatusOK, w.Code)
		s.Contains(w.Body.String(), `"message":"User promoted to admin"`)
	})
}

func (s *ControllerTestSuite) TestGetAllUsers() {
	s.Run("Success", func() {
		users := []*domain.User{{ID: "1", Username: "user1"}}
		s.mockUserUsecase.On("GetAllUsers", mock.Anything).Return(users, nil).Once()

		req, _ := http.NewRequest("GET", "/users", nil)
		w := httptest.NewRecorder()
		s.router.ServeHTTP(w, req)

		s.Equal(http.StatusOK, w.Code)
		s.Contains(w.Body.String(), `"username":"user1"`)
	})
}

func TestControllerSuite(t *testing.T) {
	suite.Run(t, new(ControllerTestSuite))
}