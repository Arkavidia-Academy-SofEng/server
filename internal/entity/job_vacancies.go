package entity

import "time"

type JobVacancy struct {
	ID           string
	RecruiterID  string
	Title        string
	Description  string
	Requirements string
	Location     string
	JobType      string
	Deadline     time.Time
	IsActive     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
