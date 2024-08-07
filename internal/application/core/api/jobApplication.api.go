package api

import (
	"github.com/tutor-connect-AA/tutor-backend/internal/application/core/domain"
	"github.com/tutor-connect-AA/tutor-backend/internal/ports/db_ports"
)

type JobApplicationAPI struct {
	jar db_ports.JobApplicationDBPort
}

func NewJobApplicationAPI(jobApplicationRepo db_ports.JobApplicationDBPort) *JobApplicationAPI {
	return &JobApplicationAPI{
		jar: jobApplicationRepo,
	}
}

func (js *JobApplicationAPI) Apply(apl domain.JobApplication) (*domain.JobApplication, error) {
	ja, err := js.jar.CreateApplicationRepo(apl)
	if err != nil {
		return nil, err
	}
	return ja, nil
}

func (js *JobApplicationAPI) GetApplicationById(id string) (*domain.JobApplication, error) {
	ja, err := js.jar.GetApplicationByIdRepo(id)
	if err != nil {
		return nil, err
	}
	return ja, nil
}

func (js *JobApplicationAPI) GetApplicationsByJob(jId string) ([]*domain.JobApplication, error) {
	apls, err := js.jar.GetApplicationsByJobRepo(jId)
	if err != nil {
		return nil, err
	}
	return apls, nil
}

func (js *JobApplicationAPI) GetApplicationsByTutor(tutorId string) ([]*domain.JobApplication, error) {

	apls, err := js.jar.GetApplicationsByJobRepo(tutorId)
	if err != nil {
		return nil, err
	}
	return apls, nil
}
func (js *JobApplicationAPI) GetApplicationsByClient(clientId string) ([]*domain.JobApplication, error) {

	apls, err := js.jar.GetApplicationsByJobRepo(clientId)
	if err != nil {
		return nil, err
	}
	return apls, nil
}

func (js *JobApplicationAPI) UpdateApplication(applicationId string, updatedApp domain.JobApplication) error {
	err := js.jar.UpdateJobApplicationRepo(applicationId, updatedApp)
	if err != nil {
		return err
	}
	return nil

}
func (js JobApplicationAPI) GetApplicationsByStatus(jId string, status domain.ApplicationStatus) ([]*domain.JobApplication, error) {

	apls, err := js.jar.GetApplicationsByStatusRepo(jId, status)
	if err != nil {
		return nil, err
	}

	return apls, err
}
func (js JobApplicationAPI) UpdateMultipleStatuses(ids []string, newStatus domain.ApplicationStatus) error {
	err := js.jar.UpdateMultipleStatusesRepo(ids, newStatus)

	return err
}

func (js JobApplicationAPI) HasApplied(jobId, tutorId string) (bool, error) {
	applied, err := js.jar.HasAppliedRepo(jobId, tutorId)

	if err != nil {
		return false, err
	}
	return applied, nil
}
