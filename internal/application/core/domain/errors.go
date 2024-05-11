package domain

import (
	"errors"
)

var ErrNoRecord = errors.New("Models: No records were found")

var ErrInvalidCredentials = errors.New("Models: Invalid credentials")
