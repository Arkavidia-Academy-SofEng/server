package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const RequestIDKey = "X-Request-ID"

func NewRequestIDMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Check if already exists in header
		requestID := c.Get(RequestIDKey)

		// If not in header, generate a new one
		if requestID == "" {
			requestID = uuid.New().String()
		}

		c.Locals(RequestIDKey, requestID)
		c.Set(RequestIDKey, requestID)

		return c.Next()
	}
}
