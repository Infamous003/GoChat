package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *application) routes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// Custom error responses
	r.MethodNotAllowed(app.methodNotAllowedResponse)
	r.NotFound(app.notfoundResponse)

	r.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello world!"))
	})

	r.Post("/register", app.registerUserHandler)
	r.Post("/login", app.loginHandler)

	return r
}
