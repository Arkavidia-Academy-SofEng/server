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

func (s *bioService) CreatePortfolio(ctx context.Context, req bio.CreatePortfolio, userID string, image *multipart.FileHeader, descriptionImage *multipart.FileHeader) error {
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
			}).Error("Failed to upload portfolio image")
			return err
		}
		imageURL = uploadedURL
	}

	var descriptionImageURL string
	if descriptionImage != nil {
		uploadedURL, err := s.s3.UploadFile(descriptionImage, descriptionImage.Filename)
		if err != nil {
			s.log.WithFields(logrus.Fields{
				"request_id": requestID,
				"error":      err.Error(),
				"user_id":    userID,
			}).Error("Failed to upload portfolio image")
			return err
		}
		descriptionImageURL = uploadedURL
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

	newPortfolio := entity.Portfolio{
		ID:               id,
		UserID:           userID,
		Image:            imageURL,
		ProjectName:      req.ProjectName,
		ProjectLocation:  req.ProjectLocation,
		DescriptionImage: descriptionImageURL,
		ProjectLink:      req.ProjectLink,
		StartDate:        req.StartDate,
		EndDate:          req.EndDate,
		Description:      req.Description,
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	if err := bioRepo.Portfolio.CreatePortfolio(ctx, newPortfolio); err != nil {
		s.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"error":      err.Error(),
			"user_id":    userID,
		}).Error("Failed to create portfolio")
		return err
	}

	s.log.WithFields(logrus.Fields{
		"request_id":       requestID,
		"id":               newPortfolio.ID,
		"user_id":          userID,
		"project_name":     newPortfolio.ProjectName,
		"project_location": newPortfolio.ProjectLocation,
	}).Info("Portfolio created successfully")

	return nil
}

func (s *bioService) GetPortfolioByID(ctx context.Context, id string) (entity.Portfolio, error) {
	requestID := contextPkg.GetRequestID(ctx)

	repo, err := s.bioRepository.NewClient(false)
	if err != nil {
		s.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"error":      err.Error(),
		}).Error("Failed to create repository client")
		return entity.Portfolio{}, err
	}

	portfolio, err := repo.Portfolio.GetPortfolioByID(ctx, id)
	if err != nil {
		s.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"error":      err.Error(),
			"id":         id,
		}).Error("Failed to get portfolio by ID")
		return entity.Portfolio{}, err
	}

	if portfolio.ID == "" {
		s.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"id":         id,
		}).Warn("Portfolio not found")
		return entity.Portfolio{}, fmt.Errorf("portfolio not found")
	}

	s.log.WithFields(logrus.Fields{
		"request_id": requestID,
		"id":         id,
	}).Debug("Portfolio retrieved successfully")

	return portfolio, nil
}

func (s *bioService) GetPortfoliosByUserID(ctx context.Context, userID string) ([]entity.Portfolio, error) {
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

	portfolios, err := bioRepo.Portfolio.GetPortfoliosByUserID(ctx, userID)
	if err != nil {
		s.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"error":      err.Error(),
			"user_id":    userID,
		}).Error("Failed to get portfolios by user ID")
		return nil, err
	}

	s.log.WithFields(logrus.Fields{
		"request_id": requestID,
		"user_id":    userID,
		"count":      len(portfolios),
	}).Debug("Portfolios retrieved successfully")

	return portfolios, nil
}

