package api

import (
	"github.com/tutor-connect-AA/tutor-backend/internal/application/core/domain"
	"github.com/tutor-connect-AA/tutor-backend/internal/ports/db_ports"
)

type TutorNotificationAPI struct {
	tutorNotificationRepo db_ports.TutorNotificationDBPort
}

func NewTutorNotificationAPI(tutorNotificationRepo db_ports.TutorNotificationDBPort) *TutorNotificationAPI {
	return &TutorNotificationAPI{
		tutorNotificationRepo: tutorNotificationRepo,
	}
}

func (tnR TutorNotificationAPI) CreateTutorNotification(newNotification domain.Notification) (*domain.Notification, error) {
	notf, err := tnR.tutorNotificationRepo.CreateTutorNotification(newNotification)

	if err != nil {
		return nil, err
	}
	return notf, nil
}
func (tnR TutorNotificationAPI) OpenedTutorNotification(id string) error {

	return tnR.tutorNotificationRepo.UpdateTutorNotificationStatus(id)
}
func (tnR TutorNotificationAPI) GetTutorNotificationById(id string) (*domain.Notification, error) {
	tntf, err := tnR.tutorNotificationRepo.GetTutorNotificationById(id)
	if err != nil {
		return nil, err
	}
	return tntf, nil
}
func (tnR TutorNotificationAPI) GetTutorNotifications(ownerId string) ([]*domain.Notification, error) {
	tNtfs, err := tnR.tutorNotificationRepo.GetTutorNotifications(ownerId)
	if err != nil {
		return nil, err
	}
	return tNtfs, nil
}

func (tnR TutorNotificationAPI) GetUnopenedTutorNotifications(ownerId string) ([]*domain.Notification, error) {
	tNtfs, err := tnR.tutorNotificationRepo.GetUnopenedTutorNotifications(ownerId)
	if err != nil {
		return nil, err
	}
	return tNtfs, nil
}

func (tnR TutorNotificationAPI) CountUnopenedTutorNotifications(ownerId string) (*int64, error) {

	count, err := tnR.tutorNotificationRepo.CountUnopenedTutorNotifications(ownerId)

	if err != nil {
		return nil, err
	}
	return count, nil
}
