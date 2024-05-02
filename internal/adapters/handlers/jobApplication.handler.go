package handlers

import (
	"fmt"
	"net/http"

	"github.com/tutor-connect-AA/tutor-backend/internal/application/core/domain"
	"github.com/tutor-connect-AA/tutor-backend/internal/ports/api_ports"
	"github.com/tutor-connect-AA/tutor-backend/internal/utils"
)

type JobApplicationHandler struct {
	jaSer api_ports.JobApplicationAPIPort
}

func NewJobApplicationHandler(jaSer api_ports.JobApplicationAPIPort) *JobApplicationHandler {
	return &JobApplicationHandler{
		jaSer: jaSer,
	}
}

// create the application object
// Get job id from the url
// Get the tutor detail from authorization header
// send notification to the job poster
func (jaH *JobApplicationHandler) Apply(w http.ResponseWriter, r *http.Request) {

	r.ParseMultipartForm(10 << 20)

	payload, err := utils.GetPayload(r)
	if err != nil {
		http.Error(w, "Could not get payload from token", http.StatusInternalServerError)
		return
	}

	fmt.Printf("Your cover letter is %v", r.PostForm.Get("coverLetter"))

	var newApplication = domain.JobApplication{
		JobId:       r.URL.Query().Get("id"),
		ApplicantId: payload["id"],
		CoverLetter: r.PostForm.Get("coverLetter"),
		// File: ,
	}
	ja, err := jaH.jaSer.Apply(newApplication)
	if err != nil {
		http.Error(w, "Could not create a new job application", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Successfully created a job application \n %v", ja)

}

func (jaH *JobApplicationHandler) ApplicationsByJob(w http.ResponseWriter, r *http.Request) {

	jId := r.URL.Query().Get("jobId")
	if jId == "" {
		http.Error(w, "Could not get job id", http.StatusInternalServerError)
		return
	}

	apls, err := jaH.jaSer.GetApplicationsbyJob(jId)
	if err != nil {
		http.Error(w, "Could not fetch applications", http.StatusInternalServerError)
		return
	}
	if len(apls) == 0 {
		fmt.Fprintf(w, "No applications fr this job yet")
		return
	}
	for _, apl := range apls {
		fmt.Fprint(w, *apl)
	}
}

func (jaH *JobApplicationHandler) ApplicationsByTutor(w http.ResponseWriter, r *http.Request) {
	tutorId := r.URL.Query().Get("tutorId")
	if tutorId == "" {
		http.Error(w, "Could not get job tutorId", http.StatusInternalServerError)
		return
	}
	apls, err := jaH.jaSer.GetApplicationsByTutor(tutorId)
	if err != nil {
		http.Error(w, "Could not fetch applications by tutor", http.StatusInternalServerError)
		return
	}
	if len(apls) == 0 {
		fmt.Fprintf(w, "No applications by this tutor yet")
		return
	}
	for _, apl := range apls {
		fmt.Fprint(w, *apl)
	}

}
func (jaH *JobApplicationHandler) ApplicationsByClient(w http.ResponseWriter, r *http.Request) {
	clientId := r.URL.Query().Get("clientId")
	if clientId == "" {
		http.Error(w, "Could not get job clientId", http.StatusInternalServerError)
		return
	}
	apls, err := jaH.jaSer.GetApplicationsByClient(clientId)
	if err != nil {
		http.Error(w, "Could not fetch applications by client", http.StatusInternalServerError)
		return
	}
	if len(apls) == 0 {
		fmt.Fprintf(w, "No applications for this client yet")
		return
	}
	for _, apl := range apls {
		fmt.Fprint(w, *apl)
	}

}
