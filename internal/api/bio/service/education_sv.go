package bioService

import (
	"ProjectGolang/internal/api/bio"
	"ProjectGolang/internal/entity"
	contextPkg "ProjectGolang/pkg/context"
	"ProjectGolang/pkg/utils"
	"fmt"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"mime/multipart"
	"time"
)

func (s *bioService) CreateEducation(ctx context.Context, req bio.CreateEducation, userID string, image *multipart.FileHeader) error {
	requestID := contextPkg.GetRequestID(ctx)

	bioRepo, err := s.bioRepository.NewClient(false)
	if err != nil {
		s.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"error":      err.Error(),
		}).Error("Failed to create repository client")
		return err
	}

	authRepo, err := s.authRepository.NewClient(false)
	if err != nil {
		s.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"error":      err.Error(),
		}).Error("Failed to create repository")
		return err
	}

	user, err := authRepo.User.GetUserByID(ctx, userID)
	if err != nil {
		s.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"error":      err.Error(),
			"user_id":    userID,
		}).Error("Failed to get user by ID")
		return err
	}

	if user.ID == "" {
		s.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"user_id":    userID,
		}).Warn("User not found")
		return fmt.Errorf("user not found")
	}

	var imageURL string
	if image != nil {
		uploadedURL, err := s.s3.UploadFile(image, image.Filename)
		if err != nil {
			s.log.WithFields(logrus.Fields{
				"request_id": requestID,
				"error":      err.Error(),
				"user_id":    userID,
			}).Error("Failed to upload education image")
			return err
		}
		imageURL = uploadedURL
	}

	now := time.Now()
	id, err := utils.NewUlidFromTimestamp(now)
	if err != nil {
		s.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"error":      err.Error(),
		}).Error("Failed to generate ULID")
		return err
	}

	newEducation := entity.Education{
		ID:                id,
		Image:             imageURL,
		UserID:            userID,
		TitleDegree:       req.TitleDegree,
		InstitutionalName: req.InstitutionalName,
		StartDate:         req.StartDate,
		EndDate:           req.EndDate,
		Description:       req.Description,
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	if err := bioRepo.Education.CreateEducation(ctx, newEducation); err != nil {
		s.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"error":      err.Error(),
			"user_id":    userID,
		}).Error("Failed to create education")
		return err
	}

	s.log.WithFields(logrus.Fields{
		"request_id":         requestID,
		"id":                 newEducation.ID,
		"user_id":            userID,
		"title_degree":       newEducation.TitleDegree,
		"institutional_name": newEducation.InstitutionalName,
	}).Info("Education created successfully")

	return nil
}

func (s *bioService) GetEducationByID(ctx context.Context, id string) (entity.Education, error) {
	requestID := contextPkg.GetRequestID(ctx)

	repo, err := s.bioRepository.NewClient(false)
	if err != nil {
		s.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"error":      err.Error(),
		}).Error("Failed to create repository client")
		return entity.Education{}, err
	}

	education, err := repo.Education.GetEducationByID(ctx, id)
	if err != nil {
		s.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"error":      err.Error(),
			"id":         id,
		}).Error("Failed to get education by ID")
		return entity.Education{}, err
	}

	if education.ID == "" {
		s.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"id":         id,
		}).Warn("Education not found")
		return entity.Education{}, fmt.Errorf("education not found")
	}

	s.log.WithFields(logrus.Fields{
		"request_id": requestID,
		"id":         id,
	}).Debug("Education retrieved successfully")

	return education, nil
}

func (s *bioService) GetEducationsByUserID(ctx context.Context, userID string) ([]entity.Education, error) {
	requestID := contextPkg.GetRequestID(ctx)

	bioRepo, err := s.bioRepository.NewClient(false)
	if err != nil {
		s.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"error":      err.Error(),
		}).Error("Failed to create repository client")
		return nil, err
	}

	authRepo, err := s.authRepository.NewClient(false)
	if err != nil {
		s.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"error":      err.Error(),
		}).Error("Failed to create repository")
		return nil, err
	}

	// Check if user exists
	user, err := authRepo.User.GetUserByID(ctx, userID)
	if err != nil {
		s.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"error":      err.Error(),
			"user_id":    userID,
		}).Error("Failed to get user by ID")
		return nil, err
	}

	if user.ID == "" {
		s.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"user_id":    userID,
		}).Warn("User not found")
		return nil, fmt.Errorf("user not found")
	}

	educations, err := bioRepo.Education.GetEducationsByUserID(ctx, userID)
	if err != nil {
		s.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"error":      err.Error(),
			"user_id":    userID,
		}).Error("Failed to get educations by user ID")
		return nil, err
	}

	s.log.WithFields(logrus.Fields{
		"request_id": requestID,
		"user_id":    userID,
		"count":      len(educations),
	}).Debug("Educations retrieved successfully")

	return educations, nil
}

