package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/tutor-connect-AA/tutor-backend/internal/application/core/domain"
)

type RegisterClient struct {
	FirstName   string      `json: "firstName"`
	FathersName string      `json: "fathersName"` //optional
	PhoneNumber string      `json:"phoneNumber"`
	Email       string      `json: "email"`
	Username    string      `json: "username"`
	Password    string      `json: "password`
	Photo       string      `json: "photo"`
	Role        domain.Role `json: "role"` // should role even exist?
	Rating      float32     `json: "rating"`
}

func (adp Adapter) Register(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Post requests only", http.StatusMethodNotAllowed)
		return
	}
	var newClient domain.Client
	err := json.NewDecoder(r.Body).Decode(&newClient)

	if err != nil {
		log.Println(err)
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
