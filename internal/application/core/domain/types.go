package domain

type Gender string

type Role string

type Education string

type ApplicationStatus string

const (
	Male   Gender = "MALE"
	Female Gender = "FEMALE"
)

const (
	client Role = "CLIENT"
	tutor  Role = "TUTOR"
	admin  Role = "ADMIN"
)

const (
	Preparatory Education = "PREPARATORY"
	Bachelors   Education = "BACHELORS"
	Masters     Education = "MASTERS"
	Phd         Education = "PHD"
)
const (
	PENDING     ApplicationStatus = "PENDING"
	SHORTLISTED ApplicationStatus = "SHORTLISTED"
	// INTERESTED  ApplicationStatus = "INTERESTED"
	HIRED ApplicationStatus = "HIRED"
)
