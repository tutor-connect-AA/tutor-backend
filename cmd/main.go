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

	var dsn = flag.String("dsn", "postgres://postgres:Maverick2020!@localhost:5432/tutor-connect", "Connection string to database")
	dbAdapter, err := db.NewAdapter(*dsn)

	if err != nil {
		log.Fatal("Can't connect to database")
		return
	}

	app := api.NewApplication(dbAdapter)

	handlers := handlers.NewHandler(app)

	mux := http.NewServeMux()

	protected := alice.New(AuthMiddleware)

	fileUpload := alice.New(FileUploadMiddleware)

	mux.Handle("/client/register", fileUpload.ThenFunc(handlers.Register))
	mux.HandleFunc("/client/listClients", handlers.GetListOfClients)
	mux.HandleFunc("/client/single", handlers.GetClientById) //make path make sense
	mux.Handle("/client/update", protected.ThenFunc(handlers.UpdateClientProfile))
	mux.HandleFunc("/client/login", handlers.LoginClient)

	log.Println("Listening on port 8080")
	http.ListenAndServe(":8080", mux)

}
