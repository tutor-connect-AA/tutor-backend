package db_ports

import "github.com/tutor-connect-AA/tutor-backend/internal/application/core/domain"

type TutorNotificationDBPort interface {
	CreateTutorNotification(newNotification domain.Notification) (*domain.Notification, error)
	UpdateTutorNotificationStatus(id string) error
	GetTutorNotificationById(id string) (*domain.Notification, error)
	GetTutorNotifications() ([]*domain.Notification, error)
}
