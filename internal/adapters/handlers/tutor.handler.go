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
		http.Error(w, "Could not parse form : "+err.Error(), http.StatusBadRequest)
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
		http.Error(w, "Invalid date : "+err.Error(), http.StatusInternalServerError)
		return
	}
	// what if conversion fails
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
			http.Error(w, "Could not upload cv : "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if photo != nil {
		photoURL, err = utils.UploadToCloudinary(photo.(multipart.File), photoPath.(string))
		fmt.Printf("CLD result is %v  and error is %v ", photoURL, err)
		if err != nil {
			photoURL = ""
			fmt.Printf("Error at upload is: %v", err)
			http.Error(w, "Could not upload photo : "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if eduCred != nil {
		eduCredURL, err = utils.UploadToCloudinary(eduCred.(multipart.File), eduCredPath.(string))
		if err != nil {
			eduCredURL = ""
			http.Error(w, "Could not upload education credential of tutor : "+err.Error(), http.StatusInternalServerError)
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
		Rating:              3,
		RateCount:           1,
		Bio:                 r.PostForm.Get("bio"),
		Username:            r.PostForm.Get("username"),
		Password:            r.PostForm.Get("password"),
		Role:                domain.TutorRole,
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
		http.Error(w, "Could not register tutor: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := Response{
		Success: true,
		Data:    tt,
	}

	err = utils.WriteJSON(w, http.StatusOK, response, nil)
	if err != nil {
		fmt.Printf("Could not encode to json %v", err)
		http.Error(w, "JSON encoding failed : "+err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Successfully registered tutor %v", tt)
}

func (th *TutorHandler) GetTutorById(w http.ResponseWriter, r *http.Request) {
	tutId := r.URL.Query().Get("tutId")

	tutor, err := th.ts.GetTutorById(tutId)

	if err != nil {
		http.Error(w, "Could not get tutor by id : "+err.Error(), http.StatusInternalServerError)
		return
	}

	res := Response{
		Success: true,
		Data:    tutor,
	}

	err = utils.WriteJSON(w, http.StatusOK, res, nil)
	if err != nil {
		http.Error(w, "Error marshalling to json "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (th *TutorHandler) SearchTutorByName(w http.ResponseWriter, r *http.Request) {

	// r.ParseMultipartForm(10 << 20)

	// searchName := r.PostForm.Get("searchName")

	searchName := r.URL.Query().Get("searchName")

	if searchName == "" {
		http.Error(w, "Can't search with an empty name ", http.StatusBadRequest)
		return
	}

	fmt.Println("searchTerm at handler : ", searchName)

	tutors, err := th.ts.SearchTutorByName(searchName)

	if err != nil {
		http.Error(w, "Could not search tutor by name : "+err.Error(), http.StatusInternalServerError)
		return
	}
	tutList := []domain.Tutor{}

	for _, tut := range tutors {
		tutList = append(tutList, *tut)
	}
	res := Response{
		Success: true,
		Data:    tutList,
	}
	err = utils.WriteJSON(w, http.StatusOK, res, nil)
	if err != nil {
		fmt.Printf("Could not encode to json %v", err)
		http.Error(w, "JSON encoding failed : "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (th *TutorHandler) GetTutors(w http.ResponseWriter, r *http.Request) {

	const pageSize = 5

	p := r.URL.Query().Get("page")
	if p == "" {
		p = "0"
	}

	pageNumber, err := strconv.Atoi(p)
	if err != nil {
		http.Error(w, `Could not get a list of jobs : `+err.Error(), http.StatusInternalServerError)
		fmt.Printf("Could not convert string to int %v", err)
		return
	}

	offset := (pageNumber - 1) * pageSize

	tutors, err := th.ts.GetTutors(offset, pageSize)

	if err != nil {
		http.Error(w, "Error getting tutors : "+err.Error(), http.StatusInternalServerError)
		return
	}

	var tuts []domain.Tutor

	for _, tutor := range tutors {
		tuts = append(tuts, *tutor)
	}

	res := Response{
		Success: true,
		Data:    tuts,
	}

	err = utils.WriteJSON(w, http.StatusOK, res, nil)
	if err != nil {
		fmt.Printf("Could not encode to json %v", err)
		http.Error(w, "JSON encoding failed : "+err.Error(), http.StatusInternalServerError)
		return
	}

}

func (th *TutorHandler) FilterTutors(w http.ResponseWriter, r *http.Request) {

	err := r.ParseMultipartForm(10 << 20)

	if err != nil {
		http.Error(w, "Error parsing form : "+err.Error(), http.StatusInternalServerError)
		return
	}

	gender := domain.Gender(r.PostForm.Get("gender"))
	city := r.PostForm.Get("city")
	education := domain.Education(r.PostForm.Get("education"))

	fmt.Println("education at filter handler is ", education)
	fieldOfStudy := r.PostForm.Get("fieldOfStudy")

	var ratingInt int

	if rating := r.PostForm.Get("rating"); rating != "" {
		ratingInt, err = strconv.Atoi(rating)
		if err != nil {
			http.Error(w, "Conversion of rating to integer failed"+err.Error(), http.StatusInternalServerError)
			return
		}
	}
	var hourlyMinInt int
	if hourlyMin := r.PostForm.Get("hourlyMin"); hourlyMin != "" {
		hourlyMinInt, err = strconv.Atoi(hourlyMin)
		if err != nil {
			http.Error(w, "Conversion of hourlyMin to integer failed"+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	var hourlyMaxInt int
	if hourlyMax := r.PostForm.Get("hourlyMax"); hourlyMax != "" {
		hourlyMaxInt, err = strconv.Atoi(hourlyMax)
		if err != nil {
			http.Error(w, "Conversion of hourlyMax to integer failed"+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	tutors, err := th.ts.FilterTutor(gender, ratingInt, hourlyMinInt, hourlyMaxInt, city, education, fieldOfStudy)

	if err != nil {
		http.Error(w, "Could not filter tutors : \n"+err.Error(), http.StatusInternalServerError)
		return
	}

	var data []domain.Tutor
	for _, tutor := range tutors {
		data = append(data, *tutor)
	}

	res := Response{
		Success: true,
		Data:    data,
	}

	err = utils.WriteJSON(w, http.StatusOK, res, nil)

	if err != nil {
		fmt.Printf("Could not encode to json %v", err)
		http.Error(w, "JSON encoding failed : "+err.Error(), http.StatusInternalServerError)
		return
	}

}
func (th *TutorHandler) RateTutor(w http.ResponseWriter, r *http.Request) {

	r.ParseMultipartForm(10 << 20)

	// ntfId := r.URL.Query().Get("ntfId")
	tutorId := r.URL.Query().Get("tutId")

	payload, err := utils.GetPayload(r)
	if err != nil {
		http.Error(w, "Coult not get payload : "+err.Error(), http.StatusInternalServerError)
		return
	}

	if payload["role"] != string(domain.ClientRole) {
		http.Error(w, "Only clients can rate tutors ", http.StatusForbidden)
		return
	}

	clientId := payload["id"]

	approved, err := th.ts.ApproveRating(clientId, tutorId)
	if err != nil {
		http.Error(w, "Could not check if client is allowed to rate : "+err.Error(), http.StatusInternalServerError)
		return
	}

	if !approved {
		http.Error(w, "Client has not worked with tutor to rate ", http.StatusForbidden)
		return
	}
	fmt.Println("tutor id at from url is ", tutorId)

	rating := r.PostForm.Get("rating")

	newRating, err := strconv.Atoi(rating)
	if err != nil {
		http.Error(w, "Could not convert hiring to integer : "+err.Error(), http.StatusInternalServerError)
		return
	}

	tutor, err := th.ts.GetTutorById(tutorId)

	if err != nil {
		http.Error(w, "Could not get tutor by id : "+err.Error(), http.StatusInternalServerError)
		return
	}

	if newRating < 1 || newRating > 5 {
		http.Error(w, "Invalid rating value", http.StatusBadRequest)
		return
	}

	var existingRatingCount float32
	if tutor.RateCount > 0 {
		existingRatingCount = float32(tutor.RateCount)
	} else {
		existingRatingCount = 1.0 // Handle no previous ratings
	}

	finalRating := (float32(tutor.Rating)*existingRatingCount + float32(newRating)) / (existingRatingCount + 1)

	updatedTutor := domain.Tutor{
		RateCount: int(existingRatingCount) + 1,
		Rating:    finalRating,
	}

	finalTutor, err := th.ts.UpdateTutor(updatedTutor, tutorId)

	if err != nil {
		http.Error(w, "Could not update rating of tutor : "+err.Error(), http.StatusInternalServerError)
		return
	}

	res := Response{
		Success: true,
		Data:    finalTutor,
	}
	err = utils.WriteJSON(w, http.StatusOK, res, nil)

	if err != nil {
		http.Error(w, "Could not marshal to json", http.StatusInternalServerError)
		return
	}

}
