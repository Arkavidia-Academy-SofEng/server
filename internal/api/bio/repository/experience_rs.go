package bioRepository

import (
	"ProjectGolang/internal/api/bio"
	"ProjectGolang/internal/entity"
	contextPkg "ProjectGolang/pkg/context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

func (r *experienceRepository) CreateExperience(ctx context.Context, experience entity.Experience) error {
	requestID := contextPkg.GetRequestID(ctx)

	query, args, err := sqlx.Named(queryCreateExperience, experience)
	if err != nil {
		r.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"error":      err.Error(),
		}).Error("Failed to build SQL query for CreateExperience")
		return err
	}

	query = r.q.Rebind(query)

	_, err = r.q.ExecContext(ctx, query, args...)
	if err != nil {
		r.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"error":      err.Error(),
		}).Error("Database error when creating experience")
		return err
	}

	return nil
}

func (r *experienceRepository) GetExperienceByID(ctx context.Context, id string) (entity.Experience, error) {
	requestID := contextPkg.GetRequestID(ctx)
	r.log.WithFields(logrus.Fields{
		"request_id": requestID,
		"id":         id,
	}).Debug("Getting experience by ID")

	query := r.q.Rebind(queryGetExperienceByID)
	r.log.Debug("Executing query to get experience by ID")

	var exp bio.ExperienceDB
	err := r.q.QueryRowxContext(ctx, query, id).Scan(
		&exp.ID,
		&exp.UserID,
		&exp.ImageURL,
		&exp.JobTitle,
		&exp.SkillUsed,
		&exp.StartDate,
		&exp.EndDate,
		&exp.Description,
		&exp.CreatedAt,
		&exp.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.log.WithFields(logrus.Fields{
				"request_id": requestID,
				"id":         id,
			}).Warn("Experience not found")
			return entity.Experience{}, nil
		}

		r.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"error":      err.Error(),
			"id":         id,
		}).Error("Database error when getting experience by ID")
		return entity.Experience{}, err
	}

	experience := r.makeExperience(exp)

	r.log.WithFields(logrus.Fields{
		"request_id": requestID,
		"id":         id,
	}).Debug("Experience retrieved successfully")

	return experience, nil
}

func (r *experienceRepository) GetExperiencesByUserID(ctx context.Context, userID string) ([]entity.Experience, error) {
	requestID := contextPkg.GetRequestID(ctx)
	r.log.WithFields(logrus.Fields{
		"request_id": requestID,
		"user_id":    userID,
	}).Debug("Getting experiences by user ID")

	query := r.q.Rebind(queryGetExperiencesByUserID)
	r.log.Debug("Executing query to get experiences by user ID")

	rows, err := r.q.QueryxContext(ctx, query, userID)
	if err != nil {
		r.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"error":      err.Error(),
			"user_id":    userID,
		}).Error("Database error when getting experiences by user ID")
		return nil, err
	}
	defer rows.Close()

	var experiences []entity.Experience
	for rows.Next() {
		var exp bio.ExperienceDB
		err := rows.Scan(
			&exp.ID,
			&exp.UserID,
			&exp.ImageURL,
			&exp.JobTitle,
			&exp.SkillUsed,
			&exp.StartDate,
			&exp.EndDate,
			&exp.Description,
			&exp.CreatedAt,
			&exp.UpdatedAt,
		)
		if err != nil {
			r.log.WithFields(logrus.Fields{
				"request_id": requestID,
				"error":      err.Error(),
			}).Error("Error scanning experience row")
			return nil, err
		}

		experience := r.makeExperience(exp)
		experiences = append(experiences, experience)
	}

	if err = rows.Err(); err != nil {
		r.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"error":      err.Error(),
		}).Error("Error iterating experience rows")
		return nil, err
	}

	r.log.WithFields(logrus.Fields{
		"request_id": requestID,
		"user_id":    userID,
		"count":      len(experiences),
	}).Debug("Experiences retrieved successfully")

	return experiences, nil
}

