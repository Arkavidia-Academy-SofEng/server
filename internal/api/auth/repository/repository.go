package authRepository

import (
	"ProjectGolang/internal/entity"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"time"
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
		User:    &userRepository{q: db, log: r.log},
		Company: &companyRepository{q: db, log: r.log},

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
	User interface {
		CreateUser(c context.Context, user entity.User) error
		GetUserByID(c context.Context, id string) (entity.User, error)
		GetUserByEmail(c context.Context, email string) (entity.User, error)
		UpdateUser(c context.Context, user entity.User) error
		CheckEmailExists(c context.Context, email string) (bool, error)
		SoftDeleteUser(c context.Context, id string, deletedAt time.Time) error
		HardDeleteExpiredUsers(c context.Context, threshold time.Time) error
	}

	Company interface {
		CreateCompany(c context.Context, company entity.Company) error
		CheckEmailExists(c context.Context, email string) (bool, error)
		GetCompanyByEmail(c context.Context, email string) (entity.Company, error)
		GetCompanyByID(c context.Context, id string) (entity.Company, error)
		UpdateCompany(c context.Context, company entity.Company) error
		SoftDeleteCompany(c context.Context, id string, deletedAt time.Time) error
		HardDeleteExpiredCompanies(c context.Context, threshold time.Time) error
	}

	Commit   func() error
	Rollback func() error
}

type userRepository struct {
	q   sqlx.ExtContext
	log *logrus.Logger
}

type companyRepository struct {
	q   sqlx.ExtContext
	log *logrus.Logger
}
