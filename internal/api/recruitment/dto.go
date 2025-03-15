package recruitment

import "time"

type CreateJobVacancy struct {
	RecruiterID  string
	Title        string
	Description  string
	Requirements string
	Location     string
	JobType      string
	Deadline     time.Time
	IsActive     bool
}

type GetJobVacancies struct {
	Page     int `json:"page" validate:"required,min=1"`
	PageSize int `json:"page_size" validate:"required,min=1,max=100"`
}

type JobVacancyResponse struct {
	ID           string    `json:"id"`
	RecruiterID  string    `json:"recruiter_id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	Requirements string    `json:"requirements"`
	Location     string    `json:"location"`
	JobType      string    `json:"job_type"`
	Deadline     time.Time `json:"deadline"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type PaginatedJobVacanciesResponse struct {
	JobVacancies []JobVacancyResponse `json:"job_vacancies"`
	TotalCount   int                  `json:"total_count"`
	TotalPages   int                  `json:"total_pages"`
	CurrentPage  int                  `json:"current_page"`
	PageSize     int                  `json:"page_size"`
}

type UpdateJobVacancy struct {
	ID           string    `json:"id" validate:"required"`
	Title        string    `json:"title" validate:"required,min=3,max=100"`
	Description  string    `json:"description" validate:"required"`
	Requirements string    `json:"requirements" validate:"required"`
	Location     string    `json:"location" validate:"required"`
	JobType      string    `json:"job_type" validate:"required,oneof=FULL_TIME PART_TIME CONTRACT REMOTE"`
	Deadline     time.Time `json:"deadline" validate:"required,gtfield=CreatedAt"`
	IsActive     bool      `json:"is_active"`
}