func (s *bioService) UpdateEducation(ctx context.Context, req bio.UpdateEducation, id string, image *multipart.FileHeader) error {
	requestID := contextPkg.GetRequestID(ctx)

	repo, err := s.bioRepository.NewClient(false)
	if err != nil {
		s.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"error":      err.Error(),
		}).Error("Failed to create repository client")
		return err
	}

	existingEducation, err := repo.Education.GetEducationByID(ctx, id)
	if err != nil {
		s.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"error":      err.Error(),
			"id":         id,
		}).Error("Failed to get education by ID")
		return err
	}

	if existingEducation.ID == "" {
		s.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"id":         id,
		}).Warn("Education not found")
		return fmt.Errorf("education not found")
	}

	if image != nil {
		err := s.s3.DeleteFile(existingEducation.Image)
		if err != nil {
			s.log.WithFields(logrus.Fields{
				"request_id": requestID,
				"error":      err.Error(),
				"id":         id,
			}).Error("Failed to delete existing education image")
		}

		imageURL, err := s.s3.UploadFile(image, image.Filename)
		if err != nil {
			s.log.WithFields(logrus.Fields{
				"request_id": requestID,
				"error":      err.Error(),
				"id":         id,
			}).Error("Failed to upload education image")
			return err
		}
		req.Image = imageURL
	}

	updatedEducation := s.updateEducationChanges(existingEducation, req)

	if err := repo.Education.UpdateEducation(ctx, updatedEducation); err != nil {
		s.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"error":      err.Error(),
			"id":         id,
		}).Error("Failed to update education")
		return err
	}

	s.log.WithFields(logrus.Fields{
		"request_id":         requestID,
		"id":                 updatedEducation.ID,
		"title_degree":       updatedEducation.TitleDegree,
		"institutional_name": updatedEducation.InstitutionalName,
	}).Info("Education updated successfully")

	return nil
}

func (s *bioService) updateEducationChanges(education entity.Education, req bio.UpdateEducation) entity.Education {
	updatedEducation := education
	updatedEducation.UpdatedAt = time.Now()

	if req.TitleDegree != "" {
		updatedEducation.TitleDegree = req.TitleDegree
	}

	if req.InstitutionalName != "" {
		updatedEducation.InstitutionalName = req.InstitutionalName
	}

	if req.StartDate != "" {
		updatedEducation.StartDate = req.StartDate
	}

	if req.EndDate != "" {
		updatedEducation.EndDate = req.EndDate
	}

	if req.Description != "" {
		updatedEducation.Description = req.Description
	}

	if req.Image != "" {
		updatedEducation.Image = req.Image
	}

	return updatedEducation
}

func (s *bioService) DeleteEducation(ctx context.Context, id string) error {
	requestID := contextPkg.GetRequestID(ctx)

	repo, err := s.bioRepository.NewClient(false)
	if err != nil {
		s.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"error":      err.Error(),
		}).Error("Failed to create repository client")
		return err
	}

	existingEducation, err := repo.Education.GetEducationByID(ctx, id)
	if err != nil {
		s.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"error":      err.Error(),
			"id":         id,
		}).Error("Failed to get education by ID")
		return err
	}

	if existingEducation.ID == "" {
		s.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"id":         id,
		}).Warn("Education not found")
		return fmt.Errorf("education not found")
	}

	if err := repo.Education.DeleteEducation(ctx, id); err != nil {
		s.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"error":      err.Error(),
			"id":         id,
		}).Error("Failed to delete education")
		return err
	}

	s.log.WithFields(logrus.Fields{
		"request_id": requestID,
		"id":         id,
	}).Info("Education deleted successfully")

	return nil
}
