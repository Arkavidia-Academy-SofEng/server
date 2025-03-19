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
