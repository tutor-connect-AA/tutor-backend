package domain

import (
	"errors"
)

var ErrNoRecord = errors.New("models: No records were found")

var ErrInvalidCredentials = errors.New("models: Invalid credentials")
