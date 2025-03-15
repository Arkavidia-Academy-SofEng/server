package authRepository

import (
	"ProjectGolang/internal/entity"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

func New(db *sqlx.DB, log *logrus.Logger) Repository {
	return &repository{
		DB:  db,
		log: log,
	}
}

type repository struct {
	DB  *sqlx.DB
	log *logrus.Logger
}

type Repository interface {
	NewClient(tx bool) (Client, error)
}

func (r *repository) NewClient(tx bool) (Client, error) {
	var db sqlx.ExtContext
	var commitFunc, rollbackFunc func() error

	r.log.WithFields(logrus.Fields{
		"transaction": tx,
	}).Debug("Creating new repository client")

	db = r.DB

	if tx {
		r.log.Debug("Starting database transaction")
		var err error
		txx, err := r.DB.Beginx()
		if err != nil {
			r.log.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Error("Failed to begin transaction")
			return Client{}, err
		}

		db = txx
		commitFunc = txx.Commit
		rollbackFunc = txx.Rollback
	} else {
		commitFunc = func() error { return nil }
		rollbackFunc = func() error { return nil }
	}

	return Client{
		Users: &userRepository{q: db, log: r.log},
		Commit: func() error {
			if tx {
				r.log.Debug("Committing transaction")
			}
			return commitFunc()
		},
		Rollback: func() error {
			if tx {
				r.log.Debug("Rolling back transaction")
			}
			return rollbackFunc()
		},
	}, nil
}

type Client struct {
	Users interface {
		CreateUser(ctx context.Context, user entity.User) error
		GetByID(ctx context.Context, id string) (entity.User, error)
		GetByEmail(ctx context.Context, email string) (entity.User, error)
		UpdateUser(ctx context.Context, user entity.User, id string) error
		DeleteUser(ctx context.Context, id string) error
	}

	Commit   func() error
	Rollback func() error
}

type userRepository struct {
	q   sqlx.ExtContext
	log *logrus.Logger
}

type sessionRepository struct {
	q   sqlx.ExtContext
	log *logrus.Logger
}

type userOauthRepository struct {
	q   sqlx.ExtContext
	log *logrus.Logger
}
