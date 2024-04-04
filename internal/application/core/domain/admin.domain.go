package domain

type Admin struct {
	id          string
	firstName   string
	fathersName string
	username    string
	password    string
	role        Role // is this needed at all?
}
