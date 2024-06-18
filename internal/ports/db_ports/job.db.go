package db_ports

import "github.com/tutor-connect-AA/tutor-backend/internal/application/core/domain"

type JobDBPort interface {
	CreateJobRepo(job domain.Job) (*domain.Job, error)
	GetJobByIdRepo(id string) (*domain.Job, error)
	GetJobsRepo(offset, limit int) ([]*domain.Job, error)
	UpdateJobRepo(jobId string, updatedJob domain.Job) (*domain.Job, error)
	GetJobByClientRepo(clientId string, offset, limit int) ([]*domain.Job, error)
	GetJobCountRepo() (*int64, error)
	GetJobCountByClientRepo(clientId string) (*int64, error)
	// UpdateJobRepo(job domain.Job) error
}
