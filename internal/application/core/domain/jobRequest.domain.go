package domain

import (
	"time"
)

type JobRequestStatus string

const (
	Pending  JobRequestStatus = "PENDING"
	Accepted JobRequestStatus = "ACCEPTED"
	Rejected JobRequestStatus = "REJECTED"
)

type JobRequest struct {
	Id          string
	Status      JobRequestStatus
	SenderId    string //id of the client who sent the request
	TutorId     string // id of the tutor for whom the request is sent
	Description string // description for the job by the client
	CreatedOn   time.Time
}
