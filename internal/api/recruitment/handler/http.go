package recruitmentHandler

import (
	recruitmentService "ProjectGolang/internal/api/recruitment/service"
	"ProjectGolang/internal/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type RecruitmentHandler struct {
	recruitmentService recruitmentService.RecruitmentService
	validator          *validator.Validate
	middleware         middleware.Middleware
	log                *logrus.Logger
}

func New(rs recruitmentService.RecruitmentService, validate *validator.Validate, middleware middleware.Middleware, log *logrus.Logger) *RecruitmentHandler {
	return &RecruitmentHandler{
		recruitmentService: rs,
		validator:          validate,
		middleware:         middleware,
		log:                log,
	}
}

func (h *RecruitmentHandler) Start(srv fiber.Router) {
	rc := srv.Group("/recruitment")
	jv := rc.Group("/job_vacancies")
	jv.Post("/", h.CreateJobVacancy)
	jv.Get("/", h.GetJobVacancies)
	jv.Put("/:id", h.UpdateJobVacancy)
	jv.Delete("/:id", h.DeleteJobVacancy)
}
