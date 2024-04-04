package domain

type Admin struct {
	firstName   string
	fathersName string
	username    string
	password    string
	role        Role // is this needed at all?
}
