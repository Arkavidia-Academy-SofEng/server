package recruitmentService

import (
	"ProjectGolang/internal/api/recruitment"
	recruitmentRepository "ProjectGolang/internal/api/recruitment/repository"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

type RecruitmentService interface {
	JobVacancy() JobVacancyDomain
	//JobApplication() AuthDomain
}

type JobVacancyDomain interface {
	CreateJobVacancy(c context.Context, req recruitment.CreateJobVacancy) error
	GetJobVacancies(c context.Context, req recruitment.GetJobVacancies) (recruitment.PaginatedJobVacanciesResponse, error)
	UpdateJobVacancy(c context.Context, req recruitment.UpdateJobVacancy) error
	DeleteJobVacancy(c context.Context, id string) error
}

//type JobApplication interface {
//	UpdatePassword(c context.Context, req auth.ResetPassword) error
//}

type recruitmentService struct {
	recruitmentRepository recruitmentRepository.Repository
	log                   *logrus.Logger

	jobVacancyDomain JobVacancyDomain
	//jobApplication JobApplication
}

func (s *recruitmentService) JobVacancy() JobVacancyDomain {
	return s.jobVacancyDomain
}

//func (r *recruitmentService) JobApplication() JobApplication {
//	return r.jobApplication
//}

type jobVacancyImpl struct {
	repo recruitmentRepository.Repository
	log  *logrus.Logger
}

//type jobApplicationImpl struct {
//	repo recruitmentRepository.Repository
//}

func New(recruitmentRepo recruitmentRepository.Repository,
	log *logrus.Logger,
) RecruitmentService {
	return &recruitmentService{
		recruitmentRepository: recruitmentRepo,
		log:                   log,

		jobVacancyDomain: &jobVacancyImpl{repo: recruitmentRepo, log: log},
		//jobApplication: &jobApplicationImpl{repo: recruitmentRepo},
	}
}
