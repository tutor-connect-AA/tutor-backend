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

func (cNH *ClientNotificationHandler) GetClientNotification(w http.ResponseWriter, r *http.Request) {
	ntfId := r.URL.Query().Get("ntfId")
	ntf, err := cNH.clientNotfService.GetClientNotificationById(ntfId)
	if err != nil {
		http.Error(w, "Could not get notification by id : "+err.Error(), http.StatusInternalServerError)
		return
	}

	payload, err := utils.GetPayload(r)
	if err != nil {
		http.Error(w, "Could not get payload : "+err.Error(), http.StatusInternalServerError)
		return
	}

	if ntf.OwnerId != payload["id"] {
		http.Error(w, "Not allowed to access this notification ", http.StatusForbidden)
		return
	}

	err = cNH.clientNotfService.OpenedClientNotification(ntfId) //changes the notification status to opened
	if err != nil {
		http.Error(w, "Could not change the status of notification to opened : "+err.Error(), http.StatusInternalServerError)
		// return this is commented because getting the notification is still successful.
	}

	res := Response{
		Success: true,
		Data:    ntf,
	}

	err = utils.WriteJSON(w, http.StatusOK, res, nil)
	if err != nil {
		http.Error(w, "Could not encode response to JSON : "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (cNH ClientNotificationHandler) GetClientNotifications(w http.ResponseWriter, r *http.Request) {
	payload, err := utils.GetPayload(r)
	if err != nil {
		http.Error(w, "Could not get payload : "+err.Error(), http.StatusInternalServerError)
		return
	}

	ntfs, err := cNH.clientNotfService.GetClientNotifications(payload["id"])

	if err != nil {
		http.Error(w, "Could not get notifications : "+err.Error(), http.StatusInternalServerError)
		return
	}

	res := Response{
		Success: true,
		Data:    ntfs,
	}

	err = utils.WriteJSON(w, http.StatusOK, res, nil)
	if err != nil {
		http.Error(w, "Could not encode response to JSON : "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (cNH ClientNotificationHandler) UnopenedClientNtfs(w http.ResponseWriter, r *http.Request) {

	payload, err := utils.GetPayload(r)
	if err != nil {
		http.Error(w, "Could not get payload : "+err.Error(), http.StatusInternalServerError)
		return
	}

	ntfs, err := cNH.clientNotfService.GetUnopenedClientNotifications(payload["id"])
	if err != nil {
		http.Error(w, "Could not get unopened notifications : "+err.Error(), http.StatusInternalServerError)
		return
	}

	res := Response{
		Success: true,
		Data:    ntfs,
	}

	err = utils.WriteJSON(w, http.StatusOK, res, nil)
	if err != nil {
		http.Error(w, "Could not encode response to JSON : "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (cNH ClientNotificationHandler) CountUnopenedClientNtfs(w http.ResponseWriter, r *http.Request) {
	payload, err := utils.GetPayload(r)
	if err != nil {
		http.Error(w, "Could not get payload : "+err.Error(), http.StatusInternalServerError)
		return
	}
	count, err := cNH.clientNotfService.CountUnopenedClientNotifications(payload["id"])

	if err != nil {
		http.Error(w, "Could not get unopened notifications count : "+err.Error(), http.StatusInternalServerError)
		return
	}
	res := Response{
		Success: true,
		Data:    count,
	}
	err = utils.WriteJSON(w, http.StatusOK, res, nil)
	if err != nil {
		http.Error(w, "Could not encode response to JSON : "+err.Error(), http.StatusInternalServerError)
		return
	}

}
