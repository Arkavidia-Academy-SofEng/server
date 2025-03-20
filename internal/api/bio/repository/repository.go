package bioRepository

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

	client := Client{
		Experience: &experienceRepository{q: db, log: r.log},
		Education:  &educationRepository{q: db, log: r.log},
		Portfolio:  &portfolioRepository{q: db, log: r.log},

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
	}
	return client, nil
}

type Client struct {
	Experience interface {
		CreateExperience(ctx context.Context, experience entity.Experience) error
		GetExperienceByID(ctx context.Context, id string) (entity.Experience, error)
		GetExperiencesByUserID(ctx context.Context, userID string) ([]entity.Experience, error)
		UpdateExperience(ctx context.Context, experience entity.Experience) error
		DeleteExperience(ctx context.Context, id string) error
		DeleteExperiencesByUserID(ctx context.Context, userID string) error
	}

	Education interface {
		CreateEducation(ctx context.Context, education entity.Education) error
		GetEducationByID(ctx context.Context, id string) (entity.Education, error)
		GetEducationsByUserID(ctx context.Context, userID string) ([]entity.Education, error)
		UpdateEducation(ctx context.Context, education entity.Education) error
		DeleteEducation(ctx context.Context, id string) error
		DeleteEducationsByUserID(ctx context.Context, userID string) error
	}

	Portfolio interface {
		CreatePortfolio(ctx context.Context, portfolio entity.Portfolio) error
		GetPortfolioByID(ctx context.Context, id string) (entity.Portfolio, error)
		GetPortfoliosByUserID(ctx context.Context, userID string) ([]entity.Portfolio, error)
		UpdatePortfolio(ctx context.Context, portfolio entity.Portfolio) error
		DeletePortfolio(ctx context.Context, id string) error
		DeletePortfoliosByUserID(ctx context.Context, userID string) error
	}

	Commit   func() error
	Rollback func() error
}

type experienceRepository struct {
	q   sqlx.ExtContext
	log *logrus.Logger
}

type educationRepository struct {
	q   sqlx.ExtContext
	log *logrus.Logger
}

type portfolioRepository struct {
	q   sqlx.ExtContext
	log *logrus.Logger
}
