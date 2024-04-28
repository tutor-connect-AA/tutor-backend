package api_ports

import "github.com/tutor-connect-AA/tutor-backend/internal/application/core/domain"

type TutorAPIPort interface {
	RegisterTutor(tutor *domain.Tutor) (*domain.Tutor, error)
	LoginTutor(email, password string) (*domain.Tutor, error)
	GetTutor(id string) (*domain.Tutor, error)
}
