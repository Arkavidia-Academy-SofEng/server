package bioHandler

import (
	"ProjectGolang/internal/api/bio"
	contextPkg "ProjectGolang/pkg/context"
	"ProjectGolang/pkg/log"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"golang.org/x/net/context"
	"mime/multipart"
	"net/http"
	"time"
)

func (h *BioHandler) CreateExperience(ctx *fiber.Ctx) error {
	requestID := h.middleware.GetRequestID(ctx)
	c, cancel := context.WithTimeout(contextPkg.FromFiberCtx(ctx), 5*time.Second)
	defer cancel()

	h.log.WithFields(log.Fields{
		"request_id": requestID,
		"path":       ctx.Path(),
	}).Debug("Processing experience creation request")

	userID := ctx.Params("userId")
	if userID == "" {
		h.log.WithFields(log.Fields{
			"request_id": requestID,
			"path":       ctx.Path(),
		}).Warn("Missing user ID in URL")
		return fiber.NewError(fiber.StatusBadRequest, "User ID is required")
	}

	req := bio.CreateExperience{
		JobTitle:    ctx.FormValue("job_title"),
		JobLocation: ctx.FormValue("job_location"),
		SkillUsed:   ctx.FormValue("skill_used"),
		StartDate:   ctx.FormValue("start_date"),
		EndDate:     ctx.FormValue("end_date"),
		Description: ctx.FormValue("description"),
	}

	imageFile, err := ctx.FormFile("image")
	if err != nil && !errors.Is(err, http.ErrMissingFile) {
		h.log.WithFields(log.Fields{
			"request_id": requestID,
			"error":      err.Error(),
		}).Warn("Failed to parse profile picture")
		return err
	}

	if err := h.validator.Struct(&req); err != nil {
		h.log.WithFields(log.Fields{
			"request_id": requestID,
			"error":      err.Error(),
			"user_id":    userID,
		}).Warn("Validation failed for experience creation")
		return err
	}

	if err := h.bioService.CreateExperience(c, req, userID, imageFile); err != nil {
		status := fiber.StatusInternalServerError
		return ctx.Status(status).JSON(fiber.Map{
			"errors": fiber.Map{"message": "Experience creation failed", "err": err.Error()},
		})
	}

	select {
	case <-c.Done():
		return ctx.Status(fiber.StatusRequestTimeout).
			JSON(utils.StatusMessage(fiber.StatusRequestTimeout))
	default:
		return ctx.SendStatus(fiber.StatusCreated)
	}
}

func (h *BioHandler) GetExperienceByID(ctx *fiber.Ctx) error {
	requestID := h.middleware.GetRequestID(ctx)
	c, cancel := context.WithTimeout(contextPkg.FromFiberCtx(ctx), 5*time.Second)
	defer cancel()

	h.log.WithFields(log.Fields{
		"request_id": requestID,
		"path":       ctx.Path(),
	}).Debug("Processing get experience request")

	id := ctx.Params("id")
	if id == "" {
		h.log.WithFields(log.Fields{
			"request_id": requestID,
			"path":       ctx.Path(),
		}).Warn("Missing experience ID in URL")
		return fiber.NewError(fiber.StatusBadRequest, "Experience ID is required")
	}

	experience, err := h.bioService.GetExperienceByID(c, id)
	if err != nil {
		status := fiber.StatusInternalServerError
		return ctx.Status(status).JSON(fiber.Map{
			"errors": fiber.Map{"message": "Failed to get experience", "err": err.Error()},
		})
	}

	if experience.ID == "" {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"errors": fiber.Map{"message": "Experience not found"},
		})
	}

	select {
	case <-c.Done():
		return ctx.Status(fiber.StatusRequestTimeout).
			JSON(utils.StatusMessage(fiber.StatusRequestTimeout))
	default:
		return ctx.Status(fiber.StatusOK).JSON(experience)
	}
}

