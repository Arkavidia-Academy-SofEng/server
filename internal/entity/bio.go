package entity

import "time"

type Experience struct {
	ID          string    `db:"id"`
	UserID      string    `db:"user_id"`
	ImageURL    string    `db:"image_url"`
	JobTitle    string    `db:"job_title"`
	JobLocation string    `db:"job_location"`
	SkillUsed   string    `db:"skill_used"`
	StartDate   string    `db:"start_date"`
	EndDate     string    `db:"end_date"`
	Description string    `db:"description"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

type Education struct {
	ID                string    `db:"id"`
	Image             string    `db:"image"`
	UserID            string    `db:"user_id"`
	TitleDegree       string    `db:"title_degree"`
	InstitutionalName string    `db:"institutional_name"`
	StartDate         string    `db:"start_date"`
	EndDate           string    `db:"end_date"`
	Description       string    `db:"description"`
	CreatedAt         time.Time `db:"created_at"`
	UpdatedAt         time.Time `db:"updated_at"`
}

type Portfolio struct {
	ID               string    `db:"id"`
	UserID           string    `db:"user_id"`
	Image            string    `db:"image"`
	ProjectName      string    `db:"project_name"`
	ProjectLocation  string    `db:"project_location"`
	DescriptionImage string    `db:"description_image"`
	ProjectLink      string    `db:"project_link"`
	StartDate        string    `db:"start_date"`
	EndDate          string    `db:"end_date"`
	Description      string    `db:"description"`
	CreatedAt        time.Time `db:"created_at"`
	UpdatedAt        time.Time `db:"updated_at"`
}
