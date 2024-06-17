package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/tutor-connect-AA/tutor-backend/internal/application/core/domain"
	"github.com/tutor-connect-AA/tutor-backend/internal/ports/api_ports"
	"github.com/tutor-connect-AA/tutor-backend/internal/utils"
)

type JobRequestHandler struct {
	jrS     api_ports.JobRequestAPIPort
	clS     api_ports.ClientAPIPort
	tS      api_ports.TutorAPIPort
	tNtfSer api_ports.TutorNotificationAPIPort
	cNtfSer api_ports.ClientNotificationAPIPort
}

func NewJobRequestHandler(jrS api_ports.JobRequestAPIPort, clS api_ports.ClientAPIPort, tS api_ports.TutorAPIPort, tNtfSer api_ports.TutorNotificationAPIPort, cNtfSer api_ports.ClientNotificationAPIPort) *JobRequestHandler {
	return &JobRequestHandler{
		jrS:     jrS,
		clS:     clS,
		tS:      tS,
		tNtfSer: tNtfSer,
		cNtfSer: cNtfSer,
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
		http.Error(w, "Unable to create a job request : "+err.Error(), http.StatusInternalServerError)
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
		http.Error(w, "Could not create job request : "+err.Error(), http.StatusInternalServerError)
		return
	}

	data := Response{
		Success: true,
		Data:    jr,
	}
	err = utils.WriteJSON(w, http.StatusOK, data, nil)
	if err != nil {
		http.Error(w, "Could not encode response to json : "+err.Error(), http.StatusInternalServerError)
		return
	}

	ntfLink := fmt.Sprintf("https://tutor-backend-schs.onrender.com/job-request/single?jrId=%v", newRequest.Id)
	message := fmt.Sprintf("You have been requested for a job service. %v", ntfLink)

	ntf := domain.Notification{
		OwnerId: newRequest.TutorId,
		Message: message,
	}

	_, err = jrH.tNtfSer.CreateTutorNotification(ntf)
	if err != nil {
		http.Error(w, "Could not send notification : "+err.Error(), http.StatusInternalServerError)
		return
	}

}