func (h *BioHandler) GetExperiencesByUserID(ctx *fiber.Ctx) error {
	requestID := h.middleware.GetRequestID(ctx)
	c, cancel := context.WithTimeout(contextPkg.FromFiberCtx(ctx), 5*time.Second)
	defer cancel()

	h.log.WithFields(log.Fields{
		"request_id": requestID,
		"path":       ctx.Path(),
	}).Debug("Processing get experiences by user ID request")

	userID := ctx.Params("userId")
	if userID == "" {
		h.log.WithFields(log.Fields{
			"request_id": requestID,
			"path":       ctx.Path(),
		}).Warn("Missing user ID in URL")
		return fiber.NewError(fiber.StatusBadRequest, "User ID is required")
	}

	experiences, err := h.bioService.GetExperiencesByUserID(c, userID)
	if err != nil {
		status := fiber.StatusInternalServerError
		return ctx.Status(status).JSON(fiber.Map{
			"errors": fiber.Map{"message": "Failed to get experiences", "err": err.Error()},
		})
	}

	select {
	case <-c.Done():
		return ctx.Status(fiber.StatusRequestTimeout).
			JSON(utils.StatusMessage(fiber.StatusRequestTimeout))
	default:
		return ctx.Status(fiber.StatusOK).JSON(experiences)
	}
}

func (h *BioHandler) UpdateExperience(ctx *fiber.Ctx) error {
	requestID := h.middleware.GetRequestID(ctx)
	c, cancel := context.WithTimeout(contextPkg.FromFiberCtx(ctx), 5*time.Second)
	defer cancel()

	h.log.WithFields(log.Fields{
		"request_id": requestID,
		"path":       ctx.Path(),
	}).Debug("Processing experience update request")

	id := ctx.Params("id")
	if id == "" {
		h.log.WithFields(log.Fields{
			"request_id": requestID,
			"path":       ctx.Path(),
		}).Warn("Missing experience ID in URL")
		return fiber.NewError(fiber.StatusBadRequest, "Experience ID is required")
	}

	req := bio.UpdateExperience{
		JobTitle:    ctx.FormValue("job_title"),
		SkillUsed:   ctx.FormValue("skill_used"),
		StartDate:   ctx.FormValue("start_date"),
		EndDate:     ctx.FormValue("end_date"),
		Description: ctx.FormValue("description"),
	}

	if err := h.validator.Struct(&req); err != nil {
		h.log.WithFields(log.Fields{
			"request_id": requestID,
			"error":      err.Error(),
			"id":         id,
		}).Warn("Validation failed for experience update")
		return err
	}

	var imageFile *multipart.FileHeader
	imageFile, err := ctx.FormFile("image")
	if err != nil && !errors.Is(err, http.ErrMissingFile) {
		h.log.WithFields(log.Fields{
			"request_id": requestID,
			"error":      err.Error(),
		}).Warn("Failed to parse profile picture")
		return err
	}

	if err := h.bioService.UpdateExperience(c, req, id, imageFile); err != nil {
		if errors.Is(err, errors.New("experience not found")) {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"errors": fiber.Map{"message": "Experience not found"},
			})
		}
		status := fiber.StatusInternalServerError
		return ctx.Status(status).JSON(fiber.Map{
			"errors": fiber.Map{"message": "Experience update failed", "err": err.Error()},
		})
	}

	select {
	case <-c.Done():
		return ctx.Status(fiber.StatusRequestTimeout).
			JSON(utils.StatusMessage(fiber.StatusRequestTimeout))
	default:
		return ctx.SendStatus(fiber.StatusOK)
	}
}

func (h *BioHandler) DeleteExperience(ctx *fiber.Ctx) error {
	requestID := h.middleware.GetRequestID(ctx)
	c, cancel := context.WithTimeout(contextPkg.FromFiberCtx(ctx), 5*time.Second)
	defer cancel()

	h.log.WithFields(log.Fields{
		"request_id": requestID,
		"path":       ctx.Path(),
	}).Debug("Processing experience deletion request")

	id := ctx.Params("id")
	if id == "" {
		h.log.WithFields(log.Fields{
			"request_id": requestID,
			"path":       ctx.Path(),
		}).Warn("Missing experience ID in URL")
		return fiber.NewError(fiber.StatusBadRequest, "Experience ID is required")
	}

	if err := h.bioService.DeleteExperience(c, id); err != nil {
		if errors.Is(err, errors.New("experience not found")) {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"errors": fiber.Map{"message": "Experience not found"},
			})
		}
		status := fiber.StatusInternalServerError
		return ctx.Status(status).JSON(fiber.Map{
			"errors": fiber.Map{"message": "Experience deletion failed", "err": err.Error()},
		})
	}

	select {
	case <-c.Done():
		return ctx.Status(fiber.StatusRequestTimeout).
			JSON(utils.StatusMessage(fiber.StatusRequestTimeout))
	default:
		return ctx.SendStatus(fiber.StatusNoContent)
	}
}
