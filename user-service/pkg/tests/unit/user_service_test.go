package service_test

import (
	"testing"
	"time"
	"user-service/pkg/model"
	"user-service/pkg/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) UserExists(username string) bool {
	args := m.Called(username)
	return args.Bool(0)
}

func (m *MockUserRepo) CreateUser(user *model.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepo) DeleteUser(userID string) error {
	args := m.Called(userID)
	return args.Error(0)
}

func (m *MockUserRepo) GetAllUsers() ([]model.User, error) {
	args := m.Called()
	return args.Get(0).([]model.User), args.Error(1)
}

type MockPublisher struct {
	mock.Mock
}

func (m *MockPublisher) Publish(subject string, data []byte) error {
	args := m.Called(subject, data)
	return args.Error(0)
}

type MockCache struct {
	mock.Mock
}

func (m *MockCache) Set(key string, value string, ttl time.Duration) error {
	args := m.Called(key, value, ttl)
	return args.Error(0)
}

func (m *MockCache) Del(key string) error {
	args := m.Called(key)
	return args.Error(0)
}

func (m *MockCache) Get(key string) (string, error) {
	args := m.Called(key)
	return args.String(0), args.Error(1)
}
func (m *MockUserRepo) GetUserByEmail(email string) (*model.User, error) {
	args := m.Called(email)
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepo) GetUserByID(userID string) (*model.User, error) {
	args := m.Called(userID)
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepo) GetUserByUsername(username string) (*model.User, error) {
	args := m.Called(username)
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepo) UpdatePassword(userID, newPassword string) error {
	args := m.Called(userID, newPassword)
	return args.Error(0)
}

func (m *MockUserRepo) UpdateUser(userID, username, email string) error {
	args := m.Called(userID, username, email)
	return args.Error(0)
}




func TestUserService_Register_Success(t *testing.T) {
	repo := &MockUserRepo{}
	cache := &MockCache{}
	nats := &MockPublisher{}

	svc := service.NewUserService(repo, nats, cache)

	username := "testuser"
	email := "maksatkarzhaubaev91@gmail.com"
	password := "password123"

	// Настраиваем мок Repo
	repo.On("UserExists", username).Return(false)
	repo.On("CreateUser", mock.AnythingOfType("*model.User")).Return(nil)

	// Настраиваем мок NATS
	nats.On("Publish", "user.registered", mock.Anything).Return(nil)

	userID, err := svc.Register(username, password, email)

	assert.NoError(t, err)
	assert.NotEmpty(t, userID)

	repo.AssertExpectations(t)
	nats.AssertExpectations(t)
}
