package handlers

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/tutor-connect-AA/tutor-backend/internal/application/core/domain"
	"github.com/tutor-connect-AA/tutor-backend/internal/ports/api_ports"
	"github.com/tutor-connect-AA/tutor-backend/internal/utils"
)

type JobHandler struct {
	jobSer api_ports.JobAPIPort
}

func NewJobHandler(js api_ports.JobAPIPort) *JobHandler {
	return &JobHandler{
		jobSer: js,
	}
}

type CreateJobReq struct {
}

func (jobH JobHandler) PostJob(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusUnauthorized)
		return
	}

	err := r.ParseMultipartForm(10 << 20) // 10 MB max size
	if err != nil {
		http.Error(w, "Could not parse form : "+err.Error(), http.StatusBadRequest)
		return
	}
	payload, err := utils.GetPayload(r)

	if err != nil {
		http.Error(w, "Could not get payload form token ", http.StatusInternalServerError)
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
		http.Error(w, "Could not parse deadline : "+err.Error(), http.StatusInternalServerError)
		return
	}
	qt, err := strconv.Atoi(r.PostForm.Get("quantity"))
	if err != nil {
		fmt.Printf("Could not convert string to int %v", err)
		http.Error(w, "Could not convert string to int for quantity: "+err.Error(), http.StatusInternalServerError)
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
		http.Error(w, "Could not convert string to int for min: "+err.Error(), http.StatusInternalServerError)
		return
	}

	max := r.PostForm.Get("max")
	ma, err := strconv.Atoi(max)
	if err != nil {
		fmt.Printf("Could not convert string to int %v", err)
		http.Error(w, "Could not convert string to int for max : "+err.Error(), http.StatusInternalServerError)
		return
	}
	ch := r.PostForm.Get("contact")
	chr, err := strconv.Atoi(ch)
	if err != nil {
		fmt.Printf("Could not convert string to int %v", err)
		http.Error(w, "Could not convert string to int for contact : "+err.Error(), http.StatusInternalServerError)
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
	jb, err := jobH.jobSer.CreateJobPost(nj)
	if err != nil {
		http.Error(w, "Could not post job ", http.StatusInternalServerError)
		return
	}
	res := Response{
		Success: true,
		Data:    jb,
	}

	err = utils.WriteJSON(w, http.StatusOK, res, nil)
	if err != nil {
		fmt.Printf("Could not encode to json %v", err)
		http.Error(w, "JSON encoding failed : "+err.Error(), http.StatusInternalServerError)
		return
	}
	// fmt.Fprintf(w, "Created job : %v", jb)
}

func (jobH JobHandler) GetJobById(w http.ResponseWriter, r *http.Request) {
	jobId := r.URL.Query().Get("id")

	jb, err := jobH.jobSer.GetJob(jobId)

	if err != nil {
		fmt.Printf("Error at get jobById handler %v", err)
		http.Error(w, "Could not get a job by id : "+err.Error(), http.StatusInternalServerError)
		return
	}
	res := Response{
		Success: true,
		Data:    jb,
	}
	err = utils.WriteJSON(w, http.StatusOK, res, nil)
	if err != nil {
		fmt.Printf("Could not get a job by id %v", err)
		http.Error(w, "Could not get a job by id : "+err.Error(), http.StatusInternalServerError)
		return
	}
	// fmt.Fprintf(w, "Got job successfully %v", jb)
}

// Use offset pagination
// limit of 2 for now at least
// first page:
// offset == 0 and limit
func (jobH JobHandler) GetJobs(w http.ResponseWriter, r *http.Request) {

	const pageSize = 4

	p := r.URL.Query().Get("page")
	if p == "" {
		p = "1"
	}

	pageNumber, err := strconv.Atoi(p)
	if err != nil {
		http.Error(w, `Could not get a list of jobs : `+err.Error(), http.StatusInternalServerError)
		fmt.Printf("Could not convert string to int %v", err)
		return
	}

	offset := (pageNumber - 1) * pageSize
	jbs, err := jobH.jobSer.GetListOfJobs(offset, pageSize)

	jbList := []domain.Job{}

	for _, jb := range jbs {
		jbList = append(jbList, *jb)
	}

	if err != nil {
		fmt.Printf("Error at getting list of jobs %v", err)
		http.Error(w, "Could not get a list of  : "+err.Error(), http.StatusInternalServerError)
		return
	}

	count, err := jobH.jobSer.GetJobCount()
	if err != nil {
		http.Error(w, "Could not get count of jobs by clients : "+err.Error(), http.StatusInternalServerError)
		return
	}

	numberOfPages := math.Ceil(float64(float32(*count) / float32(pageSize)))

	res := map[string]interface{}{
		"Success":       true,
		"Data":          jbList,
		"numberOfPages": numberOfPages,
	}

	err = utils.WriteJSON(w, http.StatusOK, res, nil)
	if err != nil {
		fmt.Printf("Could not get a list of jobs %v", err)
		http.Error(w, "Could not get a list of jobs : "+err.Error(), http.StatusInternalServerError)
		return
	}

	// fmt.Fprintf(w, "Got list of jobs successfully %v", jbs)
	// fmt.Printf("Here are the list of jobs %v", jbs)
}

func (jobH JobHandler) ComposeInterview(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Error parsing form : "+err.Error(), http.StatusInternalServerError)
		return
	}

	jobId := r.URL.Query().Get("jobId")

	questions := r.PostForm["questions"]

	var stringOfQuestions string

	for _, question := range questions {
		if stringOfQuestions == "" {
			stringOfQuestions += question
		} else {
			stringOfQuestions += "~" + question
		}
	}
	updatedJob := domain.Job{
		Interview_Questions: stringOfQuestions,
	}
	job, err := jobH.jobSer.UpdateJob(jobId, updatedJob)

	if err != nil {
		http.Error(w, "Could not compose interview : "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := Response{
		Success: true,
		Data:    job,
	}

	err = utils.WriteJSON(w, http.StatusOK, response, nil)

	if err != nil {
		http.Error(w, "Error marshalling to json  : "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (jobH JobHandler) GetJobByClient(w http.ResponseWriter, r *http.Request) {
	const pageSize = 4

	clientId := r.URL.Query().Get("cltId")
	page := r.URL.Query().Get("page")

	if page == "" {
		page = "1"
	}

	if clientId == "" {
		http.Error(w, "Client id can not be empty", http.StatusBadRequest)
		return
	}

	pageNumber, err := strconv.Atoi(page)

	if err != nil {
		http.Error(w, "Error converting to integer", http.StatusInternalServerError)
		return
	}

	offset := (pageNumber - 1) * pageSize

	jobs, err := jobH.jobSer.GetJobByClient(clientId, offset, pageSize)

	if err != nil {
		http.Error(w, "Could not get jobs by client : "+err.Error(), http.StatusInternalServerError)
		return
	}

	count, err := jobH.jobSer.GetJobCountByClient(clientId)
	if err != nil {
		http.Error(w, "Could not get count of jobs by clients : "+err.Error(), http.StatusInternalServerError)
		return
	}

	numberOfPages := math.Ceil(float64(float32(*count) / float32(pageSize)))

	res := map[string]interface{}{
		"Success":       true,
		"Data":          jobs,
		"numberOfPages": numberOfPages,
	}

	err = utils.WriteJSON(w, http.StatusOK, res, nil)

	if err != nil {
		http.Error(w, "Could not serialize to json", http.StatusInternalServerError)
		return
	}
}
