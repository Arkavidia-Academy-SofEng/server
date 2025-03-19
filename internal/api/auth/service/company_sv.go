package authService

import (
	"ProjectGolang/internal/api/auth"
	"ProjectGolang/internal/entity"
	"ProjectGolang/pkg/bcrypt"
	contextPkg "ProjectGolang/pkg/context"
	"fmt"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"mime/multipart"
	"time"
)

func (s *authService) CreateCompany(c context.Context, req auth.CreateUser) error {
	requestID := contextPkg.GetRequestID(c)

	repo, err := s.authrepository.NewClient(false)
	if err != nil {
		s.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"error":      err.Error(),
		}).Error("Failed to create repository client")
		return err
	}

	code, err := s.redis.GetOTP(c, req.Email)
	if err != nil {
		s.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"error":      err.Error(),
			"email":      req.Email,
		}).Error("Failed to get OTP from Redis")
		return err
	}

	if code != req.Code {
		s.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"email":      req.Email,
		}).Warn("Invalid OTP")
		return auth.ErrorInvalidOTP
	}

	hashedPassword, err := bcrypt.HashPassword(req.Password)
	if err != nil {
		s.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"error":      err.Error(),
		}).Error("Failed to hash password")
		return err
	}

	now := time.Now()
	id, err := NewUlidFromTimestamp(now)
	if err != nil {
		s.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"error":      err.Error(),
		}).Error("Failed to generate ULID")
		return err
	}

	newCompany := entity.Company{
		ID:          id,
		Email:       req.Email,
		Password:    hashedPassword,
		Name:        req.Name,
		PhoneNumber: req.PhoneNumber,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := repo.Company.CreateCompany(c, newCompany); err != nil {
		return err
	}

	s.log.WithFields(logrus.Fields{
		"request_id": requestID,
		"id":         newCompany.ID,
		"email":      newCompany.Email,
		"name":       newCompany.Name,
		"phone":      newCompany.PhoneNumber,
	}).Info("Company created successfully")

	return nil
}

func (s *authService) UpdateCompany(c context.Context, req auth.UpdateCompany, id string, banner *multipart.FileHeader, profile *multipart.FileHeader) error {
	requestID := contextPkg.GetRequestID(c)
	repo, err := s.authrepository.NewClient(false)
	if err != nil {
		s.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"error":      err.Error(),
		}).Error("Failed to create repository client")
		return err
	}

	existingCompany, err := repo.Company.GetCompanyByID(c, id)
	if err != nil {
		s.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"error":      err.Error(),
			"id":         id,
		}).Error("Failed to get company by ID")
		return err
	}

	if existingCompany.ID == "" {
		s.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"id":         id,
		}).Warn("Company not found")
		return auth.ErrorUserNotFound
	}

	if banner != nil {
		bannerURL, err := s.s3.UploadFile(banner, banner.Filename)
		if err != nil {
			s.log.WithFields(logrus.Fields{
				"request_id": requestID,
				"error":      err.Error(),
				"id":         id,
			}).Error("Failed to upload banner picture")
			return err
		}
		req.BannerPicture = bannerURL
	}

	if profile != nil {
		profileURL, err := s.s3.UploadFile(profile, profile.Filename)
		if err != nil {
			s.log.WithFields(logrus.Fields{
				"request_id": requestID,
				"error":      err.Error(),
				"id":         id,
			}).Error("Failed to upload profile picture")
			return err
		}
		req.ProfilePicture = profileURL
	}

	updatedCompany, err := s.updateCompanyChanges(existingCompany, req)
	if err != nil {
		s.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"error":      err.Error(),
			"id":         id,
		}).Error("Failed to update company changes")
		return err
	}

	if err := repo.Company.UpdateCompany(c, updatedCompany); err != nil {
		s.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"error":      err.Error(),
			"id":         id,
		}).Error("Failed to update company")
		return err
	}

	s.log.WithFields(logrus.Fields{
		"request_id": requestID,
		"id":         updatedCompany.ID,
		"email":      updatedCompany.Email,
		"name":       updatedCompany.Name,
	}).Info("Company updated successfully")

	return nil
}

func (s *authService) updateCompanyChanges(company entity.Company, req auth.UpdateCompany) (entity.Company, error) {
	updatedCompany := company

	updatedCompany.UpdatedAt = time.Now()

	if req.Name != "" {
		updatedCompany.Name = req.Name
	}

	if req.PhoneNumber != "" {
		updatedCompany.PhoneNumber = req.PhoneNumber
	}

	if req.ProfilePicture != "" {
		updatedCompany.ProfilePicture = req.ProfilePicture
	}

	if req.BannerPicture != "" {
		updatedCompany.BannerPicture = req.BannerPicture
	}

	if req.Location != "" {
		updatedCompany.Location = req.Location
	}

	if req.AboutUs != "" {
		updatedCompany.AboutUs = req.AboutUs
	}

	if req.IndustryTypes != "" {
		updatedCompany.IndustryTypes = req.IndustryTypes
	}

	if req.NumberEmployees != 0 {
		updatedCompany.NumberEmployees = req.NumberEmployees
	}

	if req.EstablishedDate != "" {
		establishedDate, err := time.Parse("2006-01-02", req.EstablishedDate)
		if err != nil {
			return entity.Company{}, fmt.Errorf("invalid established date format: %w", err)
		}
		updatedCompany.EstablishedDate = establishedDate
	}

	if req.CompanyURL != "" {
		updatedCompany.CompanyURL = req.CompanyURL
	}

	if req.RequiredSkill != "" {
		updatedCompany.RequiredSkill = req.RequiredSkill
	}

	return updatedCompany, nil
}

func (s *authService) DeleteCompany(c context.Context, id string) error {
	requestID := contextPkg.GetRequestID(c)
	repo, err := s.authrepository.NewClient(false)
	if err != nil {
		s.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"error":      err.Error(),
		}).Error("Failed to create repository client")
		return err
	}

	existingCompany, err := repo.Company.GetCompanyByID(c, id)
	if err != nil {
		s.log.WithFields(logrus.Fields{
			"error": err.Error(),
			"id":    id,
		}).Error("Failed to get company by ID")
		return err
	}

	if existingCompany.ID == "" {
		s.log.WithFields(logrus.Fields{
			"id": id,
		}).Warn("Company not found")
		return auth.ErrorUserNotFound
	}

	now := time.Now()
	if err := repo.Company.SoftDeleteCompany(c, id, now); err != nil {
		s.log.WithFields(logrus.Fields{
			"error": err.Error(),
			"id":    id,
		}).Error("Failed to soft delete company")
		return err
	}

	s.log.WithFields(logrus.Fields{
		"id":    id,
		"email": existingCompany.Email,
	}).Info("Company soft deleted successfully")

	return nil
}
