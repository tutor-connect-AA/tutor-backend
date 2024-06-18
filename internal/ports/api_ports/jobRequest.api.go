package api_ports

import "github.com/tutor-connect-AA/tutor-backend/internal/application/core/domain"

type JobRequestAPIPort interface {
	CreateJobRequest(newJob domain.JobRequest) (*domain.JobRequest, error)
	JobRequestById(id string) (*domain.JobRequest, error)
	JobRequests() ([]*domain.JobRequest, error)
	UpdateJobRequest(requestId string, updatedJob domain.JobRequest) error
	HasRequested(clientId, tutorId string) (bool, error)
}
