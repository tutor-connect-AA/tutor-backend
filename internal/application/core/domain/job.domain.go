package domain

import (
	"time"
)

type JobStatus string

const (
	Open   JobStatus = "OPEN"
	Closed JobStatus = "CLOSED"
)

type Job struct {
	Id                    string
	Title                 string
	Description           string
	Posted_By             string
	Posted_On             time.Time
	Deadline              time.Time
	Region                string
	City                  string
	Quantity              int    //number of tutors needed
	Grade_Of_Students     string //[]int
	Minimum_Education     Education
	Preferred_Gender      Gender
	Contact_Hour_Per_Week int
	Status                JobStatus
	Hourly_Rate_Min       int
	Hourly_Rate_Max       int
}
