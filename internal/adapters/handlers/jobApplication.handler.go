package handlers

import (
	"fmt"
	"net/http"

	"github.com/tutor-connect-AA/tutor-backend/internal/application/core/domain"
	"github.com/tutor-connect-AA/tutor-backend/internal/ports/api_ports"
	"github.com/tutor-connect-AA/tutor-backend/internal/utils"
)

type JobApplicationHandler struct {
	jaSer   api_ports.JobApplicationAPIPort
	tutSer  api_ports.TutorAPIPort
	cNtfSer api_ports.ClientNotificationAPIPort
	jobSer  api_ports.JobAPIPort
}

func NewJobApplicationHandler(jaSer api_ports.JobApplicationAPIPort, tutSer api_ports.TutorAPIPort, cNtfSer api_ports.ClientNotificationAPIPort, jobSer api_ports.JobAPIPort) *JobApplicationHandler {
	return &JobApplicationHandler{
		jaSer:   jaSer,
		tutSer:  tutSer,
		cNtfSer: cNtfSer,
		jobSer:  jobSer,
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
		http.Error(w, "Could not get payload from token : "+err.Error(), http.StatusInternalServerError)
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
		http.Error(w, "Could not create a new job application : "+err.Error(), http.StatusInternalServerError)
		return
	}

	newAppLink := fmt.Sprintf("http://localhost:8080/jobApplication/single?id=%v", ja.Id)
	message := fmt.Sprintf("Some one just applied for the job you posted. Click here to view %v", newAppLink)

	job, err := jaH.jobSer.GetJob(r.URL.Query().Get("id"))

	if err != nil {
		http.Error(w, "Could not get job from the application : "+err.Error(), http.StatusInternalServerError)
		return
	}

	newJobAppNotf := domain.Notification{
		Message: message,
		OwnerId: job.Posted_By,
	}
	_, err = jaH.cNtfSer.CreateClientNotification(newJobAppNotf)

	if err != nil {
		http.Error(w, "Could not create notification : "+err.Error(), http.StatusInternalServerError)
		return
	}

	res := Response{
		Success: true,
		Data:    ja,
	}
	err = utils.WriteJSON(w, http.StatusOK, res, nil)
	if err != nil {
		fmt.Printf("Could not encode to json %v", err)
		http.Error(w, "JSON encoding failed : "+err.Error(), http.StatusInternalServerError)
		return
	}

}

