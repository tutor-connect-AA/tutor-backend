package domain

type Admin struct {
	Id          string
	FirstName   string
	FathersName string
	Username    string
	Password    string
	Role        Role // is this needed at all?
}
