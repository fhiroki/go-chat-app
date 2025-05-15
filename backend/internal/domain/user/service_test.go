package user

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepository is a mock implementation of the Repository interface
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(ctx context.Context, user User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockRepository) FindAll(ctx context.Context) ([]*User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*User), args.Error(1)
}

func (m *MockRepository) FindByID(ctx context.Context, id int) (*User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*User), args.Error(1)
}

func (m *MockRepository) Update(ctx context.Context, user User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockRepository) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockTxManager is a mock implementation of the TxManager interface
type MockTxManager struct {
	mock.Mock
}

func (m *MockTxManager) RunInTx(ctx context.Context, fn func(ctx context.Context) error) error {
	args := m.Called(ctx, fn)
	return args.Error(0)
}

func TestUserService_FindByID(t *testing.T) {
	// Setup
	mockRepo := new(MockRepository)
	mockTxManager := new(MockTxManager)
	service := NewUserService(mockRepo, mockTxManager)

	ctx := context.Background()
	expectedUser := &User{
		ID:        "1",
		Name:      "Test User",
		Email:     "test@example.com",
		AvatarURL: "https://example.com/avatar.jpg",
		CreatedAt: time.Now(),
	}

	// Expectations
	mockRepo.On("FindByID", ctx, 1).Return(expectedUser, nil)

	// Action
	user, err := service.FindByID(ctx, 1)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
	mockRepo.AssertExpectations(t)
}
