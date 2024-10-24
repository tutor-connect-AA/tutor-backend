package db

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/tutor-connect-AA/tutor-backend/internal/application/core/domain"
	"gorm.io/gorm"
)

type client_table struct {
	Id            uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	First_Name    string
	Fathers_Name  string
	Phone_Number  string
	Email         string `gorm:"unique; not null"`
	Username      string `gorm:"unique; not null"`
	Password      string
	Role          domain.Role                 `gorm:"check:role IN ('CLIENT','TUTOR','ADMIN')"`
	Rating        float32                     `gorm:"column:rating;check:rating >= 0 AND rating <= 5"`
	Jobs          []job_table                 `gorm:"foreignKey:Posted_By;references:Id"`
	Job_requests  []job_request_table         `gorm:"foreignKey:ClientId;references:Id"`
	Notifications []client_notification_table `gorm:"foreignKey:OwnerId;references:Id"`
	Comments      []comment_table             `gorm:"foreignKey:Giver;references:Id"`
}

func (ur User) GetClientByIdPort(id string) (*domain.Client, error) {
	var clientEntity client_table
	clt := ur.db.Where("id = ?", id).First(&clientEntity)

	client := &domain.Client{
		FirstName:   clientEntity.First_Name,
		FathersName: clientEntity.Fathers_Name,
		PhoneNumber: clientEntity.Phone_Number,
		Email:       clientEntity.Email,
		Username:    clientEntity.Username,
		Password:    clientEntity.Password,
		Role:        domain.Role(clientEntity.Role),
		Rating:      clientEntity.Rating,
	}
	return client, clt.Error
}

func (ur *User) CreateClientPort(clt domain.Client) (*domain.Client, error) {

	fmt.Print("The new client is at repo is :", clt)
	newClient := &client_table{
		First_Name:   clt.FirstName,
		Fathers_Name: clt.FathersName,
		Phone_Number: clt.PhoneNumber,
		Email:        clt.Email,
		Username:     clt.Username,
		Password:     clt.Password,
		Role:         "CLIENT",
		Rating:       clt.Rating,
	}

	tx := ur.db.Begin()

	go func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return nil, err
	}

	cltRes := tx.Create(&newClient)
	if cltRes.Error != nil {
		tx.Rollback()
		return nil, cltRes.Error
	}

	newAuth := domain.Auth{
		Username: newClient.Username,
		Password: newClient.Password,
		Role:     newClient.Role,
	}

	_, err := ur.CreateAuthRepo(newAuth)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	clt.Id = newClient.Id.String()

	return &clt, nil

}

func (ur *User) GetClientsPort(offset, pageSize int) ([]*domain.Client, int64, error) {
	var clients []*client_table

	res := ur.db.
		Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&clients)

	if res.Error != nil {
		return nil, 0, res.Error
	}

	var count int64
	ur.db.Model(&client_table{}).Count(&count)

	var clientsReturn []*domain.Client

	for _, client := range clients {
		cltDomain := &domain.Client{
			Id:          client.Id.String(), //how to convert to string
			FirstName:   client.First_Name,
			FathersName: client.Fathers_Name,
			PhoneNumber: client.Phone_Number,
			Email:       client.Email,
			Username:    client.Username,
			Password:    client.Password,
			Role:        domain.Role(client.Role),
			Rating:      client.Rating,
		}
		clientsReturn = append(clientsReturn, cltDomain)
	}

	return clientsReturn, count, nil
}

func (ur *User) UpdateClientPort(updatedFieldsObj domain.Client) error {
	updtClt := &client_table{
		First_Name:   updatedFieldsObj.FirstName,
		Fathers_Name: updatedFieldsObj.FathersName,
		Phone_Number: updatedFieldsObj.PhoneNumber,
		Email:        updatedFieldsObj.Email,
		Username:     updatedFieldsObj.Username,
		Password:     updatedFieldsObj.Password,
		Role:         updatedFieldsObj.Role,
		Rating:       updatedFieldsObj.Rating,
	}
	// fmt.Println(updatedFieldsObj.Id)
	res := ur.db.Model(&client_table{}).Where("id=?", updatedFieldsObj.Id).Updates(updtClt)

	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (ur User) GetClientByUsername(username string) (*domain.Client, error) {
	var clientEntity *client_table

	clt := ur.db.Where("username = ?", username).First(&clientEntity)
	if clt.Error != nil {
		if clt.Error == gorm.ErrRecordNotFound {
			return nil, domain.ErrNoRecord
		} else {
			return nil, clt.Error
		}
	}

	client := &domain.Client{
		Id:          clientEntity.Id.String(), //how to convert the uuid to string?????
		FirstName:   clientEntity.First_Name,
		FathersName: clientEntity.Fathers_Name,
		PhoneNumber: clientEntity.Phone_Number,
		Username:    clientEntity.Username,
		Password:    clientEntity.Password,
		Email:       clientEntity.Email,
		Role:        clientEntity.Role,
		Rating:      clientEntity.Rating,
	}
	return client, nil
}
