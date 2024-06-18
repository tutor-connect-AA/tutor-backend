package api_ports

import "github.com/tutor-connect-AA/tutor-backend/internal/application/core/domain"

type TutorAPIPort interface {
	RegisterTutor(tutor *domain.Tutor) (*domain.Tutor, error)
	GetTutorById(id string) (*domain.Tutor, error)
	GetTutorByUsername(username string) (*domain.Tutor, error)
	SearchTutorByName(name string) ([]*domain.Tutor, error)
	GetTutors(offset, limit int) ([]*domain.Tutor, error)
	FilterTutor(gender domain.Gender, rating, hourlyMin, hourlyMax int, city string, education domain.Education, fieldOfStudy string) ([]*domain.Tutor, error)
	UpdateTutor(updatedTutor domain.Tutor, id string) (*domain.Tutor, error)
	ApproveRating(clientId, tutorId string) (bool, error)
	// LoginTutor(username, password string) (string, error)
}
