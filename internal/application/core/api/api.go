package api

import (
	"github.com/tutor-connect-AA/tutor-backend/internal/ports/db_ports"
)

type Application struct {
	db db_ports.DBPort
}

func NewApplication(db db_ports.DBPort) *Application {
	return &Application{
		db: db,
	}
}
