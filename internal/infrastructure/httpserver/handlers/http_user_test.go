package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/wojciechpawlinow/usermanagement/internal/application/service"
	"github.com/wojciechpawlinow/usermanagement/internal/config"
	"github.com/wojciechpawlinow/usermanagement/internal/domain"
	"github.com/wojciechpawlinow/usermanagement/internal/domain/user"
	"github.com/wojciechpawlinow/usermanagement/pkg/logger"
	serviceMock "github.com/wojciechpawlinow/usermanagement/tests/mocks/applicaion/service"
)

func ptr[T any](v T) *T {
	return &v
}

func TestGetUserByUUID(t *testing.T) {
	t.Run("get user by UUID", func(t *testing.T) {
		userID := domain.NewID()

		expectedUser := &user.User{
			ID:          userID,
			Email:       "test@example.com",
			FirstName:   "Richard",
			LastName:    "Gear",
			PhoneNumber: "664321234",
			Addresses: []*user.Address{
				{
					Type:       user.HomeAddress,
					Street:     "Main av",
					City:       "New York",
					State:      "NY",
					PostalCode: "",
					Country:    "USA",
				},
			},
		}

		s := new(serviceMock.UserServiceMock)
		s.On("GetByUUID", mock.Anything, userID.String()).Return(expectedUser, nil)

		userHandler := NewUserHTTPHandler(validator.New(), s)

		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.GET("/users/:id", userHandler.GetUser)

		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/users/%s", userID.String()), nil)
		assert.NoError(t, err)

		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)
		userJson, _ := json.Marshal(expectedUser)
		assert.Equal(t, string(userJson), recorder.Body.String())
	})

	t.Run("fail invalid request", func(t *testing.T) {
		s := new(serviceMock.UserServiceMock)
		userHandler := NewUserHTTPHandler(validator.New(), s)

		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.GET("/users/:id", userHandler.GetUser)

		req, err := http.NewRequest(http.MethodGet, "/users/adaasdasd231213", nil)
		assert.NoError(t, err)

		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("internal server error", func(t *testing.T) {
		userID := uuid.New().String()

		cfg := config.Load()
		logger.Setup(cfg)

		s := new(serviceMock.UserServiceMock)
		s.On("GetByUUID", mock.Anything, userID).Return(nil, errors.New("internal error"))

		userHandler := NewUserHTTPHandler(validator.New(), s)

		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.GET("/users/:id", userHandler.GetUser)

		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/users/%s", userID), nil)
		assert.NoError(t, err)

		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
		assert.Equal(t, `{"error":"internal server error"}`, recorder.Body.String())
	})

}

func TestGetUsers(t *testing.T) {
	t.Run("get users", func(t *testing.T) {
		expectedUsers := []*user.User{
			{
				ID:          domain.NewID(),
				Email:       "test1@example.com",
				FirstName:   "test1",
				LastName:    "test1",
				PhoneNumber: "111111111",
				Addresses: []*user.Address{
					{
						Type:       user.HomeAddress,
						Street:     "Test",
						City:       "New York",
						State:      "NY",
						PostalCode: "",
						Country:    "USA",
					},
				},
			},
			{
				ID:          domain.NewID(),
				Email:       "test2@example.com",
				FirstName:   "test2",
				LastName:    "test2",
				PhoneNumber: "111111111",
				Addresses: []*user.Address{
					{
						Type:       user.HomeAddress,
						Street:     "Test",
						City:       "New York",
						State:      "NY",
						PostalCode: "",
						Country:    "USA",
					},
				},
			},
			{
				ID:          domain.NewID(),
				Email:       "test3@example.com",
				FirstName:   "test3",
				LastName:    "test3",
				PhoneNumber: "111111111",
				Addresses: []*user.Address{
					{
						Type:       user.HomeAddress,
						Street:     "Test",
						City:       "New York",
						State:      "NY",
						PostalCode: "",
						Country:    "USA",
					},
				},
			},
		}

		s := new(serviceMock.UserServiceMock)
		s.On("Get", mock.Anything, 1, 3).Return(expectedUsers, nil)

		userHandler := NewUserHTTPHandler(validator.New(), s)

		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.GET("/users", userHandler.Get)

		req, err := http.NewRequest(http.MethodGet, "/users?size=3&page=1", nil)
		assert.NoError(t, err)

		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)
		usersJson, _ := json.Marshal(expectedUsers)
		assert.Equal(t, string(usersJson), recorder.Body.String())
	})

	t.Run("fail invalid request", func(t *testing.T) {
		s := new(serviceMock.UserServiceMock)
		userHandler := NewUserHTTPHandler(validator.New(), s)

		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.GET("/users", userHandler.Get)

		req, err := http.NewRequest(http.MethodGet, "/users?size=fail", nil)
		assert.NoError(t, err)

		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	// more use cases ...
}

func TestDeleteUser(t *testing.T) {
	t.Run("delete user", func(t *testing.T) {
		userID := uuid.New().String()

		s := new(serviceMock.UserServiceMock)
		s.On("Delete", mock.Anything, userID).Return(nil)

		userHandler := NewUserHTTPHandler(validator.New(), s)

		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.DELETE("/users/:id", userHandler.DeleteUser)

		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/users/%s", userID), nil)
		assert.NoError(t, err)

		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, `"ok"`, recorder.Body.String())
	})

	t.Run("invalid user ID", func(t *testing.T) {
		s := new(serviceMock.UserServiceMock)

		userHandler := NewUserHTTPHandler(validator.New(), s)

		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.DELETE("/users/:id", userHandler.DeleteUser)

		req, err := http.NewRequest(http.MethodDelete, "/users/asdasda23423", nil)
		assert.NoError(t, err)

		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		assert.Equal(t, `{"error":"invalid user ID"}`, recorder.Body.String())
	})

	t.Run("user not found", func(t *testing.T) {
		userID := uuid.New().String()

		s := new(serviceMock.UserServiceMock)
		s.On("Delete", mock.Anything, userID).Return(user.ErrNotFound)

		userHandler := NewUserHTTPHandler(validator.New(), s)

		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.DELETE("/users/:id", userHandler.DeleteUser)

		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/users/%s", userID), nil)
		assert.NoError(t, err)

		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusNotFound, recorder.Code)
		assert.Equal(t, `{"error":"user not found"}`, recorder.Body.String())
	})

	t.Run("internal server error", func(t *testing.T) {
		userID := uuid.New().String()

		s := new(serviceMock.UserServiceMock)
		s.On("Delete", mock.Anything, userID).Return(errors.New("internal error"))

		userHandler := NewUserHTTPHandler(validator.New(), s)

		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.DELETE("/users/:id", userHandler.DeleteUser)

		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/users/%s", userID), nil)
		assert.NoError(t, err)

		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
		assert.Equal(t, `{"error":"internal server error"}`, recorder.Body.String())
	})
}

