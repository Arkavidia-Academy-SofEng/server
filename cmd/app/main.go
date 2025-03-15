package main

import (
	"ProjectGolang/internal/config"
	"ProjectGolang/pkg/log"
	"github.com/joho/godotenv"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logger := log.NewLogger()
	if err := godotenv.Load(); err != nil {
		logger.Fatal("Error loading .env file")
	}

	fiber := config.NewFiber(logger)
	validator := config.NewValidator()
	rest, err := config.NewServer(fiber, logger, validator)
	if err != nil {
		logger.Fatal(err)
	}

	rest.RegisterHandler()

	go func() {
		if err := rest.Run(); err != nil {
			logger.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")
	rest.Shutdown()
	logger.Info("Server shutdown complete")
}
