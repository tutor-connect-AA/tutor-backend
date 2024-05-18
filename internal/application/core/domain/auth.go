package domain

type Auth struct {
	Id       string
	Username string
	Password string
	// ClientID string
	// TutorID  string
	Role Role
}
