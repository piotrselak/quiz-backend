package main

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	db "github.com/piotrselak/back/db"
	http2 "github.com/piotrselak/back/http"
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
		_, _ = w.Write([]byte("api"))
	})

	r.Mount("/quiz", quizRouter(driver))

	err := http.ListenAndServe(":3333", r)
	if err != nil {
		panic(err)
	}
}

func quizRouter(driver neo4j.DriverWithContext) http.Handler {
	r := chi.NewRouter()
	r.Use(openSession(driver))
	r.Get("/", http2.FetchAllQuizes)
	r.Post("/", http2.CreateNewQuiz)
	return r
}

func openSession(driver neo4j.DriverWithContext) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session := driver.NewSession(r.Context(), neo4j.SessionConfig{})
			defer func() {
				err := session.Close(r.Context())
				if err != nil {
					return
				}
			}()

			ctx := context.WithValue(r.Context(), "session", session)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
