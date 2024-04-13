package utils

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

var key, err = GenerateKey() //handle this better

func Tokenize(payload string) (string, error) {

	if err != nil {
		return "", nil
	}
	t := jwt.NewWithClaims(jwt.SigningMethodES256,
		jwt.MapClaims{
			"name": payload,
		})
	s, err := t.SignedString(key)
	if err != nil {
		return "", err
	}
	return s, nil
}

// look deeper into jwt token verification options for better security
func VerifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}
