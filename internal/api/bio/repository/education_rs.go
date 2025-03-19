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

func (r *educationRepository) CreateEducation(ctx context.Context, education entity.Education) error {
	requestID := contextPkg.GetRequestID(ctx)

	query, args, err := sqlx.Named(queryCreateEducation, education)
	if err != nil {
		r.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"error":      err.Error(),
		}).Error("Failed to build SQL query for CreateEducation")
		return err
	}

	query = r.q.Rebind(query)

	_, err = r.q.ExecContext(ctx, query, args...)
	if err != nil {
		r.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"error":      err.Error(),
		}).Error("Database error when creating education")
		return err
	}

	return nil
}

func (r *educationRepository) GetEducationByID(ctx context.Context, id string) (entity.Education, error) {
	requestID := contextPkg.GetRequestID(ctx)
	r.log.WithFields(logrus.Fields{
		"request_id": requestID,
		"id":         id,
	}).Debug("Getting education by ID")

	query := r.q.Rebind(queryGetEducationByID)
	r.log.Debug("Executing query to get education by ID")

	var edu bio.EducationDB
	err := r.q.QueryRowxContext(ctx, query, id).Scan(
		&edu.ID,
		&edu.Image,
		&edu.UserID,
		&edu.TitleDegree,
		&edu.InstitutionalName,
		&edu.StartDate,
		&edu.EndDate,
		&edu.Description,
		&edu.CreatedAt,
		&edu.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.log.WithFields(logrus.Fields{
				"request_id": requestID,
				"id":         id,
			}).Warn("Education not found")
			return entity.Education{}, nil
		}

		r.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"error":      err.Error(),
			"id":         id,
		}).Error("Database error when getting education by ID")
		return entity.Education{}, err
	}

	education := r.makeEducation(edu)

	r.log.WithFields(logrus.Fields{
		"request_id": requestID,
		"id":         id,
	}).Debug("Education retrieved successfully")

	return education, nil
}

func (r *educationRepository) GetEducationsByUserID(ctx context.Context, userID string) ([]entity.Education, error) {
	requestID := contextPkg.GetRequestID(ctx)
	r.log.WithFields(logrus.Fields{
		"request_id": requestID,
		"user_id":    userID,
	}).Debug("Getting educations by user ID")

	query := r.q.Rebind(queryGetEducationsByUserID)
	r.log.Debug("Executing query to get educations by user ID")

	rows, err := r.q.QueryxContext(ctx, query, userID)
	if err != nil {
		r.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"error":      err.Error(),
			"user_id":    userID,
		}).Error("Database error when getting educations by user ID")
		return nil, err
	}
	defer rows.Close()

	var educations []entity.Education
	for rows.Next() {
		var edu bio.EducationDB
		err := rows.Scan(
			&edu.ID,
			&edu.Image,
			&edu.UserID,
			&edu.TitleDegree,
			&edu.InstitutionalName,
			&edu.StartDate,
			&edu.EndDate,
			&edu.Description,
			&edu.CreatedAt,
			&edu.UpdatedAt,
		)
		if err != nil {
			r.log.WithFields(logrus.Fields{
				"request_id": requestID,
				"error":      err.Error(),
			}).Error("Error scanning education row")
			return nil, err
		}

		education := r.makeEducation(edu)
		educations = append(educations, education)
	}

	if err = rows.Err(); err != nil {
		r.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"error":      err.Error(),
		}).Error("Error iterating education rows")
		return nil, err
	}

	r.log.WithFields(logrus.Fields{
		"request_id": requestID,
		"user_id":    userID,
		"count":      len(educations),
	}).Debug("Educations retrieved successfully")

	return educations, nil
}

func (r *educationRepository) UpdateEducation(ctx context.Context, education entity.Education) error {
	requestID := contextPkg.GetRequestID(ctx)
	r.log.WithFields(logrus.Fields{
		"request_id":   requestID,
		"id":           education.ID,
		"title_degree": education.TitleDegree,
		"updated_at":   education.UpdatedAt,
	}).Debug("Updating education in database")

	query, args, err := sqlx.Named(queryUpdateEducation, education)
	if err != nil {
		r.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"error":      err.Error(),
		}).Error("Failed to build SQL query for UpdateEducation")
		return err
	}

	query = r.q.Rebind(query)
	r.log.Debug("Executing query to update education")

	result, err := r.q.ExecContext(ctx, query, args...)
	if err != nil {
		r.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"error":      err.Error(),
		}).Error("Database error when updating education")
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
			"id":         education.ID,
		}).Warn("No education was updated")
		return fmt.Errorf("education with ID %s not found", education.ID)
	}

	r.log.WithFields(logrus.Fields{
		"request_id": requestID,
		"id":         education.ID,
	}).Debug("Education updated successfully")

	return nil
}

func (r *educationRepository) DeleteEducation(ctx context.Context, id string) error {
	requestID := contextPkg.GetRequestID(ctx)
	r.log.WithFields(logrus.Fields{
		"request_id": requestID,
		"id":         id,
	}).Debug("Deleting education from database")

	query := r.q.Rebind(queryDeleteEducation)
	r.log.Debug("Executing query to delete education")

	result, err := r.q.ExecContext(ctx, query, id)
	if err != nil {
		r.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"error":      err.Error(),
			"id":         id,
		}).Error("Database error when deleting education")
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
		}).Warn("No education was deleted")
		return fmt.Errorf("education with ID %s not found", id)
	}

	r.log.WithFields(logrus.Fields{
		"request_id": requestID,
		"id":         id,
	}).Debug("Education deleted successfully")

	return nil
}

func (r *educationRepository) DeleteEducationsByUserID(ctx context.Context, userID string) error {
	requestID := contextPkg.GetRequestID(ctx)
	r.log.WithFields(logrus.Fields{
		"request_id": requestID,
		"user_id":    userID,
	}).Debug("Deleting all educations for a user")

	query := r.q.Rebind(queryDeleteEducationsByUserID)
	r.log.Debug("Executing query to delete educations by user ID")

	result, err := r.q.ExecContext(ctx, query, userID)
	if err != nil {
		r.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"error":      err.Error(),
			"user_id":    userID,
		}).Error("Database error when deleting educations by user ID")
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
	}).Info("Educations deleted successfully for user")

	return nil
}

func (r *educationRepository) makeEducation(edu bio.EducationDB) entity.Education {
	return entity.Education{
		ID:                edu.ID.String,
		Image:             edu.Image.String,
		TitleDegree:       edu.TitleDegree.String,
		InstitutionalName: edu.InstitutionalName.String,
		StartDate:         edu.StartDate.String,
		EndDate:           edu.EndDate.String,
		Description:       edu.Description.String,
		CreatedAt:         edu.CreatedAt.Time,
		UpdatedAt:         edu.UpdatedAt.Time,
	}
}
