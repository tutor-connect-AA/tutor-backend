package main

import (
	"context"
	"net/http"
)

func FileUploadMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		file, header, err := r.FormFile("photo")
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		defer file.Close()
		ctx := context.WithValue(r.Context(), "filePath", header.Filename)
		ctx = context.WithValue(ctx, "file", file)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
