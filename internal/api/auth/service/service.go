package authService

import (
	"ProjectGolang/internal/api/auth"
	authRepository "ProjectGolang/internal/api/auth/repository"
	"context"
	"github.com/sirupsen/logrus"
)

type authService struct {
	authrepository authRepository.Repository
	log            *logrus.Logger
}

type AuthService interface {
	CreateUser(c context.Context, req auth.CreateUser) (auth.UserResponse, error)
	Login(c context.Context, req auth.LoginRequest) (auth.LoginResponse, error)
	UpdateUser(c context.Context, req auth.UpdateUser) (auth.UserResponse, error)
	DeleteUser(c context.Context, id string) error
}

func New(authRepo authRepository.Repository, log *logrus.Logger) AuthService {
	return &authService{
		authrepository: authRepo,
		log:            log,
	}
}
