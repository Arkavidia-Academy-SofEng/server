package scheduler

import (
	authRepository "ProjectGolang/internal/api/auth/repository"
	"context"
	"github.com/go-co-op/gocron"
	"github.com/sirupsen/logrus"
	"time"
)

type Scheduler struct {
	scheduler *gocron.Scheduler
	repo      authRepository.Repository
	log       *logrus.Logger
}

func NewScheduler(repo authRepository.Repository, log *logrus.Logger) *Scheduler {
	return &Scheduler{
		scheduler: gocron.NewScheduler(time.UTC),
		repo:      repo,
		log:       log,
	}
}

func (s *Scheduler) Start() {
	s.scheduler.Every(1).Day().At("03:00").Do(s.cleanupSoftDeletedUsers)
	s.scheduler.StartAsync()
	s.log.Info("Scheduler started successfully")
}

func (s *Scheduler) Stop() {
	s.scheduler.Stop()
	s.log.Info("Scheduler stopped")
}

func (s *Scheduler) cleanupSoftDeletedUsers() {
	s.log.Info("Starting cleanup of soft-deleted users")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	threshold := time.Now().AddDate(0, 0, -15)

	repo, err := s.repo.NewClient(false)
	if err != nil {
		s.log.WithField("error", err.Error()).Error("Failed to create repository client for cleanup")
		return
	}

	err = repo.User.HardDeleteExpiredUsers(ctx, threshold)
	if err != nil {
		s.log.WithField("error", err.Error()).Error("Failed to hard delete expired users")
		return
	}
}
