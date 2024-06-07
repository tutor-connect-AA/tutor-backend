package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/tutor-connect-AA/tutor-backend/internal/application/core/domain"
	"github.com/tutor-connect-AA/tutor-backend/internal/ports/api_ports"
	"github.com/tutor-connect-AA/tutor-backend/internal/utils"
)

type HiringHandler struct {
	jaS api_ports.JobApplicationAPIPort
	clS api_ports.ClientAPIPort
}

func NewHiringHandler(jaS api_ports.JobApplicationAPIPort, clS api_ports.ClientAPIPort) *HiringHandler {
	return &HiringHandler{
		jaS: jaS,
		clS: clS,
	}
}

func (hH *HiringHandler) Hire(w http.ResponseWriter, r *http.Request) {

	payload, err := utils.GetPayload(r)
	if err != nil {
		http.Error(w, "Error getting payload from JWT", http.StatusInternalServerError)
		return
	}

	// fmt.Print("Client id is : ", payload["id"])

	clientInfo, err := hH.clS.GetClientById(payload["id"])
	if err != nil {
		http.Error(w, "Could not get client info", http.StatusInternalServerError)
		return
	}
	fmt.Print("Client is : ", clientInfo)
	app_id := r.URL.Query().Get("appId")

	tx_ref := utils.RandomString(20)
	return_url := "https://www.google.com"
	return_url_actual := fmt.Sprintf(`http://localhost:8080/jobApplication/verifyHire?txRef=%v&appId=%v`, tx_ref, app_id) //to be used later when deployed(b.v of verification error in url from Chapa )
	fmt.Printf("return url at verify hire is: %v", return_url)

	checkoutURL, err := utils.DoPayment("mahider3991@gmail.com", tx_ref, return_url_actual, 100)

	fmt.Printf("Checkout URL is :%v", checkoutURL)
	fmt.Println("redirected to : ", return_url_actual)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Payment redirection failed", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, return_url_actual, http.StatusSeeOther)
}

func (hH *HiringHandler) VerifyHire(w http.ResponseWriter, r *http.Request) {
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
	updatedApp := domain.JobApplication{
		Status: domain.HIRED,
		TxRef:  tx_ref,
	}
	err = hH.jaS.UpdateApplication(app_id, updatedApp)
	if err != nil {
		http.Error(w, "Could not update application status", http.StatusInternalServerError)
		return
	}
	res := Response{
		Success: true,
		Data:    "Payment successful and tutor hired.",
	}
	err = utils.WriteJSON(w, http.StatusOK, res, nil)
	if err != nil {
		fmt.Printf("Could not encode to json %v", err)
		http.Error(w, "JSON encoding failed", http.StatusInternalServerError)
		return
	}
}

func (hH *HiringHandler) Shortlist(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var questions string

	body, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, "Error reading request body: %v", http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(body, &questions)
	if err != nil {
		http.Error(w, "Could not unmarshal json", http.StatusInternalServerError)
		return
	}

	applicationId := r.URL.Query().Get("id")

	addedQuestions := &domain.JobApplication{
		InterviewQuestions: questions,
		Status:             domain.SHORTLISTED,
	}

	err = hH.jaS.UpdateApplication(applicationId, *addedQuestions)

	if err != nil {
		http.Error(w, "Could not shortlist applicant", http.StatusInternalServerError)
		return
	}

	res := Response{
		Success: true,
		Data:    "Applicant successfully shortlisted",
	}

	err = utils.WriteJSON(w, http.StatusOK, res, nil)
	if err != nil {
		fmt.Printf("Could not encode to json %v", err)
		http.Error(w, "JSON encoding failed", http.StatusInternalServerError)
		return
	}
}

func (hH *HiringHandler) SendInterview(w http.ResponseWriter, r *http.Request) {

	var videoURL string

	interviewVideoPath := r.Context().Value("interviewVideoPath")
	interviewVideo := r.Context().Value("interviewVideo")

	if interviewVideo != nil {
		videoURL, err := utils.UploadToCloudinary(interviewVideo.(multipart.File), interviewVideoPath.(string))
		fmt.Printf("CLD result is %v  and error is %v ", videoURL, err)
		if err != nil {
			videoURL = ""
			fmt.Printf("Error at upload is: %v", err)
			http.Error(w, "Could not upload interview video", http.StatusInternalServerError)
			return
		}
	}

	appId := r.URL.Query().Get("appId")

	interviewAdded := &domain.JobApplication{
		InterviewResponse: videoURL,
	}

	err := hH.jaS.UpdateApplication(appId, *interviewAdded)

	if err != nil {
		http.Error(w, "Could not upload interview response", http.StatusInternalServerError)
		return
	}

	res := Response{
		Success: true,
		Data: map[string]string{
			"interviewLink": videoURL,
		},
	}

	err = utils.WriteJSON(w, http.StatusOK, res, nil)
	if err != nil {
		fmt.Printf("Could not encode to json %v", err)
		http.Error(w, "JSON encoding failed", http.StatusInternalServerError)
		return
	}

}
