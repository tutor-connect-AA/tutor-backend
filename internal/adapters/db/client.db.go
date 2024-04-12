package db

import (
	"github.com/google/uuid"
	"github.com/tutor-connect-AA/tutor-backend/internal/application/core/domain"
	"gorm.io/gorm"
)

type Role string

const (
	client Role = "CLIENT"
	tutor  Role = "TUTOR"
	admin  Role = "ADMIN"
)

type client_table struct {
	gorm.Model
	Id           uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	First_Name   string    //`gorm :"not null"`
	Fathers_Name string    //optional
	Phone_Number string    //`gorm :"not null"`
	Email        string    `gorm:"unique; not null"`
	Username     string    `gorm:"unique; not null"`
	Password     string    //`gorm :"not null"`
	Photo        string
	Role         Role    `gorm:"check:role IN ('CLIENT','TUTOR','ADMIN')"` // should role even exist?
	Rating       float32 `gorm:"column:rating;check:rating >= 0 AND rating <= 5"`
}

func (adp Adapter) GetClientByIdPort(id string) (*domain.Client, error) {
	var clientEntity client_table
	clt := adp.db.Where("id = ?", id).First(&clientEntity)
	// clt := adp.db.First(&clientEntity, id)

	client := &domain.Client{
		// Id:          clientEntity.Id, how to convert the uuid to string?????
		FirstName:   clientEntity.First_Name,
		FathersName: clientEntity.Fathers_Name,
		PhoneNumber: clientEntity.Phone_Number,
		Email:       clientEntity.Email,
		Photo:       clientEntity.Photo,
		Role:        domain.Role(clientEntity.Role),
		Rating:      clientEntity.Rating,
	}
	return client, clt.Error
}

// func (adp Adapter) CreateClientPort(clt Client) (*domain.Client, error)
func (adp Adapter) CreateClientPort(clt domain.Client) (*domain.Client, error) {
	newClient := &client_table{
		First_Name:   clt.FirstName,
		Fathers_Name: clt.FathersName,
		Phone_Number: clt.PhoneNumber,
		Email:        clt.Email,
		Username:     clt.Username,
		Password:     clt.Password,
		Photo:        clt.Photo,
		Role:         Role(clt.Role),
		Rating:       clt.Rating,
	}
	res := adp.db.Create(&newClient)

	if res.Error != nil {
		return nil, res.Error
	}

	return &clt, nil

	// return &domain.Client{
	// 	Id:          clt.id,
	// 	FirstName:   clt.firstName,
	// 	FathersName: clt.fathersName,
	// 	PhoneNumber: clt.phoneNumber,
	// 	Email:       clt.email,
	// 	Photo:       clt.photo,
	// 	Role:        domain.Role(clt.role),
	// 	Rating:      clt.rating,
	// }, nil

}

func (adp Adapter) GetClientsPort() ([]*domain.Client, error) {
	var clients []*client_table

	res := adp.db.Find(&clients)

	if res.Error != nil {
		return nil, res.Error
	}

	var clientsReturn []*domain.Client

	for _, client := range clients {
		cltDomain := &domain.Client{
			// Id: client.Id, how to convert to string
			FirstName:   client.First_Name,
			FathersName: client.Fathers_Name,
			PhoneNumber: client.Phone_Number,
			Email:       client.Email,
			Username:    client.Username,
			Password:    client.Password,
			Photo:       client.Photo,
			Role:        domain.Role(client.Role),
			Rating:      client.Rating,
		}
		clientsReturn = append(clientsReturn, cltDomain)
	}

	return clientsReturn, nil
}

func (adp Adapter) UpdateClientPort(updatedFieldsObj domain.Client) error {
	updtClt := &client_table{
		First_Name:   updatedFieldsObj.FirstName,
		Fathers_Name: updatedFieldsObj.FathersName,
		Phone_Number: updatedFieldsObj.PhoneNumber,
		Email:        updatedFieldsObj.Email,
		Username:     updatedFieldsObj.Username,
		Password:     updatedFieldsObj.Password,
		Photo:        updatedFieldsObj.Photo,
		Role:         Role(updatedFieldsObj.Role),
		Rating:       updatedFieldsObj.Rating,
	}
	// fmt.Println(updatedFieldsObj.Id)
	res := adp.db.Model(&client_table{}).Where("id=?", updatedFieldsObj.Id).Updates(updtClt)

	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (adp Adapter) GetClientByUsername(username string) (*domain.Client, error) {
	var clientEntity *client_table

	clt := adp.db.Where("username = ?", username).First(&clientEntity)
	client := &domain.Client{
		// Id:          clientEntity.Id, how to convert the uuid to string?????
		FirstName:   clientEntity.First_Name,
		FathersName: clientEntity.Fathers_Name,
		PhoneNumber: clientEntity.Phone_Number,
		Email:       clientEntity.Email,
		Photo:       clientEntity.Photo,
		Role:        domain.Role(clientEntity.Role),
		Rating:      clientEntity.Rating,
	}
	return client, clt.Error
}
