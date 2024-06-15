package handlers

import (
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"

	"github.com/tutor-connect-AA/tutor-backend/internal/application/core/domain"
	"github.com/tutor-connect-AA/tutor-backend/internal/ports/api_ports"
	"github.com/tutor-connect-AA/tutor-backend/internal/utils"
)

type HiringHandler struct {
	jaS       api_ports.JobApplicationAPIPort
	clS       api_ports.ClientAPIPort
	tutSer    api_ports.TutorAPIPort
	tutNtfSer api_ports.TutorNotificationAPIPort
	cltNtfSer api_ports.ClientNotificationAPIPort
	jbSer     api_ports.JobAPIPort
}

func NewHiringHandler(jaS api_ports.JobApplicationAPIPort, clS api_ports.ClientAPIPort,
	tutSer api_ports.TutorAPIPort, tutNtfSer api_ports.TutorNotificationAPIPort,
	cltNtfSer api_ports.ClientNotificationAPIPort, jbSer api_ports.JobAPIPort) *HiringHandler {
	return &HiringHandler{
		jaS:       jaS,
		clS:       clS,
		tutSer:    tutSer,
		tutNtfSer: tutNtfSer,
		cltNtfSer: cltNtfSer,
		jbSer:     jbSer,
	}
}

func (hH *HiringHandler) Hire(w http.ResponseWriter, r *http.Request) {

	payload, err := utils.GetPayload(r)
	if err != nil {
		http.Error(w, "Error getting payload from JWT : "+err.Error(), http.StatusInternalServerError)
		return
	}

	// fmt.Print("Client id is : ", payload["id"])

	clientInfo, err := hH.clS.GetClientById(payload["id"])
	if err != nil {
		http.Error(w, "Could not get client info : "+err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Print("Client is : ", clientInfo)
	app_id := r.URL.Query().Get("appId")

	tx_ref := utils.RandomString(20)
	return_url := "https://www.google.com"
	return_url_actual := fmt.Sprintf(`http://localhost:8080/hiring/verifyHire?txRef=%s&appId=%s`, url.QueryEscape(tx_ref), url.QueryEscape(app_id)) //to be used later when deployed(b.v of verification error in url from Chapa )

	fmt.Println("Actual return url is :", return_url_actual)
	fmt.Printf("return url at verify hire is: %v", return_url)

	checkoutURL, err := utils.DoPayment("mahider3991@gmail.com", tx_ref, return_url_actual, 100)

	fmt.Printf("Checkout URL is :%v", checkoutURL)
	fmt.Println("redirected to : ", return_url_actual)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Payment redirection failed : "+err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, checkoutURL, http.StatusSeeOther)
}

func (hH *HiringHandler) VerifyHire(w http.ResponseWriter, r *http.Request) {
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

	app_id := params.Get("appId")
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
	data := jsonBody["data"] //Might want to change this to checking by success rather than data
	if data == nil {
		fmt.Fprintf(w, "Payment verification failed")
		return
	}
	updatedApp := domain.JobApplication{
		Status: domain.HIRED,
		TxRef:  tx_ref,
	}
	err = hH.jaS.UpdateApplication(app_id, updatedApp)
	if err != nil {
		http.Error(w, "Could not update application status : "+err.Error(), http.StatusInternalServerError)
		return
	}

	appl, err := hH.jaS.GetApplicationById(app_id)
	if err != nil {
		http.Error(w, "Could not get application by id : "+err.Error(), http.StatusInternalServerError)
		return
	}
	applicantId := appl.ApplicantId
	applicantInfo, err := hH.tutSer.GetTutorById(applicantId)

	if err != nil {
		http.Error(w, "Could not get tutor information : "+err.Error(), http.StatusInternalServerError)
		return
	}

	tutorContactInfo := map[string]string{
		"phoneNumber": applicantInfo.PhoneNumber,
		"email":       applicantInfo.Email,
	}
	res := Response{
		Success: true,
		Data:    tutorContactInfo,
	}

	link := fmt.Sprintf("http://localhost:8080?jobApplication/single?id=%v", app_id)
	message := fmt.Sprint("You have been hired.", link)

	hiredNtf := domain.Notification{
		OwnerId: applicantId,
		Message: message,
	}

	_, err = hH.tutNtfSer.CreateTutorNotification(hiredNtf)
	if err != nil {
		http.Error(w, "Could not create hiring notification : "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = utils.WriteJSON(w, http.StatusOK, res, nil)
	if err != nil {
		fmt.Printf("Could not encode to json %v", err)
		http.Error(w, "JSON encoding failed : "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (hH *HiringHandler) Shortlist(w http.ResponseWriter, r *http.Request) {

	r.ParseMultipartForm(10 << 20)

	questions := r.PostForm.Get("questions")

	applicationId := r.URL.Query().Get("id")

	addedQuestions := &domain.JobApplication{
		InterviewQuestions: questions,
		Status:             domain.SHORTLISTED,
	}

	err := hH.jaS.UpdateApplication(applicationId, *addedQuestions)

	if err != nil {
		http.Error(w, "Could not shortlist applicant", http.StatusInternalServerError)
		return
	}

	appDetail, err := hH.jaS.GetApplicationById(applicationId)

	if err != nil {
		http.Error(w, "Could not get notification by Id : "+err.Error(), http.StatusInternalServerError)
		return
	}

	ntfLink := fmt.Sprintf("http://localhost:8080/jobApplication/single?id=%v", appDetail.Id)
	message := fmt.Sprintf("You just got shortlisted for an interview. %v", ntfLink)
	shortlistedNtf := domain.Notification{
		OwnerId: appDetail.ApplicantId,
		Message: message,
	}
	_, err = hH.tutNtfSer.CreateTutorNotification(shortlistedNtf)
	if err != nil {
		http.Error(w, "Could not create shortlist notification for tutor : "+err.Error(), http.StatusInternalServerError)
		return
	}
	res := Response{
		Success: true,
		Data:    "Applicant successfully shortlisted",
	}

	err = utils.WriteJSON(w, http.StatusOK, res, nil)
	if err != nil {
		fmt.Printf("Could not encode to json %v", err)
		http.Error(w, "JSON encoding failed : "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (hH *HiringHandler) SendInterview(w http.ResponseWriter, r *http.Request) {

	appId := r.URL.Query().Get("appId")

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Could not parse form : "+err.Error(), http.StatusBadRequest)
		return
	}
	var videoURL string

	interviewVideoPath := r.Context().Value("interviewVideoPath")
	interviewVideo := r.Context().Value("interviewVideo")

	if interviewVideo != nil {

		videoURL, err = utils.UploadToCloudinary(interviewVideo.(multipart.File), interviewVideoPath.(string))
		fmt.Printf("CLD result is %v  and error is %v ", videoURL, err)

		if err != nil {
			videoURL = ""
			fmt.Printf("Error at upload is: %v", err)
			http.Error(w, "Could not upload interview video : "+err.Error(), http.StatusInternalServerError)
			return
		}

	}

	interviewAdded := domain.JobApplication{
		InterviewResponse: videoURL,
	}

	err = hH.jaS.UpdateApplication(appId, interviewAdded)

	if err != nil {
		http.Error(w, "Could not upload interview response : "+err.Error(), http.StatusInternalServerError)
		return
	}
	// -------------------------------------------------------------------------------------------------------------
	// The purpose of the code below is to get the id of the client who posted the job.
	// It is messy, since hiring handler not depends on a lot of other services.
	// Is there a better way to go about this?

	appDetail, err := hH.jaS.GetApplicationById(appId)
	if err != nil {
		http.Error(w, "Could not get application by id : "+err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("the job id is :", appDetail.JobId)
	fmt.Println("the app id is :", appDetail.Id)
	fmt.Print("\n The whole app is : ", appDetail)
	jobDetail, err := hH.jbSer.GetJob(appDetail.JobId)
	if err != nil {
		http.Error(w, "Could not get job by id : "+err.Error(), http.StatusInternalServerError)
		return
	}

	link := fmt.Sprintf("http://localhost:8080/jobApplication/single?id=%v", appId)
	message := fmt.Sprintf("A shortlisted applicant just replied for an interview. %v", link)

	//------------------------------------------------------------------------------------------------------------------

	interviewResponse := domain.Notification{
		OwnerId: jobDetail.Posted_By,
		Message: message,
	}

	_, err = hH.cltNtfSer.CreateClientNotification(interviewResponse)
	if err != nil {
		http.Error(w, "Could not create notification for interview response : "+err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("video url", videoURL)

	res := Response{
		Success: true,
		Data: map[string]string{
			"interviewLink": videoURL,
		},
	}

	err = utils.WriteJSON(w, http.StatusOK, res, nil)
	if err != nil {
		fmt.Printf("Could not encode to json %v", err)
		http.Error(w, "JSON encoding failed : "+err.Error(), http.StatusInternalServerError)
		return
	}

}
