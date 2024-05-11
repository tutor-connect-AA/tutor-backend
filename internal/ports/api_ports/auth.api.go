package api_ports

import "github.com/tutor-connect-AA/tutor-backend/internal/application/core/domain"

type AuthAPIPort interface {
	CreateAuth(domain.Auth) error
	GetAuthByUsername(username string) (*domain.Auth, error)
}
