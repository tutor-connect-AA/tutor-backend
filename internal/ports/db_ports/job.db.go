package db_ports

import "github.com/tutor-connect-AA/tutor-backend/internal/application/core/domain"

type JobDBPort interface {
	CreateJobRepo(job domain.Job) (*domain.Job, error)
	GetJobByIdRepo(id string) (*domain.Job, error)
	GetJobsRepo() ([]*domain.Job, error)
	// UpdateJobRepo(job domain.Job) error
}
