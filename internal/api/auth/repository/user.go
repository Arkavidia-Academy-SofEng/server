package authRepository

import (
	"ProjectGolang/internal/entity"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"time"
)

type userDB struct {
	ID       string `db:"db"`
	Username string `db:"username"`
	Password string `db:"password"`
	Email    string `db:"email"`
}

func (r *userRepository) CreateUser(c context.Context, user entity.User) error {
	r.log.WithFields(logrus.Fields{
		"user_id":       user.ID,
		"email":         user.Email,
		"name":          user.Name,
		"role":          user.Role,
		"is_premium":    user.IsPremium,
		"premium_until": user.PremiumUntil,
		"created_at":    user.CreatedAt,
		"updated_at":    user.UpdatedAt,
	}).Debug("Creating user in database")

	query, args, err := sqlx.Named(queryCreateUser, user)
	if err != nil {
		r.log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to build SQL query for CreateUser")
		return err
	}

	query = r.q.Rebind(query)
	r.log.Debug("Executing query to create user")

	_, err = r.q.ExecContext(c, query, args...)
	if err != nil {
		r.log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Database error when creating user")
		return err
	}

	r.log.Debug("User created successfully")

	return nil
}

//func (r *userRepository) GetByID(c context.Context, id string) (entity.User, error) {
//	r.log.WithFields(map[string]interface{}{
//		"user_id": id,
//	}).Debug("Retrieving user by ID")
//
//	var user entity.User
//
//	argsKV := map[string]interface{}{
//		"id": id,
//	}
//
//	query, args, err := sqlx.Named(queryGetById, argsKV)
//	if err != nil {
//		r.log.WithFields(map[string]interface{}{
//			"error":   err.Error(),
//			"user_id": id,
//		}).Error("Failed to build SQL query for GetByID")
//		return entity.User{}, err
//	}
//
//	query = r.q.Rebind(query)
//
//	r.log.WithFields(map[string]interface{}{
//		"user_id": id,
//	}).Debug("Executing query to get user by ID")
//
//	if err := r.q.QueryRowxContext(c, query, args...).StructScan(&user); err != nil {
//		if errors.Is(err, sql.ErrNoRows) {
//			r.log.WithFields(map[string]interface{}{
//				"user_id": id,
//			}).Warn("User not found by ID")
//			return entity.User{}, auth.ErrUserNotFound
//		}
//		r.log.WithFields(map[string]interface{}{
//			"error":   err.Error(),
//			"user_id": id,
//		}).Error("Database error when retrieving user by ID")
//		return entity.User{}, err
//	}
//
//	r.log.WithFields(map[string]interface{}{
//		"user_id":  user.ID,
//		"username": user.Username,
//	}).Debug("User retrieved successfully by ID")
//
//	return user, nil
//}
//
//func (r *userRepository) UpdateUser(c context.Context, user entity.User, id string) error {
//	r.log.WithFields(map[string]interface{}{
//		"user_id":  id,
//		"username": user.Username,
//	}).Debug("Updating user in database")
//
//	argsKV := map[string]interface{}{
//		"id":       id,
//		"username": user.Username,
//		"password": user.Password,
//	}
//
//	query, args, err := sqlx.Named(queryUpdateUser, argsKV)
//	if err != nil {
//		r.log.WithFields(map[string]interface{}{
//			"error":   err.Error(),
//			"user_id": id,
//		}).Error("Failed to build SQL query for UpdateUser")
//		return err
//	}
//
//	query = r.q.Rebind(query)
//
//	r.log.WithFields(map[string]interface{}{
//		"user_id": id,
//	}).Debug("Executing query to update user")
//
//	result, err := r.q.ExecContext(c, query, args...)
//	if err != nil {
//		r.log.WithFields(map[string]interface{}{
//			"error":   err.Error(),
//			"user_id": id,
//		}).Error("Database error when updating user")
//		return err
//	}
//
//	rowsAffected, err := result.RowsAffected()
//	if err != nil {
//		r.log.WithFields(map[string]interface{}{
//			"error":   err.Error(),
//			"user_id": id,
//		}).Error("Failed to get rows affected for user update")
//		return err
//	}
//
//	r.log.WithFields(map[string]interface{}{
//		"user_id":       id,
//		"rows_affected": rowsAffected,
//	}).Debug("User update rows affected result")
//
//	if rowsAffected == 0 {
//		r.log.WithFields(map[string]interface{}{
//			"user_id": id,
//		}).Warn("No rows affected during user update, user likely not found")
//		return auth.ErrUserNotFound
//	}
//
//	r.log.WithFields(map[string]interface{}{
//		"user_id":  id,
//		"username": user.Username,
//	}).Info("User updated successfully")
//
//	return nil
//}
//
//func (r *userRepository) DeleteUser(c context.Context, id string) error {
//	r.log.WithFields(map[string]interface{}{
//		"user_id": id,
//	}).Debug("Deleting user from database")
//
//	argsKV := map[string]interface{}{
//		"id": id,
//	}
//
//	query, args, err := sqlx.Named(queryDeleteUser, argsKV)
//	if err != nil {
//		r.log.WithFields(map[string]interface{}{
//			"error":   err.Error(),
//			"user_id": id,
//		}).Error("Failed to build SQL query for DeleteUser")
//		return err
//	}
//
//	query = r.q.Rebind(query)
//
//	r.log.WithFields(map[string]interface{}{
//		"user_id": id,
//	}).Debug("Executing query to delete user")
//
//	result, err := r.q.ExecContext(c, query, args...)
//	if err != nil {
//		r.log.WithFields(map[string]interface{}{
//			"error":   err.Error(),
//			"user_id": id,
//		}).Error("Database error when deleting user")
//		if errors.Is(err, sql.ErrNoRows) {
//			return auth.ErrUserNotFound
//		}
//		return err
//	}
//
//	rowsAffected, err := result.RowsAffected()
//	if err != nil {
//		r.log.WithFields(map[string]interface{}{
//			"error":   err.Error(),
//			"user_id": id,
//		}).Error("Failed to get rows affected for user deletion")
//		return err
//	}
//
//	if rowsAffected == 0 {
//		r.log.WithFields(map[string]interface{}{
//			"user_id": id,
//		}).Warn("No rows affected during user deletion, user likely not found")
//		return auth.ErrUserNotFound
//	}
//
//	r.log.WithFields(map[string]interface{}{
//		"user_id": id,
//	}).Info("User deleted successfully")
//
//	return nil
//}
//func (r *userRepository) GetByEmail(c context.Context, email string) (entity.User, error) {
//	r.log.WithFields(logrus.Fields{
//		"email": email,
//	}).Debug("Retrieving user by email")
//
//	var user entity.User
//
//	argsKV := map[string]interface{}{
//		"email": email,
//	}
//
//	query, args, err := sqlx.Named(queryGetByEmail, argsKV)
//	if err != nil {
//		r.log.WithFields(logrus.Fields{
//			"error": err.Error(),
//			"email": email,
//		}).Error("Failed to build SQL query for GetByEmail")
//		return entity.User{}, err
//	}
//
//	query = r.q.Rebind(query)
//
//	if err := r.q.QueryRowxContext(c, query, args...).StructScan(&user); err != nil {
//		if errors.Is(err, sql.ErrNoRows) {
//			r.log.WithFields(logrus.Fields{
//				"email": email,
//			}).Warn("User not found by email")
//			return entity.User{}, auth.ErrUserNotFound
//		}
//		r.log.WithFields(logrus.Fields{
//			"error": err.Error(),
//			"email": email,
//		}).Error("Database error when retrieving user by email")
//		return entity.User{}, err
//	}
//
//	r.log.WithFields(logrus.Fields{
//		"user_id": user.ID,
//		"email":   email,
//	}).Debug("User retrieved successfully by email")
//
//	return user, nil
//}

func (r *userRepository) CheckEmailExists(c context.Context, email string) (bool, error) {
	r.log.WithFields(logrus.Fields{
		"email": email,
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
	r.log.WithFields(logrus.Fields{
		"email": email,
	}).Debug("Getting user by email")

	query := r.q.Rebind(queryGetUserByEmail)
	r.log.Debug("Executing query to get user by email")

	var user entity.User
	err := r.q.QueryRowxContext(c, query, email).Scan(
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
	r.log.WithFields(logrus.Fields{
		"id": id,
	}).Debug("Getting user by ID")

	query := r.q.Rebind(queryGetUserByID)
	r.log.Debug("Executing query to get user by ID")

	var user entity.User
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
		if err == sql.ErrNoRows {
			r.log.WithFields(logrus.Fields{
				"id": id,
			}).Warn("User not found")
			return entity.User{}, nil
		}

		r.log.WithFields(logrus.Fields{
			"error": err.Error(),
			"id":    id,
		}).Error("Database error when getting user by ID")
		return entity.User{}, err
	}

	// Check if user is soft deleted
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
