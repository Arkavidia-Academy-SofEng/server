package recruitmentService

import (
	"ProjectGolang/internal/api/recruitment"
	"ProjectGolang/internal/entity"
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"time"
)

func (s *jobVacancyImpl) CreateJobVacancy(c context.Context, req recruitment.CreateJobVacancy) error {
	repo, err := s.repo.NewClient(false)
	if err != nil {
		s.log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to create repository client")
	}

	jobVacancy := entity.JobVacancy{
		Title:       req.Title,
		Description: req.Description,
		Location:    req.Location,
		JobType:     req.JobType,
		Deadline:    req.Deadline,
		IsActive:    req.IsActive,
	}

	if err := repo.JobVacancies.CreateJobVacancy(c, jobVacancy); err != nil {
		s.log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to create job vacancy")
		return err
	}

	s.log.WithFields(logrus.Fields{
		"title":       req.Title,
		"description": req.Description,
		"location":    req.Location,
		"job_type":    req.JobType,
		"deadline":    req.Deadline,
		"is_active":   req.IsActive,
	}).Info("Job vacancy created successfully")

	return nil
}

func (s *jobVacancyImpl) GetJobVacancies(c context.Context, req recruitment.GetJobVacancies) (recruitment.PaginatedJobVacanciesResponse, error) {
	repo, err := s.repo.NewClient(false)
	if err != nil {
		s.log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to create repository client")
		return recruitment.PaginatedJobVacanciesResponse{}, err
	}

	jobVacancies, totalCount, err := repo.JobVacancies.GetJobVacancies(c, req.Page, req.PageSize)
	if err != nil {
		s.log.WithFields(logrus.Fields{
			"error": err.Error(),
			"page":  req.Page,
			"size":  req.PageSize,
		}).Error("Failed to fetch job vacancies")
		return recruitment.PaginatedJobVacanciesResponse{}, err
	}

	totalPages := totalCount / req.PageSize
	if totalCount%req.PageSize > 0 {
		totalPages++
	}

	jobVacancyResponses := make([]recruitment.JobVacancyResponse, len(jobVacancies))
	for i, jv := range jobVacancies {
		jobVacancyResponses[i] = recruitment.JobVacancyResponse{
			ID:           jv.ID,
			RecruiterID:  jv.RecruiterID,
			Title:        jv.Title,
			Description:  jv.Description,
			Requirements: jv.Requirements,
			Location:     jv.Location,
			JobType:      jv.JobType,
			Deadline:     jv.Deadline,
			IsActive:     jv.IsActive,
			CreatedAt:    jv.CreatedAt,
			UpdatedAt:    jv.UpdatedAt,
		}
	}

	response := recruitment.PaginatedJobVacanciesResponse{
		JobVacancies: jobVacancyResponses,
		TotalCount:   totalCount,
		TotalPages:   totalPages,
		CurrentPage:  req.Page,
		PageSize:     req.PageSize,
	}

	s.log.WithFields(logrus.Fields{
		"page":       req.Page,
		"size":       req.PageSize,
		"total":      totalCount,
		"found":      len(jobVacancies),
		"totalPages": totalPages,
	}).Info("Job vacancies fetched successfully")

	return response, nil
}

func (s *jobVacancyImpl) UpdateJobVacancy(c context.Context, req recruitment.UpdateJobVacancy) error {
	repo, err := s.repo.NewClient(false)
	if err != nil {
		s.log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to create repository client")
		return err
	}

	exists, err := repo.JobVacancies.CheckJobVacancyExists(c, req.ID)
	if err != nil {
		s.log.WithFields(logrus.Fields{
			"error": err.Error(),
			"id":    req.ID,
		}).Error("Failed to check if job vacancy exists")
		return err
	}

	if !exists {
		s.log.WithFields(logrus.Fields{
			"id": req.ID,
		}).Warn("Job vacancy not found")
		return fmt.Errorf("job vacancy with ID %s not found", req.ID)
	}

	jobVacancy := entity.JobVacancy{
		ID:           req.ID,
		Title:        req.Title,
		Description:  req.Description,
		Requirements: req.Requirements,
		Location:     req.Location,
		JobType:      req.JobType,
		Deadline:     req.Deadline,
		IsActive:     req.IsActive,
		UpdatedAt:    time.Now(),
	}

	if err := repo.JobVacancies.UpdateJobVacancy(c, jobVacancy); err != nil {
		s.log.WithFields(logrus.Fields{
			"error": err.Error(),
			"id":    req.ID,
		}).Error("Failed to update job vacancy")
		return err
	}

	s.log.WithFields(logrus.Fields{
		"id":          req.ID,
		"title":       req.Title,
		"description": req.Description,
		"location":    req.Location,
		"job_type":    req.JobType,
		"deadline":    req.Deadline,
		"is_active":   req.IsActive,
	}).Info("Job vacancy updated successfully")

	return nil
}

func (s *jobVacancyImpl) DeleteJobVacancy(c context.Context, id string) error {
	repo, err := s.repo.NewClient(false)
	if err != nil {
		s.log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to create repository client")
		return err
	}

	exists, err := repo.JobVacancies.CheckJobVacancyExists(c, id)
	if err != nil {
		s.log.WithFields(logrus.Fields{
			"error": err.Error(),
			"id":    id,
		}).Error("Failed to check if job vacancy exists")
		return err
	}

	if !exists {
		s.log.WithFields(logrus.Fields{
			"id": id,
		}).Warn("Job vacancy not found")
		return fmt.Errorf("job vacancy with ID %s not found", id)
	}

	if err := repo.JobVacancies.DeleteJobVacancy(c, id); err != nil {
		s.log.WithFields(logrus.Fields{
			"error": err.Error(),
			"id":    id,
		}).Error("Failed to delete job vacancy")
		return err
	}

	s.log.WithFields(logrus.Fields{
		"id": id,
	}).Info("Job vacancy deleted successfully")

	return nil
}
