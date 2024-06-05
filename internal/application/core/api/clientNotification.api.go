package api

import (
	"github.com/tutor-connect-AA/tutor-backend/internal/application/core/domain"
	"github.com/tutor-connect-AA/tutor-backend/internal/ports/db_ports"
)

type ClientNotificationAPI struct {
	clientNotificationRepo db_ports.ClientNotificationDBPort
}

func NewClientNotificationAPI(clientNotificationRepo db_ports.ClientNotificationDBPort) *ClientNotificationAPI {
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
func (cnR ClientNotificationAPI) OpenedClientNotification(id string) error {

	return cnR.clientNotificationRepo.UpdateClientNotificationStatus(id)
}

func (cnR ClientNotificationAPI) GetClientNotificationById(id string) (*domain.Notification, error) {
	cntf, err := cnR.clientNotificationRepo.GetClientNotificationById(id)
	if err != nil {
		return nil, err
	}
	return cntf, nil
}
func (cnR ClientNotificationAPI) GetClientNotifications() ([]*domain.Notification, error) {
	cNtfs, err := cnR.clientNotificationRepo.GetClientNotifications()
	if err != nil {
		return nil, err
	}
	return cNtfs, nil
}

// CreateClientNotification(newNotification domain.Notification) (*domain.Notification, error)
// OpenedClientNotification(id string) error
