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

func (ja JobAPI) UpdateJob(jobId string, updatedJob domain.Job) (*domain.Job, error) {
	job, err := ja.jr.UpdateJobRepo(jobId, updatedJob)
	if err != nil {
		return nil, err
	}
	return job, nil
}

func (ja JobAPI) GetJobByClient(clientId string, offset, limit int) ([]*domain.Job, error) {
	jobs, err := ja.jr.GetJobByClientRepo(clientId, offset, limit)
	if err != nil {
		return nil, err
	}
	return jobs, nil
}

func (ja JobAPI) GetJobCount() (*int64, error) {
	count, err := ja.jr.GetJobCountRepo()
	if err != nil {
		return nil, err
	}
	return count, nil

}

func (ja JobAPI) GetJobCountByClient(clientId string) (*int64, error) {
	count, err := ja.jr.GetJobCountByClientRepo(clientId)
	if err != nil {
		return nil, err
	}
	return count, nil
}

// func (ja JobAPI) UpdateJobPost(updatedFieldsObj domain.Job) error {
// 	err := ja.jobDB.UpdateJobRepo(updatedFieldsObj)
// 	return err
// }
