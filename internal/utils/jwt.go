package utils

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var key, err = GenerateKey() //handle this better
// var key = []byte("secret-key")

func Tokenize(id string) (string, error) {

	if err != nil {
		return "", nil
	}

	// fmt.Println("Tokenizing key is :", key)

	expiration := time.Now().Add(time.Hour * 24)

	t := jwt.NewWithClaims(jwt.SigningMethodES256,
		jwt.MapClaims{
			"id":  id,
			"exp": expiration.Unix(),
		})
	s, err := t.SignedString(key)
	if err != nil {
		fmt.Print("err1", err)
		return "", err
	}
	return s, nil
}

// look deeper into jwt token verification options for better security
func VerifyToken(tokenString string) (map[string]interface{}, error) {
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
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	claims := token.Claims.(jwt.MapClaims)
	return claims, nil
	// return nil
}

func GetPayload(r *http.Request) (map[string]string, error) {
	// token := r.Header.Get("Authorization")
	// token = token[len("Bearer "):]
	payload := make(map[string]string)
	token := r.Header.Get("Authorization")
	token = token[len("Bearer "):]

	if err != nil {
		return payload, err
	}

	claims, err := VerifyToken(token)
	if err != nil {
		fmt.Printf("Could not verify token%v", err)
		return payload, err
	}
	id, ok := claims["id"].(string)
	fmt.Println("id in payload is ", id)
	if !ok {
		fmt.Println("Could not find user ID in claims")
		return payload, err
	}
	payload["id"] = id
	return payload, nil
}
