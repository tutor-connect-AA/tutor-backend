package domain

import (
	"time"
)

type JobStatus string

const (
	open   JobStatus = "OPEN"
	closed JobStatus = "CLOSED"
)

type Job struct {
	id                 string
	title              string
	description        string
	postedBy           string
	postedOn           time.Time
	deadline           time.Time
	region             string
	city               string
	quantity           int //number of tutors needed
	numberOfApplicants int
	gradeOfStudents    []int
	minimumEducation   Education
	preferredGender    Gender
	contactHourPerWeek int
	status             JobStatus
	applications       []string // array of ids of job applications for this job
	hourlyRateMin      int
	hourlyRateMax      int
	createdOn          time.Time
}
