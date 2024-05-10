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

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusUnauthorized)
		return
	}

	err := r.ParseMultipartForm(10 << 20) // 10 MB max size
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	payload, err := utils.GetPayload(r)

	if err != nil {
		http.Error(w, "Could not get payload form token", http.StatusInternalServerError)
		return
	}

	// token := r.Header.Get("Authorization")
	// token = token[len("Bearer "):]

	// if err != nil {
	// 	http.Error(w, "Could not get payload form token", http.StatusInternalServerError)
	// 	return
	// }

	// claims, err := utils.VerifyToken(token)
	// if err != nil {
	// 	fmt.Printf("Could not get client id from token to post a job %v", err)
	// 	http.Error(w, "Could not post job", http.StatusInternalServerError)
	// 	return
	// }
	// clientID, ok := claims["id"].(string)
	// fmt.Println("Client id is ", clientID)
	// if !ok {
	// 	fmt.Println("Could not find user ID in claims")
	// 	return
	// }
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
		Posted_By:             payload["id"],
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

func (adp JobAdapter) GetJobById(w http.ResponseWriter, r *http.Request) {
	jobId := r.URL.Query().Get("id")

	jb, err := adp.jobSer.GetJob(jobId)

	if err != nil {
		fmt.Printf("Error at get jobById handler %v", err)
		http.Error(w, "Could not get a job by id", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Got job successfully %v", jb)
}

// Use offset pagination
// limit of 2 for now at least
// first page:
// offset == 0 and limit
func (adp JobAdapter) GetJobs(w http.ResponseWriter, r *http.Request) {

	p := r.URL.Query().Get("page")
	if p == "" {
		p = "0"
	}
	pageNumber, err := strconv.Atoi(p)
	if err != nil {
		http.Error(w, `Could not get a list of jobs`, http.StatusInternalServerError)
		fmt.Printf("Could not convert string to int %v", err)
		return
	}
	const pageSize = 2
	offset := (pageNumber - 1) * pageSize
	jbs, err := adp.jobSer.GetListOfJobs(offset, pageSize)

	if err != nil {
		fmt.Printf("Error at getting list of jobs %v", err)
		http.Error(w, "Could not get a list of jobs", http.StatusInternalServerError)
		return
	}

	for _, job := range jbs {
		fmt.Fprint(w, *job)
	}

	// fmt.Fprintf(w, "Got list of jobs successfully %v", jbs)
	fmt.Printf("Here are the list of jobs %v", jbs)
}
