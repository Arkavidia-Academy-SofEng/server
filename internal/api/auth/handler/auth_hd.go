package authHandler

import (
	"ProjectGolang/internal/api/auth"
	contextPkg "ProjectGolang/pkg/context"
	"ProjectGolang/pkg/log"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"golang.org/x/net/context"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"
)

func (h *AuthHandler) RequestOTP(ctx *fiber.Ctx) error {
	requestID := h.middleware.GetRequestID(ctx)
	c, cancel := context.WithTimeout(contextPkg.FromFiberCtx(ctx), 5*time.Second)
	defer cancel()

	h.log.WithFields(log.Fields{
		"request_id": requestID,
		"path":       ctx.Path(),
	}).Debug("Processing user creation request")

	var req auth.RequestOTP
	if err := ctx.BodyParser(&req); err != nil {
		h.log.WithFields(log.Fields{
			"error": err.Error(),
			"path":  ctx.Path(),
		}).Warn("Failed to parse OTP request body")
		return err
	}

	if err := h.validator.Struct(&req); err != nil {
		h.log.WithFields(log.Fields{
			"error": err.Error(),
			"email": req.Email,
		}).Warn("Validation failed for OTP request")
		return err
	}

	if err := h.authService.RequestOTP(c, req); err != nil {
		if errors.Is(err, auth.ErrorEmailAlreadyExists) {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"errors": fiber.Map{"message": "Email already exists"},
			})
		}
		status := fiber.StatusInternalServerError
		return ctx.Status(status).JSON(fiber.Map{
			"errors": fiber.Map{"message": "Failed to request OTP", "err": err.Error()},
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

func (h *AuthHandler) CreateUser(ctx *fiber.Ctx) error {
	requestID := h.middleware.GetRequestID(ctx)
	c, cancel := context.WithTimeout(contextPkg.FromFiberCtx(ctx), 5*time.Second)
	defer cancel()

	h.log.WithFields(log.Fields{
		"request_id": requestID,
		"path":       ctx.Path(),
	}).Debug("Processing user creation request")

	var req auth.CreateUser
	if err := ctx.BodyParser(&req); err != nil {
		h.log.WithFields(log.Fields{
			"request_id": requestID,
			"error":      err.Error(),
		}).Warn("Failed to parse user creation request body")
		return err
	}

	if err := h.validator.Struct(&req); err != nil {
		h.log.WithFields(log.Fields{
			"request_id": requestID,
			"error":      err.Error(),
			"email":      req.Email,
		}).Warn("Validation failed for user creation")
		return err
	}

	if req.Role == "candidate" {
		err := h.authService.CreateUser(c, req)
		if err != nil {
			if errors.Is(err, auth.ErrorInvalidOTP) {
				return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"errors": fiber.Map{"message": "Invalid OTP code"},
				})
			}

			status := fiber.StatusInternalServerError
			return ctx.Status(status).JSON(fiber.Map{
				"errors": fiber.Map{"message": "User creation failed", "err": err.Error()},
			})
		}
	} else {
		err := h.authService.CreateCompany(c, req)
		if err != nil {
			if errors.Is(err, auth.ErrorInvalidOTP) {
				return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"errors": fiber.Map{"message": "Invalid OTP code"},
				})
			}

			status := fiber.StatusInternalServerError
			return ctx.Status(status).JSON(fiber.Map{
				"errors": fiber.Map{"message": "User creation failed", "err": err.Error()},
			})
		}
	}

	select {
	case <-c.Done():
		return ctx.Status(fiber.StatusRequestTimeout).
			JSON(utils.StatusMessage(fiber.StatusRequestTimeout))
	default:
		return ctx.SendStatus(fiber.StatusCreated)
	}
}

func (h *AuthHandler) Login(ctx *fiber.Ctx) error {
	requestID := h.middleware.GetRequestID(ctx)
	c, cancel := context.WithTimeout(contextPkg.FromFiberCtx(ctx), 5*time.Second)
	defer cancel()

	h.log.WithField("path", ctx.Path()).Debug("Processing login request")

	var req auth.LoginRequest
	if err := ctx.BodyParser(&req); err != nil {
		h.log.WithFields(log.Fields{
			"request_id": requestID,
			"error":      err.Error(),
			"path":       ctx.Path(),
		}).Warn("Failed to parse login request body")
		return err
	}

	if err := h.validator.Struct(&req); err != nil {
		h.log.WithFields(log.Fields{
			"request_id": requestID,
			"error":      err.Error(),
			"email":      req.Email,
		}).Warn("Validation failed for login request")
		return err
	}

	loginResponse, err := h.authService.Login(c, req)
	if err != nil {
		if errors.Is(err, auth.ErrorInvalidCredentials) {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"errors": fiber.Map{"message": "Invalid email or password"},
			})
		}

		status := fiber.StatusInternalServerError
		return ctx.Status(status).JSON(fiber.Map{
			"errors": fiber.Map{"message": "Login failed", "err": err.Error()},
		})
	}

	select {
	case <-c.Done():
		return ctx.Status(fiber.StatusRequestTimeout).
			JSON(utils.StatusMessage(fiber.StatusRequestTimeout))
	default:
		return ctx.Status(fiber.StatusOK).JSON(loginResponse)
	}
}

