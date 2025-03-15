package recruitmentRepository

import (
	"ProjectGolang/internal/entity"
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
)

func (r *jobVacanciesRepository) CreateJobVacancy(c context.Context, jobVacancy entity.JobVacancy) error {
	r.log.WithFields(map[string]interface{}{
		"job_vacancy_id": jobVacancy.ID,
		"title":          jobVacancy.Title,
		"recruiter_id":   jobVacancy.RecruiterID,
		"location":       jobVacancy.Location,
		"job_type":       jobVacancy.JobType,
		"deadline":       jobVacancy.Deadline,
		"is_active":      jobVacancy.IsActive,
		"created_at":     time.Now(),
		"updated_at":     time.Now(),
	}).Debug("Creating job vacancy in database")

	query, args, err := sqlx.Named(queryCreateJobVacancy, jobVacancy)
	if err != nil {
		r.log.WithFields(map[string]interface{}{
			"error": err.Error(),
		}).Error("Failed to build SQL query for CreateJobVacancy")
		return err
	}

	query = r.q.Rebind(query)
	r.log.WithFields(map[string]interface{}{}).Debug("Executing query to create job vacancy")

	_, err = r.q.ExecContext(c, query, args...)
	if err != nil {
		r.log.WithFields(map[string]interface{}{
			"error": err.Error(),
		}).Error("Database error when creating job vacancy")
		return err
	}

	r.log.WithFields(map[string]interface{}{}).Debug("Job vacancy created successfully")

	return nil
}

func (r *jobVacanciesRepository) GetJobVacancies(c context.Context, page, pageSize int) ([]entity.JobVacancy, int, error) {
	offset := (page - 1) * pageSize

	r.log.WithFields(map[string]interface{}{
		"page":     page,
		"pageSize": pageSize,
		"offset":   offset,
	}).Debug("Fetching paginated job vacancies from database")

	// First, get the total count
	var totalCount int
	countRow := r.q.QueryRowxContext(c, queryCountJobVacancies)
	if err := countRow.Scan(&totalCount); err != nil {
		r.log.WithFields(map[string]interface{}{
			"error": err.Error(),
		}).Error("Failed to get total count of job vacancies")
		return nil, 0, err
	}

	// Then get the paginated data
	query := r.q.Rebind(queryGetJobVacancies)
	r.log.Debug("Executing query to fetch paginated job vacancies")

	rows, err := r.q.QueryContext(c, query, pageSize, offset)
	if err != nil {
		r.log.WithFields(map[string]interface{}{
			"error": err.Error(),
		}).Error("Database error when fetching job vacancies")
		return nil, 0, err
	}
	defer rows.Close()

	var jobVacancies []entity.JobVacancy
	for rows.Next() {
		var jv entity.JobVacancy
		err := rows.Scan(
			&jv.ID,
			&jv.RecruiterID,
			&jv.Title,
			&jv.Description,
			&jv.Requirements,
			&jv.Location,
			&jv.JobType,
			&jv.Deadline,
			&jv.IsActive,
			&jv.CreatedAt,
			&jv.UpdatedAt,
		)
		if err != nil {
			r.log.WithFields(map[string]interface{}{
				"error": err.Error(),
			}).Error("Error scanning job vacancy row")
			return nil, 0, err
		}
		jobVacancies = append(jobVacancies, jv)
	}

	if err = rows.Err(); err != nil {
		r.log.WithFields(map[string]interface{}{
			"error": err.Error(),
		}).Error("Error iterating through job vacancy rows")
		return nil, 0, err
	}

	r.log.WithFields(map[string]interface{}{
		"count": len(jobVacancies),
		"total": totalCount,
	}).Debug("Job vacancies fetched successfully")

	return jobVacancies, totalCount, nil
}

func (r *jobVacanciesRepository) CheckJobVacancyExists(c context.Context, id string) (bool, error) {
	r.log.WithFields(map[string]interface{}{
		"job_vacancy_id": id,
	}).Debug("Checking if job vacancy exists")

	var exists bool
	query := r.q.Rebind(queryCheckJobVacancyExists)
	err := r.q.QueryRowxContext(c, query, id).Scan(&exists)
	if err != nil {
		r.log.WithFields(map[string]interface{}{
			"error": err.Error(),
			"id":    id,
		}).Error("Database error when checking job vacancy existence")
		return false, err
	}

	return exists, nil
}

func (r *jobVacanciesRepository) UpdateJobVacancy(c context.Context, jobVacancy entity.JobVacancy) error {
	r.log.WithFields(map[string]interface{}{
		"job_vacancy_id": jobVacancy.ID,
		"title":          jobVacancy.Title,
		"location":       jobVacancy.Location,
		"job_type":       jobVacancy.JobType,
		"deadline":       jobVacancy.Deadline,
		"is_active":      jobVacancy.IsActive,
		"updated_at":     jobVacancy.UpdatedAt,
	}).Debug("Updating job vacancy in database")

	query, args, err := sqlx.Named(queryUpdateJobVacancy, jobVacancy)
	if err != nil {
		r.log.WithFields(map[string]interface{}{
			"error": err.Error(),
		}).Error("Failed to build SQL query for UpdateJobVacancy")
		return err
	}

	query = r.q.Rebind(query)
	r.log.WithFields(map[string]interface{}{}).Debug("Executing query to update job vacancy")

	result, err := r.q.ExecContext(c, query, args...)
	if err != nil {
		r.log.WithFields(map[string]interface{}{
			"error": err.Error(),
		}).Error("Database error when updating job vacancy")
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.log.WithFields(map[string]interface{}{
			"error": err.Error(),
		}).Error("Failed to get rows affected after update")
		return err
	}

	if rowsAffected == 0 {
		r.log.WithFields(map[string]interface{}{
			"id": jobVacancy.ID,
		}).Warn("No job vacancy was updated")
		return fmt.Errorf("job vacancy with ID %s not found", jobVacancy.ID)
	}

	r.log.WithFields(map[string]interface{}{}).Debug("Job vacancy updated successfully")

	return nil
}

func (r *jobVacanciesRepository) DeleteJobVacancy(c context.Context, id string) error {
	r.log.WithFields(map[string]interface{}{
		"job_vacancy_id": id,
	}).Debug("Deleting job vacancy from database")

	query := r.q.Rebind(queryDeleteJobVacancy)
	r.log.WithFields(map[string]interface{}{}).Debug("Executing query to delete job vacancy")

	result, err := r.q.ExecContext(c, query, id)
	if err != nil {
		r.log.WithFields(map[string]interface{}{
			"error": err.Error(),
		}).Error("Database error when deleting job vacancy")
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.log.WithFields(map[string]interface{}{
			"error": err.Error(),
		}).Error("Failed to get rows affected after delete")
		return err
	}

	if rowsAffected == 0 {
		r.log.WithFields(map[string]interface{}{
			"id": id,
		}).Warn("No job vacancy was deleted")
		return fmt.Errorf("job vacancy with ID %s not found", id)
	}

	r.log.WithFields(map[string]interface{}{}).Debug("Job vacancy deleted successfully")

	return nil
}
