package handlers

// import (
// 	"encoding/json"
// 	"fmt"
// 	"net/http"

// 	"github.com/tutor-connect-AA/tutor-backend/internal/utils"
// )

// func PaymentVerification(w http.ResponseWriter, r *http.Request) {
// 	tx_ref := r.URL.Query().Get("tx_ref")
// 	verResult, err := utils.VerifyPayment(tx_ref)
// 	if err != nil {
// 		http.Error(w, "Could not verify payment", http.StatusInternalServerError)
// 		return
// 	}

// 	var jsonBody map[string]interface{}
// 	err = json.Unmarshal([]byte(verResult), &jsonBody)
// 	if err != nil {
// 		http.Error(w, "Could not unmarshal json", http.StatusInternalServerError)
// 		fmt.Println("Error unmarshalling json", err)
// 		return
// 	}
// 	data := jsonBody["data"]
// 	if data == nil {
// 		http.Error()
// 	}

// }
