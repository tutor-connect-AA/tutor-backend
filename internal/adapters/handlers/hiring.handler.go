package handlers

import (
	"encoding/json"
	"fmt"
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
	err = hH.jaS.UpdateApplicationStatus(app_id, updatedApp)
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