func (r *experienceRepository) UpdateExperience(ctx context.Context, experience entity.Experience) error {
	requestID := contextPkg.GetRequestID(ctx)
	r.log.WithFields(logrus.Fields{
		"request_id": requestID,
		"id":         experience.ID,
		"job_title":  experience.JobTitle,
		"updated_at": experience.UpdatedAt,
	}).Debug("Updating experience in database")

	query, args, err := sqlx.Named(queryUpdateExperience, experience)
	if err != nil {
		r.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"error":      err.Error(),
		}).Error("Failed to build SQL query for UpdateExperience")
		return err
	}

	query = r.q.Rebind(query)
	r.log.Debug("Executing query to update experience")

	result, err := r.q.ExecContext(ctx, query, args...)
	if err != nil {
		r.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"error":      err.Error(),
		}).Error("Database error when updating experience")
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"error":      err.Error(),
		}).Error("Failed to get rows affected after update")
		return err
	}

	if rowsAffected == 0 {
		r.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"id":         experience.ID,
		}).Warn("No experience was updated")
		return fmt.Errorf("experience with ID %s not found", experience.ID)
	}

	r.log.WithFields(logrus.Fields{
		"request_id": requestID,
		"id":         experience.ID,
	}).Debug("Experience updated successfully")

	return nil
}

func (r *experienceRepository) DeleteExperience(ctx context.Context, id string) error {
	requestID := contextPkg.GetRequestID(ctx)
	r.log.WithFields(logrus.Fields{
		"request_id": requestID,
		"id":         id,
	}).Debug("Deleting experience from database")

	query := r.q.Rebind(queryDeleteExperience)
	r.log.Debug("Executing query to delete experience")

	result, err := r.q.ExecContext(ctx, query, id)
	if err != nil {
		r.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"error":      err.Error(),
			"id":         id,
		}).Error("Database error when deleting experience")
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"error":      err.Error(),
			"id":         id,
		}).Error("Failed to get rows affected after delete")
		return err
	}

	if rowsAffected == 0 {
		r.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"id":         id,
		}).Warn("No experience was deleted")
		return fmt.Errorf("experience with ID %s not found", id)
	}

	r.log.WithFields(logrus.Fields{
		"request_id": requestID,
		"id":         id,
	}).Debug("Experience deleted successfully")

	return nil
}

func (r *experienceRepository) DeleteExperiencesByUserID(ctx context.Context, userID string) error {
	requestID := contextPkg.GetRequestID(ctx)
	r.log.WithFields(logrus.Fields{
		"request_id": requestID,
		"user_id":    userID,
	}).Debug("Deleting all experiences for a user")

	query := r.q.Rebind(queryDeleteExperiencesByUserID)
	r.log.Debug("Executing query to delete experiences by user ID")

	result, err := r.q.ExecContext(ctx, query, userID)
	if err != nil {
		r.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"error":      err.Error(),
			"user_id":    userID,
		}).Error("Database error when deleting experiences by user ID")
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"error":      err.Error(),
			"user_id":    userID,
		}).Error("Failed to get rows affected after delete")
		return err
	}

	r.log.WithFields(logrus.Fields{
		"request_id": requestID,
		"user_id":    userID,
		"count":      rowsAffected,
	}).Info("Experiences deleted successfully for user")

	return nil
}

func (r *experienceRepository) makeExperience(exp bio.ExperienceDB) entity.Experience {
	return entity.Experience{
		ID:          exp.ID.String,
		UserID:      exp.UserID.String,
		ImageURL:    exp.ImageURL.String,
		JobTitle:    exp.JobTitle.String,
		SkillUsed:   exp.SkillUsed.String,
		StartDate:   exp.StartDate.String,
		EndDate:     exp.EndDate.String,
		Description: exp.Description.String,
		CreatedAt:   exp.CreatedAt.Time,
		UpdatedAt:   exp.UpdatedAt.Time,
	}
}
