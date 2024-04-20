package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/tutor-connect-AA/tutor-backend/internal/application/core/domain"
	"github.com/tutor-connect-AA/tutor-backend/internal/ports/api_ports"
	"github.com/tutor-connect-AA/tutor-backend/internal/utils"
)

type ClientAdapter struct {
	ser api_ports.ClientAPIPort
}

func NewClientHandler(ser api_ports.ClientAPIPort) *ClientAdapter {
	return &ClientAdapter{
		ser: ser,
	}
}

type Client struct {
	Id          string      `form:"id"`
	FirstName   string      `form:"firstName"`
	FathersName string      `form:"fathersName,omitempty"` //optional
	PhoneNumber string      `form:"phoneNumber"`
	Email       string      `form:"email"`
	Username    string      `form:"username"`
	Password    string      `form:"password"`
	Photo       string      `form:"photo"`
	Role        domain.Role `form:"role"` // should role even exist?
	Rating      float32     `form:"rating"`
}

func (adp ClientAdapter) Register(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Post requests only", http.StatusMethodNotAllowed)
		return
	}

	var newClient domain.Client
	contentType := r.Header.Get("Content-Type")

	if contentType == "application/json" {
		err := json.NewDecoder(r.Body).Decode(&newClient)

		if err != nil {
			log.Println(err)
			http.Error(w, "Could not decode json", http.StatusInternalServerError)
			return
		}
	}

	err := r.ParseMultipartForm(10 << 20) // 10 MB max size
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	rating, err := strconv.ParseFloat(r.PostForm.Get("rating"), 64)
	if err != nil {
		fmt.Print(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	fileName := r.Context().Value("filePath")
	file := r.Context().Value("file")
	imageUrl, err := utils.UploadToCloudinary(file.(multipart.File), fileName.(string))
	if err != nil {
		http.Error(w, "Could not upload image to Cloudinary", http.StatusInternalServerError)
		return
	}
	newClient = domain.Client{
		FirstName:   r.PostForm.Get("firstName"),
		FathersName: r.PostForm.Get("fathersName"),
		PhoneNumber: r.PostForm.Get("phoneNumber"),
		Email:       r.PostForm.Get("email"),
		Username:    r.PostForm.Get("username"),
		Password:    r.PostForm.Get("password"),
		Photo:       imageUrl,
		Role:        domain.Role(r.PostForm.Get("role")),
		Rating:      float32(rating),
	}
	// fmt.Print(newClient)
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

func (adp ClientAdapter) GetClientById(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get("id")
	fmt.Printf("type of id is %T", id)

	clt, err := adp.ser.GetClientById(id)
	if err != nil {
		fmt.Fprintf(w, "Could not get client by id %v", err)
		return
	}
	fmt.Fprintf(w, "Successfully got user %v", clt)
}

func (adp ClientAdapter) GetListOfClients(w http.ResponseWriter, r *http.Request) {

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

func (adp ClientAdapter) UpdateClientProfile(w http.ResponseWriter, r *http.Request) {

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

func (adp ClientAdapter) LoginClient(w http.ResponseWriter, r *http.Request) {
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
