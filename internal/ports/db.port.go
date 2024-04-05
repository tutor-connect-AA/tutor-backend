package port

import (
	"github.com/tutor-connect-AA/tutor-backend/internal/application/core/domain"
)

type ClientPort interface {
	Get(id string) (domain.Client, error)
}
