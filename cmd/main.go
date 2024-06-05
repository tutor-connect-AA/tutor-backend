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

	userRepo := db.NewUserRepo(dbConnection)

	//Tutor Notification
	tNtfRepo := db.NewTutorNotificationRepo(dbConnection)
	tNtfSer := api.NewTutorNotificationAPI(tNtfRepo)
	tNfHandler := handlers.NewTutorNotificationHandler(tNtfSer)

	//Client Notification
	cNtfRepo := db.NewClientNotificationRepo(dbConnection)
	cNtfSer := api.NewClientNotificationAPI(cNtfRepo)
	cNfHandler := handlers.NewClientNotificationHandler(cNtfSer)

	//Client configuration
	// clientRepo := db.NewClientRepo(dbConnection)
	clientSer := api.NewClientAPI(userRepo)
	clientHandler := handlers.NewClientHandler(clientSer)

	//Job configuration
	jobRepo := db.NewJobRepo(dbConnection)
	jobSer := api.NewJobAPI(jobRepo)
	jobHandler := handlers.NewJobHandler(jobSer)

	//tutor configuration
	// tutRepo := db.NewTutorRepo(dbConnection)
	tutSer := api.NewTutorAPI(userRepo)
	tutHandler := handlers.NewTutorHandler(tutSer)

	//Job Application configuration
	jaRepo := db.NewJobApplicationRepo(dbConnection)
	jaSer := api.NewJobApplicationAPI(jaRepo)
	jaHandler := handlers.NewJobApplicationHandler(jaSer, tutSer, cNtfSer, jobSer)

	//Auth handler config(client & tutor)
	authSer := api.NewAuthService(userRepo)
	authHandler := handlers.NewAuthHandler(authSer, clientSer, tutSer)

	hireH := handlers.NewHiringHandler(jaSer, clientSer)

	//JobRequest Application configuration
	jrRepo := db.NewJobRequestRepo(dbConnection)
	jrSer := api.NewJobRequestAPI(jrRepo)
	jrHandler := handlers.NewJobRequestHandler(jrSer)

	mux := http.NewServeMux()

	protected := alice.New(AuthMiddleware)

	fileUpload := alice.New(FileUploadMiddleware)

	mux.HandleFunc("/client/register", clientHandler.Register)
	mux.HandleFunc("/client/listClients", clientHandler.GetListOfClients)
	mux.HandleFunc("/client/single", clientHandler.GetClientById) //make path make sense
	mux.Handle("/client/update", protected.ThenFunc(clientHandler.UpdateClientProfile))
	// mux.HandleFunc("/client/login", clientHandler.LoginClient)

	mux.Handle("/job/post", protected.ThenFunc(jobHandler.PostJob))
	mux.HandleFunc("/job/single", jobHandler.GetJobById)
	mux.HandleFunc("/job/all", jobHandler.GetJobs)

	mux.Handle("/tutor/register", fileUpload.ThenFunc(tutHandler.RegisterTutor))
	// mux.HandleFunc("/tutor/login", tutHandler.LoginTutor)

	mux.HandleFunc("/jobApplication/newJob", jaHandler.Apply)
	mux.HandleFunc("/jobApplication/single", jaHandler.GetApplicationById)
	mux.HandleFunc("/jobApplication/job", jaHandler.ApplicationsByJob)
	mux.HandleFunc("/jobApplication/tutor", jaHandler.ApplicationsByJob)
	mux.HandleFunc("/jobApplication/client", jaHandler.ApplicationsByClient)

	mux.HandleFunc("/jobApplication/hire", hireH.Hire)
	mux.HandleFunc("/jobApplication/verifyHire", hireH.VerifyHire)

	mux.HandleFunc("/login", authHandler.Login)

	mux.Handle("/jobRequest/new", protected.ThenFunc(jrHandler.RequestJob))
	mux.HandleFunc("/jobRequest/single", jrHandler.GetJobRequest)
	// mux.HandleFunc("/jobRequest/multiple",jrHandler.)

	mux.HandleFunc("/tutorNotification/single", tNfHandler.GetTutorNotification)
	mux.HandleFunc("/tutorNotifications", tNfHandler.GetTutorNotifications)

	mux.HandleFunc("/clientNotification/single", cNfHandler.GetClientNotification)
	mux.HandleFunc("/clientNotifications", cNfHandler.GetClientNotifications)

	log.Println("Listening on port 8080")
	http.ListenAndServe(":8080", mux)

}
