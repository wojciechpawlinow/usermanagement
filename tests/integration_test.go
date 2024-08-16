package tests

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/wojciechpawlinow/usermanagement/internal/config"
	"github.com/wojciechpawlinow/usermanagement/internal/infrastructure/container"
	"github.com/wojciechpawlinow/usermanagement/internal/infrastructure/httpserver/handlers"
	"github.com/wojciechpawlinow/usermanagement/pkg/logger"
)

func TestIntegration(t *testing.T) {
	cfg := config.Load()
	logger.Setup(cfg)
	ctn := container.New()
	userHandler := ctn.Get("http-user").(*handlers.UserHTTPHandler)
	router := gin.Default()
	router.POST("/users", userHandler.CreateUser)
	router.PUT("/users/:id", userHandler.UpdateUser)
	router.DELETE("/users/:id", userHandler.DeleteUser)
	router.GET("/users/:id", userHandler.GetUser)
	router.GET("/users", userHandler.Get)

	createReqBody := `{
	  "email": "test999@myemailxx.com",
	  "password": "secure123",
	  "first_name": "FirstName",
	  "last_name": "LastName",
	  "phone_number": "1234567890",
	  "addresses": [
		{
		  "type": 1,
		  "street": "Test avenue",
		  "city": "New York",
		  "state": "NY",
		  "postal_code": "55010",
		  "country": "USA"
		}
	  ]
	}`

	createReq, _ := http.NewRequest(http.MethodPost, "/users", io.NopCloser(strings.NewReader(createReqBody)))
	createReq.Header.Set("Content-Type", "application/json")
	createRec := httptest.NewRecorder()
	router.ServeHTTP(createRec, createReq)

	assert.Equal(t, http.StatusCreated, createRec.Code)

	var resp map[string]string
	_ = json.Unmarshal([]byte(createRec.Body.String()), &resp)

	userID := resp["uuid"]

	getReq, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/users/%s", userID), nil)
	getRec := httptest.NewRecorder()
	router.ServeHTTP(getRec, getReq)

	assert.Equal(t, http.StatusOK, getRec.Code)
	assert.Contains(t, getRec.Body.String(), `"email":"test999@myemailxx.com"`)

	updateReqBody := `{
		"first_name": "New test name",
		"addresses": [
		  {
	        "type": 3,
		    "street": "New address type",
		    "city": "Warszawa",
		    "postal_code": "62702"
		  }
	    ]
	}`
	updateReq, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/users/%s", userID), io.NopCloser(strings.NewReader(updateReqBody)))
	updateReq.Header.Set("Content-Type", "application/json")

	updateRec := httptest.NewRecorder()
	router.ServeHTTP(updateRec, updateReq)

	assert.Equal(t, http.StatusOK, updateRec.Code)

	deleteReq, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/users/%s", userID), nil)
	deleteRec := httptest.NewRecorder()
	router.ServeHTTP(deleteRec, deleteReq)

	assert.Equal(t, http.StatusOK, deleteRec.Code)
}
