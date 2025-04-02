package usecases

import (
	"context"
	"task_manager/domain"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateUser(ctx context.Context, user domain.User) error {
	return m.Called(ctx, user).Error(0)
}

func (m *MockUserRepository) FindUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	args := m.Called(ctx, username)
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) PromoteUser(ctx context.Context, username string) error {
	return m.Called(ctx, username).Error(0)
}

func (m *MockUserRepository) IsFirstUser(ctx context.Context) (bool, error) {
	args := m.Called(ctx)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) GetAllUsers(ctx context.Context) ([]*domain.User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*domain.User), args.Error(1)
}

func (m *MockUserRepository) DeleteUser(ctx context.Context, id string) error {
	return nil
}

type MockPasswordService struct {
	mock.Mock
}

func (m *MockPasswordService) HashPassword(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

func (m *MockPasswordService) ComparePassword(hashed, plain string) error {
	return m.Called(hashed, plain).Error(0)
}

type MockJWTService struct {
	mock.Mock
}

func (m *MockJWTService) GenerateToken(userID, username, role string) (string, error) {
	args := m.Called(userID, username, role)
	return args.String(0), args.Error(1)
}

func (m *MockJWTService) ValidateToken(tokenString string) (*jwt.Token, error) {
	return nil, nil
}

type UserUsecaseTestSuite struct {
	suite.Suite
	mockRepo *MockUserRepository
	mockPass *MockPasswordService
	mockJWT  *MockJWTService
	usecase  domain.UserUsecase
	ctx      context.Context
}

func (s *UserUsecaseTestSuite) SetupTest() {
	s.mockRepo = &MockUserRepository{}
	s.mockPass = &MockPasswordService{}
	s.mockJWT = &MockJWTService{}
	s.usecase = NewUserUsecase(s.mockRepo, s.mockPass, s.mockJWT)
	s.ctx = context.Background()
}

func (s *UserUsecaseTestSuite) TearDownTest() {
	s.mockRepo.AssertExpectations(s.T())
	s.mockPass.AssertExpectations(s.T())
	s.mockJWT.AssertExpectations(s.T())
}

func (s *UserUsecaseTestSuite) TestRegister() {
	s.Run("SuccessFirstUser", func() {
		user := domain.User{Username: "testuser", Password: "plain"}
		s.mockRepo.On("FindUserByUsername", s.ctx, "testuser").Return((*domain.User)(nil), nil).Once()
		s.mockPass.On("HashPassword", "plain").Return("hashed", nil).Once()
		s.mockRepo.On("IsFirstUser", s.ctx).Return(true, nil).Once()
		s.mockRepo.On("CreateUser", s.ctx, mock.Anything).Return(nil).Once()

		err := s.usecase.Register(s.ctx, user)
		s.NoError(err)
	})

	s.Run("SuccessNotFirstUser", func() {
		user := domain.User{Username: "testuser2", Password: "plain"}
		s.mockRepo.On("FindUserByUsername", s.ctx, "testuser2").Return((*domain.User)(nil), nil).Once()
		s.mockPass.On("HashPassword", "plain").Return("hashed", nil).Once()
		s.mockRepo.On("IsFirstUser", s.ctx).Return(false, nil).Once()
		s.mockRepo.On("CreateUser", s.ctx, mock.Anything).Return(nil).Once()

		err := s.usecase.Register(s.ctx, user)
		s.NoError(err)
	})

	s.Run("UsernameExists", func() {
		user := domain.User{Username: "testuser", Password: "plain"}
		existing := &domain.User{Username: "testuser"}
		s.mockRepo.On("FindUserByUsername", s.ctx, "testuser").Return(existing, nil).Once()

		err := s.usecase.Register(s.ctx, user)
		s.Error(err)
		s.Equal("username already exists", err.Error())
	})
}

func (s *UserUsecaseTestSuite) TestLogin() {
	s.Run("Success", func() {
		user := &domain.User{ID: "1", Username: "testuser", Password: "hashed", Role: "user"}
		s.mockRepo.On("FindUserByUsername", s.ctx, "testuser").Return(user, nil).Once()
		s.mockPass.On("ComparePassword", "hashed", "plain").Return(nil).Once()
		s.mockJWT.On("GenerateToken", "1", "testuser", "user").Return("token", nil).Once()

		token, err := s.usecase.Login(s.ctx, "testuser", "plain")
		s.NoError(err)
		s.Equal("token", token)
	})

	s.Run("InvalidCredentials", func() {
		s.mockRepo.On("FindUserByUsername", s.ctx, "testuser").Return((*domain.User)(nil), nil).Once()

		token, err := s.usecase.Login(s.ctx, "testuser", "plain")
		s.Error(err)
		s.Equal("invalid credentials", err.Error())
		s.Empty(token)
	})
}

func (s *UserUsecaseTestSuite) TestPromoteUser() {
	s.Run("Success", func() {
		user := &domain.User{Username: "testuser"}
		s.mockRepo.On("FindUserByUsername", s.ctx, "testuser").Return(user, nil).Once()
		s.mockRepo.On("PromoteUser", s.ctx, "testuser").Return(nil).Once()

		err := s.usecase.PromoteUser(s.ctx, "testuser")
		s.NoError(err)
	})

	s.Run("UserNotFound", func() {
		s.mockRepo.On("FindUserByUsername", s.ctx, "unknown").Return((*domain.User)(nil), nil).Once()

		err := s.usecase.PromoteUser(s.ctx, "unknown")
		s.Error(err)
		s.Equal("user not found", err.Error())
	})
}

func (s *UserUsecaseTestSuite) TestGetAllUsers() {
	s.Run("Success", func() {
		users := []*domain.User{
			{ID: "1", Username: "user1"},
			{ID: "2", Username: "user2"},
		}
		s.mockRepo.On("GetAllUsers", s.ctx).Return(users, nil).Once()

		result, err := s.usecase.GetAllUsers(s.ctx)
		s.NoError(err)
		s.Len(result, 2)
		s.Equal("user1", result[0].Username)
	})

	s.Run("NoUsers", func() {
		s.mockRepo.On("GetAllUsers", s.ctx).Return([]*domain.User{}, nil).Once()

		result, err := s.usecase.GetAllUsers(s.ctx)
		s.Error(err)
		s.Equal("no users found", err.Error())
		s.Nil(result)
	})
}

func TestUserUsecaseSuite(t *testing.T) {
	suite.Run(t, new(UserUsecaseTestSuite))
}