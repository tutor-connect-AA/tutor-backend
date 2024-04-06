package api_ports

import (
	"github.com/tutor-connect-AA/tutor-backend/internal/application/core/domain"
)

type APIPort interface {
	RegisterClient(usr domain.Client) (*domain.Client, error)
	GetClientById(id string) (*domain.Client, error)
}
