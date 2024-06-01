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
		http.Error(w, "Could not get notification by id", http.StatusInternalServerError)
		return
	}

	err = tNH.tutorNotfService.OpenedTutorNotification(ntfId) //changes the notification status to opened
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

func (tNH TutorNotificationHandler) GetTutorNotifications(w http.ResponseWriter, r *http.Request) {
	ntfs, err := tNH.tutorNotfService.GetTutorNotifications()

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
