package handlers

import (
	"fmt"
	"net/http"

	"github.com/tutor-connect-AA/tutor-backend/internal/application/core/domain"
	"github.com/tutor-connect-AA/tutor-backend/internal/ports/api_ports"
	"github.com/tutor-connect-AA/tutor-backend/internal/utils"
)

type JobRequestHandler struct {
	jrS api_ports.JobRequestAPIPort
}

func NewJobRequestHandler(jrS api_ports.JobRequestAPIPort) *JobRequestHandler {
	return &JobRequestHandler{
		jrS: jrS,
	}
}
func (jrH JobRequestHandler) RequestJob(w http.ResponseWriter, r *http.Request) {
	//get the tutor id from url query parameter
	//get client id from jwt payload
	//Default status of "Requested"

	r.ParseMultipartForm(10 << 20)
	tutorId := r.URL.Query().Get("id")

	payload, err := utils.GetPayload(r)
	if err != nil {
		http.Error(w, "Unable to create a job request", http.StatusInternalServerError)
		return
	}
	clientId := payload["id"]

	fmt.Println("Tutor id is : ", tutorId)
	fmt.Println("Client id is: ", clientId)

	status := domain.REQUESTED

	description := r.PostForm.Get("description")

	newRequest := domain.JobRequest{
		Status:      status,
		Description: description,
		ClientId:    clientId,
		TutorId:     tutorId,
	}

	jr, err := jrH.jrS.CreateJobRequest(newRequest)

	if err != nil {
		http.Error(w, "Could not create job request", http.StatusInternalServerError)
		return
	}
	data := Response{
		Success: true,
		Data:    jr,
	}
	err = utils.WriteJSON(w, http.StatusOK, data, nil)
	if err != nil {
		http.Error(w, "Could not encode response to json", http.StatusInternalServerError)
		return
	}
}

func (jrH JobRequestHandler) GetJobRequest(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	jr, err := jrH.jrS.JobRequestById(id)
	if err != nil {
		http.Error(w, "Could not get job request by id", http.StatusInternalServerError)
		return
	}
	data := Response{
		Success: true,
		Data:    jr,
	}
	err = utils.WriteJSON(w, http.StatusOK, data, nil)
	if err != nil {
		http.Error(w, "Could not encode to json", http.StatusInternalServerError)
		return
	}
}
