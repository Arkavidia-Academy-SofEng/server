package auth

import (
	"ProjectGolang/internal/entity"
	"database/sql"
	"time"
)

type RequestOTP struct {
	Email string `json:"email" validate:"required,email"`
	Role  string `json:"role" validate:"required,oneof= recruiter candidate"`
}

type CreateUser struct {
	Code        string          `json:"code" validate:"required"`
	Email       string          `json:"email" validate:"required,email"`
	Password    string          `json:"password" validate:"required,min=8"`
	Name        string          `json:"name" validate:"required"`
	PhoneNumber string          `json:"phone_number" validate:"required"`
	Role        entity.UserRole `json:"role" validate:"required,oneof= recruiter candidate"`
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
	Name           string `form:"name" validate:"omitempty"`
	ProfilePicture string `form:"-"` // Stores file path after upload
	BannerPicture  string `form:"-"` // Stores file path after upload
	PhoneNumber    string `form:"phone_number" validate:"omitempty"`
	Location       string `form:"location" validate:"omitempty"`
	Headline       string `form:"headline" validate:"omitempty"`
}

type UpdateCompany struct {
	Name            string `form:"name" validate:"omitempty"`
	ProfilePicture  string `form:"-"` // Stores file path after upload
	BannerPicture   string `form:"-"` // Stores file path after upload
	PhoneNumber     string `form:"phone_number" validate:"omitempty"`
	Location        string `form:"location" validate:"omitempty"`
	AboutUs         string `form:"about_us" validate:"omitempty"`
	IndustryTypes   string `form:"industry_types" validate:"omitempty"`
	NumberEmployees int    `form:"number_employees" validate:"omitempty"`
	EstablishedDate string `form:"established_date" validate:"omitempty"`
	CompanyURL      string `form:"company_url" validate:"omitempty"`
	RequiredSkill   string `form:"required_skill" validate:"omitempty"`
}

type UserDB struct {
	ID             sql.NullString `db:"id"`
	Email          sql.NullString `db:"email"`
	Password       sql.NullString `db:"password"`
	PhoneNumber    sql.NullString `db:"phone_number"`
	Name           sql.NullString `db:"name"`
	Role           sql.NullString `db:"role"`
	Location       sql.NullString `db:"location"`
	ProfilePicture sql.NullString `db:"profile_picture"`
	BannerPicture  sql.NullString `db:"banner_picture"`
	IsPremium      sql.NullBool   `db:"is_premium"`
	PremiumUntil   sql.NullTime   `db:"premium_until"`
	Headline       sql.NullString `db:"headline"`
	Address        sql.NullString `db:"address"`
	CreatedAt      sql.NullTime   `db:"created_at"`
	UpdatedAt      sql.NullTime   `db:"updated_at"`
	DeletedAt      sql.NullTime   `db:"deleted_at"`
}

type CompanyDB struct {
	ID              sql.NullString `db:"id"`
	Email           sql.NullString `db:"email"`
	Password        sql.NullString `db:"password"`
	PhoneNumber     sql.NullString `db:"phone_number"`
	Name            sql.NullString `db:"name"`
	Location        sql.NullString `db:"location"`
	ProfilePicture  sql.NullString `db:"profile_picture"`
	BannerPicture   sql.NullString `db:"banner_picture"`
	AboutUs         sql.NullString `db:"about_us"`
	IndustryTypes   sql.NullString `db:"industry_types"`
	NumberEmployees sql.NullInt64  `db:"number_employees"`
	EstablishedDate sql.NullTime   `db:"established_date"`
	CompanyURL      sql.NullString `db:"company_url"`
	RequiredSkill   sql.NullString `db:"required_skill"`
	CreatedAt       sql.NullTime   `db:"created_at"`
	UpdatedAt       sql.NullTime   `db:"updated_at"`
	DeletedAt       sql.NullTime   `db:"deleted_at"`
}
