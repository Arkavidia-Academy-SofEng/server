package authRepository

import (
	"ProjectGolang/internal/api/auth"
	"ProjectGolang/internal/entity"
	contextPkg "ProjectGolang/pkg/context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"time"
)

func (r *companyRepository) CreateCompany(c context.Context, company entity.Company) error {
	requestID := contextPkg.GetRequestID(c)

	query, args, err := sqlx.Named(queryCreateCompany, company)
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
		}).Error("Database error when creating company")
		return err
	}

	return nil
}

func (r *companyRepository) CheckEmailExists(c context.Context, email string) (bool, error) {
	requestID := contextPkg.GetRequestID(c)
	r.log.WithFields(logrus.Fields{
		"request_id": requestID,
		"email":      email,
	}).Debug("Checking if email exists")

	var exists bool
	query := r.q.Rebind(queryCheckCompanyEmailExists)
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

func (r *companyRepository) GetCompanyByEmail(c context.Context, email string) (entity.Company, error) {
	requestID := contextPkg.GetRequestID(c)
	r.log.WithFields(logrus.Fields{
		"request_id": requestID,
		"email":      email,
	}).Debug("Getting company by email")

	query := r.q.Rebind(queryGetCompanyByEmail)
	r.log.Debug("Executing query to get company by email")

	var res auth.CompanyDB
	err := r.q.QueryRowxContext(c, query, email).Scan(
		&res.ID,
		&res.Email,
		&res.Password,
		&res.Name,
		&res.ProfilePicture,
		&res.BannerPicture,
		&res.AboutUs,
		&res.IndustryTypes,
		&res.NumberEmployees,
		&res.EstablishedDate,
		&res.CompanyURL,
		&res.RequiredSkill,
		&res.Location,
		&res.PhoneNumber,
		&res.CreatedAt,
		&res.UpdatedAt,
		&res.DeletedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.log.WithFields(logrus.Fields{
				"email": email,
			}).Warn("Company not found")
			return entity.Company{}, nil
		}

		r.log.WithFields(logrus.Fields{
			"error": err.Error(),
			"email": email,
		}).Error("Database error when getting company by email")
		return entity.Company{}, err
	}

	company := r.makeCompany(res)

	if company.DeletedAt != nil {
		r.log.WithFields(logrus.Fields{
			"id":         company.ID,
			"deleted_at": company.DeletedAt,
		}).Warn("Company is soft deleted")
		return entity.Company{}, nil
	}

	r.log.WithFields(logrus.Fields{
		"id":    company.ID,
		"email": company.Email,
	}).Debug("Company retrieved successfully")

	return company, nil
}

func (r *companyRepository) GetCompanyByID(c context.Context, id string) (entity.Company, error) {
	requestID := contextPkg.GetRequestID(c)
	r.log.WithFields(logrus.Fields{
		"request_id": requestID,
		"id":         id,
	}).Debug("Getting company by ID")

	query := r.q.Rebind(queryGetCompanyByID)
	r.log.Debug("Executing query to get company by ID")

	var company auth.CompanyDB
	err := r.q.QueryRowxContext(c, query, id).Scan(
		&company.ID,
		&company.ProfilePicture,
		&company.BannerPicture,
		&company.Email,
		&company.Password,
		&company.Name,
		&company.AboutUs,
		&company.IndustryTypes,
		&company.NumberEmployees,
		&company.EstablishedDate,
		&company.CompanyURL,
		&company.RequiredSkill,
		&company.Location,
		&company.PhoneNumber,
		&company.CreatedAt,
		&company.UpdatedAt,
		&company.DeletedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.log.WithFields(logrus.Fields{
				"request_id": requestID,
				"id":         id,
			}).Warn("Company not found")
			return entity.Company{}, nil
		}

		r.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"error":      err.Error(),
			"id":         id,
		}).Error("Database error when getting company by ID")
		return entity.Company{}, err
	}

	companyRes := r.makeCompany(company)

	if companyRes.DeletedAt != nil {
		r.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"id":         company.ID,
			"deleted_at": company.DeletedAt,
		}).Warn("Company is soft deleted")
		return entity.Company{}, nil
	}

	r.log.WithFields(logrus.Fields{
		"request_id": requestID,
		"id":         company.ID,
		"email":      company.Email,
	}).Debug("Company retrieved successfully")

	return companyRes, nil
}

