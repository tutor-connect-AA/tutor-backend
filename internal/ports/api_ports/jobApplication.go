package api_ports

import "github.com/tutor-connect-AA/tutor-backend/internal/application/core/domain"

type JobApplicationAPIPort interface {
	Apply(domain.JobApplication) (*domain.JobApplication, error)
	GetApplication(id string) (*domain.JobApplication, error)
	GetApplicationsbyJob(jId string) ([]*domain.JobApplication, error)
	GetApplicationsByTutor(tutorId string) ([]*domain.JobApplication, error)
	GetApplicationsByClient(clientId string) ([]*domain.JobApplication, error)
}
