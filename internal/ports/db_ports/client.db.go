package db_ports

import (
	"github.com/tutor-connect-AA/tutor-backend/internal/application/core/domain"
)

type DBPort interface {
	GetClientByIdPort(id string) (*domain.Client, error)
	CreateClientPort(clt domain.Client) (*domain.Client, error)
	GetClientsPort() ([]*domain.Client, error)
	UpdateClientPort(updatedFieldsObj domain.Client) error
}