func (h *AuthHandler) UpdateUser(ctx *fiber.Ctx) error {
	requestID := h.middleware.GetRequestID(ctx)
	c, cancel := context.WithTimeout(contextPkg.FromFiberCtx(ctx), 5*time.Second)
	defer cancel()

	h.log.WithField("path", ctx.Path()).Debug("Processing user update request")

	id := ctx.Params("id")
	if id == "" {
		h.log.WithFields(log.Fields{
			"path": ctx.Path(),
		}).Warn("Missing user ID in URL")
		return fiber.NewError(fiber.StatusBadRequest, "User ID is required")
	}

	req := auth.UpdateUser{
		Name:        ctx.FormValue("name"),
		PhoneNumber: ctx.FormValue("phone_number"),
		Location:    ctx.FormValue("location"),
		Headline:    ctx.FormValue("headline"),
	}

	profileFile, err := ctx.FormFile("profile_picture")
	if err != nil {
		h.log.WithFields(log.Fields{
			"request_id": requestID,
			"error":      err.Error(),
			"file":       profileFile.Filename,
		}).Warn("Failed to parse profile picture")
		return err
	}

	bannerFile, err := ctx.FormFile("banner_picture")
	if err != nil {
		h.log.WithFields(log.Fields{
			"request_id": requestID,
			"error":      err.Error(),
			"file":       bannerFile.Filename,
		}).Warn("Failed to parse banner picture")
		return err
	}

	if err := h.validator.Struct(&req); err != nil {
		h.log.WithFields(log.Fields{
			"request_id": requestID,
			"error":      err.Error(),
			"id":         id,
		}).Warn("Validation failed for user update")
		return err
	}

	err = h.authService.UpdateUser(c, req, id, bannerFile, profileFile)
	if err != nil {
		if errors.Is(err, auth.ErrorUserNotFound) {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"errors": fiber.Map{"message": "User not found"},
			})
		}

		status := fiber.StatusInternalServerError
		return ctx.Status(status).JSON(fiber.Map{
			"errors": fiber.Map{"message": "User update failed", "err": err.Error()},
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

func (h *AuthHandler) DeleteUser(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	h.log.WithField("path", ctx.Path()).Debug("Processing user deletion request")

	id := ctx.Params("id")
	if id == "" {
		h.log.WithFields(log.Fields{
			"path": ctx.Path(),
		}).Error("Missing user ID in URL")
		return fiber.NewError(fiber.StatusBadRequest, "User ID is required")
	}

	if err := h.authService.DeleteUser(c, id); err != nil {
		h.log.WithFields(log.Fields{
			"error": err.Error(),
			"id":    id,
		}).Error("User deletion failed")
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

func (h *AuthHandler) UpdateCompany(ctx *fiber.Ctx) error {
	requestID := h.middleware.GetRequestID(ctx)
	c, cancel := context.WithTimeout(contextPkg.FromFiberCtx(ctx), 5*time.Second)
	defer cancel()

	h.log.WithField("path", ctx.Path()).Debug("Processing company update request")

	id := ctx.Params("id")
	if id == "" {
		h.log.WithFields(log.Fields{
			"path": ctx.Path(),
		}).Warn("Missing company ID in URL")
		return fiber.NewError(fiber.StatusBadRequest, "Company ID is required")
	}

	req := auth.UpdateCompany{
		Name:            ctx.FormValue("name"),
		PhoneNumber:     ctx.FormValue("phone_number"),
		Location:        ctx.FormValue("location"),
		AboutUs:         ctx.FormValue("about_us"),
		IndustryTypes:   ctx.FormValue("industry_types"),
		EstablishedDate: ctx.FormValue("established_date"),
		CompanyURL:      ctx.FormValue("company_url"),
		RequiredSkill:   ctx.FormValue("required_skill"),
	}

	// Parse number_employees as int if provided
	if numEmployeesStr := ctx.FormValue("number_employees"); numEmployeesStr != "" {
		numEmployees, err := strconv.Atoi(numEmployeesStr)
		if err != nil {
			h.log.WithFields(log.Fields{
				"request_id": requestID,
				"error":      err.Error(),
				"value":      numEmployeesStr,
			}).Warn("Invalid number_employees value")
			return fiber.NewError(fiber.StatusBadRequest, "Invalid number of employees")
		}
		req.NumberEmployees = numEmployees
	}

	var profileFile *multipart.FileHeader
	profileFile, err := ctx.FormFile("profile_picture")
	if err != nil && !errors.Is(err, http.ErrMissingFile) {
		h.log.WithFields(log.Fields{
			"request_id": requestID,
			"error":      err.Error(),
		}).Warn("Failed to parse profile picture")
		return err
	}

	var bannerFile *multipart.FileHeader
	bannerFile, err = ctx.FormFile("banner_picture")
	if err != nil && !errors.Is(err, http.ErrMissingFile) {
		h.log.WithFields(log.Fields{
			"request_id": requestID,
			"error":      err.Error(),
		}).Warn("Failed to parse banner picture")
		return err
	}

	if err := h.validator.Struct(&req); err != nil {
		h.log.WithFields(log.Fields{
			"request_id": requestID,
			"error":      err.Error(),
			"id":         id,
		}).Warn("Validation failed for company update")
		return err
	}

	err = h.authService.UpdateCompany(c, req, id, bannerFile, profileFile)
	if err != nil {
		if errors.Is(err, auth.ErrorUserNotFound) {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"errors": fiber.Map{"message": "Company not found"},
			})
		}

		status := fiber.StatusInternalServerError
		return ctx.Status(status).JSON(fiber.Map{
			"errors": fiber.Map{"message": "Company update failed", "err": err.Error()},
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
