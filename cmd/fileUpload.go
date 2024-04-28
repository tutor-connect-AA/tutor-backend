package main

import (
	"context"
	"net/http"
)

func FileUploadMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		photo, photoHeader, err := r.FormFile("photo")
		if err != nil {
			http.Error(w, "Could not get photo file from request", http.StatusInternalServerError)
			return
		}
		defer photo.Close()
		ctx := context.WithValue(r.Context(), "photoPath", photoHeader.Filename)
		ctx = context.WithValue(ctx, "photo", photo)

		edu, eduHeader, err := r.FormFile("cv")
		if err != nil {
			http.Error(w, "Could not get education credential files", http.StatusInternalServerError)
			return
		}
		defer photo.Close()
		ctx = context.WithValue(r.Context(), "cvPath", eduHeader.Filename)
		ctx = context.WithValue(ctx, "cv", edu)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
