package main

import "net/http"

// CORS Cross-Origin Resource Sharing is a security feature that controls which web origins can access API
func (app *application) enableCORS(h http.Handler) http.Handler {
	// type in GO that adapts regular functions to serve HTTP requests
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// sets headers allowing requests from front-end localhost
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")

		// gets the HTTP method from the request (browser makes cross-origin requests first)
		if r.Method == "OPTIONS" {
			// permits credential sharing
			w.Header().Set("Access-Control-Allow-Crendentials", "true")
			// lists the allowed HTTP methods
			w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
			// specifies the permitted request headers
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, X-CSRF-Token, Authorization")
			return
		} else {
			h.ServeHTTP(w, r)
		}
	})
}