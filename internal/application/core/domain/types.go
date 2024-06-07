package domain

type Gender string

const (
	Male   Gender = "MALE"
	Female Gender = "FEMALE"
)

type Role string

const (
	ClientRole Role = "CLIENT"
	TutorRole  Role = "TUTOR"
	AdminRole  Role = "ADMIN"
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
	INTERVIEWED ApplicationStatus = "INTERVIEWED"
	HIRED       ApplicationStatus = "HIRED"
)

type JobRequestStatus string

const (
	REQUESTED  JobRequestStatus  = "REQUESTED"
	INTERESTED ApplicationStatus = "INTERESTED"
	ACCEPTED   JobRequestStatus  = "ACCEPTED"
	REJECTED   JobRequestStatus  = "REJECTED"
)
