package bioHandler

import (
	bioService "ProjectGolang/internal/api/bio/service"
	"ProjectGolang/internal/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type BioHandler struct {
	bioService bioService.BioService
	validator  *validator.Validate
	middleware middleware.Middleware
	log        *logrus.Logger
}

func New(bs bioService.BioService, validate *validator.Validate, middleware middleware.Middleware, log *logrus.Logger) *BioHandler {
	return &BioHandler{
		bioService: bs,
		validator:  validate,
		middleware: middleware,
		log:        log,
	}
}
func (h *BioHandler) Start(srv fiber.Router) {
	experiences := srv.Group("/experiences")
	experiences.Get("/:id", h.GetExperienceByID)
	experiences.Put("/:id", h.middleware.NewTokenMiddleware, h.UpdateExperience)
	experiences.Delete("/:id", h.middleware.NewTokenMiddleware, h.DeleteExperience)

	userExperiences := srv.Group("/users/:userId/experiences")
	userExperiences.Post("/", h.middleware.NewTokenMiddleware, h.CreateExperience)
	userExperiences.Get("/", h.GetExperiencesByUserID)

	educations := srv.Group("/educations")
	educations.Get("/:id", h.GetEducationByID)
	educations.Put("/:id", h.middleware.NewTokenMiddleware, h.UpdateEducation)
	educations.Delete("/:id", h.middleware.NewTokenMiddleware, h.DeleteEducation)

	userEducations := srv.Group("/users/:userId/educations")
	userEducations.Post("/", h.middleware.NewTokenMiddleware, h.CreateEducation)
	userEducations.Get("/", h.GetEducationsByUserID)

	portfolios := srv.Group("/portfolios")
	portfolios.Get("/:id", h.GetPortfolioByID)
	portfolios.Put("/:id", h.middleware.NewTokenMiddleware, h.UpdatePortfolio)
	portfolios.Delete("/:id", h.middleware.NewTokenMiddleware, h.DeletePortfolio)

	userPortfolios := srv.Group("/users/:userId/portfolios")
	userPortfolios.Post("/", h.middleware.NewTokenMiddleware, h.CreatePortfolio)
	userPortfolios.Get("/", h.GetPortfoliosByUserID)
}
