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

func (r *portfolioRepository) CreatePortfolio(ctx context.Context, portfolio entity.Portfolio) error {
	requestID := contextPkg.GetRequestID(ctx)

	query, args, err := sqlx.Named(queryCreatePortfolio, portfolio)
	if err != nil {
		r.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"error":      err.Error(),
		}).Error("Failed to build SQL query for CreatePortfolio")
		return err
	}

	query = r.q.Rebind(query)

	_, err = r.q.ExecContext(ctx, query, args...)
	if err != nil {
		r.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"error":      err.Error(),
		}).Error("Database error when creating portfolio")
		return err
	}

	return nil
}

func (r *portfolioRepository) GetPortfolioByID(ctx context.Context, id string) (entity.Portfolio, error) {
	requestID := contextPkg.GetRequestID(ctx)
	r.log.WithFields(logrus.Fields{
		"request_id": requestID,
		"id":         id,
	}).Debug("Getting portfolio by ID")

	query := r.q.Rebind(queryGetPortfolioByID)
	r.log.Debug("Executing query to get portfolio by ID")

	var port bio.PortfolioDB
	err := r.q.QueryRowxContext(ctx, query, id).Scan(
		&port.ID,
		&port.UserID,
		&port.Image,
		&port.ProjectName,
		&port.ProjectLocation,
		&port.DescriptionImage,
		&port.ProjectLink,
		&port.StartDate,
		&port.EndDate,
		&port.Description,
		&port.CreatedAt,
		&port.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.log.WithFields(logrus.Fields{
				"request_id": requestID,
				"id":         id,
			}).Warn("Portfolio not found")
			return entity.Portfolio{}, nil
		}

		r.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"error":      err.Error(),
			"id":         id,
		}).Error("Database error when getting portfolio by ID")
		return entity.Portfolio{}, err
	}

	portfolio := r.makePortfolio(port)

	r.log.WithFields(logrus.Fields{
		"request_id": requestID,
		"id":         id,
	}).Debug("Portfolio retrieved successfully")

	return portfolio, nil
}

func (r *portfolioRepository) GetPortfoliosByUserID(ctx context.Context, userID string) ([]entity.Portfolio, error) {
	requestID := contextPkg.GetRequestID(ctx)
	r.log.WithFields(logrus.Fields{
		"request_id": requestID,
		"user_id":    userID,
	}).Debug("Getting portfolios by user ID")

	query := r.q.Rebind(queryGetPortfoliosByUserID)
	r.log.Debug("Executing query to get portfolios by user ID")

	rows, err := r.q.QueryxContext(ctx, query, userID)
	if err != nil {
		r.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"error":      err.Error(),
			"user_id":    userID,
		}).Error("Database error when getting portfolios by user ID")
		return nil, err
	}
	defer rows.Close()

	var portfolios []entity.Portfolio
	for rows.Next() {
		var port bio.PortfolioDB
		err := rows.Scan(
			&port.ID,
			&port.UserID,
			&port.Image,
			&port.ProjectName,
			&port.ProjectLocation,
			&port.DescriptionImage,
			&port.ProjectLink,
			&port.StartDate,
			&port.EndDate,
			&port.Description,
			&port.CreatedAt,
			&port.UpdatedAt,
		)
		if err != nil {
			r.log.WithFields(logrus.Fields{
				"request_id": requestID,
				"error":      err.Error(),
			}).Error("Error scanning portfolio row")
			return nil, err
		}

		portfolio := r.makePortfolio(port)
		portfolios = append(portfolios, portfolio)
	}

	if err = rows.Err(); err != nil {
		r.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"error":      err.Error(),
		}).Error("Error iterating portfolio rows")
		return nil, err
	}

	r.log.WithFields(logrus.Fields{
		"request_id": requestID,
		"user_id":    userID,
		"count":      len(portfolios),
	}).Debug("Portfolios retrieved successfully")

	return portfolios, nil
}

func (r *portfolioRepository) UpdatePortfolio(ctx context.Context, portfolio entity.Portfolio) error {
	requestID := contextPkg.GetRequestID(ctx)
	r.log.WithFields(logrus.Fields{
		"request_id":   requestID,
		"id":           portfolio.ID,
		"project_name": portfolio.ProjectName,
		"updated_at":   portfolio.UpdatedAt,
	}).Debug("Updating portfolio in database")

	query, args, err := sqlx.Named(queryUpdatePortfolio, portfolio)
	if err != nil {
		r.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"error":      err.Error(),
		}).Error("Failed to build SQL query for UpdatePortfolio")
		return err
	}

	query = r.q.Rebind(query)
	r.log.Debug("Executing query to update portfolio")

	result, err := r.q.ExecContext(ctx, query, args...)
	if err != nil {
		r.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"error":      err.Error(),
		}).Error("Database error when updating portfolio")
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
			"id":         portfolio.ID,
		}).Warn("No portfolio was updated")
		return fmt.Errorf("portfolio with ID %s not found", portfolio.ID)
	}

	r.log.WithFields(logrus.Fields{
		"request_id": requestID,
		"id":         portfolio.ID,
	}).Debug("Portfolio updated successfully")

	return nil
}

func (r *portfolioRepository) DeletePortfolio(ctx context.Context, id string) error {
	requestID := contextPkg.GetRequestID(ctx)
	r.log.WithFields(logrus.Fields{
		"request_id": requestID,
		"id":         id,
	}).Debug("Deleting portfolio from database")

	query := r.q.Rebind(queryDeletePortfolio)
	r.log.Debug("Executing query to delete portfolio")

	result, err := r.q.ExecContext(ctx, query, id)
	if err != nil {
		r.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"error":      err.Error(),
			"id":         id,
		}).Error("Database error when deleting portfolio")
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
		}).Warn("No portfolio was deleted")
		return fmt.Errorf("portfolio with ID %s not found", id)
	}

	r.log.WithFields(logrus.Fields{
		"request_id": requestID,
		"id":         id,
	}).Debug("Portfolio deleted successfully")

	return nil
}

func (r *portfolioRepository) DeletePortfoliosByUserID(ctx context.Context, userID string) error {
	requestID := contextPkg.GetRequestID(ctx)
	r.log.WithFields(logrus.Fields{
		"request_id": requestID,
		"user_id":    userID,
	}).Debug("Deleting all portfolios for a user")

	query := r.q.Rebind(queryDeletePortfoliosByUserID)
	r.log.Debug("Executing query to delete portfolios by user ID")

	result, err := r.q.ExecContext(ctx, query, userID)
	if err != nil {
		r.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"error":      err.Error(),
			"user_id":    userID,
		}).Error("Database error when deleting portfolios by user ID")
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
	}).Info("Portfolios deleted successfully for user")

	return nil
}

func (r *portfolioRepository) makePortfolio(port bio.PortfolioDB) entity.Portfolio {
	return entity.Portfolio{
		ID:               port.ID.String,
		UserID:           port.UserID.String,
		Image:            port.Image.String,
		ProjectName:      port.ProjectName.String,
		ProjectLocation:  port.ProjectLocation.String,
		DescriptionImage: port.DescriptionImage.String,
		ProjectLink:      port.ProjectLink.String,
		StartDate:        port.StartDate.String,
		EndDate:          port.EndDate.String,
		Description:      port.Description.String,
		CreatedAt:        port.CreatedAt.Time,
		UpdatedAt:        port.UpdatedAt.Time,
	}
}
