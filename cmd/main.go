package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"net/http"

	"github.com/joho/godotenv"
	"github.com/justinas/alice"
	"github.com/tutor-connect-AA/tutor-backend/internal/adapters/db"
	"github.com/tutor-connect-AA/tutor-backend/internal/adapters/handlers"
	"github.com/tutor-connect-AA/tutor-backend/internal/application/core/api"
)

func init() {
	// Load environment variables from .env file if it exists
	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(".env"); err != nil {
			log.Fatal("Error loading .env file")
		}
	}
}

func enableCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {

	//Login endpoint should be shared by tutor and client?

	// err := godotenv.Load(".env")

	// if err != nil {
	// 	log.Fatal("Could not access dotenv variables")
	// 	log.Print(err)
	// 	return
	// }

	dbString := os.Getenv("DB_URL")
	fmt.Println("DB STRING : ", dbString)

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	fmt.Println("db string is :", dbString)

	var dsn = flag.String("dsn", dbString, "Connection string to database")
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

	hireH := handlers.NewHiringHandler(jaSer, clientSer, tutSer, tNtfSer, cNtfSer, jobSer)

	//JobRequest Application configuration
	jrRepo := db.NewJobRequestRepo(dbConnection)
	jrSer := api.NewJobRequestAPI(jrRepo)
	jrHandler := handlers.NewJobRequestHandler(jrSer, clientSer, tutSer)

	mux := http.NewServeMux()

	// handler := enableCors(mux)

	protected := alice.New(AuthMiddleware)

	fileUpload := alice.New(FileUploadMiddleware)

	mux.HandleFunc("/clients/register", clientHandler.Register)
	mux.HandleFunc("/clients", clientHandler.GetListOfClients)
	mux.HandleFunc("/clients/single", clientHandler.GetClientById) //make path make sense
	mux.Handle("/clients/update", protected.ThenFunc(clientHandler.UpdateClientProfile))
	// mux.HandleFunc("/client/login", clientHandler.LoginClient)

	mux.Handle("/jobs/post", protected.ThenFunc(jobHandler.PostJob))
	mux.HandleFunc("/jobs/single", jobHandler.GetJobById)
	mux.HandleFunc("/jobs", jobHandler.GetJobs)
	mux.HandleFunc("/jobs/Interview", jobHandler.ComposeInterview)

	mux.Handle("/tutor/register", fileUpload.ThenFunc(tutHandler.RegisterTutor))
	// mux.HandleFunc("/tutor/login", tutHandler.LoginTutor)

	mux.HandleFunc("/job-applications/create", jaHandler.Apply)
	mux.HandleFunc("/job-applications/single", jaHandler.GetApplicationById)
	mux.HandleFunc("/job-applications/jobs", jaHandler.ApplicationsByJob)
	mux.HandleFunc("/job-applications/tutors", jaHandler.ApplicationsByTutor)
	mux.HandleFunc("/job-applications/clients", jaHandler.ApplicationsByClient)
	mux.HandleFunc("/job-applications/status", jaHandler.GetApplicationByStatus)

	mux.HandleFunc("/hiring/hire", hireH.Hire)
	mux.HandleFunc("/hiring/verify-hire", hireH.VerifyHire)
	mux.HandleFunc("/hiring/shortlist", hireH.Shortlist)
	mux.Handle("/hiring/reply", fileUpload.ThenFunc(hireH.SendInterview))

	mux.HandleFunc("/login", authHandler.Login)

	mux.Handle("/job-request/new", protected.ThenFunc(jrHandler.RequestJob))
	mux.HandleFunc("/job-request/single", jrHandler.GetJobRequest)
	mux.Handle("/job-request/update", protected.ThenFunc(jrHandler.ChangeJobRequestStatus))
	mux.Handle("/job-request/hire", protected.ThenFunc(jrHandler.HireFromRequest))
	mux.HandleFunc("/job-request/verify-hire", jrHandler.VerifyHireFromRequest)
	// mux.HandleFunc("/jobRequest/multiple",jrHandler.)

	mux.HandleFunc("/tutor-notification/single", tNfHandler.GetTutorNotification)
	mux.Handle("/tutor-notifications", protected.ThenFunc(tNfHandler.GetTutorNotifications))
	mux.Handle("/tutor-notifications/unopened", protected.ThenFunc(tNfHandler.UnopenedTutorNtfs))
	mux.Handle("/tutor-notifications/count", protected.ThenFunc(tNfHandler.CountUnopenedTutorNtfs))

	mux.HandleFunc("/client-notification/single", cNfHandler.GetClientNotification)
	mux.Handle("/client-notifications", protected.ThenFunc(cNfHandler.GetClientNotifications))
	mux.Handle("/client-notifications/unopened", protected.ThenFunc(cNfHandler.UnopenedClientNtfs))
	mux.Handle("/client-notifications/count", protected.ThenFunc(cNfHandler.CountUnopenedClientNtfs))

	log.Println("Listening on port:", port)
	http.ListenAndServe(":"+port, enableCors(mux))

}
