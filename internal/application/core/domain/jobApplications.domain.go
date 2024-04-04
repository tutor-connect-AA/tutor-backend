package domain

import (
	"time"
)

type ApplicationStatus string

type JobApplication struct {
	id          string
	jobId       string //id of job
	applicantId string //id of the applicant
	coverLetter string // should this be a link to the file or the text itself
	file        string //video, other documents
	createdOn   time.Time
	updatedOn   time.Time
}
