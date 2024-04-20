package api

import (
	"github.com/tutor-connect-AA/tutor-backend/internal/application/core/domain"
	"github.com/tutor-connect-AA/tutor-backend/internal/ports/db_ports"
	"github.com/tutor-connect-AA/tutor-backend/internal/utils"
)

type ClientAPI struct {
	cr db_ports.ClientDBPort
}

func NewClientAPI(cr db_ports.ClientDBPort) *ClientAPI {
	return &ClientAPI{
		cr: cr,
	}
}

func (ca ClientAPI) GetClientById(id string) (*domain.Client, error) {
	client, err := ca.cr.GetClientByIdPort(id)

	if err != nil {
		return nil, err
	}
	return client, nil
}

func (ca ClientAPI) RegisterClient(usr domain.Client) (*domain.Client, error) {
	hashedPass, err := utils.HashPass(usr.Password)
	if err != nil {
		return nil, err
	}
	usr.Password = hashedPass
	client, err := ca.cr.CreateClientPort(usr)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (ca ClientAPI) GetListOfClients() ([]*domain.Client, error) {
	clt, err := ca.cr.GetClientsPort()

	if err != nil {
		return nil, err
	}
	return clt, nil
}

func (ca ClientAPI) UpdateClientProfile(updatedClt domain.Client) error {
	err := ca.cr.UpdateClientPort(updatedClt)

	if err != nil {
		return err
	}
	return nil
}

func (ca ClientAPI) LoginClient(username, password string) (string, error) {
	clt, err := ca.cr.GetClientByUsername(username)
	// fmt.Printf("client at client login service is %v", clt)

	if err != nil {
		return "", err
	}

	//Handle different login errors differently

	err = utils.CheckPass(clt.Password, password)
	if err != nil {
		return "", err
	}
	jwtToken, err := utils.Tokenize(clt.Id)

	if err != nil {
		return "", err
	}
	return jwtToken, nil
}
