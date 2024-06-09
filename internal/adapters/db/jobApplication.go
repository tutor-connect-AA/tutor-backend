package db

import (
	"github.com/google/uuid"
	"github.com/tutor-connect-AA/tutor-backend/internal/application/core/domain"
	"gorm.io/gorm"
)

type job_application_table struct {
	gorm.Model
	Id                 uuid.UUID                `gorm:"type:uuid;default:uuid_generate_v4()"`
	JobId              string                   //`gorm:"not null;foreignKey:job_table(Id)"`
	Job_Table          job_table                `gorm:"foreignKey:job_id;references:Id"`
	ApplicantId        string                   //`gorm:"not null;foreignKey:tutor_table(Id)"`
	Tutor_table        tutor_table              `gorm:"foreignKey:applicant_id;references:Id"`
	CoverLetter        string                   `gorm:"type:text"` // Assuming text storage for cover letter
	Status             domain.ApplicationStatus `gorm:"type:text"`
	TxRef              string                   `gorm:"text"` //this was unique so that one tx_ref can't be used to hire multiple times (removed now find other ways)
	InterviewQuestions string                   `gorm:"text"`
	InterviewResponse  string
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
	apl.Id = newApplication.Id.String()
	apl.Status = newApplication.Status
	return &apl, nil
}
func (jar JobApplicationRepo) GetApplicationByIdRepo(id string) (*domain.JobApplication, error) {
	var application job_application_table
	res := jar.db.Where("id =?", id).First(&application)

	if res.Error != nil {
		return nil, res.Error
	}
	return &domain.JobApplication{
		Id:                 application.Id.String(),
		JobId:              application.JobId,
		ApplicantId:        application.ApplicantId,
		CoverLetter:        application.CoverLetter,
		Status:             application.Status,
		InterviewQuestions: application.InterviewQuestions,
		InterviewResponse:  application.InterviewResponse,
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
			Id:                 apl.Id.String(),
			JobId:              apl.JobId,
			ApplicantId:        apl.ApplicantId,
			CoverLetter:        apl.CoverLetter,
			Status:             apl.Status,
			InterviewQuestions: apl.InterviewQuestions,
			InterviewResponse:  apl.InterviewResponse,
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
			Id:                 apl.Id.String(),
			JobId:              apl.JobId,
			ApplicantId:        apl.ApplicantId,
			CoverLetter:        apl.CoverLetter,
			Status:             apl.Status,
			InterviewQuestions: apl.InterviewQuestions,
			InterviewResponse:  apl.InterviewResponse,
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
			Id:                 apl.Id.String(),
			JobId:              apl.JobId,
			ApplicantId:        apl.ApplicantId,
			CoverLetter:        apl.CoverLetter,
			Status:             apl.Status,
			InterviewQuestions: apl.InterviewQuestions,
			InterviewResponse:  apl.InterviewResponse,
		}
		applications = append(applications, &newApl)
	}

	return applications, nil
}

func (jar JobApplicationRepo) UpdateJobApplicationRepo(applicationId string, updatedApp domain.JobApplication) error {
	res := jar.db.Model(&job_application_table{}).Where("id = ?", applicationId).Updates(updatedApp)

	return res.Error
}
func (jar JobApplicationRepo) GetApplicationsByStatusRepo(jId string, status domain.ApplicationStatus) ([]*domain.JobApplication, error) {

	var aplsByStatus []job_application_table
	res := jar.db.
		Where("status = ?", status).
		Where("job_id = ?", jId).
		Order("updated_at").
		Find(&aplsByStatus)

	if res.Error != nil {
		return nil, res.Error
	}

	var applications []*domain.JobApplication
	for _, apl := range aplsByStatus {
		var newApl = domain.JobApplication{
			Id:                 apl.Id.String(),
			JobId:              apl.JobId,
			ApplicantId:        apl.ApplicantId,
			CoverLetter:        apl.CoverLetter,
			Status:             apl.Status,
			InterviewQuestions: apl.InterviewQuestions,
			InterviewResponse:  apl.InterviewResponse,
		}
		applications = append(applications, &newApl)
	}

	return applications, nil

}