func (jaH *JobApplicationHandler) GetApplicationById(w http.ResponseWriter, r *http.Request) {

	appId := r.URL.Query().Get("id")

	application, err := jaH.jaSer.GetApplicationById(appId)

	if err != nil {
		http.Error(w, "Could not get application : "+err.Error(), http.StatusInternalServerError)
		return
	}
	applicant, err := jaH.tutSer.GetTutorById(application.ApplicantId)

	if err != nil {
		http.Error(w, "Could not get application : "+err.Error(), http.StatusInternalServerError)
		return
	}

	applicantFullName := fmt.Sprintf("%v %v", applicant.FirstName, applicant.FathersName)

	applicantData := map[string]interface{}{
		"fullName":                  applicantFullName,
		"gender":                    applicant.Gender,
		"photo":                     applicant.Photo,
		"rating":                    applicant.Rating,
		"hourlyRate":                applicant.HourlyRate,
		"highestCompletedEducation": applicant.Education,
		"fieldOfStudy":              applicant.FieldOfStudy,
		"cv":                        applicant.CV,
		"city":                      applicant.City,
		"region":                    applicant.Region,
		"preferredLocation":         applicant.PreferredWorkLocation,
	}
	data := map[string]interface{}{
		"coverLetter":        application.CoverLetter,
		"interviewQuestions": application.InterviewQuestions,
		"interviewResponse":  application.InterviewResponse,
		"applicant":          applicantData,
	}

	res := Response{
		Success: true,
		Data:    data,
	}

	err = utils.WriteJSON(w, http.StatusOK, res, nil)
	if err != nil {
		fmt.Printf("Could not encode to json %v", err)
		http.Error(w, "JSON encoding failed : "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (jaH *JobApplicationHandler) ApplicationsByJob(w http.ResponseWriter, r *http.Request) {

	jId := r.URL.Query().Get("jobId")
	if jId == "" {
		http.Error(w, "Could not get job id", http.StatusInternalServerError)
		return
	}

	apls, err := jaH.jaSer.GetApplicationsByJob(jId)
	if err != nil {
		http.Error(w, "Could not fetch applications : "+err.Error(), http.StatusInternalServerError)
		return
	}
	res := Response{}
	if len(apls) == 0 {
		res.Success = true
		res.Data = "No applications to be displayed yet"

		err = utils.WriteJSON(w, http.StatusNoContent, res, nil)
		if err != nil {
			fmt.Printf("Could not encode to json %v", err)
			http.Error(w, "JSON encoding failed ", http.StatusInternalServerError)
			return
		}
		return
	}
	aplList := []domain.JobApplication{}

	for _, apl := range apls {
		aplList = append(aplList, *apl)
	}
	res.Success = true
	res.Data = aplList
	err = utils.WriteJSON(w, http.StatusOK, res, nil)
	if err != nil {
		fmt.Printf("Could not encode to json %v", err)
		http.Error(w, "JSON encoding failed", http.StatusInternalServerError)
		return
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
		http.Error(w, "Could not fetch applications by tutor : "+err.Error(), http.StatusInternalServerError)
		return
	}
	if len(apls) == 0 {
		fmt.Fprintf(w, "No applications by this tutor yet")
		return
	}
	aplList := []domain.JobApplication{}

	for _, apl := range apls {
		aplList = append(aplList, *apl)
	}
	res := Response{
		Success: true,
		Data:    aplList,
	}
	err = utils.WriteJSON(w, http.StatusOK, res, nil)
	if err != nil {
		fmt.Printf("Could not encode to json %v", err)
		http.Error(w, "JSON encoding failed", http.StatusInternalServerError)
		return
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
		http.Error(w, "Could not fetch applications by client : "+err.Error(), http.StatusInternalServerError)
		return
	}
	if len(apls) == 0 {
		fmt.Fprintf(w, "No applications for this client yet")
		return
	}
	aplList := []domain.JobApplication{}

	for _, apl := range apls {
		aplList = append(aplList, *apl)
	}
	res := Response{
		Success: true,
		Data:    aplList,
	}
	err = utils.WriteJSON(w, http.StatusOK, res, nil)
	if err != nil {
		fmt.Printf("Could not encode to json %v", err)
		http.Error(w, "JSON encoding failed : "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (jaH *JobApplicationHandler) GetApplicationByStatus(w http.ResponseWriter, r *http.Request) {

	jobId := r.URL.Query().Get("jobId")
	status := domain.ApplicationStatus(r.URL.Query().Get("status"))

	apls, err := jaH.jaSer.GetApplicationsByStatus(jobId, status)
	if err != nil {
		http.Error(w, "Could not fetch applications by status : "+err.Error(), http.StatusInternalServerError)
		return
	}

	res := Response{
		Success: true,
		Data:    "No applications with of this status for the given job.",
	}

	if len(apls) == 0 {
		err = utils.WriteJSON(w, http.StatusOK, res, nil)
		if err != nil {
			fmt.Printf("Could not encode to json %v", err)
			http.Error(w, "JSON encoding failed : "+err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}
	aplList := []domain.JobApplication{}

	for _, apl := range apls {
		aplList = append(aplList, *apl)
	}
	res = Response{
		Success: true,
		Data:    aplList,
	}
	err = utils.WriteJSON(w, http.StatusOK, res, nil)
	if err != nil {
		fmt.Printf("Could not encode to json %v", err)
		http.Error(w, "JSON encoding failed : "+err.Error(), http.StatusInternalServerError)
		return
	}
}
