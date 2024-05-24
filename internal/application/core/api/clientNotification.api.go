package api

import (
	"github.com/tutor-connect-AA/tutor-backend/internal/application/core/domain"
	"github.com/tutor-connect-AA/tutor-backend/internal/ports/db_ports"
)

type ClientNotificationAPI struct {
	clientNotificationRepo db_ports.ClientNotificationDBPort
}

func NewNotificationAPI(clientNotificationRepo db_ports.ClientNotificationDBPort) *ClientNotificationAPI {
	return &ClientNotificationAPI{
		clientNotificationRepo: clientNotificationRepo,
	}
}

func (cnR ClientNotificationAPI) CreateClientNotification(newNotification domain.Notification) (*domain.Notification, error) {
	notf, err := cnR.clientNotificationRepo.CreateClientNotification(newNotification)

	if err != nil {
		return nil, err
	}
	return notf, nil
}
func (cnR ClientNotificationAPI) OpenedNotification(id string) error {

	return cnR.clientNotificationRepo.UpdateClientNotificationStatus(id)
}

// CreateClientNotification(newNotification domain.Notification) (*domain.Notification, error)
// OpenedClientNotification(id string) error
