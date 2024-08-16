package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/wojciechpawlinow/usermanagement/internal/config"
	"github.com/wojciechpawlinow/usermanagement/internal/domain"
	"github.com/wojciechpawlinow/usermanagement/internal/domain/user"
	"github.com/wojciechpawlinow/usermanagement/pkg/logger"
	domainMock "github.com/wojciechpawlinow/usermanagement/tests/mocks/domain"
	repoMock "github.com/wojciechpawlinow/usermanagement/tests/mocks/infrastructure/database/mysql"
)

func TestCreate(t *testing.T) {
	t.Run("create user", func(t *testing.T) {
		mockRepo := new(repoMock.UserRepositoryMock)

		mockTimeProvider := new(domainMock.TimeProviderMock)
		mockTimeProvider.On("UtcNow").Return(time.Now())

		userSrv := NewUserService(mockRepo, mockTimeProvider)

		dto := &CreateUserDTO{
			ID:          domain.NewID(),
			Email:       "test@example.com",
			Password:    "admin123",
			FirstName:   "Test",
			LastName:    "Test",
			PhoneNumber: "1234567890",
			Addresses: []*CreateUserAddress{
				{
					Type:       1,
					Street:     "123 Main St",
					City:       "New York",
					State:      "NY",
					PostalCode: "10001",
					Country:    "USA",
				},
			},
		}

		mockRepo.On("Create", mock.Anything, mock.Anything, mock.Anything).Return(nil)

		err := userSrv.Create(context.Background(), dto)
		assert.NoError(t, err)
	})
	t.Run("email already exists", func(t *testing.T) {
		mockRepo := new(repoMock.UserRepositoryMock)

		mockTimeProvider := new(domainMock.TimeProviderMock)
		mockTimeProvider.On("UtcNow").Return(time.Now())

		userSrv := NewUserService(mockRepo, mockTimeProvider)

		dto := &CreateUserDTO{
			ID:          domain.NewID(),
			Email:       "test@example.com",
			Password:    "admin123",
			FirstName:   "Test",
			LastName:    "Test",
			PhoneNumber: "1234567890",
			Addresses: []*CreateUserAddress{
				{
					Type:       1,
					Street:     "Test",
					City:       "New York",
					State:      "NY",
					PostalCode: "",
					Country:    "USA",
				},
			},
		}

		mockRepo.On("Create", mock.Anything, mock.Anything, mock.Anything).Return(user.ErrEmailAlreadyExists)

		err := userSrv.Create(context.Background(), dto)
		assert.ErrorIs(t, err, user.ErrEmailAlreadyExists)
	})

	t.Run("repository error", func(t *testing.T) {
		cfg := config.Load()
		logger.Setup(cfg)

		mockRepo := new(repoMock.UserRepositoryMock)

		mockTimeProvider := new(domainMock.TimeProviderMock)
		mockTimeProvider.On("UtcNow").Return(time.Now())

		userSrv := NewUserService(mockRepo, mockTimeProvider)

		dto := &CreateUserDTO{
			ID:          domain.NewID(),
			Email:       "test@example.com",
			Password:    "admin123",
			FirstName:   "Test",
			LastName:    "Test",
			PhoneNumber: "1234567890",
			Addresses: []*CreateUserAddress{
				{
					Type:       1,
					Street:     "Test",
					City:       "New York",
					State:      "NY",
					PostalCode: "",
					Country:    "USA",
				},
			},
		}

		mockRepo.On("Create", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("some repository error"))

		err := userSrv.Create(context.Background(), dto)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed creating user")
	})
}

