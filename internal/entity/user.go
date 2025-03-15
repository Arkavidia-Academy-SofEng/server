package entity

import (
	"time"
)

type UserRole string

const (
	RoleAdmin     UserRole = "admin"
	RoleRecruiter UserRole = "recruiter"
	RoleCandidate UserRole = "candidate"
)

type User struct {
	ID             string
	Email          string
	Password       string
	Name           string
	Role           UserRole
	ProfilePicture string
	IsPremium      bool
	PremiumUntil   time.Time
	Headline       string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      *time.Time
}

type UserLoginData struct {
	ID       string
	Username string
	Email    string
}
