package entity

import "time"

type Company struct {
	ID              string     `db:"id"`
	Email           string     `db:"email"`
	Password        string     `db:"password"`
	Name            string     `db:"name"`
	ProfilePicture  string     `db:"profile_picture"`
	BannerPicture   string     `db:"banner_picture"`
	PhoneNumber     string     `db:"phone_number"`
	Location        string     `db:"location"`
	AboutUs         string     `db:"about_us"`
	IndustryTypes   string     `db:"industry_types"`
	NumberEmployees int        `db:"number_employees"`
	EstablishedDate time.Time  `db:"established_date"`
	CompanyURL      string     `db:"company_url"`
	RequiredSkill   string     `db:"required_skill"`
	CreatedAt       time.Time  `db:"created_at"`
	UpdatedAt       time.Time  `db:"updated_at"`
	DeletedAt       *time.Time `db:"deleted_at"`
}
