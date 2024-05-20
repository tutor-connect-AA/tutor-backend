package domain

type Gender string

const (
	Male   Gender = "MALE"
	Female Gender = "FEMALE"
)

type Role string

const (
	client Role = "CLIENT"
	tutor  Role = "TUTOR"
	admin  Role = "ADMIN"
)

type Education string

const (
	Preparatory Education = "PREPARATORY"
	Bachelors   Education = "BACHELORS"
	Masters     Education = "MASTERS"
	Phd         Education = "PHD"
)

type ApplicationStatus string

const (
	PENDING     ApplicationStatus = "PENDING"
	SHORTLISTED ApplicationStatus = "SHORTLISTED"
	// INTERESTED  ApplicationStatus = "INTERESTED"
	HIRED ApplicationStatus = "HIRED"
)

type JobRequestStatus string

const (
	REQUESTED JobRequestStatus = "REQUESTED"
	ACCEPTED  JobRequestStatus = "ACCEPTED"
	REJECTED  JobRequestStatus = "REJECTED"
)
