package api

import "github.com/tutor-connect-AA/tutor-backend/internal/application/core/domain"

func (app Application) GetClientById(id string) (*domain.Client, error) {
	client, err := app.db.GetClientByIdPort(id)

	if err != nil {
		return nil, err
	}
	return client, nil
}

func (app Application) RegisterClient(usr domain.Client) (*domain.Client, error) {
	client, err := app.db.CreateClientPort(usr)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (app Application) GetClient(id string) (*domain.Client, error) {
	client, err := app.db.GetClientByIdPort(id)

	if err != nil {
		return nil, err
	}

	return client, nil
}
