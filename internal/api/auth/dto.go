package auth

import (
	"ProjectGolang/internal/entity"
	"time"
)

type CreateUser struct {
	Email          string          `json:"email" validate:"required,email"`
	Password       string          `json:"password" validate:"required,min=8"`
	Name           string          `json:"name" validate:"required"`
	Role           entity.UserRole `json:"role" validate:"required,oneof=admin recruiter candidate"`
	ProfilePicture string          `json:"profile_picture"`
	IsPremium      bool            `json:"is_premium"`
	PremiumUntil   time.Time       `json:"premium_until"`
	Headline       string          `json:"headline"`
}

type UserResponse struct {
	ID             string          `json:"id"`
	Email          string          `json:"email"`
	Name           string          `json:"name"`
	Role           entity.UserRole `json:"role"`
	ProfilePicture string          `json:"profile_picture"`
	IsPremium      bool            `json:"is_premium"`
	PremiumUntil   time.Time       `json:"premium_until"`
	Headline       string          `json:"headline"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

type UpdateUser struct {
	ID             string          `json:"-"` // Will be set from URL parameter
	Name           string          `json:"name" validate:"required"`
	Role           entity.UserRole `json:"role" validate:"required,oneof=admin recruiter candidate"`
	ProfilePicture string          `json:"profile_picture"`
	IsPremium      bool            `json:"is_premium"`
	PremiumUntil   time.Time       `json:"premium_until"`
	Headline       string          `json:"headline"`
}
