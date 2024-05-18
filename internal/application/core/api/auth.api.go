package api

import (
	"github.com/tutor-connect-AA/tutor-backend/internal/application/core/domain"
	"github.com/tutor-connect-AA/tutor-backend/internal/ports/db_ports"
)

type AuthService struct {
	authRepo db_ports.AuthDBPort
}

func NewAuthService(ar db_ports.AuthDBPort) *AuthService {
	return &AuthService{
		authRepo: ar,
	}
}

func (as AuthService) CreateAuth(auth domain.Auth) error {
	if _, err := as.authRepo.CreateAuthRepo(auth); err != nil {
		return err
	}
	return nil
}

func (as AuthService) GetAuthByUsername(username string) (*domain.Auth, error) {
	auth, err := as.authRepo.GetAuthByUsernameRepo(username)
	if err != nil {
		return nil, err
	}
	return auth, nil
}
