package utils

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var key = []byte("secret-key")

func Tokenize(id string, role string) (string, error) {

	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = "subject"
	claims["id"] = id
	claims["role"] = role
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix()

	// Sign the token with the secret key
	tokenString, err := token.SignedString(key)
	if err != nil {
		fmt.Println("Error signing token:", err)
		return "", err
	}

	fmt.Println("JWT token:", tokenString)
	return tokenString, nil
}

func VerifyToken(tokenString string) (map[string]interface{}, error) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Check signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return key, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	claims := token.Claims.(jwt.MapClaims)
	return claims, nil

}

func GetPayload(r *http.Request) (map[string]string, error) {
	// token := r.Header.Get("Authorization")
	// token = token[len("Bearer "):]
	payload := make(map[string]string)
	token := r.Header.Get("Authorization")
	token = token[len("Bearer "):]

	claims, err := VerifyToken(token)
	if err != nil {
		fmt.Printf("Could not verify token\n %v ", err)
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
