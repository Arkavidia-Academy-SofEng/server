package config

import (
	"ProjectGolang/database/postgres"
	authHandler "ProjectGolang/internal/api/auth/handler"
	authRepository "ProjectGolang/internal/api/auth/repository"
	authService "ProjectGolang/internal/api/auth/service"
	"ProjectGolang/internal/middleware"
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
	db         *sqlx.DB
	log        *logrus.Logger
	middleware middleware.Middleware
	validator  *validator.Validate
	handlers   []handler
}

type handler interface {
	Start(srv fiber.Router)
}

func NewServer(fiberApp *fiber.App, log *logrus.Logger, validator *validator.Validate) (*Server, error) {
	db, err := postgres.NewPostgresConnection()
	if err != nil {
		log.Errorf("Failed to connect to database: %v", err)
		return nil, err
	}

	bootstrap := &Server{
		engine:     fiberApp,
		db:         db,
		log:        log,
		validator:  validator,
		middleware: middleware.New(log),
	}

	return bootstrap, nil
}

func (s *Server) RegisterHandler() {
	//Auth Domain
	authRepo := authRepository.New(s.db, s.log)
	authServices := authService.New(authRepo, s.log)
	authHandlers := authHandler.New(authServices, s.validator, s.middleware, s.log)

	//Another Domain

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
