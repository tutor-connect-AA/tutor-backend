package db_ports

import "github.com/tutor-connect-AA/tutor-backend/internal/application/core/domain"

type JobRequestDBPort interface {
	CreateJobRequestRepo(newJob domain.JobRequest) (*domain.JobRequest, error)
	JobRequestByIdRepo(id string) (*domain.JobRequest, error)
	JobRequestsRepo() ([]*domain.JobRequest, error)
	UpdateRequestRepo(requestId string, updatedRequest domain.JobRequest) error
}
