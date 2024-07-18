package utils

import (
	"fmt"
	"net/smtp"
)

func SendEmail(to []string, message []byte) {
	from := "mahider3991@gmail.com"
	password := "Maverick2020!"

	// Receiver email address.
	// to := []string{
	//   "sender@example.com",
	// }

	// to := reveiver

	// smtp server configuration.
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Sending email.
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Email Sent Successfully!")
}
