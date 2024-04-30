package domain

import (
	"time"
)

type Tutor struct {
	Id                  string
	FirstName           string
	FathersName         string
	Email               string
	PhoneNumber         string
	Gender              Gender
	Photo               string
	Rating              float32
	Bio                 string
	Username            string
	Password            string
	Role                Role
	CV                  string //link to file or blob?
	HourlyRate          float32
	Region              string
	City                string
	Education           Education
	FieldOfStudy        string //could be an enum?
	EducationCredential string //link to file
	CurrentlyEnrolled   Education
	// ProofOfCurrentEnrollment string //link to file
	GraduationDate        time.Time
	PreferredSubjects     string //should be an enum limited to 2
	PreferredWorkLocation string // should this exist
	CreatedOn             time.Time
	UpdatedOn             time.Time
}
