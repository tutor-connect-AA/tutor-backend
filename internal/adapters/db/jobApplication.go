package db

import (
	"github.com/google/uuid"
	"github.com/tutor-connect-AA/tutor-backend/internal/application/core/domain"
	"gorm.io/gorm"
)

type job_application_table struct {
	gorm.Model
	Id          uuid.UUID                `gorm:"type:uuid;default:uuid_generate_v4()"`
	JobId       string                   `gorm:"not null;foreignKey:job_table(Id);primaryKey"`
	ApplicantId string                   `gorm:"not null;foreignKey:tutor_table(Id)";primaryKey`
	CoverLetter string                   `gorm:"type:text"` // Assuming text storage for cover letter
	Status      domain.ApplicationStatus `gorm:"type:text"`
	// Tutors      []*tutor_table `foreignKey:applicant_id"`
	// File        string    `gorm:"type:text"` // Assuming text storage for other documents (link can also be used)

}

type JobApplicationRepo struct {
	db *gorm.DB
}

func NewJobApplicationRepo(db *gorm.DB) *JobApplicationRepo {
	return &JobApplicationRepo{
		db: db,
	}
}

func (jar JobApplicationRepo) CreateApplicationRepo(apl domain.JobApplication) (*domain.JobApplication, error) {
	var newApplication = &job_application_table{
		JobId:       apl.JobId,
		ApplicantId: apl.ApplicantId,
		CoverLetter: apl.CoverLetter,
		Status:      domain.PENDING,
		// File:        apl.File,
	}

	res := jar.db.Create(&newApplication)
	if res.Error != nil {
		return nil, res.Error
	}
	return &apl, nil
}
func (jar JobApplicationRepo) GetApplicationByIdRepo(id string) (*domain.JobApplication, error) {
	var application job_application_table
	res := jar.db.Where("id = ?", id).First(&application)

	if res.Error != nil {
		return nil, res.Error
	}
	return &domain.JobApplication{
		Id:          application.Id.String(),
		JobId:       application.JobId,
		ApplicantId: application.ApplicantId,
		CoverLetter: application.CoverLetter,
		Status:      application.Status,
		// File:        application.File,
	}, nil
}

func (jar JobApplicationRepo) GetApplicationsByJobRepo(jId string) ([]*domain.JobApplication, error) {
	var aplsByJob []job_application_table
	res := jar.db.Where("job_id = ?", jId).Find(&aplsByJob)
	if res.Error != nil {
		return nil, res.Error
	}
	var applications []*domain.JobApplication
	for _, apl := range aplsByJob {
		var newApl = domain.JobApplication{
			Id:          apl.Id.String(),
			JobId:       apl.JobId,
			ApplicantId: apl.ApplicantId,
			CoverLetter: apl.CoverLetter,
			Status:      apl.Status,
		}
		applications = append(applications, &newApl)
	}

	return applications, nil
}

func (jar JobApplicationRepo) GetApplicationsByTutorRepo(tId string) ([]*domain.JobApplication, error) {
	var aplsByTutor []job_application_table
	// res := jar.db.Where("applicant_id = ?", tId).Find(&aplsByTutor)
	res := jar.db.Where("applicant_id = ?", tId).Order("created_at").Find(&aplsByTutor)

	if res.Error != nil {
		return nil, res.Error
	}
	var applications []*domain.JobApplication
	for _, apl := range aplsByTutor {
		var newApl = domain.JobApplication{
			Id:          apl.Id.String(),
			JobId:       apl.JobId,
			ApplicantId: apl.ApplicantId,
			CoverLetter: apl.CoverLetter,
			Status:      apl.Status,
		}
		applications = append(applications, &newApl)
	}

	return applications, nil
}

func (jar JobApplicationRepo) GetApplicationsByClientRepo(cltId string) ([]*domain.JobApplication, error) {
	var aplsByClt []job_application_table
	// res := jar.db.Where("applicant_id = ?", tId).Find(&aplsByTutor)
	res := jar.db.Where("posted_by = ?", cltId).Order("created_at").Find(&aplsByClt)

	if res.Error != nil {
		return nil, res.Error
	}
	var applications []*domain.JobApplication
	for _, apl := range aplsByClt {
		var newApl = domain.JobApplication{
			Id:          apl.Id.String(),
			JobId:       apl.JobId,
			ApplicantId: apl.ApplicantId,
			CoverLetter: apl.CoverLetter,
			Status:      apl.Status,
		}
		applications = append(applications, &newApl)
	}

	return applications, nil
}

func (jar JobApplicationRepo) UpdateApplicationStatusRepo(applicationId string, status domain.ApplicationStatus) error {
	res := jar.db.Table("job_application_tables").Where("id = ?", applicationId).Update("status", status)

	if res.Error != nil {
		return res.Error
	}
	return nil
}
