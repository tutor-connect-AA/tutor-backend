package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var key, err = GenerateKey() //handle this better
// var key = []byte("secret-key")

func Tokenize(payload string) (string, error) {

	if err != nil {
		return "", nil
	}

	// fmt.Println("Tokenizing key is :", key)

	expiration := time.Now().Add(time.Second * 30)

	t := jwt.NewWithClaims(jwt.SigningMethodES256,
		jwt.MapClaims{
			"name": payload,
			"exp":  expiration.Unix(),
		})
	s, err := t.SignedString(key)
	if err != nil {
		fmt.Print("err1", err)
		return "", err
	}
	return s, nil
}

// look deeper into jwt token verification options for better security
func VerifyToken(tokenString string) error {
	// token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
	// 	return key, nil
	// })
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Check signing method
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// Extract public key from token
		return &key.PublicKey, nil
	})
	// fmt.Println("Verifying key is :", key)
	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}
