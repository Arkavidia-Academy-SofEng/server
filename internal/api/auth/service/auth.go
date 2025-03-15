package authService

import (
	"ProjectGolang/internal/api/auth"
	"ProjectGolang/internal/entity"
	"ProjectGolang/pkg/bcrypt"
	jwtPkg "ProjectGolang/pkg/jwt"
	"ProjectGolang/pkg/response"
	"context"
	"github.com/sirupsen/logrus"
	"time"
)

func (s *authService) RegisterUser(c context.Context, req auth.CreateUserRequest) error {
	s.log.WithFields(logrus.Fields{
		"email":    req.Email,
		"username": req.Username,
	}).Debug("Starting user registration process")

	repo, err := s.authrepository.NewClient(false)
	if err != nil {
		s.log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to create repository client")
		return err
	}

	s.log.Debug("Hashing password")
	hashedPassword, err := bcrypt.HashPassword(req.Password)
	if err != nil {
		s.log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to hash password")
		return err
	}

	s.log.Debug("Generating ULID for new user")
	ulid, err := NewUlidFromTimestamp(time.Now())
	if err != nil {
		s.log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to generate ULID")
		return err
	}

	user := entity.User{
		ID:       ulid,
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
	}

	s.log.WithFields(logrus.Fields{
		"user_id":  ulid,
		"username": req.Username,
		"email":    req.Email,
	}).Info("Creating new user in database")

	if err := repo.Users.CreateUser(c, user); err != nil {
		s.log.WithFields(logrus.Fields{
			"error":    err.Error(),
			"user_id":  ulid,
			"username": req.Username,
			"email":    req.Email,
		}).Error("Failed to create user in database")
		return err
	}

	s.log.WithFields(logrus.Fields{
		"user_id":  ulid,
		"username": req.Username,
		"email":    req.Email,
	}).Info("User created successfully")

	return nil
}

func (s *authService) Login(c context.Context, req auth.LoginUserRequest) (auth.LoginUserResponse, error) {
	s.log.WithFields(logrus.Fields{
		"email": req.Email,
	}).Debug("Processing login request")

	repo, err := s.authrepository.NewClient(false)
	if err != nil {
		s.log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to create repository client for login")
		return auth.LoginUserResponse{}, err
	}

	user, err := repo.Users.GetByEmail(c, req.Email)
	if err != nil {
		s.log.WithFields(logrus.Fields{
			"error": err.Error(),
			"email": req.Email,
		}).Warn("User not found or database error during login")
		return auth.LoginUserResponse{}, err
	}

	s.log.WithFields(logrus.Fields{
		"user_id": user.ID,
		"email":   req.Email,
	}).Debug("User found, comparing passwords")

	if err := bcrypt.ComparePassword(user.Password, req.Password); err != nil {
		s.log.WithFields(logrus.Fields{
			"user_id": user.ID,
			"email":   req.Email,
		}).Warn("Invalid password during login")
		return auth.LoginUserResponse{}, response.New(401, "Invalid email or password")
	}

	userData := MakeUserData(user)
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
		return auth.LoginUserResponse{}, err
	}

	s.log.WithFields(logrus.Fields{
		"user_id": userData["id"],
		"expires": expired,
	}).Info("Login successful, token created")

	res := auth.LoginUserResponse{
		AccessToken:      token,
		ExpiresInMinutes: time.Until(time.Unix(expired, 0)).Minutes(),
	}

	return res, nil
}

func (s *authService) UpdateUser(c context.Context, user entity.UserLoginData, req auth.UpdateUserRequest) error {
	s.log.WithFields(logrus.Fields{
		"user_id": user.ID,
	}).Debug("Processing user update request")

	repo, err := s.authrepository.NewClient(false)
	if err != nil {
		s.log.WithFields(logrus.Fields{
			"error":   err.Error(),
			"user_id": user.ID,
		}).Error("Failed to create repository client for update")
		return err
	}

	s.log.WithFields(logrus.Fields{
		"user_id": user.ID,
	}).Debug("Retrieving current user data")

	userData, err := repo.Users.GetByID(c, user.ID)
	if err != nil {
		s.log.WithFields(logrus.Fields{
			"error":   err.Error(),
			"user_id": user.ID,
		}).Error("Failed to retrieve user data for update")
		return err
	}

	s.log.Debug("Calculating user data differences")
	newUser, err := GetUserDifferenceData(userData, req)
	if err != nil {
		s.log.WithFields(logrus.Fields{
			"error":   err.Error(),
			"user_id": user.ID,
		}).Error("Failed to process user data differences")
		return err
	}

	s.log.WithFields(logrus.Fields{
		"user_id":  user.ID,
		"username": newUser.Username,
	}).Info("Updating user in database")

	if err := repo.Users.UpdateUser(c, newUser, userData.ID); err != nil {
		s.log.WithFields(logrus.Fields{
			"error":   err.Error(),
			"user_id": user.ID,
		}).Error("Failed to update user in database")
		return err
	}

	s.log.WithFields(logrus.Fields{
		"user_id": user.ID,
	}).Info("User updated successfully")

	return nil
}

func (s *authService) DeleteUser(c context.Context, id string) error {
	s.log.WithFields(logrus.Fields{
		"user_id": id,
	}).Debug("Processing user deletion request")

	repo, err := s.authrepository.NewClient(false)
	if err != nil {
		s.log.WithFields(logrus.Fields{
			"error":   err.Error(),
			"user_id": id,
		}).Error("Failed to create repository client for deletion")
		return err
	}

	s.log.WithFields(logrus.Fields{
		"user_id": id,
	}).Info("Deleting user from database")

	if err := repo.Users.DeleteUser(c, id); err != nil {
		s.log.WithFields(logrus.Fields{
			"error":   err.Error(),
			"user_id": id,
		}).Error("Failed to delete user from database")
		return err
	}

	s.log.WithFields(logrus.Fields{
		"user_id": id,
	}).Info("User deleted successfully")

	return nil
}
