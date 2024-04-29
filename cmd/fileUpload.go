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
		// fmt.Printf("The uploaded photo header: %v", photo)

		cv, cvHeader, err := r.FormFile("cv")
		if err != nil {
			http.Error(w, "Could not get education credential files", http.StatusInternalServerError)
			return
		}
		defer cv.Close()
		ctx = context.WithValue(ctx, "cvPath", cvHeader.Filename)
		ctx = context.WithValue(ctx, "cv", cv)
		// fmt.Printf("The uploaded cv header %v", cv)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
