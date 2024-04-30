package db_ports

import "github.com/tutor-connect-AA/tutor-backend/internal/application/core/domain"

type TutorDBPort interface {
	CreateTutorRepo(tutor *domain.Tutor) (*domain.Tutor, error)
	GetTutorByIdRepo(id string) (*domain.Tutor, error)
	GetTutorByEmail(email string) (*domain.Tutor, error)
}
