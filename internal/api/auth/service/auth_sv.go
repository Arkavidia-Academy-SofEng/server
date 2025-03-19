package authService

import (
	"ProjectGolang/internal/api/auth"
	"ProjectGolang/internal/entity"
	"ProjectGolang/pkg/bcrypt"
	contextPkg "ProjectGolang/pkg/context"
	jwtPkg "ProjectGolang/pkg/jwt"
	"context"
	"github.com/sirupsen/logrus"
	"mime/multipart"
	"time"
)

func (s *authService) RequestOTP(c context.Context, req auth.RequestOTP) error {
	requestID := contextPkg.GetRequestID(c)
	repo, err := s.authrepository.NewClient(false)
	if err != nil {
		s.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"error":      err.Error(),
		}).Error("Failed to create repository client")
		return err
	}

	if req.Role == "candidate" {
		exists, err := repo.User.CheckEmailExists(c, req.Email)
		if err != nil {
			s.log.WithFields(logrus.Fields{
				"request_id": requestID,
				"error":      err.Error(),
				"email":      req.Email,
			}).Error("Failed to check if email exists")
			return err
		}

		if exists {
			s.log.WithFields(logrus.Fields{
				"request_id": requestID,
				"email":      req.Email,
			}).Warn("Email already exists")
			return auth.ErrorEmailAlreadyExists
		}
	} else {
		exists, err := repo.Company.CheckEmailExists(c, req.Email)
		if err != nil {
			s.log.WithFields(logrus.Fields{
				"request_id": requestID,
				"error":      err.Error(),
				"email":      req.Email,
			}).Error("Failed to check if email exists")
			return err
		}

		if exists {
			s.log.WithFields(logrus.Fields{
				"request_id": requestID,
				"email":      req.Email,
			}).Warn("Email already exists")
			return auth.ErrorEmailAlreadyExists
		}
	}

	otp, err := generateOTP(6)
	if err != nil {
		s.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"error":      err.Error(),
			"email":      req.Email,
		}).Error("Failed to generate OTP")
		return err
	}

	if err := s.redis.SetOTP(c, req.Email, otp); err != nil {
		s.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"error":      err.Error(),
			"email":      req.Email,
			"otp":        otp,
		}).Error("Failed to save OTP to Redis")
		return err
	}

	go func() {
		if err := s.smtp.CreateSmtp(req.Email, otp); err != nil {
			s.log.WithFields(logrus.Fields{
				"request_id": requestID,
				"error":      err.Error(),
				"user":       req.Email,
			}).Error("Failed to send OTP email")
		} else {
			s.log.WithFields(logrus.Fields{
				"request_id": requestID,
				"email":      req.Email,
				"otp":        "[SECRET]",
			}).Info("OTP email sent successfully")
		}
	}()

	s.log.WithFields(logrus.Fields{
		"request_id": requestID,
		"email":      req.Email,
		"otp":        "[SECRET]",
	}).Info("OTP generated and sent successfully")

	return nil
}

func (s *authService) CreateUser(c context.Context, req auth.CreateUser) error {
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

	newUser := entity.User{
		ID:          id,
		Email:       req.Email,
		Password:    hashedPassword,
		Name:        req.Name,
		Role:        req.Role,
		PhoneNumber: req.PhoneNumber,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := repo.User.CreateUser(c, newUser); err != nil {
		return err
	}

	s.log.WithFields(logrus.Fields{
		"request_id": requestID,
		"id":         newUser.ID,
		"email":      newUser.Email,
		"name":       newUser.Name,
		"role":       newUser.Role,
		"phone":      newUser.PhoneNumber,
	}).Info("User created successfully")

	return nil
}

func (s *authService) Login(c context.Context, req auth.LoginRequest) (auth.LoginResponse, error) {
	requestID := contextPkg.GetRequestID(c)
	repo, err := s.authrepository.NewClient(false)
	if err != nil {
		s.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"error":      err.Error(),
		}).Error("Failed to create repository client")
		return auth.LoginResponse{}, err
	}

	var foundComp entity.Company
	foundUser, err := repo.User.GetUserByEmail(c, req.Email)
	if err != nil {
		s.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"error":      err.Error(),
			"email":      req.Email,
		}).Error("Failed to get user by email")
		return auth.LoginResponse{}, err
	}

	if foundUser.ID == "" {
		foundComp, err = repo.Company.GetCompanyByEmail(c, req.Email)
		if err != nil {
			s.log.WithFields(logrus.Fields{
				"request_id": requestID,
				"error":      err.Error(),
				"email":      req.Email,
			}).Error("Failed to get company by email")
			return auth.LoginResponse{}, err
		}

		if foundComp.ID == "" {
			s.log.WithFields(logrus.Fields{
				"request_id": requestID,
				"email":      req.Email,
			}).Warn("User not found")
			return auth.LoginResponse{}, auth.ErrorInvalidCredentials
		}
	}

	if foundUser.ID == "" {
		foundUser.ID = foundComp.ID
		foundUser.Name = foundComp.Name
		foundUser.Role = "recruiter"
		foundUser.IsPremium = false
		foundUser.Password = foundComp.Password
	}

	err = bcrypt.ComparePassword(foundUser.Password, req.Password)
	if err != nil {
		s.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"error":      err.Error(),
			"email":      req.Email,
		}).Warn("Invalid password")
		return auth.LoginResponse{}, auth.ErrorInvalidCredentials
	}

	userData := makeUserData(foundUser)
	s.log.WithFields(logrus.Fields{
		"request_id": requestID,
		"user_id":    userData["id"],
		"username":   userData["name"],
		"role":       userData["role"],
		"is_premium": userData["is_premium"],
	}).Debug("User authenticated successfully, generating token")

	token, expired, err := jwtPkg.Sign(userData, time.Hour*1)
	if err != nil {
		s.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"error":      err.Error(),
			"user_id":    userData["id"],
		}).Error("Error signing JWT token")
		return auth.LoginResponse{}, err
	}

	loginResponse := auth.LoginResponse{
		Token:     token,
		ExpiresAt: expired,
	}

	s.log.WithFields(logrus.Fields{
		"request_id": requestID,
		"id":         foundUser.ID,
		"email":      foundUser.Email,
		"name":       foundUser.Name,
	}).Info("User logged in successfully")

	return loginResponse, nil
}

