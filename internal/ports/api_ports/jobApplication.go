package api_ports

import "github.com/tutor-connect-AA/tutor-backend/internal/application/core/domain"

type JobApplicationAPIPort interface {
	Apply(domain.JobApplication) (*domain.JobApplication, error)
	GetApplicationById(id string) (*domain.JobApplication, error)
	GetApplicationsByJob(jId string) ([]*domain.JobApplication, error)
	GetApplicationsByTutor(tutorId string) ([]*domain.JobApplication, error)
	GetApplicationsByClient(clientId string) ([]*domain.JobApplication, error)
	UpdateApplicationStatus(applicationId string, updatedApp domain.JobApplication) error
}
