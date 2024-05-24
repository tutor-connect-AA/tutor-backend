package domain

import "time"

type Notification struct {
	Id        string
	Message   string
	OwnerId   string
	Opened    bool
	CreatedAt time.Time
}
