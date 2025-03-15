package authHandler

import (
	authService "ProjectGolang/internal/api/auth/service"
	"ProjectGolang/internal/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type AuthHandler struct {
	authService authService.AuthService
	validator   *validator.Validate
	middleware  middleware.Middleware
	log         *logrus.Logger
}

func New(as authService.AuthService, validate *validator.Validate, middleware middleware.Middleware, log *logrus.Logger) *AuthHandler {
	return &AuthHandler{
		authService: as,
		validator:   validate,
		middleware:  middleware,
		log:         log,
	}
}
func (h *AuthHandler) Start(srv fiber.Router) {
	auth := srv.Group("/auth")
	auth.Post("/register", h.HandleRegister)
	auth.Post("/login", h.HandleLogin)
	auth.Patch("/update", h.middleware.NewTokenMiddleware, h.HandleUpdateUser)
	auth.Delete("/delete/:id", h.HandleDeleteUser)
}