func (jrH JobRequestHandler) GetJobRequest(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("jrId")
	jr, err := jrH.jrS.JobRequestById(id)
	if err != nil {
		http.Error(w, "Could not get job request by id :"+err.Error(), http.StatusInternalServerError)
		return
	}
	data := Response{
		Success: true,
		Data:    jr,
	}
	err = utils.WriteJSON(w, http.StatusOK, data, nil)
	if err != nil {
		http.Error(w, "Could not encode to json : "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (jrH JobRequestHandler) ChangeJobRequestStatus(w http.ResponseWriter, r *http.Request) {

	r.ParseMultipartForm(10 << 20)
	newJobRequestStatus := domain.JobRequestStatus(r.PostForm.Get("statusUpdate"))

	jrId := r.URL.Query().Get("jrId")

	jr, err := jrH.jrS.JobRequestById(jrId)

	if err != nil {
		http.Error(w, "Could not update job request status : "+err.Error(), http.StatusInternalServerError)
		return
	}

	payload, err := utils.GetPayload(r)

	if err != nil {
		http.Error(w, "Could not update job request status : "+err.Error(), http.StatusInternalServerError)
		return
	}

	if payload["id"] != jr.TutorId {
		http.Error(w, "Not allowed to make this change", http.StatusForbidden)
		return
	}
	updatedJR := domain.JobRequest{
		Status: newJobRequestStatus,
	}

	err = jrH.jrS.UpdateJobRequest(jrId, updatedJR)

	if err != nil {
		http.Error(w, "Could not update job request status : "+err.Error(), http.StatusInternalServerError)
		return
	}

	res := Response{
		Success: true,
		Data:    "Successfully updated job request status",
	}

	err = utils.WriteJSON(w, http.StatusOK, res, nil)

	if err != nil {
		http.Error(w, "Could not encode json", http.StatusInternalServerError)
		return
	}

	ntfLink := fmt.Sprintf("https://tutor-backend-schs.onrender.com/job-request/single?jrId=%v", jr.Id)
	message := fmt.Sprintf("The status of a request you were involved in just changed to %v.  %v", newJobRequestStatus, ntfLink)

	ntf := domain.Notification{
		OwnerId: jr.ClientId,
		Message: message,
	}
	_, err = jrH.cNtfSer.CreateClientNotification(ntf)

	if err != nil {
		http.Error(w, "Could not send notification : "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (jrH JobRequestHandler) HireFromRequest(w http.ResponseWriter, r *http.Request) {
	req_id := r.URL.Query().Get("reqId")

	jReq, err := jrH.jrS.JobRequestById(req_id)

	if err != nil {
		http.Error(w, "Could not get request : "+err.Error(), http.StatusInternalServerError)
		return
	}

	if jReq.Status != domain.ACCEPTED {
		http.Error(w, "The tutor has not accepted your job request yet to be hired", http.StatusForbidden)
		return
	}

	payload, err := utils.GetPayload(r)
	if err != nil {
		http.Error(w, "Error getting payload from JWT : "+err.Error(), http.StatusInternalServerError)
		return
	}

	clientInfo, err := jrH.clS.GetClientById(payload["id"])
	if err != nil {
		http.Error(w, "Could not get client info : "+err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Print("Client is : ", clientInfo)

	tx_ref := utils.RandomString(20)

	// change url
	return_url_actual := fmt.Sprintf(`https://tutor-backend-schs.onrender.com/jobRequest/verifyHire?txRef=%s&reqId=%s`, url.QueryEscape(tx_ref), url.QueryEscape(req_id)) //to be used later when deployed(b.v of verification error in url from Chapa )

	fmt.Println("Actual return url is :", return_url_actual)

	checkoutURL, err := utils.DoPayment("mahider3991@gmail.com", tx_ref, return_url_actual, 100)

	fmt.Println("Checkout URL is : ", checkoutURL)
	fmt.Println("redirected to : ", return_url_actual)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Payment redirection failed : "+err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, checkoutURL, http.StatusSeeOther)
}

func (jrH JobRequestHandler) VerifyHireFromRequest(w http.ResponseWriter, r *http.Request) {
	tx_ref := r.URL.Query().Get("txRef")

	rawQuery := r.URL.RawQuery
	fmt.Println("Raw query is:", rawQuery)

	// Replace &amp; with & in the raw query string
	fixedQuery := strings.ReplaceAll(rawQuery, "&amp;", "&")
	fmt.Println("Fixed query is:", fixedQuery)

	// Parse the fixed query string
	params, err := url.ParseQuery(fixedQuery)
	if err != nil {
		http.Error(w, "Could not parse query parameters : "+err.Error(), http.StatusInternalServerError)
		return
	}

	req_id := params.Get("reqId")
	verResult, err := utils.VerifyPayment(tx_ref)
	if err != nil {
		http.Error(w, "Could not verify payment : "+err.Error(), http.StatusInternalServerError)
		return
	}

	var jsonBody map[string]interface{}
	err = json.Unmarshal([]byte(verResult), &jsonBody)
	if err != nil {
		http.Error(w, "Could not unmarshal json : "+err.Error(), http.StatusInternalServerError)
		fmt.Println("Error unmarshalling json", err)
		return
	}
	data := jsonBody["data"]
	if data == nil {
		fmt.Fprintf(w, "Payment verification failed")
		return
	}
	// updatedApp := domain.JobRequest{
	// 	Status: domain.PAID,
	// 	TxRef:  tx_ref,
	// }

	request, err := jrH.jrS.JobRequestById(req_id)
	if err != nil {
		http.Error(w, "Could not get application by id : "+err.Error(), http.StatusInternalServerError)
		return
	}
	tutorId := request.TutorId
	applicantInfo, err := jrH.tS.GetTutorById(tutorId)
	if err != nil {
		http.Error(w, "Could not get tutor information : "+err.Error(), http.StatusInternalServerError)
		return
	}

	tutorContactInfo := map[string]string{
		"phoneNumber": applicantInfo.PhoneNumber,
		"email":       applicantInfo.Email,
	}
	tutorContact := fmt.Sprintf("email : %v , phoneNumber : %v", applicantInfo.Email, applicantInfo.PhoneNumber)

	updatedReq := domain.JobRequest{
		Status:           domain.PAID,
		TutorContactInfo: tutorContact,
		TxRef:            tx_ref,
	}

	err = jrH.jrS.UpdateJobRequest(req_id, updatedReq)
	if err != nil {
		http.Error(w, "Could not update application status : "+err.Error(), http.StatusInternalServerError)
		return
	}

	res := Response{
		Success: true,
		Data:    tutorContactInfo,
	}
	err = utils.WriteJSON(w, http.StatusOK, res, nil)
	if err != nil {
		fmt.Printf("Could not encode to json %v", err)
		http.Error(w, "JSON encoding failed", http.StatusInternalServerError)
		return
	}

	tNtfLink := fmt.Sprintf("https://tutor-backend-schs.onrender.com/job-request/single?jrId=%v", request.Id)
	tMsg := fmt.Sprintf("You have been hired. %v", tNtfLink)

	tNtf := domain.Notification{
		Message: tMsg,
		OwnerId: applicantInfo.Id,
	}
	_, err = jrH.tNtfSer.CreateTutorNotification(tNtf)
	if err != nil {
		http.Error(w, "Could not add notification for tutor :"+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, tNtfLink, http.StatusSeeOther)
}
