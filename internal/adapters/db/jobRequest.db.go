package db

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/tutor-connect-AA/tutor-backend/internal/application/core/domain"
	"gorm.io/gorm"
)

type JobRequestRepo struct {
	db *gorm.DB
}

func NewJobRequestRepo(db *gorm.DB) *JobRequestRepo {
	return &JobRequestRepo{
		db: db,
	}
}

type job_request_table struct {
	gorm.Model
	Id               uuid.UUID               `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Description      string                  `gorm:"not null"` // description for the job by the client
	Status           domain.JobRequestStatus `gorm:"type:text"`
	ClientId         string                  //`gorm:"foreignKey:client_table(id)"` //id of the client who sent the request
	Client_table     client_table            `gorm:"foreignKey:ClientId;references:Id"`
	TutorId          string                  //`gorm:"foreignKey:tutor_table(id)"` // id of the tutor for whom the request is sent
	Tutor_table      tutor_table             `gorm:"foreignKey:TutorId;references:Id"`
	TxRef            string                  `gorm:"text"`
	TutorContactInfo string                  `gorm:"text"`
}

func (jrr JobRequestRepo) CreateJobRequestRepo(newJob domain.JobRequest) (*domain.JobRequest, error) {
	jr := job_request_table{
		Description: newJob.Description,
		Status:      domain.REQUESTED,
		ClientId:    newJob.ClientId,
		TutorId:     newJob.TutorId,
	}
	fmt.Println(jr.ClientId, jr.TutorId)
	if res := jrr.db.Create(&jr); res.Error != nil {
		return nil, res.Error

	}
	newJob.Id = jr.Id.String()
	return &newJob, nil
}

func (jrr JobRequestRepo) JobRequestByIdRepo(id string) (*domain.JobRequest, error) {
	var jr job_request_table
	if res := jrr.db.Where("id=?", id).First(&jr); res.Error != nil {
		return nil, res.Error
	}
	return &domain.JobRequest{
		Id:               jr.Id.String(),
		Description:      jr.Description,
		Status:           jr.Status,
		ClientId:         jr.ClientId,
		TutorId:          jr.TutorId,
		CreatedOn:        jr.CreatedAt,
		TutorContactInfo: jr.TutorContactInfo,
		TxRef:            jr.TxRef,
	}, nil
}

func (jrr JobRequestRepo) JobRequestsRepo() ([]*domain.JobRequest, error) {
	var jrs []job_request_table
	if res := jrr.db.Order("created_at DESC").Find(&jrs); res.Error != nil {
		return nil, res.Error
	}

	var requests []*domain.JobRequest

	for _, request := range jrs {
		domainRequest := &domain.JobRequest{
			Id:               request.Id.String(),
			Description:      request.Description,
			Status:           request.Status,
			ClientId:         request.ClientId,
			TutorId:          request.TutorId,
			TutorContactInfo: request.TutorContactInfo,
			TxRef:            request.TxRef,
		}
		requests = append(requests, domainRequest)
	}

	return requests, nil

}

func (jrr JobRequestRepo) UpdateRequestRepo(requestId string, updatedRequest domain.JobRequest) error {
	res := jrr.db.Model(&job_request_table{}).Where("id = ?", requestId).Updates(updatedRequest)

	return res.Error
}
