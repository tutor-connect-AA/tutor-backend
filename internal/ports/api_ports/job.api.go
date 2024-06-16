package api_ports

import (
	"github.com/tutor-connect-AA/tutor-backend/internal/application/core/domain"
)

type JobAPIPort interface {
	CreateJobPost(job domain.Job) (*domain.Job, error)
	GetJob(id string) (*domain.Job, error)
	GetListOfJobs(offset, limit int) ([]*domain.Job, error)
	UpdateJob(jobId string, updatedJob domain.Job) (*domain.Job, error)
	// UpdateJobPost(updatedFieldsObj domain.Job) error
}
