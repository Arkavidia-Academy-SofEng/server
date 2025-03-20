package bio

import "database/sql"

type CreateExperience struct {
	ImageURL    string `form:"image"`
	JobTitle    string `form:"job_title" validate:"required"`
	JobLocation string `form:"job_location"`
	SkillUsed   string `form:"skill_used"`
	StartDate   string `form:"start_date" validate:"required"`
	EndDate     string `form:"end_date"`
	Description string `form:"description"`
}

type ExperienceDB struct {
	ID          sql.NullString `db:"id"`
	UserID      sql.NullString `db:"user_id"`
	ImageURL    sql.NullString `db:"image_url"`
	JobTitle    sql.NullString `db:"job_title"`
	JobLocation sql.NullString `db:"job_location"`
	SkillUsed   sql.NullString `db:"skill_used"`
	StartDate   sql.NullString `db:"start_date"`
	EndDate     sql.NullString `db:"end_date"`
	Description sql.NullString `db:"description"`
	CreatedAt   sql.NullTime   `db:"created_at"`
	UpdatedAt   sql.NullTime   `db:"updated_at"`
}

type UpdateExperience struct {
	ImageURL    string `form:"image"`
	JobTitle    string `form:"job_title"`
	JobLocation string `form:"job_location"`
	SkillUsed   string `form:"skill_used"`
	StartDate   string `form:"start_date"`
	EndDate     string `form:"end_date"`
	Description string `form:"description"`
}

type CreateEducation struct {
	Image             string `form:"image"`
	TitleDegree       string `form:"title_degree" validate:"required"`
	InstitutionalName string `form:"institutional_name" validate:"required"`
	StartDate         string `form:"start_date" validate:"required"`
	EndDate           string `form:"end_date"`
	Description       string `form:"description"`
}

type EducationDB struct {
	ID                sql.NullString `db:"id"`
	Image             sql.NullString `db:"image"`
	UserID            sql.NullString `db:"user_id"`
	TitleDegree       sql.NullString `db:"title_degree"`
	InstitutionalName sql.NullString `db:"institutional_name"`
	StartDate         sql.NullString `db:"start_date"`
	EndDate           sql.NullString `db:"end_date"`
	Description       sql.NullString `db:"description"`
	CreatedAt         sql.NullTime   `db:"created_at"`
	UpdatedAt         sql.NullTime   `db:"updated_at"`
}

type UpdateEducation struct {
	Image             string `form:"image"`
	TitleDegree       string `form:"title_degree"`
	InstitutionalName string `form:"institutional_name"`
	StartDate         string `form:"start_date"`
	EndDate           string `form:"end_date"`
	Description       string `form:"description"`
}

type CreatePortfolio struct {
	Image            string `form:"image"`
	ProjectName      string `form:"project_name" validate:"required"`
	ProjectLocation  string `form:"project_location" validate:"required"`
	DescriptionImage string `form:"description_image"`
	ProjectLink      string `form:"project_link"`
	StartDate        string `form:"start_date" validate:"required"`
	EndDate          string `form:"end_date"`
	Description      string `form:"description"`
}

type PortfolioDB struct {
	ID               sql.NullString `db:"id"`
	UserID           sql.NullString `db:"user_id"`
	Image            sql.NullString `db:"image"`
	ProjectName      sql.NullString `db:"project_name"`
	ProjectLocation  sql.NullString `db:"project_location"`
	DescriptionImage sql.NullString `db:"description_image"`
	ProjectLink      sql.NullString `db:"project_link"`
	StartDate        sql.NullString `db:"start_date"`
	EndDate          sql.NullString `db:"end_date"`
	Description      sql.NullString `db:"description"`
	CreatedAt        sql.NullTime   `db:"created_at"`
	UpdatedAt        sql.NullTime   `db:"updated_at"`
}

type UpdatePortfolio struct {
	Image            string `form:"image"`
	ProjectName      string `form:"project_name"`
	ProjectLocation  string `form:"project_location"`
	DescriptionImage string `form:"description_image"`
	ProjectLink      string `form:"project_link"`
	StartDate        string `form:"start_date"`
	EndDate          string `form:"end_date"`
	Description      string `form:"description"`
}
