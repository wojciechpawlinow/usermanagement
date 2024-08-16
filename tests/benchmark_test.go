package tests

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/wojciechpawlinow/usermanagement/internal/config"
	"github.com/wojciechpawlinow/usermanagement/internal/infrastructure/container"
	"github.com/wojciechpawlinow/usermanagement/internal/infrastructure/httpserver/handlers"
	"github.com/wojciechpawlinow/usermanagement/pkg/logger"
)

func BenchmarkCreateUser(b *testing.B) {
	cfg := config.Load()
	cfg.Set("LOG_LEVEL", "info")
	logger.Setup(cfg)
	ctn := container.New()
	userHandler := ctn.Get("http-user").(*handlers.UserHTTPHandler)

	router := gin.Default()
	gin.SetMode(gin.ReleaseMode)
	router.POST("/users", userHandler.CreateUser)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		createReqBody := fmt.Sprintf(`{
				  "email": "%s",
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
				}`, generateRandomEmail())

		createReq, _ := http.NewRequest(http.MethodPost, "/users", io.NopCloser(strings.NewReader(createReqBody)))
		createReq.Header.Set("Content-Type", "application/json")
		createRec := httptest.NewRecorder()
		router.ServeHTTP(createRec, createReq)
	}
}

func generateRandomEmail() string {
	rand.Seed(time.Now().UnixNano())
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	usernameLength := rand.Intn(6) + 5
	username := make([]byte, usernameLength)
	for i := range username {
		username[i] = letters[rand.Intn(len(letters))]
	}
	domains := []string{"example.com", "test.com", "mail.com", "domain.com"}
	domain := domains[rand.Intn(len(domains))]

	return fmt.Sprintf("%s@%s", strings.ToLower(string(username)), domain)
}
