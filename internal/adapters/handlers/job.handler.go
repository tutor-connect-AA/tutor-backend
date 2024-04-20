package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/tutor-connect-AA/tutor-backend/internal/application/core/domain"
	"github.com/tutor-connect-AA/tutor-backend/internal/ports/api_ports"
	"github.com/tutor-connect-AA/tutor-backend/internal/utils"
)

type JobAdapter struct {
	jobSer api_ports.JobAPIPort
}

func NewJobHandler(js api_ports.JobAPIPort) *JobAdapter {
	return &JobAdapter{
		jobSer: js,
	}
}

type CreateJobReq struct {
}

func (adp JobAdapter) PostJob(w http.ResponseWriter, r *http.Request) {

	err := r.ParseMultipartForm(10 << 20) // 10 MB max size
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	token := r.Header.Get("Authorization")
	token = token[len("Bearer "):]

	claims, err := utils.VerifyToken(token)
	if err != nil {
		fmt.Printf("Could not get client id from token to post a job %v", err)
		http.Error(w, "Could not post job", http.StatusInternalServerError)
		return
	}
	clientID, ok := claims["id"].(string)
	fmt.Println("Client id is ", clientID)
	if !ok {
		fmt.Println("Could not find user ID in claims")
		return
	}
	dl, err := time.Parse("2006-01-02", r.PostForm.Get("deadline"))
	if err != nil {
		fmt.Printf("Could not parse deadline %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	qt, err := strconv.Atoi(r.PostForm.Get("quantity"))
	if err != nil {
		fmt.Printf("Could not convert string to int %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	grades := r.PostForm["grades"]
	// var intGrades []int
	var strGrades string
	for _, grade := range grades {
		if strGrades == "" {
			strGrades += grade
		} else {
			strGrades += "," + grade
		}
		// intValue, err := strconv.Atoi(grade)
		// if err != nil {
		// 	fmt.Printf("Error converting value to integer: %v", err)
		// 	fmt.Fprintf(w, "Error converting value to integer: %v", err)
		// 	continue
		// }
		// intGrades = append(intGrades, intValue)

	}

	min := r.PostForm.Get("min")
	mi, err := strconv.Atoi(min)
	if err != nil {
		fmt.Printf("Could not convert string to int %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	max := r.PostForm.Get("max")
	ma, err := strconv.Atoi(max)
	if err != nil {
		fmt.Printf("Could not convert string to int %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	ch := r.PostForm.Get("contact")
	chr, err := strconv.Atoi(ch)
	if err != nil {
		fmt.Printf("Could not convert string to int %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var nj = domain.Job{
		Title:                 r.PostForm.Get("title"),
		Description:           r.PostForm.Get("description"),
		Posted_By:             clientID,
		Deadline:              dl,
		Region:                r.PostForm.Get("region"),
		City:                  r.PostForm.Get("city"),
		Quantity:              qt,
		Grade_Of_Students:     strGrades,
		Minimum_Education:     domain.Education(r.PostForm.Get("education")),
		Preferred_Gender:      domain.Gender(r.PostForm.Get("gender")),
		Contact_Hour_Per_Week: chr,
		Hourly_Rate_Min:       mi,
		Hourly_Rate_Max:       ma,
	}
	jb, err := adp.jobSer.CreateJobPost(nj)
	if err != nil {
		http.Error(w, "Could not post job ", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Created job : %v", jb)
}
