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

func (js *JobApplicationAPI) UpdateApplicationStatus(applicationId string, updatedApp domain.JobApplication) error {
	err := js.jar.UpdateApplicationStatusRepo(applicationId, updatedApp)
	if err != nil {
		return err
	}
	return nil

}
