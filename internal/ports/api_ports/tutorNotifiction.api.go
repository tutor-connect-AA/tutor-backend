package api_ports

import "github.com/tutor-connect-AA/tutor-backend/internal/application/core/domain"

type TutorNotificationAPIPort interface {
	CreateTutorNotification(newNotification domain.Notification) (*domain.Notification, error)
	OpenedTutorNotification(id string) error
	GetTutorNotificationById(id string) (*domain.Notification, error)
	GetTutorNotifications() ([]*domain.Notification, error)
	GetUnopenedTutorNotifications() ([]*domain.Notification, error)
	CountUnopenedTutorNotifications() (*int64, error)
}
