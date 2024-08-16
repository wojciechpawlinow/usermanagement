package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/wojciechpawlinow/usermanagement/internal/config"
	"github.com/wojciechpawlinow/usermanagement/internal/infrastructure/container"
	"github.com/wojciechpawlinow/usermanagement/internal/infrastructure/httpserver"
	"github.com/wojciechpawlinow/usermanagement/pkg/logger"
)

func main() {

	// load configuration from a file or fallback to defaults
	cfg := config.Load()

	// create logger
	logger.Setup(cfg)

	// build dependencies
	ctn := container.New()

	var err error
	errChan := make(chan error, 1)

	// create and run HTTP server
	srv := httpserver.Run(cfg, ctn, errChan)

	// wait for interrupt signal to gracefully shut down the server with a timeout of 10 seconds
	// use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt)
	signal.Notify(quit, syscall.SIGTERM)

	select {
	case <-quit:
		logger.Info("initiating graceful shutdown...")
	case err = <-errChan:
		logger.Error(fmt.Errorf("server error: %w", err))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// gracefully shut down the server
	if err = srv.Shutdown(ctx); err != nil {
		logger.Error(fmt.Errorf("server shutdown failed: %w", err))
	}
}
