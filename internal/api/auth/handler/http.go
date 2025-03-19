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
	users := srv.Group("/users")
	users.Post("/otp", h.RequestOTP)
	users.Post("/", h.CreateUser)
	users.Post("/login", h.Login)
	users.Put("/:id", h.middleware.NewTokenMiddleware, h.UpdateUser)

	companies := srv.Group("/companies")
	companies.Put("/:id", h.middleware.NewTokenMiddleware, h.UpdateCompany)

}
