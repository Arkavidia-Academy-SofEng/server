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