func TestUpdateUser(t *testing.T) {
	t.Run("update user", func(t *testing.T) {
		userID := uuid.New().String()

		reqBody := `{
			"first_name": "Test",
			"last_name": "Test",
			"phone_number": "1234567890",
			"addresses": [
				{
					"type": 1,
					"street": "Test",
					"city": "New York",
					"state": "NY",
					"country": "USA"
				}
			]
		}`

		updateUserDTO := &service.UpdateUserDTO{
			FirstName:   ptr("Test"),
			LastName:    ptr("Test"),
			PhoneNumber: ptr("1234567890"),
			Addresses: []*service.UpdateUserAddress{
				{
					Type:    1,
					Street:  ptr("Test"),
					City:    ptr("New York"),
					State:   ptr("NY"),
					Country: ptr("USA"),
				},
			},
		}

		s := new(serviceMock.UserServiceMock)
		s.On("Update", mock.Anything, userID, updateUserDTO).Return(nil)

		userHandler := NewUserHTTPHandler(validator.New(), s)

		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.PUT("/users/:id", userHandler.UpdateUser)

		req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/users/%s", userID), nil)
		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")
		req.Body = io.NopCloser(strings.NewReader(reqBody))

		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, `"ok"`, recorder.Body.String())
	})

	t.Run("invalid user ID", func(t *testing.T) {
		reqBody := `{
			"first_name": "John",
			"last_name": "Doe"
		}`

		s := new(serviceMock.UserServiceMock)

		userHandler := NewUserHTTPHandler(validator.New(), s)

		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.PUT("/users/:id", userHandler.UpdateUser)

		req, err := http.NewRequest(http.MethodPut, "/users/zdcwcwe23234", nil)
		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")
		req.Body = io.NopCloser(strings.NewReader(reqBody))

		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		assert.Equal(t, `{"error":"invalid user ID"}`, recorder.Body.String())
	})

	t.Run("validation error", func(t *testing.T) {
		userID := uuid.New().String()

		reqBody := `{
			"first_name": "",
			"last_name": ""
		}`

		s := new(serviceMock.UserServiceMock)

		userHandler := NewUserHTTPHandler(validator.New(), s)

		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.PUT("/users/:id", userHandler.UpdateUser)

		req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/users/%s", userID), nil)
		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")
		req.Body = io.NopCloser(strings.NewReader(reqBody))

		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("user not found", func(t *testing.T) {
		userID := uuid.New().String()
		reqBody := `{
			"first_name": "Test",
			"last_name": "Test"
		}`

		s := new(serviceMock.UserServiceMock)

		s.On("Update", mock.Anything, userID, mock.Anything).Return(user.ErrNotFound)

		userHandler := NewUserHTTPHandler(validator.New(), s)

		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.PUT("/users/:id", userHandler.UpdateUser)

		req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/users/%s", userID), nil)
		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")
		req.Body = io.NopCloser(strings.NewReader(reqBody))

		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusNotFound, recorder.Code)
		assert.Equal(t, `{"error":"user not found"}`, recorder.Body.String())
	})
}

