package db_ports

import "github.com/tutor-connect-AA/tutor-backend/internal/application/core/domain"

type AuthDBPort interface {
	CreateAuthRepo(domain.Auth) (string, error)
	GetAuthByUsernameRepo(username string) (*domain.Auth, error)
}
