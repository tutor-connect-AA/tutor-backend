package domain

import (
	"time"
)

type JobApplication struct {
	Id                string
	JobId             string //id of job
	ApplicantId       string //id of the applicant
	CoverLetter       string // should this be a link to the file or the text itself
	Status            ApplicationStatus
	InterviewResponse string
	TxRef             string
	TutorContactInfo  string
	// File        string //video, other documents
	CreatedOn time.Time
	UpdatedOn time.Time
}
