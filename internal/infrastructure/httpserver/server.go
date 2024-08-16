package httpserver

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
	"github.com/sarulabs/di"

	"github.com/wojciechpawlinow/usermanagement/internal/config"
	"github.com/wojciechpawlinow/usermanagement/internal/infrastructure/database/mysql"
	"github.com/wojciechpawlinow/usermanagement/internal/infrastructure/httpserver/handlers"
)

type Server struct {
	*http.Server
	shutdownDeps
}

type shutdownDeps struct {
	conns *mysql.Connections
}

// Run is a Server constructor that starts the HTTP server in a goroutine and enables routing
func Run(cfg config.Provider, ctn di.Container, errChan chan error) *Server {

	userHandler := ctn.Get("http-user").(*handlers.UserHTTPHandler)

	// define routes
	router := gin.Default()
	router.POST("/users", userHandler.CreateUser)
	router.PUT("/users/:id", userHandler.UpdateUser)
	router.DELETE("/users/:id", userHandler.DeleteUser)
	router.GET("/users/:id", userHandler.GetUser)
	router.GET("/users", userHandler.Get)

	s := &Server{
		&http.Server{
			Addr:              fmt.Sprintf(":%s", cfg.GetString("port")),
			Handler:           router,
			ReadHeaderTimeout: 5 * time.Second,
		},
		shutdownDeps{
			conns: ctn.Get("mysql-conns").(*mysql.Connections),
		},
	}

	go func() {
		if err := s.Server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errChan <- err
		}
	}()

	_, _ = color.New(color.FgHiGreen).Printf("\n=> an HTTP server listening at: %s\n\n", cfg.GetString("port"))

	return s
}

// Shutdown is a Shutdown function overload
func (srv *Server) Shutdown(ctx context.Context) error {
	srv.shutdownDeps.conns.Read.Close()
	srv.shutdownDeps.conns.Write.Close()

	return srv.Server.Shutdown(ctx)
}
