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
    mux.Route("/api", func(r chi.Router) {
        // Post routes
        r.Get("/posts", app.GetPosts)
        r.Post("/posts", app.CreatePost)

        // Tag routes
        r.Get("/tags", app.GetTags)
        r.Post("/tags", app.CreateTag)
        r.Patch("/tags/{id}/increment-search", app.UpdateTagSearchCount) // New endpoint for incrementing search count

        // Authentication routes
        r.Post("/login", app.Login)
    })

    return mux
}