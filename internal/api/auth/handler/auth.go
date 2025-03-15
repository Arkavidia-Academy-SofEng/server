package authHandler

import (
	"ProjectGolang/internal/api/auth"
	jwtPkg "ProjectGolang/pkg/jwt"
	"ProjectGolang/pkg/log"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"golang.org/x/net/context"
	"time"
)

func (h *AuthHandler) HandleRegister(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	h.log.WithField("path", ctx.Path()).Debug("Processing registration request")

	var req auth.CreateUserRequest
	if err := ctx.BodyParser(&req); err != nil {
		h.log.WithFields(log.Fields{
			"error": err.Error(),
			"path":  ctx.Path(),
		}).Error("Failed to parse registration request body")
		return err
	}

	if err := h.validator.Struct(&req); err != nil {
		h.log.WithFields(log.Fields{
			"error": err.Error(),
			"email": req.Email,
		}).Warn("Validation failed for user registration")
		return err
	}

	requestLog := req
	requestLog.Password = "[REDACTED]"

	h.log.WithFields(log.Fields{
		"username": req.Username,
		"email":    req.Email,
	}).Info("Processing user registration request")

	if err := h.authService.RegisterUser(c, req); err != nil {
		h.log.WithFields(log.Fields{
			"error": err.Error(),
			"email": req.Email,
		}).Error("User registration failed")
		return err
	}

	h.log.WithFields(log.Fields{
		"email":    req.Email,
		"username": req.Username,
	}).Info("User registered successfully")

	select {
	case <-c.Done():
		return ctx.Status(fiber.StatusRequestTimeout).
			JSON(utils.StatusMessage(fiber.StatusRequestTimeout))
	default:
		return ctx.SendStatus(fiber.StatusCreated)
	}
}

func (h *AuthHandler) HandleLogin(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	h.log.WithField("path", ctx.Path()).Debug("Processing login request")

	var req auth.LoginUserRequest
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
		}).Warn("Validation failed for user login")
		return err
	}

	h.log.WithFields(log.Fields{
		"email": req.Email,
	}).Debug("Attempting user login")

	res, err := h.authService.Login(ctx.Context(), req)
	if err != nil {
		h.log.WithFields(log.Fields{
			"error": err.Error(),
			"email": req.Email,
		}).Error("User login failed")
		return err
	}

	h.log.WithFields(log.Fields{
		"email":                 req.Email,
		"token_expires_minutes": res.ExpiresInMinutes,
	}).Info("User logged in successfully")

	select {
	case <-c.Done():
		return ctx.Status(fiber.StatusRequestTimeout).
			JSON(utils.StatusMessage(fiber.StatusRequestTimeout))
	default:
		return ctx.Status(fiber.StatusOK).JSON(res)
	}
}

func (h *AuthHandler) HandleUpdateUser(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	h.log.WithField("path", ctx.Path()).Debug("Processing user update request")

	var req auth.UpdateUserRequest
	if err := ctx.BodyParser(&req); err != nil {
		h.log.WithFields(log.Fields{
			"error": err.Error(),
			"path":  ctx.Path(),
		}).Error("Failed to parse user update request body")
		return err
	}

	if err := h.validator.Struct(&req); err != nil {
		h.log.WithFields(log.Fields{
			"error": err.Error(),
		}).Warn("Validation failed for user update")
		return err
	}

	userData, err := jwtPkg.GetUserLoginData(ctx)
	if err != nil {
		h.log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Failed to get user data from token")
		return err
	}

	h.log.WithFields(log.Fields{
		"user_id":  userData.ID,
		"username": userData.Username,
	}).Debug("Updating user")

	if err := h.authService.UpdateUser(ctx.Context(), userData, req); err != nil {
		h.log.WithFields(log.Fields{
			"error":   err.Error(),
			"user_id": userData.ID,
		}).Error("Error updating user")
		return err
	}

	h.log.WithFields(log.Fields{
		"user_id": userData.ID,
	}).Info("User updated successfully")

	select {
	case <-c.Done():
		return ctx.Status(fiber.StatusRequestTimeout).
			JSON(utils.StatusMessage(fiber.StatusRequestTimeout))
	default:
		return ctx.SendStatus(fiber.StatusOK)
	}
}

func (h *AuthHandler) HandleDeleteUser(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	id := ctx.Params("id")
	h.log.WithFields(log.Fields{
		"user_id": id,
		"path":    ctx.Path(),
	}).Debug("Processing user deletion request")

	if err := h.authService.DeleteUser(ctx.Context(), id); err != nil {
		h.log.WithFields(log.Fields{
			"error":   err.Error(),
			"user_id": id,
		}).Error("Failed to delete user")
		return err
	}

	h.log.WithFields(log.Fields{
		"user_id": id,
	}).Info("User deleted successfully")

	select {
	case <-c.Done():
		return ctx.Status(fiber.StatusRequestTimeout).
			JSON(utils.StatusMessage(fiber.StatusRequestTimeout))
	default:
		return ctx.SendStatus(fiber.StatusOK)
	}
}
