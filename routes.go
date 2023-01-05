package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/piotrselak/back/handlers"
)

func quizRouter() http.Handler {
	r := chi.NewRouter()
	//r.Use(AdminOnly)
	r.Get("/", handlers.FetchAllQuizes())
	return r
}
