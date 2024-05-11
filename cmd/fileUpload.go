package main

import (
	"context"
	"net/http"
)

func FileUploadMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var ctx context.Context

		// Check if the photo file was uploaded
		photo, photoHeader, err := r.FormFile("photo")
		if err == nil && photo != nil {
			defer photo.Close()
			ctx = context.WithValue(r.Context(), "photoPath", photoHeader.Filename)
			ctx = context.WithValue(ctx, "photo", photo)
		}

		// Check if the cv file was uploaded
		cv, cvHeader, err := r.FormFile("cv")
		if err == nil && cv != nil {
			defer cv.Close()
			ctx = context.WithValue(r.Context(), "cvPath", cvHeader.Filename)
			ctx = context.WithValue(ctx, "cv", cv)
		}

		// Check if the education credential file was uploaded
		eduCred, eduCredHeader, err := r.FormFile("eduCred")
		if err == nil && eduCred != nil {
			defer eduCred.Close()
			ctx = context.WithValue(r.Context(), "eduCredPath", eduCredHeader.Filename)
			ctx = context.WithValue(ctx, "eduCred", eduCred)
		}

		// Call the next handler with the updated context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
