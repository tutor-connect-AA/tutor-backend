package utils

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPass(password string) (string, error) {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), 12) // cost
	if err != nil {
		return "", err
	}

	return string(hashedPass), nil
}

func CheckPass(hashedPassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err
}
