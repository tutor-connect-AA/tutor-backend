package db

import (
	"time"

	"github.com/google/uuid"
	"github.com/tutor-connect-AA/tutor-backend/internal/application/core/domain"
	"gorm.io/gorm"
)

type JobRepo struct {
	db *gorm.DB
}

func NewJobRepo(db *gorm.DB) *JobRepo {
	return &JobRepo{
		db: db,
	}
}

type job_table struct {
	gorm.Model
	Id                    uuid.UUID    `gorm:"type:uuid;default:uuid_generate_v4()"`
	Title                 string       `gorm:"size: 255"`
	Description           string       `gorm:"not null"`
	Posted_By             string       //`gorm:"foreignKey:Posted_By` fk
	Client_table          client_table `gorm:"foreignKey:Posted_By;references:Id"`
	Deadline              time.Time
	Region                string
	City                  string
	Quantity              int    //number of tutors needed
	Grade_Of_Students     string //[]int `gorm:"type:integer[]"`
	Minimum_Education     domain.Education
	Preferred_Gender      domain.Gender
	Contact_Hour_Per_Week int
	Status                domain.JobStatus
	Hourly_Rate_Min       int                     `gorm:"column:hourly_rate_min;check:Hourly_Rate_Min > 0 "`
	Hourly_Rate_Max       int                     `gorm:"column:hourly_mate_max;check:Hourly_Rate_Min > 0 "`
	Applications          []job_application_table `gorm:"foreignKey:job_id;references:Id"`
	Interview_Questions   string
	// Clt                   client_table
}

func (jr JobRepo) CreateJobRepo(job domain.Job) (*domain.Job, error) {
	j := &job_table{
		Title:                 job.Title,
		Description:           job.Description,
		Posted_By:             job.Posted_By,
		Deadline:              job.Deadline,
		Region:                job.Region,
		City:                  job.City,
		Quantity:              job.Quantity,
		Grade_Of_Students:     job.Grade_Of_Students,
		Minimum_Education:     job.Minimum_Education,
		Preferred_Gender:      job.Preferred_Gender,
		Contact_Hour_Per_Week: job.Contact_Hour_Per_Week,
		Status:                job.Status,
		Hourly_Rate_Min:       job.Hourly_Rate_Min,
		Hourly_Rate_Max:       job.Hourly_Rate_Max,
	}

	res := jr.db.Create(&j)
	if res.Error != nil {
		return nil, res.Error
	}

	job.Id = j.Id.String()
	return &job, nil
}

func (jr JobRepo) GetJobByIdRepo(id string) (*domain.Job, error) {
	var jb *job_table
	res := jr.db.Where("id = ?", id).First(&jb)
	if res.Error != nil {
		return nil, res.Error
	}
	return &domain.Job{
		Id:                    jb.Id.String(),
		Title:                 jb.Title,
		Description:           jb.Description,
		Posted_By:             jb.Posted_By,
		Deadline:              jb.Deadline,
		Region:                jb.Region,
		City:                  jb.City,
		Quantity:              jb.Quantity,
		Grade_Of_Students:     jb.Grade_Of_Students,
		Minimum_Education:     jb.Minimum_Education,
		Preferred_Gender:      jb.Preferred_Gender,
		Contact_Hour_Per_Week: jb.Contact_Hour_Per_Week,
		Status:                jb.Status,
		Hourly_Rate_Min:       jb.Hourly_Rate_Min,
		Hourly_Rate_Max:       jb.Hourly_Rate_Max,
		Interview_Questions:   jb.Interview_Questions,
	}, nil
}

func (jr JobRepo) GetJobsRepo(offset, limit int) ([]*domain.Job, error) {
	var jbs []*job_table

	var jobList []*domain.Job

	if err := jr.db.Order("created_at DESC").Offset(offset).Limit(limit).Find(&jbs).Error; err != nil {
		return nil, err
	}

	// res := jr.db.Find(&jbs)
	// if res.Error != nil {
	// 	return nil, res.Error
	// }

	for _, job := range jbs {
		oneJob := &domain.Job{
			Id:                    job.Id.String(),
			Title:                 job.Title,
			Description:           job.Description,
			Posted_By:             job.Posted_By,
			Deadline:              job.Deadline,
			Region:                job.Region,
			City:                  job.City,
			Quantity:              job.Quantity,
			Grade_Of_Students:     job.Grade_Of_Students,
			Minimum_Education:     job.Minimum_Education,
			Preferred_Gender:      job.Preferred_Gender,
			Contact_Hour_Per_Week: job.Contact_Hour_Per_Week,
			Status:                job.Status,
			Hourly_Rate_Min:       job.Hourly_Rate_Min,
			Hourly_Rate_Max:       job.Hourly_Rate_Max,
			Interview_Questions:   job.Interview_Questions,
		}
		jobList = append(jobList, oneJob)
	}
	return jobList, nil
}

func (jr JobRepo) UpdateJobRepo(jobId string, updatedJob domain.Job) (*domain.Job, error) {
	if res := jr.db.Model(&job_table{}).Where("id = ?", jobId).Updates(updatedJob); res.Error != nil {
		return nil, res.Error
	}

	var updatedRecord domain.Job
	jr.db.First(&updatedRecord, jobId)

	return &updatedRecord, nil
}

func (jr JobRepo) GetJobByClientRepo(clientId string, offset, limit int) ([]*domain.Job, error) {
	var jbs []*job_table

	var jobList []*domain.Job

	if err := jr.db.Order("created_at DESC").
		Where("posted_by = ?", clientId).
		Offset(offset).
		Limit(limit).
		Find(&jbs).Error; err != nil {
		return nil, err
	}

	for _, job := range jbs {
		oneJob := &domain.Job{
			Id:                    job.Id.String(),
			Title:                 job.Title,
			Description:           job.Description,
			Posted_By:             job.Posted_By,
			Deadline:              job.Deadline,
			Region:                job.Region,
			City:                  job.City,
			Quantity:              job.Quantity,
			Grade_Of_Students:     job.Grade_Of_Students,
			Minimum_Education:     job.Minimum_Education,
			Preferred_Gender:      job.Preferred_Gender,
			Contact_Hour_Per_Week: job.Contact_Hour_Per_Week,
			Status:                job.Status,
			Hourly_Rate_Min:       job.Hourly_Rate_Min,
			Hourly_Rate_Max:       job.Hourly_Rate_Max,
			Interview_Questions:   job.Interview_Questions,
		}
		jobList = append(jobList, oneJob)
	}
	return jobList, nil
}

func (jr JobRepo) GetJobCountRepo() (*int64, error) {
	var count int64
	if err := jr.db.Model(&job_table{}).Count(&count).Error; err != nil {
		return nil, err
	}
	return &count, nil
}
func (jr JobRepo) GetJobCountByClientRepo(clientId string) (*int64, error) {
	var count int64
	if err := jr.db.Model(&job_table{}).
		Where("posted_by = ?", clientId).
		Count(&count).Error; err != nil {
		return nil, err
	}
	return &count, nil

}
