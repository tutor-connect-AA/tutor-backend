package db_ports

import "github.com/tutor-connect-AA/tutor-backend/internal/application/core/domain"

type ClientNotificationDBPort interface {
	CreateClientNotification(newNotification domain.Notification) (*domain.Notification, error)
	UpdateClientNotificationStatus(id string) error
	GetClientNotificationById(id string) (*domain.Notification, error)
	GetClientNotifications() ([]*domain.Notification, error)
	GetUnopenedClientNotifications() ([]*domain.Notification, error)
	CountUnopenedClientNotifications() (*int64, error)
}
