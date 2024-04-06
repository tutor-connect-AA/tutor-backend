package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/tutor-connect-AA/tutor-backend/internal/application/core/domain"
)

type RegisterClient struct {
	firstName   string      `json: "firstName"`
	fathersName string      `json: "fathersName"` //optional
	phoneNumber string      `json:"string"`
	email       string      `json: "email"`
	photo       string      `json: "photo"`
	role        domain.Role `json: "role"` // should role even exist?
	rating      float32     `json: "rating"`
}

func (adp Adapter) Register(w http.ResponseWriter, r *http.Request) {
	var newClient domain.Client
	err := json.NewDecoder(r.Body).Decode(&newClient)

	if err != nil {
		http.Error(w, "Could not decode json", http.StatusInternalServerError)
		return
	}
	clt, err := adp.ser.RegisterClient(newClient)

	if err != nil {
		http.Error(w, "Could not register client", http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
	fmt.Fprintf(w, "Client registered successfully : %v", clt)

}
