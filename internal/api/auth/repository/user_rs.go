package authRepository

import (
	"ProjectGolang/internal/api/auth"
	"ProjectGolang/internal/entity"
	contextPkg "ProjectGolang/pkg/context"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"time"
)

func (r *userRepository) CreateUser(c context.Context, user entity.User) error {
	requestID := contextPkg.GetRequestID(c)

	query, args, err := sqlx.Named(queryCreateUser, user)
	if err != nil {
		r.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"error":      err.Error(),
		}).Error("Failed to build SQL query for CreateUser")
		return err
	}

	query = r.q.Rebind(query)

	_, err = r.q.ExecContext(c, query, args...)
	if err != nil {
		r.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"error":      err.Error(),
		}).Error("Database error when creating user")
		return err
	}

	return nil
}

func (r *userRepository) CheckEmailExists(c context.Context, email string) (bool, error) {
	requestID := contextPkg.GetRequestID(c)
	r.log.WithFields(logrus.Fields{
		"request_id": requestID,
		"email":      email,
	}).Debug("Checking if email exists")

	var exists bool
	query := r.q.Rebind(queryCheckEmailExists)
	err := r.q.QueryRowxContext(c, query, email).Scan(&exists)
	if err != nil {
		r.log.WithFields(logrus.Fields{
			"error": err.Error(),
			"email": email,
		}).Error("Database error when checking email existence")
		return false, err
	}

	return exists, nil
}

func (r *userRepository) GetUserByEmail(c context.Context, email string) (entity.User, error) {
	requestID := contextPkg.GetRequestID(c)
	r.log.WithFields(logrus.Fields{
		"request_id": requestID,
		"email":      email,
	}).Debug("Getting user by email")

	query := r.q.Rebind(queryGetUserByEmail)
	r.log.Debug("Executing query to get user by email")

	var res auth.UserDB
	err := r.q.QueryRowxContext(c, query, email).Scan(
		&res.ID,
		&res.Email,
		&res.Password,
		&res.Name,
		&res.PhoneNumber,
		&res.Role,
		&res.ProfilePicture,
		&res.IsPremium,
		&res.PremiumUntil,
		&res.Headline,
		&res.CreatedAt,
		&res.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.log.WithFields(logrus.Fields{
				"email": email,
			}).Warn("User not found")
			return entity.User{}, nil
		}

		r.log.WithFields(logrus.Fields{
			"error": err.Error(),
			"email": email,
		}).Error("Database error when getting user by email")
		return entity.User{}, err
	}

	user := r.makeUser(res)

	if user.DeletedAt != nil {
		r.log.WithFields(logrus.Fields{
			"id":         user.ID,
			"deleted_at": user.DeletedAt,
		}).Warn("User is soft deleted")
		return entity.User{}, nil
	}

	r.log.WithFields(logrus.Fields{
		"id":    user.ID,
		"email": user.Email,
	}).Debug("User retrieved successfully")

	return user, nil
}

func (r *userRepository) GetUserByID(c context.Context, id string) (entity.User, error) {
	requestID := contextPkg.GetRequestID(c)
	r.log.WithFields(logrus.Fields{
		"request_id": requestID,
		"id":         id,
	}).Debug("Getting user by ID")

	query := r.q.Rebind(queryGetUserByID)
	r.log.Debug("Executing query to get user by ID")

	var user auth.UserDB
	err := r.q.QueryRowxContext(c, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Name,
		&user.Role,
		&user.ProfilePicture,
		&user.IsPremium,
		&user.PremiumUntil,
		&user.Headline,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.log.WithFields(logrus.Fields{
				"request_id": requestID,
				"id":         id,
			}).Warn("User not found")
			return entity.User{}, nil
		}

		r.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"error":      err.Error(),
			"id":         id,
		}).Error("Database error when getting user by ID")
		return entity.User{}, err
	}

	userRes := r.makeUser(user)

	if userRes.DeletedAt != nil {
		r.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"id":         user.ID,
			"deleted_at": user.DeletedAt,
		}).Warn("User is soft deleted")
		return entity.User{}, nil
	}

	r.log.WithFields(logrus.Fields{
		"request_id": requestID,
		"id":         user.ID,
		"email":      user.Email,
	}).Debug("User retrieved successfully")

	return userRes, nil
}

