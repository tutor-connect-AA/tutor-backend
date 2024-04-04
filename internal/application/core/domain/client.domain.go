package domain

type Client struct {
	id          string
	firstName   string
	fathersName string //optional
	phoneNumber string
	email       string
	photo       string
	role        Role // should role even exist?
	rating      int
}
