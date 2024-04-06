package db_ports

import (
	"github.com/tutor-connect-AA/tutor-backend/internal/application/core/domain"
)

type DBPort interface {
	CreateClientPort(domain.Client) (*domain.Client, error)
	GetClientByIdPort(id string) (*domain.Client, error)
}
