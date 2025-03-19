package authService

import (
	"ProjectGolang/internal/api/auth"
	authRepository "ProjectGolang/internal/api/auth/repository"
	"ProjectGolang/pkg/redis"
	"ProjectGolang/pkg/s3"
	"ProjectGolang/pkg/smtp"
	"context"
	"github.com/sirupsen/logrus"
	"mime/multipart"
)

type authService struct {
	authrepository authRepository.Repository
	log            *logrus.Logger
	smtp           smtp.ItfSmtp
	redis          redis.ItfRedis
	s3             s3.ItfS3
}

type AuthService interface {
	RequestOTP(c context.Context, req auth.RequestOTP) error
	CreateUser(c context.Context, req auth.CreateUser) error
	Login(c context.Context, req auth.LoginRequest) (auth.LoginResponse, error)
	UpdateUser(c context.Context, req auth.UpdateUser, id string, banner *multipart.FileHeader, profile *multipart.FileHeader) error
	DeleteUser(c context.Context, id string) error

	CreateCompany(c context.Context, req auth.CreateUser) error
	UpdateCompany(c context.Context, req auth.UpdateCompany, id string, banner *multipart.FileHeader, profile *multipart.FileHeader) error
	DeleteCompany(c context.Context, id string) error
}

func New(authRepo authRepository.Repository,
	log *logrus.Logger,
	smtp smtp.ItfSmtp,
	redis redis.ItfRedis,
	s3 s3.ItfS3) AuthService {
	return &authService{
		authrepository: authRepo,
		log:            log,
		smtp:           smtp,
		redis:          redis,
		s3:             s3,
	}
}
