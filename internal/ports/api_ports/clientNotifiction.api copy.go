package api_ports

import "github.com/tutor-connect-AA/tutor-backend/internal/application/core/domain"

type ClientNotificationAPIPort interface {
	CreateClientNotification(newNotification domain.Notification) (*domain.Notification, error)
	OpenedClientNotification(id string) error
	GetClientNotificationById(id string) (*domain.Notification, error)
	GetClientNotifications(ownerId string) ([]*domain.Notification, error)
	GetUnopenedClientNotifications(ownerId string) ([]*domain.Notification, error)
	CountUnopenedClientNotifications(ownerId string) (*int64, error)
}
