package handlers

import (
	"fmt"
	"net/http"

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
		http.Error(w, "Login failed", http.StatusInternalServerError)
		return
	}

	err = utils.CheckPass(usr.Password, password)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	var token string

	if usr.Role == "CLIENT" {
		clt, err := aH.cS.GetClientByUsername(usr.Username)
		if err != nil {
			http.Error(w, "Login failed", http.StatusInternalServerError)
			return
		}
		token, err = utils.Tokenize(clt.Id, string(clt.Role))
		if err != nil {
			http.Error(w, "Login failed", http.StatusInternalServerError)
			return
		}

	} else if usr.Role == "TUTOR" {
		ttr, err := aH.tS.GetTutorByUsername(usr.Username)
		if err != nil {
			http.Error(w, "Login failed", http.StatusInternalServerError)
			return
		}
		token, err = utils.Tokenize(ttr.Id, string(ttr.Role))
		if err != nil {
			http.Error(w, "Login failed", http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "Unknown error type", http.StatusInternalServerError)
		return
	}

	// token, err := adp.ts.LoginTutor(username, password)

	w.Header().Set("Authorization", "Bearer "+token)

	fmt.Fprintf(w, "Successfully logged in.: %v", token)
}
