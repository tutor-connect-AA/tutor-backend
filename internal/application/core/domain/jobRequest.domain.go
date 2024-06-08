package domain

import (
	"time"
)

type JobRequest struct {
	Id          string
	Description string // description for the job by the client
	Status      JobRequestStatus
	ClientId    string //id of the client who sent the request
	TutorId     string // id of the tutor for whom the request is sent
	TxRef       string
	CreatedOn   time.Time
}
