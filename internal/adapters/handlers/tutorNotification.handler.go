package handlers

import (
	"net/http"

	"github.com/tutor-connect-AA/tutor-backend/internal/ports/api_ports"
	"github.com/tutor-connect-AA/tutor-backend/internal/utils"
)

type TutorNotificationHandler struct {
	tutorNotfService api_ports.TutorNotificationAPIPort
}

func NewTutorNotificationHandler(tutorNotfService api_ports.TutorNotificationAPIPort) *TutorNotificationHandler {
	return &TutorNotificationHandler{
		tutorNotfService: tutorNotfService,
	}
}

func (tNH TutorNotificationHandler) GetTutorNotification(w http.ResponseWriter, r *http.Request) {
	ntfId := r.URL.Query().Get("ntfId")
	ntf, err := tNH.tutorNotfService.GetTutorNotificationById(ntfId)
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
		http.Error(w, "Not allowed to access this notification : "+err.Error(), http.StatusForbidden)
		return
	}

	err = tNH.tutorNotfService.OpenedTutorNotification(ntfId) //changes the notification status to opened
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

func (tNH TutorNotificationHandler) GetTutorNotifications(w http.ResponseWriter, r *http.Request) {

	payload, err := utils.GetPayload(r)
	if err != nil {
		http.Error(w, "Could not get payload : "+err.Error(), http.StatusInternalServerError)
		return
	}

	ntfs, err := tNH.tutorNotfService.GetTutorNotifications(payload["id"])

	if err != nil {
		http.Error(w, "Could not get notifications :"+err.Error(), http.StatusInternalServerError)
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

func (tNH TutorNotificationHandler) UnopenedTutorNtfs(w http.ResponseWriter, r *http.Request) {

	payload, err := utils.GetPayload(r)

	if err != nil {
		http.Error(w, "could not get payload", http.StatusInternalServerError)
		return
	}

	ntfs, err := tNH.tutorNotfService.GetUnopenedTutorNotifications(payload["id"])
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

func (tNH TutorNotificationHandler) CountUnopenedTutorNtfs(w http.ResponseWriter, r *http.Request) {

	payload, err := utils.GetPayload(r)
	if err != nil {
		http.Error(w, "Could not get payload : "+err.Error(), http.StatusInternalServerError)
		return
	}

	count, err := tNH.tutorNotfService.CountUnopenedTutorNotifications(payload["id"])

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
