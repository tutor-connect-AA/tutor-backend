package domain

import (
	"time"
)

type Gender string

type Role string

type Education string

const (
	male   Gender = "MALE"
	female Gender = "FEMALE"
)

const (
	client Role = "CLIENT"
	tutor  Role = "TUTOR"
	admin  Role = "ADMIN"
)

const (
	preparatory Education = "PREPARATORY"
	bachelors   Education = "BACHELORS"
	masters     Education = "MASTERS"
	phd         Education = "PHD"
)

type User struct {
	id                       string
	firstName                string
	fathersName              string
	email                    string
	phoneNumber              string
	gender                   Gender
	photo                    string
	rating                   int
	bio                      string
	username                 string
	password                 string
	role                     Role
	cv                       string //link to file or blob?
	hourlyRate               int
	region                   string
	city                     string
	education                Education
	fieldOfStudy             string //could be an enum?
	educationCredential      string //link to file
	currentlyEnrolled        Education
	proofOfCurrentEnrollment string //link to file
	graduationDate           time.Time
	preferredSubjects        []string //should be an enum limited to 2
	preferredWorkLocation    string   // should this exist
	createdOn                time.Time
	updatedOn                time.Time
}
