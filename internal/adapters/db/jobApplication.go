package db

import (
	"github.com/google/uuid"
	"github.com/tutor-connect-AA/tutor-backend/internal/application/core/domain"
	"gorm.io/gorm"
)

type job_application_table struct {
	gorm.Model
	Id                uuid.UUID                `gorm:"type:uuid;default:uuid_generate_v4()"`
	JobId             string                   //`gorm:"not null;foreignKey:job_table(Id)"`
	Job_Table         job_table                `gorm:"foreignKey:job_id;references:Id"`
	ApplicantId       string                   //`gorm:"not null;foreignKey:tutor_table(Id)"`
	Tutor_table       tutor_table              `gorm:"foreignKey:applicant_id;references:Id"`
	CoverLetter       string                   `gorm:"type:text"` // Assuming text storage for cover letter
	Status            domain.ApplicationStatus `gorm:"type:text"`
	TxRef             string                   `gorm:"text"` //this was unique so that one tx_ref can't be used to hire multiple times (removed now find other ways)
	InterviewResponse string
	TutorContactInfo  string `gorm:"text"`
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
	// applied, err := jar.HasAppliedRepo(apl.JobId, apl.ApplicantId)
	// if err != nil {
	// 	return nil, err
	// }
	// if applied {
	// 	return nil, errors.New("a tutor can only apply once for a job")
	// }
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
		Id:                application.Id.String(),
		JobId:             application.JobId,
		ApplicantId:       application.ApplicantId,
		CoverLetter:       application.CoverLetter,
		Status:            application.Status,
		InterviewResponse: application.InterviewResponse,
		TutorContactInfo:  application.TutorContactInfo,
		// File:        application.File,
	}, nil
}

func (jar JobApplicationRepo) GetApplicationsByJobRepo(jId string) ([]*domain.JobApplication, error) {
	var aplsByJob []job_application_table
	res := jar.db.Order("created_at DESC").
		Where("job_id = ?", jId).
		Find(&aplsByJob)

	if res.Error != nil {
		return nil, res.Error
	}
	var applications []*domain.JobApplication
	for _, apl := range aplsByJob {
		var newApl = domain.JobApplication{
			Id:                apl.Id.String(),
			JobId:             apl.JobId,
			ApplicantId:       apl.ApplicantId,
			CoverLetter:       apl.CoverLetter,
			Status:            apl.Status,
			InterviewResponse: apl.InterviewResponse,
			TutorContactInfo:  apl.TutorContactInfo,
		}
		applications = append(applications, &newApl)
	}

	return applications, nil
}

func (jar JobApplicationRepo) GetApplicationsByTutorRepo(tId string) ([]*domain.JobApplication, error) {
	var aplsByTutor []job_application_table
	// res := jar.db.Where("applicant_id = ?", tId).Find(&aplsByTutor)
	res := jar.db.Where("applicant_id = ?", tId).
		Order("created_at DESC").
		Find(&aplsByTutor)

	if res.Error != nil {
		return nil, res.Error
	}
	var applications []*domain.JobApplication
	for _, apl := range aplsByTutor {
		var newApl = domain.JobApplication{
			Id:                apl.Id.String(),
			JobId:             apl.JobId,
			ApplicantId:       apl.ApplicantId,
			CoverLetter:       apl.CoverLetter,
			Status:            apl.Status,
			InterviewResponse: apl.InterviewResponse,
			TutorContactInfo:  apl.TutorContactInfo,
		}
		applications = append(applications, &newApl)
	}

	return applications, nil
}

func (jar JobApplicationRepo) GetApplicationsByClientRepo(cltId string) ([]*domain.JobApplication, error) {
	var aplsByClt []job_application_table
	// res := jar.db.Where("applicant_id = ?", tId).Find(&aplsByTutor)
	res := jar.db.Where("posted_by = ?", cltId).
		Order("created_at DESC").
		Find(&aplsByClt)

	if res.Error != nil {
		return nil, res.Error
	}
	var applications []*domain.JobApplication
	for _, apl := range aplsByClt {
		var newApl = domain.JobApplication{
			Id:                apl.Id.String(),
			JobId:             apl.JobId,
			ApplicantId:       apl.ApplicantId,
			CoverLetter:       apl.CoverLetter,
			Status:            apl.Status,
			InterviewResponse: apl.InterviewResponse,
			TutorContactInfo:  apl.TutorContactInfo,
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
		Order("updated_at DESC").
		Find(&aplsByStatus)

	if res.Error != nil {
		return nil, res.Error
	}

	var applications []*domain.JobApplication
	for _, apl := range aplsByStatus {
		var newApl = domain.JobApplication{
			Id:                apl.Id.String(),
			JobId:             apl.JobId,
			ApplicantId:       apl.ApplicantId,
			CoverLetter:       apl.CoverLetter,
			Status:            apl.Status,
			InterviewResponse: apl.InterviewResponse,
			TutorContactInfo:  apl.TutorContactInfo,
		}
		applications = append(applications, &newApl)
	}

	return applications, nil
}

func (jar JobApplicationRepo) UpdateMultipleStatusesRepo(ids []string, newStatus domain.ApplicationStatus) error {

	res := jar.db.Model(&job_application_table{}).Where("id IN ?", ids).Update("status", newStatus)

	return res.Error

}

func (jar JobApplicationRepo) HasAppliedRepo(jobId, tutorId string) (bool, error) {
	var count int64
	err := jar.db.Model(&job_application_table{}).
		Count(&count).
		Where("job_id = ?", jobId).
		Where("applicant_id = ?", tutorId).Error

	if err != nil {
		return false, err
	}
	return count > 0, nil
}