func TestUpdate(t *testing.T) {
	t.Run("update user", func(t *testing.T) {
		mockRepo := new(repoMock.UserRepositoryMock)
		mockTimeProvider := new(domainMock.TimeProviderMock)

		userSrv := NewUserService(mockRepo, mockTimeProvider)

		userID := domain.NewID().String()
		dto := &UpdateUserDTO{
			FirstName:   ptr("Test"),
			LastName:    ptr("Test"),
			PhoneNumber: ptr("1234567890"),
			Addresses: []*UpdateUserAddress{
				{
					Type:       1,
					Street:     ptr("Test"),
					City:       ptr("New York"),
					State:      ptr("NY"),
					PostalCode: ptr(""),
					Country:    ptr("USA"),
				},
			},
		}

		mockRepo.On("UpdateBasicFields", mock.Anything, mock.Anything, mock.Anything).Return(nil)
		mockRepo.On("UpdateAddress", mock.Anything, mock.Anything, 1, mock.Anything).Return(nil)

		err := userSrv.Update(context.Background(), userID, dto)
		assert.NoError(t, err)
	})

	t.Run("error parsing userID", func(t *testing.T) {
		mockRepo := new(repoMock.UserRepositoryMock)
		mockTimeProvider := new(domainMock.TimeProviderMock)

		userSrv := NewUserService(mockRepo, mockTimeProvider)

		invalidUserID := "invalid-uuid"
		dto := &UpdateUserDTO{}

		err := userSrv.Update(context.Background(), invalidUserID, dto)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed parsing uuid")
	})

	t.Run("error updating basic fields", func(t *testing.T) {
		mockRepo := new(repoMock.UserRepositoryMock)
		mockTimeProvider := new(domainMock.TimeProviderMock)

		userSrv := NewUserService(mockRepo, mockTimeProvider)

		userID := domain.NewID().String()
		dto := &UpdateUserDTO{
			FirstName: ptr("Test"),
		}

		mockRepo.On("UpdateBasicFields", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("some error"))

		err := userSrv.Update(context.Background(), userID, dto)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed updating user personal data")
	})

	t.Run("error updating address", func(t *testing.T) {
		mockRepo := new(repoMock.UserRepositoryMock)
		mockTimeProvider := new(domainMock.TimeProviderMock)

		userSrv := NewUserService(mockRepo, mockTimeProvider)

		userID := domain.NewID().String()
		dto := &UpdateUserDTO{
			Addresses: []*UpdateUserAddress{
				{
					Type:       1,
					Street:     ptr("Test"),
					City:       ptr("New York"),
					State:      ptr("NY"),
					PostalCode: ptr(""),
					Country:    ptr("USA"),
				},
			},
		}

		mockRepo.On("UpdateAddress", mock.Anything, mock.Anything, 1, mock.Anything).Return(errors.New("some error"))

		err := userSrv.Update(context.Background(), userID, dto)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed updating user address data")
	})

	t.Run("error inserting new address if address not found", func(t *testing.T) {
		mockRepo := new(repoMock.UserRepositoryMock)
		mockTimeProvider := new(domainMock.TimeProviderMock)

		mockTimeProvider.On("UtcNow").Return(time.Now())

		userSrv := NewUserService(mockRepo, mockTimeProvider)

		userID := domain.NewID().String()
		dto := &UpdateUserDTO{
			Addresses: []*UpdateUserAddress{
				{
					Type:       1,
					Street:     ptr("Test"),
					City:       ptr("New York"),
					State:      ptr("NY"),
					PostalCode: ptr(""),
					Country:    ptr("USA"),
				},
			},
		}

		mockRepo.On("UpdateAddress", mock.Anything, mock.Anything, 1, mock.Anything).Return(user.ErrAddressNotFound)
		mockRepo.On("InsertAddress", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("insert address error"))

		err := userSrv.Update(context.Background(), userID, dto)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed inserting additional address")
	})
}

