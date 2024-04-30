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
			http.Error(w, "Could not get tutor's cv", http.StatusInternalServerError)
			return
		}
		defer cv.Close()
		ctx = context.WithValue(ctx, "cvPath", cvHeader.Filename)
		ctx = context.WithValue(ctx, "cv", cv)

		eduCred, eduCredPath, err := r.FormFile("eduCred")
		if err != nil {
			http.Error(w, "Could not get education credential file", http.StatusInternalServerError)
			return
		}
		defer eduCred.Close()
		ctx = context.WithValue(ctx, "eduCredPath", eduCredPath.Filename)
		ctx = context.WithValue(ctx, "eduCred", eduCred)

		// fmt.Printf("The uploaded cv header %v", cv)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
