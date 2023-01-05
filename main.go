package main

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/piotrselak/back/internal/handlers"
	"github.com/piotrselak/back/pkg/db"
)

func main() {
	driver := db.InitNeo4j()
	ctx := context.Background()
	defer driver.Close(ctx)

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("api"))
	})

	r.Mount("/quiz", quizRouter(driver))

	http.ListenAndServe(":3333", r)
}

func quizRouter(driver neo4j.DriverWithContext) http.Handler {
	r := chi.NewRouter()
	r.Use(db.OpenSession(driver))
	r.Get("/", handlers.FetchAllQuizes)
	return r
}
