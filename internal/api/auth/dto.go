package auth

import (
	"ProjectGolang/internal/entity"
)

type CreateUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=32"`
	Username string `json:"username" validate:"required,min=3,max=255"`
}

type LoginUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type SinginUserResponse struct {
	AccessToken      string `json:"accessToken"`
	RefreshToken     string `json:"refreshToken"`
	ExpiresInSeconds int64  `json:"expiresInSeconds"`
	SessionID        string `json:"sessionID"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}

type LoginUserResponse struct {
	AccessToken      string  `json:"accessToken"`
	ExpiresInMinutes float64 `json:"expiresInHour"`
}

type UserClaims struct {
	Email    string              `json:"email"`
	ID       string              `json:"id"`
	Provider entity.AuthProvider `json:"provider"`
}

type UpdateUserRequest struct {
	Username string `json:"username" validate:"required,min=3,max=255"`
	Password string `json:"password" validate:"required,min=8,max=32"`
}

type UpdateEmailUserRequest struct {
	Email string `json:"email" validate:"required,email"`
}
