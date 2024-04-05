package domain

type Client struct {
	Id          string
	FirstName   string
	FathersName string //optional
	PhoneNumber string
	Email       string
	Photo       string
	Role        Role // should role even exist?
	Rating      float32
}
