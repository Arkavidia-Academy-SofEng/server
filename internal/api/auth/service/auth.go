package authService

import (
	"ProjectGolang/internal/api/auth"
	"ProjectGolang/internal/entity"
	"ProjectGolang/pkg/bcrypt"
	jwtPkg "ProjectGolang/pkg/jwt"
	"context"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"time"
)

func (s *authService) CreateUser(c context.Context, req auth.CreateUser) (auth.UserResponse, error) {
	repo, err := s.authrepository.NewClient(false)
	if err != nil {
		s.log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to create repository client")
		return auth.UserResponse{}, err
	}

	// Check if email already exists
	exists, err := repo.Users.CheckEmailExists(c, req.Email)
	if err != nil {
		s.log.WithFields(logrus.Fields{
			"error": err.Error(),
			"email": req.Email,
		}).Error("Failed to check if email exists")
		return auth.UserResponse{}, err
	}

	if exists {
		s.log.WithFields(logrus.Fields{
			"email": req.Email,
		}).Warn("Email already exists")
		return auth.UserResponse{}, auth.ErrorEmailAlreadyExists
	}

	// Hash the password
	hashedPassword, err := bcrypt.HashPassword(req.Password)
	if err != nil {
		s.log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to hash password")
		return auth.UserResponse{}, err
	}

	now := time.Now()
	newUser := entity.User{
		ID:             uuid.New().String(),
		Email:          req.Email,
		Password:       hashedPassword,
		Name:           req.Name,
		Role:           req.Role,
		ProfilePicture: req.ProfilePicture,
		IsPremium:      req.IsPremium,
		PremiumUntil:   req.PremiumUntil,
		Headline:       req.Headline,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if err := repo.Users.CreateUser(c, newUser); err != nil {
		s.log.WithFields(logrus.Fields{
			"error": err.Error(),
			"email": req.Email,
		}).Error("Failed to create user")
		return auth.UserResponse{}, err
	}

	response := auth.UserResponse{
		ID:             newUser.ID,
		Email:          newUser.Email,
		Name:           newUser.Name,
		Role:           newUser.Role,
		ProfilePicture: newUser.ProfilePicture,
		IsPremium:      newUser.IsPremium,
		PremiumUntil:   newUser.PremiumUntil,
		Headline:       newUser.Headline,
		CreatedAt:      newUser.CreatedAt,
		UpdatedAt:      newUser.UpdatedAt,
	}

	s.log.WithFields(logrus.Fields{
		"id":    newUser.ID,
		"email": newUser.Email,
		"name":  newUser.Name,
		"role":  newUser.Role,
	}).Info("User created successfully")

	return response, nil
}

func (s *authService) Login(c context.Context, req auth.LoginRequest) (auth.LoginResponse, error) {
	repo, err := s.authrepository.NewClient(false)
	if err != nil {
		s.log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to create repository client")
		return auth.LoginResponse{}, err
	}

	foundUser, err := repo.Users.GetUserByEmail(c, req.Email)
	if err != nil {
		s.log.WithFields(logrus.Fields{
			"error": err.Error(),
			"email": req.Email,
		}).Error("Failed to get user by email")
		return auth.LoginResponse{}, err
	}

	if foundUser.ID == "" {
		s.log.WithFields(logrus.Fields{
			"email": req.Email,
		}).Warn("User not found")
		return auth.LoginResponse{}, auth.ErrorInvalidCredentials
	}

	err = bcrypt.ComparePassword(foundUser.Password, req.Password)
	if err != nil {
		s.log.WithFields(logrus.Fields{
			"error": err.Error(),
			"email": req.Email,
		}).Warn("Invalid password")
		return auth.LoginResponse{}, auth.ErrorInvalidCredentials
	}

	userData := MakeUserData(foundUser)
	s.log.WithFields(logrus.Fields{
		"user_id":  userData["id"],
		"username": userData["Username"],
	}).Debug("User authenticated successfully, generating token")

	token, expired, err := jwtPkg.Sign(userData, time.Hour*1)
	if err != nil {
		s.log.WithFields(logrus.Fields{
			"error":   err.Error(),
			"user_id": userData["id"],
		}).Error("Error signing JWT token")
		return auth.LoginResponse{}, err
	}

	loginResponse := auth.LoginResponse{
		Token:     token,
		ExpiresAt: expired,
	}

	s.log.WithFields(logrus.Fields{
		"id":    foundUser.ID,
		"email": foundUser.Email,
		"name":  foundUser.Name,
	}).Info("User logged in successfully")

	return loginResponse, nil
}

func (s *authService) UpdateUser(c context.Context, req auth.UpdateUser) (auth.UserResponse, error) {
	repo, err := s.authrepository.NewClient(false)
	if err != nil {
		s.log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to create repository client")
		return auth.UserResponse{}, err
	}

	// Check if user exists
	existingUser, err := repo.Users.GetUserByID(c, req.ID)
	if err != nil {
		s.log.WithFields(logrus.Fields{
			"error": err.Error(),
			"id":    req.ID,
		}).Error("Failed to get user by ID")
		return auth.UserResponse{}, err
	}

	if existingUser.ID == "" {
		s.log.WithFields(logrus.Fields{
			"id": req.ID,
		}).Warn("User not found")
		return auth.UserResponse{}, auth.ErrorUserNotFound
	}

	// Update user fields
	updatedUser := existingUser
	updatedUser.Name = req.Name
	updatedUser.Role = req.Role
	updatedUser.ProfilePicture = req.ProfilePicture
	updatedUser.IsPremium = req.IsPremium
	updatedUser.PremiumUntil = req.PremiumUntil
	updatedUser.Headline = req.Headline
	updatedUser.UpdatedAt = time.Now()

	if err := repo.Users.UpdateUser(c, updatedUser); err != nil {
		s.log.WithFields(logrus.Fields{
			"error": err.Error(),
			"id":    req.ID,
		}).Error("Failed to update user")
		return auth.UserResponse{}, err
	}

	// Create response (excluding password)
	response := auth.UserResponse{
		ID:             updatedUser.ID,
		Email:          updatedUser.Email,
		Name:           updatedUser.Name,
		Role:           updatedUser.Role,
		ProfilePicture: updatedUser.ProfilePicture,
		IsPremium:      updatedUser.IsPremium,
		PremiumUntil:   updatedUser.PremiumUntil,
		Headline:       updatedUser.Headline,
		CreatedAt:      updatedUser.CreatedAt,
		UpdatedAt:      updatedUser.UpdatedAt,
	}

	s.log.WithFields(logrus.Fields{
		"id":    updatedUser.ID,
		"email": updatedUser.Email,
		"name":  updatedUser.Name,
	}).Info("User updated successfully")

	return response, nil
}
func (s *authService) DeleteUser(c context.Context, id string) error {
	repo, err := s.authrepository.NewClient(false)
	if err != nil {
		s.log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to create repository client")
		return err
	}

	existingUser, err := repo.Users.GetUserByID(c, id)
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
	if err := repo.Users.SoftDeleteUser(c, id, now); err != nil {
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
