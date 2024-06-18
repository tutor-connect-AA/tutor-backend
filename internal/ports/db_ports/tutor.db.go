package db_ports

import "github.com/tutor-connect-AA/tutor-backend/internal/application/core/domain"

type TutorDBPort interface {
	CreateTutorRepo(tutor *domain.Tutor) (*domain.Tutor, error)
	GetTutorByIdRepo(id string) (*domain.Tutor, error)
	GetTutorByUsername(username string) (*domain.Tutor, error)
	SearchTutorByNameRepo(name string) ([]*domain.Tutor, error)
	GetTutorsRepo(offset, limit int) ([]*domain.Tutor, error)
	FilterTutorRepo(gender domain.Gender, rating, hourlyMin, hourlyMax int, city string, education domain.Education, fieldOfStudy string) ([]*domain.Tutor, error)
	UpdateTutorRepo(updatedTutor domain.Tutor, id string) (*domain.Tutor, error)
	ApproveRatingRepo(clientId, tutorId string) (bool, error)
}
