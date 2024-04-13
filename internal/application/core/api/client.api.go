package api

import (
	"github.com/tutor-connect-AA/tutor-backend/internal/application/core/domain"
	"github.com/tutor-connect-AA/tutor-backend/internal/utils"
)

func (app Application) GetClientById(id string) (*domain.Client, error) {
	client, err := app.db.GetClientByIdPort(id)

	if err != nil {
		return nil, err
	}
	return client, nil
}

func (app Application) RegisterClient(usr domain.Client) (*domain.Client, error) {
	hashedPass, err := utils.HashPass(usr.Password)
	if err != nil {
		return nil, err
	}
	usr.Password = hashedPass
	client, err := app.db.CreateClientPort(usr)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (app Application) GetListOfClients() ([]*domain.Client, error) {
	clt, err := app.db.GetClientsPort()

	if err != nil {
		return nil, err
	}
	return clt, nil
}

func (app Application) UpdateClientProfile(updatedClt domain.Client) error {
	err := app.db.UpdateClientPort(updatedClt)

	if err != nil {
		return err
	}
	return nil
}

func (app Application) LoginClient(username, password string) (string, error) {
	clt, err := app.db.GetClientByUsername(username)

	if err != nil {
		return "", err
	}

	//Handle different login errors differently

	err = utils.CheckPass(clt.Password, password)
	if err != nil {
		return "", err
	}
	jwtToken, err := utils.Tokenize(clt.Username)

	if err != nil {
		return "", err
	}
	return jwtToken, nil
}
