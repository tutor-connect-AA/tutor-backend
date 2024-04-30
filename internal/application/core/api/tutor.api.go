package api

import (
	"github.com/tutor-connect-AA/tutor-backend/internal/application/core/domain"
	"github.com/tutor-connect-AA/tutor-backend/internal/ports/db_ports"
	"github.com/tutor-connect-AA/tutor-backend/internal/utils"
)

type TutorService struct {
	tutorRepo db_ports.TutorDBPort
}

func NewTutorAPI(tr db_ports.TutorDBPort) *TutorService {
	return &TutorService{
		tutorRepo: tr,
	}
}

// RegisterTutor(tutor *domain.Tutor) (*domain.Tutor, error)
// LoginTutor(email, password string) (domain.Tutor, error)
// GetTutor(id string) (*domain.Tutor, error)
func (ts *TutorService) RegisterTutor(tutor *domain.Tutor) (*domain.Tutor, error) {
	tutor, err := ts.tutorRepo.CreateTutorRepo(tutor)
	if err != nil {
		return nil, err
	}
	return tutor, nil
}

func (ts *TutorService) GetTutor(id string) (*domain.Tutor, error) {
	tutor, err := ts.tutorRepo.GetTutorByIdRepo(id)
	if err != nil {
		return nil, err
	}
	return tutor, nil
}

func (ts *TutorService) LoginTutor(email, password string) (*domain.Tutor, error) {
	tutor, err := ts.tutorRepo.GetTutorByEmail(email)
	if err != nil {
		return nil, err
	}
	if checkPass := utils.CheckPass(tutor.Password, password); checkPass != nil {
		return nil, err
	}

	return tutor, nil
}