func TestDelete(t *testing.T) {
	t.Run("delete user", func(t *testing.T) {
		mockRepo := new(repoMock.UserRepositoryMock)
		userSrv := NewUserService(mockRepo, new(domainMock.TimeProviderMock))

		userID := domain.NewID().String()

		mockRepo.On("Delete", mock.Anything, mock.Anything).Return(nil)

		err := userSrv.Delete(context.Background(), userID)
		assert.NoError(t, err)
	})

	t.Run("error parsing userID", func(t *testing.T) {
		mockRepo := new(repoMock.UserRepositoryMock)
		userSrv := NewUserService(mockRepo, new(domainMock.TimeProviderMock))

		invalidUserID := "sdasdasd31231"

		err := userSrv.Delete(context.Background(), invalidUserID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed parsing uuid")
	})

	t.Run("user not found", func(t *testing.T) {
		mockRepo := new(repoMock.UserRepositoryMock)
		userSrv := NewUserService(mockRepo, new(domainMock.TimeProviderMock))

		userID := domain.NewID().String()

		mockRepo.On("Delete", mock.Anything, mock.Anything).Return(user.ErrNotFound)

		err := userSrv.Delete(context.Background(), userID)
		assert.ErrorIs(t, err, user.ErrNotFound)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := new(repoMock.UserRepositoryMock)
		userSrv := NewUserService(mockRepo, new(domainMock.TimeProviderMock))

		userID := domain.NewID().String()

		mockRepo.On("Delete", mock.Anything, mock.Anything).Return(errors.New("some repository error"))

		err := userSrv.Delete(context.Background(), userID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed deleting user")
	})
}

func TestGet(t *testing.T) {
	t.Run("get user", func(t *testing.T) {
		mockRepo := new(repoMock.UserRepositoryMock)
		userSrv := NewUserService(mockRepo, nil)

		expectedUsers := []*user.User{
			{
				ID:          domain.NewID(),
				Email:       "test1@example.com",
				FirstName:   "Test1",
				LastName:    "Test1",
				PhoneNumber: "111111111",
			},
			{
				ID:          domain.NewID(),
				Email:       "test2@example.com",
				FirstName:   "Test2",
				LastName:    "Test2",
				PhoneNumber: "222222222",
			},
		}

		mockRepo.On("Get", mock.Anything, 1, 2).Return(expectedUsers, nil)

		users, err := userSrv.Get(context.Background(), 1, 2)
		assert.NoError(t, err)
		assert.Equal(t, expectedUsers, users)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := new(repoMock.UserRepositoryMock)
		userSrv := NewUserService(mockRepo, nil)

		mockRepo.On("Get", mock.Anything, 1, 2).Return(nil, errors.New("some repository error"))

		users, err := userSrv.Get(context.Background(), 1, 2)
		assert.Error(t, err)
		assert.Nil(t, users)
		assert.Contains(t, err.Error(), "some repository error")
	})
}

func TestGetByUUID(t *testing.T) {
	t.Run("get by uuid", func(t *testing.T) {
		mockRepo := new(repoMock.UserRepositoryMock)
		userSrv := NewUserService(mockRepo, nil)

		userID := domain.NewID()
		expectedUser := &user.User{
			ID:          userID,
			Email:       "test@example.com",
			FirstName:   "Test",
			LastName:    "Test",
			PhoneNumber: "1234567890",
		}

		mockRepo.On("GetByUUID", mock.Anything, userID).Return(expectedUser, nil)

		resultUser, err := userSrv.GetByUUID(context.Background(), userID.String())
		assert.NoError(t, err)
		assert.Equal(t, expectedUser, resultUser)
	})

	t.Run("error parsing userID", func(t *testing.T) {
		mockRepo := new(repoMock.UserRepositoryMock)
		userSrv := NewUserService(mockRepo, nil)

		invalidUserID := "invalid-uuid"

		resultUser, err := userSrv.GetByUUID(context.Background(), invalidUserID)
		assert.Error(t, err)
		assert.Nil(t, resultUser)
		assert.Contains(t, err.Error(), "failed parsing uuid")
	})

	t.Run("user not found", func(t *testing.T) {
		mockRepo := new(repoMock.UserRepositoryMock)
		userSrv := NewUserService(mockRepo, nil)

		userID := domain.NewID()

		mockRepo.On("GetByUUID", mock.Anything, userID).Return(nil, user.ErrNotFound)

		resultUser, err := userSrv.GetByUUID(context.Background(), userID.String())
		assert.ErrorIs(t, err, user.ErrNotFound)
		assert.Nil(t, resultUser)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := new(repoMock.UserRepositoryMock)
		userSrv := NewUserService(mockRepo, nil)

		userID := domain.NewID()

		mockRepo.On("GetByUUID", mock.Anything, userID).Return(nil, errors.New("some repository error"))

		resultUser, err := userSrv.GetByUUID(context.Background(), userID.String())
		assert.Error(t, err)
		assert.Nil(t, resultUser)
		assert.Contains(t, err.Error(), "some repository error")
	})
}

func ptr(s string) *string {
	return &s
}
