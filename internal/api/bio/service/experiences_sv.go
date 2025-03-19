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

func (s *bioService) CreateExperience(ctx context.Context, req bio.CreateExperience, userID string, image *multipart.FileHeader) error {
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
			}).Error("Failed to upload experience image")
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

	newExperience := entity.Experience{
		ID:          id,
		UserID:      userID,
		ImageURL:    imageURL,
		JobTitle:    req.JobTitle,
		JobLocation: req.JobLocation,
		SkillUsed:   req.SkillUsed,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		Description: req.Description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := bioRepo.Experience.CreateExperience(ctx, newExperience); err != nil {
		s.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"error":      err.Error(),
			"user_id":    userID,
		}).Error("Failed to create experience")
		return err
	}

	s.log.WithFields(logrus.Fields{
		"request_id": requestID,
		"id":         newExperience.ID,
		"user_id":    userID,
		"job_title":  newExperience.JobTitle,
	}).Info("Experience created successfully")

	return nil
}

func (s *bioService) GetExperienceByID(ctx context.Context, id string) (entity.Experience, error) {
	requestID := contextPkg.GetRequestID(ctx)

	repo, err := s.bioRepository.NewClient(false)
	if err != nil {
		s.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"error":      err.Error(),
		}).Error("Failed to create repository client")
		return entity.Experience{}, err
	}

	experience, err := repo.Experience.GetExperienceByID(ctx, id)
	if err != nil {
		s.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"error":      err.Error(),
			"id":         id,
		}).Error("Failed to get experience by ID")
		return entity.Experience{}, err
	}

	if experience.ID == "" {
		s.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"id":         id,
		}).Warn("Experience not found")
		return entity.Experience{}, fmt.Errorf("experience not found")
	}

	s.log.WithFields(logrus.Fields{
		"request_id": requestID,
		"id":         id,
		"user_id":    experience.UserID,
	}).Debug("Experience retrieved successfully")

	return experience, nil
}

func (s *bioService) GetExperiencesByUserID(ctx context.Context, userID string) ([]entity.Experience, error) {
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

	experiences, err := bioRepo.Experience.GetExperiencesByUserID(ctx, userID)
	if err != nil {
		s.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"error":      err.Error(),
			"user_id":    userID,
		}).Error("Failed to get experiences by user ID")
		return nil, err
	}

	s.log.WithFields(logrus.Fields{
		"request_id": requestID,
		"user_id":    userID,
		"count":      len(experiences),
	}).Debug("Experiences retrieved successfully")

	return experiences, nil
}

func (s *bioService) UpdateExperience(ctx context.Context, req bio.UpdateExperience, id string, image *multipart.FileHeader) error {
	requestID := contextPkg.GetRequestID(ctx)

	repo, err := s.bioRepository.NewClient(false)
	if err != nil {
		s.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"error":      err.Error(),
		}).Error("Failed to create repository client")
		return err
	}

	existingExperience, err := repo.Experience.GetExperienceByID(ctx, id)
	if err != nil {
		s.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"error":      err.Error(),
			"id":         id,
		}).Error("Failed to get experience by ID")
		return err
	}

	if existingExperience.ID == "" {
		s.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"id":         id,
		}).Warn("Experience not found")
		return fmt.Errorf("experience not found")
	}

	if image != nil {
		err := s.s3.DeleteFile(existingExperience.ImageURL)
		if err != nil {
			s.log.WithFields(logrus.Fields{
				"request_id": requestID,
				"error":      err.Error(),
				"id":         id,
			}).Error("Failed to delete existing experience image")

		}

		imageURL, err := s.s3.UploadFile(image, image.Filename)
		if err != nil {
			s.log.WithFields(logrus.Fields{
				"request_id": requestID,
				"error":      err.Error(),
				"id":         id,
			}).Error("Failed to upload experience image")
			return err
		}
		req.ImageURL = imageURL
	}

	updatedExperience := s.updateExperienceChanges(existingExperience, req)

	if err := repo.Experience.UpdateExperience(ctx, updatedExperience); err != nil {
		s.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"error":      err.Error(),
			"id":         id,
		}).Error("Failed to update experience")
		return err
	}

	s.log.WithFields(logrus.Fields{
		"request_id": requestID,
		"id":         updatedExperience.ID,
		"user_id":    updatedExperience.UserID,
		"job_title":  updatedExperience.JobTitle,
	}).Info("Experience updated successfully")

	return nil
}

func (s *bioService) updateExperienceChanges(experience entity.Experience, req bio.UpdateExperience) entity.Experience {
	updatedExperience := experience
	updatedExperience.UpdatedAt = time.Now()

	if req.JobTitle != "" {
		updatedExperience.JobTitle = req.JobTitle
	}

	if req.SkillUsed != "" {
		updatedExperience.SkillUsed = req.SkillUsed
	}

	if req.StartDate != "" {
		updatedExperience.StartDate = req.StartDate
	}

	if req.EndDate != "" {
		updatedExperience.EndDate = req.EndDate
	}

	if req.Description != "" {
		updatedExperience.Description = req.Description
	}

	if req.ImageURL != "" {
		updatedExperience.ImageURL = req.ImageURL
	}

	return updatedExperience
}

func (s *bioService) DeleteExperience(ctx context.Context, id string) error {
	requestID := contextPkg.GetRequestID(ctx)

	repo, err := s.bioRepository.NewClient(false)
	if err != nil {
		s.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"error":      err.Error(),
		}).Error("Failed to create repository client")
		return err
	}

	existingExperience, err := repo.Experience.GetExperienceByID(ctx, id)
	if err != nil {
		s.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"error":      err.Error(),
			"id":         id,
		}).Error("Failed to get experience by ID")
		return err
	}

	if existingExperience.ID == "" {
		s.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"id":         id,
		}).Warn("Experience not found")
		return fmt.Errorf("experience not found")
	}

	if err := repo.Experience.DeleteExperience(ctx, id); err != nil {
		s.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"error":      err.Error(),
			"id":         id,
		}).Error("Failed to delete experience")
		return err
	}

	s.log.WithFields(logrus.Fields{
		"request_id": requestID,
		"id":         id,
		"user_id":    existingExperience.UserID,
	}).Info("Experience deleted successfully")

	return nil
}
