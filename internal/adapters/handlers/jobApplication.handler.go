package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/tutor-connect-AA/tutor-backend/internal/application/core/domain"
	"github.com/tutor-connect-AA/tutor-backend/internal/ports/api_ports"
	"github.com/tutor-connect-AA/tutor-backend/internal/utils"
)

type JobApplicationHandler struct {
	jaSer api_ports.JobApplicationAPIPort
	clSer api_ports.ClientAPIPort
}

func NewJobApplicationHandler(jaSer api_ports.JobApplicationAPIPort, clSer api_ports.ClientAPIPort) *JobApplicationHandler {
	return &JobApplicationHandler{
		jaSer: jaSer,
		clSer: clSer,
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

	apls, err := jaH.jaSer.GetApplicationsByJob(jId)
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
func (jaH *JobApplicationHandler) Hire(w http.ResponseWriter, r *http.Request) {

	// payload, err := utils.GetPayload(r)
	// if err != nil {
	// 	http.Error(w, "Error getting payload from JWT", http.StatusInternalServerError)
	// 	return
	// }

	// clientInfo, err := jaH.clSer.GetClientById(payload["id"])
	// if err != nil {
	// 	http.Error(w, "Could not get client info", http.StatusInternalServerError)
	// 	return
	// }
	app_id := r.URL.Query().Get("appId")

	tx_ref := utils.RandomString(20)
	return_url := "https://www.google.com"
	return_url_actual := fmt.Sprintf(`localhost:8080/jobApplication/verifyHire?txRef=%v&appId=%v`, tx_ref, app_id) //to be used later when deployed(b.v of verification error in url from Chapa )
	fmt.Printf("return url at verify hire is: %v", return_url)

	checkoutURL, err := utils.DoPayment( /*clientInfo.Email,*/ tx_ref, return_url, 100)

	fmt.Printf("Checkout URL is :%v", checkoutURL)
	fmt.Println("redirected to : ", return_url_actual)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Payment redirection failed", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, return_url_actual, http.StatusSeeOther)
	// applicationId := r.URL.Query().Get("id")
	// err = jaH.jaSer.UpdateApplicationStatus(applicationId, domain.HIRED)
	// if err != nil {
	// 	fmt.Println(err)
	// 	http.Error(w, "Could not perform hiring operation", http.StatusInternalServerError)
	// 	return
	// }
	// fmt.Fprint(w, "Successfully updated status of job application to HIRED")
}

func (jaH *JobApplicationHandler) VerifyHire(w http.ResponseWriter, r *http.Request) {
	tx_ref := r.URL.Query().Get("txRef")
	app_id := r.URL.Query().Get("appId")
	verResult, err := utils.VerifyPayment(tx_ref)
	if err != nil {
		http.Error(w, "Could not verify payment", http.StatusInternalServerError)
		return
	}

	var jsonBody map[string]interface{}
	err = json.Unmarshal([]byte(verResult), &jsonBody)
	if err != nil {
		http.Error(w, "Could not unmarshal json", http.StatusInternalServerError)
		fmt.Println("Error unmarshalling json", err)
		return
	}
	data := jsonBody["data"]
	if data == nil {
		fmt.Fprintf(w, "Payment verification failed")
		return
	}
	jaH.jaSer.UpdateApplicationStatus(app_id, domain.HIRED)
	fmt.Fprintf(w, "Payment successful and tutor hired.")
}
