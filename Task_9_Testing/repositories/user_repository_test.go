package repositories

import (
	"context"
	"testing"
	"task_manager/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
	return m.Called(ctx, id).Error(0)
}

func TestCreateUser(t *testing.T) {
	mockRepo := &MockUserRepository{}
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		user := domain.User{
			ID:       "1",
			Username: "testuser",
			Password: "hashed",
			Role:     "user",
		}
		mockRepo.On("CreateUser", ctx, user).Return(nil).Once()

		err := mockRepo.CreateUser(ctx, user)
		assert.NoError(t, err, "CreateUser should succeed")

		mockRepo.AssertExpectations(t)
	})
}

func TestFindUserByUsername(t *testing.T) {
	mockRepo := &MockUserRepository{}
	ctx := context.Background()

	t.Run("UserFound", func(t *testing.T) {
		user := &domain.User{
			ID:       "1",
			Username: "testuser",
			Password: "hashed",
			Role:     "user",
		}
		mockRepo.On("FindUserByUsername", ctx, "testuser").Return(user, nil).Once()

		result, err := mockRepo.FindUserByUsername(ctx, "testuser")
		assert.NoError(t, err, "FindUserByUsername should succeed")
		assert.NotNil(t, result, "Result should not be nil")
		assert.Equal(t, "testuser", result.Username, "Username should match")

		mockRepo.AssertExpectations(t)
	})

	t.Run("UserNotFound", func(t *testing.T) {
		mockRepo.On("FindUserByUsername", ctx, "unknown").Return((*domain.User)(nil), nil).Once()

		result, err := mockRepo.FindUserByUsername(ctx, "unknown")
		assert.NoError(t, err, "FindUserByUsername should succeed")
		assert.Nil(t, result, "Result should be nil when user not found")

		mockRepo.AssertExpectations(t)
	})
}

func TestPromoteUser(t *testing.T) {
	mockRepo := &MockUserRepository{}
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("PromoteUser", ctx, "testuser").Return(nil).Once()

		err := mockRepo.PromoteUser(ctx, "testuser")
		assert.NoError(t, err, "PromoteUser should succeed")

		mockRepo.AssertExpectations(t)
	})
}

func TestIsFirstUser(t *testing.T) {
	mockRepo := &MockUserRepository{}
	ctx := context.Background()

	t.Run("IsFirst", func(t *testing.T) {
		mockRepo.On("IsFirstUser", ctx).Return(true, nil).Once()

		isFirst, err := mockRepo.IsFirstUser(ctx)
		assert.NoError(t, err, "IsFirstUser should succeed")
		assert.True(t, isFirst, "Should return true for first user")

		mockRepo.AssertExpectations(t)
	})

	t.Run("NotFirst", func(t *testing.T) {
		mockRepo.On("IsFirstUser", ctx).Return(false, nil).Once()

		isFirst, err := mockRepo.IsFirstUser(ctx)
		assert.NoError(t, err, "IsFirstUser should succeed")
		assert.False(t, isFirst, "Should return false when not first user")

		mockRepo.AssertExpectations(t)
	})
}

func TestGetAllUsers(t *testing.T) {
	mockRepo := &MockUserRepository{}
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		users := []*domain.User{
			{ID: "1", Username: "user1", Password: "hashed1", Role: "user"},
			{ID: "2", Username: "user2", Password: "hashed2", Role: "admin"},
		}
		mockRepo.On("GetAllUsers", ctx).Return(users, nil).Once()

		result, err := mockRepo.GetAllUsers(ctx)
		assert.NoError(t, err, "GetAllUsers should succeed")
		assert.Len(t, result, 2, "Should return two users")
		assert.Equal(t, "user1", result[0].Username, "First user’s username should match")
		assert.Equal(t, "user2", result[1].Username, "Second user’s username should match")

		mockRepo.AssertExpectations(t)
	})
}

func TestDeleteUser(t *testing.T) {
	mockRepo := &MockUserRepository{}
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("DeleteUser", ctx, "1").Return(nil).Once()

		err := mockRepo.DeleteUser(ctx, "1")
		assert.NoError(t, err, "DeleteUser should succeed")

		mockRepo.AssertExpectations(t)
	})
}