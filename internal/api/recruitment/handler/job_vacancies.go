package recruitmentHandler

import (
	"ProjectGolang/internal/api/recruitment"
	"ProjectGolang/pkg/log"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"golang.org/x/net/context"
	"time"
)

func (h *RecruitmentHandler) CreateJobVacancy(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	h.log.WithField("path", ctx.Path()).Debug("Processing job vacancy creation request")

	var req recruitment.CreateJobVacancy
	if err := ctx.BodyParser(&req); err != nil {
		h.log.WithFields(log.Fields{
			"error": err.Error(),
			"path":  ctx.Path(),
		}).Error("Failed to parse job vacancy creation request body")
		return err
	}

	if err := h.validator.Struct(&req); err != nil {
		h.log.WithFields(log.Fields{
			"error": err.Error(),
			"title": req.Title,
		}).Warn("Validation failed for job vacancy creation")
		return err
	}

	if err := h.recruitmentService.JobVacancy().CreateJobVacancy(c, req); err != nil {
		h.log.WithFields(log.Fields{
			"error": err.Error(),
			"title": req.Title,
		}).Error("Job vacancy creation failed")
		return err
	}

	select {
	case <-c.Done():
		return ctx.Status(fiber.StatusRequestTimeout).
			JSON(utils.StatusMessage(fiber.StatusRequestTimeout))
	default:
		return ctx.SendStatus(fiber.StatusCreated)
	}
}

func (h *RecruitmentHandler) GetJobVacancies(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	h.log.WithField("path", ctx.Path()).Debug("Processing job vacancies fetch request")

	var req recruitment.GetJobVacancies
	if err := ctx.QueryParser(&req); err != nil {
		h.log.WithFields(log.Fields{
			"error": err.Error(),
			"path":  ctx.Path(),
		}).Error("Failed to parse job vacancies query parameters")
		return err
	}

	if err := h.validator.Struct(&req); err != nil {
		h.log.WithFields(log.Fields{
			"error": err.Error(),
			"page":  req.Page,
			"size":  req.PageSize,
		}).Warn("Validation failed for job vacancies request")
		return err
	}

	result, err := h.recruitmentService.JobVacancy().GetJobVacancies(c, req)
	if err != nil {
		h.log.WithFields(log.Fields{
			"error": err.Error(),
			"page":  req.Page,
			"size":  req.PageSize,
		}).Error("Job vacancies fetch failed")
		return err
	}

	select {
	case <-c.Done():
		return ctx.Status(fiber.StatusRequestTimeout).
			JSON(utils.StatusMessage(fiber.StatusRequestTimeout))
	default:
		return ctx.Status(fiber.StatusOK).JSON(result)
	}
}

func (h *RecruitmentHandler) UpdateJobVacancy(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	h.log.WithField("path", ctx.Path()).Debug("Processing job vacancy update request")

	// Get ID from URL parameters
	id := ctx.Params("id")
	if id == "" {
		h.log.WithFields(log.Fields{
			"path": ctx.Path(),
		}).Error("Missing job vacancy ID in URL")
		return fiber.NewError(fiber.StatusBadRequest, "Job vacancy ID is required")
	}

	var req recruitment.UpdateJobVacancy
	if err := ctx.BodyParser(&req); err != nil {
		h.log.WithFields(log.Fields{
			"error": err.Error(),
			"path":  ctx.Path(),
		}).Error("Failed to parse job vacancy update request body")
		return err
	}

	req.ID = id

	if err := h.validator.Struct(&req); err != nil {
		h.log.WithFields(log.Fields{
			"error": err.Error(),
			"id":    req.ID,
			"title": req.Title,
		}).Warn("Validation failed for job vacancy update")
		return err
	}

	if err := h.recruitmentService.JobVacancy().UpdateJobVacancy(c, req); err != nil {
		h.log.WithFields(log.Fields{
			"error": err.Error(),
			"id":    req.ID,
			"title": req.Title,
		}).Error("Job vacancy update failed")
		return err
	}

	select {
	case <-c.Done():
		return ctx.Status(fiber.StatusRequestTimeout).
			JSON(utils.StatusMessage(fiber.StatusRequestTimeout))
	default:
		return ctx.SendStatus(fiber.StatusOK)
	}
}

func (h *RecruitmentHandler) DeleteJobVacancy(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	h.log.WithField("path", ctx.Path()).Debug("Processing job vacancy deletion request")

	// Get ID from URL parameters
	id := ctx.Params("id")
	if id == "" {
		h.log.WithFields(log.Fields{
			"path": ctx.Path(),
		}).Error("Missing job vacancy ID in URL")
		return fiber.NewError(fiber.StatusBadRequest, "Job vacancy ID is required")
	}

	if err := h.recruitmentService.JobVacancy().DeleteJobVacancy(c, id); err != nil {
		h.log.WithFields(log.Fields{
			"error": err.Error(),
			"id":    id,
		}).Error("Job vacancy deletion failed")
		return err
	}

	select {
	case <-c.Done():
		return ctx.Status(fiber.StatusRequestTimeout).
			JSON(utils.StatusMessage(fiber.StatusRequestTimeout))
	default:
		return ctx.SendStatus(fiber.StatusNoContent)
	}
}
