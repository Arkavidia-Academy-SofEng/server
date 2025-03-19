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
	ID             string     `db:"id"`
	Email          string     `db:"email"`
	Password       string     `db:"password"`
	Name           string     `db:"name"`
	Role           UserRole   `db:"role"`
	ProfilePicture string     `db:"profile_picture"`
	BannerPicture  string     `db:"banner_picture"`
	PhoneNumber    string     `db:"phone_number"`
	IsPremium      bool       `db:"is_premium"`
	PremiumUntil   time.Time  `db:"premium_until"`
	Location       string     `db:"location"`
	Headline       string     `db:"headline"`
	CreatedAt      time.Time  `db:"created_at"`
	UpdatedAt      time.Time  `db:"updated_at"`
	DeletedAt      *time.Time `db:"deleted_at"`
}

type UserLoginData struct {
	ID        string
	Name      string
	Email     string
	Role      UserRole
	IsPremium bool
}