func (s *authService) UpdateUser(c context.Context, req auth.UpdateUser, id string, banner *multipart.FileHeader, profile *multipart.FileHeader) error {
	requestID := contextPkg.GetRequestID(c)
	repo, err := s.authrepository.NewClient(false)
	if err != nil {
		s.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"error":      err.Error(),
		}).Error("Failed to create repository client")
		return err
	}

	existingUser, err := repo.User.GetUserByID(c, id)
	if err != nil {
		s.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"error":      err.Error(),
			"id":         id,
		}).Error("Failed to get user by ID")
		return err
	}

	if existingUser.ID == "" {
		s.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"id":         id,
		}).Warn("User not found")
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

	updatedUser, err := s.updateUserChanges(existingUser, req)
	if err != nil {
		s.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"error":      err.Error(),
			"id":         id,
		}).Error("Failed to update user changes")
		return err
	}

	if err := repo.User.UpdateUser(c, updatedUser); err != nil {
		s.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"error":      err.Error(),
			"id":         id,
		}).Error("Failed to update user")
		return err
	}

	s.log.WithFields(logrus.Fields{
		"request_id": requestID,
		"id":         updatedUser.ID,
		"email":      updatedUser.Email,
		"name":       updatedUser.Name,
	}).Info("User updated successfully")

	return nil
}

func (s *authService) DeleteUser(c context.Context, id string) error {
	requestID := contextPkg.GetRequestID(c)
	repo, err := s.authrepository.NewClient(false)
	if err != nil {
		s.log.WithFields(logrus.Fields{
			"request_id": requestID,
			"error":      err.Error(),
		}).Error("Failed to create repository client")
		return err
	}

	existingUser, err := repo.User.GetUserByID(c, id)
	if err != nil {
		s.log.WithFields(logrus.Fields{
			"error": err.Error(),
			"id":    id,
		}).Error("Failed to get user by ID")
		return err
	}

	if existingUser.ID == "" {
		s.log.WithFields(logrus.Fields{
			"id": id,
		}).Warn("User not found")
		return auth.ErrorUserNotFound
	}

	now := time.Now()
	if err := repo.User.SoftDeleteUser(c, id, now); err != nil {
		s.log.WithFields(logrus.Fields{
			"error": err.Error(),
			"id":    id,
		}).Error("Failed to soft delete user")
		return err
	}

	s.log.WithFields(logrus.Fields{
		"id":    id,
		"email": existingUser.Email,
	}).Info("User soft deleted successfully")

	return nil
}

func (s *authService) updateUserChanges(existingUser entity.User, req auth.UpdateUser) (entity.User, error) {
	updatedUser := existingUser

	if req.Name != "" {
		updatedUser.Name = req.Name
	}
	if req.PhoneNumber != "" {
		updatedUser.PhoneNumber = req.PhoneNumber
	}

	if req.Location != "" {
		updatedUser.Location = req.Location
	}

	if req.BannerPicture != "" {
		updatedUser.BannerPicture = req.BannerPicture
	}

	if req.ProfilePicture != "" {
		updatedUser.ProfilePicture = req.ProfilePicture
	}

	if req.Headline != "" {
		updatedUser.Headline = req.Headline
	}

	updatedUser.UpdatedAt = time.Now()

	return updatedUser, nil
}