func (r *companyRepository) UpdateCompany(c context.Context, company entity.Company) error {
	r.log.WithFields(logrus.Fields{
		"id":               company.ID,
		"name":             company.Name,
		"about_us":         company.AboutUs,
		"industry_types":   company.IndustryTypes,
		"number_employees": company.NumberEmployees,
		"established_date": company.EstablishedDate,
		"company_url":      company.CompanyURL,
		"updated_at":       company.UpdatedAt,
	}).Debug("Updating company in database")

	query, args, err := sqlx.Named(queryUpdateCompanies, company)
	if err != nil {
		r.log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to build SQL query for UpdateCompany")
		return err
	}

	query = r.q.Rebind(query)
	r.log.Debug("Executing query to update company")

	result, err := r.q.ExecContext(c, query, args...)
	if err != nil {
		r.log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Database error when updating company")
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
			"id": company.ID,
		}).Warn("No company was updated")
		return fmt.Errorf("company with ID %s not found", company.ID)
	}

	r.log.WithFields(logrus.Fields{
		"id": company.ID,
	}).Debug("Company updated successfully")

	return nil
}

func (r *companyRepository) SoftDeleteCompany(c context.Context, id string, deletedAt time.Time) error {
	r.log.WithFields(logrus.Fields{
		"id":         id,
		"deleted_at": deletedAt,
	}).Debug("Soft deleting company in database")

	query := r.q.Rebind(querySoftDeleteCompany)
	r.log.Debug("Executing query to soft delete company")

	result, err := r.q.ExecContext(c, query, deletedAt, id)
	if err != nil {
		r.log.WithFields(logrus.Fields{
			"error": err.Error(),
			"id":    id,
		}).Error("Database error when soft deleting company")
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
		}).Warn("No company was soft deleted")
		return fmt.Errorf("company with ID %s not found", id)
	}

	r.log.WithFields(logrus.Fields{
		"id": id,
	}).Debug("Company soft deleted successfully")

	return nil
}

func (r *companyRepository) HardDeleteExpiredCompanies(c context.Context, threshold time.Time) error {
	r.log.WithFields(logrus.Fields{
		"threshold": threshold,
	}).Debug("Hard deleting expired companies from database")

	query := r.q.Rebind(queryHardDeleteExpiredCompanies)
	r.log.Debug("Executing query to hard delete expired companies")

	result, err := r.q.ExecContext(c, query, threshold)
	if err != nil {
		r.log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Database error when hard deleting expired companies")
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
	}).Info("Hard deleted expired companies successfully")

	return nil
}

func (r *companyRepository) makeCompany(company auth.CompanyDB) entity.Company {
	companyRes := entity.Company{
		ID:              company.ID.String,
		Email:           company.Email.String,
		Password:        company.Password.String,
		Name:            company.Name.String,
		ProfilePicture:  company.ProfilePicture.String,
		BannerPicture:   company.BannerPicture.String,
		PhoneNumber:     company.PhoneNumber.String,
		Location:        company.Location.String,
		AboutUs:         company.AboutUs.String,
		IndustryTypes:   company.IndustryTypes.String,
		NumberEmployees: int(company.NumberEmployees.Int64),
		EstablishedDate: company.EstablishedDate.Time,
		CompanyURL:      company.CompanyURL.String,
		RequiredSkill:   company.RequiredSkill.String,
		CreatedAt:       company.CreatedAt.Time,
		UpdatedAt:       company.UpdatedAt.Time,
	}

	if company.DeletedAt.Valid {
		companyRes.DeletedAt = &company.DeletedAt.Time
	} else {
		companyRes.DeletedAt = nil
	}
	return companyRes
}
