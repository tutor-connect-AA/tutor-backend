package handlers

import (
	"fmt"
	"net/http"

	"github.com/tutor-connect-AA/tutor-backend/internal/application/core/domain"
	"github.com/tutor-connect-AA/tutor-backend/internal/ports/api_ports"
	"github.com/tutor-connect-AA/tutor-backend/internal/utils"
)

type AuthHandler struct {
	aS api_ports.AuthAPIPort
	cS api_ports.ClientAPIPort
	tS api_ports.TutorAPIPort
}

func NewAuthHandler(aS api_ports.AuthAPIPort, cS api_ports.ClientAPIPort, tS api_ports.TutorAPIPort) *AuthHandler {
	return &AuthHandler{
		aS: aS,
		cS: cS,
		tS: tS,
	}
}

func (aH *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)
	if r.Method != http.MethodPost {
		http.Error(w, "Post requests only", http.StatusMethodNotAllowed)
		return
	}
	username := r.PostForm.Get("username")
	password := r.PostForm.Get("password")

	usr, err := aH.aS.GetAuthByUsername(username)
	if err != nil {
		http.Error(w, "Login failed : "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = utils.CheckPass(usr.Password, password)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	userId := ""

	var token string

	if usr.Role == "CLIENT" {
		clt, err := aH.cS.GetClientByUsername(usr.Username)
		if err != nil {
			http.Error(w, "Login failed : "+err.Error(), http.StatusInternalServerError)
			return
		}
		if userId == "" {
			userId = clt.Id
		}
		token, err = utils.Tokenize(clt.Id, string(clt.Role))
		if err != nil {
			http.Error(w, "Login failed : "+err.Error(), http.StatusInternalServerError)
			return
		}

	} else if usr.Role == "TUTOR" {
		ttr, err := aH.tS.GetTutorByUsername(usr.Username)
		if err != nil {
			http.Error(w, "Login failed : "+err.Error(), http.StatusInternalServerError)
			return
		}
		if userId == "" {
			userId = ttr.Id
		}
		token, err = utils.Tokenize(ttr.Id, string(ttr.Role))
		if err != nil {
			http.Error(w, "Login failed : "+err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "Unknown error type ", http.StatusInternalServerError)
		return
	}

	// token, err := adp.ts.LoginTutor(username, password)
	// w.Header().Set("Authorization", "Bearer "+token)

	headers := http.Header{}
	headers.Set("Authorization", "Bearer "+token)

	data := map[string]interface{}{
		"userId": userId,
		"header": token,
		"role":   usr.Role,
	}

	res := Response{
		Success: true,
		Data:    data,
	}
	// res := map[string]interface{}{
	// 	"success": true,
	// 	"data":    "Successfully logged in",
	// }

	err = utils.WriteJSON(w, http.StatusOK, res, headers)

	if err != nil {
		http.Error(w, "Could not marshal to json", http.StatusInternalServerError)
		return
	}
}

func (aH *AuthHandler) ViewSelf(w http.ResponseWriter, r *http.Request) {
	payload, err := utils.GetPayload(r)
	if err != nil {
		http.Error(w, "Could not get payload : "+err.Error(), http.StatusInternalServerError)
		return
	}
	userId := payload["id"]
	userRole := payload["role"]

	fmt.Println("id", userId)
	fmt.Println("role", userRole)

	if domain.Role(userRole) == domain.ClientRole {
		clt, err := aH.cS.GetClientById(userId)
		if err != nil {
			http.Error(w, "Could not get client : "+err.Error(), http.StatusInternalServerError)
			return
		}
		res := Response{
			Success: true,
			Data:    clt,
		}

		err = utils.WriteJSON(w, http.StatusOK, res, nil)

		if err != nil {
			http.Error(w, "Could not marshal to json", http.StatusInternalServerError)
			return
		}

	} else {
		tut, err := aH.tS.GetTutorById(userId)
		if err != nil {
			http.Error(w, "Could not get tutor : "+err.Error(), http.StatusInternalServerError)
			return
		}
		res := Response{
			Success: true,
			Data:    tut,
		}

		err = utils.WriteJSON(w, http.StatusOK, res, nil)

		if err != nil {
			http.Error(w, "Could not marshal to json", http.StatusInternalServerError)
			return
		}

	}
}