func (s *bioService) UpdatePortfolio(ctx context.Context, req bio.UpdatePortfolio, id string, image *multipart.FileHeader, descriptionImage *multipart.FileHeader) error {
	requestID := contextPkg.GetRequestID(ctx)

	repo, err := s.bioRepository.NewClient(false)
	if err != nil {
		s.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"error":      err.Error(),
		}).Error("Failed to create repository client")
		return err
	}

	existingPortfolio, err := repo.Portfolio.GetPortfolioByID(ctx, id)
	if err != nil {
		s.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"error":      err.Error(),
			"id":         id,
		}).Error("Failed to get portfolio by ID")
		return err
	}

	if existingPortfolio.ID == "" {
		s.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"id":         id,
		}).Warn("Portfolio not found")
		return fmt.Errorf("portfolio not found")
	}

	if image != nil {
		err := s.s3.DeleteFile(existingPortfolio.Image)
		if err != nil {
			s.log.WithFields(logrus.Fields{
				"request_id": requestID,
				"error":      err.Error(),
				"id":         id,
			}).Error("Failed to delete existing portfolio image")
		}

		imageURL, err := s.s3.UploadFile(image, image.Filename)
		if err != nil {
			s.log.WithFields(logrus.Fields{
				"request_id": requestID,
				"error":      err.Error(),
				"id":         id,
			}).Error("Failed to upload portfolio image")
			return err
		}
		req.Image = imageURL
	}

	if descriptionImage != nil {
		err := s.s3.DeleteFile(existingPortfolio.Image)
		if err != nil {
			s.log.WithFields(logrus.Fields{
				"request_id": requestID,
				"error":      err.Error(),
				"id":         id,
			}).Error("Failed to delete existing description image")
		}

		descriptionImageURL, err := s.s3.UploadFile(descriptionImage, descriptionImage.Filename)
		if err != nil {
			s.log.WithFields(logrus.Fields{
				"request_id": requestID,
				"error":      err.Error(),
				"id":         id,
			}).Error("Failed to upload description image")
			return err
		}
		req.DescriptionImage = descriptionImageURL
	}

	updatedPortfolio := s.updatePortfolioChanges(existingPortfolio, req)

	if err := repo.Portfolio.UpdatePortfolio(ctx, updatedPortfolio); err != nil {
		s.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"error":      err.Error(),
			"id":         id,
		}).Error("Failed to update portfolio")
		return err
	}

	s.log.WithFields(logrus.Fields{
		"request_id":       requestID,
		"id":               updatedPortfolio.ID,
		"project_name":     updatedPortfolio.ProjectName,
		"project_location": updatedPortfolio.ProjectLocation,
	}).Info("Portfolio updated successfully")

	return nil
}

func (s *bioService) updatePortfolioChanges(portfolio entity.Portfolio, req bio.UpdatePortfolio) entity.Portfolio {
	updatedPortfolio := portfolio
	updatedPortfolio.UpdatedAt = time.Now()

	if req.ProjectName != "" {
		updatedPortfolio.ProjectName = req.ProjectName
	}

	if req.ProjectLocation != "" {
		updatedPortfolio.ProjectLocation = req.ProjectLocation
	}

	if req.DescriptionImage != "" {
		updatedPortfolio.DescriptionImage = req.DescriptionImage
	}

	if req.ProjectLink != "" {
		updatedPortfolio.ProjectLink = req.ProjectLink
	}

	if req.StartDate != "" {
		updatedPortfolio.StartDate = req.StartDate
	}

	if req.EndDate != "" {
		updatedPortfolio.EndDate = req.EndDate
	}

	if req.Description != "" {
		updatedPortfolio.Description = req.Description
	}

	if req.Image != "" {
		updatedPortfolio.Image = req.Image
	}

	return updatedPortfolio
}

func (s *bioService) DeletePortfolio(ctx context.Context, id string) error {
	requestID := contextPkg.GetRequestID(ctx)

	repo, err := s.bioRepository.NewClient(false)
	if err != nil {
		s.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"error":      err.Error(),
		}).Error("Failed to create repository client")
		return err
	}

	existingPortfolio, err := repo.Portfolio.GetPortfolioByID(ctx, id)
	if err != nil {
		s.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"error":      err.Error(),
			"id":         id,
		}).Error("Failed to get portfolio by ID")
		return err
	}

	if existingPortfolio.ID == "" {
		s.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"id":         id,
		}).Warn("Portfolio not found")
		return fmt.Errorf("portfolio not found")
	}

	if existingPortfolio.Image != "" {
		err := s.s3.DeleteFile(existingPortfolio.Image)
		if err != nil {
			s.log.WithFields(logrus.Fields{
				"request_id": requestID,
				"error":      err.Error(),
				"id":         id,
			}).Error("Failed to delete portfolio image from S3")
		}
	}

	if err := repo.Portfolio.DeletePortfolio(ctx, id); err != nil {
		s.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"error":      err.Error(),
			"id":         id,
		}).Error("Failed to delete portfolio")
		return err
	}

	s.log.WithFields(logrus.Fields{
		"request_id": requestID,
		"id":         id,
	}).Info("Portfolio deleted successfully")

	return nil
}
