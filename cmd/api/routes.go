package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)


func (app *application) routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(app.enableCORS)
	mux.Get("/", app.Home)

	// routes for the API calls
    mux.Route("/api", func (r chi.Router)  {
        r.Get("/posts", app.GetPosts)
		r.Get("/tags", app.GetTags)
		r.Post("/login", app.Login) 
    })

	return mux
}