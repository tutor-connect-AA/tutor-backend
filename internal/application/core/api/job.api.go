package api

import (
	"github.com/tutor-connect-AA/tutor-backend/internal/application/core/domain"
	"github.com/tutor-connect-AA/tutor-backend/internal/ports/db_ports"
)

type JobAPI struct {
	jr db_ports.JobDBPort
}

func NewJobAPI(jr db_ports.JobDBPort) *JobAPI {
	return &JobAPI{
		jr: jr,
	}
}

func (ja JobAPI) CreateJobPost(job domain.Job) (*domain.Job, error) {
	j, err := ja.jr.CreateJobRepo(job)
	if err != nil {
		return nil, err
	}
	return j, nil
}

func (ja JobAPI) GetJob(id string) (*domain.Job, error) {
	j, err := ja.jr.GetJobByIdRepo(id)
	if err != nil {
		return nil, err
	}
	return j, nil
}

func (ja JobAPI) GetListOfJobs(offset, limit int) ([]*domain.Job, error) {
	js, err := ja.jr.GetJobsRepo(offset, limit)
	if err != nil {
		return nil, err
	}
	return js, nil
}

// func (ja JobAPI) UpdateJobPost(updatedFieldsObj domain.Job) error {
// 	err := ja.jobDB.UpdateJobRepo(updatedFieldsObj)
// 	return err
// }
