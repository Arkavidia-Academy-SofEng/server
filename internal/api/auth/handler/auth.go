package authHandler

import (
	"ProjectGolang/internal/api/auth"
	"ProjectGolang/pkg/log"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"golang.org/x/net/context"
	"time"
)

func (h *AuthHandler) CreateUser(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	h.log.WithField("path", ctx.Path()).Debug("Processing user creation request")

	var req auth.CreateUser
	if err := ctx.BodyParser(&req); err != nil {
		h.log.WithFields(log.Fields{
			"error": err.Error(),
			"path":  ctx.Path(),
		}).Error("Failed to parse user creation request body")
		return err
	}

	if err := h.validator.Struct(&req); err != nil {
		h.log.WithFields(log.Fields{
			"error": err.Error(),
			"email": req.Email,
		}).Warn("Validation failed for user creation")
		return err
	}

	newUser, err := h.authService.CreateUser(c, req)
	if err != nil {
		h.log.WithFields(log.Fields{
			"error": err.Error(),
			"email": req.Email,
		}).Error("User creation failed")
		return err
	}

	select {
	case <-c.Done():
		return ctx.Status(fiber.StatusRequestTimeout).
			JSON(utils.StatusMessage(fiber.StatusRequestTimeout))
	default:
		return ctx.Status(fiber.StatusCreated).JSON(newUser)
	}
}

func (h *AuthHandler) Login(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	h.log.WithField("path", ctx.Path()).Debug("Processing login request")

	var req auth.LoginRequest
	if err := ctx.BodyParser(&req); err != nil {
		h.log.WithFields(log.Fields{
			"error": err.Error(),
			"path":  ctx.Path(),
		}).Error("Failed to parse login request body")
		return err
	}

	if err := h.validator.Struct(&req); err != nil {
		h.log.WithFields(log.Fields{
			"error": err.Error(),
			"email": req.Email,
		}).Warn("Validation failed for login request")
		return err
	}

	loginResponse, err := h.authService.Login(c, req)
	if err != nil {
		h.log.WithFields(log.Fields{
			"error": err.Error(),
			"email": req.Email,
		}).Error("Login failed")
		return err
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
	c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	h.log.WithField("path", ctx.Path()).Debug("Processing user update request")

	// Get ID from URL parameters
	id := ctx.Params("id")
	if id == "" {
		h.log.WithFields(log.Fields{
			"path": ctx.Path(),
		}).Error("Missing user ID in URL")
		return fiber.NewError(fiber.StatusBadRequest, "User ID is required")
	}

	var req auth.UpdateUser
	if err := ctx.BodyParser(&req); err != nil {
		h.log.WithFields(log.Fields{
			"error": err.Error(),
			"path":  ctx.Path(),
		}).Error("Failed to parse user update request body")
		return err
	}

	// Set the ID from URL parameter
	req.ID = id

	if err := h.validator.Struct(&req); err != nil {
		h.log.WithFields(log.Fields{
			"error": err.Error(),
			"id":    req.ID,
		}).Warn("Validation failed for user update")
		return err
	}

	updatedUser, err := h.authService.UpdateUser(c, req)
	if err != nil {
		h.log.WithFields(log.Fields{
			"error": err.Error(),
			"id":    req.ID,
		}).Error("User update failed")
		return err
	}

	select {
	case <-c.Done():
		return ctx.Status(fiber.StatusRequestTimeout).
			JSON(utils.StatusMessage(fiber.StatusRequestTimeout))
	default:
		return ctx.Status(fiber.StatusOK).JSON(updatedUser)
	}
}

func (h *AuthHandler) DeleteUser(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	h.log.WithField("path", ctx.Path()).Debug("Processing user deletion request")

	// Get ID from URL parameters
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
