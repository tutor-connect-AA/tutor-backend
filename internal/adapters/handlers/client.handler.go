package handlers

import (
	"fmt"
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

func (adp ClientAdapter) Register(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Post requests only", http.StatusMethodNotAllowed)
		return
	}

	var newClient domain.Client

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Could not parse form : "+err.Error(), http.StatusBadRequest)
		return
	}

	newClient = domain.Client{
		FirstName:   r.PostForm.Get("firstName"),
		FathersName: r.PostForm.Get("fathersName"),
		PhoneNumber: r.PostForm.Get("phoneNumber"),
		Email:       r.PostForm.Get("email"),
		Username:    r.PostForm.Get("username"),
		Password:    r.PostForm.Get("password"),
		Role:        domain.Role("CLIENT"),
		Rating:      3,
	}
	fmt.Print("New Client at handler is : ", newClient)
	clt, err := adp.ser.RegisterClient(newClient)

	if err != nil {
		http.Error(w, "Could not register client : "+err.Error(), http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	res := Response{
		Success: true,
		Data:    clt,
	}

	err = utils.WriteJSON(w, http.StatusOK, res, nil)

	if err != nil {
		http.Error(w, "Could not send json : "+err.Error(), http.StatusInternalServerError)
		return
	}

	// fmt.Fprintf(w, "Client registered successfully : %v", clt)

}

func (adp ClientAdapter) GetClientById(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get("id")
	fmt.Printf("type of id is %T", id)

	clt, err := adp.ser.GetClientById(id)
	if err != nil {
		fmt.Fprintf(w, "Could not get client by id %v", err)
		return
	}
	res := Response{
		Success: true,
		Data:    clt,
	}
	err = utils.WriteJSON(w, http.StatusOK, res, nil)
	if err != nil {
		http.Error(w, "Could not send json : "+err.Error(), http.StatusInternalServerError)
		return
	}

	// fmt.Fprintf(w, "Successfully got user %v", clt)
}

func (adp ClientAdapter) GetListOfClients(w http.ResponseWriter, r *http.Request) {

	offsetStr := r.URL.Query().Get("offset")
	pageSizeStr := r.URL.Query().Get("pageSize")

	offset := 0
	pageSize := 10

	if offsetStr != "" {
		parsedOffset, err := strconv.Atoi(offsetStr)
		if err == nil {
			offset = parsedOffset
		}
	}

	if pageSizeStr != "" {
		parsedPageSize, err := strconv.Atoi(pageSizeStr)
		if err == nil {
			pageSize = parsedPageSize
		}
	}

	clts, count, err := adp.ser.GetListOfClients(offset, pageSize)

	if err != nil {
		fmt.Fprintf(w, "Could not get list of clients : %v", err)
		return
	}

	data := make(map[string]interface{})

	cltList := []domain.Client{}
	for _, clt := range clts {
		cltList = append(cltList, *clt)
	}
	data["items"] = cltList
	data["total"] = count
	err = utils.WriteJSON(w, http.StatusOK, data, nil)
	if err != nil {
		http.Error(w, "Could not send json : "+err.Error(), http.StatusInternalServerError)
		return
	}

}

func (adp ClientAdapter) UpdateClientProfile(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPut {
		http.Error(w, "Post requests only", http.StatusMethodNotAllowed)
		return
	}
	err := r.ParseMultipartForm(10 << 20) // 10 MB max size
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var updatedClient = domain.Client{}
	//check that the user is updating his own profile by checking against the id in the auth token

	clientId := r.URL.Query().Get("id")

	payload, err := utils.GetPayload(r)
	if err != nil {
		http.Error(w, "Could not update profile :"+err.Error(), http.StatusInternalServerError)
		return
	}

	if clientId != payload["id"] {
		http.Error(w, "Not allowed to update this profile ", http.StatusForbidden)
		return
	}

	updatedClient.Id = clientId

	if firstName := r.PostForm.Get("firstName"); firstName != "" {
		updatedClient.FirstName = firstName

	}

	fmt.Println("First name is : ", updatedClient.FirstName)
	fmt.Println("Client id  and payload id is the same : ", clientId == payload["id"])
	fmt.Printf("updatedProfile : %v", updatedClient)

	err = adp.ser.UpdateClientProfile(domain.Client(updatedClient))

	if err != nil {
		fmt.Fprintf(w, "Could not update client profile")
		return
	}
	res := Response{
		Success: true,
		Data:    "Client profile updated successfully",
	}
	utils.WriteJSON(w, http.StatusOK, res, nil)
	fmt.Fprintf(w, "Updated Client successfully")
}

// type LoginReq struct {
// 	Username string `json:"username"`
// 	Password string `json:"password"`
// }

// func (adp ClientAdapter) LoginClient(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != http.MethodPost {
// 		http.Error(w, "Post requests only", http.StatusMethodNotAllowed)
// 		return
// 	}
// 	var loginData *LoginReq
// 	err := json.NewDecoder(r.Body).Decode(&loginData)
// 	if err != nil {
// 		http.Error(w, "Could not decode json", http.StatusInternalServerError)
// 		return
// 	}

// 	token, err := adp.ser.LoginClient(loginData.Username, loginData.Password)

// 	if err != nil {
// 		fmt.Println(err)
// 		fmt.Fprintf(w, "Could not login")
// 		return
// 	}

// 	w.Header().Set("Authorization", "Bearer "+token)

// 	fmt.Fprintf(w, "Successfully logged in.: %v", token)
// }
