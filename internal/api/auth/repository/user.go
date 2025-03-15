package authRepository

import (
	"ProjectGolang/internal/api/auth"
	"ProjectGolang/internal/entity"
	"context"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
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
	r.log.WithFields(map[string]interface{}{
		"user_id":  user.ID,
		"email":    user.Email,
		"username": user.Username,
	}).Debug("Creating new user in database")

	argsKV := map[string]interface{}{
		"id":         user.ID,
		"email":      user.Email,
		"username":   user.Username,
		"password":   user.Password,
		"created_at": time.Now(),
	}

	r.log.Debug("Building SQL query for user creation")
	query, args, err := sqlx.Named(queryCreateUser, argsKV)
	if err != nil {
		r.log.WithFields(map[string]interface{}{
			"error":   err.Error(),
			"user_id": user.ID,
		}).Error("Failed to build SQL query for user creation")
		return err
	}
	query = r.q.Rebind(query)

	r.log.WithFields(map[string]interface{}{
		"user_id": user.ID,
	}).Debug("Executing SQL query to create user")

	_, err = r.q.ExecContext(c, query, args...)
	if err != nil {
		r.log.WithFields(map[string]interface{}{
			"error":   err.Error(),
			"user_id": user.ID,
		}).Error("Database error when creating user")

		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			switch pqErr.Code {
			case "23505":
				if pqErr.Constraint == "users_email_key" {
					r.log.WithFields(map[string]interface{}{
						"email":   user.Email,
						"user_id": user.ID,
					}).Warn("Attempted to create user with existing email")
					return auth.ErrEmailAlreadyExists
				}
			}
		}
		return err
	}

	r.log.WithFields(map[string]interface{}{
		"user_id":  user.ID,
		"username": user.Username,
	}).Info("User created successfully in database")

	return nil
}

func (r *userRepository) GetByID(c context.Context, id string) (entity.User, error) {
	r.log.WithFields(map[string]interface{}{
		"user_id": id,
	}).Debug("Retrieving user by ID")

	var user entity.User

	argsKV := map[string]interface{}{
		"id": id,
	}

	query, args, err := sqlx.Named(queryGetById, argsKV)
	if err != nil {
		r.log.WithFields(map[string]interface{}{
			"error":   err.Error(),
			"user_id": id,
		}).Error("Failed to build SQL query for GetByID")
		return entity.User{}, err
	}

	query = r.q.Rebind(query)

	r.log.WithFields(map[string]interface{}{
		"user_id": id,
	}).Debug("Executing query to get user by ID")

	if err := r.q.QueryRowxContext(c, query, args...).StructScan(&user); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.log.WithFields(map[string]interface{}{
				"user_id": id,
			}).Warn("User not found by ID")
			return entity.User{}, auth.ErrUserNotFound
		}
		r.log.WithFields(map[string]interface{}{
			"error":   err.Error(),
			"user_id": id,
		}).Error("Database error when retrieving user by ID")
		return entity.User{}, err
	}

	r.log.WithFields(map[string]interface{}{
		"user_id":  user.ID,
		"username": user.Username,
	}).Debug("User retrieved successfully by ID")

	return user, nil
}

func (r *userRepository) UpdateUser(c context.Context, user entity.User, id string) error {
	r.log.WithFields(map[string]interface{}{
		"user_id":  id,
		"username": user.Username,
	}).Debug("Updating user in database")

	argsKV := map[string]interface{}{
		"id":       id,
		"username": user.Username,
		"password": user.Password,
	}

	query, args, err := sqlx.Named(queryUpdateUser, argsKV)
	if err != nil {
		r.log.WithFields(map[string]interface{}{
			"error":   err.Error(),
			"user_id": id,
		}).Error("Failed to build SQL query for UpdateUser")
		return err
	}

	query = r.q.Rebind(query)

	r.log.WithFields(map[string]interface{}{
		"user_id": id,
	}).Debug("Executing query to update user")

	result, err := r.q.ExecContext(c, query, args...)
	if err != nil {
		r.log.WithFields(map[string]interface{}{
			"error":   err.Error(),
			"user_id": id,
		}).Error("Database error when updating user")
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.log.WithFields(map[string]interface{}{
			"error":   err.Error(),
			"user_id": id,
		}).Error("Failed to get rows affected for user update")
		return err
	}

	r.log.WithFields(map[string]interface{}{
		"user_id":       id,
		"rows_affected": rowsAffected,
	}).Debug("User update rows affected result")

	if rowsAffected == 0 {
		r.log.WithFields(map[string]interface{}{
			"user_id": id,
		}).Warn("No rows affected during user update, user likely not found")
		return auth.ErrUserNotFound
	}

	r.log.WithFields(map[string]interface{}{
		"user_id":  id,
		"username": user.Username,
	}).Info("User updated successfully")

	return nil
}

func (r *userRepository) DeleteUser(c context.Context, id string) error {
	r.log.WithFields(map[string]interface{}{
		"user_id": id,
	}).Debug("Deleting user from database")

	argsKV := map[string]interface{}{
		"id": id,
	}

	query, args, err := sqlx.Named(queryDeleteUser, argsKV)
	if err != nil {
		r.log.WithFields(map[string]interface{}{
			"error":   err.Error(),
			"user_id": id,
		}).Error("Failed to build SQL query for DeleteUser")
		return err
	}

	query = r.q.Rebind(query)

	r.log.WithFields(map[string]interface{}{
		"user_id": id,
	}).Debug("Executing query to delete user")

	result, err := r.q.ExecContext(c, query, args...)
	if err != nil {
		r.log.WithFields(map[string]interface{}{
			"error":   err.Error(),
			"user_id": id,
		}).Error("Database error when deleting user")
		if errors.Is(err, sql.ErrNoRows) {
			return auth.ErrUserNotFound
		}
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.log.WithFields(map[string]interface{}{
			"error":   err.Error(),
			"user_id": id,
		}).Error("Failed to get rows affected for user deletion")
		return err
	}

	if rowsAffected == 0 {
		r.log.WithFields(map[string]interface{}{
			"user_id": id,
		}).Warn("No rows affected during user deletion, user likely not found")
		return auth.ErrUserNotFound
	}

	r.log.WithFields(map[string]interface{}{
		"user_id": id,
	}).Info("User deleted successfully")

	return nil
}
func (r *userRepository) GetByEmail(c context.Context, email string) (entity.User, error) {
	r.log.WithFields(logrus.Fields{
		"email": email,
	}).Debug("Retrieving user by email")

	var user entity.User

	argsKV := map[string]interface{}{
		"email": email,
	}

	query, args, err := sqlx.Named(queryGetByEmail, argsKV)
	if err != nil {
		r.log.WithFields(logrus.Fields{
			"error": err.Error(),
			"email": email,
		}).Error("Failed to build SQL query for GetByEmail")
		return entity.User{}, err
	}

	query = r.q.Rebind(query)

	if err := r.q.QueryRowxContext(c, query, args...).StructScan(&user); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.log.WithFields(logrus.Fields{
				"email": email,
			}).Warn("User not found by email")
			return entity.User{}, auth.ErrUserNotFound
		}
		r.log.WithFields(logrus.Fields{
			"error": err.Error(),
			"email": email,
		}).Error("Database error when retrieving user by email")
		return entity.User{}, err
	}

	r.log.WithFields(logrus.Fields{
		"user_id": user.ID,
		"email":   email,
	}).Debug("User retrieved successfully by email")

	return user, nil
}
