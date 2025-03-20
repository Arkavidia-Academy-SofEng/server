package bioHandler

import (
	"ProjectGolang/internal/api/bio"
	contextPkg "ProjectGolang/pkg/context"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"mime/multipart"
	"net/http"
	"time"
)

func (h *BioHandler) CreatePortfolio(ctx *fiber.Ctx) error {
	requestID := h.middleware.GetRequestID(ctx)
	c, cancel := context.WithTimeout(contextPkg.FromFiberCtx(ctx), 5*time.Second)
	defer cancel()

	h.log.WithFields(logrus.Fields{
		"request_id": requestID,
		"path":       ctx.Path(),
	}).Debug("Processing portfolio creation request")

	userID := ctx.Params("userId")
	if userID == "" {
		h.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"path":       ctx.Path(),
		}).Warn("Missing user ID in URL")
		return fiber.NewError(fiber.StatusBadRequest, "User ID is required")
	}

	req := bio.CreatePortfolio{
		ProjectName:     ctx.FormValue("project_name"),
		ProjectLocation: ctx.FormValue("project_location"),
		ProjectLink:     ctx.FormValue("project_link"),
		StartDate:       ctx.FormValue("start_date"),
		EndDate:         ctx.FormValue("end_date"),
		Description:     ctx.FormValue("description"),
	}

	imageFile, err := ctx.FormFile("image")
	if err != nil && !errors.Is(err, http.ErrMissingFile) {
		h.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"error":      err.Error(),
		}).Warn("Failed to parse portfolio image")
		return err
	}

	descriptionImage, err := ctx.FormFile("description_image")
	if err != nil && !errors.Is(err, http.ErrMissingFile) {
		h.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"error":      err.Error(),
		}).Warn("Failed to parse description image")
		return err
	}

	if err := h.validator.Struct(&req); err != nil {
		h.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"error":      err.Error(),
			"user_id":    userID,
		}).Warn("Validation failed for portfolio creation")
		return err
	}

	if err := h.bioService.CreatePortfolio(c, req, userID, imageFile, descriptionImage); err != nil {
		status := fiber.StatusInternalServerError
		return ctx.Status(status).JSON(fiber.Map{
			"errors": fiber.Map{"message": "Portfolio creation failed", "err": err.Error()},
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

func (h *BioHandler) GetPortfolioByID(ctx *fiber.Ctx) error {
	requestID := h.middleware.GetRequestID(ctx)
	c, cancel := context.WithTimeout(contextPkg.FromFiberCtx(ctx), 5*time.Second)
	defer cancel()

	h.log.WithFields(logrus.Fields{
		"request_id": requestID,
		"path":       ctx.Path(),
	}).Debug("Processing get portfolio request")

	id := ctx.Params("id")
	if id == "" {
		h.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"path":       ctx.Path(),
		}).Warn("Missing portfolio ID in URL")
		return fiber.NewError(fiber.StatusBadRequest, "Portfolio ID is required")
	}

	portfolio, err := h.bioService.GetPortfolioByID(c, id)
	if err != nil {
		status := fiber.StatusInternalServerError
		return ctx.Status(status).JSON(fiber.Map{
			"errors": fiber.Map{"message": "Failed to get portfolio", "err": err.Error()},
		})
	}

	if portfolio.ID == "" {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"errors": fiber.Map{"message": "Portfolio not found"},
		})
	}

	select {
	case <-c.Done():
		return ctx.Status(fiber.StatusRequestTimeout).
			JSON(utils.StatusMessage(fiber.StatusRequestTimeout))
	default:
		return ctx.Status(fiber.StatusOK).JSON(portfolio)
	}
}

