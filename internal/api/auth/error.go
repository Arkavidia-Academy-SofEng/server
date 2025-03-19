package auth

import (
	"ProjectGolang/pkg/response"
	"github.com/gofiber/fiber/v2"
)

var (
	ErrorEmailAlreadyExists = response.New(fiber.StatusBadRequest, "email already exists")
	ErrorInvalidCredentials = response.New(fiber.StatusBadRequest, "invalid credentials")
	ErrorUserNotFound       = response.New(fiber.StatusNotFound, "user not found")
	ErrorInvalidOTP         = response.New(fiber.StatusBadRequest, "invalid otp")
)
