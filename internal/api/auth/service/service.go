package authService

import (
	"ProjectGolang/internal/api/auth"
	authRepository "ProjectGolang/internal/api/auth/repository"
	"ProjectGolang/internal/entity"
	"context"
	"github.com/sirupsen/logrus"
)

type authService struct {
	authrepository authRepository.Repository
	log            *logrus.Logger
}

type AuthService interface {
	RegisterUser(ctx context.Context, req auth.CreateUserRequest) error
	Login(ctx context.Context, req auth.LoginUserRequest) (auth.LoginUserResponse, error)
	UpdateUser(ctx context.Context, user entity.UserLoginData, req auth.UpdateUserRequest) error
	DeleteUser(ctx context.Context, id string) error
}

func New(authRepo authRepository.Repository, log *logrus.Logger) AuthService {
	return &authService{
		authrepository: authRepo,
		log:            log,
	}
}