func (r *userRepository) UpdateUser(c context.Context, user entity.User) error {
	r.log.WithFields(logrus.Fields{
		"id":            user.ID,
		"name":          user.Name,
		"role":          user.Role,
		"is_premium":    user.IsPremium,
		"premium_until": user.PremiumUntil,
		"updated_at":    user.UpdatedAt,
	}).Debug("Updating user in database")

	query, args, err := sqlx.Named(queryUpdateUser, user)
	if err != nil {
		r.log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to build SQL query for UpdateUser")
		return err
	}

	query = r.q.Rebind(query)
	r.log.Debug("Executing query to update user")

	result, err := r.q.ExecContext(c, query, args...)
	if err != nil {
		r.log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Database error when updating user")
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to get rows affected after update")
		return err
	}

	if rowsAffected == 0 {
		r.log.WithFields(logrus.Fields{
			"id": user.ID,
		}).Warn("No user was updated")
		return fmt.Errorf("user with ID %s not found", user.ID)
	}

	r.log.WithFields(logrus.Fields{
		"id": user.ID,
	}).Debug("User updated successfully")

	return nil
}

func (r *userRepository) SoftDeleteUser(c context.Context, id string, deletedAt time.Time) error {
	r.log.WithFields(logrus.Fields{
		"id":         id,
		"deleted_at": deletedAt,
	}).Debug("Soft deleting user in database")

	query := r.q.Rebind(querySoftDeleteUser)
	r.log.Debug("Executing query to soft delete user")

	result, err := r.q.ExecContext(c, query, deletedAt, id)
	if err != nil {
		r.log.WithFields(logrus.Fields{
			"error": err.Error(),
			"id":    id,
		}).Error("Database error when soft deleting user")
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.log.WithFields(logrus.Fields{
			"error": err.Error(),
			"id":    id,
		}).Error("Failed to get rows affected after soft delete")
		return err
	}

	if rowsAffected == 0 {
		r.log.WithFields(logrus.Fields{
			"id": id,
		}).Warn("No user was soft deleted")
		return fmt.Errorf("user with ID %s not found", id)
	}

	r.log.WithFields(logrus.Fields{
		"id": id,
	}).Debug("User soft deleted successfully")

	return nil
}

func (r *userRepository) HardDeleteExpiredUsers(c context.Context, threshold time.Time) error {
	r.log.WithFields(logrus.Fields{
		"threshold": threshold,
	}).Debug("Hard deleting expired users from database")

	query := r.q.Rebind(queryHardDeleteExpiredUsers)
	r.log.Debug("Executing query to hard delete expired users")

	result, err := r.q.ExecContext(c, query, threshold)
	if err != nil {
		r.log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Database error when hard deleting expired users")
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to get rows affected after hard delete")
		return err
	}

	r.log.WithFields(logrus.Fields{
		"count": rowsAffected,
	}).Info("Hard deleted expired users successfully")

	return nil
}

func (r *userRepository) makeUser(user auth.UserDB) entity.User {
	userRes := entity.User{
		ID:             user.ID.String,
		Email:          user.Email.String,
		Password:       user.Password.String,
		Name:           user.Name.String,
		Role:           entity.UserRole(user.Role.String),
		ProfilePicture: user.ProfilePicture.String,
		IsPremium:      user.IsPremium.Bool,
		PremiumUntil:   user.PremiumUntil.Time,
		Headline:       user.Headline.String,
		CreatedAt:      user.CreatedAt.Time,
		UpdatedAt:      user.UpdatedAt.Time,
		Location:       user.Location.String,
	}

	if user.DeletedAt.Valid {
		userRes.DeletedAt = &user.DeletedAt.Time
	} else {
		userRes.DeletedAt = nil
	}
	return userRes
}
