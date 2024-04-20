package main

import (
	"fmt"
	"net/http"

	"github.com/tutor-connect-AA/tutor-backend/internal/utils"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		// fmt.Println(token)
		if token == "" {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Missing authorization header")
			return
		}
		token = token[len("Bearer "):]
		_, err := utils.VerifyToken(token)

		if err != nil {
			http.Error(w, "Not authorized", http.StatusUnauthorized)
			fmt.Print("Invalid token")
			return
		}
		fmt.Fprintf(w, "You are authorized")
		next.ServeHTTP(w, r)
	})
}
