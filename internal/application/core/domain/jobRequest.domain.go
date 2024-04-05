package domain

import (
	"time"
)

type JobRequestStatus string

const (
	pending  JobRequestStatus = "PENDING"
	accepted JobRequestStatus = "ACCEPTED"
	rejected JobRequestStatus = "REJECTED"
)

type JobRequest struct {
	id          string
	status      JobRequestStatus
	senderId    string //id of the client who sent the request
	tutorId     string // id of the tutor for whom the request is sent
	description string // description for the job by the client
	createdOn   time.Time
}
