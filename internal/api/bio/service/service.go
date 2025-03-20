package bioService

import (
	authRepository "ProjectGolang/internal/api/auth/repository"
	"ProjectGolang/internal/api/bio"
	bioRepository "ProjectGolang/internal/api/bio/repository"
	"ProjectGolang/internal/entity"
	"ProjectGolang/pkg/redis"
	"ProjectGolang/pkg/s3"
	"ProjectGolang/pkg/smtp"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"mime/multipart"
)

type bioService struct {
	authRepository authRepository.Repository
	bioRepository  bioRepository.Repository
	log            *logrus.Logger
	smtp           smtp.ItfSmtp
	redis          redis.ItfRedis
	s3             s3.ItfS3
}

type BioService interface {
	CreateExperience(ctx context.Context, req bio.CreateExperience, userID string, image *multipart.FileHeader) error
	GetExperienceByID(ctx context.Context, id string) (entity.Experience, error)
	GetExperiencesByUserID(ctx context.Context, userID string) ([]entity.Experience, error)
	UpdateExperience(ctx context.Context, req bio.UpdateExperience, id string, image *multipart.FileHeader) error
	DeleteExperience(ctx context.Context, id string) error

	CreateEducation(ctx context.Context, req bio.CreateEducation, userID string, image *multipart.FileHeader) error
	GetEducationByID(ctx context.Context, id string) (entity.Education, error)
	GetEducationsByUserID(ctx context.Context, userID string) ([]entity.Education, error)
	UpdateEducation(ctx context.Context, req bio.UpdateEducation, id string, image *multipart.FileHeader) error
	DeleteEducation(ctx context.Context, id string) error

	CreatePortfolio(ctx context.Context, req bio.CreatePortfolio, userID string, image *multipart.FileHeader, descriptionImage *multipart.FileHeader) error
	GetPortfolioByID(ctx context.Context, id string) (entity.Portfolio, error)
	GetPortfoliosByUserID(ctx context.Context, userID string) ([]entity.Portfolio, error)
	UpdatePortfolio(ctx context.Context, req bio.UpdatePortfolio, id string, image *multipart.FileHeader, descriptionImage *multipart.FileHeader) error
	DeletePortfolio(ctx context.Context, id string) error
}

func New(authRepo authRepository.Repository, bioRepo bioRepository.Repository,
	log *logrus.Logger,
	smtp smtp.ItfSmtp,
	redis redis.ItfRedis,
	s3 s3.ItfS3) BioService {
	return &bioService{
		authRepository: authRepo,
		bioRepository:  bioRepo,
		log:            log,
		smtp:           smtp,
		redis:          redis,
		s3:             s3,
	}
}
