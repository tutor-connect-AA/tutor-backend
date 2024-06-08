package api

import (
	"github.com/tutor-connect-AA/tutor-backend/internal/application/core/domain"
	"github.com/tutor-connect-AA/tutor-backend/internal/ports/db_ports"
)

type JobRequestAPI struct {
	jrR db_ports.JobRequestDBPort
}

func NewJobRequestAPI(jrR db_ports.JobRequestDBPort) *JobRequestAPI {
	return &JobRequestAPI{
		jrR: jrR,
	}
}

func (jrS JobRequestAPI) CreateJobRequest(newJob domain.JobRequest) (*domain.JobRequest, error) {
	jr, err := jrS.jrR.CreateJobRequestRepo(newJob)
	if err != nil {
		return nil, err
	}
	return jr, nil
}

func (jrS JobRequestAPI) JobRequestById(id string) (*domain.JobRequest, error) {
	jr, err := jrS.jrR.JobRequestByIdRepo(id)
	if err != nil {
		return nil, err
	}
	return jr, nil
}

func (jrs JobRequestAPI) JobRequests() ([]*domain.JobRequest, error) {
	jrList, err := jrs.jrR.JobRequestsRepo()
	if err != nil {
		return nil, err
	}
	return jrList, err
}

func (jrs JobRequestAPI) UpdateJobRequest(requestId string, updatedJob domain.JobRequest) error {
	err := jrs.jrR.UpdateRequestRepo(requestId, updatedJob)
	return err
}
