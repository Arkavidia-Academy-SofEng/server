package recruitmentRepository

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
		JobVacancies: &jobVacanciesRepository{q: db, log: r.log},
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
	JobVacancies interface {
		CreateJobVacancy(c context.Context, jobVacancy entity.JobVacancy) error
		GetJobVacancies(c context.Context, page, pageSize int) ([]entity.JobVacancy, int, error)
		CheckJobVacancyExists(c context.Context, id string) (bool, error)
		UpdateJobVacancy(c context.Context, jobVacancy entity.JobVacancy) error
		DeleteJobVacancy(c context.Context, id string) error
	}

	Commit   func() error
	Rollback func() error
}

type jobVacanciesRepository struct {
	q   sqlx.ExtContext
	log *logrus.Logger
}
