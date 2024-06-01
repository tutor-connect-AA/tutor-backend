package api_ports

import "github.com/tutor-connect-AA/tutor-backend/internal/application/core/domain"

type TutorNotificationAPIPort interface {
	CreateTutorNotification(newNotification domain.Notification) (*domain.Notification, error)
	OpenedTutorNotification(id string) error
	GetTutorNotificationById(id string) (*domain.Notification, error)
}