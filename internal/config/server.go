package config

import (
	"ProjectGolang/database/postgres"
	authHandler "ProjectGolang/internal/api/auth/handler"
	authRepository "ProjectGolang/internal/api/auth/repository"
	authService "ProjectGolang/internal/api/auth/service"
	"ProjectGolang/internal/middleware"
	"ProjectGolang/pkg/s3"
	"ProjectGolang/pkg/scheduler"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"os"
)

type Server struct {
	engine     *fiber.App
	DB         *sqlx.DB
	log        *logrus.Logger
	middleware middleware.Middleware
	validator  *validator.Validate
	s3         s3.ItfS3
	scheduler  *scheduler.Scheduler
	handlers   []handler
}

type handler interface {
	Start(srv fiber.Router)
}

func NewServer(fiberApp *fiber.App, log *logrus.Logger, validator *validator.Validate) (*Server, error) {
	DB, err := postgres.NewPostgresConnection()
	if err != nil {
		log.Errorf("Failed to connect to database: %v", err)
		return nil, err
	}

	objectDB, err := s3.New()
	if err != nil {
		log.Errorf("Failed to connect to S3: %v", err)
		return nil, err
	}

	bootstrap := &Server{
		engine:     fiberApp,
		DB:         DB,
		log:        log,
		validator:  validator,
		middleware: middleware.New(log),
		s3:         objectDB,
	}

	return bootstrap, nil
}

func (s *Server) RegisterHandler() {
	//Auth Domain
	authRepo := authRepository.New(s.DB, s.log)
	authServices := authService.New(authRepo, s.log)
	authHandlers := authHandler.New(authServices, s.validator, s.middleware, s.log)
	timeScheduler := scheduler.NewScheduler(authRepo, s.log)

	//Another Domain

	timeScheduler.Start()
	s.scheduler = timeScheduler
	s.log.Info("Scheduler started successfully")
	s.checkHealth()
	s.handlers = append(s.handlers, authHandlers)
}

func (s *Server) Run() error {
	s.engine.Use(cors.New())
	s.engine.Use(s.middleware.NewLoggingMiddleware)
	router := s.engine.Group("/api/v1")

	for _, h := range s.handlers {
		h.Start(router)
	}

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000"
	}

	s.log.Infof("Starting server on port %s", port)

	if err := s.engine.Listen(fmt.Sprintf(":%s", port)); err != nil {
		return err
	}
	return nil
}

func (s *Server) checkHealth() {
	s.engine.Get("/", func(ctx *fiber.Ctx) error {
		s.log.Info("Health check endpoint called")
		return ctx.JSON(fiber.Map{
			"message": "Server is Healthy!",
		})
	})
}

func (s *Server) Shutdown() {
	// Stop the scheduler
	if s.scheduler != nil {
		s.scheduler.Stop()
		s.log.Info("Scheduler stopped")
	}

	// Close database connections
	if s.DB != nil {
		err := s.DB.Close()
		if err != nil {
			return
		}

		s.log.Info("Database connection closed")
	}

	// Shutdown fiber app
	if s.engine != nil {
		err := s.engine.Shutdown()
		if err != nil {
			return
		}
		s.log.Info("Fiber app shutdown")
	}
}
