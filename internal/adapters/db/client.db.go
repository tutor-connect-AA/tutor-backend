package db

import (
	"github.com/tutor-connect-AA/tutor-backend/internal/application/core/domain"
	"gorm.io/gorm"
)

type Role string

const (
	client Role = "CLIENT"
	tutor  Role = "TUTOR"
	admin  Role = "ADMIN"
)

type Client struct {
	gorm.Model
	id          string
	firstName   string
	fathersName string //optional
	phoneNumber string
	email       string
	photo       string
	role        Role // should role even exist?
	rating      float32
}

func (adp Adapter) GetClient(id string) (domain.Client, error) {
	var clientEntity Client
	clt := adp.db.First(&clientEntity, id)

	client := domain.Client{
		Id:          clientEntity.id,
		FirstName:   clientEntity.firstName,
		FathersName: clientEntity.fathersName,
		PhoneNumber: clientEntity.phoneNumber,
		Email:       clientEntity.email,
		Photo:       clientEntity.photo,
		Role:        domain.Role(clientEntity.role),
		Rating:      clientEntity.rating,
	}
	return client, clt.Error
}