func (h *BioHandler) GetPortfoliosByUserID(ctx *fiber.Ctx) error {
	requestID := h.middleware.GetRequestID(ctx)
	c, cancel := context.WithTimeout(contextPkg.FromFiberCtx(ctx), 5*time.Second)
	defer cancel()

	h.log.WithFields(logrus.Fields{
		"request_id": requestID,
		"path":       ctx.Path(),
	}).Debug("Processing get portfolios by user ID request")

	userID := ctx.Params("userId")
	if userID == "" {
		h.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"path":       ctx.Path(),
		}).Warn("Missing user ID in URL")
		return fiber.NewError(fiber.StatusBadRequest, "User ID is required")
	}

	portfolios, err := h.bioService.GetPortfoliosByUserID(c, userID)
	if err != nil {
		status := fiber.StatusInternalServerError
		return ctx.Status(status).JSON(fiber.Map{
			"errors": fiber.Map{"message": "Failed to get portfolios", "err": err.Error()},
		})
	}

	select {
	case <-c.Done():
		return ctx.Status(fiber.StatusRequestTimeout).
			JSON(utils.StatusMessage(fiber.StatusRequestTimeout))
	default:
		return ctx.Status(fiber.StatusOK).JSON(portfolios)
	}
}

func (h *BioHandler) UpdatePortfolio(ctx *fiber.Ctx) error {
	requestID := h.middleware.GetRequestID(ctx)
	c, cancel := context.WithTimeout(contextPkg.FromFiberCtx(ctx), 5*time.Second)
	defer cancel()

	h.log.WithFields(logrus.Fields{
		"request_id": requestID,
		"path":       ctx.Path(),
	}).Debug("Processing portfolio update request")

	id := ctx.Params("id")
	if id == "" {
		h.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"path":       ctx.Path(),
		}).Warn("Missing portfolio ID in URL")
		return fiber.NewError(fiber.StatusBadRequest, "Portfolio ID is required")
	}

	req := bio.UpdatePortfolio{
		ProjectName:      ctx.FormValue("project_name"),
		ProjectLocation:  ctx.FormValue("project_location"),
		DescriptionImage: ctx.FormValue("description_image"),
		ProjectLink:      ctx.FormValue("project_link"),
		StartDate:        ctx.FormValue("start_date"),
		EndDate:          ctx.FormValue("end_date"),
		Description:      ctx.FormValue("description"),
	}

	if err := h.validator.Struct(&req); err != nil {
		h.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"error":      err.Error(),
			"id":         id,
		}).Warn("Validation failed for portfolio update")
		return err
	}

	var imageFile *multipart.FileHeader
	imageFile, err := ctx.FormFile("image")
	if err != nil && !errors.Is(err, http.ErrMissingFile) {
		h.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"error":      err.Error(),
		}).Warn("Failed to parse portfolio image")
		return err
	}

	var descriptionFile *multipart.FileHeader
	descriptionFile, err = ctx.FormFile("description_image")
	if err != nil && !errors.Is(err, http.ErrMissingFile) {
		h.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"error":      err.Error(),
		}).Warn("Failed to parse description image")
		return err
	}

	if err := h.bioService.UpdatePortfolio(c, req, id, imageFile, descriptionFile); err != nil {
		if errors.Is(err, errors.New("portfolio not found")) {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"errors": fiber.Map{"message": "Portfolio not found"},
			})
		}
		status := fiber.StatusInternalServerError
		return ctx.Status(status).JSON(fiber.Map{
			"errors": fiber.Map{"message": "Portfolio update failed", "err": err.Error()},
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

func (h *BioHandler) DeletePortfolio(ctx *fiber.Ctx) error {
	requestID := h.middleware.GetRequestID(ctx)
	c, cancel := context.WithTimeout(contextPkg.FromFiberCtx(ctx), 5*time.Second)
	defer cancel()

	h.log.WithFields(logrus.Fields{
		"request_id": requestID,
		"path":       ctx.Path(),
	}).Debug("Processing portfolio deletion request")

	id := ctx.Params("id")
	if id == "" {
		h.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"path":       ctx.Path(),
		}).Warn("Missing portfolio ID in URL")
		return fiber.NewError(fiber.StatusBadRequest, "Portfolio ID is required")
	}

	if err := h.bioService.DeletePortfolio(c, id); err != nil {
		if errors.Is(err, errors.New("portfolio not found")) {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"errors": fiber.Map{"message": "Portfolio not found"},
			})
		}
		status := fiber.StatusInternalServerError
		return ctx.Status(status).JSON(fiber.Map{
			"errors": fiber.Map{"message": "Portfolio deletion failed", "err": err.Error()},
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
