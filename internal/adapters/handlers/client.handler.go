package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/tutor-connect-AA/tutor-backend/internal/application/core/domain"
)

type Client struct {
	Id          string      `json:"id"`
	FirstName   string      `json:"firstName"`
	FathersName string      `json:"fathersName,omitempty"` //optional
	PhoneNumber string      `json:"phoneNumber"`
	Email       string      `json:"email"`
	Username    string      `json:"username"`
	Password    string      `json:"password"`
	Photo       string      `json:"photo"`
	Role        domain.Role `json:"role"` // should role even exist?
	Rating      float32     `json:"rating"`
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

type GetClientIdReq struct {
	Id string `json:"id"`
}

func (adp Adapter) GetClientById(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get("id")
	fmt.Printf("type of id is %T", id)

	clt, err := adp.ser.GetClientById(id)
	if err != nil {
		fmt.Fprintf(w, "Could not get client by id %v", err)
		return
	}
	fmt.Fprintf(w, "Successfully got user %v", clt)
}

func (adp Adapter) GetListOfClients(w http.ResponseWriter, r *http.Request) {

	clts, err := adp.ser.GetListOfClients()

	if err != nil {
		fmt.Fprintf(w, "Could not get list of clients : %v", err)
		return
	}

	for _, clt := range clts {
		fmt.Fprintf(w, "%v \n", *clt)
	}
	// fmt.Fprintf(w, "Here are all the clients, %v", clts)

}

func (adp Adapter) UpdateClientProfile(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPut {
		http.Error(w, "Post requests only", http.StatusMethodNotAllowed)
		return
	}
	var updatedClient Client
	err := json.NewDecoder(r.Body).Decode(&updatedClient)

	if err != nil {
		log.Println(err)
		http.Error(w, "Could not decode json", http.StatusInternalServerError)
		return
	}
	fmt.Printf("%v", updatedClient)
	err = adp.ser.UpdateClientProfile(domain.Client(updatedClient))

	if err != nil {
		fmt.Fprintf(w, "Could not update client profile")
		return
	}
	fmt.Fprintf(w, "Updated Client successfully")
}

type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (adp Adapter) LoginClient(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Post requests only", http.StatusMethodNotAllowed)
		return
	}
	var loginData *LoginReq
	err := json.NewDecoder(r.Body).Decode(&loginData)
	if err != nil {
		http.Error(w, "Could not decode json", http.StatusInternalServerError)
		return
	}

	token, err := adp.ser.LoginClient(loginData.Username, loginData.Password)

	if err != nil {
		fmt.Println(err)
		fmt.Fprintf(w, "Could not login")
		return
	}

	w.Header().Set("Authorization", "Bearer "+token)

	fmt.Fprintf(w, "Successfully logged in.: %v", token)
}
