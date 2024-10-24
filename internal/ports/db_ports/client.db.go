package db_ports

import (
	"github.com/tutor-connect-AA/tutor-backend/internal/application/core/domain"
)

type ClientDBPort interface {
	GetClientByIdPort(id string) (*domain.Client, error)
	CreateClientPort(clt domain.Client) (*domain.Client, error)
	GetClientsPort(offset, pageSize int) ([]*domain.Client, int64, error)
	UpdateClientPort(updatedFieldsObj domain.Client) error
	GetClientByUsername(username string) (*domain.Client, error)
}
