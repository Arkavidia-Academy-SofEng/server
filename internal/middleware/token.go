package middleware

import (
	"ProjectGolang/internal/entity"
	jwtPkg "ProjectGolang/pkg/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"reflect"
	"strings"
)

const (
	AccessTokenSecret = "JWT_ACCESS_TOKEN_SECRET"
)

type tokenMiddleware struct {
}

func newTokenMiddleware() *tokenMiddleware {
	return &tokenMiddleware{}
}

func (m *middleware) NewTokenMiddleware(ctx *fiber.Ctx) error {
	requestID := ctx.Locals("request_id")
	authHeader := ctx.Get("Authorization")

	if authHeader == "" {
		m.log.WithFields(
			logrus.Fields{
				"request_id":  requestID,
				"auth_header": authHeader,
			}).Warn("Authorization header not found")
		return ctx.Status(fiber.StatusUnauthorized).SendString("unauthorized, access token invalid or expired")
	}

	headerParts := strings.Split(authHeader, " ")
	m.log.WithFields(logrus.Fields{
		"auth_type": headerParts[0],
		"has_token": len(headerParts) > 1,
	}).Debug("Authorization header check")

	if !strings.HasPrefix(authHeader, "Bearer ") {
		m.log.WithFields(
			logrus.Fields{
				"request_id":  requestID,
				"auth_header": authHeader,
			}).Warn("Invalid Authorization format - must start with 'Bearer '")
		return ctx.Status(fiber.StatusUnauthorized).SendString("unauthorized, access token invalid or expired")
	}

	userToken, err := jwtPkg.VerifyTokenHeader(ctx, AccessTokenSecret)
	if err != nil {
		m.log.WithFields(
			logrus.Fields{
				"request_id":  requestID,
				"auth_header": authHeader,
			}).Warn("Token verification failed")
		return ctx.Status(fiber.StatusUnauthorized).SendString("unauthorized, access token invalid or expired")
	}

	claims, ok := userToken.Claims.(jwt.MapClaims)
	if !ok {
		m.log.WithFields(
			logrus.Fields{
				"request_id":  requestID,
				"auth_header": authHeader,
			}).Warn("Couold not extract claims from token")
		return ctx.Status(fiber.StatusUnauthorized).SendString("unauthorized, access token invalid or expired")
	}

	m.log.WithFields(logrus.Fields{
		"claim_keys":   reflect.ValueOf(claims).MapKeys(),
		"exp":          claims["exp"],
		"id_exists":    claims["id"] != nil,
		"email_exists": claims["email"] != nil,
	}).Debug("Token claims")

	user := entity.UserLoginData{
		ID:        claims["id"].(string),
		Email:     claims["email"].(string),
		Name:      claims["name"].(string),
		Role:      entity.UserRole(claims["role"].(string)),
		IsPremium: claims["is_premium"].(bool),
	}
	ctx.Locals("user", user)

	m.log.WithFields(logrus.Fields{
		"request_id": requestID,
		"user":       user,
	}).Info("Authentication successful")

	return ctx.Next()
}
