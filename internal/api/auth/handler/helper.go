package authHandler

import (
	"ProjectGolang/internal/api/auth"
	"github.com/gofiber/fiber/v2"
)

func (h *AuthHandler) parseAndBindRequest(ctx *fiber.Ctx) (auth.CreateUserRequest, error) {
	var req auth.CreateUserRequest
	if err := ctx.BodyParser(&req); err != nil {
		return req, err
	}

	if err := h.validator.Struct(req); err != nil {
		return req, err
	}

	return req, nil
}
