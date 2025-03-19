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

func (h *BioHandler) CreateEducation(ctx *fiber.Ctx) error {
	requestID := h.middleware.GetRequestID(ctx)
	c, cancel := context.WithTimeout(contextPkg.FromFiberCtx(ctx), 5*time.Second)
	defer cancel()

	h.log.WithFields(log.Fields{
		"request_id": requestID,
		"path":       ctx.Path(),
	}).Debug("Processing education creation request")

	userID := ctx.Params("userId")
	if userID == "" {
		h.log.WithFields(log.Fields{
			"request_id": requestID,
			"path":       ctx.Path(),
		}).Warn("Missing user ID in URL")
		return fiber.NewError(fiber.StatusBadRequest, "User ID is required")
	}

	req := bio.CreateEducation{
		TitleDegree:       ctx.FormValue("title_degree"),
		InstitutionalName: ctx.FormValue("institutional_name"),
		StartDate:         ctx.FormValue("start_date"),
		EndDate:           ctx.FormValue("end_date"),
		Description:       ctx.FormValue("description"),
	}

	imageFile, err := ctx.FormFile("image")
	if err != nil && !errors.Is(err, http.ErrMissingFile) {
		h.log.WithFields(log.Fields{
			"request_id": requestID,
			"error":      err.Error(),
		}).Warn("Failed to parse education image")
		return err
	}

	if err := h.validator.Struct(&req); err != nil {
		h.log.WithFields(log.Fields{
			"request_id": requestID,
			"error":      err.Error(),
			"user_id":    userID,
		}).Warn("Validation failed for education creation")
		return err
	}

	if err := h.bioService.CreateEducation(c, req, userID, imageFile); err != nil {
		status := fiber.StatusInternalServerError
		return ctx.Status(status).JSON(fiber.Map{
			"errors": fiber.Map{"message": "Education creation failed", "err": err.Error()},
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

func (h *BioHandler) GetEducationByID(ctx *fiber.Ctx) error {
	requestID := h.middleware.GetRequestID(ctx)
	c, cancel := context.WithTimeout(contextPkg.FromFiberCtx(ctx), 5*time.Second)
	defer cancel()

	h.log.WithFields(log.Fields{
		"request_id": requestID,
		"path":       ctx.Path(),
	}).Debug("Processing get education request")

	id := ctx.Params("id")
	if id == "" {
		h.log.WithFields(log.Fields{
			"request_id": requestID,
			"path":       ctx.Path(),
		}).Warn("Missing education ID in URL")
		return fiber.NewError(fiber.StatusBadRequest, "Education ID is required")
	}

	education, err := h.bioService.GetEducationByID(c, id)
	if err != nil {
		status := fiber.StatusInternalServerError
		return ctx.Status(status).JSON(fiber.Map{
			"errors": fiber.Map{"message": "Failed to get education", "err": err.Error()},
		})
	}

	if education.ID == "" {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"errors": fiber.Map{"message": "Education not found"},
		})
	}

	select {
	case <-c.Done():
		return ctx.Status(fiber.StatusRequestTimeout).
			JSON(utils.StatusMessage(fiber.StatusRequestTimeout))
	default:
		return ctx.Status(fiber.StatusOK).JSON(education)
	}
}

func (h *BioHandler) GetEducationsByUserID(ctx *fiber.Ctx) error {
	requestID := h.middleware.GetRequestID(ctx)
	c, cancel := context.WithTimeout(contextPkg.FromFiberCtx(ctx), 5*time.Second)
	defer cancel()

	h.log.WithFields(log.Fields{
		"request_id": requestID,
		"path":       ctx.Path(),
	}).Debug("Processing get educations by user ID request")

	userID := ctx.Params("userId")
	if userID == "" {
		h.log.WithFields(log.Fields{
			"request_id": requestID,
			"path":       ctx.Path(),
		}).Warn("Missing user ID in URL")
		return fiber.NewError(fiber.StatusBadRequest, "User ID is required")
	}

	educations, err := h.bioService.GetEducationsByUserID(c, userID)
	if err != nil {
		status := fiber.StatusInternalServerError
		return ctx.Status(status).JSON(fiber.Map{
			"errors": fiber.Map{"message": "Failed to get educations", "err": err.Error()},
		})
	}

	select {
	case <-c.Done():
		return ctx.Status(fiber.StatusRequestTimeout).
			JSON(utils.StatusMessage(fiber.StatusRequestTimeout))
	default:
		return ctx.Status(fiber.StatusOK).JSON(educations)
	}
}

func (h *BioHandler) UpdateEducation(ctx *fiber.Ctx) error {
	requestID := h.middleware.GetRequestID(ctx)
	c, cancel := context.WithTimeout(contextPkg.FromFiberCtx(ctx), 5*time.Second)
	defer cancel()

	h.log.WithFields(log.Fields{
		"request_id": requestID,
		"path":       ctx.Path(),
	}).Debug("Processing education update request")

	id := ctx.Params("id")
	if id == "" {
		h.log.WithFields(log.Fields{
			"request_id": requestID,
			"path":       ctx.Path(),
		}).Warn("Missing education ID in URL")
		return fiber.NewError(fiber.StatusBadRequest, "Education ID is required")
	}

	req := bio.UpdateEducation{
		TitleDegree:       ctx.FormValue("title_degree"),
		InstitutionalName: ctx.FormValue("institutional_name"),
		StartDate:         ctx.FormValue("start_date"),
		EndDate:           ctx.FormValue("end_date"),
		Description:       ctx.FormValue("description"),
	}

	if err := h.validator.Struct(&req); err != nil {
		h.log.WithFields(log.Fields{
			"request_id": requestID,
			"error":      err.Error(),
			"id":         id,
		}).Warn("Validation failed for education update")
		return err
	}

	var imageFile *multipart.FileHeader
	imageFile, err := ctx.FormFile("image")
	if err != nil && !errors.Is(err, http.ErrMissingFile) {
		h.log.WithFields(log.Fields{
			"request_id": requestID,
			"error":      err.Error(),
		}).Warn("Failed to parse education image")
		return err
	}

	if err := h.bioService.UpdateEducation(c, req, id, imageFile); err != nil {
		if errors.Is(err, errors.New("education not found")) {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"errors": fiber.Map{"message": "Education not found"},
			})
		}
		status := fiber.StatusInternalServerError
		return ctx.Status(status).JSON(fiber.Map{
			"errors": fiber.Map{"message": "Education update failed", "err": err.Error()},
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

func (h *BioHandler) DeleteEducation(ctx *fiber.Ctx) error {
	requestID := h.middleware.GetRequestID(ctx)
	c, cancel := context.WithTimeout(contextPkg.FromFiberCtx(ctx), 5*time.Second)
	defer cancel()

	h.log.WithFields(log.Fields{
		"request_id": requestID,
		"path":       ctx.Path(),
	}).Debug("Processing education deletion request")

	id := ctx.Params("id")
	if id == "" {
		h.log.WithFields(log.Fields{
			"request_id": requestID,
			"path":       ctx.Path(),
		}).Warn("Missing education ID in URL")
		return fiber.NewError(fiber.StatusBadRequest, "Education ID is required")
	}

	if err := h.bioService.DeleteEducation(c, id); err != nil {
		if errors.Is(err, errors.New("education not found")) {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"errors": fiber.Map{"message": "Education not found"},
			})
		}
		status := fiber.StatusInternalServerError
		return ctx.Status(status).JSON(fiber.Map{
			"errors": fiber.Map{"message": "Education deletion failed", "err": err.Error()},
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
