package db_ports

import "github.com/tutor-connect-AA/tutor-backend/internal/application/core/domain"

type TutorNotificationDBPort interface {
	CreateTutorNotification(newNotification domain.Notification) (*domain.Notification, error)
	UpdateTutorNotificationStatus(id string) error
	GetTutorNotificationById(id string) (*domain.Notification, error)
	GetTutorNotifications(ownerId string) ([]*domain.Notification, error)
	GetUnopenedTutorNotifications(ownerId string) ([]*domain.Notification, error)
	CountUnopenedTutorNotifications(ownerId string) (*int64, error)
}
