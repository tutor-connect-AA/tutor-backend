package handlers

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"

	// "sync"
	"time"

	"github.com/tutor-connect-AA/tutor-backend/internal/application/core/domain"
	"github.com/tutor-connect-AA/tutor-backend/internal/ports/api_ports"
	"github.com/tutor-connect-AA/tutor-backend/internal/utils"
)

type TutorHandler struct {
	ts api_ports.TutorAPIPort
}

func NewTutorHandler(ts api_ports.TutorAPIPort) *TutorHandler {
	return &TutorHandler{
		ts: ts,
	}
}

func (th *TutorHandler) RegisterTutor(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20) // 10 MB max size
	fmt.Printf("Preferred subjects are %v", r.Form["prefSubjects"])
	stringPrefSubj := strings.Join(r.Form["prefSubjects"], ",")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	rating, err := strconv.ParseFloat(r.PostForm.Get("rating"), 32)
	if err != nil {
		http.Error(w, "Rating field data type mismatch", http.StatusInternalServerError)
		return
	}

	hourlyRate, err := strconv.ParseFloat(r.PostForm.Get("hourlyRate"), 32)
	if err != nil {
		http.Error(w, "Hourly rate field data type mismatch", http.StatusInternalServerError)
		return
	}

	gradDate, err := time.Parse("2006-01-02", r.PostForm.Get("graduationDate"))
	if err != nil {
		fmt.Printf("Could not parse deadline %v", err)
		http.Error(w, "Invalid date", http.StatusInternalServerError)
		return
	}

	role := domain.Role(r.PostForm.Get("role"))                                // what if conversion fails
	education := domain.Education(r.PostForm.Get("education"))                 //what if conversion fails
	currentlyEnrolled := domain.Education(r.PostForm.Get("currentlyEnrolled")) //what if conversion fails
	gender := domain.Gender(r.PostForm.Get("gender"))

	photoPath := r.Context().Value("photoPath")
	photo := r.Context().Value("photo")

	cvPath := r.Context().Value("cvPath")
	cv := r.Context().Value("cv")

	eduCredPath := r.Context().Value("eduCredPath")
	eduCred := r.Context().Value("eduCred")

	var cvURL string
	var photoURL string
	var eduCredURL string

	if cv != nil {
		cvURL, err = utils.UploadToCloudinary(cv.(multipart.File), cvPath.(string))
		if err != nil {
			cvURL = ""
			http.Error(w, "Could not upload cv", http.StatusInternalServerError)
			return
		}
	}

	if photo != nil {
		photoURL, err = utils.UploadToCloudinary(photo.(multipart.File), photoPath.(string))
		fmt.Printf("CLD result is %v  and error is %v ", photoURL, err)
		if err != nil {
			photoURL = ""
			fmt.Printf("Error at upload is: %v", err)
			http.Error(w, "Could not upload photo", http.StatusInternalServerError)
			return
		}
	}

	if eduCred != nil {
		eduCredURL, err = utils.UploadToCloudinary(eduCred.(multipart.File), eduCredPath.(string))
		if err != nil {
			eduCredURL = ""
			http.Error(w, "Could not upload education credential of tutor", http.StatusInternalServerError)
			return
		}
	}

	var newTutor = &domain.Tutor{
		FirstName:           r.PostForm.Get("firstName"),
		FathersName:         r.PostForm.Get("fathersName"),
		Email:               r.PostForm.Get("email"),
		PhoneNumber:         r.PostForm.Get("phoneNumber"),
		Gender:              gender,
		Photo:               photoURL,
		Rating:              float32(rating),
		Bio:                 r.PostForm.Get("bio"),
		Username:            r.PostForm.Get("username"),
		Password:            r.PostForm.Get("password"),
		Role:                role,
		CV:                  cvURL,
		HourlyRate:          float32(hourlyRate),
		Region:              r.PostForm.Get("region"),
		City:                r.PostForm.Get("city"),
		Education:           education,
		FieldOfStudy:        r.PostForm.Get("fieldOfStudy"),
		EducationCredential: eduCredURL,
		CurrentlyEnrolled:   currentlyEnrolled,
		GraduationDate:      gradDate,
		PreferredSubjects:   stringPrefSubj,
	}
	tt, err := th.ts.RegisterTutor(newTutor)
	if err != nil {
		http.Error(w, "Could not register tutor", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Successfully registered tutor %v", tt)

}
