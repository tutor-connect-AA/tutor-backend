package db_ports

import "github.com/tutor-connect-AA/tutor-backend/internal/application/core/domain"

type JobApplicationDBPort interface {
	GetApplicationByIdRepo(id string) (*domain.JobApplication, error)
	CreateApplicationRepo(apl domain.JobApplication) (*domain.JobApplication, error)
	GetApplicationsByJobRepo(jId string) ([]*domain.JobApplication, error)
	GetApplicationsByTutorRepo(tId string) ([]*domain.JobApplication, error)
	GetApplicationsByClientRepo(tId string) ([]*domain.JobApplication, error)
	UpdateJobApplicationRepo(applicationId string, updatedApp domain.JobApplication) error
	GetApplicationsByStatusRepo(jId string, status domain.ApplicationStatus) ([]*domain.JobApplication, error)
	UpdateMultipleStatusesRepo(ids []string, newStatus domain.ApplicationStatus) error
}
