package handlers

import (
	"github.com/tutor-connect-AA/tutor-backend/internal/ports/api_ports"
)

type Adapter struct {
	ser api_ports.APIPort
}