func TestCreateUser(t *testing.T) {
	t.Run("successful creation", func(t *testing.T) {
		reqBody := `{
			"email": "test@example.com",
			"password": "securePassword123",
			"first_name": "Test",
			"last_name": "Test",
			"phone_number": "1234567890",
			"addresses": [
				{
					"type": 1,
					"street": "Test",
					"city": "New York",
					"state": "NY",
					"postal_code": "55010",
					"country": "USA"
				}
			]
		}`

		s := new(serviceMock.UserServiceMock)
		s.On("Create", mock.Anything, mock.Anything).Return(nil)

		userHandler := NewUserHTTPHandler(validator.New(), s)

		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.POST("/users", userHandler.CreateUser)

		req, err := http.NewRequest(http.MethodPost, "/users", nil)
		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")
		req.Body = io.NopCloser(strings.NewReader(reqBody))

		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusCreated, recorder.Code)
		assert.Contains(t, recorder.Body.String(), `"uuid"`)
	})

	t.Run("validation error", func(t *testing.T) {
		reqBody := `{
			"email": "asdasd2323423",
			"password": "",
			"first_name": "",
			"last_name": "",
			"phone_number": "",
			"addresses": []
		}`

		s := new(serviceMock.UserServiceMock)
		userHandler := NewUserHTTPHandler(validator.New(), s)

		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.POST("/users", userHandler.CreateUser)

		req, err := http.NewRequest(http.MethodPost, "/users", nil)
		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")
		req.Body = io.NopCloser(strings.NewReader(reqBody))

		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("email already exists", func(t *testing.T) {
		reqBody := `{
			"email": "test@example.com",
			"password": "securePassword123",
			"first_name": "Test",
			"last_name": "Test",
			"phone_number": "1234567890",
			"addresses": [
				{
					"type": 1,
					"street": "Test",
					"city": "New York",
					"state": "NY",
					"postal_code": "55010",
					"country": "USA"
				}
			]
		}`

		s := new(serviceMock.UserServiceMock)
		s.On("Create", mock.Anything, mock.Anything).Return(user.ErrEmailAlreadyExists)

		userHandler := NewUserHTTPHandler(validator.New(), s)

		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.POST("/users", userHandler.CreateUser)

		req, err := http.NewRequest(http.MethodPost, "/users", nil)
		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")
		req.Body = io.NopCloser(strings.NewReader(reqBody))

		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusConflict, recorder.Code)
		assert.Equal(t, `{"error":"email already exists"}`, recorder.Body.String())
	})
}
