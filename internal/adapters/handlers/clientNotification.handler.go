package handlers

import (
	"net/http"

	"github.com/tutor-connect-AA/tutor-backend/internal/ports/api_ports"
	"github.com/tutor-connect-AA/tutor-backend/internal/utils"
)

type ClientNotificationHandler struct {
	clientNotfService api_ports.ClientNotificationAPIPort
}

func NewClientNotificationHandler(clientNotfService api_ports.ClientNotificationAPIPort) *ClientNotificationHandler {
	return &ClientNotificationHandler{
		clientNotfService: clientNotfService,
	}
}

func (tNH ClientNotificationHandler) GetClientNotification(w http.ResponseWriter, r *http.Request) {
	ntfId := r.URL.Query().Get("ntfId")
	ntf, err := tNH.clientNotfService.GetClientNotificationById(ntfId)
	if err != nil {
		http.Error(w, "Could not get notification by id", http.StatusInternalServerError)
		return
	}

	err = tNH.clientNotfService.OpenedClientNotification(ntfId) //changes the notification status to opened
	if err != nil {
		http.Error(w, "Could not change the status of notification to opened", http.StatusInternalServerError)
		// return this is commented because getting the notification is still successful.
	}

	res := Response{
		Success: true,
		Data:    ntf,
	}

	err = utils.WriteJSON(w, http.StatusOK, res, nil)
	if err != nil {
		http.Error(w, "Could not encode response to JSON", http.StatusInternalServerError)
		return
	}
}

func (cNH ClientNotificationHandler) GetClientNotifications(w http.ResponseWriter, r *http.Request) {
	ntfs, err := cNH.clientNotfService.GetClientNotifications()

	if err != nil {
		http.Error(w, "Could not get notifications", http.StatusInternalServerError)
		return
	}

	res := Response{
		Success: true,
		Data:    ntfs,
	}

	err = utils.WriteJSON(w, http.StatusOK, res, nil)
	if err != nil {
		http.Error(w, "Could not encode response to JSON", http.StatusInternalServerError)
		return
	}
}
