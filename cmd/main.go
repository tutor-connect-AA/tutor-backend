package main

import (
	"flag"
	"log"

	"net/http"

	"github.com/justinas/alice"
	"github.com/tutor-connect-AA/tutor-backend/internal/adapters/db"
	"github.com/tutor-connect-AA/tutor-backend/internal/adapters/handlers"
	"github.com/tutor-connect-AA/tutor-backend/internal/application/core/api"
)

func main() {

	//Login endpoint should be shared by tutor and client?

	var dsn = flag.String("dsn", "postgres://postgres:Maverick2020!@localhost:5432/tutor-connect", "Connection string to database")
	dbConnection, err := db.ConnectDB(*dsn)

	if err != nil {
		log.Fatal("Can't connect to database")
		return
	}

	//Client configuration
	clientRepo := db.NewClientRepo(dbConnection)
	clientAPI := api.NewClientAPI(clientRepo)
	clientHandler := handlers.NewClientHandler(clientAPI)

	//Job configuration
	jobRepo := db.NewJobRepo(dbConnection)
	jobAPI := api.NewJobAPI(jobRepo)
	jobHandler := handlers.NewJobHandler(jobAPI)

	//tutor configuration
	tutRepo := db.NewTutorRepo(dbConnection)
	tutAPI := api.NewTutorAPI(tutRepo)
	tutHandler := handlers.NewTutorHandler(tutAPI)

	//Job Application configuration
	jaRepo := db.NewJobApplicationRepo(dbConnection)
	jaAPI := api.NewJobApplicationAPI(jaRepo)
	jaHandler := handlers.NewJobApplicationHandler(jaAPI)

	mux := http.NewServeMux()

	protected := alice.New(AuthMiddleware)

	fileUpload := alice.New(FileUploadMiddleware)

	mux.Handle("/client/register", fileUpload.ThenFunc(clientHandler.Register))
	mux.HandleFunc("/client/listClients", clientHandler.GetListOfClients)
	mux.HandleFunc("/client/single", clientHandler.GetClientById) //make path make sense
	mux.Handle("/client/update", protected.ThenFunc(clientHandler.UpdateClientProfile))
	mux.HandleFunc("/client/login", clientHandler.LoginClient)

	mux.Handle("/job/post", protected.ThenFunc(jobHandler.PostJob))
	mux.HandleFunc("/job/single", jobHandler.GetJobById)
	mux.HandleFunc("/job/all", jobHandler.GetJobs)

	mux.Handle("/tutor/register", fileUpload.ThenFunc(tutHandler.RegisterTutor))
	mux.HandleFunc("/tutor/login", tutHandler.LoginTutor)

	mux.HandleFunc("/jobApplication/newJob", jaHandler.Apply)
	mux.HandleFunc("/jobApplication/job", jaHandler.ApplicationsByJob)
	mux.HandleFunc("/jobApplication/tutor", jaHandler.ApplicationsByJob)
	mux.HandleFunc("/jobApplication/client", jaHandler.ApplicationsByClient)

	log.Println("Listening on port 8080")
	http.ListenAndServe(":8080", mux)

}
