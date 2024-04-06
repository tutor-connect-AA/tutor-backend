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

func (adp Adapter) GetClient(id string) (*domain.Client, error) {
	var clientEntity Client
	clt := adp.db.First(&clientEntity, id)

	client := &domain.Client{
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

func (adp Adapter) CreateClient(clt Client) (*domain.Client, error) {
	res := adp.db.Create(&clt)

	if res.Error != nil {
		return nil, res.Error
	}

	return &domain.Client{
		Id:          clt.id,
		FirstName:   clt.firstName,
		FathersName: clt.fathersName,
		PhoneNumber: clt.phoneNumber,
		Email:       clt.email,
		Photo:       clt.photo,
		Role:        domain.Role(clt.role),
		Rating:      clt.rating,
	}, nil

}
